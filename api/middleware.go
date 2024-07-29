package api

import (
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func enableCORS(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// If the request method is OPTIONS, return a 200 status (pre-flight request)
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		handler.ServeHTTP(w, r)
	})
}

func attachMiddleware(router *mux.Router) {
	router.Use(enableCORS)

	// Handle cancelled requests
	router.Use(func(next http.Handler) http.Handler {
		return http.TimeoutHandler(next, 30*time.Second, "Request timed out")
	})
}
