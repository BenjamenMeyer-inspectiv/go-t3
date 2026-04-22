// Package api provides the HTTP router and handlers for go-t3.
package api

import "net/http"

// NewRouter returns an HTTP ServeMux with all routes registered.
func NewRouter() http.Handler {
	mux := http.NewServeMux()
	// Game routes will be added in Phase 2
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/ping", pingHandler)
	mux.HandleFunc("/pong", pongHandler)
	return mux
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"ping":"pong"}`))
}

func pongHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"pong":"ping"}`))
}
