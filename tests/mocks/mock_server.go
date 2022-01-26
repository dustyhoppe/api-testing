package mocks

import (
	"encoding/json"
	"fmt"
	"github.com/happyreturns/gohelpers/log"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

type MockServer struct {
	server *httptest.Server
	logger *log.Logger
	mocks  []*Mock
}

func NewMockServer(logger *log.Logger) *MockServer {
	return &MockServer{
		logger: logger,
		mocks:  []*Mock{},
	}
}

func (m *MockServer) AddMock(mock *Mock) {
	m.mocks = append(m.mocks, mock)
}

func (m *MockServer) GetBaseUrl() string {
	if m.server == nil {
		return ""
	}

	return m.server.URL
}

// Start the mock server
func (m *MockServer) Start(t *testing.T) {
	m.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		m.logger.Info("Mock Server attempting to handle request")

		if len(m.mocks) == 0 {
			m.logger.Warn("No mocks configured for server")
			return
		}

		path := request.URL.RequestURI()
		verb := request.Method
		m.logger.Info(fmt.Sprintf("REQUEST - %s %s", verb, path))

		matchedMocks := make([]*Mock, 0)

		for _, mock := range m.mocks {
			if urlsMatch(path, mock.path) && hasVerb(mock.supportedVerbs, verb) {
				matchedMocks = append(matchedMocks, mock)
			}
		}

		if len(matchedMocks) == 0 {
			assert.Fail(t, "No matching mock routes for request")
			return
		}

		if len(matchedMocks) > 1 {
			assert.Fail(t, "Multiple matched routes. Ambiguous mock routes.")
			return
		}

		m.logger.Info("MOCK FOUND")

		targetMock := matchedMocks[0]
		if len(targetMock.requestHeaders) > 0 {
			for headerKey, headerValue := range targetMock.requestHeaders {
				if !hasHeaderWithValue(request.Header, headerKey, headerValue) {
					assert.Fail(t, fmt.Sprintf("Missing request header - %s:%s", headerKey, headerValue))
					return
				}
			}
		}

		if targetMock.latencyMs > 0 {
			duration, err := time.ParseDuration(fmt.Sprintf("%dms", targetMock.latencyMs))
			if err == nil {
				m.logger.Info(fmt.Sprintf("SLEEPING - %dms", targetMock.latencyMs))
				time.Sleep(duration)
			}
		}

		m.logger.Info(fmt.Sprintf("STATUS CODE - %d", targetMock.statusCode))

		w.WriteHeader(targetMock.statusCode)

		if err := json.NewEncoder(w).Encode(targetMock.responseBody); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}))
}

// Stop the mock server
func (m *MockServer) Stop() {
	if m.server != nil {
		m.server.Close()
	}
}

func urlsMatch(expected, actual string) bool {
	expectedUrl, _ := url.Parse(expected)
	actualUrl, _ := url.Parse(actual)

	expectedQuery := expectedUrl.Query()
	actualQuery := actualUrl.Query()

	if !strings.EqualFold(expectedUrl.String(), actualUrl.String()) {
		return false
	}

	if len(expectedQuery) != len(actualQuery) {
		return false
	}

	for key, _ := range expectedQuery {
		expectedValue := expectedQuery.Get(key)
		actualValue := actualQuery.Get(key)

		if !strings.EqualFold(expectedValue, actualValue) {
			return false
		}
	}

	return true
}

func hasVerb(allowed []string, actual string) bool {

	for _, verb := range allowed {
		if strings.EqualFold(verb, actual) {
			return true
		}
	}

	return false
}

func hasHeaderWithValue(header http.Header, key string, value string) bool {
	headerValue := header.Get(key)
	return headerValue == value
}
