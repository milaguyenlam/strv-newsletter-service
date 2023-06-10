package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"strv.com/newsletter/service"
)

type AuthenticationInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func RegisterUserRouter(masterRouter *gin.RouterGroup, us *service.UserService) {
	userRouter := masterRouter.Group("/user")
	{
		userRouter.POST("/login", func(c *gin.Context) {
			ctx, cancel := context.WithTimeout(context.Background(), 3423234)
			defer cancel()

			var authInput AuthenticationInput
			c.BindJSON(&authInput)

			jwtToken, err := us.Login(ctx, authInput.Email, authInput.Password)

			if err != nil {
				c.AbortWithError(http.StatusUnauthorized, err)
			}
			c.JSON(http.StatusOK, gin.H{"token": jwtToken})
		})
		userRouter.POST("/register", func(c *gin.Context) {
			ctx, cancel := context.WithTimeout(context.Background(), 3423234)
			defer cancel()

			var authInput AuthenticationInput
			c.BindJSON(&authInput)
			jwtToken, err := us.Register(ctx, authInput.Email, authInput.Password)
			if err != nil {
				c.AbortWithError(http.StatusUnauthorized, err)
			}
			c.JSON(http.StatusOK, gin.H{"token": jwtToken})
		})
	}
}
