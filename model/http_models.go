package model

type MessageResponse struct {
	Message string `json:"message"`
}

func NewMessageResponse(message string) MessageResponse {
	return MessageResponse{Message: message}
}

type AuthenticationResponse struct {
	Token   string `json:"token"`
	Message string `json:"message"`
}

func NewAuthenticationResponse(token, message string) AuthenticationResponse {
	return AuthenticationResponse{
		Token:   token,
		Message: message,
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
