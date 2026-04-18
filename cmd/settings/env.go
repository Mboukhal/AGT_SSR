package settings

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "modernc.org/sqlite"
)

const (
	STATIC_DIR = "ui"
)

func DevelopmentSettings(r chi.Router) {

	// everything else → proxy to SvelteKit DEV
	r.NotFound(func(w http.ResponseWriter, req *http.Request) {

		if strings.HasPrefix(req.URL.String(), "/?token=") {
			http.NotFound(w, req)
			return
		}

		path := req.URL.Path

		ui := "http://localhost:" + os.Getenv("APP_PORT")
		proxyReq, _ := http.NewRequest(req.Method, ui+path, req.Body)
		proxyReq.Header = req.Header

		resp, err := http.DefaultClient.Do(proxyReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		for k, v := range resp.Header {
			w.Header()[k] = v
		}
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	})

	// r.NotFound(func(w http.ResponseWriter, req *http.Request) {

	// 	// skip API
	// 	if strings.HasPrefix(req.URL.Path, "/api") {
	// 		http.NotFound(w, req)
	// 		return
	// 	}

	// 	// only GET should fallback to SvelteKit
	// 	if req.Method != http.MethodGet {
	// 		http.NotFound(w, req)
	// 		return
	// 	}

	// 	ui := "http://localhost:" + os.Getenv("APP_PORT")

	// 	u, err := url.Parse(ui)

	// 	if err != nil {
	// 		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 			http.Error(w, "Invalid UI URL", http.StatusInternalServerError)
	// 		})
	// 		return
	// 	}
	// 	proxy := httputil.NewSingleHostReverseProxy(u)
	// 	// proxy to sveltekit
	// 	proxy.ServeHTTP(w, req)
	// })
}

func ProductionSettings(r chi.Router) {
	// Apply gzip middleware to all responses
	r.Use(middleware.Compress(5))

	// Serve static files from /_/{path...}
	r.Get("/_app/*", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		path := req.URL.Path

		if strings.HasPrefix(path, "/_app/immutable") {
			// Cache immutable assets for 1 hour
			w.Header().Set("Cache-Control", "public, max-age=3600, immutable")
		}

		http.FileServer(http.Dir(STATIC_DIR)).ServeHTTP(w, req)
	}))

	// SPA fallback - serve index.html for all other routes
	r.NotFound(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Try to serve the requested file if it exists

		path := req.URL.Path
		// Prevent serving files from .well-known directory
		// This is a security measure to avoid exposing sensitive files
		if strings.HasPrefix(path, "/.well-known/") {
			http.NotFound(w, req)
			return
		}

		cleanPath := filepath.Clean(req.URL.Path)
		// log.Println("cleanPath:", cleanPath)
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
		w.Header().Set("Cache-Control", "no-cache")
		http.ServeContent(w, req, "index.html", info.ModTime(), bytes.NewReader(content))
	}))
}
