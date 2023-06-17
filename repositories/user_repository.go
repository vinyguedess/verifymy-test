package repositories

import (
	"context"

	"gorm.io/gorm"

	"verifymy-golang-test/models"
)

type UserRepository interface {
	Create(context.Context, models.User) (*models.User, error)
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

type userRepository struct {
	db *gorm.DB
}

func (s *userRepository) Create(
	ctx context.Context, user models.User,
) (*models.User, error) {
	err := s.db.WithContext(ctx).Create(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}
