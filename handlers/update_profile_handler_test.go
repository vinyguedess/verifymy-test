package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	mock_services "verifymy-golang-test/mocks/services"
	"verifymy-golang-test/models"
)

type updateProfileHandlerTestSuite struct {
	suite.Suite
	ctrl            *gomock.Controller
	userServiceMock *mock_services.MockUserService
	handler         Handler
}

func TestUpdateProfileHandlerTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(updateProfileHandlerTestSuite))
}

func (s *updateProfileHandlerTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.userServiceMock = mock_services.NewMockUserService(s.ctrl)
	s.handler = NewUpdateProfileHandler(s.userServiceMock)
}

func (s *updateProfileHandlerTestSuite) TestMethod() {
	s.Equal([]string{"PUT"}, s.handler.Method())
}

func (s *updateProfileHandlerTestSuite) TestRoute() {
	s.Equal("/profile", s.handler.Route())
}

func (s *updateProfileHandlerTestSuite) TestServeHTTP() {
	tests := []struct {
		description        string
		payload            string
		expectedPayload    models.User
		updateProfileError error
		expectedStatusCode int
		invalidJsonError   bool
	}{
		{
			description:        "Success",
			payload:            `{"name": "new name"}`,
			expectedPayload:    models.User{Name: "new name"},
			expectedStatusCode: http.StatusNoContent,
		},
		{
			description:        "Invalid JSON",
			payload:            `{"name": "new name"`,
			expectedStatusCode: http.StatusUnprocessableEntity,
			invalidJsonError:   true,
		},
		{
			description:        "Unexpected error",
			payload:            `{"name": "new name"}`,
			expectedPayload:    models.User{Name: "new name"},
			updateProfileError: errors.New("unexpected error"),
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		s.Run(test.description, func() {
			request := httptest.NewRequest(
				http.MethodPut, "/profile", bytes.NewReader([]byte(test.payload)),
			)
			response := httptest.NewRecorder()

			if !test.invalidJsonError {
				s.userServiceMock.EXPECT().UpdateProfile(
					request.Context(), test.expectedPayload,
				).Return(test.updateProfileError)
			}

			s.handler.ServeHTTP(response, request)
			if test.invalidJsonError {
				s.Equal(
					response.Body.String(),
					`{"message":"Invalid JSON","details":["unexpected EOF"]}`,
				)
			} else if test.updateProfileError != nil {
				s.Equal(
					response.Body.String(),
					fmt.Sprintf(
						`{"message":"unexpected error","details":["%s"]}`,
						test.updateProfileError.Error(),
					),
				)
			}
			s.Equal(test.expectedStatusCode, response.Code)
		})
	}
}
