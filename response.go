package bowtie

import (
	"encoding/json"
	"net/http"
)

type ResponseWriter interface {
	http.ResponseWriter

	// Add error safely adds a new error to the context, converting it to bowtie.Error if appropriate
	AddError(err error)

	// Errors returns an array that contains any error assigned to the response writer
	Errors() []Error

	// Status returns the HTTP status code of the writer. You can set this by using `WriteHeader()`
	Status() int

	// Written returns true if any data (including a status code) has been written to the writer's
	// output stream
	Written() bool

	// WriteOrError checks if `err` is not nil, in which case it adds it to the context's error
	// list and returns. If `err` is nil, `p` is written to the output stream instead. This is a
	// convenient way of dealing with functions that return (data, error) tuples inside a middleware
	WriteOrError(p []byte, err error) (int, error)

	// WriteString is a convenience method that outputs a string
	WriteString(s string) (int, error)

	// WriteStringOrError is a convenience method that outputs a string or adds an error to the writer.
	// It works like `WriteOrError`, but takes string instead of a byte array
	WriteStringOrError(s string, err error) (int, error)

	// WriteJSON writes data in JSON format to the output stream. The output Content-Type header
	// is also automatically set to `application/json`
	WriteJSON(data interface{}) (int, error)

	// WriteJSONOrError checks if `err` is not nil, in which case it adds it to the context's error
	// list and returns. If `err` is nil, `data` is serialized to JSON and written to the output
	// stream instead; the Content-Type of the response is also set to `application/json` automatically.
	// This is a convenient way of dealing with functions that return (data, error) tuples inside
	// a middleware
	WriteJSONOrError(data interface{}, err error) (int, error)
}

type ResponseWriterInstance struct {
	http.ResponseWriter
	written bool
	errors  []Error
	status  int
}

var _ ResponseWriter = &ResponseWriterInstance{}

func NewResponseWriter(w http.ResponseWriter) ResponseWriter {
	return &ResponseWriterInstance{
		ResponseWriter: w,
		errors:         []Error{},
		status:         200,
	}
}

// Errors returns an array that contains any error assigned to the response writer
func (r *ResponseWriterInstance) Errors() []Error {
	return r.errors
}

// Add error safely adds a new error to the context, converting it to bowtie.Error if appropriate
func (r *ResponseWriterInstance) AddError(err error) {
	if e, ok := err.(Error); ok {
		r.WriteHeader(e.StatusCode())
	} else {
		r.WriteHeader(500)
	}

	r.errors = append(r.errors, NewErrorWithError(err))
}

// Status returns the HTTP status code of the writer. You can set this by using `WriteHeader()`
func (r *ResponseWriterInstance) Status() int {
	return r.status
}

// WriteHeader writes a status header
func (r *ResponseWriterInstance) WriteHeader(status int) {
	r.ResponseWriter.WriteHeader(status)
	r.status = status
	r.written = true
}

// Written returns true if any data (including a status code) has been written to the writer's
// output stream
func (r *ResponseWriterInstance) Written() bool {
	return r.written
}

// Write implements io.Writer and outputs data to the HTTP stream
func (r *ResponseWriterInstance) Write(p []byte) (int, error) {
	n, err := r.ResponseWriter.Write(p)

	if err != nil {
		r.written = true
	}

	return n, err
}

// WriteOrError checks if `err` is not nil, in which case it adds it to the context's error
// list and returns. If `err` is nil, `p` is written to the output stream instead. This is a
// convenient way of dealing with functions that return (data, error) tuples inside a middleware
func (r *ResponseWriterInstance) WriteOrError(p []byte, err error) (int, error) {
	if err != nil {
		r.AddError(err)
		return 0, err
	}

	return r.Write(p)
}

// WriteString is a convenience method that outputs a string
func (r *ResponseWriterInstance) WriteString(s string) (int, error) {
	return r.Write([]byte(s))
}

// WriteStringOrError is a convenience method that outputs a string or adds an error to the writer.
// It works like `WriteOrError`, but takes string instead of a byte array
func (r *ResponseWriterInstance) WriteStringOrError(s string, err error) (int, error) {
	return r.WriteOrError([]byte(s), err)
}

// WriteJSON writes data in JSON format to the output stream. The output Content-Type header
// is also automatically set to `application/json`
func (r *ResponseWriterInstance) WriteJSON(data interface{}) (int, error) {
	return r.WriteOrError(json.Marshal(data))
}

// WriteJSONOrError checks if `err` is not nil, in which case it adds it to the context's error
// list and returns. If `err` is nil, `data` is serialized to JSON and written to the output
// stream instead; the Content-Type of the response is also set to `application/json` automatically.
// This is a convenient way of dealing with functions that return (data, error) tuples inside
// a middleware
func (r *ResponseWriterInstance) WriteJSONOrError(data interface{}, err error) (int, error) {
	if err != nil {
		r.AddError(err)
		return 0, err
	}

	return r.WriteJSON(data)
}
