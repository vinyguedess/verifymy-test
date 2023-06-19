package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"

	"verifymy-golang-test/entities"
	mock_services "verifymy-golang-test/mocks/services"
	"verifymy-golang-test/models"
)

type deleteUserByIdHandlerTestSuite struct {
	suite.Suite
	ctrl            *gomock.Controller
	userServiceMock *mock_services.MockUserService
	handler         Handler
}

func TestDeleteUserByIdHandlerTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(deleteUserByIdHandlerTestSuite))
}

func (s *deleteUserByIdHandlerTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.userServiceMock = mock_services.NewMockUserService(s.ctrl)
	s.handler = NewDeleteUserByIdHandler(s.userServiceMock)
}

func (s *deleteUserByIdHandlerTestSuite) TestGetMethod() {
	s.Equal([]string{"DELETE"}, s.handler.Method())
}

func (s *deleteUserByIdHandlerTestSuite) TestRoute() {
	s.Equal("/users/{user_id}", s.handler.Route())
}

func (s *deleteUserByIdHandlerTestSuite) TestServeHTTP() {
	userId := uuid.New()
	user := &models.User{
		ID: userId,
	}

	tests := []struct {
		description        string
		findByIdResponse   *models.User
		findByIdError      error
		deleteByIdError    error
		expectedPayload    map[string]interface{}
		expectedStatusCode int
	}{
		{
			description:        "Success",
			findByIdResponse:   user,
			expectedStatusCode: http.StatusNoContent,
		},
		{
			description:   "User not found",
			findByIdError: entities.NewItemNotFoundError("User", userId.String()),
			expectedPayload: map[string]interface{}{
				"message": "User not found", "details": []interface{}{userId.String()},
			},
			expectedStatusCode: http.StatusNotFound,
		},
		{
			description:   "Unexpected error getting user by id",
			findByIdError: errors.New("unexpected error getting user by id"),
			expectedPayload: map[string]interface{}{
				"message": "unexpected error",
				"details": []interface{}{"unexpected error getting user by id"},
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			description:      "Unexpected error deleting user by id",
			findByIdResponse: user,
			deleteByIdError:  errors.New("unexpected error deleting user by id"),
			expectedPayload: map[string]interface{}{
				"message": "unexpected error",
				"details": []interface{}{"unexpected error deleting user by id"},
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		s.Run(test.description, func() {
			response := httptest.NewRecorder()
			request := httptest.NewRequest(
				http.MethodDelete,
				fmt.Sprintf("/users/%s", userId.String()),
				nil,
			)
			request = mux.SetURLVars(
				request, map[string]string{"user_id": userId.String()},
			)

			s.userServiceMock.EXPECT().FindById(
				request.Context(), userId.String(),
			).Return(test.findByIdResponse, test.findByIdError)

			if test.findByIdError == nil {
				s.userServiceMock.EXPECT().DeleteById(
					request.Context(), userId.String(),
				).Return(test.deleteByIdError)
			}

			s.handler.ServeHTTP(response, request)
			if test.expectedPayload != nil {
				var payload map[string]interface{}
				_ = json.NewDecoder(response.Body).Decode(&payload)

				s.Equal(test.expectedPayload, payload)
			}
			s.Equal(test.expectedStatusCode, response.Code)
		})
	}
}
