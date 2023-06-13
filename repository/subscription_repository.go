package repository

import (
	"context"

	"cloud.google.com/go/firestore"
	"strv.com/newsletter/model"
)

const subscriptionsCollectionName = "subscriptions"

// SubscriptionRepository is a struct that allows interaction with the Firestore database to manage subscriptions.
// It embeds FirebaseRepository which contains a Firestore client.
type SubscriptionRepository struct {
	FirebaseRepository // This embedded struct contains the Firestore client
}

// NewSubscriptionRepository creates a new SubscriptionRepository with the provided firestore client.
func NewSubscriptionRepository(client *firestore.Client) *SubscriptionRepository {
	return &SubscriptionRepository{FirebaseRepository: FirebaseRepository{client: client}}
}

// Create adds a new subscription to the Firestore database.
// It returns the ID of the created subscription and any error encountered during the operation.
func (sr *SubscriptionRepository) Create(ctx context.Context, subscription *model.Subscription) (string, error) {
	subscriptionID := subscription.GetID()
	_, err := sr.client.Collection(subscriptionsCollectionName).Doc(subscriptionID).Create(ctx, subscription)
	if err != nil {
		return "", err
	}
	return subscriptionID, nil
}

// Set updates an existing subscription in the Firestore database.
// It returns any error encountered during the operation.
func (sr *SubscriptionRepository) Set(ctx context.Context, subscription *model.Subscription) error {
	_, err := sr.client.Collection(subscriptionsCollectionName).Doc(subscription.GetID()).Set(ctx, subscription)
	if err != nil {
		return err
	}
	return nil
}

// Get retrieves a subscription from the Firestore database by its documentId.
// It returns a pointer to the retrieved Subscription and any error encountered during the operation.
func (sr *SubscriptionRepository) Get(ctx context.Context, documentID string) (*model.Subscription, error) {
	subscription := &model.Subscription{}
	snapshot, err := sr.client.Collection(subscriptionsCollectionName).Doc(documentID).Get(ctx)
	if err != nil {
		return nil, err
	}
	err = snapshot.DataTo(subscription)
	if err != nil {
		return nil, err
	}
	return subscription, nil
}
