package main

import (
	"bytes"
	"io"
	"log"
	"maps"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	router "github.com/Mboukhal/SvGoPg/core"
	"github.com/Mboukhal/SvGoPg/core/auth"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
)

func init() {
	_ = godotenv.Load()
}

func main() {

	app_env := os.Getenv("APP_ENV")
	if app_env == "" {
		app_env = os.Getenv("NODE_ENV")
	}
	if app_env == "" {
		app_env = "production"
	}
	isProduction := app_env == "production"

	r := chi.NewRouter()

	cookieSecret := os.Getenv("COOKIE_SECRET")
	if cookieSecret == "" {
		cookieSecret = "dev-secret-change-in-production"
	}

	store := sessions.NewCookieStore([]byte(cookieSecret))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}

	svc := auth.NewService(store)

	// A good base middleware stack
	if isProduction {
		r.Use(middleware.RequestID)
		r.Use(middleware.RealIP)
		r.Use(middleware.Recoverer)
		r.Use(middleware.Logger)
	}
	r.Use(middleware.Logger)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Use(auth.PageAuth(svc))

	if isProduction {
		productionSettings(r)
	} else {
		developmentSettings(r)
	}

	router.RegisterRoutes(r, svc)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Starting server on port %s in %s mode", port, app_env)
	erro := http.ListenAndServe(":"+port, r)
	if erro != nil {
		log.Fatal(erro)
	}
}

func developmentSettings(r chi.Router) {

	// everything else → proxy to SvelteKit DEV
	r.NotFound(func(w http.ResponseWriter, req *http.Request) {

		if strings.HasPrefix(req.URL.String(), "/?token=") {
			http.NotFound(w, req)
			return
		}

		path := req.URL.Path

		ui := "http://localhost:1337"
		proxyReq, _ := http.NewRequest(req.Method, ui+path, req.Body)
		proxyReq.Header = req.Header

		resp, err := http.DefaultClient.Do(proxyReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		// go template chage for login page --------------------------------
		if path == "/login" {

			login_page := responseToString(resp)
			data := map[string]string{
				"Name": "Login Page 1",
			}
			var buf bytes.Buffer
			tmpl, err := template.New("login").Parse(login_page)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			err = tmpl.Execute(&buf, data)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write(buf.Bytes())
		}
		// ------------------------------------------------------

		maps.Copy(w.Header(), resp.Header)
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	})

}

func responseToString(resp *http.Response) string {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	return string(bodyBytes)
}

const (
	STATIC_DIR = "static"
)

func productionSettings(r chi.Router) {
	// Apply gzip middleware to all responses
	r.Use(middleware.Compress(5))

	// Serve static files from /_/{path...}
	r.Get("/_astro/*", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		path := req.URL.Path

		if strings.HasPrefix(path, "/_astro") {
			// Cache immutable assets for 1 hour
			w.Header().Set("Cache-Control", "public, max-age=3600, immutable")
		}

		http.FileServer(http.Dir(STATIC_DIR)).ServeHTTP(w, req)
	}))

	// SPA fallback - serve index.html for all other routes
	r.NotFound(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		cleanPath := filepath.Clean(req.URL.Path)
		// If the path is not root, check if the file exists
		if cleanPath != "/" {
			if f, err := os.Open(STATIC_DIR + cleanPath); err == nil {
				defer f.Close()
				stat, _ := f.Stat()
				if !stat.IsDir() {
					http.FileServer(http.Dir(STATIC_DIR)).ServeHTTP(w, req)
					return
				}
			}
		}

		// Fallback to index.html for SPA routing
		indexFile, err := os.Open(STATIC_DIR + "/index.html")
		if err != nil {
			http.Error(w, "index.html not found", http.StatusNotFound)
			return
		}
		defer indexFile.Close()

		info, _ := indexFile.Stat()
		content, _ := io.ReadAll(indexFile)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		// w.Header().Set("Cache-Control", "no-cache")
		http.ServeContent(w, req, "index.html", info.ModTime(), bytes.NewReader(content))
	}))
}
