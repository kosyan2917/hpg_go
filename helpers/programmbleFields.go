package helper

import (
	"go.mongodb.org/mongo-driver/bson"
)

func ProgramField(username string, fieldId int, fieldName string) (string, bool, error) {
	var msg string
	var rollGame bool
	var err error
	switch fieldName {
	case "Старт":
		msg, rollGame, err = Start(username)
	}
	return msg, rollGame, err
}

func Start(username string) (string, bool, error) {
	filter := bson.D{{"username", username}}
	update := bson.D{{"$set", bson.D{{"can_roll", true}, {"position", 2}}}}
	msg := "По правилам при попадании на клетку старт вы переходите на одну клетку вперед"
	_, err := userCollection.UpdateOne(nil, filter, update)
	return msg, false, err
}

func Podlyanka(username string) error {
	filter := bson.D{{"username", username}}
	update := bson.D{{"$set", bson.D{{"can_roll", true}}}}
	_, err := userCollection.UpdateOne(nil, filter, update)
	return err
}
