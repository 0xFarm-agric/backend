package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FarmTechnician struct {
	ID                     primitive.ObjectID `bson:"_id,omitempty"`
	Name                   string             `bson:"name"`
	IDNumber               string             `bson:"id_number"`
	IDPhoto                string             `bson:"id_photo"` // path or URL to photo
	School                 string             `bson:"school"`
	Certificate            string             `bson:"certificate"`
	CourseOfStudy          string             `bson:"course_of_study"`
	FoodCropSpecialization string             `bson:"food_crop_specialization"`
	CreatedAt              time.Time          `bson:"created_at"`
}
