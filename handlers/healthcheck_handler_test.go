package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
)

type healthCheckHandlerTestSuite struct {
	suite.Suite
	handler Handler
}

func TestHealthCheckHandlerTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(healthCheckHandlerTestSuite))
}

func (s *healthCheckHandlerTestSuite) SetupTest() {
	s.handler = NewHealthCheckHandler()
}

func (s *healthCheckHandlerTestSuite) TestMethod() {
	s.Equal([]string{"GET"}, s.handler.Method())
}

func (s *healthCheckHandlerTestSuite) TestRoute() {
	s.Equal("/", s.handler.Route())
}

func (s *healthCheckHandlerTestSuite) TestServeHTTP() {
	s.T().Setenv("SERVICE_NAME", "verifymy-test")
	s.T().Setenv("VERSION", "1.0.0")

	request := httptest.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()

	s.handler.ServeHTTP(response, request)

	var payload map[string]interface{}
	_ = json.NewDecoder(response.Body).Decode(&payload)

	s.Equal(
		map[string]interface{}{
			"service_name": "verifymy-test",
			"version":      "1.0.0",
		},
		payload,
	)
	s.Equal(http.StatusOK, response.Code)
}
