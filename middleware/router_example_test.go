package middleware

import (
	"github.com/mtabini/go-bowtie"
	"net/http"
	"strconv"
)

// echoValue is a handler that retrieves the parameter `id` from the
// router's context and outputs it back to the user.
//
// Note how the router creates its own Context instance, which
// allows it to add the new property `Params` that we can then use
// by re-casting the generic context.
func EchoValue(c bowtie.Context) {
	id := c.(*RouterContext).Params.ByName("id")

	c.Response().WriteString("The ID is " + id)
}

// validateValue checks that the parameter `id` supplied to the
// router is, in fact, an integer. Because Bowtie's router supports
// handler chanining, we can place this in front of `echoValue` and
// cause it to interrupt the chain if it detects an error.
//
// This makes splitting functionality and reusing code easier.
func ValidateValue(c bowtie.Context) {
	id := c.(*RouterContext).Params.ByName("id")

	if _, err := strconv.ParseInt(id, 10, 64); err != nil {
		c.Response().AddError(bowtie.NewError(400, "Invalid, non-numeric ID %s", id))
	}
}

// Output:
// > GET /test/123
//
// < HTTP/1.1 200 OK
// < The ID is 123
//
//
// > GET /validate/1234
//
// < HTTP/1.1 200 OK
// < The ID is 1234
//
// > GET /validate/invalid
//
// < HTTP/1.1 400 Bad Request
// < [{"message":"Invalid, non-numeric ID 123s","statusCode":400}]
func ExampleServer_routing() {
	// Create a new Bowtie server
	s := bowtie.NewServer()

	// Register middlewares

	r := NewRouter()

	// Define our routes

	r.GET("/test/:id", EchoValue)
	r.GET("/validate/:id", ValidateValue, EchoValue)

	s.AddMiddleware(ErrorReporter)
	s.AddMiddlewareProvider(r)

	// Server can be used directly with http.ListenAndServe

	http.ListenAndServe(":8000", s)
}
