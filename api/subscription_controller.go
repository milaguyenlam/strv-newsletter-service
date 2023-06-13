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

// SubcriptionController is a struct that contains subscription services and user services.
type SubcriptionController struct {
	ss *service.SubscriptionService // subscription service instance
	us *service.UserService         // user service instance
}

// NewSubscriptionController is a constructor function that initializes a new SubscriptionController.
func NewSubscriptionController(ss *service.SubscriptionService, us *service.UserService) *SubcriptionController {
	return &SubcriptionController{
		ss: ss,
		us: us,
	}
}

// RegisterSubscriptionRouter is a method that sets up routes for subscription related requests.
func (sc *SubcriptionController) RegisterSubscriptionRouter(masterRouter *gin.RouterGroup) {
	subscriptionRouter := masterRouter.Group("/subscription") // create a new router group for subscription
	{
		// setup endpoints for create, send, subscribe and unsubscribe actions
		subscriptionRouter.POST("/create", middleware.CreateAuthMiddleware(sc.us, timeoutPeriod), sc.Create)
		subscriptionRouter.POST("/:subscriptionID/send", middleware.CreateAuthMiddleware(sc.us, timeoutPeriod), sc.Send)
		subscriptionRouter.GET("/:subscriptionID/subscribe", sc.Subscribe)
		subscriptionRouter.GET("/:subscriptionID/unsubscribe", sc.Unsubscribe)
	}
}

// Create a new subscription
// @Summary Create a new subscription
// @Description Create a new subscription with the given name and description
// @ID create-subscription
// @Security Bearer
// @Accept  json
// @Produce  json
// @Param   input body model.CreateSubscriptionInput true "Subscription input"
// @Success 200 {object} model.MessageResponse "Subscription ID - composed of its name and the editor's email divided by underscore."
// @Failure 500 {object} model.MessageResponse "Error message"
// @Router /subscription/create [post]
func (sc *SubcriptionController) Create(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutPeriod)
	defer cancel()

	var createSubscriptionInput model.CreateSubscriptionInput
	c.BindJSON(&createSubscriptionInput)

	currentUser, err := getCurrentUser(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewMessageResponse("Error while authenticating user"))
		return
	}

	subscriptionID, err := sc.ss.CreateSubscription(ctx, createSubscriptionInput.Name, currentUser.Email, createSubscriptionInput.Description)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewMessageResponse(fmt.Sprintf("Error while creating newsletter subscription: %v", err)))
		return
	}

	c.JSON(http.StatusOK, model.NewMessageResponse(fmt.Sprintf("Subscription created with id: %s", subscriptionID)))
}

// Send newsletter
// @Summary Send a newsletter
// @Description Send a newsletter to all subscribers
// @ID send-newsletter
// @Security Bearer
// @Accept  json
// @Produce  json
// @Param subscriptionID path string true "Subscription ID"
// @Param email body model.Email true "Email details"
// @Success 200 {object} model.MessageResponse "Message"
// @Failure 500 {object} model.MessageResponse "Error message"
// @Router /subscription/{subscriptionID}/send [post]
func (sc *SubcriptionController) Send(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutPeriod)
	defer cancel()

	var email model.Email
	c.BindJSON(&email)

	subcriptionID, subscriptionIDPresent := c.Params.Get("subscriptionID")
	if !subscriptionIDPresent {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewMessageResponse("Invalid request - subscriptionID has to be specified."))
		return
	}

	err := sc.ss.SendNewsletterEmail(ctx, subcriptionID, &email)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewMessageResponse(fmt.Sprintf("Sending newsletter failed: %s", err.Error())))
		return
	}

	c.JSON(http.StatusOK, model.NewMessageResponse("Email successfully sent!"))
}

// Subscribe to a newsletter
// @Summary Subscribe to a newsletter
// @Description Subscribe to a newsletter with a given subscription ID
// @ID subscribe-newsletter
// @Produce  json
// @Param subscriptionID path string true "Subscription ID"
// @Param email query string true "Email to subscribe"
// @Success 200 {object} model.MessageResponse "Message"
// @Failure 500 {object} model.MessageResponse "Error message"
// @Router /subscription/{subscriptionID}/subscribe [get]
func (sc *SubcriptionController) Subscribe(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutPeriod)
	defer cancel()

	email, emailPresent := c.GetQuery("email")
	subscriptionID, subscriptionIDPresent := c.Params.Get("subscriptionID")
	if !emailPresent || !subscriptionIDPresent {
		c.JSON(http.StatusBadRequest, model.NewMessageResponse("Invalid request - email and subscription ID has to be specified."))
		return
	}

	err := sc.ss.Subscribe(ctx, subscriptionID, email, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewMessageResponse(fmt.Sprintf("Subscription failed: %v", err.Error())))
		return
	}
	c.JSON(http.StatusOK, model.NewMessageResponse(fmt.Sprintf("%s successfully subscribed to %s", email, subscriptionID)))
}

// Unsubscribe from a newsletter
// @Summary Unsubscribe from a newsletter
// @Description Unsubscribe from a newsletter with a given subscription ID
// @ID unsubscribe-newsletter
// @Produce  json
// @Param subscriptionID path string true "Subscription ID"
// @Param email query string true "Email to unsubscribe"
// @Success 200 {object} model.MessageResponse "Message"
// @Failure 500 {object} model.MessageResponse "Error message"
// @Router /subscription/{subscriptionID}/unsubscribe [get]
func (sc *SubcriptionController) Unsubscribe(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutPeriod)
	defer cancel()

	email, emailPresent := c.GetQuery("email")
	subscriptionID, subscriptionIDPresent := c.Params.Get("subscriptionID")
	if !emailPresent || !subscriptionIDPresent {
		c.JSON(http.StatusBadRequest, model.NewMessageResponse("Invalid request - email and subscription ID has to be specified."))
		return
	}

	err := sc.ss.Unsubscribe(ctx, subscriptionID, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewMessageResponse(fmt.Sprintf("Unsubscription failed: %v", err.Error())))
		return
	}
	c.JSON(http.StatusOK, model.NewMessageResponse(fmt.Sprintf("%s successfully unsubscribed from %s", email, subscriptionID)))
}
