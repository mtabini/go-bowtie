package middleware

import (
	"github.com/mtabini/go-bowtie"
)

// ErrorReporter is a middleware that safely handles error reporting
// by outputting the errors that have accumulated in the context's response
// writer. It computes the status of a request from the maximum response
// status of all the errors (if any are present).
func ErrorReporter(c bowtie.Context, next func()) {
	next()

	res := c.Response()

	errs := res.Errors()
	outErrs := []bowtie.Error{}

	if len(errs) > 0 {
		maxStatus := 0

		for _, err := range errs {
			if err.StatusCode() < 500 {
				outErrs = append(outErrs, err)
			}
		}

		if maxStatus >= 500 {
			outErrs = append(outErrs, bowtie.NewError(500, "A server error has occurred"))
		}

		c.Response().WriteJSON(outErrs)
	}
}
