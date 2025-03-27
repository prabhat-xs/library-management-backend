package controllers

import (
	// "log"
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
		Name:           input.Name,
		Email:          input.Email,
		Password:       string(hashedPassword),
		Contact_number: input.ContactNumber,
		Role:           "Owner",
		LibID:          library.LibID,
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
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
	token, _ := utils.GenerateJWT(user.ID, user.LibID, user.Email, user.Role)

	// FOR SETTING SECURE SITE
	prodMode := os.Getenv("PROD_MODE") == "true"
	c.SetCookie("token", token, 3600*72, "/", "localhost", prodMode, true)

	var returnUser = models.User{
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user":    returnUser,
		"token": token,
	})
}

// CREATING ADMIN USER
func CreateAdminUser(c *gin.Context) {
	var input struct {
		Name          string `json:"name" binding:"required"`
		Email         string `json:"email" binding:"required"`
		Password      string `json:"password" binding:"required"`
		ContactNumber string `json:"contactNumber" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// VALIDATING OWNER USING COOKIES
	ownerEmail, _ := c.Get("email")
	var owner models.User
	if err := config.DB.Where("email = ?", ownerEmail).First(&owner).Error; err != nil || owner.Role != "Owner" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	adminUser := models.User{
		Name:           input.Name,
		Email:          input.Email,
		Password:       string(hashedPassword),
		Contact_number: input.ContactNumber,
		Role:           "Admin",
		LibID:          owner.LibID,
	}

	if err := config.DB.Create(&adminUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create admin user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Admin user created successfully"})
}

// CREATING READER USER
func CreateReaderUser(c *gin.Context) {
	var input struct {
		Name          string `json:"name" binding:"required"`
		Email         string `json:"email" binding:"required"`
		Password      string `json:"password" binding:"required"`
		ContactNumber string `json:"contactNumber" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// VALIDATING ADMIN USER USING COOKIES
	adminEmail, _ := c.Get("email")
	var admin models.User
	if err := config.DB.Where("email = ?", adminEmail).First(&admin).Error; err != nil || !(admin.Role == "Admin" || admin.Role == "Owner") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	readerUser := models.User{
		Name:           input.Name,
		Email:          input.Email,
		Password:       string(hashedPassword),
		Contact_number: input.ContactNumber,
		Role:           "Reader",
		LibID:          admin.LibID,
	}

	if err := config.DB.Create(&readerUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create reader user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reader user created successfully"})
}

// UPDATE USER PASSWORD
func UpdatePassword(c *gin.Context) {
	var input struct {
		OldPassword string `binding:"required"`
		NewPassword string `binding:"required"`
	}

	// INPUT VALIDATION
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	email, _ := c.Get("email")
	// FETCHING USER INFORMATION
	var user models.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// EXISTING PASSWORD VERIFICATION
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.OldPassword)) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong Password"})
		return
	}

	// NEW PASSWORD HASHING AND SAVING
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost) //; err != nil {
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Bcrypt failed to generate password!",
		})
		return
	}

	user.Password = string(hashedPassword)
	config.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{
		"message": "Password updated successfully!",
	})

}

// LIST ALL READERS
func ListAllUsers(c *gin.Context) {
	libid, _ := c.Get("libid")
	var users []models.User
	if err := config.DB.Where("lib_id = ? ", libid).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"admins": users,
	})

}

// UPDATE USER DETAILS
func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	libid, _ := c.Get("libid")
	role, _ := c.Get("role")

	var input struct {
		Name          string
		Role          string
		Email         string
		Password      string
		ContactNumber string
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Nothing to update!",
		})
		return
	}
	var user models.User
	if err := config.DB.Where("lib_id = ?", libid).First(&user, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if (role == "Admin" && user.Role != "Reader") || role == "Reader" || (role == "Admin" && input.Role != "") {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Admin can't perform this operation!",
		})
		return
	}
	flag := true
	if input.Role != "" {
		user.Role = input.Role
		flag = false
	}
	if input.Name != "" {
		user.Name = input.Name
		flag = false
	}
	if input.Email != "" {
		user.Email = input.Email
		flag = false
	}
	if input.Password != "" {
		hashed, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		user.Password = string(hashed)
		flag = false
	}
	if input.ContactNumber != "" {
		user.Contact_number = input.ContactNumber
		flag = false
	}

	if flag {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Nothing to update",
		})
		return
	}

	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "User details updated successfully!",
	})
}

// DELETE USER
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	role, _ := c.Get("role")
	libid, _ := c.Get("libid")

	var user models.User
	if err := config.DB.Where("libid = ?", libid).First(&user, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User does not exists",
		})
	}
	if (role == "admin" && user.Role != "reader") || role == "reader" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized!",
		})
	}

	if err := config.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully!",
	})
}

// LOGOUT FUNCTIONALITY
func Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
