package handler

import (
	"context"
	"gcs/model"
	"gcs/utils"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

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
	role := c.Query("role")

	userCollection := handler.mongodb.Collection("user")
	filer := make(bson.D, 1)
	if name != "" {
		filer = append(filer, bson.E{
			"name", primitive.Regex{Pattern: name, Options: ""},
		})
	}
	if role != "" {
		filer = append(filer, bson.E{
			"role", primitive.Regex{Pattern: role, Options: ""},
		})
	}
	cursor, err := userCollection.Find(handler.ctx, filer, options.Find().SetProjection(bson.M{"password": 0}))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusInternalServerError,
			"msg":    "数据库操作失败",
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

func (handler *UserHandler) NewUser(c *gin.Context) {
	var json model.User
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"msg":    "参数错误",
			"error":  err.Error(),
			"data":   gin.H{},
		})
		return
	}

	sysUsername := viper.GetString("admin.username")
	if len(sysUsername) > 0 && json.Name == sysUsername {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"msg":    "用户名已存在",
			"data":   gin.H{},
		})
		return
	}

	var exitUser model.User
	userCollection := handler.mongodb.Collection("user")

	err := userCollection.FindOne(handler.ctx, bson.M{
		"name": json.Name,
	}).Decode(&exitUser)

	if err == mongo.ErrNoDocuments {
		json.ID = primitive.NewObjectID()
		json.Password = utils.Sha256(json.Password)
		result, err := userCollection.InsertOne(handler.ctx, json)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusInternalServerError,
				"msg":    "数据库操作失败",
				"error":  err.Error(),
				"data":   gin.H{},
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": 0,
			"msg":    "ok",
			"data":   result,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusBadRequest,
		"msg":    "用户名已存在",
		"data":   gin.H{},
	})
}

func (handler *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var json model.User
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"msg":    "参数错误",
			"error":  err.Error(),
			"data":   gin.H{},
		})
		return
	}

	objectId, _ := primitive.ObjectIDFromHex(id)
	userCollection := handler.mongodb.Collection("user")

	setElements := bson.D{
		{"role", json.Role},
	}
	if len(json.Password) > 0 {
		setElements = append(setElements, bson.E{"password", utils.Sha256(json.Password)})
	}

	update := bson.D{
		{
			"$set",
			setElements,
		},
	}
	_, err := userCollection.UpdateByID(handler.ctx, objectId, update)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusInternalServerError,
			"msg":    "数据库操作失败",
			"error":  err.Error(),
			"data":   gin.H{},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "ok",
		"data":   gin.H{},
	})
}

func (handler *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	userCollection := handler.mongodb.Collection("user")
	_, err := userCollection.DeleteOne(handler.ctx, bson.M{
		"_id": objectId,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusInternalServerError,
			"msg":    "数据库操作失败",
			"error":  err.Error(),
			"data":   gin.H{},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "ok",
		"data":   gin.H{},
	})
}
