package handler

import (
	"encoding/json"
	"net/http"

	"github.com/lissy93/go-apod/shared"
)

func Handler(w http.ResponseWriter, r *http.Request) {
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
