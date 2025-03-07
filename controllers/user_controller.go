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

func RegisterUser(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := config.DB.Where("Email=?", user.Email).First(&user).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User with this email already exists!"})
		return
	}

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

	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	
	user.Password = string(hashed)

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(http.StatusCreated, user)

}

func LoginUser(c *gin.Context) {
	var user models.User
	var fetchedUser models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := config.DB.First(&fetchedUser, user.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User does not exits!",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(fetchedUser.Password), []byte(user.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid password",
		})
	}

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



