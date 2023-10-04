package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type board struct {
	Title  string  `json:"title"`
	Fields []field `json:"fields"`
}

type field struct {
	ID       primitive.ObjectID `bson:"_id"`
	Name     string             `json:"name"`
	Low      int                `json:"low"`
	High     int                `json:"high"`
	Tags     []string           `json:"tags"`
	ImageUrl string             `json:"image"`
	Points   int                `json:"points"`
	Id       int                `json:"id"`
	Games    []string           `json:"games"`
	Rating   int                `json:"rating"`
}
