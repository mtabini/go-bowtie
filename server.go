// Package bowtie provides a Web middleware for Go apps.
//
// More information at https://github.com/mtabini/go-bowtie
//
// For a quick start, check out the examples at http://godoc.org/github.com/mtabini/go-bowtie/quick
package bowtie

import (
	"net/http"
)

// Middleware is a function that encapsulate a Bowtie middleware. It receives an execution
// context and a next function it can optionally call to delay its own execution until the
// end of the request (useful for logging, error handling, etc.)
type Middleware func(c Context, next func())

// Interface MiddleProvider can be implemented by structs that want to offer both
// a middleware and a context factory. They can be installed onto a server by calling
// AddMiddlewareProvider()
type MiddlewareProvider interface {
	Middleware() Middleware
	ContextFactory() ContextFactory
}

// Struct Server is a Bowtie server. It provides a handler compatible with http.ListenAndServe
// that creates a context and executes any attached middleware.
type Server struct {
	middlewares           []Middleware
	contextFactories      []ContextFactory
	ResponseWriterFactory ResponseWriterFactory
}

// NewServer initializes and returns a new Server instance.
func NewServer() *Server {
	return &Server{
		middlewares:           []Middleware{},
		contextFactories:      []ContextFactory{},
		ResponseWriterFactory: NewResponseWriter,
	}
}

// SetContextFactory changes the context factory used by the server. This allows you
// to create your own Context structs and use them inside your apps.
func (s *Server) AddContextFactory(value ContextFactory) {
	s.contextFactories = append(s.contextFactories, value)
}

// AddMiddleware adds a new middleware handler. Handlers are executed in the order
// in which they are added to the server
func (s *Server) AddMiddleware(f Middleware) {
	s.middlewares = append(s.middlewares, f)
}

// AddMiddlewareProvider registers a new middleware provider
func (s *Server) AddMiddlewareProvider(p MiddlewareProvider) {
	if mw := p.Middleware(); mw != nil {
		s.middlewares = append(s.middlewares, mw)
	}

	if cf := p.ContextFactory(); cf != nil {
		s.AddContextFactory(cf)
	}
}

// NewContext creates a new basic server context. You should not need to call this
// except for testing purposes. Instead, you should extend the server context
// with your struct and provide a context factory to the server
func (s *Server) NewContext(r *http.Request, w http.ResponseWriter) Context {
	c := NewContext(r, s.ResponseWriterFactory(w))

	for _, factory := range s.contextFactories {
		c = factory(c)
	}

	return c
}

// Run is the server's main entry point. It executes each middleware in sequence
// until one of them causes data to be written to the output
func (s *Server) Run(c Context) {
	mwIndex := -1
	mwCount := len(s.middlewares)

	if body := c.Request().Body; body != nil {
		defer body.Close()
	}

	var next func()

	next = func() {
		mwIndex += 1

		for mwIndex < mwCount {
			s.middlewares[mwIndex](c, next)
			mwIndex += 1

			if c.Response().Written() {
				return
			}
		}
	}

	next()
}

// ServeHTTP handles requests and can be used as a handler for http.Server
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Body != nil {
		defer r.Body.Close()
	}

	s.Run(s.NewContext(r, w))
}
