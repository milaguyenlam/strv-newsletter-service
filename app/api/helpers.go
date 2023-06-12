package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"strv.com/newsletter/middleware"
	"strv.com/newsletter/model"
)

type MessageResponse struct {
	Message string
}

func getCurrentUser(c *gin.Context) (*model.User, error) {
	userValue, exists := c.Get(middleware.UserContextKey)
	if !exists {
		return nil, fmt.Errorf("Getting current user from gin context failed: %v", c.Request.Header["Authorization"])
	}
	user, ok := userValue.(*model.User)
	if !ok {
		return nil, fmt.Errorf("Getting current user from gin context failed (type assertion failed): %v", c.Request.Header["Authorization"])
	}
	return user, nil
}

func createMessageResponse(message string) gin.H {
	return gin.H{"message": message}
}
