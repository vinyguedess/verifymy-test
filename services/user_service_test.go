package services

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"verifymy-golang-test/common"
	mock_repositories "verifymy-golang-test/mocks/repositories"
	"verifymy-golang-test/models"
)

type userServiceTestSuite struct {
	suite.Suite
	ctrl               *gomock.Controller
	userRepositoryMock *mock_repositories.MockUserRepository
	service            UserService
}

func TestUserServiceTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(userServiceTestSuite))
}

func (s *userServiceTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.userRepositoryMock = mock_repositories.NewMockUserRepository(s.ctrl)
	s.service = NewUserService(s.userRepositoryMock)
}

func (s *userServiceTestSuite) TestFindById() {
	userId := uuid.New()
	user := &models.User{
		ID: userId,
	}

	tests := []struct {
		description      string
		findByIdResponse *models.User
		findByIdError    error
		userNotFound     bool
	}{
		{
			description:      "Success",
			findByIdResponse: user,
		},
		{
			description:   "Error finding user by id",
			findByIdError: errors.New("error"),
		},
		{
			description:  "User not found",
			userNotFound: true,
		},
	}

	for _, test := range tests {
		s.Run(test.description, func() {
			ctx := context.Background()

			s.userRepositoryMock.EXPECT().FindById(ctx, userId.String()).Return(
				test.findByIdResponse, test.findByIdError,
			)

			foundUser, err := s.service.FindById(ctx, userId.String())
			if test.findByIdError != nil {
				s.Error(err)
				s.Nil(foundUser)
			} else if test.userNotFound {
				s.Error(err)
				s.ErrorContains(err, "User not found")
			} else {
				s.NoError(err)
				s.Equal(test.findByIdResponse, foundUser)
			}
		})
	}
}

func (s *userServiceTestSuite) TestUpdateProfile() {
	userId := uuid.New()
	user := &models.User{
		ID: userId,
	}

	tests := []struct {
		description                   string
		attributes                    models.User
		updatingPassword              bool
		updateAttributesByUserIdError error
	}{
		{
			description: "Success",
			attributes: models.User{
				Name: "John Doe",
			},
		},
		{
			description: "Success changing password",
			attributes: models.User{
				Password: "my-password",
			},
			updatingPassword: true,
		},
		{
			description: "Error updating attributes by user id",
			attributes: models.User{
				Name: "John Doe",
			},
			updateAttributesByUserIdError: errors.New("error"),
		},
	}

	for _, test := range tests {
		s.Run(test.description, func() {
			ctx := context.Background()
			ctx = context.WithValue(ctx, common.AuthUser, user)

			var attributesCopy interface{} = test.attributes
			if test.updatingPassword {
				attributesCopy = gomock.Any()
			}

			s.userRepositoryMock.EXPECT().UpdateAttributesByUserId(
				ctx, userId.String(), attributesCopy,
			).Return(test.updateAttributesByUserIdError)

			err := s.service.UpdateProfile(ctx, test.attributes)
			if test.updateAttributesByUserIdError != nil {
				s.Error(err)
			} else {
				s.NoError(err)
			}
		})
	}
}

func (s *userServiceTestSuite) TestDeleteById() {
	userId := uuid.New()
	ctx := context.Background()

	s.userRepositoryMock.EXPECT().UpdateAttributesByUserId(
		ctx, userId.String(), gomock.Any(),
	).Return(nil)

	err := s.service.DeleteById(ctx, userId.String())

	s.NoError(err)
}
