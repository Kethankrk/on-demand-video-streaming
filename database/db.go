package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func InitMongoDb() {
	var dbConnectErr error
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	Client, dbConnectErr = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://127.0.0.1:27017"))

	if dbConnectErr != nil {
		fmt.Println(dbConnectErr.Error())
		log.Fatal("Error while connecting to mongodb")
	}
}
