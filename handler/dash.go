package handler

import (
	"context"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type DashHandler struct {
	ctx context.Context
}

func NewDashHandler(ctx context.Context) *DashHandler {
	return &DashHandler{
		ctx: ctx,
	}
}

func (handler *DashHandler) ShowDash(c *gin.Context) {
	session := sessions.Default(c)

	c.HTML(http.StatusOK, "dash.html", gin.H{
		"title":    "dash",
		"username": session.Get("username"),
		"role":     session.Get("role"),
	})
}
