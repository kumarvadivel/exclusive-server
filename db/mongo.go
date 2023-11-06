package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var EkartDb *mongo.Database

// collection instances for Ekartdb
var UsersEkart *mongo.Collection

func MongoConnect() {
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URI"))

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		panic(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		panic(err)
	}

	Client = client

	EkartDb = Client.Database("Ekart")
	UsersEkart = EkartDb.Collection("Users")
	// Create a change stream for the entire database
	databaseStreamOptions := options.ChangeStream().SetFullDocument(options.UpdateLookup)
	databaseStream, err := EkartDb.Watch(context.Background(), mongo.Pipeline{}, databaseStreamOptions)
	if err != nil {
		fmt.Println("Failed to create database change stream:", err)
		return
	}

	defer databaseStream.Close(context.Background())

	for {
		if databaseStream.Next(context.Background()) {
			var changeEvent bson.M
			err := databaseStream.Decode(&changeEvent)
			if err != nil {
				fmt.Println("Error decoding change event:", err)
			} else {
				setTimestamps(changeEvent)
				// Print the updated document with timestamps
				fmt.Println("Updated document:", changeEvent["fullDocument"])
			}
		}
	}

	//return Client
}
func setTimestamps(eventDocument bson.M) {
	now := time.Now()

	// Set 'created_at' timestamp for new documents
	if eventDocument["operationType"] == "insert" {
		eventDocument["fullDocument"].(bson.M)["created_at"] = now
		eventDocument["fullDocument"].(bson.M)["updated_at"] = now
	}

	// Set 'updated_at' timestamp for updates
	if eventDocument["operationType"] == "update" {
		eventDocument["updateDescription"].(bson.M)["updatedFields"].(bson.M)["updated_at"] = now
	}
}

func setupHooks(collection *mongo.Collection) {
	// Set up a `beforeInsert` hook to automatically set the `CreatedAt` and `UpdatedAt` timestamps
	// beforeInsert := func(ctx context.Context, doc bson.D) {
	// 	createdAt := time.Now()
	// 	doc = append(doc, bson.E{Key: "created_at", Value: createdAt})
	// 	doc = append(doc, bson.E{Key: "updated_at", Value: createdAt})
	// }

	// _, err := collection.Indexes().CreateOne(
	// 	context.TODO(),
	// 	mongo.IndexModel{
	// 		Keys: bson.M{"created_at": 1},
	// 	},
	// )

	// if err != nil {
	// 	panic(err)
	// }
}

func SelectDb(dbName string) *mongo.Database {
	return Client.Database(dbName)
}

func SelectCollection(db *mongo.Database, collectionName string) *mongo.Collection {
	return db.Collection(collectionName)
}
