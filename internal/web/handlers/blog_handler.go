package handlers

import (
	"0xFarms-backend/internal/core/domain"
	"0xFarms-backend/internal/core/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BlogHandler struct {
	blogService *services.BlogService
}

// NewCommitHandler creates a new instance of CommitHandler with the given services
func NewBlogHandler(blogService *services.BlogService) *BlogHandler {
	return &BlogHandler{
		blogService: blogService,
	}
}

func (h *BlogHandler) SaveBlog(c *gin.Context) {

	var blog domain.Blog

	// Retrieve the top N commit authors from the repository service
	ok, err := h.blogService.AddBlog(&blog)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusCode": http.StatusInternalServerError, "message": "Server error"})
		return
	}
	if ok {
		c.JSON(http.StatusOK, gin.H{"statusCode": http.StatusOK, "data": ok})
		return
	}
}

func (h *BlogHandler) GetABlog(c *gin.Context) {

	id := c.Param("id")

	// Retrieve the top N commit authors from the repository service
	blog, err := h.blogService.GetBlog(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusCode": http.StatusInternalServerError, "message": "Server error"})
		return
	}
	if blog == nil {
		c.JSON(http.StatusNotFound, gin.H{"statusCode": http.StatusNotFound, "message": "No blogs found"})
		return
	}

	// Return the top authors as JSON
	c.JSON(http.StatusOK, blog)
}

func (h *BlogHandler) GetAllBlogs(c *gin.Context) {

	// Parse the "n" parameter to determine the number of top authors
	n, err := strconv.Atoi(c.Param("n"))
	if err != nil || n <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid number of authors"})
		return
	}

	// Retrieve the top N commit authors from the repository service
	blogs, err := h.blogService.GetAllBlogs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusCode": http.StatusInternalServerError, "message": "Server error"})
		return
	}
	if len(blogs) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"statusCode": http.StatusNotFound, "message": "No blogs found"})
		return
	}

	// Return the top authors as JSON
	c.JSON(http.StatusOK, blogs)
}
func (h *BlogHandler) DeleteBlog(c *gin.Context) {
	id := c.Param("id")

	// Validate input
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": "ID is required"})
		return
	}

	// Remove all commits for the specified repository
	ok, err := h.blogService.DeleteBlog(id)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"statusCode": http.StatusInternalServerError, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"statusCode": http.StatusOK, "message": "Blog commits removed successfully"})
}
