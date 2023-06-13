package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/dgrijalva/jwt-go"
	"strv.com/newsletter/model"
	"strv.com/newsletter/repository"
)

// UserService holds dependencies for managing user authentication and registration
type UserService struct {
	ur        *repository.UserRepository // user repository for database operations
	svc       *ses.SES                   // Amazon SES for email verification
	jwtSecret string                     // secret used to sign JWT tokens
}

// NewUserService creates a new UserService instance with provided dependencies
func NewUserService(ur *repository.UserRepository, svc *ses.SES, jwtSecret string) *UserService {
	return &UserService{
		ur:        ur,
		svc:       svc,
		jwtSecret: jwtSecret,
	}
}

// Login authenticates a user with their email and password, and returns a JWT token if successful
func (us *UserService) Login(ctx context.Context, email string, password string) (string, error) {
	// Retrieve the user from the repository
	foundUser, err := us.ur.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	// Verify the provided password
	err = foundUser.VerifyPassword(password)
	if err != nil {
		return "", err
	}

	// Create a new JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	})

	// Sign the token with the secret
	tokenString, err := token.SignedString([]byte(us.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Register creates a new user with the provided email and password, verifies the email, and returns a JWT token
func (us *UserService) Register(ctx context.Context, email string, password string) (string, error) {
	// Create a new user in the repository
	err := us.ur.CreateUser(ctx, email, password)
	if err != nil {
		return "", err
	}

	// Verify the email address with Amazon SES
	input := &ses.VerifyEmailAddressInput{
		EmailAddress: aws.String(email),
	}
	_, err = us.svc.VerifyEmailAddress(input)
	if err != nil {
		return "", err
	}

	// Log the user in to get a JWT token
	token, err := us.Login(ctx, email, password)
	if err != nil {
		return "", err
	}

	return token, nil
}

// GetByEmail retrieves a user by their email address
func (us *UserService) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	user, err := us.ur.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// VerifyToken verifies a JWT token and returns the associated user
func (us *UserService) VerifyToken(ctx context.Context, token string) (*model.User, error) {
	// Parse the JWT token
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Ensure the token uses the expected signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		// Return the secret for validation
		return []byte(us.jwtSecret), nil
	})

	// Validate the token and retrieve the claims
	var claims jwt.MapClaims
	var ok bool

	if claims, ok = parsedToken.Claims.(jwt.MapClaims); !ok || !parsedToken.Valid {
		return nil, errors.New("Invalid token")
	}

	user, err := us.GetByEmail(ctx, claims["email"].(string))

	if err != nil {
		return nil, fmt.Errorf("Retriving user from token: %w", err)
	}
	return user, nil
}
