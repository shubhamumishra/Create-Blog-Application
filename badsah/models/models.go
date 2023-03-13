package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BlogPost struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Title     string             `bson:"title"`
	Content   string             `bson:"content"`
	CreatedAt time.Time          `bson:"createdAt"`
}

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"username"`
	Password string             `bson:"password"`
	Role     string             `bson:"role"`
}
