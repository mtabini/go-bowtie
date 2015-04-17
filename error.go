package bowtie

import (
	"encoding/json"
	"fmt"
)

// Interface Error represents a Bowtie error, which extends the standard error interface to provide
// additional Web-friendly functionality.
//
// Instances of Error must be able to provide both a public and private representation; the latter
// may include sensitive information like stack traces and private messages, while the former
// can be safely outputted to an end user.
type Error interface {
	error
	fmt.Stringer
	json.Marshaler
	// StatusCode return the error's status code
	StatusCode() int
	// Message returns the error's message
	Message() string
	// Data returns the error's associated data
	Data() interface{}
	// SetData sets the error's associated data
	SetData(interface{})
	// GetPrivateRepresentation a private representation of the error. Useful for logging.
	PrivateRepresentation() map[string]interface{}
	// GetStackTrace returns the stack trace associated with this error, if any
	StackTrace() []StackFrame
	// RecordStackTrace captures a stack track and return the error instance
	CaptureStackTrace() Error
}

// Struct ErrorInstance incorporates an error and associates it with an HTTP status code (assumed to be 500
// if not present. Arbitrary data can also be added to the error for logging purposes, as can a stack
// trace.
//
// The ErrorInstance struct is smart enough that, when asked for serialization either through error's Error()
// method or by JSON marshalling, it does not leak any sensitive information if the StatusCode is >=
// 500 (which indicates a server error).
//
// For status codes that indicate user errors ([400-499]), the struct allows public consumers to see
// the actual message that was provided at initialization time.
type ErrorInstance struct {
	statusCode int          // The HTTP status code
	message    string       // A message associated with the error. May be overwritten if the status code is >= 500
	data       interface{}  // Assorted data associated with the error, for logging purposes
	stackTrace []StackFrame // The stack trace associated with the error, for logging purposes
}

// NewError builds a new Error instance; the `format` and `arguments` parameters work as in `fmt.Sprintf()`
func NewError(statusCode int, format string, arguments ...interface{}) Error {
	return &ErrorInstance{
		statusCode: statusCode,
		message:    fmt.Sprintf(format, arguments...),
	}
}

// NewErrorFromError builds a new Error instance starting from a regular Go error (or something that
// can be cast to it). If an instance of Error is passed to it, the function returns a copy thereof
// (and not the original), but _not_ of the associated data, which may be copied by reference.
//
// If the error
func NewErrorWithError(err error) Error {
	if e, ok := err.(Error); ok {
		return &ErrorInstance{
			statusCode: e.StatusCode(),
			message:    e.Message(),
			data:       e.Data(),
			stackTrace: e.StackTrace(),
		}
	}

	return &ErrorInstance{
		statusCode: 500,
		message:    err.Error(),
	}
}

// Ensure that ErrorInstance always satisfies Error

var _ Error = &ErrorInstance{}

// Satisfy the error, fmt.Stringer, and json.Marshaler interfaces
func (e *ErrorInstance) Error() string {
	if e.statusCode > 499 {
		return "An server error has occurred."
	}

	return e.message
}

func (e *ErrorInstance) String() string {
	return e.Error()
}

func (e *ErrorInstance) MarshalJSON() ([]byte, error) {
	result := map[string]interface{}{
		"statusCode": e.statusCode,
		"message":    e.Error(),
	}

	return json.Marshal(result)
}

// Satisfy the Error interface

// Returns the status code associated with e
func (e *ErrorInstance) StatusCode() int {
	return e.statusCode
}

// Returns the message associated with e
func (e *ErrorInstance) Message() string {
	return e.message
}

// Returns the data associated with e
func (e *ErrorInstance) Data() interface{} {
	return e.data
}

// Sets the data associated with e
func (e *ErrorInstance) SetData(data interface{}) {
	e.data = data
}

// Returns a private representation of e
func (e *ErrorInstance) PrivateRepresentation() map[string]interface{} {
	return map[string]interface{}{
		"statusCode": e.statusCode,
		"message":    e.message,
		"data":       e.data,
		"stackTrace": e.stackTrace,
	}
}

func (e *ErrorInstance) StackTrace() []StackFrame {
	return e.stackTrace
}
