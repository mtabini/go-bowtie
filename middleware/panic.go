// Borrowed from https://github.com/go-martini/martini/blob/master/recovery.go
package middleware

import (
	"github.com/mtabini/go-bowtie"
	"net/http"
)

// Recovery returns a middleware that recovers from any panics and writes a 500 if there was one.
// While Martini is in development mode, Recovery will also output the panic as HTML.
func Recovery(c bowtie.Context, next func()) {
	defer func() {
		if err := recover(); err != nil {
			e := bowtie.NewError(http.StatusInternalServerError, "panic: %#v", err)
			e.CaptureStackTrace()

			c.Response().AddError(e)
		}
	}()

	next()
}
