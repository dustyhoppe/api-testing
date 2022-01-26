package mocks

type Mock struct {
	// Request arguments
	path           string
	supportedVerbs []string
	requestHeaders map[string]string

	// Response
	statusCode      int
	latencyMs       int
	responseBody    interface{}
}

func NewMock(path string, supportedVerbs []string, requestHeaders map[string]string, statusCode, latencyMs int, responseBody interface{}) *Mock {

	return &Mock{
		path:           path,
		supportedVerbs: supportedVerbs,
		requestHeaders: requestHeaders,

		statusCode:      statusCode,
		latencyMs:       latencyMs,
		responseBody:    responseBody,
	}
}
