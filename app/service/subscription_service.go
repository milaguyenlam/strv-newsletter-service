package service

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"strv.com/newsletter/model"
	"strv.com/newsletter/repository"
	"strv.com/newsletter/utils"
)

type SubscriptionService struct {
	SR *repository.SubscriptionRepository
}

func NewSubscriptionService(sr *repository.SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{
		SR: sr,
	}
}

func (ss *SubscriptionService) CreateSubscription(name string, editorEmail string, description string) error {
	subsription := &model.Subscription{
		Name:             name,
		EditorEmail:      editorEmail,
		Description:      description,
		SubscribedEmails: nil,
	}
	_, err := ss.SR.Set(subsription)
	return err
}

func (ss *SubscriptionService) Subscribe(subscribedEmail string, subscriptionName string, subscriptionEditorEmail string) error {
	subscription, err := ss.SR.Get(subscriptionName, subscriptionEditorEmail)
	if err != nil {
		return err
	}
	subscription.AddSubscribedEmail(subscribedEmail)
	_, err = ss.SR.Set(subscription)
	if err != nil {
		return err
	}
	//TODO: Send email with confirmation and a link to unsubscribe
	return err
}

func (ss *SubscriptionService) Unsubscribe(subscribedEmail string, subscriptionName string, subscriptionEditorEmail string) error {
	subscription, err := ss.SR.Get(subscriptionName, subscriptionEditorEmail)
	if err != nil {
		return err
	}
	subscription.RemoveSubscribedEmail(subscribedEmail)
	_, err = ss.SR.Set(subscription)
	return err
}

func (ss *SubscriptionService) SendEmail(subscriptionName string, editorEmail string, email model.Email) error {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("YourRegion"),
		Credentials: credentials.NewStaticCredentials("YourAccessKey", "YourSecretKey", ""),
	})
	if err != nil {
		return err
	}
	subscription, err := ss.SR.Get(subscriptionName, editorEmail)
	if err != nil {
		return err
	}
	toAddresses := utils.GetKeys(subscription.SubscribedEmails)
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: toAddresses,
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Text: &ses.Content{
					Data: aws.String(email.Body),
				},
			},
			Subject: &ses.Content{
				Data: aws.String(email.Subject),
			},
		},
		Source: aws.String(editorEmail),
	}
	svc := ses.New(sess)
	_, err = svc.SendEmail(input)
	//TODO: implement retries?
	return err
}
