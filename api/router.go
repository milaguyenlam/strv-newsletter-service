package api

import (
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const BasePath = "/api/v1"
const timeoutPeriod = 10 * time.Second //10s

// SetupRoutes is a function to setup all routes for the application.
// It takes gin engine instance, userController and subscriptionController as arguments.
func SetupRoutes(app *gin.Engine, userController *UserController, subscriptionController *SubcriptionController) {
	masterRouter := app.Group(BasePath)
	setupSwagger(app)                                               // setup Swagger UI
	userController.RegisterUserRouter(masterRouter)                 // register user routes
	subscriptionController.RegisterSubscriptionRouter(masterRouter) // register subscription routes
}

func setupSwagger(app *gin.Engine) {
	app.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
