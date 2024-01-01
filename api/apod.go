package handler

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/lissy93/go-apod/shared"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	corsAllowedOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	if corsAllowedOrigins == "" {
		corsAllowedOrigins = "*" // Default to allowing all origins
	}

	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", corsAllowedOrigins)
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		return
	}

	conf, err := shared.NewConfig()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	client := &http.Client{}
	apodResponse, err := shared.FetchApod(r.Context(), client, conf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apodResponse)
}
