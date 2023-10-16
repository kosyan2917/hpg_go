package helper

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"hpg_backend_go/services/db_client"
	"math/rand"
	"time"
)

type oldPos struct {
	Pos int `json:"position" bson:"position"`
}

func Roll() (int, int) {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	return r.Intn(6) + 1, r.Intn(6) + 1
}

func Move(username string, dice1 int, dice2 int) (newPos int, err error) {
	client := db_client.CreateClient()
	collection := client.Client.Database("hpg").Collection("users")
	filter := bson.D{{"username", username}}
	var pos oldPos
	err = collection.FindOne(context.TODO(), filter).Decode(&pos)
	if err != nil {
		return 0, err
	}

	newPos = pos.Pos + dice1 + dice2
	if newPos > 40 {
		newPos = newPos - 40
	}
	update := bson.D{{"$set", bson.D{{"position", newPos}}}}
	_, err = collection.UpdateOne(context.TODO(), filter, update)
	return newPos, err
}

func CanRoll(username string) (bool, int) {
	client := db_client.CreateClient()
	collection := client.Client.Database("hpg").Collection("users")
	filter := bson.D{{"username", username}}
	var user struct {
		CanRoll bool `json:"can_roll" bson:"can_roll"`
		Pos     int  `json:"position" bson:"position"`
	}
	err := collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return false, 0
	}
	fmt.Println(user.CanRoll)
	return user.CanRoll, user.Pos
}
