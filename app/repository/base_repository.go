package repository

import (
	"cloud.google.com/go/firestore"
	"gorm.io/gorm"
)

type PostgresRepository struct {
	DB *gorm.DB
}

type FirebaseRepository struct {
	Client *firestore.Client
}
