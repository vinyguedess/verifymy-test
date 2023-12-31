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
	FindAll(ctx context.Context, limit int, offset int) ([]models.User, int64, error)
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
	err := repo.db.WithContext(ctx).
		Where("email", email).
		Where("deleted_at IS NULL").
		First(&user).
		Error
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
	err := repo.db.WithContext(ctx).
		Where("id", id).
		Where("deleted_at IS NULL").
		First(&user).
		Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func (repo *userRepository) FindAll(
	ctx context.Context, limit int, offset int,
) ([]models.User, int64, error) {
	var users []models.User
	err := repo.db.WithContext(ctx).
		Where("deleted_at IS NULL").
		Limit(limit).
		Offset(offset).
		Find(&users).
		Error
	if err != nil {
		return nil, 0, err
	}

	var totalResults int64
	err = repo.db.WithContext(ctx).
		Where("deleted_at IS NULL").
		Model(&models.User{}).
		Count(&totalResults).
		Error
	if err != nil {
		return nil, 0, err
	}

	return users, totalResults, nil
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
