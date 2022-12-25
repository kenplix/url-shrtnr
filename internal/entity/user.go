package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User entity information
//
//	@Description	User entity information
type User struct {
	ID        primitive.ObjectID `json:"id" example:"63a75a2574ef628a127ee972"`
	Username  string             `json:"username" example:"kenplix"`
	Email     string             `json:"email" example:"tolstoi.job@gmail.com"`
	CreatedAt time.Time          `json:"createdAt" example:"2022-12-24T21:49:33.072726+02:00"`
	// UpdatedAt is a date of last user personal information modification
	UpdatedAt time.Time `json:"updatedAt" example:"2022-12-24T21:58:27.072726+02:00"`
	// SuspendedAt is a date when user was suspended through certain reasons (optional)
	SuspendedAt *time.Time `json:"suspendedAt,omitempty" example:"2022-12-25T14:25:58.821989+02:00"`
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
