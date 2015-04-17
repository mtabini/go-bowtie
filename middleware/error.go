package middleware

import (
	"github.com/mtabini/bowtie"
)

func ErrorReporter(c bowtie.Context, next func()) {
	next()

	res := c.Response()

	errs := res.Errors
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
