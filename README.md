# Bowtie - Web middleware for Go

Bowtie is an HTTP middleware for Go. It makes heavy use of interfaces to provide a programming model that is both idiomatic and easily extensible (with, hopefully, relatively minimal overhead).

## Getting started

Standing up a [basic Bowtie](#The-quick-server) server with a default router and a few useful middlewares only requires a few lines of code:

```go
package main

import (
    "github.com/mtabini/go-bowtie"
    "github.com/mtabini/go-bowtie/middleware"
    "github.com/mtabini/go-bowtie/quick"
    "net/http"
)

func main() {
    // Create a new Bowtie server
    s := quick.New()

    s.GET("/test/:id", func(c bowtie.Context) {
        id := c.(*middleware.RouterContext).Params.ByName("id")

        c.Response().WriteString("The ID is " + id)
    })

    http.ListenAndServe(":8000", s)
}
```

Like all comparable systems, however, Bowtie becomes powerful once you start adding middlewares to it. Let's start, then with the most basic server—one that does, well, _nothing:_

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

## Contexts

Bowtie works by taking an HTTP request associating it with an [_execution context_](http://godoc.org/github.com/mtabini/go-bowtie#Context) that encapsulates its primary elements (that is, a request object and a response writer). It then executes a series of zero or more middlewares against this context until either some data is written to the output or an error occurs.

The server exposes an [AddContextFactory](http://godoc.org/github.com/mtabini/go-bowtie#Server.AddContextFactory) function that can be used to extend the functionality associated with the context with your own structs. You can add multiple context factories and have each build a new context on top of the previous one.

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

In order to do really do anything with Bowtie, you will need to write (or, at least, use) one or more middlewares. Luckily, that's easy:

```go
func MyMiddleware(c bowtie.Context, next func()) {
    // Cast the context to our context and get the DB URL

    myC := c.(*MyContext)

    // Output the URL to the client

    c.Response().WriteString(myC.DBURL)
}
```

As you can see, the middleware is simply a function that receives a context interface as its first argument. We can then cast the context to our own specialized struct and take advantage of its functionality if necessary.

The second argument to a middleware is a reference to a function that can be called to delay the execution of the middleware until after all other middlewares have run. This is handy for things like logging and error management, but you can ignore it in most cases—in fact, your middleware does _not_ need to return anything; it can simply exit when it's done.

## Middleware providers

As you can imagine, the creation of new context types and middlewares often goes hand-in-hand; therefore, a Bowtie server also defines a [`MiddlewareProvider`](http://godoc.org/github.com/mtabini/go-bowtie#MiddlewareProvider) interface that can be used to handily couple these two operations in a single call.

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

## Bowtie's request and response writer

Bowtie extends `http.Request` with a handful of functions designed to make reading the request's data a bit easier.

Also replaced is `http.ResponseWriter`; Bowtie comes with [its own interface](http://godoc.org/github.com/mtabini/go-bowtie#ResponseWriter) that provides a few additional convenience methods for writing strings and serializing JSON, as well as managing errors. In particular, you may find the functions in the form

```go
ResponseWriter.WriteXXXOrError(data, error)
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

## Routing

Although you are free to use your own router, Bowtie comes with a slightly modified copy of Julien Schmidt's trie-based [httprouter](https://github.com/julienschmidt/httprouter) that is ready to go:

```go
package main

import (
    "github.com/mtabini/go-bowtie"
    "github.com/mtabini/go-bowtie/middleware"
    "net/http"
    "strconv"
)

func echoValue(c bowtie.Context) {
    id := c.(*middleware.RouterContext).Params.ByName("id")

    c.Response().WriteString("The ID is " + id)
}

func validateValue(c bowtie.Context) {
    id := c.(*middleware.RouterContext).Params.ByName("id")

    if _, err := strconv.ParseInt(id, 10, 64); err != nil {
        c.Response().AddError(bowtie.NewError(400, "Invalid, non-numeric ID %s", id))
    }
}

func ExampleServer_routing() {
    // Create a new Bowtie server
    s := bowtie.NewServer()

    // Register middlewares

    r := middleware.NewRouter()

    r.GET("/test/:id", echoValue)
    r.GET("/validate/:id", validateValue, echoValue)

    s.AddMiddleware(middleware.ErrorReporter)
    s.AddMiddlewareProvider(r)

    // bowtie.Server can be used directly with http.ListenAndServe
    http.ListenAndServe(":8000", s)
}
```

The main changes from the Schmidt's original router are as follows:

- Bowtie's router defines its own context, called `RouterContext`; this contains a `Params` property that encapsulates the router's parameters.
- Bowtie's version passes `bowtie.Context` to its handlers instead of instances of the HTTP request and response writer.
- Bowtie's version supports multiple handlers per router, which are chained together and executed in sequence until either the end of the list is reached or one of the handlers writes to the output stream. The example above takes advantage of this feature by prepending the `validateValue` handler to `echoValue` and only allowing the latter to run if the data passed to by the client satisfies certain criteria.
- Finally, Bowtie's router introduces a `GetSupportedMethods` function that can be used to determine which HTTP methods are supported for a given route. This, in turn, is used by the [CORS middleware](https://godoc.org/github.com/mtabini/go-bowtie/middleware#NewCORSHandler) to respond to `OPTIONS` requests properly.

## Bundled middlewares

Bowtie comes bundled with a few more middlewares:

- [CORSHandler](https://godoc.org/github.com/mtabini/go-bowtie/middleware#CORSHandler) makes quick work of handling CORS requests, and information from the router to provide precise answers to `OPTIONS` pre-flight requests.
- [ErrorReporter](https://godoc.org/github.com/mtabini/go-bowtie/middleware#ErrorReporter) handles errors and outputs them safely to the output stream.
- [Logger](https://godoc.org/github.com/mtabini/go-bowtie/middleware#Logger) handles logging. It comes with both plaintext and [Bunyan](https://github.com/mtabini/go-bunyan) output handlers.
- [Recovery](https://godoc.org/github.com/mtabini/go-bowtie/middleware#Recovery) handles panics gracefully, turning them into 500 errors and capturing all the appropriate details for later logging.

## The quick server

What about the `quick` server listed at the beginning? It provides a simple set of defaults that gets you going very quickly; by calling `New()`, you get a struct that contains both a `bowtie.Server` and a `bowtie/middleware.Router`, plus a few pre-set middlewares, roughly equivalent to:

```go
func New() *QuickServer {
    r := middleware.NewRouter()

    s := bowtie.NewServer()

    s.AddMiddleware(middleware.NewLogger(middleware.MakePlaintextLogger()))
    s.AddMiddleware(middleware.Recovery)
    s.AddMiddleware(middleware.ErrorReporter)

    cors := middleware.NewCORSHandler(r)

    cors.SetDefaults()

    s.AddMiddlewareProvider(cors)

    s.AddMiddlewareProvider(r)

    return &QuickServer{
        s,
        r,
    }
}
```

## Bowtie compared to other Go frameworks

Bowtie borrows liberally—both in ideas and in code—from several other Go frameworks. Adopting [httprouter](https://github.com/julienschmidt/httprouter) as the default seemed like a good idea, given its raw speed. The only changes made were integrating Julien Schmidt's code with Bowtie's context-based execution, and allowing multiple handlers to be attached to a particular route.

Bowtie was also inspired by [Go-Martini](https://github.com/go-martini/martini)'s simplicity and immediateness. Even though Martini's design is not idiomatic to Go, it is perhaps one of the easiest frameworks to pick up and use. In addition to adopting bits of its code, Bowtie strives for the same kind of approachable and immediateness.

Finally, the idea of a running context is borrowed from [gin-gonic](https://github.com/gin-gonic/gin). The main difference between the two is that Bowtie forces the use of Go interface for carrying custom information through the context, whereas Gin allows you to append arbitrary data to it. This seemed like a more sensible approach that allows Go to do its job by providing type security and strictness.

Bowtie's error management comes from my personal fixation with safety. I don't want to leak information, and I don't want developers to worry about whether they will.

## Benchmarks

None. Most of the slowness in a web app comes from places other than the basic framework used to run it, and therefore benchmarks that are only concerned with the raw speed of a router are a bit misleading.

That said, Bowtie reliance on httprouter should mean that you can expect comparable speed from it—and perhaps a level of efficiency equivalent with Gin's.





