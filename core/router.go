package router

import (
	"github.com/Mboukhal/SvGoPg/core/auth"
	"github.com/go-chi/chi"
)

// RegisterRoutes sets up the OAuth routes on the given router.
func RegisterRoutes(r chi.Router) {

	r.Route("/api", func(r chi.Router) {
		auth.RouterHandler(r)
	})
}
