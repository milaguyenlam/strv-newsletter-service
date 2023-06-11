package model

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email          string `gorm:"type:varchar(100);unique_index"`
	HashedPassword []byte `gorm:"type:varchar(100)"`
}

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

func (u *User) VerifyPassword(password string) error {
	return bcrypt.CompareHashAndPassword(u.HashedPassword, []byte(password))
}
