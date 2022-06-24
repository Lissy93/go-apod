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

package main

import (
	"context"
	"embed"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/kelseyhightower/envconfig"
)

//go:embed static
var static embed.FS

func spaHandler(data embed.FS, root string) http.HandlerFunc {
	contentStatic, err := fs.Sub(fs.FS(data), root)
	if err != nil {
		log.Fatalf("failed to create static site content: %v", err)
	}

	fs := http.FileServer(http.FS(contentStatic))
	return func(w http.ResponseWriter, r *http.Request) {
		if _, err := data.ReadFile(root + r.RequestURI); os.IsNotExist(err) {
			http.StripPrefix(r.RequestURI, fs).ServeHTTP(w, r)
		} else {
			fs.ServeHTTP(w, r)
		}
	}
}

type config struct {
	Port               string `envconfig:"PORT" default:"8080"`
	CORSAllowedOrigins string `envconfig:"CORS_ALLOWED_ORIGINS" default:"*"`
	NASAAPIKey         string `envconfig:"NASA_API_KEY" required:"true"`
	NASABaseURL        string `envconfig:"NASA_BASE_URL" default:"https://api.nasa.gov/planetary/apod"`
}

type server struct {
	router *chi.Mux
	client *http.Client
	conf   *config
}

func newServer(conf *config) *server {
	return &server{
		router: chi.NewRouter(),
		client: &http.Client{Timeout: 60 * time.Second},
		conf:   conf,
	}
}

func (s *server) routes() *chi.Mux {
	s.router.Use(
		cors.New(cors.Options{
			AllowedOrigins: strings.Split(s.conf.CORSAllowedOrigins, ","),
			AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodOptions, http.MethodPut, http.MethodDelete},
			AllowedHeaders: []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		}).Handler,
	)
	s.router.Get("/image", s.handleImage())
	s.router.Get("/apod", s.handleApod())
	s.router.Get("/*", spaHandler(static, "static"))
	return s.router
}

type apodResponse struct {
	Copyright      string `json:"copyright,omitempty"`
	Date           string `json:"date,omitempty"`
	Explanation    string `json:"explanation,omitempty"`
	HdUrl          string `json:"hdurl,omitempty"`
	MediaType      string `json:"media_type,omitempty"`
	ServiceVersion string `json:"service_version,omitempty"`
	Title          string `json:"title,omitempty"`
	Url            string `json:"url,omitempty"`
}

// fetch the url of nasa apod
func (s *server) apod(ctx context.Context) (*apodResponse, error) {
	url := fmt.Sprintf("%s?api_key=%s", s.conf.NASABaseURL, s.conf.NASAAPIKey)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	result := apodResponse{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *server) handleImage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apod, err := s.apod(r.Context())
		if err != nil {
			log.Printf("apodUrl: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, apod.Url, nil)
		if err != nil {
			log.Printf("image: http.NewRequestWithContext: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		resp, err := s.client.Do(req)
		if err != nil {
			log.Printf("image: client.Do: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		w.Header().Set("Content-Length", fmt.Sprint(resp.ContentLength))
		w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
		_, err = io.Copy(w, resp.Body)
		if err != nil {
			log.Printf("image: io.Copy: %v", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}
	}
}

func (s *server) handleApod() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		results, reqErr := s.apod(r.Context())
		if reqErr != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
	}
}

// Start web server
func main() {
	var conf config
	err := envconfig.Process("apod", &conf)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\033[1;92m🌌 Go-APOD running at http://localhost:" + conf.Port + "/\033[0m")
	log.Fatal(http.ListenAndServe(":"+conf.Port, newServer(&conf).routes()))
}
