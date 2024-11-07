package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Blog represents a blog post
type Blog struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title     string             `bson:"title" json:"title"`
	Content   string             `bson:"content" json:"content"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	Category  string             `bson:"category" json:"category"`
	Photo     string             `bson:"photo" json:"photo"`
	Comments  []Comment          `bson:"comments" json:"comments"`
}

// Comment represents a comment on a blog post
type Comment struct {
	ID         string    `bson:"id" json:"id"`
	AuthorName string    `bson:"authorName" json:"authorName"`
	Content    string    `bson:"content" json:"content"`
	Photo      string    `bson:"photo" json:"photo"`
	CreatedAt  time.Time `bson:"createdAt" json:"createdAt"`
	Votes      int       `bson:"votes" json:"votes"`
}
type BlogFilters struct {
	Category string
	Title    string
}
