package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strings"
)

type apiConfig struct {
	FileserverHits int
}

var config = apiConfig{}

var profaneWords = map[string]bool{"kerfuffle": true, "sharbert": true, "fornax": true}

func cleanProfanity(sentence string) string {
	words := strings.Split(sentence, " ")
	for i, word := range words {
		if profaneWords[strings.ToLower(word)] {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	type validBody struct {
		Body string `json:"body"`
	}
	type invalidBody struct {
		Error string `json:"error"`
	}
	type resBody struct {
		CleanedBody string `json:"cleaned_body"`
	}
	var respBody validBody
	err := json.NewDecoder(r.Body).Decode(&respBody)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	if len(respBody.Body) > 140 {
		errorRes := invalidBody{Error: "Chirp is too long"}
		errorData, err := json.Marshal(errorRes)
		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorData)
		return
	}

	res := resBody{CleanedBody: cleanProfanity(respBody.Body)}
	response, err := json.Marshal(res)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	w.Write(response)
}

func accessCounterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		config.FileserverHits++
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
	tmpl, err := template.ParseFiles("pages/metrics.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := apiConfig{
		FileserverHits: config.FileserverHits,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func resetMetricsHandler(w http.ResponseWriter, r *http.Request) {
	config.FileserverHits = 0
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Metrics reset"))
}

func main() {
	mux := http.NewServeMux()

	mux.Handle("GET /app/", accessCounterMiddleware(http.StripPrefix("/app/", http.FileServer(http.Dir("./")))))
	mux.Handle("GET /app/assets/logo.png", accessCounterMiddleware(http.FileServer(http.Dir("./"))))
	mux.HandleFunc("GET /api/reset", resetMetricsHandler)
	mux.HandleFunc("POST /api/validate_chirp", validateChirpHandler)

	mux.HandleFunc("GET /admin/metrics", metricsHandler)
	mux.HandleFunc("GET /api/healthz", healthCheckHandler)
	http.ListenAndServe(":8080", mux)
}
