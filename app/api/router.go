package api

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"strv.com/newsletter/docs"
)

const BasePath = "api/v1"

func SetupRoutes(app *gin.Engine, userController *UserController, subscriptionController *SubcriptionController) {
	masterRouter := app.Group(BasePath)
	setupSwagger(app)
	userController.RegisterUserRouter(masterRouter)
	subscriptionController.RegisterSubscriptionRouter(masterRouter)

}

func setupSwagger(app *gin.Engine) {
	docs.SwaggerInfo.BasePath = BasePath
	app.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
