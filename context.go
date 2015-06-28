package bowtie

import (
	"net/http"
	"sync/atomic"
	"time"
)

// ContextFactory is a function that processes a context.
// Your application (and each middleware) can provide its own factory when the server is created,
// thus allowing you to set new values into the context as needed
type ContextFactory func(context Context)

type ContextKey int64

var currentContextKey int64 = 0

func GenerateContextKey() ContextKey {
	return ContextKey(atomic.AddInt64(&currentContextKey, 1))
}

// Interface Context represents a server's context, which provides information used by the
// middleware. The basic context deals primarily with providing an interface to the request
// and response
type Context interface {
	// Get returns a property set into the context
	Get(ContextKey) interface{}

	// Set sets a new property into the context
	Set(ContextKey, interface{})

	// Request returns the request object associated with this request
	Request() *Request

	// Response returns the response writer associated with this request
	Response() ResponseWriter

	// GetRunningTime returns the amount of time during which this request has been running
	GetRunningTime() time.Duration
}

var _ Context = &ContextInstance{}

// Struct ContextInstance is a concrete implementation of the base server context. Your application
// can safely incorporate it into its own structs to extend the functionality provided by
// Bowtie
type ContextInstance struct {
	r         *Request
	w         ResponseWriter
	values    map[ContextKey]interface{}
	startTime time.Time
}

// NewContext is a ContextFactory that creates a basic context. You will probably want to create
// your own context and context factory that extends the basic context for your uses
func NewContext(r *http.Request, w http.ResponseWriter) Context {
	return &ContextInstance{
		r:         NewRequest(r),
		w:         NewResponseWriter(w),
		values:    map[ContextKey]interface{}{},
		startTime: time.Now(),
	}
}

// Request returns the request associated with the context
func (c *ContextInstance) Request() *Request {
	return c.r
}

func (c *ContextInstance) Get(key ContextKey) interface{} {
	return c.values[key]
}

func (c *ContextInstance) Set(key ContextKey, value interface{}) {
	c.values[key] = value
}

// Response returns the response writer assocaited with the context
func (c *ContextInstance) Response() ResponseWriter {
	return c.w
}

// GetRunningTime returns the amount of time during which this request has been running
func (c *ContextInstance) GetRunningTime() time.Duration {
	return time.Now().Sub(c.startTime)
}
