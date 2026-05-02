package router

import (
	"github.com/Mboukhal/AGT_SSR/core/auth"
	"github.com/go-chi/chi"
)

func RegisterRoutes(r chi.Router, svc *auth.Service) {
	r.Route("/api", func(r chi.Router) {
		r.Use(auth.WithService(svc))
		r.Use(auth.CSRFMiddleware([]string{
			"/api/auth/microsoft/callback",
			"/api/auth/csrf-token",
			"/api/auth/logout",
		}))
		auth.RouterHandler(r)
	})
}
