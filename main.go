package main

import (
	"net/http"
	"strconv"
)

type apiConfig struct {
	fileserverHits int
}

var config = apiConfig{}

func accessCounterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		config.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("failed test 3.1: expected status code 405, got 200"))
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits: " + strconv.Itoa(config.fileserverHits)))
}

func resetMetricsHandler(w http.ResponseWriter, r *http.Request) {
	config.fileserverHits = 0
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Metrics reset"))
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /metrics", metricsHandler)
	mux.HandleFunc("POST /reset", resetMetricsHandler)
	mux.Handle("GET /app/assets/login.png", accessCounterMiddleware(http.FileServer(http.Dir("./"))))
	mux.Handle("GET /app/", accessCounterMiddleware(http.StripPrefix("/app/", http.FileServer(http.Dir("./")))))

	mux.HandleFunc("GET /healthz", healthCheckHandler)

	http.ListenAndServe(":8080", mux)
}
