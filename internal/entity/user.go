package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	FirstName    string             `json:"firstName" bson:"firstName"`
	LastName     string             `json:"lastName" bson:"lastName"`
	Email        string             `json:"email" bson:"email"`
	Password     string             `json:"password" bson:"password"`
	RegisteredAt time.Time          `json:"registeredAt" bson:"registeredAt"`
	LastVisitAt  time.Time          `json:"lastVisitAt" bson:"lastVisitAt"`
}
