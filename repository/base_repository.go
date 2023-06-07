package repository

import (
	"database/sql"

	"cloud.google.com/go/firestore"
)

type PostgresRepository struct {
	DB *sql.DB
}

type FirebaseRepository struct {
	Client *firestore.Client
}
