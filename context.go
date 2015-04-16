package bowtie

import (
	"net/http"
)

// ContextFactory is a function that creates a context starting from previous context.
// Your application (and each middleware) can provide its own factory when the server is created,
// thus allowing you to create your own custom context with ease
type ContextFactory func(previous Context) Context

// Interface Context represents a server's context, which provides information used by the
// middleware. The basic context deals primarily with providing an interface to the request
// and response, as well as error management
type Context interface {
	// Request returns the request object associated with this request
	Request() *http.Request

	// Response returns the response writer associated with this request
	Response() *ResponseWriter
}

var _ Context = &ContextInstance{}

// Struct ContextInstance is a concrete implementation of the base server context. Your application
// can safely incorporate it into its own structs to extend the functionality provided by
// Bowtie
type ContextInstance struct {
	r *http.Request
	w *ResponseWriter
}

// NewContext is a ContextFactory that creates a basic context. You will probably want to create
// your own context and context factory that extends the basic context for your uses
func NewContext(r *http.Request, w http.ResponseWriter) Context {
	return &ContextInstance{
		r: r,
		w: NewResponseWriter(w),
	}
}

// Request returns the request associated with the context
func (c *ContextInstance) Request() *http.Request {
	return c.r
}

// Response returns the response writer assocaited with the context
func (c *ContextInstance) Response() *ResponseWriter {
	return c.w
}
