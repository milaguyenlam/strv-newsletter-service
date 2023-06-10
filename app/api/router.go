package api

import (
	"github.com/gin-gonic/gin"
	"strv.com/newsletter/service"
)

func SetupRoutes(app *gin.Engine, userService *service.UserService, subscriptionService *service.SubscriptionService) {
	masterRouter := app.Group("/api/v1")
	RegisterUserRouter(masterRouter, userService)
	RegisterSubscriptionRouter(masterRouter, subscriptionService, userService)
}
