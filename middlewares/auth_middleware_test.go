package middlewares

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"verifymy-golang-test/common"
	mock_services "verifymy-golang-test/mocks/services"
	"verifymy-golang-test/models"
)

type authMiddlewareTestSuite struct {
	suite.Suite
	ctrl        *gomock.Controller
	authService *mock_services.MockAuthService
	middleware  func(http.Handler) http.Handler
	nextHandler http.Handler
}

func TestAuthMiddlewareTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(authMiddlewareTestSuite))
}

func (s *authMiddlewareTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.authService = mock_services.NewMockAuthService(s.ctrl)
	s.middleware = AuthMiddleware(s.authService)
	s.nextHandler = http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			user := r.Context().Value(common.AuthUser)

			jsonPayload, _ := json.Marshal(user)
			w.Write(jsonPayload)
		},
	)
}

func (s *authMiddlewareTestSuite) TestAuthMiddleware() {
	userId := uuid.New()
	user := &models.User{
		ID:   userId,
		Name: "Lebron James",
		DateOfBirth: models.Date(
			time.Date(1984, 12, 30, 0, 0, 0, 0, time.UTC),
		),
		Email:    "king.james@nba.com",
		Password: "password",
		Address:  "1111 S Figueroa St, Los Angeles",
	}

	test := []struct {
		description                      string
		route                            string
		authorizationHeader              string
		accessToken                      string
		expectedGetUserFromTokenResponse *models.User
		expectedGetUserFromTokenError    error
		expectedStatusCode               int
		expectedResponse                 map[string]interface{}
		isPublicURL                      bool
		hasAccessTokenProblem            bool
	}{
		{
			description:                      "Success",
			route:                            "/me",
			authorizationHeader:              "Bearer ACCESS_TOKEN",
			accessToken:                      "ACCESS_TOKEN",
			expectedGetUserFromTokenResponse: user,
			expectedGetUserFromTokenError:    nil,
			expectedStatusCode:               http.StatusOK,
			expectedResponse: map[string]interface{}{
				"id":            userId.String(),
				"name":          "Lebron James",
				"date_of_birth": "1984-12-30",
				"email":         "king.james@nba.com",
				"password":      nil,
				"address":       "1111 S Figueroa St, Los Angeles",
			},
		},
		{
			description:        "Success with a public URL",
			route:              "/",
			expectedStatusCode: http.StatusOK,
			expectedResponse:   nil,
			isPublicURL:        true,
		},
		{
			description:         "Malformed authorization header",
			route:               "/me",
			authorizationHeader: "Bearer",
			expectedStatusCode:  http.StatusUnauthorized,
			expectedResponse: map[string]interface{}{
				"message": "Malformed authorization header",
			},
			hasAccessTokenProblem: true,
		},
		{
			description:         "Invalid token type",
			route:               "/me",
			authorizationHeader: "Basic ACCESS_TOKEN",
			expectedStatusCode:  http.StatusUnauthorized,
			expectedResponse: map[string]interface{}{
				"message": "Authorization header must be a bearer token",
			},
			hasAccessTokenProblem: true,
		},
		{
			description:                   "Invalid access token",
			route:                         "/me",
			authorizationHeader:           "Bearer ACCESS_TOKEN",
			accessToken:                   "ACCESS_TOKEN",
			expectedGetUserFromTokenError: errors.New("invalid access token"),
			expectedStatusCode:            http.StatusUnauthorized,
			expectedResponse: map[string]interface{}{
				"message": "Invalid access token",
			},
		},
	}

	for _, test := range test {
		s.Run(test.description, func() {
			request := httptest.NewRequest("GET", test.route, nil)
			request.Header.Set("Authorization", test.authorizationHeader)

			response := httptest.NewRecorder()

			if !test.isPublicURL && !test.hasAccessTokenProblem {
				s.authService.EXPECT().GetUserFromToken(
					request.Context(), test.accessToken,
				).Return(
					test.expectedGetUserFromTokenResponse, test.expectedGetUserFromTokenError,
				)
			}

			midHandler := s.middleware(s.nextHandler)
			midHandler.ServeHTTP(response, request)

			var responseBody map[string]interface{}
			json.Unmarshal(response.Body.Bytes(), &responseBody)

			s.Equal(test.expectedResponse, responseBody)
			s.Equal(test.expectedStatusCode, response.Code)
		})
	}
}
