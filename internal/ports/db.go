package ports

import "0xFarms-backend/internal/core/domain"

type MongoDB interface {
	SaveBlog(blog *domain.Blog) (bool, error)
	RetrieveBlog(id string) (*domain.Blog, error)
	RetrieveAllBlogs() ([]domain.Blog, error)
	UpdateBlog(id string, updatedBlog *domain.Blog) (bool, error)
	RemoveBlog(id string) (bool, error)
	AddComment(blogID string, comment *domain.Comment) (bool, error)
	UpdateCommentVote(blogID string, commentID string, upvote bool) (bool, error)
	RetrieveBlogsByFilters(filters *domain.BlogFilters) ([]domain.Blog, error)
	// New farm-related operations
	CreateFarm(farm *domain.VerticalFarm) (string, error)
	GetFarm(id string) (*domain.VerticalFarm, error)
	UpdateFarm(id string, farm *domain.VerticalFarm) (bool, error)
	AddIoTReading(farmID string, reading *domain.IoTReading) error
	GetCropSpecification(cropType string) (domain.CropSpecification, error)
}
