package repositories

import (
	"context"

	"gorm.io/gorm"

	"verifymy-golang-test/models"
)

type UserRepository interface {
	Create(context.Context, models.User) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

type userRepository struct {
	db *gorm.DB
}

func (repo *userRepository) Create(
	ctx context.Context, user models.User,
) (*models.User, error) {
	err := repo.db.WithContext(ctx).Create(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := repo.db.WithContext(ctx).Where("email", email).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}
