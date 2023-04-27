package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math/rand"
	"time"
)

func main() {
	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(context.TODO())

	// Select database and collection
	db := client.Database("iot")

	// Create a temperature struct
	type Temp struct {
		ID        primitive.ObjectID `bson:"_id,omitempty"`
		Location  string             `bson:"location"`
		Temp      float64            `bson:"temp"`
		CreatedAt time.Time          `bson:"createdAt"`
	}

	// Create a sensor object
	sensor := Sensor{}

	// Create 50 temperature documents
	var count int

	// Read temperature data from sensor
	for temp := range sensor.Readings() {
		// Create a temperature document
		t := Temp{
			Location:  sensor.Location,
			Temp:      temp,
			CreatedAt: time.Now(),
		}

		// Determine the bucket
		bucket := getBucket(t.Temp)

		// Insert the temperature document into the appropriate bucket
		_, err = db.Collection("temp_"+bucket).InsertOne(context.TODO(), t)
		if err != nil {
			panic(err)
		}

		/*
			Count the number of documents created.
			If 50 documents have been created, stop
		*/
		count++
		if count == 50 {
			break
		}

	}

	fmt.Println(getTemps(10))

}

// Sensor struct represents a sensor
type Sensor struct {
	ID       string `bson:"_id"`
	Location string
}

// Readings method returns a channel of float64 values
func (s *Sensor) Readings() <-chan float64 {
	// Create a channel of float64 values
	readings := make(chan float64)

	// Start a goroutine to generate random temperature readings
	go func() {
		for {
			reading := rand.Float64()*60 - 20
			readings <- reading
			// time.Sleep(500 * time.Microsecond)
		}
	}()

	return readings
}

// C
func getBucket(temp float64) string {
	if temp < 0 {
		return "lt_0"
	} else if temp < 10 {
		return "lt_10"
	} else if temp < 20 {
		return "lt_20"
	} else {
		return "gte_20"
	}
}

// Temp struct represents a temperature document
type Temp struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `bson:"name,omitempty"`
	Location string             `bson:"location,omitempty"`
	Temp     float64            `bson:"temp,omitempty"`
	Time     time.Time          `bson:"time,omitempty"`
}

// getTemps function returns a slice of Temp structs
func getTemps(limit int) ([]Temp, error) {
	// connect to MongoDB
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(context.TODO())

	// Select database and collection
	db := client.Database("iot")

	ctx := context.TODO()
	collection := db.Collection("temp_gte_20")

	pipeline := mongo.Pipeline{
		bson.D{{"$match", bson.D{{"temp", bson.D{{"$gt", 1}}}}}},
		{{"$bucket", bson.M{ // lookup the reviews by product ID
			"groupBy":    "$temp",
			"boundaries": []int{0, 10, 20, 30},
			"default":    "other",
			"output": bson.M{
				"count":     bson.M{"$sum": 1},
				"documents": bson.M{"$push": "$$ROOT"},
			},
		}}},
		{{"$limit", limit}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	var temps []Temp
	for cursor.Next(ctx) {
		var t struct {
			Documents []Temp `bson:"documents"`
		}
		err := cursor.Decode(&t)
		if err != nil {
			cursor.Close(ctx)
			return nil, err
		}
		temps = append(temps, t.Documents...)
	}
	if err := cursor.Err(); err != nil {
		cursor.Close(ctx)
		return nil, err
	}
	if err := cursor.Close(ctx); err != nil {
		return nil, err
	}
	return temps, nil
}
