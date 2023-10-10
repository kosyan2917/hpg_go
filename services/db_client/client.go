package db_client

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
)

const url = "mongodb://localhost:27017"

type DBClientInterface interface {
	CreateClient() (*mongo.Client, error)
}

type DbClient struct {
	Client *mongo.Client
}

func CreateClient() DbClient {
	var db_client *DbClient
	var once sync.Once
	once.Do(func() {
		serverAPI := options.ServerAPI(options.ServerAPIVersion1)
		opts := options.Client().ApplyURI(url).SetServerAPIOptions(serverAPI)
		client, err := mongo.Connect(context.TODO(), opts)
		if err != nil {
			panic(err)
		}
		db_client = &DbClient{Client: client}
	})
	return *db_client
}
