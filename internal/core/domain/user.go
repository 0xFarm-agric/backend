package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrdinaryUser struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Email     string             `bson:"email"`
	Name      string             `bson:"name"`
	GoogleID  string             `bson:"google_id"` // Google unique ID
	CreatedAt time.Time          `bson:"created_at"`
	UserType  string             `bson:"user_type"` // "technician" or "ordinary_user"
}
