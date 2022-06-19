package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

var client = http.Client{}

type ApodResponse struct {
	Copyright      string `json:"copyright,omitempty"`
	Date           string `json:"date,omitempty"`
	Explanation    string `json:"explanation,omitempty"`
	HdUrl          string `json:"hdurl,omitempty"`
	MediaType      string `json:"media_type,omitempty"`
	ServiceVersion string `json:"service_version,omitempty"`
	Title          string `json:"title,omitempty"`
	Url            string `json:"url,omitempty"`
}

// Returns either value of environmental variable with a given key, or fallback value
func get_env(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

// Makes request to APOD API, returns results and error
func make_request() ([]byte, error) {
	// Get API key, and form request URL
	apiKey := os.Getenv("NASA_API_KEY")
	baseUrl := "https://api.nasa.gov/planetary/apod"
	url := fmt.Sprintf("%s?api_key=%s", baseUrl, apiKey)
	if apiKey == "" || len(apiKey) < 40 {
		return nil, fmt.Errorf("A valid API key must be specified")
	}

	// Make request
	response, responseErr := http.Get(url)
	if responseErr != nil {
		return nil, responseErr
	}
	defer response.Body.Close()
	body, responseErr := ioutil.ReadAll(response.Body)

	// Return response body and/or response error
	return body, responseErr
}

// Converts string results into APOD object
func parse_request(response []byte, err error) (ApodResponse, error) {
	if err != nil {
		return ApodResponse{}, err
	}
	var result ApodResponse
	maybeParseError := json.Unmarshal(response, &result)
	if maybeParseError != nil {
		return ApodResponse{}, maybeParseError
	}
	return result, nil
}

// Gets the URL for todays image
func get_image_url() (string, error) {
	response, reqErr := make_request()
	if reqErr != nil {
		return "", reqErr
	}
	results, parseErr := parse_request(response, reqErr)
	if parseErr != nil {
		return "", parseErr
	}
	return string(results.Url), nil
}

// HTTP handler for the /image route
func image_handler(res http.ResponseWriter, req *http.Request) {
	imagePath, nasaError := get_image_url()
	if nasaError != nil {
		http.Error(res, "API Error: "+nasaError.Error(), 400)
	}
	reqImg, err := client.Get(imagePath)
	if err != nil {
		fmt.Fprintf(res, "Error %d", err)
		return
	}
	res.Header().Set("Content-Length", fmt.Sprint(reqImg.ContentLength))
	res.Header().Set("Content-Type", reqImg.Header.Get("Content-Type"))
	if _, err = io.Copy(res, reqImg.Body); err != nil {
		http.Error(res, "API Error: "+err.Error(), 400)
	}
	reqImg.Body.Close()
}

// Sets headers for API responses
func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Content-Type", "application/json")
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

// HTTP hander for the APOD info route
func apod_handler(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	results, reqErr := make_request()
	if reqErr != nil {
		http.Error(w, "API Error: "+reqErr.Error(), 400)
	}
	fmt.Fprintf(w, string(results))
}

// Start web server
func main() {
	port := ":" + get_env("PORT", "8080")
	http.HandleFunc("/apod", apod_handler)
	http.HandleFunc("/image", image_handler)
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.ListenAndServe(port, nil)
}
