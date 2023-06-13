package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	firebase "firebase.google.com/go"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"strv.com/newsletter/api"
	"strv.com/newsletter/config"
	"strv.com/newsletter/model"
	"strv.com/newsletter/repository"
	"strv.com/newsletter/service"
)

// @title           STRV Newsletter Subscription API
// @version         1.0
// @description     This is a newsletter subscription API service.

// @contact.name   Nguyen Thanh Lam
// @contact.url    https://github.com/milaguyenlam
// @contact.email  milaguyenlam@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080

// @SecurityDefinitions.apikey Bearer
// @in header
// @name Authorization

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

const gracefulShutdownTimeout = 5 * time.Second

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Loading configuration: %v", err)
	}

	// Connect to Postgres
	db, err := gorm.Open(postgres.Open(cfg.DatabaseDSN), &gorm.Config{Logger: logger.Default})
	if err != nil {
		log.Fatalf("Opening database: %v", err)
	}
	// Create User table defined in model.User
	err = db.AutoMigrate(&model.User{})
	if err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}

	// Create a Firebase App
	opt := option.WithCredentialsJSON(cfg.FirebaseJSON)
	firebaseApp, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("Initializing Firebase app: %v", err)
	}
	firestoreClient, err := firebaseApp.Firestore(context.Background())
	if err != nil {
		log.Fatalf("Getting Firestore client: %v", err)
	}

	// Create an AWS session
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(cfg.AwsRegion),
	})
	if err != nil {
		log.Fatalf("Creating AWS session: %v", err)
	}
	ses := ses.New(awsSession)

	// Create repositories
	userRepository := repository.NewUserRepository(db)
	subscriptionRepository := repository.NewSubscriptionRepository(firestoreClient)

	// Create services
	userService := service.NewUserService(userRepository, ses, cfg.JWTSecret)
	subscriptionService := service.NewSubscriptionService(subscriptionRepository, ses)

	// Create controllers
	userController := api.NewUserController(userService)
	subscriptionController := api.NewSubscriptionController(subscriptionService, userService)

	app := gin.Default()
	api.SetupRoutes(app, userController, subscriptionController)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.ServerPort),
		Handler: app,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen and serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Gracefully shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), gracefulShutdownTimeout)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server shut down.")
}
