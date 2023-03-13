package middleware

import (
	"badsah/controllers"
	"badsah/models"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		client, err := controllers.GetMongoClient()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		username, password, ok := c.Request.BasicAuth()
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header"})
			c.Abort()
			return
		}

		var user models.User
		if err := client.Database("blog").Collection("users").FindOne(context.Background(), bson.M{"username": username, "password": password}).Decode(&user); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			c.Abort()
			return
		}

		if user.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient privileges"})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}
