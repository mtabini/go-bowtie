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

type Logger func(c bowtie.Context)

func PlaintextLogger(c bowtie.Context) {
	req := c.Request()
	res := c.Response()

	log.Printf("%s %s %s %d %f", req.RemoteAddr, req.Method, req.URL, res.Status(), float64(c.GetRunningTime())/float64(time.Second))
}

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

func NewLogger(logger Logger) bowtie.Middleware {
	return func(c bowtie.Context, next func()) {
		next()

		logger(c)
	}
}
