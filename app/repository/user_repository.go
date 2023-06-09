package repository

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"strv.com/newsletter/model"
)

type UserRepository struct {
	PostgresRepository
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		PostgresRepository{DB: db},
	}
}

func (ur *UserRepository) CreateUser(email string, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user := &model.User{
		Email:          email,
		HashedPassword: hashedPassword,
	}
	return ur.DB.Create(user).Error
}

func (ur *UserRepository) GetByEmail(email string) (*model.User, error) {
	user := &model.User{}
	if err := ur.DB.Where("Email = ?", email).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}
