package model

type MessageResponse struct {
	Message string `json:"message"`
}

func NewMessageResponse(message string) MessageResponse {
	return MessageResponse{Message: message}
}

type AuthenticationResponse struct {
	Token string `json:"token"`
}

func NewAuthenticationResponse(token string) AuthenticationResponse {
	return AuthenticationResponse{
		Token: token,
	}
}

type AuthenticationInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateSubscriptionInput struct {
	Name        string `json:"subscriptionName"`
	Description string `json:"description"`
}

type SendEmailInput struct {
	Email
}
