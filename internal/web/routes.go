package web

import (
	"0xFarms-backend/internal/web/handlers"

	"github.com/gin-gonic/gin"
)

// SetupAPIRoutes sets up the API routes for the application.
func SetupAPIRoutes(r *gin.Engine, blogHandler *handlers.BlogHandler, farmHandler *handlers.FarmHandler) {

	r.GET("/blog/save", blogHandler.SaveBlog)
	r.GET("/blog/:id/get_one_blog", blogHandler.GetABlog)
	r.GET("/blog/all_blog", blogHandler.GetAllBlogs)
}
