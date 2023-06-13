package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"strv.com/newsletter/model"
	"strv.com/newsletter/service"
)

// UserController is a struct that contains user service.
type UserController struct {
	us *service.UserService // user service instance
}

// NewUserController is a constructor function that initializes a new UserController.
func NewUserController(us *service.UserService) *UserController {
	return &UserController{
		us: us,
	}
}

// RegisterUserRouter is a method that sets up routes for user related requests.
func (uc *UserController) RegisterUserRouter(masterRouter *gin.RouterGroup) {
	userRouter := masterRouter.Group("/user") // create a new router group for user
	{
		// setup endpoints for login and register actions
		userRouter.POST("/login", uc.Login)
		userRouter.POST("/register", uc.Register)
	}
}

// @Summary User Login
// @Description Logs in a user
// @ID login
// @Accept  json
// @Produce  json
// @Param authInput body model.AuthenticationInput true "Login credentials"
// @Success 200 {object} model.AuthenticationResponse "Token"
// @Failure 401 {object} model.MessageResponse "Message"
// @Router /user/login [post]
func (uc *UserController) Login(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutPeriod)
	defer cancel()

	var authInput model.AuthenticationInput
	c.BindJSON(&authInput)

	jwtToken, err := uc.us.Login(ctx, authInput.Email, authInput.Password)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewMessageResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": jwtToken})
}

// @Summary User Registration
// @Description Registers a new user
// @ID register
// @Accept  json
// @Produce  json
// @Param authInput body model.AuthenticationInput true "Registration details"
// @Success 200 {object} model.AuthenticationResponse "Token"
// @Failure 401 {object} model.MessageResponse "Message"
// @Router /user/register [post]
func (uc *UserController) Register(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutPeriod)
	defer cancel()

	var authInput model.AuthenticationInput
	c.BindJSON(&authInput)
	jwtToken, err := uc.us.Register(ctx, authInput.Email, authInput.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewMessageResponse(fmt.Sprintf("Registration failed: %v", err.Error())))
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": jwtToken})
}
