package repositories

import (
	"context"

	"gorm.io/gorm"

	"verifymy-golang-test/models"
)

type UserRepository interface {
	Create(context.Context, models.User) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindById(ctx context.Context, id string) (*models.User, error)
	UpdateAttributesByUserId(ctx context.Context, userId string, data models.User) error
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
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func (repo *userRepository) FindById(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	err := repo.db.WithContext(ctx).Where("id", id).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func (repo *userRepository) UpdateAttributesByUserId(
	ctx context.Context, userId string, data models.User,
) error {
	err := repo.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id", userId).
		Updates(data).
		Error
	if err != nil {
		return err
	}

	return nil
}
