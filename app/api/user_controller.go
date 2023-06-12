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

type AuthenticationResponse struct {
	Token string `json:"token"`
}

type UserController struct {
	us *service.UserService
}

func NewUserController(us *service.UserService) *UserController {
	return &UserController{
		us: us,
	}
}

func (uc *UserController) RegisterUserRouter(masterRouter *gin.RouterGroup) {
	userRouter := masterRouter.Group("/user")
	{
		userRouter.POST("/login", uc.Login)
		userRouter.POST("/register", uc.Register)
	}
}

// @Summary User Login
// @Description Logs in a user
// @ID login
// @Accept  json
// @Produce  json
// @Param authInput body AuthenticationInput true "Login credentials"
// @Success 200 {object} AuthenticationResponse "Token"
// @Failure 401 {object} MessageResponse "Message"
// @Router /user/login [post]
func (uc *UserController) Login(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), TimeoutPeriod)
	defer cancel()

	var authInput AuthenticationInput
	c.BindJSON(&authInput)

	jwtToken, err := uc.us.Login(ctx, authInput.Email, authInput.Password)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, createMessageResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": jwtToken})
}

// @Summary User Registration
// @Description Registers a new user
// @ID register
// @Accept  json
// @Produce  json
// @Param authInput body AuthenticationInput true "Registration details"
// @Success 200 {object} AuthenticationResponse "Token"
// @Failure 401 {object} MessageResponse "Message"
// @Router /user/register [post]
func (uc *UserController) Register(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), TimeoutPeriod)
	defer cancel()

	var authInput AuthenticationInput
	c.BindJSON(&authInput)
	jwtToken, err := uc.us.Register(ctx, authInput.Email, authInput.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, createMessageResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": jwtToken})
}
