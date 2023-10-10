package serializers

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"hpg_backend_go/models"
)

func BoardSerializer(client *mongo.Client) (models.Board, error) {

	filter := bson.D{{"title", "Standart board"}}
	var board models.Board
	err := client.Database("hpg").Collection("boards").FindOne(nil, filter).Decode(&board)
	fmt.Println(board.Fields)
	if err != nil {
		panic(err)
	}
	return board, err
}
