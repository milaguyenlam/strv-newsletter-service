package model

type CustomError struct {
	HttpCode int
	Message  string
}

func NewCustomError(httpCode int, message string) *CustomError {
	return &CustomError{
		HttpCode: httpCode,
		Message:  message,
	}
}

func (ce *CustomError) Error() string {
	return ce.Message
}

func (ce *CustomError) ToMessageResponse() MessageResponse {
	return NewMessageResponse(ce.Message)
}
