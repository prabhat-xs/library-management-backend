package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/prabhat-xs/library-management-backend/controllers"
	"github.com/prabhat-xs/library-management-backend/middleware"
)

func SetupRoutes(r *gin.Engine) {

	auth := r.Group("v1/auth/")
	{
		auth.POST("/signup", controllers.Signup)
		auth.POST("/login", controllers.Login)
		auth.GET("/logout", controllers.Logout)
	}

	owner := r.Group("v1/owner/").Use(middleware.AuthMiddleware("Owner"))
	{
		owner.POST("/password", controllers.UpdatePassword)
		owner.POST("/create-admin", controllers.CreateAdminUser)
	}

	admin := r.Group("v1/admin/").Use(middleware.AuthMiddleware("Admin", "Owner"))
	{
		admin.POST("/password", controllers.UpdatePassword)
		admin.POST("/create-reader", controllers.CreateReaderUser)
		admin.GET("/books/search", controllers.SearchBook)
		admin.POST("/books/add", controllers.AddBook)
		admin.PATCH("/books/:isbn", controllers.UpdateBook)
		admin.DELETE("/books/:isbn", controllers.DeleteBook)
		admin.GET("/requests/all", controllers.ListRequests)
		admin.POST("/requests/process", controllers.ProcessRequest)
	}

	reader := r.Group("v1/reader/").Use(middleware.AuthMiddleware("Reader"))
	{
		reader.POST("/password", controllers.UpdatePassword)
		reader.GET("/books/search", controllers.SearchBook)
		reader.POST("/books/requests", controllers.RaiseBookRequest)
	}
}
