package main

import (
	"context"
	"fmt"
	"github.com/panhongrainbow/designPatternsExplainedz/mongodb/bucket/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// main function gets a slice of Celsius structs
func main() {
	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = client.Disconnect(context.TODO())
	}()

	// Select database and collection
	db := client.Database("iot")
	ctx := context.TODO()
	collection := db.Collection("temp_gte_20")

	// Get the latest ten documents
	latest, err := GetLatestTenDocuments(collection)
	if err != nil {
		panic(err)
	}

	// Print the latest ten documents
	fmt.Println()
	fmt.Println("print the latest ten documents")
	for i, document := range latest {
		fmt.Println(i, document)
	}

	// Get the latest ten documents with a temperature greater than 20 by using the $bucket operator
	pipeline := mongo.Pipeline{
		bson.D{{"$match", bson.D{{"temperature", bson.D{{"$gt", 20}}}}}},
		bson.D{{"$match", bson.D{{"$expr", bson.D{{"$lt", bson.A{"$NOW", "createdAt"}}}}}}},
		{{"$sort", bson.D{{"createdAt", -1}}}},
		{{"$limit", 10}},
		{{"$bucket", bson.M{ // lookup the reviews by product ID
			"groupBy":    "$temperature",
			"boundaries": []int{20, 25, 30, 35},
			"default":    "other",
			"output": bson.M{
				"count":     bson.M{"$sum": 1},
				"documents": bson.M{"$push": "$$ROOT"},
			},
		}}},
	}

	// Aggregate the documents
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		panic(err)
	}

	// Make a slice of slices of group documents
	var temps = make([][]model.Temperature, 4)
	var i = 0
	for cursor.Next(ctx) {
		var t struct {
			Documents []model.Temperature `bson:"documents"`
		}
		err = cursor.Decode(&t)
		if err != nil {
			_ = cursor.Close(ctx)
			panic(err)
		}
		temps[i] = t.Documents
		i++
	}
	if err = cursor.Err(); err != nil {
		_ = cursor.Close(ctx)
		panic(err)
	}
	if err = cursor.Close(ctx); err != nil {
		panic(err)
	}

	// Print the documents in each bucket
	fmt.Println()
	fmt.Println("Print the documents in each bucket")
	for i := range temps {
		fmt.Println("group ", i+1, temps[i])
	}
}
