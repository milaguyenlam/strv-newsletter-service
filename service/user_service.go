package service

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/dgrijalva/jwt-go"
	"strv.com/newsletter/model"
	"strv.com/newsletter/repository"
)

type UserService struct {
	ur  *repository.UserRepository
	svc *ses.SES
}

func NewUserService(ur *repository.UserRepository, svc *ses.SES) *UserService {
	return &UserService{
		ur:  ur,
		svc: svc,
	}
}

func (us *UserService) Login(ctx context.Context, email string, password string) (string, error) {
	foundUser, err := us.ur.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	err = foundUser.VerifyPassword(password)
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

func (us *UserService) Register(ctx context.Context, email string, password string) (string, error) {
	err := us.ur.CreateUser(ctx, email, password)
	if err != nil {
		return "", err
	}

	input := &ses.VerifyEmailAddressInput{
		EmailAddress: aws.String(email),
	}
	_, err = us.svc.VerifyEmailAddress(input)
	if err != nil {
		return "", err
	}

	token, err := us.Login(ctx, email, password)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (us *UserService) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	user, err := us.ur.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}
