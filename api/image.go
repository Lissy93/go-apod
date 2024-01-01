package handler

import (
	"io"
	"net/http"

	"github.com/lissy93/go-apod/shared" // Ensure this import path is correct
)

func Handler(w http.ResponseWriter, r *http.Request) {
	conf, err := shared.NewConfig()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	client := &http.Client{}
	apodData, err := shared.FetchApod(r.Context(), client, conf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Fetch the image from the URL in the APOD data
	imageResp, err := http.Get(apodData.Url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer imageResp.Body.Close()

	// Copy the content type and content from the response
	w.Header().Set("Content-Type", imageResp.Header.Get("Content-Type"))
	w.WriteHeader(imageResp.StatusCode)
	_, err = io.Copy(w, imageResp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
