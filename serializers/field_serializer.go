package serializers

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func BoardSerializer(client *mongo.Client) (bson.M, error) {

	filter := bson.D{{"title", "Standart board"}}
	var board bson.M
	err := client.Database("hpg").Collection("boards").FindOne(nil, filter).Decode(&board)
	fmt.Println(board)
	if err != nil {
		panic(err)
	}
	return board, err
}
