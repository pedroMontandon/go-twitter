package main

import (
	"net/http"
)

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", healthCheckHandler)
	mux.Handle("/app/assets/login.png", http.FileServer(http.Dir("./")))
	mux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir("./"))))

	http.ListenAndServe(":8080", mux)
}
