package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"strv.com/newsletter/repository"
)

func AuthMiddleware(ur *repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// get JWT from the header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Authorization header is required"})
			return
		}

		// validate JWT format
		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid authorization token format"})
			return
		}

		// validate JWT
		token, _ := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
			// check signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte("secret"), nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			_, err := ur.GetByEmail(claims["email"].(string))
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization token"})
				return
			}
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization token"})
			return
		}
	}
}
