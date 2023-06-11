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

type SubcriptionController struct {
	ss *service.SubscriptionService
	us *service.UserService
}

func NewSubscriptionController(ss *service.SubscriptionService, us *service.UserService) *SubcriptionController {
	return &SubcriptionController{
		ss: ss,
		us: us,
	}
}

func (sc *SubcriptionController) RegisterSubscriptionRouter(masterRouter *gin.RouterGroup) {
	subscriptionRouter := masterRouter.Group("/subscription")
	{
		subscriptionRouter.POST("/create", middleware.AuthMiddleware(sc.us), sc.Create)
		subscriptionRouter.POST("/:subscriptionID/send", middleware.AuthMiddleware(sc.us), sc.Send)
		subscriptionRouter.GET(":subscriptionID/subscribe", sc.Subscribe)
		subscriptionRouter.GET(":subscriptionID/unsubscribe", sc.Unsubscribe)
	}
}

// Create a new subscription
// @Summary Create a new subscription
// @Description Create a new subscription with the given name and description
// @ID create-subscription
// @Accept  json
// @Produce  json
// @Param   input body CreateSubscriptionInput true "subscription input"
// @Success 200 {object} string "subscriptionId"
// @Failure 500 {object} string "Error message"
// @Router /subscription [post]
func (sc *SubcriptionController) Create(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 3423234)
	defer cancel()

	var createSubscriptionInput CreateSubscriptionInput
	c.BindJSON(&createSubscriptionInput)

	currentUser := getCurrentUser(c)
	subscriptionID, err := sc.ss.CreateSubscription(ctx, createSubscriptionInput.Name, currentUser.Email, createSubscriptionInput.Description)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Error while subscribing to the newsletter"})
	}

	c.JSON(http.StatusOK, gin.H{"subscriptionId": subscriptionID})
}

// Send newsletter
// @Summary Send a newsletter
// @Description Send a newsletter to all subscribers
// @ID send-newsletter
// @Accept  json
// @Produce  json
// @Param subscriptionID path string true "Subscription ID"
// @Param email body model.Email true "Email details"
// @Success 200 {object} string "Message"
// @Failure 500 {object} string "Error message"
// @Router /subscription/{subscriptionID}/send [post]
func (sc *SubcriptionController) Send(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 3423234)
	defer cancel()

	var email model.Email
	c.BindJSON(&email)

	subcriptionID, subscriptionIDPresent := c.Params.Get("subscriptionID")
	if !subscriptionIDPresent {
		c.AbortWithStatusJSON(http.StatusBadRequest, createMessageResponse(""))
		return
	}

	err := sc.ss.SendNewsletterEmail(ctx, subcriptionID, &email)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, createMessageResponse(""))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subscription created successfully!"})
}

// Subscribe to a newsletter
// @Summary Subscribe to a newsletter
// @Description Subscribe to a newsletter with a given subscription ID
// @ID subscribe-newsletter
// @Produce  json
// @Param subscriptionID path string true "Subscription ID"
// @Param email query string true "Email to subscribe"
// @Success 200 {object} string "Message"
// @Failure 500 {object} string "Error message"
// @Router /subscription/{subscriptionID}/subscribe [get]
func (sc *SubcriptionController) Subscribe(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 4324234)
	defer cancel()

	email, emailPresent := c.GetQuery("email")
	subscriptionID, subscriptionIDPresent := c.Params.Get("subscriptionID")
	if !emailPresent || !subscriptionIDPresent {
		c.JSON(http.StatusBadRequest, createMessageResponse(""))
		return
	}

	err := sc.ss.Subscribe(ctx, subscriptionID, email, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, createMessageResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, createMessageResponse(""))
}

// Unsubscribe from a newsletter
// @Summary Unsubscribe from a newsletter
// @Description Unsubscribe from a newsletter with a given subscription ID
// @ID unsubscribe-newsletter
// @Produce  json
// @Param subscriptionID path string true "Subscription ID"
// @Param email query string true "Email to unsubscribe"
// @Success 200 {object} string "Message"
// @Failure 500 {object} string "Error message"
// @Router /subscription/{subscriptionID}/unsubscribe [get]
func (sc *SubcriptionController) Unsubscribe(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 4324234)
	defer cancel()

	email, emailPresent := c.GetQuery("email")
	subscriptionID, subscriptionIDPresent := c.Params.Get("subscriptionID")
	if !emailPresent || !subscriptionIDPresent {
		c.JSON(http.StatusBadRequest, createMessageResponse(""))
		return
	}

	err := sc.ss.Unsubscribe(ctx, subscriptionID, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createMessageResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, createMessageResponse(fmt.Sprintf("%s successfully unsubscribed")))
}
