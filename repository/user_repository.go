package repository

import (
	"database/sql"

	"strv.com/newsletter/model"
)

type UserRepository struct {
	PostgresRepository
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		PostgresRepository{DB: db},
	}
}

func (ur *UserRepository) SignUp(email string, password string) (*model.User, error) {
	return nil, nil
}

func (ur *UserRepository) Login(email string, password string) (*model.User, error) {
	return nil, nil
}
