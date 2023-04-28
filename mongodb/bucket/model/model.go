package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// Temperature struct represents a temperature document
type Temperature struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Location  string             `bson:"location,omitempty"`
	Celsius   float64            `bson:"temperature,omitempty"`
	CreatedAt time.Time          `bson:"createdAt,omitempty"`
}
