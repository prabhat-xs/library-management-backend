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
        owner.POST("/books/add", controllers.AddBook)
        owner.POST("/requests/approve", controllers.ApproveIssueRequest)
        owner.POST("/create-admin", controllers.CreateAdminUser)
        owner.POST("/create-reader", controllers.CreateReaderUser)
    }
    admin := r.Group("v1/admin/").Use(middleware.AuthMiddleware( "Admin"))
    {
        admin.POST("/books/add", controllers.AddBook)
        admin.GET("/requests/all", controllers.ListRequests)
        admin.POST("/requests/approve", controllers.ApproveIssueRequest)
        admin.POST("/create-reader", controllers.CreateReaderUser)
        admin.GET("/books/search", controllers.SearchBook)
    }

    reader := r.Group("v1/reader/").Use(middleware.AuthMiddleware("Reader"))
    {
        reader.GET("/books/search", controllers.SearchBook)
        reader.POST("/books/requests", controllers.RaiseIssueRequest)
    }
}