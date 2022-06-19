/*
A go app that surfaces NASA's Astronomy Picture of the Day.
The web server exposes three routes:
	/				- Serves up static docs site as homepage
	/apod		- Fetches and returns JSON from APOD API
	/image	- Returns raw image from today's APOD img URL

To deploy, see docs at: https://github.com/lissy93/go-apod
*/

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

/*
Copyright © 2022 Alicia Sykes <https://aliciasykes.com>

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the “Software”), to deal in
the Software without restriction, including without limitation the rights to use,
copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the
Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A
PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/
