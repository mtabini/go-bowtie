// Package quick provides the easiest way to stand up a Bowtie server.
//
// The resulting server includes a plaintext logger, error handler, panic recovery, and
// a router.
package quick

import (
	"github.com/mtabini/go-bowtie"
	"github.com/mtabini/go-bowtie/middleware"
)

// Struct QuickServer encapsulates a Bowtie server and router. It can be passed
// directly to net/http.ListenAndServe, and exposes all the router's methods.
type QuickServer struct {
	*bowtie.Server
	*middleware.Router
}

// New creates a new QuickServer.
func New() *QuickServer {
	r := middleware.NewRouter()

	s := bowtie.NewServer()

	s.AddMiddleware(middleware.NewLogger(middleware.MakePlaintextLogger()))
	s.AddMiddleware(middleware.Recovery)
	s.AddMiddleware(middleware.ErrorReporter)

	cors := middleware.NewCORSHandler(r)

	cors.SetDefaults()

	s.AddMiddlewareProvider(cors)

	s.AddMiddlewareProvider(r)

	return &QuickServer{
		s,
		r,
	}
}
