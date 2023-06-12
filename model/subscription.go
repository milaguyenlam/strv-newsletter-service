package model

type Subscription struct {
	Name             string          `firestore:"name"`
	EditorEmail      string          `firestore:"editorEmail"`
	Description      string          `firestore:"description"`
	SubscribedEmails map[string]bool `firestore:"subscribedEmails"`
}

func NewSubscription(name, editorEmail, description string) *Subscription {
	return &Subscription{
		Name:             name,
		EditorEmail:      editorEmail,
		Description:      description,
		SubscribedEmails: make(map[string]bool),
	}
}

func (s *Subscription) AddSubscribedEmail(email string) {
	if s.SubscribedEmails == nil {
		s.SubscribedEmails = make(map[string]bool)
	}
	s.SubscribedEmails[email] = true
}

func (s *Subscription) RemoveSubscribedEmail(email string) {
	if s.SubscribedEmails == nil {
		s.SubscribedEmails = make(map[string]bool)
	}
	delete(s.SubscribedEmails, email)
}

func (s *Subscription) GetSubscribedEmailsAsSlice() []*string {
	emails := make([]*string, 0, len(s.SubscribedEmails))
	for k, _ := range s.SubscribedEmails {
		emails = append(emails, &k)
	}
	return emails
}

func (s *Subscription) GetID() string {
	return s.Name + s.EditorEmail
}
