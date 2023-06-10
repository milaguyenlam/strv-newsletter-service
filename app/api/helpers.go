package api

import (
	"github.com/gin-gonic/gin"
	"strv.com/newsletter/middleware"
	"strv.com/newsletter/model"
)

func getCurrentUser(ctx *gin.Context) *model.User {
	userValue, exists := ctx.Get(middleware.UserContextKey)
	if !exists {
		return nil
	}
	user, ok := userValue.(*model.User)
	if !ok {
		return nil
	}
	return user
}

func createMessageResponse(message string) gin.H {
	return gin.H{"message": message}
}
