package middleware

import (
	"github.com/mtabini/go-bowtie"
	"net/http"
	"strings"
)

// Struct CORSHandler provides CORS support. It can automatically use an instance of
// Router to provide accurate responses to an OPTION preflight request, and supports
// restricting which headers are allowed in input and output.
//
// CORSHandler conforms to the bowtie.MiddlewareProvided interface.
//
// A set of sensible defaults can be installed by calling the SetDefaults() method.
type CORSHandler struct {
	router         *Router
	AllowedOrigins []string
	AllowedHeaders []string
	ExposedHeaders []string
}

func (h *CORSHandler) handle(c bowtie.Context, next func()) {
	req := c.Request()
	res := c.Response()

	header := res.Header()

	origin := req.Header.Get("Origin")

	if len(h.AllowedOrigins) > 0 {
		found := false

		for _, o := range h.AllowedOrigins {
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

	if len(h.AllowedHeaders) > 0 {
		header.Set("Access-Control-Allow-Headers", strings.Join(h.AllowedHeaders, ", "))
	}

	if len(h.ExposedHeaders) > 0 {
		header.Set("Access-Control-Expose-Headers", strings.Join(h.ExposedHeaders, ", "))
	}

	if req.Method == "OPTIONS" {
		header.Set("Access-Control-Allow-Methods", strings.Join(h.router.GetSupportedMethods(req.URL.Path), ", "))

		res.WriteHeader(http.StatusNoContent)
	}
}

// SetDefaults sets a basic set of defaults. Allows any origin and exposes commonly-used headers both
// in input and output
func (c *CORSHandler) SetDefaults() {
	c.AllowedHeaders = []string{"Content-Type", "Content-Length", "Accept-Encoding", "Authorization", "If-None-Match", "Range"}
	c.ExposedHeaders = []string{"Accept-Range", "Content-Type", "Content-Length", "Content-Range", "ETag"}
}

func (h *CORSHandler) Middleware() bowtie.Middleware {
	return h.handle
}

func (h *CORSHandler) ContextFactory() bowtie.ContextFactory {
	return nil
}

// NewCORSHandler creates a new CORS Handler that uses `router` to determine which HTTP methods
// are acceptable for a given route.
func NewCORSHandler(router *Router) *CORSHandler {
	return &CORSHandler{
		router:         router,
		AllowedOrigins: []string{},
		AllowedHeaders: []string{},
		ExposedHeaders: []string{},
	}
}
