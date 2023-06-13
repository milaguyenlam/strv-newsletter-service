package repository

import (
	"cloud.google.com/go/firestore"
	"gorm.io/gorm"
)

// PostgresRepository is a struct that holds the database client for a PostgreSQL database.
// It serves as the base repository for all repositories interacting with a PostgreSQL database.
type PostgresRepository struct {
	db *gorm.DB // db is the GORM DB client
}

// FirebaseRepository is a struct that holds the client for a Firestore database.
// It serves as the base repository for all repositories interacting with a Firestore database.
type FirebaseRepository struct {
	client *firestore.Client // client is the Firestore client
}
