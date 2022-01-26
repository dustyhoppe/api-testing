package tests

import (
	"fmt"
	api2 "github.com/happyreturns/api-testing/api"
	conf2 "github.com/happyreturns/api-testing/conf"
	"github.com/happyreturns/api-testing/models"
	"github.com/happyreturns/api-testing/tests/mocks"
	"github.com/happyreturns/gohelpers/log"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

const TMDB_API_KEY = "random-api-key"

func NewTestServer(logger *log.Logger, mockServer *mocks.MockServer) *httptest.Server {
	conf := &conf2.Conf{
		TMDBUrl:    mockServer.GetBaseUrl(),
		TMDBApiKey: TMDB_API_KEY,
	}
	api := api2.NewApi(conf, logger)
	api.Initialize()

	server := httptest.NewServer(api.Handler)

	return server
}

func Test_SearchMovies(t *testing.T) {

	logger := log.NewLogger("api-test", "api-test")

	mockServer := mocks.NewMockServer(logger)
	mockServer.Start(t)

	server := NewTestServer(logger, mockServer)
	client := NewHttpTestClient(server)

	defer server.Close()
	defer mockServer.Stop()

	t.Run("Search Movies", func(t *testing.T) {
		// Arrange
		query := "crazy"
		path := fmt.Sprintf("/3/search/movie?api_key=%s&query=%s", TMDB_API_KEY, query)
		supportedVerbs := []string{"GET"}
		requestHeaders := make(map[string]string)
		apiResponse := &models.SearchMoviesResponse{
			Page: 1,
			Results: []models.MovieSearchResult{
				{
					Id:            1,
					OriginalTitle: "Interstellar",
				},
				{
					Id:            2,
					OriginalTitle: "The Green Mile",
				},
			},
		}

		mock := mocks.NewMock(path, supportedVerbs, requestHeaders, 200, 250, apiResponse)
		mockServer.AddMock(mock)

		response := &models.SearchMoviesResponse{}

		// Act
		err := client.ExecuteGet(fmt.Sprintf("v1/movies/search?q=%s", query), response)

		// Assert
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Results))

		assert.Equal(t, 1, response.Results[0].Id)
		assert.Equal(t, "Interstellar", response.Results[0].OriginalTitle)

		assert.Equal(t, 2, response.Results[1].Id)
		assert.Equal(t, "The Green Mile", response.Results[1].OriginalTitle)

	})

	t.Run("Search Movies exceeds HTTP timeout", func(t *testing.T) {
		// Arrange
		query := "the+green"
		path := fmt.Sprintf("/3/search/movie?api_key=%s&query=%s", TMDB_API_KEY, query)
		supportedVerbs := []string{"GET"}
		requestHeaders := make(map[string]string)
		apiResponse := &models.SearchMoviesResponse{
			Page: 1,
			Results: []models.MovieSearchResult{
				{
					Id:            1,
					OriginalTitle: "Interstellar",
				},
				{
					Id:            2,
					OriginalTitle: "The Green Mile",
				},
			},
		}

		mock := mocks.NewMock(path, supportedVerbs, requestHeaders, 200, 500, apiResponse)
		mockServer.AddMock(mock)

		response := &models.SearchMoviesResponse{}

		// Act
		err := client.ExecuteGet(fmt.Sprintf("v1/movies/search?q=%s", query), response)

		// Assert
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Results))

		assert.Equal(t, 1, response.Results[0].Id)
		assert.Equal(t, "Interstellar", response.Results[0].OriginalTitle)

		assert.Equal(t, 2, response.Results[1].Id)
		assert.Equal(t, "The Green Mile", response.Results[1].OriginalTitle)

	})
}

func Test_CreateGuestSession(t *testing.T) {

	logger := log.NewLogger("api-test", "api-test")

	mockServer := mocks.NewMockServer(logger)
	mockServer.Start(t)

	server := NewTestServer(logger, mockServer)
	client := NewHttpTestClient(server)

	defer server.Close()
	defer mockServer.Stop()

	t.Run("Create Guest Session", func(t *testing.T) {
		// Arrange
		path := fmt.Sprintf("/3/authentication/guest_session/new?api_key=%s", TMDB_API_KEY)
		supportedVerbs := []string{"GET"}
		requestHeaders := make(map[string]string)
		apiResponse := &models.CreateGuestSessionResponse{
			ExpiresAt:      "expires-at",
			GuestSessionId: "abc123",
			Success:        true,
		}

		mock := mocks.NewMock(path, supportedVerbs, requestHeaders, 200, 250, apiResponse)
		mockServer.AddMock(mock)

		response := &models.CreateGuestSessionResponse{}

		// Act
		err := client.ExecutePost("v1/sessions", nil, response)

		// Assert
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, apiResponse.Success, response.Success)
		assert.Equal(t, apiResponse.ExpiresAt, response.ExpiresAt)
		assert.Equal(t, apiResponse.GuestSessionId, response.GuestSessionId)
	})
}
