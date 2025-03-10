package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prabhat-xs/library-management-backend/models"
	"github.com/prabhat-xs/library-management-backend/config"
)


func AddBook(c *gin.Context) {
    var book models.Books
    adminEmail, _ := c.Get("email")
    var admin models.User
    if err := config.DB.Where("email = ?", adminEmail).First(&admin).Error; err != nil || admin.Role != "Admin" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    book.LibID = admin.LibID

    if err := c.ShouldBindJSON(&book); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }


    var existing models.Books
    if err := config.DB.Where("isbn = ? AND lib_id = ?", book.ISBN, book.LibID).First(&existing).Error; err == nil {
        existing.Total_copies += book.Total_copies
        existing.Available_copies += book.Total_copies
        config.DB.Save(&existing)
        c.JSON(http.StatusOK, gin.H{"message": "Book copies updated"})
        return
    }

    book.Available_copies = book.Total_copies
    config.DB.Create(&book)
    c.JSON(http.StatusOK, gin.H{"message": "Book added successfully"})
}

func SearchBook(c *gin.Context) {
    title := c.Query("title")
    author := c.Query("author")
    publisher := c.Query("publisher")

    var books []models.Books
    query := config.DB.Model(&models.Books{})

    if title != "" {
        query = query.Where("title ILIKE ?", "%"+title+"%")
    }
    if author != "" {
        query = query.Where("authors ILIKE ?", "%"+author+"%")
    }
    if publisher != "" {
        query = query.Where("publisher ILIKE ?", "%"+publisher+"%")
    }

    query.Find(&books)
    
    c.JSON(http.StatusOK, books)
}