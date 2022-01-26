package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	httptools "github.com/happyreturns/gohelpers/http"
)

type HttpTestClient struct {
	client *http.Client
	server *httptest.Server
}

func NewHttpTestClient(server *httptest.Server) *HttpTestClient {

	return &HttpTestClient{
		client: &http.Client{Timeout: 60 * time.Second},
		server: server,
	}
}

func (h *HttpTestClient) ExecuteGet(path string, response interface{}) error {
	buf := &bytes.Buffer{}
	url := fmt.Sprintf("%s/%s", h.server.URL, path)

	// Create GET request
	req, err := http.NewRequest("GET", url, buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")


	// Do GET request
	resp, err := httptools.CheckResponse(h.client.Do(req))
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

func (h *HttpTestClient) ExecutePost(path string, body interface{}, response interface{}) error {
	// Encode request body
	buf := bytes.Buffer{}
	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		return fmt.Errorf("failure encoding request body: %w", err)
	}

	// Create POST request
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", h.server.URL, path), &buf)
	if err != nil {
		return fmt.Errorf("unable to create http request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Do POST request
	resp, err := httptools.CheckResponse(h.client.Do(req))
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
