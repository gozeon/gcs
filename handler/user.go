package handler

import (
	"context"
	"gcs/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"

	"github.com/gin-contrib/sessions"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	ctx     context.Context
	mongodb *mongo.Database
}

func NewUserHandler(ctx context.Context, mongodb *mongo.Database) *UserHandler {
	return &UserHandler{
		ctx:     ctx,
		mongodb: mongodb,
	}
}

func (handler *UserHandler) Me(c *gin.Context) {
	session := sessions.Default(c)

	c.JSON(http.StatusOK, gin.H{
		"username": session.Get("username"),
		"role":     session.Get("role"),
	})
}

func (handler *UserHandler) Logout(c *gin.Context) {
	session := sessions.Default(c)

	session.Clear()

	c.JSON(http.StatusOK, gin.H{
		"username": session.Get("username"),
		"role":     session.Get("role"),
	})
}

func (handler *UserHandler) Users(c *gin.Context) {
	name := c.Query("name")

	userCollection := handler.mongodb.Collection("user")
	filer := make(bson.D, 1)
	if name != "" {
		filer = append(filer, bson.E{
			"name", primitive.Regex{Pattern: name, Options: ""},
		})
	}
	cursor, err := userCollection.Find(handler.ctx, filer)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusInternalServerError,
			"msg":    "系统错误",
			"error":  err.Error(),
			"data":   gin.H{},
		})
		return
	}
	defer cursor.Close(handler.ctx)

	users := make([]model.User, 0)
	for cursor.Next(handler.ctx) {
		var node model.User
		cursor.Decode(&node)
		users = append(users, node)
	}

	// if cache in redis
	// cacheData, _ := json.Marshal(users)
	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "ok",
		"data":   users,
	})
}
