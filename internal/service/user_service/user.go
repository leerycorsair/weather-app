package userservice

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"time"
	"weather-app/internal/models"
	"weather-app/internal/repository"
	"weather-app/internal/service"

	"github.com/dgrijalva/jwt-go"
)

const (
	salt       = "qwerty123456"
	signingKey = "1234567890"
	tokenTTL   = 12 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"  db:"user_id"`
}

type UserService struct {
	cityService service.CityService
	userRep     repository.UserRepository
}

func NewUserService(cityService service.CityService, userRep repository.UserRepository) *UserService {
	return &UserService{
		cityService: cityService,
		userRep:     userRep,
	}
}

func (s *UserService) CreateUser(user models.User) (int, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.userRep.CreateUser(user)
}

func (s *UserService) GenerateToken(login, password string) (string, error) {
	user, err := s.userRep.GetUser(login, generatePasswordHash(password))
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix()},
		user.Id,
	})

	return token.SignedString([]byte(signingKey))
}

func (s *UserService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, nil
	}
	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type")
	}
	return claims.UserId, nil
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func (s *UserService) GetFavorites(userId int) ([]int, error) {
	return s.userRep.GetFavorites(userId)
}

func (s *UserService) AddFavorite(userId int, cityId int) (int, error) {
	return s.userRep.AddFavorite(userId, cityId)
}

func (s *UserService) DeleteFavorite(userId int, cityId int) error {
	return s.userRep.DeleteFavorite(userId, cityId)
}
