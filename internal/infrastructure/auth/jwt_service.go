package auth

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/wonyus/backend-challenge/internal/domain/entities"
	domainErrors "github.com/wonyus/backend-challenge/internal/domain/errors"
	"github.com/wonyus/backend-challenge/internal/domain/repositories"
	"github.com/wonyus/backend-challenge/internal/domain/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type jwtService struct {
	secret   []byte
	userRepo repositories.UserRepository
}

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func NewJWTService(secret string, userRepo repositories.UserRepository) services.AuthService {
	return &jwtService{
		secret:   []byte(secret),
		userRepo: userRepo,
	}
}

func (s *jwtService) GenerateToken(ctx context.Context, user *entities.User) (string, error) {
	if string(s.secret) == "" {
		return "", domainErrors.ErrInvalidTokenSecret
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: user.ID.Hex(),
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

func (s *jwtService) ValidateToken(ctx context.Context, tokenString string) (*entities.User, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return s.secret, nil
	})

	if err != nil {
		return nil, domainErrors.ErrInvalidToken
	}

	if !token.Valid {
		return nil, domainErrors.ErrInvalidToken
	}

	userID, err := primitive.ObjectIDFromHex(claims.UserID)
	if err != nil {
		return nil, domainErrors.ErrInvalidToken
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, domainErrors.ErrUserNotFound
	}

	return user, nil
}

func (s *jwtService) HashPassword(password string) (string, error) {
	if password == "" {
		return "", domainErrors.ErrEmptyPassword
	}
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {

		return "", err
	}
	return string(hashedBytes), nil
}

func (s *jwtService) ComparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
