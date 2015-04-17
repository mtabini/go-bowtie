package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mtabini/go-bowtie"
	"github.com/mtabini/go-bunyan"
	"log"
	"time"
)

// Logger is a generic logger function that takes a bowtie context and
// is tasked with logging the request it encapsulates. You pass this
// value to NewLogger and then add the resulting middleware to your server.
//
// The logger middleware comes with two output handlers: plaintext and Bunyan.
// For example:
//
//    s := bowtie.NewServer()
//
//    s.AddMiddleware(middleware.NewLogger(middleware.MakePlaintextLogger()))
type Logger func(c bowtie.Context)

// MakePlaintextLogger logs requests to standard output using this space-limited simple format:
// RemoteAddress Method URL Status RunningTime
func MakePlaintextLogger() Logger {
	return func(c bowtie.Context) {
		req := c.Request()
		res := c.Response()

		log.Printf("%s %s %s %d %f", req.RemoteAddr, req.Method, req.URL, res.Status(), float64(c.GetRunningTime())/float64(time.Second))
	}
}

// BunyanLogger logs requests using a Bunyan logger. See https://github.com/mtabini/go-bunyan
// for more information
func MakeBunyanLogger(logger bunyan.Logger) Logger {
	return func(c bowtie.Context) {
		req := c.Request()
		res := c.Response()

		e := bunyan.NewLogEntry(bunyan.Info, fmt.Sprintf("%s %s", req.Method, req.URL.RequestURI()))

		e.SetRequest(req)
		e.SetResponseStatusCode(res.Status())

		e.SetCompletedIn(fmt.Sprintf("%v", c.GetRunningTime()))

		errs := res.Errors()

		if len(errs) > 0 {
			outErrs := make([]map[string]interface{}, len(errs))

			for index, err := range errs {
				outErrs[index] = err.PrivateRepresentation()
			}

			outErr, _ := json.Marshal(outErrs)

			e.Level = bunyan.Error
			e.SetResponseError(errors.New(string(outErr)))
		}

		logger.Log(e)
	}
}

// NewLogger creates a new logger middleware. It waits until all other
// middlewares have finished running, then calls `logger` with the
// request's context.
func NewLogger(logger Logger) bowtie.Middleware {
	return func(c bowtie.Context, next func()) {
		next()

		logger(c)
	}
}
