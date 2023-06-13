package model

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User is a struct representing a user in the database.
// It includes the ORM model, an email which is unique and not null, and the hashed password.
type User struct {
	gorm.Model            // Embedded GORM model providing basic fields
	Email          string `gorm:"type:varchar(100);not null; unique"` // User's unique email
	HashedPassword []byte `gorm:"type:varchar(100);not null;"`        // Hashed password of the user
}

// NewUser is a function that takes an email and a password, hashes the password, and returns a new User struct.
// If there's any error during these operations, it returns the error. Otherwise, it returns the user.
func NewUser(email string, password string) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &User{
		Email:          email,
		HashedPassword: hashedPassword,
	}, err
}

// VerifyPassword is a method that compares the provided password with the hashed password stored in the user struct.
// It returns an error if the password doesn't match, otherwise returns nil.
func (u *User) VerifyPassword(password string) error {
	return bcrypt.CompareHashAndPassword(u.HashedPassword, []byte(password))
}
