package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	mock_services "verifymy-golang-test/mocks/services"
	"verifymy-golang-test/models"
)

type listUsersHandlerTestSuite struct {
	suite.Suite
	ctrl            *gomock.Controller
	userServiceMock *mock_services.MockUserService
	handler         Handler
}

func TestNewListUsersHandler(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(listUsersHandlerTestSuite))
}

func (s *listUsersHandlerTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.userServiceMock = mock_services.NewMockUserService(s.ctrl)
	s.handler = NewListUsersHandler(s.userServiceMock)
}

func (s *listUsersHandlerTestSuite) TestMethod() {
	s.Equal([]string{"GET"}, s.handler.Method())
}

func (s *listUsersHandlerTestSuite) TestRoute() {
	s.Equal("/users", s.handler.Route())
}

func (s *listUsersHandlerTestSuite) TestServeHTTP() {
	users := []models.User{
		{
			ID: uuid.New(),
		},
	}

	tests := []struct {
		description        string
		queryString        map[string]string
		expectedLimit      int
		expectedPage       int
		findAllResponse    []models.User
		findAllCount       int64
		findAllError       error
		expectedResponse   interface{}
		expectedStatusCode int
	}{
		{
			description:     "Success with default query params",
			expectedLimit:   10,
			expectedPage:    1,
			findAllResponse: users,
			findAllCount:    int64(1),
			expectedResponse: []interface{}{
				map[string]interface{}{
					"id":            users[0].ID.String(),
					"name":          users[0].Name,
					"email":         users[0].Email,
					"password":      nil,
					"date_of_birth": time.Time{}.Format(models.DateFormat),
					"address":       users[0].Address,
				},
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			description: "Success setting page and limit",
			queryString: map[string]string{
				"limit": "3",
				"page":  "1",
			},
			expectedLimit:   3,
			expectedPage:    1,
			findAllResponse: users,
			findAllCount:    int64(1),
			expectedResponse: []interface{}{
				map[string]interface{}{
					"id":            users[0].ID.String(),
					"name":          users[0].Name,
					"email":         users[0].Email,
					"password":      nil,
					"date_of_birth": time.Time{}.Format(models.DateFormat),
					"address":       users[0].Address,
				},
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			description:   "Error fetching list of users",
			expectedLimit: 10,
			expectedPage:  1,
			findAllError:  errors.New("error fetching list of users"),
			expectedResponse: map[string]interface{}{
				"message": "unexpected error",
				"details": []interface{}{"error fetching list of users"},
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		s.Run(test.description, func() {
			request := httptest.NewRequest(http.MethodGet, "/users", nil)
			query := request.URL.Query()
			for key, value := range test.queryString {
				query.Add(key, value)
			}
			request.URL.RawQuery = query.Encode()

			response := httptest.NewRecorder()

			s.userServiceMock.EXPECT().FindAll(
				request.Context(),
				test.expectedLimit,
				test.expectedPage,
			).Return(test.findAllResponse, test.findAllCount, test.findAllError)

			s.handler.ServeHTTP(response, request)

			var payload interface{}
			_ = json.Unmarshal(response.Body.Bytes(), &payload)

			s.Equal(test.expectedResponse, payload)
			s.Equal(test.expectedStatusCode, response.Code)
			if test.findAllError == nil {
				s.Equal(fmt.Sprintf("%d", test.findAllCount), response.Header().Get("X-Total-Count"))

			}
		})
	}

}
