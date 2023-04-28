package iot

import (
	"math/rand"
	"time"
)

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
			time.Sleep(1 * time.Second)
		}
	}()

	// Return the channel
	return readings
}
