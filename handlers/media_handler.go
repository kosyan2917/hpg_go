package handlers

import (
	"github.com/gin-gonic/gin"
	"hpg_backend_go/serializers"
)

func Board(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")
	data, err := serializers.Serialize()
	//fmt.Println(string(data))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	}
	//unquoted, err := strconv.Unquote(string(data))
	c.JSON(200, gin.H{"data": data})
}
