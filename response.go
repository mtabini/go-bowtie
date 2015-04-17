package bowtie

import (
	"encoding/json"
	"net/http"
)

type ResponseWriter struct {
	http.ResponseWriter
	Written bool
	Errors  []Error
	Status  int
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		Errors:         []Error{},
		Status:         200,
	}
}

// Add error safely adds a new error to the context, converting it to bowtie.Error if appropriate
func (r *ResponseWriter) AddError(err error) {
	if e, ok := err.(Error); ok {
		r.WriteHeader(e.StatusCode())
	} else {
		r.WriteHeader(500)
	}

	r.Errors = append(r.Errors, NewErrorWithError(err))
}

// WriteHeader writes a status header
func (r *ResponseWriter) WriteHeader(status int) {
	r.ResponseWriter.WriteHeader(status)
	r.Status = status
	r.Written = true
}

// Write implements io.Writer and outputs data to the HTTP stream
func (r *ResponseWriter) Write(p []byte) (int, error) {
	n, err := r.ResponseWriter.Write(p)

	if err != nil {
		r.Written = true
	}

	return n, err
}

// WriteOrError checks if `err` is not nil, in which case it adds it to the context's error
// list and returns. If `err` is nil, `p` is written to the output stream instead. This is a
// convenient way of dealing with functions that return (data, error) tuples inside a middleware
func (r *ResponseWriter) WriteOrError(p []byte, err error) (int, error) {
	if err != nil {
		r.AddError(err)
		return 0, err
	}

	return r.Write(p)
}

// WriteJSON writes data in JSON format to the output stream. The output Content-Type header
// is also automatically set to `application/json`
func (r *ResponseWriter) WriteJSON(data interface{}) (int, error) {
	return r.WriteOrError(json.Marshal(data))
}

// WriteJSONOrError checks if `err` is not nil, in which case it adds it to the context's error
// list and returns. If `err` is nil, `data` is serialized to JSON and written to the output
// stream instead; the Content-Type of the response is also set to `application/json` automatically.
// This is a convenient way of dealing with functions that return (data, error) tuples inside
// a middleware
func (r *ResponseWriter) WriteJSONOrError(data interface{}, err error) (int, error) {
	if err != nil {
		r.AddError(err)
		return 0, err
	}

	return r.WriteJSON(data)
}
