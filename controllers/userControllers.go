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
		token, refreshToken, _ := helper.GenerateToken(user.Email, user.Username, user.ID.Hex(), user.Role)
		user.Token = &token
		user.RefreshToken = &refreshToken
		user.Effects = []models.Effect{}
		user.Items = []models.Item{}
		user.Rolls = [][]int{}
		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			msg := fmt.Sprintf("User Details were not Saved")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, resultInsertionNumber)
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
		}
		token, refreshToken, _ := helper.GenerateToken(foundUser.Email, foundUser.Username, foundUser.ID.Hex(), foundUser.Role)
		helper.UpdateAllTokens(token, refreshToken, foundUser.Username)
		err = userCollection.FindOne(ctx, bson.M{}).Decode(&foundUser)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, foundUser)
	}
}

func Me() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		defer cancel()
		err := userCollection.FindOne(ctx, bson.M{"username": name}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}
