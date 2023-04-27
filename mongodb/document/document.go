package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// User struct represents a user
type User struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Name    string             `bson:"name"`
	Age     int                `bson:"age"`
	Address string             `bson:"address"`
}

func main() {
	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	// Make connection is available
	collection := client.Database("property").Collection("user")

	// Insert a user
	user := User{
		Name:    "Tom",
		Age:     18,
		Address: "New York",
	}
	result, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		panic(err)
	}

	// Select the user
	var queriedUser User
	id := result.InsertedID.(primitive.ObjectID)
	err = collection.FindOne(context.TODO(), User{ID: id}).Decode(&queriedUser)
	if err != nil {
		panic(err)
	}
	println(queriedUser.Name, queriedUser.Age, queriedUser.Address)

	// Filter and update the user
	filter := User{ID: id}
	update := User{
		Age: 25,
	}
	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		panic(err)
	}

	// Read the updated user
	var updatedUser User
	err = collection.FindOne(context.TODO(), filter).Decode(&updatedUser)
	if err != nil {
		panic(err)
	}
	println(updatedUser.Name, updatedUser.Age, updatedUser.Address)
}
