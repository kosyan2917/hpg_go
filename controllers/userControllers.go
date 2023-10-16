package controllers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"hpg_backend_go/helpers"
	"hpg_backend_go/models"
	"hpg_backend_go/services/db_client"
	"log"
	"net/http"
	"time"
)

type Tokens struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type ForMe struct {
	Username    string          `json:"username"`
	Avatar      *string         `json:"avatar"`
	Color       string          `json:"color"`
	CurrentGame string          `json:"current_game"`
	Items       []models.Item   `json:"items"`
	Effects     []models.Effect `json:"effects"`
}

type Color struct {
	Color string `json:"color"`
}

var validate = validator.New()
var userCollection = db_client.CreateClient().Client.Database("hpg").Collection("users")

func Signup() gin.HandlerFunc {

	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error sasi": err.Error()})
			return
		}

		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error detected while fetching the email"})
			log.Panic(err)
			return
		}

		password, err := helper.HashPassword(user.Password)
		user.Password = password

		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "The mentioned E-Mail or Phone Number already exists"})
			return
		}
		count, err = userCollection.CountDocuments(ctx, bson.M{"username": user.Username})
		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "The mentioned Username already exists"})
			return
		}

		user.ID = primitive.NewObjectID()
		user.Role = "PLAYER"
		user.Position = 1
		user.CanRoll = true
		token, refreshToken, _ := helper.GenerateToken(user.Email, user.Username, user.ID.Hex(), user.Role)
		user.Token = &token
		user.Color = "#000000"
		user.RefreshToken = &refreshToken
		user.Effects = []models.Effect{}
		user.Items = []models.Item{}
		user.Rolls = [][]int{}
		_, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			msg := fmt.Sprintf("User Details were not Saved")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		var tokens Tokens
		tokens.Token = token
		tokens.RefreshToken = refreshToken
		c.JSON(http.StatusOK, tokens)
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "email or password is incorrect"})
			return
		}

		passwordIsValid := helper.CheckPasswordHash(user.Password, foundUser.Password)
		defer cancel()
		if passwordIsValid != true {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Incorrect Password"})
			return
		}

		if foundUser.Email == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
			return
		}
		token, refreshToken, _ := helper.GenerateToken(foundUser.Email, foundUser.Username, foundUser.ID.Hex(), foundUser.Role)
		helper.UpdateAllTokens(token, refreshToken, foundUser.Username)
		err = userCollection.FindOne(ctx, bson.M{}).Decode(&foundUser)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var tokens Tokens
		tokens.Token = token
		tokens.RefreshToken = refreshToken
		c.JSON(http.StatusOK, tokens)
	}
}

func Me() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.GetString("name")
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		defer cancel()
		err := userCollection.FindOne(ctx, bson.M{"username": name}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var forMe ForMe
		fmt.Println(user)
		forMe.Username = user.Username
		forMe.Avatar = user.AvatarUrl
		forMe.Color = user.Color
		forMe.Items = user.Items
		forMe.Effects = user.Effects
		forMe.CurrentGame = user.CurrentGame
		c.JSON(http.StatusOK, forMe)
	}
}

func Refresh() gin.HandlerFunc {
	return func(c *gin.Context) {
		refresh := c.Request.Header.Get("refresh_token")
		//fmt.Println(refresh)
		claims, err, status := helper.ValidateToken(refresh)
		fmt.Println(claims)
		if err != "" {
			c.JSON(status, gin.H{"error": "Invalid Refresh Token"})
			return
		}
		token, refreshToken, _ := helper.GenerateToken(claims.Email, claims.Name, claims.Uid, claims.Role)
		helper.UpdateAllTokens(token, refreshToken, claims.Name)
		c.JSON(http.StatusOK, gin.H{"token": token, "refresh_token": refreshToken})
	}
}

func ChangeAvatar() gin.HandlerFunc {
	return func(c *gin.Context) {
		// single file
		file, _ := c.FormFile("file")
		log.Println(file.Filename)
		client := db_client.CreateClient()
		collection := client.Client.Database("hpg").Collection("users")
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		avatarUrl := "/media/avatars/" + file.Filename
		_, err := collection.UpdateOne(ctx, bson.M{"username": c.GetString("name")}, bson.M{"$set": bson.M{"avatar_url": avatarUrl}})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// Upload the file to specific dst.
		err = c.SaveUploadedFile(file, "media/avatars/"+file.Filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
}

func ChangeColor() gin.HandlerFunc {
	return func(c *gin.Context) {
		newColor, _ := c.GetPostForm("color")
		fmt.Println(newColor)
		client := db_client.CreateClient()
		collection := client.Client.Database("hpg").Collection("users")
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		_, err := collection.UpdateOne(ctx, bson.M{"username": c.GetString("name")}, bson.M{"$set": bson.M{"color": newColor}})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		defer cancel()
	}
}
