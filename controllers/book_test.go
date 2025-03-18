package controllers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/prabhat-xs/library-management-backend/config"
	"github.com/prabhat-xs/library-management-backend/models"
	"github.com/prabhat-xs/library-management-backend/testutils"
	"github.com/stretchr/testify/assert"
)

func TestAddBook(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testutils.SetupTestDB()

	if config.DB == nil {
		t.Fatalf("config.DB is nil! Database not initialized.")
	}

	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Set("libid", uint(1))
		c.Next()
	})

	router.POST("/books/add", AddBook)

	lib := models.Library{Name: "Test Library"}
	result := config.DB.Create(&lib)

	if result.Error != nil {
		t.Fatalf("❌ Failed to insert test library: %v", result.Error)
	}

	bookPayload := `{
		"title": "Go Programming",
		"authors": "John Doe",
		"publisher": "Tech Press",
		"version": "1st",
		"isbn": 123456,
		"total_copies": 5
	}`

	req, _ := http.NewRequest("POST", "/books/add", bytes.NewBuffer([]byte(bookPayload)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected HTTP 200 OK response")

	var book models.Books
	err := config.DB.First(&book, "isbn = ?", 123456).Error
	assert.Nil(t, err, "❌ Book should exist in the database")
	assert.Equal(t, "Go Programming", book.Title)
}
func TestAddBook_BookExists(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testutils.SetupTestDB()

	if config.DB == nil {
		t.Fatalf("config.DB is nil! Database not initialized.")
	}

	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Set("libid", uint(1))
		c.Next()
	})
	booka := models.Books{
		ISBN:             123456,
		Title:            "Go Programming",
		Authors:          "John Doe",
		Publisher:        "Tech Press",
		Version:          "1st",
		LibID:            1,
		Total_copies:     5,
		Available_copies: 5,
	}
	config.DB.Create(&booka)
	router.POST("/books/add", AddBook)

	lib := models.Library{Name: "Test Library"}
	result := config.DB.Create(&lib)

	if result.Error != nil {
		t.Fatalf("❌ Failed to insert test library: %v", result.Error)
	}

	bookPayload := `{
		"title": "Go Programming",
		"authors": "John Doe",
		"publisher": "Tech Press",
		"version": "1st",
		"isbn": 123456,
		"total_copies": 5
	}`

	req, _ := http.NewRequest("POST", "/books/add", bytes.NewBuffer([]byte(bookPayload)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Book copies updated")

}

func TestAddBook_LibraryNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testutils.SetupTestDB()

	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set("libid", uint(999))
		c.Next()
	})
	router.POST("/books/add", AddBook)

	bookPayload := `{
		"title": "Go Programming",
		"authors": "John Doe",
		"publisher": "Tech Press",
		"version": "1st",
		"isbn": 123456,
		"total_copies": 5
	}`
	req, _ := http.NewRequest("POST", "/books/add", bytes.NewBuffer([]byte(bookPayload)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}

func TestAddBook_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testutils.SetupTestDB()

	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set("libid", uint(1))
		c.Next()
	})
	router.POST("/books/add", AddBook)

	bookPayload := `{invalid_json}`
	req, _ := http.NewRequest("POST", "/books/add", bytes.NewBuffer([]byte(bookPayload)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}

func TestSearchBook(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testutils.SetupTestDB()

	router := gin.Default()

	book := models.Books{
		ISBN:             123456,
		Title:            "Go Programming",
		Authors:          "John Doe",
		Publisher:        "Tech Press",
		Version:          "1st",
		LibID:            1,
		Total_copies:     5,
		Available_copies: 5,
	}
	config.DB.Create(&book)

	router.Use(func(c *gin.Context) {
		c.Set("libid", uint(1))
		c.Next()
	})

	router.POST("/books/search", SearchBook)

	searchPayload := `{
		"title": "Go Programming"
	}`

	req, _ := http.NewRequest("POST", "/books/search", bytes.NewBuffer([]byte(searchPayload)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Go Programming")
}
func TestSearchBook_BadJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testutils.SetupTestDB()

	router := gin.Default()

	book := models.Books{
		ISBN:             123456,
		Title:            "Go Programming",
		Authors:          "John Doe",
		Publisher:        "Tech Press",
		Version:          "1st",
		LibID:            1,
		Total_copies:     5,
		Available_copies: 5,
	}
	config.DB.Create(&book)

	router.Use(func(c *gin.Context) {
		c.Set("libid", uint(1))
		c.Next()
	})

	router.POST("/books/search", SearchBook)

	searchPayload := `{
	}`

	req, _ := http.NewRequest("POST", "/books/search", bytes.NewBuffer([]byte(searchPayload)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

}

func TestSearchBook_NoSearchFields(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testutils.SetupTestDB()

	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set("libid", uint(1)) 
		c.Next()
	})
	router.POST("/books/search", SearchBook)

	searchPayload := `{}`
	req, _ := http.NewRequest("POST", "/books/search", bytes.NewBuffer([]byte(searchPayload)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "At least one search field is required")
}

func TestSearchBook_BookNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testutils.SetupTestDB()

	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set("libid", uint(1)) 
		c.Next()
	})
	router.POST("/books/search", SearchBook)

	searchPayload := `{"title": "Nonexistent Book"}`
	req, _ := http.NewRequest("POST", "/books/search", bytes.NewBuffer([]byte(searchPayload)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}

func TestUpdateBook(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testutils.SetupTestDB()

	router := gin.Default()

	book := models.Books{
		ISBN:             123456,
		Title:            "Go Programming",
		Authors:          "John Doe",
		Publisher:        "Tech Press",
		Version:          "1st",
		LibID:            1,
		Total_copies:     5,
		Available_copies: 5,
	}
	config.DB.Create(&book)

	router.Use(func(c *gin.Context) {
		c.Set("libid", uint(1))
		c.Next()
	})

	router.PUT("/books/update", UpdateBook)

	updatePayload := `{
		"isbn": 123456,
		"title": "Advanced Go Programming",
		"authors": "Jane Doe",
		"publisher": "Tech Press 2nd Edition",
		"version": "2nd",
		"total_copies": 10
	}`

	req, _ := http.NewRequest("PUT", "/books/update", bytes.NewBuffer([]byte(updatePayload)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Book updated successfully")

	var updatedBook models.Books
	config.DB.Where("isbn = ?", 123456).First(&updatedBook)
	assert.Equal(t, "Advanced Go Programming", updatedBook.Title)
	assert.Equal(t, "Jane Doe", updatedBook.Authors)
	assert.Equal(t, "Tech Press 2nd Edition", updatedBook.Publisher)
	assert.Equal(t, "2nd", updatedBook.Version)
	assert.Equal(t, uint(10), updatedBook.Total_copies)
}
func TestUpdateBook_NothingToUpdate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testutils.SetupTestDB()

	router := gin.Default()

	book := models.Books{
		ISBN:             123456,
		Title:            "Go Programming",
		Authors:          "John Doe",
		Publisher:        "Tech Press",
		Version:          "1st",
		LibID:            1,
		Total_copies:     5,
		Available_copies: 5,
	}
	config.DB.Create(&book)

	router.Use(func(c *gin.Context) {
		c.Set("libid", uint(1))
		c.Next()
	})

	router.PUT("/books/update", UpdateBook)

	updatePayload := `{
		"ISBN":123456
	}`

	req, _ := http.NewRequest("PUT", "/books/update", bytes.NewBuffer([]byte(updatePayload)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Nothing to update")

}

func TestUpdateBook_BookNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testutils.SetupTestDB()

	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set("libid", uint(1)) 
		c.Next()
	})
	router.PUT("/books/update", UpdateBook)

	updatePayload := `{"isbn": 999999}`
	req, _ := http.NewRequest("PUT", "/books/update", bytes.NewBuffer([]byte(updatePayload)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Book not found")
}

func TestUpdateBook_AvailableCopiesExceedTotalCopies(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testutils.SetupTestDB()

	book := models.Books{
		ISBN:             123456,
		Title:            "Go Programming",
		Authors:          "John Doe",
		Publisher:        "Tech Press",
		Version:          "1st",
		LibID:            1,
		Total_copies:     5,
		Available_copies: 5,
	}
	config.DB.Create(&book)

	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set("libid", uint(1)) 
		c.Next()
	})
	router.PUT("/books/update", UpdateBook)

	updatePayload := `{"isbn": 123456, "available_copies": 10}`
	req, _ := http.NewRequest("PUT", "/books/update", bytes.NewBuffer([]byte(updatePayload)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Available copies can not be more than total copies")
}

func TestDeleteBook(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testutils.SetupTestDB()

	router := gin.Default()

	admin := models.User{
		Name:  "Admin User",
		Email: "admin@example.com",
		Role:  "Admin",
		LibID: 1,
	}
	config.DB.Create(&admin)

	book := models.Books{
		ISBN:             123456,
		Title:            "Go Programming",
		Authors:          "John Doe",
		Publisher:        "Tech Press",
		Version:          "1st",
		LibID:            1,
		Total_copies:     5,
		Available_copies: 5,
	}
	config.DB.Create(&book)

	router.Use(func(c *gin.Context) {
		c.Set("email", "admin@example.com")
		c.Next()
	})

	isbnStr := strconv.Itoa(123456)
	router.DELETE("/books/delete/:isbn", DeleteBook)

	req, _ := http.NewRequest("DELETE", "/books/delete/"+isbnStr, nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Book deleted successfully")

	var deletedBook models.Books
	err := config.DB.Where("isbn = ?", 123456).First(&deletedBook).Error
	assert.NotNil(t, err)
}
func TestDeleteBook_NotAvailable(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testutils.SetupTestDB()

	router := gin.Default()

	admin := models.User{
		Name:  "Admin User",
		Email: "admin@example.com",
		Role:  "Admin",
		LibID: 1,
	}
	config.DB.Create(&admin)

	book := models.Books{
		ISBN:             123456,
		Title:            "Go Programming",
		Authors:          "John Doe",
		Publisher:        "Tech Press",
		Version:          "1st",
		LibID:            1,
		Total_copies:     5,
		Available_copies: 5,
	}
	config.DB.Create(&book)

	router.Use(func(c *gin.Context) {
		c.Set("email", "admin@example.com")
		c.Next()
	})

	isbnStr := strconv.Itoa(54321)
	router.DELETE("/books/delete/:isbn", DeleteBook)

	req, _ := http.NewRequest("DELETE", "/books/delete/"+isbnStr, nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Book not found")

}
func TestDeleteBook_InvalidISBNFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testutils.SetupTestDB()

	router := gin.Default()
	router.DELETE("/books/delete/:isbn", DeleteBook)

	req, _ := http.NewRequest("DELETE", "/books/delete/invalid_isbn_format", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid ISBN format")
}
