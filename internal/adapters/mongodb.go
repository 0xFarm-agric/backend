package adapters

import (
	"0xFarms-backend/internal/core/domain"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// BlogService handles blog operations
type DB struct {
	blogCollection     *mongo.Collection
	farmCollection     *mongo.Collection
	cropSpecCollection *mongo.Collection
	userCollection     *mongo.Collection
}

// NewBlogService creates a new instance of the blog service
func NewMongoAdapter(mongoURI string) (*DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, err
	}

	// Ping the database
	if err = client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	// Get the necessary collections
	blogCollection := client.Database("0xFarms").Collection("blogs")
	farmCollection := client.Database("0xFarms").Collection("vertical_farms")
	cropSpecCollection := client.Database("0xFarms").Collection("crop_specs")
	userCollection := client.Database("0xFarms").Collection("users")
	return &DB{
		blogCollection:     blogCollection,
		farmCollection:     farmCollection,
		cropSpecCollection: cropSpecCollection,
		userCollection:     userCollection,
	}, nil
}

// RegisterFarmTechnician registers a new farm technician with basic authentication
func (db *DB) RegisterFarmTechnician(technician *domain.FarmTechnician) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	technician.CreatedAt = time.Now()

	result, err := db.userCollection.InsertOne(ctx, technician)
	if err != nil {
		return "", err
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

// RegisterOrdinaryUser registers a new ordinary user using Google authentication
func (db *DB) RegisterOrdinaryUser(user *domain.OrdinaryUser) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user.CreatedAt = time.Now()
	user.UserType = "ordinary_user"

	result, err := db.userCollection.InsertOne(ctx, user)
	if err != nil {
		return "", err
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

// GetUser retrieves a user by ID
func (db *DB) GetUser(id string) (*domain.OrdinaryUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user domain.OrdinaryUser
	err = db.userCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// DeleteUser deletes a user by ID
func (db *DB) DeleteUser(id string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}

	result, err := db.userCollection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return false, err
	}

	if result.DeletedCount == 0 {
		return false, errors.New("user not found")
	}

	return true, nil
}

// AddBlog creates a new blog post
func (db *DB) SaveBlog(blog *domain.Blog) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := db.blogCollection.InsertOne(ctx, blog)
	if err != nil {
		return false, err
	}

	blog.ID = result.InsertedID.(primitive.ObjectID)
	return true, nil
}

// GetBlog retrieves a single blog post by ID
func (db *DB) RetrieveBlog(id string) (*domain.Blog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var blog domain.Blog
	err = db.blogCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&blog)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("blog not found")
		}
		return nil, err
	}

	return &blog, nil
}

// GetAllBlogs retrieves all blog posts
func (db *DB) RetrieveAllBlogs() ([]domain.Blog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := db.blogCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var blogs []domain.Blog
	if err = cursor.All(ctx, &blogs); err != nil {
		return nil, err
	}

	return blogs, nil
}

// UpdateBlog updates an existing blog post
func (db *DB) UpdateBlog(id string, updatedBlog *domain.Blog) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}

	updatedBlog.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"title":      updatedBlog.Title,
			"content":    updatedBlog.Content,
			"updated_at": updatedBlog.UpdatedAt,
		},
	}

	result, err := db.blogCollection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return false, err
	}

	if result.MatchedCount == 0 {
		return false, errors.New("blog not found")
	}

	return true, nil
}

// DeleteBlog deletes a blog post
func (s *DB) RemoveBlog(id string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}

	result, err := s.blogCollection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return false, err
	}

	if result.DeletedCount == 0 {
		return false, errors.New("blog not found")
	}

	return true, nil
}

// AddComment adds a comment to a blog post
func (db *DB) AddComment(blogID string, comment *domain.Comment) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return false, err
	}

	update := bson.M{
		"$push": bson.M{
			"comments": comment,
		},
	}

	result, err := db.blogCollection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		update,
	)

	if err != nil {
		return false, err
	}

	if result.MatchedCount == 0 {
		return false, errors.New("blog not found")
	}

	return true, nil
}

// UpdateCommentVote updates the vote count for a comment
func (db *DB) UpdateCommentVote(blogID string, commentID string, upvote bool) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return false, err
	}

	// Calculate vote increment
	voteIncrement := 1
	if !upvote {
		voteIncrement = -1
	}

	update := bson.M{
		"$inc": bson.M{
			"comments.$[elem].votes": voteIncrement,
		},
	}

	arrayFilters := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"elem.id": commentID},
		},
	})

	result, err := db.blogCollection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		update,
		arrayFilters,
	)

	if err != nil {
		return false, err
	}

	if result.MatchedCount == 0 {
		return false, errors.New("blog or comment not found")
	}

	return true, nil
}

// RetrieveBlogsByFilters retrieves blog posts based on the provided filters
func (db *DB) RetrieveBlogsByFilters(filters *domain.BlogFilters) ([]domain.Blog, error) {
	var blogs []domain.Blog
	// query := bson.M{}

	// if filters.Category != "" {
	// 	query["category"] = filters.Category
	// }

	// if filters.Title != "" {
	// 	query["title"] = bson.M{"$regex": primitive.Regex{Pattern: filters.Title, Options: "i"}}
	// }

	// err := db.blogCollection.Find(query).All(&blogs)
	// if err != nil {
	// 	return nil, err
	// }
	return blogs, nil
}

// CreateFarm creates a new vertical farm in the database
func (db *DB) CreateFarm(farm *domain.VerticalFarm) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := db.farmCollection.InsertOne(ctx, farm)
	if err != nil {
		return "", err
	}

	farm.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return farm.ID, nil
}

// GetFarm retrieves a vertical farm by ID
func (db *DB) GetFarm(id string) (*domain.VerticalFarm, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var farm domain.VerticalFarm
	err = db.farmCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&farm)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("farm not found")
		}
		return nil, err
	}

	return &farm, nil
}

// UpdateFarm updates an existing vertical farm
func (db *DB) UpdateFarm(id string, farm *domain.VerticalFarm) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}

	farm.LastUpdated = time.Now()

	update := bson.M{
		"$set": farm,
	}

	result, err := db.farmCollection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return false, err
	}

	if result.MatchedCount == 0 {
		return false, errors.New("farm not found")
	}

	return true, nil
}

// AddIoTReading adds a new IoT sensor reading to a farm
func (db *DB) AddIoTReading(farmID string, reading *domain.IoTReading) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(farmID)
	if err != nil {
		return err
	}

	update := bson.M{
		"$push": bson.M{
			"iot_data": reading,
		},
	}

	_, err = db.farmCollection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}

// GetCropSpecification retrieves a crop specification by type
func (db *DB) GetCropSpecification(cropType string) (domain.CropSpecification, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var cropSpec domain.CropSpecification
	err := db.cropSpecCollection.FindOne(ctx, bson.M{"name": cropType}).Decode(&cropSpec)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.CropSpecification{}, errors.New("crop specification not found")
		}
		return domain.CropSpecification{}, err
	}

	return cropSpec, nil
}
