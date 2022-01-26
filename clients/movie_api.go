package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	conf2 "github.com/happyreturns/api-testing/conf"
	"github.com/happyreturns/api-testing/models"
	httptools "github.com/happyreturns/gohelpers/http"
	"github.com/happyreturns/gohelpers/log"
	"net/http"
	"net/url"
	"time"
)

type MovieApiClient struct {
	client  *http.Client
	logger  *log.Logger
	baseUrl string
	apiKey  string
}

func NewMovieApiClient(conf *conf2.Conf, logger *log.Logger) *MovieApiClient {
	return &MovieApiClient{
		baseUrl: conf.TMDBUrl,
		apiKey:  conf.TMDBApiKey,
		client:  &http.Client{Timeout: 3 * time.Second},
		logger:  logger,
	}
}

func (m *MovieApiClient) SearchMovies(query string) (error, *models.SearchMoviesResponse) {
	path := "3/search/movie"
	response := &models.SearchMoviesResponse{}

	params := url.Values{
		"query": []string{query},
	}

	err := m.get(path, params, response)
	if err != nil {
		return err, nil
	}

	return nil, response
}

func (m *MovieApiClient) CreateGuestSession() (error, *models.CreateGuestSessionResponse) {
	path := "3/authentication/guest_session/new"
	response := &models.CreateGuestSessionResponse{}

	params := url.Values{}

	err := m.get(path, params, response)
	if err != nil {
		return err, nil
	}

	return nil, response
}


func (m *MovieApiClient) get(path string, queryParams url.Values, response interface{}) error {

	queryParams.Add("api_key", m.apiKey)

	buf := &bytes.Buffer{}
	url := fmt.Sprintf("%s/%s?%s", m.baseUrl, path, queryParams.Encode())
	m.logger.Info(url)

	// Create GET request
	req, err := http.NewRequest("GET", url, buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	// Do GET request
	resp, err := httptools.CheckResponse(m.client.Do(req))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Parse response
	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return fmt.Errorf("unable to unmarshal movie db response: %w", err)
	}

	return nil
}

func (m *MovieApiClient) post(path string, queryParams url.Values, body interface{}, response interface{}) error {
	queryParams.Add("api_key", m.apiKey)

	// Encode request body
	buf := bytes.Buffer{}
	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		return fmt.Errorf("failure encoding request body: %w", err)
	}

	// Create POST request
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s?%s", m.baseUrl, path, queryParams.Encode()), &buf)
	if err != nil {
		return fmt.Errorf("unable to create http request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Do POST request
	resp, err := httptools.CheckResponse(m.client.Do(req))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Parse response
	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return fmt.Errorf("unable to unmarshal response: %w", err)
	}

	return nil
}
