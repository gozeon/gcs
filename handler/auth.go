package handler

import (
	"context"
	"gcs/cont"
	"gcs/model"
	"gcs/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

type AuthHandler struct {
	ctx     context.Context
	mongodb *mongo.Database
}

func NewAuthHandler(ctx context.Context, mongodb *mongo.Database) *AuthHandler {
	return &AuthHandler{
		ctx:     ctx,
		mongodb: mongodb,
	}
}

func (handler *AuthHandler) ShowLogin(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()

	next := c.DefaultQuery("next", "/portal/r/w/dash")
	errMsg := c.Query("errMsg")

	c.HTML(http.StatusOK, "login.html", gin.H{
		"title":  "login",
		"next":   next,
		"errMsg": errMsg,
	})
}

func (handler *AuthHandler) Login(c *gin.Context) {
	var json model.LoginUser
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusUnauthorized,
			"msg":    "参数错误",
		})
		return
	}
	session := sessions.Default(c)

	// find in conf
	if json.Username == viper.GetString("admin.username") && json.Password == viper.GetString("admin.password") {
		session.Set("username", json.Username)
		var role cont.Role = cont.Admin
		session.Set("role", role.String())
		session.Save()

		c.JSON(http.StatusOK, gin.H{
			"status": 0,
			"msg":    "登录成功",
			"data": gin.H{
				"username": json.Username,
				"role":     role.String(),
			},
		})

	} else {
		var exitUser model.User
		// find in db
		userCollection := handler.mongodb.Collection("user")
		err := userCollection.FindOne(handler.ctx, bson.M{
			"name":     json.Username,
			"password": utils.Sha256(json.Password),
		}).Decode(&exitUser)

		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusBadRequest,
				"msg":    "用户名或密码不正确",
				"error":  err.Error(),
				"data":   gin.H{},
			})
			return
		}

		session.Set("username", exitUser.Name)
		session.Set("role", exitUser.Role)
		session.Save()

		c.JSON(http.StatusOK, gin.H{
			"status": 0,
			"msg":    "登录成功",
			"data": gin.H{
				"username": json.Username,
				"role":     exitUser.Role,
			},
		})
	}
}

func (handler *AuthHandler) Logout(c *gin.Context) {
	session := sessions.Default(c)

	session.Clear()
	session.Save()

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "ok",
		"data":   gin.H{},
	})
}

func AuthMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		username := session.Get("username")

		// if login
		if username == nil {
			c.Redirect(http.StatusFound, "/portal/r/w/login?next="+c.Request.RequestURI+"&errMsg="+"获取登录信息失败，请重新登录")
			return
		}

		role := session.Get("role").(string)

		// auth
		if len(roles) > 0 {
			if !utils.Contains(roles, role) {
				c.Redirect(http.StatusFound, "/portal/r/w/login?next="+c.Request.RequestURI+"&errMsg="+"无权访问，请重新登录")

				//c.AbortWithStatusJSON(http.StatusOK, gin.H{
				//	"status": http.StatusUnauthorized,
				//	"msg":    "无权访问",
				//	"data":   gin.H{},
				//})
				return
			}
		}

		c.Next()
	}
}
