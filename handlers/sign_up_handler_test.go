package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"verifymy-golang-test/entities"
	mock_services "verifymy-golang-test/mocks/services"
	"verifymy-golang-test/models"
)

type signUpHandlerTestSuite struct {
	suite.Suite
	ctrl            *gomock.Controller
	authServiceMock *mock_services.MockAuthService
	handler         Handler
}

func TestSignUpHandlerTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(signUpHandlerTestSuite))
}

func (s *signUpHandlerTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.authServiceMock = mock_services.NewMockAuthService(s.ctrl)
	s.handler = NewSignUpHandler(s.authServiceMock)
}

func (s *signUpHandlerTestSuite) TestMethod() {
	s.Equal([]string{"POST"}, s.handler.Method())
}

func (s *signUpHandlerTestSuite) TestRoute() {
	s.Equal("/auth/sign_up", s.handler.Route())
}

func (s *signUpHandlerTestSuite) TestServeHTTP() {
	tests := []struct {
		description         string
		payload             string
		signUpResponse      *entities.Credentials
		signUpError         error
		expectedResponse    map[string]interface{}
		expectedStatusCode  int
		invalidPayloadError bool
	}{
		{
			description: "Success",
			payload: `{
				"name": "Bruce Wayne",
				"email": "bruce.wayne@jleague.io",
				"date_of_birth": "1939-05-01",
				"password": "lov3u4lfr3d",
				"address": "Gotham City"
			}`,
			signUpResponse: &entities.Credentials{
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
			description:         "Invalid JSON",
			payload:             `{"name": "Bruce Wayne"`,
			invalidPayloadError: true,
			expectedResponse: map[string]interface{}{
				"message": "Invalid JSON",
				"details": []interface{}{"unexpected EOF"},
			},
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			description: "Email already in use",
			payload: `{
				"name": "Bruce Wayne",
				"email": "bruce.wayne@jleague.io",
				"date_of_birth": "1939-05-01",
				"password": "lov3u4lfr3d",
				"address": "Gotham City"
			}`,
			signUpError: entities.NewEmailAlreadyInUseError("bruce.wayne@jleague.io"),
			expectedResponse: map[string]interface{}{
				"message": "e-mail is already in use",
				"details": []interface{}{"bruce.wayne@jleague.io"},
			},
			expectedStatusCode: http.StatusForbidden,
		},
		{
			description: "Unexpected error",
			payload: `{
				"name": "Bruce Wayne",
				"email": "bruce.wayne@jleague.io",
				"date_of_birth": "1939-05-01",
				"password": "lov3u4lfr3d",
				"address": "Gotham City"
			}`,
			signUpError: errors.New("unexpected error raised"),
			expectedResponse: map[string]interface{}{
				"message": "unexpected error",
				"details": []interface{}{"unexpected error raised"},
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		s.Run(test.description, func() {
			request := httptest.NewRequest(
				http.MethodPost,
				"/auth/sign_up",
				bytes.NewReader([]byte(test.payload)),
			)
			response := httptest.NewRecorder()

			if !test.invalidPayloadError {
				s.authServiceMock.EXPECT().SignUp(
					request.Context(),
					models.User{
						Name:  "Bruce Wayne",
						Email: "bruce.wayne@jleague.io",
						DateOfBirth: models.Date(
							time.Date(
								1939, time.May, 1, 0, 0, 0, 0, time.UTC,
							),
						),
						Password: "lov3u4lfr3d",
						Address:  "Gotham City",
					},
				).Return(
					test.signUpResponse, test.signUpError,
				)
			}

			s.handler.ServeHTTP(response, request)

			var payload map[string]interface{}
			_ = json.NewDecoder(response.Body).Decode(&payload)

			s.Equal(test.expectedResponse, payload)
			s.Equal(test.expectedStatusCode, response.Code)
		})
	}
}
