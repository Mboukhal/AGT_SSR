package auth

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/csrf"
)

type contextKey string

const (
	AuthServiceKey contextKey = "auth_service"
)

func WithService(svc *Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), AuthServiceKey, svc)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetService(r *http.Request) *Service {
	svc, _ := r.Context().Value(AuthServiceKey).(*Service)
	return svc
}

func CSRFMiddleware(ignorePaths []string) func(http.Handler) http.Handler {
	csrfKey := os.Getenv("CSRF_SECRET")
	if csrfKey == "" {
		csrfKey = "csrf-secret-key-change-in-production"
	}

	csrfMiddleware := csrf.Protect(
		[]byte(csrfKey),
		csrf.CookieName("_csrf"),
		csrf.Secure(false),
		csrf.Path("/"),
		csrf.HttpOnly(true),
		csrf.SameSite(csrf.SameSiteLaxMode),
	)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for _, path := range ignorePaths {
				if r.URL.Path == path {
					next.ServeHTTP(w, r)
					return
				}
			}
			csrfMiddleware(next).ServeHTTP(w, r)
		})
	}
}

func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc := GetService(r)
		if svc == nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		session, err := svc.Store.Get(r, SessionName)
		if err != nil {
			http.Error(w, "invalid session", http.StatusUnauthorized)
			return
		}

		_, _, ok := svc.GetUserFromSession(session)
		if !ok {
			http.Error(w, "not authenticated", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

func PageAuth(svc *Service) func(http.Handler) http.Handler {
	publicPaths := map[string]bool{
		"/login": true,
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/_astro") || strings.HasPrefix(r.URL.Path, "/api") {
				next.ServeHTTP(w, r)
				return
			}

			session, _ := svc.Store.Get(r, SessionName)
			_, _, loggedIn := svc.GetUserFromSession(session)

			if !loggedIn && !publicPaths[r.URL.Path] {
				http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
				return
			}

			if loggedIn && publicPaths[r.URL.Path] {
				http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
