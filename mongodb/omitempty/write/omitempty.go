package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// User struct with three fields: Name, Age and Email.serialized if it is empty
type User struct {
	Name  string `bson:"name"`
	Age   int    `bson:"age"`
	Email string `bson:"email"` // <<<<< notice here !
}

func main() {
	/*
		Connect to mongodb using the given URI
		Handle any errors and defer the disconnection
	*/
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	// Select the database and collection to use
	db := client.Database("omitempty")
	collection := db.Collection("users")

	// Create two User objects, one with Email and one without
	user1 := User{
		Name:  "Alice",
		Age:   25,
		Email: "alice@example.com",
	}
	user2 := User{
		Name: "Bob",
		Age:  30,
	}

	// Insert the user1 and user2 object into the collection
	_, err = collection.InsertMany(context.Background(), []interface{}{user1, user2})
	if err != nil {
		log.Fatal(err)
	}

	// Query the collection for documents that have an empty email field
	filter := bson.M{"email": bson.M{"$in": []string{""}}} // <<<<< notice here !
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	// Iterate over the query results and print each document
	for cursor.Next(context.Background()) {
		var user User
		err := cursor.Decode(&user)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(user)
	}

	// Drop the collection to clean up the data
	err = collection.Drop(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}
