package service

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
	"strv.com/newsletter/model"
	"strv.com/newsletter/repository"
)

const sendEmailRetries = 10

// SubscriptionService is a structure to encapsulate subscriptions' related functionalities.
type SubscriptionService struct {
	sr  *repository.SubscriptionRepository // SubscriptionRepository to perform operations on the subscription data
	svc *ses.SES                           // AWS Simple Email Service client
}

// NewSubscriptionService creates and returns a new SubscriptionService.
func NewSubscriptionService(sr *repository.SubscriptionRepository, svc *ses.SES) *SubscriptionService {
	return &SubscriptionService{
		sr:  sr,
		svc: svc,
	}
}

// CreateSubscription creates a new subscription and stores it in the repository.
func (ss *SubscriptionService) CreateSubscription(ctx context.Context, name string, editorEmail string, description string) (string, error) {
	subscription := model.NewSubscription(name, editorEmail, description)
	subscriptionID, err := ss.sr.Create(ctx, subscription)
	return subscriptionID, err
}

// Subscribe adds an email to a subscription and sends a confirmation email.
func (ss *SubscriptionService) Subscribe(ctx context.Context, subscriptionID string, subscribedEmail string, unsubscribeLink string) error {
	subscription, err := ss.sr.Get(ctx, subscriptionID)
	if err != nil {
		return err
	}
	subscription.AddSubscribedEmail(subscribedEmail)
	err = ss.sr.Set(ctx, subscription)
	if err != nil {
		return err
	}
	err = ss.sendConfirmationEmail(ctx, subscription, subscribedEmail, unsubscribeLink)
	if err != nil {
		return err
	}
	return nil
}

// Unsubscribe removes an email from a subscription.
func (ss *SubscriptionService) Unsubscribe(ctx context.Context, subscriptionID string, subscribedEmail string) error {
	subscription, err := ss.sr.Get(ctx, subscriptionID)
	if err != nil {
		return err
	}
	subscription.RemoveSubscribedEmail(subscribedEmail)
	err = ss.sr.Set(ctx, subscription)
	if err != nil {
		return err
	}
	return nil
}

// SendNewsletterEmail sends a newsletter email to all subscribed emails.
func (ss *SubscriptionService) SendNewsletterEmail(ctx context.Context, subscriptionID string, email *model.Email) error {
	subscription, err := ss.sr.Get(ctx, subscriptionID)
	if err != nil {
		return err
	}
	err = ss.sendEmail(ctx, createSendEmailInput(subscription.GetSubscribedEmailsAsSlice(), subscription.EditorEmail, email.Subject, email.Body))
	if err != nil {
		return err
	}
	return nil
}

// sendConfirmationEmail sends a confirmation email for the subscription.
func (ss *SubscriptionService) sendConfirmationEmail(ctx context.Context, subscription *model.Subscription, subscribedEmail string, unsubscribeLink string) error {
	err := ss.sendEmail(ctx, createSendEmailInput([]*string{&subscribedEmail}, subscription.EditorEmail, fmt.Sprintf("Subscription confirmed: %s", subscription.Name), fmt.Sprintf("You've successfully subscribed to %s newsletter by %s\nDescription: %s\nUse this link to unsubscribe: %s", subscription.Name, subscription.EditorEmail, subscription.Description, unsubscribeLink)))
	if err != nil {
		return err
	}
	return nil
}

// sendEmail sends an email via AWS SES, with retry mechanism for failure.
func (ss *SubscriptionService) sendEmail(ctx context.Context, input *ses.SendEmailInput) (err error) {
	for i := 0; i < sendEmailRetries; i++ {
		_, err := ss.svc.SendEmailWithContext(ctx, input)
		if err == nil {
			return nil
		}
	}
	return err
}

// createSendEmailInput creates an input for SES SendEmail API.
func createSendEmailInput(toAddresses []*string, source string, subject string, body string) *ses.SendEmailInput {
	return &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: toAddresses,
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Text: &ses.Content{
					Data: aws.String(body),
				},
			},
			Subject: &ses.Content{
				Data: aws.String(subject),
			},
		},
		Source: aws.String(source),
	}
}
