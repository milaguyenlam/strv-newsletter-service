package model

type Subscription struct {
	Name             string          `firestore:"name"`
	EditorEmail      string          `firestore:"editorEmail"`
	Description      string          `firestore:"description"`
	SubscribedEmails map[string]bool `firestore:"subscribedEmails"`
}

func (s *Subscription) AddSubscribedEmail(email string) {
	s.SubscribedEmails[email] = true
}

func (s *Subscription) RemoveSubscribedEmail(email string) {
	delete(s.SubscribedEmails, email)
}
