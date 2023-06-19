package services

import (
	"context"

	"verifymy-golang-test/common"
	"verifymy-golang-test/entities"
	"verifymy-golang-test/models"
	"verifymy-golang-test/repositories"
	"verifymy-golang-test/utils"
)

type UserService interface {
	FindById(ctx context.Context, userId string) (*models.User, error)
	UpdateProfile(ctx context.Context, attributes models.User) error
}

type userService struct {
	userRepository repositories.UserRepository
}

func NewUserService(
	userRepository repositories.UserRepository,
) UserService {
	return &userService{
		userRepository: userRepository,
	}
}

func (s *userService) FindById(ctx context.Context, userId string) (*models.User, error) {
	user, err := s.userRepository.FindById(ctx, userId)
	if err != nil {
		return nil, err
	} else if user == nil {
		return nil, entities.NewItemNotFoundError("User", userId)
	}

	return user, nil
}

func (s *userService) UpdateProfile(ctx context.Context, attributes models.User) error {
	user := ctx.Value(common.AuthUser).(*models.User)
	if attributes.Password != "" {
		hashedPassword, err := utils.PasswordHash(string(attributes.Password))
		if err != nil {
			return err
		}

		attributes.Password = models.SecretValue(hashedPassword)
	}

	return s.userRepository.UpdateAttributesByUserId(ctx, user.ID.String(), attributes)
}
