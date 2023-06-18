package services

import (
	"context"
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
	SignIn(ctx context.Context, email string, password string) (*entities.Credentials, error)
	GetUserFromToken(ctx context.Context, accessToken string) (*models.User, error)
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
		return nil, entities.NewEmailAlreadyInUseError(user.Email)
	}

	hashedPassword, err := s.hashPassword(string(user.Password))
	if err != nil {
		return nil, err
	}

	user.Password = models.SecretValue(hashedPassword)
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

func (s *authService) SignIn(
	ctx context.Context, email string, password string,
) (*entities.Credentials, error) {
	user, err := s.userRepository.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	} else if user == nil {
		return nil, entities.NewInvalidEmailAndOrPasswordError()
	}

	if err = s.compareHashAndPassword(string(user.Password), password); err != nil {
		return nil, entities.NewInvalidEmailAndOrPasswordError()
	}

	return s.getCredentialsFromUser(user)
}

func (s *authService) compareHashAndPassword(
	hashedPassword string, password string,
) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (s *authService) GetUserFromToken(
	ctx context.Context, token string,
) (*models.User, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		return nil, err
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, entities.NewInvalidTokenError()
	}

	user, err := s.userRepository.FindById(ctx, userID)
	if err != nil {
		return nil, err
	} else if user == nil {
		return nil, entities.NewInvalidTokenError()
	}

	return user, nil
}
