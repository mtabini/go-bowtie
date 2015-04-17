package quick

import (
	"github.com/mtabini/go-bowtie"
	"github.com/mtabini/go-bowtie/middleware"
)

type QuickServer struct {
	*bowtie.Server
	*middleware.Router
}

func New() *QuickServer {
	r := middleware.NewRouter()

	s := bowtie.NewServer()

	s.AddMiddleware(middleware.NewLogger(middleware.PlaintextLogger))
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
