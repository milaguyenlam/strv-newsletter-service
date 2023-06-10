package repository

import (
	"cloud.google.com/go/firestore"
	"gorm.io/gorm"
)

type PostgresRepository struct {
	db *gorm.DB
}

type FirebaseRepository struct {
	client *firestore.Client
}
