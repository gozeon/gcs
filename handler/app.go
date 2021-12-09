package handler

import (
	"context"
	"gcs/model"
	"github.com/gin-contrib/sessions"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AppHandler struct {
	ctx     context.Context
	mongodb *mongo.Database
}

func NewAppHandler(ctx context.Context, mongodb *mongo.Database) *AppHandler {
	return &AppHandler{
		ctx:     ctx,
		mongodb: mongodb,
	}
}

func (handler *AppHandler) ShowApps(c *gin.Context) {
	session := sessions.Default(c)

	c.HTML(http.StatusOK, "app-card.html", gin.H{
		"title":    "应用中心",
		"username": session.Get("username"),
		"role":     session.Get("role"),
	})
}

func (handler *AppHandler) ShowApp(c *gin.Context) {
	id := c.Param("id")
	c.HTML(http.StatusOK, "app-detail.html", gin.H{
		"title": "",
		"id":    id,
	})
}

func (handler *AppHandler) AppInfo(c *gin.Context) {
	id := c.Param("id")
	var app model.App

	objectId, _ := primitive.ObjectIDFromHex(id)
	appCollection := handler.mongodb.Collection("app")

	err := appCollection.FindOne(handler.ctx, bson.M{
		"_id": objectId,
	}).Decode(&app)

	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusBadRequest,
			"msg":    "数据库操作未找到",
			"error":  err.Error(),
			"data":   gin.H{},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "ok",
		"data":   app,
	})
}

func (handler *AppHandler) Apps(c *gin.Context) {
	name := c.Query("name")

	appCollection := handler.mongodb.Collection("app")
	filer := make(bson.D, 1)
	if name != "" {
		filer = append(filer, bson.E{
			"name", primitive.Regex{Pattern: name, Options: ""},
		})
	}

	cursor, err := appCollection.Find(handler.ctx, filer, options.Find())
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

	apps := make([]model.App, 0)
	for cursor.Next(handler.ctx) {
		var node model.App
		cursor.Decode(&node)
		apps = append(apps, node)
	}

	// if cache in redis
	// cacheData, _ := json.Marshal(users)
	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "ok",
		"data":   apps,
	})
}

func (handler *AppHandler) NewApp(c *gin.Context) {
	var json model.App
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusBadRequest,
			"msg":    "参数错误",
			"error":  err.Error(),
			"data":   gin.H{},
		})
		return
	}

	appCollection := handler.mongodb.Collection("app")
	json.ID = primitive.NewObjectID()
	result, err := appCollection.InsertOne(handler.ctx, json)
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
}

func (handler *AppHandler) DeleteApp(c *gin.Context) {
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	appCollection := handler.mongodb.Collection("app")
	_, err := appCollection.DeleteOne(handler.ctx, bson.M{
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

func (handler *AppHandler) UpdateApp(c *gin.Context) {
	id := c.Param("id")
	var json model.App
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
	appCollection := handler.mongodb.Collection("app")

	json.ID = objectId

	update := bson.D{
		{
			"$set",
			json,
		},
	}

	_, err := appCollection.UpdateByID(handler.ctx, objectId, update)
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
