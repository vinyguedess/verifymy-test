package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"verifymy-golang-test/entities"
	mock_services "verifymy-golang-test/mocks/services"
)

type signInHandlerTestSuite struct {
	suite.Suite
	ctrl        *gomock.Controller
	authService *mock_services.MockAuthService
	handler     Handler
}

func TestSignInHandlerTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(signInHandlerTestSuite))
}

func (s *signInHandlerTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.authService = mock_services.NewMockAuthService(s.ctrl)
	s.handler = NewSignInHandler(s.authService)
}

func (s *signInHandlerTestSuite) TestMethod() {
	s.Equal([]string{"POST"}, s.handler.Method())
}

func (s *signInHandlerTestSuite) TestRoute() {
	s.Equal("/auth/sign_in", s.handler.Route())
}

func (s *signInHandlerTestSuite) TestServeHTTP() {
	tests := []struct {
		description         string
		payload             string
		signInResponse      *entities.Credentials
		signInError         error
		expectedResponse    map[string]interface{}
		expectedStatusCode  int
		invalidPayloadError bool
	}{
		{
			description: "Success",
			payload:     `{"email":"clark.kent@jleague.io","password":"lois_lane"}`,
			signInResponse: &entities.Credentials{
				AccessToken: "ACCESS_TOKEN",
				ExpiresAt:   1,
			},
			expectedResponse: map[string]interface{}{
				"access_token": "ACCESS_TOKEN",
				"expires_at":   float64(1),
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			description: "Invalid JSON",
			payload:     `{"email":"`,
			expectedResponse: map[string]interface{}{
				"message": "Invalid JSON",
				"details": []interface{}{"unexpected EOF"},
			},
			expectedStatusCode:  http.StatusUnprocessableEntity,
			invalidPayloadError: true,
		},
		{
			description: "Invalid email and/or password",
			payload:     `{"email":"clark.kent@jleague.io","password":"lois_lane"}`,
			signInError: entities.NewInvalidEmailAndOrPasswordError(),
			expectedResponse: map[string]interface{}{
				"message": "invalid e-mail and/or password",
				"details": nil,
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			description: "Unexpected error",
			payload:     `{"email":"clark.kent@jleague.io","password":"lois_lane"}`,
			signInError: errors.New("unexpected error was raised"),
			expectedResponse: map[string]interface{}{
				"message": "unexpected error",
				"details": []interface{}{"unexpected error was raised"},
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		s.Run(test.description, func() {
			request := httptest.NewRequest(
				"POST", "/auth/sign_in", bytes.NewReader([]byte(test.payload)),
			)
			response := httptest.NewRecorder()

			if !test.invalidPayloadError {
				s.authService.EXPECT().SignIn(
					request.Context(), "clark.kent@jleague.io", "lois_lane",
				).Return(test.signInResponse, test.signInError)
			}

			s.handler.ServeHTTP(response, request)

			var payload map[string]interface{}
			_ = json.NewDecoder(response.Body).Decode(&payload)

			s.Equal(test.expectedResponse, payload)
			s.Equal(test.expectedStatusCode, response.Code)
		})
	}
}
