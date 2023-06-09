package repository

import (
	"context"
	"log"

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

func (sr *SubscriptionRepository) Set(subscription *model.Subscription) (*firestore.WriteResult, error) {
	ctx := context.Background()
	documentId := subscription.Name + subscription.EditorEmail
	return sr.Client.Collection("subscriptions").Doc(documentId).Set(ctx, subscription)
}

func (sr *SubscriptionRepository) Get(name string, editorEmail string) (*model.Subscription, error) {
	ctx := context.Background()
	subscription := &model.Subscription{}
	documentId := name + editorEmail
	snapshot, err := sr.Client.Collection("subscriptions").Doc(documentId).Get(ctx)
	if err != nil {
		log.Println("")
		return nil, err
	}
	snapshot.DataTo(subscription)
	return subscription, err
}
