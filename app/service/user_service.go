package service

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"strv.com/newsletter/repository"
)

type UserService struct {
	UR *repository.UserRepository
}

func NewUserService(ur *repository.UserRepository) *UserService {
	return &UserService{
		UR: ur,
	}
}

func (us *UserService) Login(c *gin.Context, email string, password string) (string, error) {
	foundUser, err := us.UR.GetByEmail(email)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword(foundUser.HashedPassword, []byte(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (us *UserService) Register(c *gin.Context, email string, password string) (string, error) {
	err := us.UR.CreateUser(email, password)
	if err != nil {
		return "", err
	}
	return us.Login(c, email, password)
}
