package controllers

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"hpg_backend_go/models"
	"hpg_backend_go/serializers"
	"hpg_backend_go/services/db_client"
)

type BoardResponse struct {
	Title  string             `json:"title"`
	Fields []models.Field     `json:"fields"`
	ID     primitive.ObjectID `bson:"_id"`
	User   User               `json:"user"`
}

type User struct {
	Name string `json:"name"`
	Role string `json:"role"`
}

func Board() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")
		client := db_client.CreateClient()
		data, err := serializers.BoardSerializer(client.Client)
		var user User
		user.Name = c.GetString("name")
		user.Role = c.GetString("role")
		//fmt.Println(string(data))
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
		}
		var boardResponse BoardResponse
		boardResponse.Title = data.Title
		boardResponse.Fields = data.Fields
		boardResponse.ID = data.ID
		boardResponse.User = user

		//unquoted, err := strconv.Unquote(string(data))
		c.JSON(200, boardResponse)
	}
}
