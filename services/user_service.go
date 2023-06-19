package services

import (
	"context"
	"database/sql"
	"time"

	"verifymy-golang-test/common"
	"verifymy-golang-test/entities"
	"verifymy-golang-test/models"
	"verifymy-golang-test/repositories"
	"verifymy-golang-test/utils"
)

type UserService interface {
	FindById(ctx context.Context, userId string) (*models.User, error)
	FindAll(ctx context.Context, limit int, page int) ([]models.User, int64, error)
	UpdateProfile(ctx context.Context, attributes models.User) error
	DeleteById(ctx context.Context, userId string) error
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

func (s *userService) FindAll(ctx context.Context, limit int, page int) ([]models.User, int64, error) {
	offset := (page - 1) * limit
	return s.userRepository.FindAll(ctx, limit, offset)
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

func (s *userService) DeleteById(ctx context.Context, userId string) error {
	return s.userRepository.UpdateAttributesByUserId(
		ctx,
		userId,
		models.User{
			DeletedAt: sql.NullTime{Time: time.Now().UTC()},
		},
	)
}
