package service

import (
	"context"
	"fmt"
	"log"
	"net/http"

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
func (ss *SubscriptionService) CreateSubscription(ctx context.Context, name string, editorEmail string, description string) (string, *model.CustomError) {
	subscription := model.NewSubscription(name, editorEmail, description)
	subscriptionID, err := ss.sr.Create(ctx, subscription)
	if err != nil {
		log.Println(err)
		return "", model.NewCustomError(http.StatusInternalServerError, fmt.Sprintf("Error occured while creating new subscription"))
	}

	return subscriptionID, nil
}

// Subscribe adds an email to a subscription and sends a confirmation email.
func (ss *SubscriptionService) Subscribe(ctx context.Context, subscriptionID string, subscribedEmail string, unsubscribeLink string) *model.CustomError {
	subscription, err := ss.sr.Get(ctx, subscriptionID)
	if err != nil {
		log.Println(err)
		return model.NewCustomError(http.StatusBadRequest, fmt.Sprintf("Error occured while fetching subscription with ID %s", subscriptionID))
	}
	subscription.AddSubscribedEmail(subscribedEmail)
	err = ss.sr.Set(ctx, subscription)
	if err != nil {
		log.Println(err)
		return model.NewCustomError(http.StatusBadRequest, fmt.Sprintf("Error occured while adding %s to %s", subscribedEmail, subscriptionID))
	}
	err = ss.sendConfirmationEmail(ctx, subscription, subscribedEmail, unsubscribeLink)
	if err != nil {
		log.Println(err)
		return model.NewCustomError(http.StatusInternalServerError, fmt.Sprintf("Error occured while sending confirmation email"))
	}
	return nil
}

// Unsubscribe removes an email from a subscription.
func (ss *SubscriptionService) Unsubscribe(ctx context.Context, subscriptionID string, subscribedEmail string) *model.CustomError {
	subscription, err := ss.sr.Get(ctx, subscriptionID)
	if err != nil {
		log.Println(err)
		return model.NewCustomError(http.StatusBadRequest, fmt.Sprintf("Error occured while fetching subscription with ID %s", subscriptionID))
	}
	subscription.RemoveSubscribedEmail(subscribedEmail)
	err = ss.sr.Set(ctx, subscription)
	if err != nil {
		log.Println(err)
		return model.NewCustomError(http.StatusInternalServerError, fmt.Sprintf("Error occured while removing %s from %s", subscribedEmail, subscriptionID))
	}
	return nil
}

// SendNewsletterEmail sends a newsletter email to all subscribed emails.
func (ss *SubscriptionService) SendNewsletterEmail(ctx context.Context, subscriptionID string, email *model.Email) *model.CustomError {
	subscription, err := ss.sr.Get(ctx, subscriptionID)
	if err != nil {
		log.Println(err)
		return model.NewCustomError(http.StatusBadRequest, fmt.Sprintf("Error occured while fetching subscription with ID %s", subscriptionID))
	}
	err = ss.sendEmail(ctx, createSendEmailInput(subscription.GetSubscribedEmailsAsSlice(), subscription.EditorEmail, email.Subject, email.Body))
	if err != nil {
		log.Println(err)
		return model.NewCustomError(http.StatusInternalServerError, "Error occured while sending newsletter emails!")
	}
	return nil
}

// sendConfirmationEmail sends a confirmation email for the subscription.
func (ss *SubscriptionService) sendConfirmationEmail(ctx context.Context, subscription *model.Subscription, subscribedEmail string, unsubscribeLink string) error {
	err := ss.sendEmail(ctx, createSendEmailInput([]*string{&subscribedEmail}, subscription.EditorEmail, fmt.Sprintf("Subscription confirmed: %s", subscription.Name), fmt.Sprintf("You've successfully subscribed to %s newsletter by %s\nDescription: %s\nUse this link to unsubscribe: %s", subscription.Name, subscription.EditorEmail, subscription.Description, unsubscribeLink)))
	if err != nil {
		return fmt.Errorf("Sending confirmation email: %w", err)
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
	return fmt.Errorf("Sending emails from %s to %#v: ", *input.Source, input.Destination.ToAddresses)
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
