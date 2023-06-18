package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"

	mock_repositories "verifymy-golang-test/mocks"
	"verifymy-golang-test/models"
)

type authServiceTestSuite struct {
	suite.Suite
	ctrl               *gomock.Controller
	ctx                context.Context
	userRepositoryMock *mock_repositories.MockUserRepository
	authService        AuthService
}

func TestAuthService(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(authServiceTestSuite))
}

func (s *authServiceTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.ctx = context.Background()
	s.userRepositoryMock = mock_repositories.NewMockUserRepository(s.ctrl)
	s.authService = NewAuthService(s.userRepositoryMock)
}

func (s *authServiceTestSuite) TestSignUp() {
	s.T().Setenv("SECRET_KEY", "MY_SECRET_KEY")

	user := models.User{
		ID:    uuid.New(),
		Name:  "John Doe",
		Email: "john.doe@gmail.com",
		DateOfBirth: time.Now().UTC().Add(
			time.Hour * (24 * 365 * 18 * -1),
		),
		Password: "my-password",
		Address:  "Jl. Raya Bogor",
	}

	payload := models.User{
		Name:  "John Doe",
		Email: "john.doe@gmail.com",
		DateOfBirth: time.Now().UTC().Add(
			time.Hour * (24 * 365 * 18 * -1),
		),
		Password: "my-password",
		Address:  "Jl. Raya Bogor",
	}

	tests := []struct {
		description         string
		findByEmailResponse *models.User
		findByEmailError    error
		createUserResponse  *models.User
		createUserError     error
	}{
		{
			description:        "Success",
			createUserResponse: &user,
		},
		{
			description:      "Failed to fetch user by e-mail",
			findByEmailError: errors.New("failed to fetch user by e-mail"),
		},
		{
			description:         "E-mail is already in use",
			findByEmailError:    nil,
			findByEmailResponse: &user,
		},
		{
			description:     "Failed to create user",
			createUserError: errors.New("failed to create user"),
		},
	}

	for _, test := range tests {
		s.Run(test.description, func() {
			s.SetupTest()

			s.userRepositoryMock.EXPECT().FindByEmail(s.ctx, gomock.Any()).Return(
				test.findByEmailResponse, test.findByEmailError,
			)

			if test.findByEmailResponse == nil && test.findByEmailError == nil {
				s.userRepositoryMock.EXPECT().Create(s.ctx, gomock.Any()).Return(
					test.createUserResponse, test.createUserError,
				)
			}

			credentials, err := s.authService.SignUp(
				s.ctx,
				payload,
			)
			if test.findByEmailResponse != nil {
				s.NotNil(err)
				s.ErrorContains(err, "e-mail is already in use")
				s.Nil(credentials)
			} else if test.findByEmailError != nil {
				s.NotNil(err)
				s.ErrorContains(err, test.findByEmailError.Error())
				s.Nil(credentials)
			} else if test.createUserError != nil {
				s.NotNil(err)
				s.ErrorContains(err, test.createUserError.Error())
				s.Nil(credentials)
			} else {
				s.Nil(err)
				s.NotNil(credentials)
			}
		})
	}
}

func (s *authServiceTestSuite) TestSignIn() {
	password := "my-password"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := models.User{
		ID:    uuid.New(),
		Name:  "John Doe",
		Email: "john.doe@gmail.com",
		DateOfBirth: time.Now().UTC().Add(
			time.Hour * (24 * 365 * 18 * -1),
		),
		Password: string(hashedPassword),
		Address:  "Jl. Raya Bogor",
	}

	tests := []struct {
		description             string
		email                   string
		password                string
		findUserByEmailResponse *models.User
		findUserByEmailError    error
		invalidPasswordError    bool
	}{
		{
			description:             "Success",
			email:                   user.Email,
			password:                password,
			findUserByEmailResponse: &user,
		},
		{
			description:          "Failed to fetch user by e-mail",
			email:                user.Email,
			password:             password,
			findUserByEmailError: errors.New("failed to fetch user by e-mail"),
		},
		{
			description: "User not found",
			email:       user.Email,
			password:    password,
		},
		{
			description:             "Invalid password",
			email:                   user.Email,
			password:                "invalid-password",
			findUserByEmailResponse: &user,
			invalidPasswordError:    true,
		},
	}

	for _, test := range tests {
		s.Run(test.description, func() {
			s.SetupTest()

			s.userRepositoryMock.EXPECT().FindByEmail(s.ctx, test.email).Return(
				test.findUserByEmailResponse,
				test.findUserByEmailError,
			)

			credentials, err := s.authService.SignIn(
				s.ctx,
				test.email,
				test.password,
			)
			if test.findUserByEmailError != nil {
				s.NotNil(err)
				s.ErrorContains(err, test.findUserByEmailError.Error())
				s.Nil(credentials)
			} else if test.findUserByEmailResponse == nil || test.invalidPasswordError {
				s.NotNil(err)
				s.ErrorContains(err, "invalid e-mail and/or password")
				s.Nil(credentials)
			} else {
				s.NoError(err)
				s.NotNil(credentials)
			}
		})
	}
}