package router

import (
	ep "github.com/Mboukhal/SvGoPg/core/auth/ep"
	"github.com/go-chi/chi/v5"
)

// RegisterRoutes sets up the OAuth routes on the given router.
func RegisterRoutes(r chi.Router) {

	r.Route("/api", func(r chi.Router) {
		ep.RouterHandler(r)
	})
}
