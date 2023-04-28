package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Address struct represents an address
type Address struct {
	Street string `bson:"street,omitempty" default:"Unknown"` // cannot figure out the usage of Unknown tag
	City   string `bson:"city,omitempty" default:"Unknown"`   // cannot figure out the usage of Unknown tag
}

// User struct represents a user
type User struct {
	Name    string   `bson:"name"`
	Age     int      `bson:"age,omitempty"`
	Hobbies []string `bson:"hobbies,omitempty"`
	Address *Address `bson:"address"`
}

// main function
func main() {
	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}

	// Disconnect when finish
	defer func() {
		_ = client.Disconnect(context.TODO())
	}()

	// Drop the collection created before if exists
	collection := client.Database("test").Collection("user")
	_ = collection.Drop(context.Background())

	// Create user data, including a pointer to an address
	john := User{
		Name:    "John",
		Age:     30,
		Hobbies: []string{"music", "skiing"},
		Address: &Address{
			Street: "Main St",
			City:   "New York",
		},
	}

	mary := User{
		Name:    "Mary",
		Age:     40,
		Address: &Address{City: "Chicago"},
	}

	tom := User{
		Name: "Tom",
	}

	// Insert data
	_, err = collection.InsertOne(context.TODO(), john)
	_, err = collection.InsertOne(context.TODO(), mary)
	_, err = collection.InsertOne(context.TODO(), tom)

	// Find all users
	cursor, err := collection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		panic(err)
	}

	// Iterate the cursor and print each document
	for cursor.Next(context.TODO()) {
		var user User
		err = cursor.Decode(&user)
		if err != nil {
			panic(err)
		}

		fmt.Println(user.Name, user.Age, user.Hobbies, user.Address)
	}
}
