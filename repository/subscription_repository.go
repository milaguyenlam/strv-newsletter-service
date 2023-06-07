package repository

import (
	"context"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
	"strv.com/newsletter/model"
)

type SubscriptionRepository struct {
	Client *firestore.Client
}

func NewSubscriptionRepository() (*SubscriptionRepository, error) {
	ctx := context.Background()
	app, err := firebase.NewApp(ctx, nil, option.WithCredentialsFile("path/to/serviceAccountKey.json"))
	if err != nil {
		return nil, err
	}
	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, err
	}
	return &SubscriptionRepository{Client: client}, nil
}

func (sr *SubscriptionRepository) Create(subscription *model.Subscription) (*firestore.DocumentRef, *firestore.WriteResult, error) {
	ctx := context.Background()
	return sr.Client.Collection("subscriptions").Add(ctx, subscription)
}

func (sr *SubscriptionRepository) Remove()
