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
	"strv.com/newsletter/api"
	"strv.com/newsletter/config"
	"strv.com/newsletter/repository"
	"strv.com/newsletter/service"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Connect to Postgres
	db, err := gorm.Open(postgres.Open(cfg.DatabaseDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	// Create a Firebase App
	opt := option.WithCredentialsFile(cfg.FirebaseCredentialsFile)
	firebaseApp, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("Error initializing Firebase app: %v", err)
	}
	firestoreClient, err := firebaseApp.Firestore(context.Background())
	if err != nil {
		log.Fatalf("Error getting Firestore client: %v", err)
	}

	// Create an AWS session
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(cfg.AwsRegion),
	})
	if err != nil {
		log.Fatalf("Error creating AWS session: %v", err)
	}

	// Create repositories
	userRepository := repository.NewUserRepository(db)
	subscriptionRepository := repository.NewSubscriptionRepository(firestoreClient)

	// Create services
	userService := service.NewUserService(userRepository)
	subscriptionService := service.NewSubscriptionService(subscriptionRepository, ses.New(awsSession))

	app := gin.Default()
	api.SetupRoutes(app, userService, subscriptionService)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.ServerPort),
		Handler: app,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen and serve: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Gracefully shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
