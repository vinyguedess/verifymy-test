package services

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"verifymy-golang-test/entities"
	"verifymy-golang-test/models"
	"verifymy-golang-test/repositories"
)

type AuthService interface {
	SignUp(ctx context.Context, user models.User) (*entities.Credentials, error)
}

func NewAuthService(
	userRepository repositories.UserRepository,
) AuthService {
	return &authService{
		userRepository: userRepository,
	}
}

type authService struct {
	userRepository repositories.UserRepository
}

func (s *authService) SignUp(
	ctx context.Context, user models.User,
) (*entities.Credentials, error) {
	foundUser, err := s.userRepository.FindByEmail(ctx, user.Email)
	if err != nil {
		return nil, err
	} else if foundUser != nil {
		return nil, errors.New("e-mail is already in use")
	}

	hashedPassword, err := s.hashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	user.Password = hashedPassword
	signedUser, err := s.userRepository.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return s.getCredentialsFromUser(signedUser)
}

func (s *authService) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func (s *authService) getCredentialsFromUser(
	user *models.User,
) (*entities.Credentials, error) {
	expiresAt := time.Now().UTC().Add(time.Hour * 24).Unix()

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.String(),
		"exp":     expiresAt,
	})
	accessTokenString, err := accessToken.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return nil, err
	}

	return &entities.Credentials{
		AccessToken: accessTokenString,
		ExpiresAt:   expiresAt,
	}, nil
}
