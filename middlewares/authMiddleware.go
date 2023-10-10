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

		claims, err := helper.ValidateToken(clientToken)
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
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
		claims, err := helper.ValidateToken(clientToken)
		fmt.Println(claims)
		if err != "" {
			c.Set("name", "guest")
			c.Set("role", "GUEST")
		} else {
			c.Set("name", claims.Name)
			c.Set("role", claims.Role)
		}
		c.Next()

	}
}
