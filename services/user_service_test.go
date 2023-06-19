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

func (s *userServiceTestSuite) TestFindAll() {
	tests := []struct {
		description    string
		limit          int
		page           int
		expectedOffset int
	}{
		{
			description:    "Page 1 Limit 10",
			limit:          10,
			page:           1,
			expectedOffset: 0,
		},
		{
			description:    "Page 2 Limit 5",
			limit:          10,
			page:           2,
			expectedOffset: 10,
		},
		{
			description:    "Page 3 Limit 7",
			limit:          7,
			page:           3,
			expectedOffset: 14,
		},
	}

	for _, test := range tests {
		s.Run(test.description, func() {
			ctx := context.Background()

			s.userRepositoryMock.EXPECT().FindAll(ctx, test.limit, test.expectedOffset).Return(
				[]models.User{}, int64(0), nil,
			)

			_, _, err := s.service.FindAll(ctx, test.limit, test.page)
			s.NoError(err)
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
