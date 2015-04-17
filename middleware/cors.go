package middleware

import (
	"github.com/mtabini/go-bowtie"
	"net/http"
	"strings"
)

func NewCORSHandler(router *Router, allowedOrigins []string) bowtie.Middleware {
	return func(c bowtie.Context, next func()) {
		req := c.Request()
		res := c.Response()

		header := res.Header()

		origin := req.Header.Get("Origin")

		if len(allowedOrigins) > 0 {
			found := false

			for _, o := range allowedOrigins {
				if o == origin {
					found = true
					break
				}
			}

			if !found {
				res.WriteHeader(http.StatusForbidden)
				return
			}
		}

		if origin == "" {
			origin = "*"
		}

		header.Set("Access-Control-Allow-Credentials", "true")
		header.Set("Access-Control-Allow-Origin", origin)
		header.Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, Authorization, If-None-Match, Range")
		header.Set("Access-Control-Expose-Headers", "Accept-Range, Content-Type, Content-Length, Content-Range, ETag, User-ETag,")

		if req.Method == "OPTIONS" {
			header.Set("Access-Control-Allow-Methods", strings.Join(router.GetSupportedMethods(req.URL.Path), ", "))

			res.WriteHeader(http.StatusNoContent)
		}
	}
}
