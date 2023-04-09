# MongoDB Model Design

## Polymorphic Pattern

### Introduction

MongoDB’s polymorphic schema allows for storing different types of documents in the same collection and adding new types of documents without affecting existing ones.

### Pros and Cons

Pros

- It is **convenient** to store documents of **different types in the same collection** without needing to pre-define the schema of the collection.
- Allowing the use of **the same set of APIs** to handle different types of documents makes for **simpler usage**.
- Allowing retrieval of **multiple types** of documents in **a single query** significantly **reduces query time**
- Storing different types of documents in the same collection can **reduce the storage space** needed for **indexing** and avoid creating multiple collections for different document types.

Cons

1. Switching between **different types** of data leads to **some performance loss**
2. There may be a problem when **querying a collection** that contains some documents with **a specific field** while others do not have that field
3. If you need to add **a new type** of document, you may need to modify **the program** to accommodate the new type.
4. Due to the need for **type conversion**, there may also be **an increase in code complexity** and the possibility of errors.

### Example

#### In Mongosh

 Connect to MongoDB

```bash
$ systemctl start mongod.service

$ mongosh

$ use model
```

Insert the following data.

```json
$ db.model.insertMany([{
    "name": "Savannah",
    "landAnimals": [
        {
            "name": "Lion"
        },
        {
            "name": "Zebra"
        },
        {
            "name": "Giraffe"
        }
    ]
},
{
    "name": "Coral Reef",
    "aquaticAnimals": [
        {
            "name": "Clownfish"
        },
        {
            "name": "Starfish"
        },
        {
            "name": "Seahorse"
        }
    ]
}])

```

delete the data
```json
$ db.model.deleteMany({ $or: [{ name: "Savannah" }, { name: "Coral Reef" }] })
```

#### In Golang

```go
package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
Define the structure for the data model,
including optional fields for land and aquatic animals
*/
type model struct {
	Name           string          `bson:"name"`
	LandAnimals    []LandAnimal    `bson:"landAnimals,omitempty"`
	AquaticAnimals []AquaticAnimal `bson:"aquaticAnimals,omitempty"`
}


// Define the structure for land animals
type LandAnimal struct {
	Name string `bson:"name"`
}

// Define the structure for aquatic animals
type AquaticAnimal struct {
	Name string `bson:"name"`
}

func main() {
	// Establish a connection to the MongoDB instance
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Verify that the connection was successful
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	// Get a handle to the "model" collection within the "model" database
	modelCollection := client.Database("model").Collection("polymorphic")

	// Insert two documents into the collection, each with a different set of animals
	_, err = modelCollection.InsertMany(context.Background(), []interface{}{
        // First document with land animals
		model{
			Name: "Savannah",
			LandAnimals: []LandAnimal{
				{Name: "Lion"},
				{Name: "Zebra"},
				{Name: "Giraffe"},
			},
		},
        // Second document with aquatic animals
		model{
			Name: "Coral Reef",
			AquaticAnimals: []AquaticAnimal{
				{Name: "Clownfish"},
				{Name: "Starfish"},
				{Name: "Seahorse"},
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Data inserted successfully!")

	// Query the collection for the documents that match the specified filter (by name)
    filter := bson.M{"name": bson.M{"$in": []string{"Savannah", "Coral Reef"}}}
	cur, err := modelCollection.Find(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(context.Background())
    
    /*
        Iterate over the cursor and decode each document into a "model" object,
        adding it to a slice of models
    */
	var models []model
	for cur.Next(context.Background()) {
		var model model
		err := cur.Decode(&model)
		if err != nil {
			log.Fatal(err)
		}
		models = append(models, model)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Data queried successfully: %v\n", models)

	// Delete all documents that match the specified filter
	_, err = modelCollection.DeleteMany(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Data deleted successfully!")
}
```





























