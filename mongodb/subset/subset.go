package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

// Product struct represents a product in a store
type Product struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `bson:"name,omitempty"`
	Price    float64            `bson:"price,omitempty"`
	Reviews1 []Review           `bson:"reviews1,omitempty"`
	Reviews2 []Review           `bson:"reviews2,omitempty"`
}

// Review represents a customer review of a product type Review
type Review struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	ProductID primitive.ObjectID `bson:"product_id,omitempty"`
	Rating    int                `bson:"rating,omitempty"`
	Comment   string             `bson:"comment,omitempty"`
}

func main() {
	// Establish a connection to the MongoDB instance
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(context.Background())

	// productCollection is a variable that holds a reference to the products collection in the subset database
	productCollection := client.Database("subset").Collection("products")
	// reviewCollection is a variable that holds a reference to the reviews collection in the subset database
	reviewCollection := client.Database("subset").Collection("reviews")

	// product is a variable that holds a Product struct with some values
	product := Product{
		ID:    primitive.NewObjectID(),
		Name:  "iPhone 12",
		Price: 799.99,
		Reviews1: []Review{
			{
				Rating:  5,
				Comment: "highly recommend!",
			},
			{
				Rating:  5,
				Comment: "it easy to read and watch videos.",
			},
		},
	}
	// productResult is a variable that holds the result of inserting the product into the productCollection
	productResult, err := productCollection.InsertOne(context.Background(), product)
	if err != nil {
		panic(err)
	}
	fmt.Println("ID of the inserted product", productResult)

	// review1, review2, and review3 are variables that hold Review structs with some values
	review1 := Review{
		ID:        primitive.NewObjectID(),
		ProductID: product.ID,
		Rating:    5,
		Comment:   "Awesome phone!",
	}
	review2 := Review{
		ID:        primitive.NewObjectID(),
		ProductID: product.ID,
		Rating:    4,
		Comment:   "Good camera quality.",
	}
	review3 := Review{
		ID:        primitive.NewObjectID(),
		ProductID: product.ID,
		Rating:    3,
		Comment:   "Battery life could be better.",
	}

	// reviewResult is a variable that holds the result of inserting the three reviews into the reviewCollection
	reviewResult, err := reviewCollection.InsertMany(context.Background(), []interface{}{review1, review2, review3})
	if err != nil {
		panic(err)
	}
	fmt.Println("IDs of the inserted product reviews", reviewResult.InsertedIDs)

	// Update one document in the productCollection
	_, err = productCollection.UpdateOne(
		context.Background(),
		// bson.M is a map that specifies the filter condition for finding the product by its ID
		bson.M{"_id": product.ID},
		// bson.D is a slice of key-value pairs that specifies the update operation to set the reviews field to a slice of Review structs
		bson.D{
			{"$set", bson.D{{"reviews", []Review{review1, review2, review3}}}},
		},
	)
	if err != nil {
		panic(err)
	}

	// pipeline is a variable that holds an aggregation pipeline to join the products and reviews collections
	pipeline := mongo.Pipeline{
		{{"$match", bson.D{{"_id", product.ID}}}}, // match the product by ID
		{{"$lookup", bson.D{ // lookup the reviews by product ID
			{"from", "reviews"},
			{"localField", "_id"},
			{"foreignField", "product_id"},
			{"as", "reviews2"},
		}}},
	}

	// cursor is a variable that holds the result of running the aggregation pipeline on the productCollection
	cursor, err := productCollection.Aggregate(context.Background(), pipeline)
	if err != nil {
		panic(err)
	}
	defer cursor.Close(context.Background())

	// Iterate over the cursor and print out each product with its reviews
	for cursor.Next(context.Background()) {
		var p Product
		err := cursor.Decode(&p)
		if err != nil {
			panic(err)
		}
		fmt.Println("Product:", p.Name, "Price:", p.Price)
		fmt.Println("Reviews1:")
		for _, r := range p.Reviews1 {
			fmt.Println("Rating:", r.Rating, "Comment:", r.Comment)
		}
		fmt.Println("Reviews2:")
		for _, r := range p.Reviews2 {
			fmt.Println("Rating:", r.Rating, "Comment:", r.Comment)
		}
	}

	// Drop the collection to clean up the data
	err = productCollection.Drop(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	err = reviewCollection.Drop(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}
