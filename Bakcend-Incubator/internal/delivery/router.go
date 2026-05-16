package delivery

import (
	"net/http"

	httphandler "backend-incubator/internal/delivery/http"

	"github.com/gorilla/mux"
)

// NewRouter membuat router dengan semua endpoint terdaftar.
func NewRouter(h *httphandler.Handler) http.Handler {
	r := mux.NewRouter()

	// CORS middleware agar bisa diakses Vue dev server.
	r.Use(corsMiddleware)

	api := r.PathPrefix("/api").Subrouter()

	// GET — baca state terkini (untuk initial load atau polling fallback).
	api.HandleFunc("/state", h.HandleGetState).Methods(http.MethodGet)

	// GET — SSE stream untuk data realtime.
	api.HandleFunc("/stream", h.HandleSSE).Methods(http.MethodGet)

	// POST — action dari frontend.
	api.HandleFunc("/mode", h.HandleSetMode).Methods(http.MethodPost)
	api.HandleFunc("/lamp", h.HandleSetLamp).Methods(http.MethodPost)

	return r
}

// corsMiddleware menambahkan header CORS untuk dev environment.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Preflight request dari browser.
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
