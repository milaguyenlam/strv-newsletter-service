package repository

import (
	"context"

	"gorm.io/gorm"
	"strv.com/newsletter/model"
)

// UserRepository is a struct that provides methods for interaction with the Postgres database to manage users.
// It embeds the PostgresRepository which contains a reference to the database client.
type UserRepository struct {
	PostgresRepository // Embedded PostgresRepository
}

// NewUserRepository creates a new UserRepository with the provided GORM DB client.
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		PostgresRepository{db: db},
	}
}

// CreateUser creates a new user in the Postgres database.
// It takes an email and a password, creates a new user model, and then inserts it into the database.
// If there's any error during these operations, it returns the error. Otherwise, it returns nil.
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

// GetByEmail retrieves a user from the Postgres database using the provided email.
// It returns a pointer to the retrieved User and any error encountered during the operation.
func (ur *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	user := &model.User{}
	err := ur.db.WithContext(ctx).Where("Email = ?", email).First(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}
