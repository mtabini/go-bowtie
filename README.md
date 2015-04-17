# Bowtie - Web middleware for Go

Bowtie is an HTTP middleware for Go. It makes heavy use of interfaces to provide a programming model that is both idiomatic and easily extensible (with, hopefully, relatively minimal overhead).

## Getting started

Standing up a basic Bowtie server only requires a few lines of code:

```go
package main

import (
    "github.com/mtabini/go-bowtie"
    "net/http"
)

func main() {
    s := bowtie.NewServer()

    http.ListenAndServe(":8000", s)
}
```

Like all comparable systems, however, Bowtie becomes powerful once you start adding middlewares to it.

## Contexts

Bowtie works by listening for an HTTP(S) connection and then creating an execution context that encapsulates its primary elements (that is, a request object and a response writer). It then executes a series of zero or more middlewares against this context until either some data is written to the output or an error occurs.

The [context](http://godoc.org/github.com/mtabini/go-bowtie#Context) is encapsulate by a simple interface that exposes the HTTP request and response writer. The server exposes an [AddContextFactory](http://godoc.org/github.com/mtabini/go-bowtie#Server.AddContextFactory) function that can be used to extend the functionality associated with the context with your own structs.

For example, suppose that you want to store come configuration parameters inside your own context—say, the URL to your database. You can create a custom context like such:

```go
package main

import (
    "github.com/mtabini/go-bowtie"
    "net/http"
)

type MyContext struct {
    bowtie.Context
    DBURL string
}

func MyContextFactory(previous bowtie.Context) bowtie.Context {
    // Return an instance of our context that encapsulates the previous
    // context created for the server

    return &MyContext{
        Context: previous,
        DBURL:   "db://mydb/table1",
    }
}

func main() {
    s := bowtie.NewServer()

    s.AddContextFactory(MyContextFactory)

    http.ListenAndServe(":8000", s)
}
```

In this fashion, you (or your middlewares) can extend the context as many times as you need in order to add more functionality to it. At runtime, you can then cast the generic instance of `bowtie.Context` that the server passes around to the specific object you need and access it.

## Middlewares

In order to do so, you will need to write (or, at least, use) a middleware. Luckily, that's easy:

```go
func MyMiddleware(c bowtie.Context, next func()) {
    // Cast the context to our context and get the DB URL

    myC := c.(*MyContext)

    // Output the URL to the client

    c.Response().WriteString(myC.DBURL)
}
```

As you can see, the middleware is simply a function that receives a context interface as its first argument. We can then cast the context to our own specialized struct and take advantage of its functionality.

The second argument to a middleware is a reference to a function that can be called to delay the execution of the middleware until after all other middlewares have run. This is handy for things like logging and error management, but you can ignore it in most cases—in fact, your middleware does _not_ need to return anything; it can simply exit when it's done.

## Middleware providers

As you can imagine, the creation of new context types and middlewares often goes hand-in-hand; therefore, a Bowtie server also defines [`MiddlewareProvider`](http://godoc.org/github.com/mtabini/go-bowtie#MiddlewareProvider) interface that can be used to handily couple these two operations in a single call.

For example, we could refactor our context factory and middleware into a unified struct that gives us more reusability:

```go
type MyMiddlewareProvider struct {
    DBURL string
}

func (m *MyMiddlewareProvider) ContextFactory() bowtie.ContextFactory {
    return func(previous bowtie.Context) bowtie.Context {
        // Return an instance of our context that encapsulates the previous
        // context created for the server

        return &MyContext{
            Context: previous,
            DBURL:   m.DBURL,
        }
    }
}

func (m *MyMiddlewareProvider) Middleware() bowtie.Middleware {
    return func(c bowtie.Context, next func()) {
        // Cast the context to our context and get the DB URL

        myC := c.(*MyContext)

        // Output the URL to the client

        c.Response().WriteString(myC.DBURL)
    }
}

func main() {
    s := bowtie.NewServer()

    s.AddMiddlewareProvider(&MyMiddlewareProvider{DBURL: "db:/my/database"})

    http.ListenAndServe(":8000", s)
}
```

As you can see, in addition to keeping things simple by not having to add both a middleware and a context, we also gain the ability to let the developer choose the database URL when she instantiates the provider at runtime.

## Bowtie's response writer

Bowtie replaces `http.ResponseWriter` with [its own interface](http://godoc.org/github.com/mtabini/go-bowtie#ResponseWriter) that provides a few additional convenience methods for writing strings and serializing JSON, as well as managing errors. In particular, you may find the functions in the form

```go
ResponseWriter.f(data, error)
```

useful for dealing with the common Go pattern of returning a `(data, error)` tuple from a function call. If the error is not nil, it is added to the response writer; otherwise, the data is written to the output stream. For example:

```go

func f() (string, error) {
    return "test", nil
}

func middleware (c bowtie.Context, next func()) {
    c.Response().WriteStringOrError(f())
}
```

Since the response writer is an interface (downwards compatible with `http.ResponseWriter`, of course), you are free to extend its functionality as needed. Simply create your own [`ResponseWriterFactory`](http://godoc.org/github.com/mtabini/go-bowtie#ResponseWriterFactory) function and set your server's [eponymous property](http://godoc.org/github.com/mtabini/go-bowtie#Server) to it.

## Error management

In an attempt to minimize the likelihood of sensitive data leakage, Bowtie encapsulates errors inside a special [`Error`](http://godoc.org/github.com/mtabini/go-bowtie#Error) interface that can be safely written to the output stream. 

Each error has an associated status code, which is automatically written to the response writer. If the status code is greater than 499, the error is assumed to be a server-side problem, and a generic message is written to the stream instead of the actual error message. Naturally, this is also true if the error is marshalled to JSON.

New `Error` instances can be created by calling [`NewError`](http://godoc.org/github.com/mtabini/go-bowtie#NewError), which also offers `fmt.Printf`-like templating capabilities. You can also convert an existing Go `error` to `Error` by calling [`NewErrorWithError`](http://godoc.org/github.com/mtabini/go-bowtie#NewErrorWithError); in this case, the new error instance is assigned a status code of 500.

Bowtie's response writer maintains its own error array, which is populated either by calling one of the error-aware convenience functions, or by explicitly adding new instances with a call to `AddError`. Errors passed into the array are automatically encapsulated by an `Error`-compliant struct that renders them safe for output.

Note, however, that the server doesn't actually do anything with the error array; in order to output the errors to the users, you will need to use the `Error` middleware (see below).









