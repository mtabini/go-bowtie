package bowtie

import (
	"net/http"
)

// Struct MyContext extends Bowtie's context and adds more features to it
type MyContext struct {
	Context
	DBURL string
}

// Struct MyMiddlewareProvider satisfies the MiddlewareProvider interface
// and provides both a middleware and a context factory
type MyMiddlewareProvider struct {
	DBURL string
}

// ContextFactory, when added to a server, “wraps” our context around an
// existing context. At execution time, the middleware can then cast
// the context that the server passes to it to MyContext and take
// advantage of its functionality.
func (m *MyMiddlewareProvider) ContextFactory() ContextFactory {
	return func(previous Context) Context {
		// Return an instance of our context that encapsulates the previous
		// context created for the server

		return &MyContext{
			Context: previous,
			DBURL:   m.DBURL,
		}
	}
}

// A middleware is simply a function that takes a context, which it can
// use to manipulate the current HTTP request, and a next function that
// can be called to delay the middleware's execution until after all
// other middlewares have run.
func (m *MyMiddlewareProvider) Middleware() Middleware {
	return func(c Context, next func()) {
		// Cast the context to our context and get the DB URL

		myC := c.(*MyContext)

		// Output the URL to the client

		c.Response().WriteString(myC.DBURL)
	}
}

func ExampleServer_middleware() {
	// Create a new Bowtie server
	s := NewServer()

	// Register our new middleware provider. This adds our context factory
	// to it, and injects our middleware into its execution queue.
	s.AddMiddlewareProvider(&MyMiddlewareProvider{DBURL: "db:/my/database"})

	// Server can be used directly with http.ListenAndServe
	http.ListenAndServe(":8000", s)
}
