package handler

import (
	"context"
	"gcs/cont"
	"gcs/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
)

type AuthHandler struct {
	ctx context.Context
}

func NewAuthHandler(ctx context.Context) *AuthHandler {
	return &AuthHandler{
		ctx: ctx,
	}
}

func (handler *AuthHandler) ShowLogin(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()

	next := c.DefaultQuery("next", "/portal/r/w/dash")
	errMsg := c.Query("errMsg")

	c.HTML(http.StatusOK, "login.html", gin.H{
		"title":  session.Get("username"),
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
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusUnauthorized,
			"msg":    "用户名或密码不正确",
		})
	}
}

func AuthMiddleware(role ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		username := session.Get("username")

		if username == nil {
			c.Redirect(http.StatusFound, "/portal/r/w/login?next="+c.Request.RequestURI+"&errMsg="+"获取信息失败，请重新登录")
			return
		}

		c.Next()

		// if len(role) > 0 {
		// 	currentRole := session.Get("role").(string)

		// 	if utils.Contains(role, currentRole) {
		// 		c.Next()
		// 	} else {
		// 		c.JSON(http.StatusUnauthorized, gin.H{
		// 			"message": "unauthorized",
		// 		})
		// 		return
		// 	}
		// } else {
		// 	if username == nil {
		// 		c.JSON(http.StatusUnauthorized, gin.H{
		// 			"message": "unauthorized",
		// 		})
		// 		return
		// 	} else {
		// 		c.Next()
		// 	}
		// }
	}
}
