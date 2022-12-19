package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID          primitive.ObjectID `json:"id"`
	Username    string             `json:"username"`
	Email       string             `json:"email"`
	CreatedAt   time.Time          `json:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt"`
	SuspendedAt *time.Time         `json:"suspendedAt,omitempty"`
}

type UserModel struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username     string             `json:"username" bson:"username"`
	Email        string             `json:"email" bson:"email"`
	PasswordHash string             `json:"passwordHash" bson:"passwordHash"`
	CreatedAt    time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt    time.Time          `json:"updatedAt" bson:"updatedAt"`
	SuspendedAt  *time.Time         `json:"suspendedAt,omitempty" bson:"suspendedAt"`
}

func (u UserModel) Filter() User {
	return User{
		ID:          u.ID,
		Username:    u.Username,
		Email:       u.Email,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
		SuspendedAt: u.SuspendedAt,
	}
}
