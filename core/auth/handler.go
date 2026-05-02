package auth

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/gorilla/csrf"
)

func RouterHandler(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		r.Get("/microsoft/login", MicrosoftLoginHandler)
		r.Get("/microsoft/callback", MicrosoftCallbackHandler)
		r.Get("/logout", LogoutHandler)
		r.Post("/logout", LogoutHandler)
		r.Get("/me", RequireAuth(MeHandler))
		r.Get("/csrf-token", CSRFTokenHandler)
	})
}

func MicrosoftLoginHandler(w http.ResponseWriter, r *http.Request) {
	svc := GetService(r)
	if svc == nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	authURL := svc.AuthURL("random-state-placeholder")
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func MicrosoftCallbackHandler(w http.ResponseWriter, r *http.Request) {
	svc := GetService(r)
	if svc == nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		http.Redirect(w, r, "/login?error=oauth_denied", http.StatusTemporaryRedirect)
		return
	}

	token, err := svc.Exchange(r.Context(), code)
	if err != nil {
		http.Redirect(w, r, "/login?error=oauth_failed", http.StatusTemporaryRedirect)
		return
	}

	msUser, err := svc.GetMicrosoftUser(r.Context(), token)
	if err != nil {
		http.Redirect(w, r, "/login?error=user_fetch_failed", http.StatusTemporaryRedirect)
		return
	}

	session, err := svc.CreateSession(r, msUser)
	if err != nil {
		http.Redirect(w, r, "/login?error=session_error", http.StatusTemporaryRedirect)
		return
	}

	if err := session.Save(r, w); err != nil {
		http.Redirect(w, r, "/login?error=session_save", http.StatusTemporaryRedirect)
		return
	}

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	svc := GetService(r)

	if svc != nil {
		session, err := svc.Store.Get(r, SessionName)
		if err == nil {
			session.Options.MaxAge = -1
			_ = session.Save(r, w)
		}
	}

	if r.Method == http.MethodGet {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{"message": "logged out"})
}

func MeHandler(w http.ResponseWriter, r *http.Request) {
	svc := GetService(r)
	if svc == nil {
		WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	session, _ := svc.Store.Get(r, SessionName)
	userID, username, ok := svc.GetUserFromSession(session)
	if !ok {
		WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "not authenticated"})
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"user": map[string]string{
			"id":       userID,
			"username": username,
			"email":    session.Values[EmailKey].(string),
		},
	})
}

func CSRFTokenHandler(w http.ResponseWriter, r *http.Request) {
	token := csrf.Token(r)
	WriteJSON(w, http.StatusOK, map[string]string{
		"csrfToken": token,
	})
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
