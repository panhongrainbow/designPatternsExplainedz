package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Product struct represents a product in a store
type Product struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Name    string             `bson:"name,omitempty"`
	Price   float64            `bson:"price,omitempty"`
	Reviews []Review           `bson:"reviews,omitempty"`
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

	//
	productCollection.UpdateOne(
		context.Background(),
		// bson.M is a map that specifies the filter condition for finding the product by its ID
		bson.M{"_id": product.ID},
		// bson.D is a slice of key-value pairs that specifies the update operation to set the reviews field to a slice of Review structs
		bson.D{
			{"$set", bson.D{{"reviews", []Review{review1, review2, review3}}}},
			{"$push", bson.D{{"review_ids", bson.D{{"$each", reviewResult.InsertedIDs}}}}},
		},
	)

	// 查詢產品
	var result Product
	err = productCollection.FindOne(context.Background(), bson.M{"_id": product.ID}).Decode(&result)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)

	pipeline := mongo.Pipeline{{{"$lookup", bson.D{{"from", "reviews"}, {"localField", "review_ids"}, {"foreignField", "_id"}, {"as", "reviews"}}}}}
	cursor, err := productCollection.Aggregate(context.Background(), pipeline)
	if err != nil {
		panic(err)
	}
	defer cursor.Close(context.Background())

	var products []Product
	if err = cursor.All(context.Background(), &products); err != nil {
		panic(err)
	}
	fmt.Println(products)
}
