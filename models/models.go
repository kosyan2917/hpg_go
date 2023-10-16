package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Board struct {
	Title  string             `json:"title"`
	Fields []Field            `json:"fields"`
	ID     primitive.ObjectID `bson:"_id"`
}

type Field struct {
	Name      *string   `bson:"name" json:"name"`
	Low       *int      `bson:"low" json:"low"`
	High      *int      `bson:"high" json:"high"`
	Tags      *[]string `bson:"tags" json:"tags"`
	ImageUrl  *string   `bson:"image" json:"imageUrl"`
	Points    *int      `bson:"points" json:"points"`
	Id        *int      `bson:"id" json:"id"`
	Games     *[]string `bson:"games" json:"games"`
	Rating    *string   `bson:"rating" json:"rating"`
	Condition *string   `bson:"condition" json:"condition"`
	Type      *string   `bson:"type" json:"type"`
}

type User struct {
	ID           primitive.ObjectID `bson:"_id"`
	Username     string             `json:"username"`
	Password     string             `json:"password"`
	Email        string             `json:"email"`
	Role         string             `json:"role"`
	Rolls        [][]int            `json:"rolls"`
	Points       int                `json:"points"`
	Token        *string            `json:"token"`
	RefreshToken *string            `bson:"refresh_token" json:"refresh_token"`
	Items        []Item             `json:"items"`
	Effects      []Effect           `json:"effects"`
	AvatarUrl    *string            `bson:"avatar_url" json:"avatar_url"`
	Color        string             `json:"color"`
	Position     int                `json:"position"`
	CurrentGame  string             `bson:"current_game" json:"current_game"`
	CanRoll      bool               `bson:"can_roll" json:"can_roll"`
}

type Item struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageUrl    string `json:"imageUrl"`
	Count       int    `json:"count"`
}

type Effect struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
	Source      Item   `json:"source"`
	ImageUrl    string `json:"imageUrl"`
}
