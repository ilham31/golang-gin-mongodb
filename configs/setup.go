package configs

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Client = ConnectDB()

func ConnectDB() *mongo.Client {
	client, error := mongo.NewClient(options.Client().ApplyURI(EnvironmentMongoUri()))
	if error != nil {
		log.Fatal(error)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	error = client.Connect(ctx)
	if error != nil {
		log.Fatal(error)
	}

	error = client.Ping(ctx, nil)
	if error != nil {
		log.Fatal(error)
	}

	fmt.Println("Connected to MongoDB")
	return client
}

func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("golangGinMongo").Collection(collectionName)
	return collection
}
