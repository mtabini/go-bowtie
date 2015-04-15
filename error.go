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
	// GetStatusCode return the error's status code
	GetStatusCode() int
	// GetMessage returns the error's message
	GetMessage() string
	// GetData returns the error's associated data
	GetData() interface{}
	// GetPrivateRepresentation a private representation of the error. Useful for logging.
	GetPrivateRepresentation() map[string]interface{}
	// GetStackTrace returns the stack trace associated with this error, if any
	GetStackTrace() []StackFrame
	// RecordStackTrace captures a stack track and return the error instance
	RecordStackTrace() Error
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
	StatusCode int          // The HTTP status code
	Message    string       // A message associated with the error. May be overwritten if the status code is >= 500
	Data       interface{}  // Assorted data associated with the error, for logging purposes
	StackTrace []StackFrame // The stack trace associated with the error, for logging purposes
}

// NewError builds a new Error instance; the `format` and `arguments` parameters work as in `fmt.Sprintf()`
func NewError(statusCode int, format string, arguments ...interface{}) Error {
	return &ErrorInstance{
		StatusCode: statusCode,
		Message:    fmt.Sprintf(format, arguments...),
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
			StatusCode: e.GetStatusCode(),
			Message:    e.GetMessage(),
			Data:       e.GetData(),
			StackTrace: e.GetStackTrace(),
		}
	}

	return &ErrorInstance{
		StatusCode: 500,
		Message:    err.Error(),
	}
}

// Ensure that ErrorInstance always satisfies Error

var _ Error = &ErrorInstance{}

// Satisfy the error, fmt.Stringer, and json.Marshaler interfaces
func (e *ErrorInstance) Error() string {
	if e.StatusCode > 499 {
		return "An server error has occurred."
	}

	return e.Message
}

func (e *ErrorInstance) String() string {
	return e.Error()
}

func (e *ErrorInstance) MarshalJSON() ([]byte, error) {
	result := map[string]interface{}{
		"statusCode": e.StatusCode,
		"message":    e.Error(),
	}

	return json.Marshal(result)
}

// Satisfy the Error interface

// Returns the status code associated with e
func (e *ErrorInstance) GetStatusCode() int {
	return e.StatusCode
}

// Returns the message associated with e
func (e *ErrorInstance) GetMessage() string {
	return e.Message
}

// Returns the data associated with e
func (e *ErrorInstance) GetData() interface{} {
	return e.Data
}

// Returns a private representation of e
func (e *ErrorInstance) GetPrivateRepresentation() map[string]interface{} {
	return map[string]interface{}{
		"statusCode": e.StatusCode,
		"message":    e.Message,
		"data":       e.Data,
		"stackTrace": e.StackTrace,
	}
}

func (e *ErrorInstance) GetStackTrace() []StackFrame {
	return e.StackTrace
}
