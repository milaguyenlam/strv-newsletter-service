package middleware

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Claims struct {
	Email string
	jwt.StandardClaims
}

var jwtKey = []byte("SecretKey")

func AuthMiddleware(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, nil)
		c.Abort()
		return
	}

	c.Next()
}
