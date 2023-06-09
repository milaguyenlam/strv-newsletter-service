package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"strv.com/newsletter/model"
	"strv.com/newsletter/service"
	"strv.com/newsletter/utils"
)

func CreateAuthMiddleware(us *service.UserService, timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		// get JWT from the header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, model.NewMessageResponse("Missing authorization token"))
			return
		}

		// validate JWT format
		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusBadRequest, model.NewMessageResponse(fmt.Sprintf("Invalid authorization token format: %s", authHeader)))
			return
		}

		// validate JWT
		user, err := us.VerifyToken(ctx, bearerToken[1])
		if err != nil {
			utils.AbortWithStatusJSONFromError(c, err)
			return
		}

		c.Set(service.UserContextKey, user)
		c.Next()
	}
}
