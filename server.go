package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	srv := gin.Default()
	srv.Static("/media", "./media")
	srv.Run(":8080")
}
