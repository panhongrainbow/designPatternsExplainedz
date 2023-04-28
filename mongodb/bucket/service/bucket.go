package main

import (
	"context"
	"fmt"
	"github.com/panhongrainbow/designPatternsExplainedz/mongodb/bucket/iot"
	"github.com/panhongrainbow/designPatternsExplainedz/mongodb/bucket/model"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// SwitchBucket function returns a string representing the bucket
func SwitchBucket(temp float64) string {
	switch {
	case temp < 0:
		return "lt_0"
	case temp < 10:
		return "lt_10"
	case temp < 20:
		return "lt_20"
	default:
		return "gte_20"
	}
}

// main function
func main() {
	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}

	// Disconnect from MongoDB after main function finishes
	defer func() {
		_ = client.Disconnect(context.TODO())
	}()

	// Select database and collection
	db := client.Database("iot")

	// Create a sensor object
	sensor := iot.Sensor{
		Location: "laboratory",
	}

	// Read temperature data from sensor
	for temp := range sensor.Readings() {
		// Create a temperature document
		t := model.Temperature{
			Location:  sensor.Location,
			Celsius:   temp,
			CreatedAt: time.Now(),
		}

		fmt.Println(t.Celsius, t.CreatedAt)

		// Determine the bucket
		bucket := SwitchBucket(t.Celsius)

		// Insert the temperature document into the appropriate bucket
		_, err = db.Collection("temp_"+bucket).InsertOne(context.TODO(), t)
		if err != nil {
			panic(err)
		}
	}
}
