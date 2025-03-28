package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/prabhat-xs/library-management-backend/controllers"
	"github.com/prabhat-xs/library-management-backend/middleware"
	"github.com/prabhat-xs/library-management-backend/utils"
)

func SetupRoutes(r *gin.Engine) {
	auth := r.Group("v1/auth/")
	{
		auth.POST("/signup", controllers.Signup)
		auth.POST("/login", controllers.Login)
	}

	owner := r.Group("v1/owner/").Use(middleware.AuthMiddleware(utils.ValidateJWT, "Owner"))
	{
		owner.POST("/password", controllers.UpdatePassword)
		owner.POST("/create-admin", controllers.CreateAdminUser)
		owner.GET("/users", controllers.ListAllUsers)
		owner.PATCH("/users", controllers.UpdateUser)
		owner.GET("/logout", controllers.Logout)
	}

	admin := r.Group("v1/admin/").Use(middleware.AuthMiddleware(utils.ValidateJWT, "Admin", "Owner"))
	{
		admin.POST("/password", controllers.UpdatePassword)
		admin.POST("/create-reader", controllers.CreateReaderUser)
		admin.PATCH("/users/", controllers.UpdateUser)
		admin.GET("/books/search", controllers.SearchBook)
		admin.POST("/books/add", controllers.AddBook)
		admin.PATCH("/books/:isbn", controllers.UpdateBook)
		admin.DELETE("/books/:isbn", controllers.DeleteBook)
		admin.GET("/requests/all", controllers.ListRequests)
		admin.POST("/requests/process", controllers.ProcessRequest)
		admin.GET("/users", controllers.ListAllUsers)
		admin.GET("/logout", controllers.Logout)
		admin.GET("/books", controllers.ShowBooks)
	}

	reader := r.Group("v1/reader/").Use(middleware.AuthMiddleware(utils.ValidateJWT, "Reader"))
	{
		reader.POST("/password", controllers.UpdatePassword)
		reader.POST("/books/search", controllers.SearchBook)
		reader.GET("/books", controllers.ShowBooks)
		reader.POST("/books/requests", controllers.RaiseBookRequest)
		reader.GET("/books/requests", controllers.MyRequests)
		reader.GET("/books/my", controllers.BorrowedBooks)
		reader.GET("/logout", controllers.Logout)
	}
	r.GET("/books", controllers.ShowBooks)
}
