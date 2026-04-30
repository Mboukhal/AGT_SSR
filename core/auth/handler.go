package auth

import (
	"net/http"

	"github.com/go-chi/chi"
)

func RouterHandler(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", LoginHandler)
	})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	return
}
