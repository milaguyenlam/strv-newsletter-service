package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"strv.com/newsletter/model"
)

func init() {
	if err := godotenv.Load("test.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func getURL(host, path string) string {
	return "https://" + host + "/api/v1" + path
}

func TestAPI(t *testing.T) {
	host := os.Getenv("HOST_URL")
	userEmail := os.Getenv("TEST_USER_EMAIL")
	userPassword := os.Getenv("TEST_USER_PASSWORD")
	subscribedEmail := os.Getenv("TEST_SUBSCRIBED_EMAIL")
	subscriptionName := os.Getenv("TEST_SUBSCRIPTION_NAME")
	subscriptionDesc := os.Getenv("TEST_SUBSCRIPTION_DESC")

	var jwtToken string

	t.Run("User Registration", func(t *testing.T) {
		requestBody, _ := json.Marshal(map[string]string{
			"email":    userEmail,
			"password": userPassword,
		})
		resp, err := http.Post(getURL(host, "/user/register"), "application/json", bytes.NewBuffer(requestBody))
		if err != nil {
			t.Fatalf("Could not make a request: %v", err)
		}
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		var authResponse model.AuthenticationResponse
		json.Unmarshal(bodyBytes, &authResponse)

		jwtToken = authResponse.Token
	})

	t.Run("User Login", func(t *testing.T) {
		requestBody, _ := json.Marshal(map[string]string{
			"email":    userEmail,
			"password": userPassword,
		})
		resp, err := http.Post(getURL(host, "/user/login"), "application/json", bytes.NewBuffer(requestBody))
		if err != nil {
			t.Fatalf("Could not make a request: %v", err)
		}
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		var authResponse model.AuthenticationResponse
		json.Unmarshal(bodyBytes, &authResponse)
		jwtToken = authResponse.Token
	})

	t.Run("Create a Subscription", func(t *testing.T) {
		requestBody, _ := json.Marshal(map[string]string{
			"description":      subscriptionDesc,
			"subscriptionName": subscriptionName,
		})
		req, _ := http.NewRequest("POST", getURL(host, "/subscription/create"), bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+jwtToken)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Could not make a request: %v", err)
		}
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		var messageResponse model.MessageResponse
		json.Unmarshal(bodyBytes, &messageResponse)
		assert.Equal(t, fmt.Sprintf("Subscription created with id: %s_%s", subscriptionName, userEmail), messageResponse.Message)
	})

	t.Run("Subscribe to a Newsletter", func(t *testing.T) {
		subscriprionID := subscriptionName + "_" + userEmail
		escapedSubscriptionID := url.QueryEscape(subscriprionID)
		req, _ := http.NewRequest("GET", getURL(host, fmt.Sprintf("/subscription/%s/subscribe?email=%s", escapedSubscriptionID, url.QueryEscape(subscribedEmail))), nil)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Could not make a request: %v", err)
		}
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		var messageResponse model.MessageResponse
		json.Unmarshal(bodyBytes, &messageResponse)
		assert.Equal(t, fmt.Sprintf("%s successfully subscribed to %s_%s", subscribedEmail, subscriptionName, userEmail), messageResponse.Message)
	})

	t.Run("Unubscribe Newsletter", func(t *testing.T) {
		subscriprionID := subscriptionName + "_" + userEmail
		escapedSubscriptionID := url.QueryEscape(subscriprionID)
		req, _ := http.NewRequest("GET", getURL(host, fmt.Sprintf("/subscription/%s/unsubscribe?email=%s", escapedSubscriptionID, url.QueryEscape(subscribedEmail))), nil)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Could not make a request: %v", err)
		}
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		var messageResponse model.MessageResponse
		json.Unmarshal(bodyBytes, &messageResponse)
		assert.Equal(t, fmt.Sprintf("%s successfully unsubscribed from %s_%s", subscribedEmail, subscriptionName, userEmail), messageResponse.Message)
	})
}
