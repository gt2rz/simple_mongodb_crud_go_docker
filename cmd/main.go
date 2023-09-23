package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	mongoOptions := options.Client().ApplyURI("mongodb://simple-mongo-crud-db:27017")
	mongoClient, err := mongo.Connect(ctx, mongoOptions)

	defer func() {
		cancel()
		if err = mongoClient.Disconnect(ctx); err != nil {
			log.Fatalf("Cannot disconnect from MongoDB: %v", err)
		}
	}()

	if err != nil {
		log.Fatalf("Cannot connect to MongoDB: %v", err)
	}

	err = mongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatalf("Cannot ping MongoDB: %v", err)
		return
	}

	log.Println("Connected to MongoDB!")

	database := mongoClient.Database("demo")
	demoCollection := database.Collection("demo")
	demoCollection.Drop(ctx)

	// insert one document
	document := bson.M{
		"name":       "John",
		"content":    "Hello World!",
		"bank_money": 1000000,
		"created_at": time.Now(),
	}

	result, err := demoCollection.InsertOne(ctx, document)
	if err != nil {
		log.Fatalf("Cannot insert document: %v", err)
		return
	}

	log.Printf("Inserted document with ID: %v\n", result.InsertedID)

	// query all data
	fmt.Println("Query all data")
	options := options.Find()
	cursor, err := demoCollection.Find(ctx, options)
	if err != nil {
		log.Fatalf("Cannot query all data: %v", err)
		return
	}

	results := []bson.M{}
	if err = cursor.All(ctx, &results); err != nil {
		log.Fatalf("Cannot decode cursor to results: %v", err)
		return
	}

	for _, result := range results {
		fmt.Println(result)
	}

	// query one data
	fmt.Println("Query one data")
	filter := bson.M{"name": "John"}
	var resultOne bson.M
	if err = demoCollection.FindOne(ctx, filter).Decode(&resultOne); err != nil {
		log.Fatalf("Cannot decode cursor to results: %v", err)
		return
	}

	fmt.Println(resultOne)

	// update one data
	fmt.Println("Update one data")
	update := bson.M{"$set": bson.M{"name": "John Doe"}}
	updateResult, err := demoCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Fatalf("Cannot update one data: %v", err)
		return
	}

	fmt.Println(updateResult)

	// update many data
	fmt.Println("Update many data")
	updateManyResult, err := demoCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		log.Fatalf("Cannot update many data: %v", err)
		return
	}

	fmt.Println(updateManyResult)

	// delete one data
	fmt.Println("Delete one data")
	deleteResult, err := demoCollection.DeleteOne(ctx, filter)
	if err != nil {
		log.Fatalf("Cannot delete one data: %v", err)
		return
	}

	fmt.Println(deleteResult)

	// delete many data
	fmt.Println("Delete many data")
	deleteManyResult, err := demoCollection.DeleteMany(ctx, filter)
	if err != nil {
		log.Fatalf("Cannot delete many data: %v", err)
		return
	}

	fmt.Println(deleteManyResult)

	// drop collection
	fmt.Println("Drop collection")
	if err = demoCollection.Drop(ctx); err != nil {
		log.Fatalf("Cannot drop collection: %v", err)
		return
	}

	fmt.Println("Done!")

}
