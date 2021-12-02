package main

import (
	"context"
	"fmt"
	"gcs/cont"
	"gcs/handler"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/gin-gonic/gin"
)

var authHandler *handler.AuthHandler
var dashHandler *handler.DashHandler
var userHandler *handler.UserHandler

func init() {
	ctx := context.Background()

	// config
	viper.SetConfigFile("config.json")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Prepare connected to MongoDB: " + viper.GetString("mongo.uri") + ", " + viper.GetString("mongo.db"))

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(viper.GetString("mongo.uri")))
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB: " + viper.GetString("mongo.uri") + ", " + viper.GetString("mongo.db"))
	mongodb := client.Database(viper.GetString("mongo.db"))

	authHandler = handler.NewAuthHandler(ctx, mongodb)
	dashHandler = handler.NewDashHandler(ctx)
	userHandler = handler.NewUserHandler(ctx, mongodb)
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/**.html")
	router.Static("/assets", "./assets")
	router.StaticFile("/favicon.ico", "./assets/favicon.ico")

	// session
	store := cookie.NewStore([]byte(viper.GetString("secret.session")))
	router.Use(sessions.Sessions(viper.GetString("secret.cookie.name"), store))

	v1 := router.Group(fmt.Sprintf("/%s", viper.Get("server.base")))
	{
		// ping
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})

		webRouter := v1.Group("/r/w")
		{
			webRouter.GET("/login", authHandler.ShowLogin)
			webRouter.GET("/dash", handler.AuthMiddleware(), dashHandler.ShowDash)
		}

		// json router
		jsonRouter := v1.Group("/r/jd")
		{
			jsonRouter.POST("/login", authHandler.Login)
			jsonRouter.GET("/logout", authHandler.Logout)
			jsonRouter.GET("/me", handler.AuthMiddleware(), userHandler.Me)

			userRouter := jsonRouter.Group("/user")
			userRouter.Use(handler.AuthMiddleware(cont.Admin.String()))
			{
				userRouter.GET("/", userHandler.Users)
				userRouter.POST("/", userHandler.NewUser)
				userRouter.DELETE("/:id", userHandler.DeleteUser)
				userRouter.PUT("/:id", userHandler.UpdateUser)
			}
		}
	}

	router.Run(fmt.Sprintf(":%d", viper.GetInt("server.port")))

}
