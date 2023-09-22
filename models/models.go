package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type board struct {
	Title  string  `bson:"title"`
	Fields []field `bson:"fields"`
}

type field struct {
	ID       primitive.ObjectID `bson:"_id"`
	Name     string             `bson:"name"`
	Low      int                `bson:"low"`
	High     int                `bson:"high"`
	Tags     []string           `bson:"tags"`
	ImageUrl string             `bson:"image"`
	Points   int                `bson:"points"`
	Id       int                `bson:"id"`
	Games    []string           `bson:"games"`
	Rating   int                `bson:"rating"`
}
