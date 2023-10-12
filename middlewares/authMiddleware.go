package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	helper "hpg_backend_go/helpers"
	"net/http"
)

func UserAuthenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("No Authorization Header Provided")})
			c.Abort()
			return
		}

		claims, err, status := helper.ValidateToken(clientToken)
		if err == "token is expired" {
			c.JSON(status, gin.H{"error": err})
			c.Abort()
			return
		}
		if err != "" {
			c.JSON(status, gin.H{"error": err})
			c.Abort()
			return
		}
		c.Set("email", claims.Email)
		c.Set("name", claims.Name)
		c.Next()
	}
}

func BoardAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token")
		claims, err, _ := helper.ValidateToken(clientToken)
		if err != "" {
			c.Set("name", "guest")
			c.Set("role", "GUEST")
			c.Set("status", 403)
		} else {
			c.Set("name", claims.Name)
			c.Set("role", claims.Role)
			c.Set("status", 200)
		}
		c.Next()

	}
}
