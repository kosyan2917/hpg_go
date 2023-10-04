package main

import (
	"github.com/gin-gonic/gin"
	"hpg_backend_go/handlers"
	"hpg_backend_go/services/db_client"
)

func main() {
	srv := gin.Default()
	srv.Use(CORSMiddleware())
	client, err := db_client.CreateClient()
	if err != nil {
		panic(err)
	}
	srv.Static("/media", "./media")
	srv.GET("/board", handlers.Board)
	srv.Run(":8080")
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
