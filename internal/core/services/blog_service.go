package services

import (
	"0xfarms-backend/internal/core/domain"
	"0xfarms-backend/internal/ports"
	"errors"

	"github.com/google/uuid"
	"time"
)

// BlogService handles blog operations
type BlogService struct {
	db ports.MongoDB
}

// NewBlogService creates a new instance of the blog service
func NewBlogService(mongoDB ports.MongoDB) *BlogService {

	return &BlogService{db: mongoDB}
}

// AddBlog creates a new blog post
func (s *BlogService) AddBlog(blog *domain.Blog) (bool, error) {

	blog.CreatedAt = time.Now()
	blog.UpdatedAt = time.Now()

	ok, err := s.db.SaveBlog(blog)
	if err != nil {
		return false, err
	}
	return ok, nil
}

// GetBlog retrieves a single blog post by ID
func (s *BlogService) GetBlog(id string) (*domain.Blog, error) {

	blog, err := s.db.RetrieveBlog(id)
	if err != nil {
		return &domain.Blog{}, err
	}

	return blog, nil
}

// GetAllBlogs retrieves all blog posts
func (s *BlogService) GetAllBlogs() ([]domain.Blog, error) {

	blogs, err := s.db.RetrieveAllBlogs()
	if err != nil {
		return []domain.Blog{}, err
	}

	return blogs, nil
}

// UpdateBlog updates an existing blog post
func (s *BlogService) UpdateBlog(id string, updatedBlog *domain.Blog) (bool, error) {
	ok, err := s.db.UpdateBlog(id, updatedBlog)
	if err != nil {
		return false, err
	}
	return ok, nil
}

// DeleteBlog deletes a blog post
func (s *BlogService) DeleteBlog(id string) (bool, error) {

	ok, err := s.db.RemoveBlog(id)
	if err != nil {
		return false, err
	}

	return ok, nil
}

// AddComment adds a comment to a blog post
func (s *BlogService) AddComment(blogID, authorName, content, photo string) (*domain.Comment, error) {
	comment := &domain.Comment{
		ID:         uuid.New().String(),
		AuthorName: authorName,
		Content:    content,
		Photo:      photo,
		CreatedAt:  time.Now(),
		Votes:      0,
	}

	ok, err := s.db.AddComment(blogID, comment)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("failed to add comment")
	}

	return comment, nil
}

// VoteComment updates the vote count for a comment
func (s *BlogService) VoteComment(blogID string, commentID string, upvote bool) (bool, error) {
	ok, err := s.db.UpdateCommentVote(blogID, commentID, upvote)
	if err != nil {
		return false, err
	}
	return ok, nil
}
