package serializers

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"hpg_backend_go/models"
)

type user_av struct {
	Username  string `json:"username"`
	Pos       int    `json:"pos"`
	Points    int    `json:"points"`
	AvatarUrl string `json:"avatar_url"`
	Color     string `json:"color"`
}

func BoardSerializer(client *mongo.Client) (models.Board, error) {
	filter := bson.D{{"title", "Standart board"}}
	var board models.Board
	err := client.Database("hpg").Collection("boards").FindOne(nil, filter).Decode(&board)
	if err != nil {
		panic(err)
	}
	return board, err
}

func UsersSerializer(client *mongo.Client) ([]user_av, error) {
	var users []user_av
	cur, err := client.Database("hpg").Collection("users").Find(context.TODO(), bson.D{})
	if err != nil {
		panic(err)
	}
	for cur.Next(context.TODO()) {
		var user user_av
		var dbUser models.User
		err := cur.Decode(&dbUser)
		if err != nil {
			panic(err)
		}
		user.Username = dbUser.Username
		user.Pos = dbUser.Position
		user.Points = dbUser.Points
		user.AvatarUrl = *dbUser.AvatarUrl
		user.Color = dbUser.Color
		users = append(users, user)
	}
	return users, err
}
