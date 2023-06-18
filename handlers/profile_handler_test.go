package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"verifymy-golang-test/common"
	"verifymy-golang-test/models"
)

type profileHandlerTestSuite struct {
	suite.Suite
	ctrl    *gomock.Controller
	handler Handler
}

func TestProfileHandlerTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(profileHandlerTestSuite))
}

func (s *profileHandlerTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.handler = NewProfileHandler()
}

func (s *profileHandlerTestSuite) TestMethod() {
	s.Equal([]string{"GET"}, s.handler.Method())
}

func (s *profileHandlerTestSuite) TestRoute() {
	s.Equal("/profile", s.handler.Route())
}

func (s *profileHandlerTestSuite) TestServeHTTP() {
	userId := uuid.New()
	user := &models.User{
		ID:   userId,
		Name: "Peter Parker",
		DateOfBirth: models.Date(
			time.Now().UTC().AddDate(-20, 0, 0),
		),
		Email:    "peter.parker@nyork.co",
		Password: "sp00der",
		Address:  "20 Ingram Street",
	}

	request := httptest.NewRequest("GET", "/profile", nil)
	response := httptest.NewRecorder()

	s.handler.ServeHTTP(
		response,
		request.WithContext(
			context.WithValue(request.Context(), common.AuthUser, user),
		),
	)

	var jsonPayload map[string]interface{}
	_ = json.Unmarshal(response.Body.Bytes(), &jsonPayload)

	s.Equal(
		map[string]interface{}{
			"id":            userId.String(),
			"name":          "Peter Parker",
			"date_of_birth": time.Now().UTC().AddDate(-20, 0, 0).Format("2006-01-02"),
			"email":         "peter.parker@nyork.co",
			"password":      nil,
			"address":       "20 Ingram Street",
		},
		jsonPayload,
	)
	s.Equal(http.StatusOK, response.Code)
	s.Equal("application/json", response.Header().Get("Content-Type"))
}
