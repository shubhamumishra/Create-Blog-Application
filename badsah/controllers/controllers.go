package controllers

import (
	"context"
	"net/http"
	"time"
	"badsah/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client

func GetMongoClient() (*mongo.Client, error) {
	if mongoClient != nil {
		return mongoClient, nil
	}

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		return nil, err
	}

	mongoClient = client

	return client, nil
}

func CreateBlogPost(c *gin.Context) {
	var blogPost models.BlogPost
	if err := c.ShouldBindJSON(&blogPost); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	blogPost.ID = primitive.NewObjectID()
	blogPost.CreatedAt = time.Now()

	client, err := GetMongoClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	collection := client.Database("blog").Collection("blogPosts")
	if _, err := collection.InsertOne(context.Background(), blogPost); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, blogPost)
}

func GetBlogPost(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blog post ID"})
		return
	}

	client, err := GetMongoClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var blogPost models.BlogPost
	if err := client.Database("blog").Collection("blogPosts").FindOne(context.Background(), bson.M{"_id": id}).Decode(&blogPost); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Blog post not found"})
		return
	}

	c.JSON(http.StatusOK, blogPost)
}

func GetAllBlogPosts(c *gin.Context) {
	client, err := GetMongoClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var blogPosts []models.BlogPost
	cursor, err := client.Database("blog").Collection("blogPosts").Find(context.Background(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(context.Background())

	if err := cursor.All(context.Background(), &blogPosts); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, blogPosts)
}

func UpdateBlogPost(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blog post ID"})
	return
	}
	var blogPost models.BlogPost
	if err := c.ShouldBindJSON(&blogPost); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client, err := GetMongoClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result, err := client.Database("blog").Collection("blogPosts").UpdateOne(context.Background(), bson.M{"_id": id}, bson.M{"$set": blogPost})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.ModifiedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Blog post not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Blog post updated successfully"})
}

func DeleteBlogPost(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blog post ID"})
	return
	}
	client, err := GetMongoClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result, err := client.Database("blog").Collection("blogPosts").DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Blog post not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Blog post deleted successfully"})
}


