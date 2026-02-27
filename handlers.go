package main

import (
	"fmt"
	"net/http"
)

func readinessHandler(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (cfg *apiConfig) countHandler(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	outHTML := fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", cfg.fileserverHits.Load())

	w.Write([]byte(outHTML))
}
