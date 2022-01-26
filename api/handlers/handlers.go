package handlers

import (
	"encoding/json"
	"github.com/happyreturns/api-testing/clients"
	"github.com/happyreturns/gohelpers/log"
	"net/http"
)

type MovieHandlers struct {
	apiClient *clients.MovieApiClient
	logger    *log.Logger
}

func NewMovieHandlers(movieApi *clients.MovieApiClient, logger *log.Logger) *MovieHandlers {
	return &MovieHandlers{
		apiClient: movieApi,
		logger:    logger,
	}
}

func (m *MovieHandlers) SearchMoviesHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()["q"]

	err, response := m.apiClient.SearchMovies(query[0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (m *MovieHandlers) CreateGuestSession(w http.ResponseWriter, r *http.Request) {

	// Create a guest session
	err, sessionResponse := m.apiClient.CreateGuestSession()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(sessionResponse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
