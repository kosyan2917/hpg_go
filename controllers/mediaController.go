package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
	"path/filepath"
	"strings"
)

func Media() gin.HandlerFunc {
	return func(c *gin.Context) {
		dir := c.Param("dir")
		asset := c.Param("asset")
		if strings.TrimPrefix(asset, "/") == "" {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		fullName := filepath.Join("media/", dir, filepath.FromSlash(path.Clean("/"+asset)))
		fmt.Println(fullName)
		fmt.Println("")
		c.File(fullName)
		c.JSON(http.StatusOK, gin.H{"message": fullName})
	}
}
