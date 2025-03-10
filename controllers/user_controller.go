package controllers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/prabhat-xs/library-management-backend/config"
	"github.com/prabhat-xs/library-management-backend/models"
	"github.com/prabhat-xs/library-management-backend/utils"
	"golang.org/x/crypto/bcrypt"
)

var library models.Library

// CREATING OWNER ACCOUNT
func Signup(c *gin.Context) {
	var input struct {
		Name, Email, Password, ContactNumber, LibraryName string `binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// CHECKING IF LIBRARY WITH SAME NAME EXISTS OR NOT
	var existing models.Library
	if err := config.DB.Where("name = ?", input.LibraryName).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Library already exists"})
		return
	}

	// LIBRARY CREATION
	library := models.Library{Name: input.LibraryName}
	config.DB.Create(&library)

	// PASSWORD HASHING
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	user := models.User{
		Name:          input.Name,
		Email:         input.Email,
		Password:      string(hashedPassword),
		Contact_number: input.ContactNumber,
		Role:          "Owner",
		LibID:         library.ID,
	}
	config.DB.Create(&user)
	c.JSON(http.StatusOK, gin.H{"message": "Owner account created successfully"})
}

// USER LOGIN
func Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// CHECKING IF USER EXISTS
	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// PASSWORD VERIFICATION
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// TOKEN GENERATION
	token, _ := utils.GenerateJWT(user.Email, user.Role,user.ID)

	// FOR SETTING SECURE SITE
	prodMode := os.Getenv("PROD_MODE")=="true"
	c.SetCookie("token", token, 3600*72, "/", "localhost",prodMode, true) 

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}
