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
	// OUTPUT
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
}
