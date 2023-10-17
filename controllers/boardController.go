package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	helper "hpg_backend_go/helpers"
	"hpg_backend_go/models"
	"hpg_backend_go/serializers"
	"hpg_backend_go/services/db_client"
	"strconv"
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
			return
		}
		var boardResponse BoardResponse
		boardResponse.Title = data.Title
		boardResponse.Fields = data.Fields
		boardResponse.ID = data.ID
		boardResponse.User = user

		//unquoted, err := strconv.Unquote(string(data))
		c.JSON(c.GetInt("status"), boardResponse)
	}
}

func Users() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")
		client := db_client.CreateClient()
		data, err := serializers.UsersSerializer(client.Client)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		fmt.Println(data)
		c.JSON(200, data)
	}
}

func Roll() gin.HandlerFunc {
	return func(c *gin.Context) {
		canRoll, pos := helper.CanRoll(c.GetString("name"))
		fmt.Println(canRoll)
		if !canRoll {
			c.JSON(500, gin.H{"error": "You can't roll", "pos": pos})
			return
		}
		client := db_client.CreateClient()
		collection := client.Client.Database("hpg").Collection("users")
		filter := bson.D{{"username", c.GetString("name")}}
		update := bson.D{{"$set", bson.D{{"can_roll", false}}}}
		_, err := collection.UpdateOne(nil, filter, update)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		dice1, dice2 := helper.Roll()
		newPos, err := helper.Move(c.GetString("name"), dice1, dice2)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		collection = client.Client.Database("hpg").Collection("boards")
		filter = bson.D{{"title", "Standart board"}}
		var board models.Board
		err = collection.FindOne(nil, filter).Decode(&board)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		var msg = ""
		err = nil
		if *board.Fields[newPos-1].Type == "program" {
			msg, canRoll, err = helper.ProgramField(c.GetString("name"), newPos, *board.Fields[newPos-1].Name)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			if canRoll {
				update := bson.D{{"$set", bson.D{{"can_roll", true}}}}
				_, err = collection.UpdateOne(nil, filter, update)
				if err != nil {
					c.JSON(500, gin.H{"error": err.Error()})
					return
				}
			}
		}
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"dice1": fmt.Sprint("static/dices/", dice1, ".png"),
			"dice2":  fmt.Sprint("static/dices/", dice2, ".png"),
			"newPos": newPos, "msg": msg})
	}
}

func SetCurrentGame() gin.HandlerFunc {
	return func(c *gin.Context) {
		client := db_client.CreateClient()
		collection := client.Client.Database("hpg").Collection("users")
		game, _ := c.GetPostForm("game")
		fmt.Println(game)
		filter := bson.D{{"username", c.GetString("name")}}
		update := bson.D{{"$set", bson.D{{"current_game", game}, {"can_roll", false}}}}
		_, err := collection.UpdateOne(nil, filter, update)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"status": "success"})
	}
}

func GetCurrentGame() gin.HandlerFunc {
	return func(c *gin.Context) {
		client := db_client.CreateClient()
		collection := client.Client.Database("hpg").Collection("users")
		var user models.User
		filter := bson.D{{"username", c.GetString("name")}}
		err := collection.FindOne(nil, filter).Decode(&user)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"current_game": user.CurrentGame})
	}
}

func Field() gin.HandlerFunc {
	return func(c *gin.Context) {
		client := db_client.CreateClient()
		collection := client.Client.Database("hpg").Collection("boards")
		var board models.Board
		filter := bson.D{{"title", "Standart board"}}
		err := collection.FindOne(nil, filter).Decode(&board)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(500, gin.H{"error": "id must be integer"})
			return
		}
		if id > 40 || id < 1 {
			c.JSON(500, gin.H{"error": "id must be between 1 and 40"})
			return
		}
		var field models.Field
		field = board.Fields[id-1]
		c.JSON(200, field)
	}
}
