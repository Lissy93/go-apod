package shared

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port               string `envconfig:"PORT" default:"8080"`
	CORSAllowedOrigins string `envconfig:"CORS_ALLOWED_ORIGINS" default:"*"`
	NASAAPIKey         string `envconfig:"NASA_API_KEY" required:"true"`
	NASABaseURL        string `envconfig:"NASA_BASE_URL" default:"https://api.nasa.gov/planetary/apod"`
}

func NewConfig() (*Config, error) {
	var conf Config
	err := envconfig.Process("apod", &conf)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}

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

func FetchApod(ctx context.Context, client *http.Client, conf *Config) (*ApodResponse, error) {
	url := fmt.Sprintf("%s?api_key=%s", conf.NASABaseURL, conf.NASAAPIKey)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	result := ApodResponse{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
