package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/prabhat-xs/library-management-backend/config"
	"github.com/prabhat-xs/library-management-backend/models"
	"golang.org/x/crypto/bcrypt"
)

var library models.Library

// User registration controller
func RegisterUser(c *gin.Context) {
	var user models.User

	// Biding user model with request data
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Checking if the user already exists
	if err := config.DB.Where("Email=?", user.Email).First(&user).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User with this email already exists!"})
		return
	}

	// If new user is of type owner, new library should be created
	if user.Role == "owner" {
		if user.LibName == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":       "Library name missing!",
				"description": "New owners must provide a name for their library!",
			})
			return
		}
		library.Name = user.LibName

		if err := config.DB.Create(&library).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		user.LibID = library.ID

	}

	// Encrypting user password
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	
	user.Password = string(hashed)

	// User creation in database
	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(http.StatusCreated, user)

}

// User login controller
func LoginUser(c *gin.Context) {
	var user models.User
	var fetchedUser models.User

	// Data binding and validation 
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Checking whether user is registered
	if err := config.DB.First(&fetchedUser, user.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User does not exits!",
		})
	}

	// Verifying password
	if err := bcrypt.CompareHashAndPassword([]byte(fetchedUser.Password), []byte(user.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid password",
		})
	}

	// JWT Token
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"user": fetchedUser.ID,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})

}



