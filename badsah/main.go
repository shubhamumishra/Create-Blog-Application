package main

import (
	"badsah/models"
	"badsah/middleware"
	"context"
	"log"
	"badsah/controllers"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

func initialize() {
	client, err := controllers.GetMongoClient()
	if err != nil {
		log.Fatalf("Error initializing MongoDB client: %v", err)
	}
	userCollection := client.Database("blog").Collection("users")
	if count, _ := userCollection.CountDocuments(context.Background(), bson.M{}); count == 0 {
		_, err := userCollection.InsertOne(context.Background(), models.User{
			Username: "admin",
			Password: "password",
			Role:     "admin",
		})
		if err != nil {
			log.Fatalf("Error seeding admin user: %v", err)
		}
	}

	router := gin.Default()
	router.POST("/create/blog-posts", controllers.CreateBlogPost)
	router.GET("/blog-posts", controllers.GetAllBlogPosts)
	router.GET("/blog-posts/:id", controllers.GetBlogPost)
	adminRoutes := router.Group("/admin").Use(middleware.Authenticate())
	adminRoutes.POST("/blog-posts", controllers.CreateBlogPost)
	adminRoutes.PUT("/blog-posts/:id", controllers.UpdateBlogPost)
	adminRoutes.DELETE("/blog-posts/:id", controllers.DeleteBlogPost)

	log.Fatal(router.Run(":8080"))
}

func main() {
	initialize()
}
