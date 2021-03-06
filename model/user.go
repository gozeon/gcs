package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LoginUser struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type User struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Name     string             `json:"name" bson:"name" binding:"required"`
	Password string             `json:"password" bson:"password" binding:"required"`
	Role     string             `json:"role" bson:"role" binding:"required"`
}
