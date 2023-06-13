package utils

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"strv.com/newsletter/model"
	"strv.com/newsletter/service"
)

func GetCurrentUser(c *gin.Context) (*model.User, *model.CustomError) {
	userValue, exists := c.Get(service.UserContextKey)
	if !exists {
		return nil, model.NewCustomError(http.StatusInternalServerError, fmt.Sprintf("Getting current user from gin context failed: %v", c.Request.Header["Authorization"]))
	}
	user, ok := userValue.(*model.User)
	if !ok {
		return nil, model.NewCustomError(http.StatusInternalServerError, fmt.Sprintf("Getting current user from gin context failed (type assertion failed): %v", c.Request.Header["Authorization"]))
	}
	return user, nil
}

func AbortWithStatusJSONFromError(c *gin.Context, err *model.CustomError) {
	c.AbortWithStatusJSON(err.HttpCode, err.ToMessageResponse())
}

func CreateUnsubscribeLink(subscribeRequest *http.Request) string {
	protocol := "http://"
	if subscribeRequest.TLS != nil {
		protocol = "https://"
	}
	fullURL := protocol + subscribeRequest.Host + subscribeRequest.RequestURI
	return strings.Replace(fullURL, "/subscribe?", "/unsubscribe?", 1)
}
