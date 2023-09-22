package serializers

import (
	"go.mongodb.org/mongo-driver/bson"
	"hpg_backend_go/services/db_client"
)

func Serialize() (bson.M, error) {
	client, err := db_client.CreateClient()
	if err != nil {
		panic(err)
	}
	filter := bson.D{{"title", "Standart board"}}
	var board bson.M
	err = client.Database("hpg").Collection("boards").FindOne(nil, filter).Decode(&board)
	if err != nil {
		panic(err)
	}
	return board, err
}
