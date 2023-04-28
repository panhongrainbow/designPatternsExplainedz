package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetLatestTenDocuments function returns the latest ten documents
func GetLatestTenDocuments(collection *mongo.Collection) (results []bson.M, err error) {
	// Create a context
	ctx := context.Background()

	// Set options to sort by createdAt field in descending order and limit to 10 documents
	findOptions := options.Find().SetSort(bson.D{{"createdAt", -1}}).SetLimit(10)

	// Find documents in the collection using the options
	cursor, err := collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return
	}

	// Iterate through the cursor
	if err = cursor.All(ctx, &results); err != nil {
		return
	}

	// Close the cursor
	return
}
