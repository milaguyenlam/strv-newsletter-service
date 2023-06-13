package model

// Subscription is a struct representing a newsletter subscription.
// It includes the name of the subscription, editor's email, description and a map of subscribed emails.
type Subscription struct {
	Name             string          `firestore:"name"`             // Name of the subscription
	EditorEmail      string          `firestore:"editorEmail"`      // Email of the editor of the subscription
	Description      string          `firestore:"description"`      // Description of the subscription
	SubscribedEmails map[string]bool `firestore:"subscribedEmails"` // Map of emails that are subscribed to this subscription
}

// NewSubscription is a function that creates a new subscription with given name, editor's email and description.
// It initializes an empty map for subscribed emails.
func NewSubscription(name, editorEmail, description string) *Subscription {
	return &Subscription{
		Name:             name,
		EditorEmail:      editorEmail,
		Description:      description,
		SubscribedEmails: make(map[string]bool),
	}
}

// AddSubscribedEmail is a method that adds a new email to the subscription's map of subscribed emails.
func (s *Subscription) AddSubscribedEmail(email string) {
	if s.SubscribedEmails == nil {
		s.SubscribedEmails = make(map[string]bool)
	}
	s.SubscribedEmails[email] = true
}

// RemoveSubscribedEmail is a method that removes a specified email from the subscription's map of subscribed emails.
func (s *Subscription) RemoveSubscribedEmail(email string) {
	if s.SubscribedEmails == nil {
		s.SubscribedEmails = make(map[string]bool)
	}
	delete(s.SubscribedEmails, email)
}

// GetSubscribedEmailsAsSlice is a method that returns all subscribed emails as a slice of strings.
func (s *Subscription) GetSubscribedEmailsAsSlice() []*string {
	emails := make([]*string, 0, len(s.SubscribedEmails))
	for k, _ := range s.SubscribedEmails {
		emails = append(emails, &k)
	}
	return emails
}

// GetID is a method that returns the unique identifier for a subscription, composed of its name and the editor's email.
func (s *Subscription) GetID() string {
	return s.Name + "_" + s.EditorEmail
}
