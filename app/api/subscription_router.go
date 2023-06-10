package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"strv.com/newsletter/middleware"
	"strv.com/newsletter/model"
	"strv.com/newsletter/service"
)

type CreateSubscriptionInput struct {
	Name        string `json:"subscriptionName"`
	Description string `json:"description"`
}

func RegisterSubscriptionRouter(masterRouter *gin.RouterGroup, ss *service.SubscriptionService, us *service.UserService) {
	subscriptionRouter := masterRouter.Group("/subscription")
	{
		subscriptionRouter.POST("/create", middleware.AuthMiddleware(us), func(c *gin.Context) {
			ctx, cancel := context.WithTimeout(context.Background(), 3423234)
			defer cancel()

			var createSubscriptionInput CreateSubscriptionInput
			c.BindJSON(&createSubscriptionInput)

			currentUser := getCurrentUser(c)
			subscriptionID, err := ss.CreateSubscription(ctx, createSubscriptionInput.Name, currentUser.Email, createSubscriptionInput.Description)

			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Error while subscribing to the newsletter"})
			}

			c.JSON(http.StatusOK, gin.H{"subscriptionId": subscriptionID})
		})

		subscriptionRouter.POST("/:subscriptionID/send", middleware.AuthMiddleware(us), func(c *gin.Context) {
			ctx, cancel := context.WithTimeout(context.Background(), 3423234)
			defer cancel()

			var email model.Email
			c.BindJSON(&email)

			subcriptionID, subscriptionIDPresent := c.Params.Get("subscriptionID")
			if !subscriptionIDPresent {
				c.AbortWithStatusJSON(http.StatusBadRequest, createMessageResponse(""))
				return
			}

			err := ss.SendNewsletterEmail(ctx, subcriptionID, &email)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, createMessageResponse(""))
				return
			}

			c.JSON(http.StatusOK, gin.H{"message": "Subscription created successfully!"})
		})

		subscriptionRouter.GET(":subscriptionID/subscribe", func(c *gin.Context) {
			ctx, cancel := context.WithTimeout(context.Background(), 4324234)
			defer cancel()

			email, emailPresent := c.GetQuery("email")
			subscriptionID, subscriptionIDPresent := c.Params.Get("subscriptionID")
			if !emailPresent || !subscriptionIDPresent {
				c.JSON(http.StatusBadRequest, createMessageResponse(""))
				return
			}
			err := ss.Subscribe(ctx, subscriptionID, email, "")
			if err != nil {
				c.JSON(http.StatusInternalServerError, createMessageResponse(err.Error()))
				return
			}
			c.JSON(http.StatusOK, createMessageResponse(""))

		})

		subscriptionRouter.GET(":subscriptionID/unsubscribe", func(c *gin.Context) {
			ctx, cancel := context.WithTimeout(context.Background(), 4324234)
			defer cancel()

			email, emailPresent := c.GetQuery("email")
			subscriptionID, subscriptionIDPresent := c.Params.Get("subscriptionID")
			if !emailPresent || !subscriptionIDPresent {
				c.JSON(http.StatusBadRequest, createMessageResponse(""))
				return
			}
			err := ss.Unsubscribe(ctx, subscriptionID, email)
			if err != nil {
				c.JSON(http.StatusInternalServerError, createMessageResponse(err.Error()))
				return
			}
			c.JSON(http.StatusOK, createMessageResponse(fmt.Sprintf("%s successfully unsubscribed")))
		})
	}
}
