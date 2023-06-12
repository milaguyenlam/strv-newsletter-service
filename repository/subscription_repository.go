package repository

import (
	"context"

	"cloud.google.com/go/firestore"
	"strv.com/newsletter/model"
)

type SubscriptionRepository struct {
	client *firestore.Client
}

func NewSubscriptionRepository(client *firestore.Client) *SubscriptionRepository {
	return &SubscriptionRepository{client: client}
}

func (sr *SubscriptionRepository) Create(ctx context.Context, subscription *model.Subscription) (string, error) {
	subscriptionID := subscription.GetID()
	_, err := sr.client.Collection("subscriptions").Doc(subscriptionID).Create(ctx, subscription)
	if err != nil {
		return "", err
	}
	return subscriptionID, nil
}

func (sr *SubscriptionRepository) Set(ctx context.Context, documentId string, subscription *model.Subscription) error {
	_, err := sr.client.Collection("subscriptions").Doc(documentId).Set(ctx, subscription)
	if err != nil {
		return err
	}
	return nil
}

func (sr *SubscriptionRepository) Get(ctx context.Context, documentId string) (*model.Subscription, error) {
	subscription := &model.Subscription{}
	snapshot, err := sr.client.Collection("subscriptions").Doc(documentId).Get(ctx)
	if err != nil {
		return nil, err
	}
	err = snapshot.DataTo(subscription)
	if err != nil {
		return nil, err
	}
	return subscription, nil
}
