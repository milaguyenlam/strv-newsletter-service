package repository

import (
	"context"

	"gorm.io/gorm"
	"strv.com/newsletter/model"
)

type UserRepository struct {
	PostgresRepository
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		PostgresRepository{db: db},
	}
}

func (ur *UserRepository) CreateUser(ctx context.Context, email string, password string) error {
	user, err := model.NewUser(email, password)
	if err != nil {
		return err
	}
	err = ur.db.WithContext(ctx).Create(user).Error
	if err != nil {
		return err
	}
	return nil
}

func (ur *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	user := &model.User{}
	err := ur.db.WithContext(ctx).Where("Email = ?", email).First(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}
