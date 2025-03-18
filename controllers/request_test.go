package controllers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prabhat-xs/library-management-backend/config"
	"github.com/prabhat-xs/library-management-backend/models"
	"github.com/prabhat-xs/library-management-backend/testutils"

	"github.com/stretchr/testify/assert"
)

type CustomResponseRecorder struct {
	*httptest.ResponseRecorder
}

func (c *CustomResponseRecorder) CloseNotify() <-chan bool {
	return make(chan bool)
}

func TestRaiseBookRequest_Issue(t *testing.T) {
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
		c.Set("id", uint(2))
		c.Set("libid", uint(1))
		c.Next()
	})

	router.POST("/requests/raise", RaiseBookRequest)

	issuePayload := `{
		"isbn": 123456,
		"requestType": "issue"
	}`

	req, _ := http.NewRequest("POST", "/requests/raise", bytes.NewBuffer([]byte(issuePayload)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Issue request raised successfully")
}
func TestRaiseBookRequest_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testutils.SetupTestDB()

	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set("libid", uint(1))
		c.Next()
	})

	router.POST("/requests/raise", RaiseBookRequest)

	issuePayload := `{invalid	   	}`

	req, _ := http.NewRequest("POST", "/requests/raise", bytes.NewBuffer([]byte(issuePayload)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}

func TestRaiseBookRequest_InvalidReqType(t *testing.T) {
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
		c.Set("id", uint(2))
		c.Set("libid", uint(1))
		c.Next()
	})

	router.POST("/requests/raise", RaiseBookRequest)

	issuePayload := `{
		"isbn": 123456,
		"requestType": "tissue"
	}`

	req, _ := http.NewRequest("POST", "/requests/raise", bytes.NewBuffer([]byte(issuePayload)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Request types can be issue OR return only!")
}
func TestRaiseBookRequest_BookNotAvailable(t *testing.T) {
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
		c.Set("id", uint(2))
		c.Set("libid", uint(1))
		c.Next()
	})

	router.POST("/requests/raise", RaiseBookRequest)

	issuePayload := `{
		"isbn": 54321,
		"requestType": "tissue"
	}`

	req, _ := http.NewRequest("POST", "/requests/raise", bytes.NewBuffer([]byte(issuePayload)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Request types can be issue OR return only!")
}

func TestRaiseBookRequest_DuplicateIssueRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testutils.SetupTestDB()

	
	request := models.RequestEvents{
		BookID:      123456,
		ReaderID:    2,
		LibID:       1,
		RequestType: "issue",
		RequestDate: time.Now(),
	}
	config.DB.Create(&request)

	router := gin.Default()
	router.POST("/requests/raise", RaiseBookRequest)

	payload := `{"isbn": 123456, "requestType": "issue"}`
	req, _ := http.NewRequest("POST", "/requests/raise", bytes.NewBuffer([]byte(payload)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Book unavailable")
}

func TestRaiseBookRequest_NoIssueForReturn(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testutils.SetupTestDB()

	router := gin.Default()
	router.POST("/requests/raise", RaiseBookRequest)

	payload := `{"isbn": 123456, "requestType": "return"}`
	req, _ := http.NewRequest("POST", "/requests/raise", bytes.NewBuffer([]byte(payload)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Book unavailable")
}

func TestListRequests_DBError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testutils.SetupTestDB()

	
	sqlDB, _ := config.DB.DB()
	sqlDB.Close()
	router := gin.Default()
	router.GET("/requests/list", ListRequests)

	req, _ := http.NewRequest("GET", "/requests/list", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}

func TestRaiseBookRequest_Return(t *testing.T) {
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
		Available_copies: 4,
	}
	config.DB.Create(&book)

	issueRegistry := models.IssueRegistry{
		ISBN:      book.ISBN,
		LibID:     book.LibID,
		ReaderID:  2,
		Status:    "issued",
		IssueDate: time.Now(),
	}
	config.DB.Create(&issueRegistry)

	router.Use(func(c *gin.Context) {
		c.Set("id", uint(2))
		c.Set("libid", uint(1))
		c.Next()
	})

	router.POST("/requests/raise", RaiseBookRequest)

	returnPayload := `{
		"isbn": 123456,
		"requestType": "return"
	}`

	req, _ := http.NewRequest("POST", "/requests/raise", bytes.NewBuffer([]byte(returnPayload)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Return request raised successfully")
}

func TestListRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testutils.SetupTestDB()

	router := gin.Default()
	router.GET("/requests/list", ListRequests)

	request1 := models.RequestEvents{BookID: 123456, ReaderID: 2, LibID: 1, RequestType: "issue", RequestDate: time.Now()}
	request2 := models.RequestEvents{BookID: 123457, ReaderID: 3, LibID: 1, RequestType: "return", RequestDate: time.Now()}
	config.DB.Create(&request1)
	config.DB.Create(&request2)

	req, _ := http.NewRequest("GET", "/requests/list", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "requests")
}

func TestProcessRequest_ApproveIssue(t *testing.T) {
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

	request := models.RequestEvents{
		BookID:      book.ISBN,
		ReaderID:    2,
		LibID:       book.LibID,
		RequestType: "issue",
		RequestDate: time.Now(),
	}
	config.DB.Create(&request)

	router.Use(func(c *gin.Context) {
		c.Set("id", uint(1))
		c.Next()
	})

	router.POST("/requests/process", ProcessRequest)

	approvePayload := `{
		"action": "approve",
		"reqtype": "issue",
		"reqid": 1
	}`

	req, _ := http.NewRequest("POST", "/requests/process", bytes.NewBuffer([]byte(approvePayload)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Issue request approved successfully")
}
func TestProcessRequest_BADJson(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testutils.SetupTestDB()

	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Set("id", uint(1))
		c.Next()
	})

	router.POST("/requests/process", ProcessRequest)

	approvePayload := `{
		bad_json
	}`

	req, _ := http.NewRequest("POST", "/requests/process", bytes.NewBuffer([]byte(approvePayload)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestProcessRequest_Reject(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testutils.SetupTestDB()

	router := gin.Default()

	request := models.RequestEvents{
		BookID:      123456,
		ReaderID:    2,
		LibID:       1,
		RequestType: "issue",
		RequestDate: time.Now(),
	}
	config.DB.Create(&request)

	router.Use(func(c *gin.Context) {
		c.Set("id", uint(1))
		c.Next()
	})

	router.POST("/requests/process", ProcessRequest)

	rejectPayload := `{
		"action": "reject",
		"reqtype": "issue",
		"reqid": 1
	}`

	req, _ := http.NewRequest("POST", "/requests/process", bytes.NewBuffer([]byte(rejectPayload)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Request processed succesfully! Issue req rejected!")
}
func TestProcessRequest_InvalidAction(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testutils.SetupTestDB()

	router := gin.Default()

	request := models.RequestEvents{
		BookID:      123456,
		ReaderID:    2,
		LibID:       1,
		RequestType: "issue",
		RequestDate: time.Now(),
	}
	config.DB.Create(&request)

	router.Use(func(c *gin.Context) {
		c.Set("id", uint(1))
		c.Next()
	})

	router.POST("/requests/process", ProcessRequest)

	rejectPayload := `{
		"action": "breakfast",
		"reqtype": "issue",
		"reqid": 1
	}`

	req, _ := http.NewRequest("POST", "/requests/process", bytes.NewBuffer([]byte(rejectPayload)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandleReturnRequest_Success(t *testing.T) {
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
		Available_copies: 4,
	}
	config.DB.Create(&book)

	issueRegistry := models.IssueRegistry{
		ISBN:      book.ISBN,
		LibID:     book.LibID,
		ReaderID:  2,
		Status:    "issued",
		IssueDate: time.Now(),
	}
	config.DB.Create(&issueRegistry)

	returnRequest := models.RequestEvents{
		BookID:      book.ISBN,
		ReaderID:    2,
		LibID:       book.LibID,
		RequestType: "return",
		RequestDate: time.Now(),
	}
	config.DB.Create(&returnRequest)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set("id", uint(1))
	c.Params = []gin.Param{{Key: "reqId", Value: "1"}}

	handleReturnRequest(c, 1, returnRequest.ReqID)

	var updatedIssueRegistry models.IssueRegistry
	err := config.DB.Where("isbn = ? AND reader_id = ?", book.ISBN, 2).First(&updatedIssueRegistry).Error
	assert.Nil(t, err)
	assert.Equal(t, "returned", updatedIssueRegistry.Status)

	var updatedBook models.Books
	err = config.DB.Where("isbn = ?", book.ISBN).First(&updatedBook).Error
	assert.Nil(t, err)
	assert.Equal(t, uint(5), updatedBook.Available_copies)

	var deletedRequest models.RequestEvents
	err = config.DB.Where("req_id = ?", returnRequest.ReqID).First(&deletedRequest).Error
	assert.NotNil(t, err)
}

func TestHandleReturnRequest_IssueRegistryNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testutils.SetupTestDB()

	returnRequest := models.RequestEvents{
		BookID:      123456,
		ReaderID:    2,
		LibID:       1,
		RequestType: "return",
		RequestDate: time.Now(),
	}
	config.DB.Create(&returnRequest)

	w := &CustomResponseRecorder{ResponseRecorder: httptest.NewRecorder()}
	c, _ := gin.CreateTestContext(w)
	c.Set("id", uint(1))

	handleReturnRequest(c, 1, returnRequest.ReqID)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}
func TestHandleReturnRequest_RequestNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testutils.SetupTestDB()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handleReturnRequest(c, 1, 999)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}

func TestHandleReturnRequest_UpdateIssueRegistryError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testutils.SetupTestDB()

	request := models.RequestEvents{
		BookID:      123456,
		ReaderID:    2,
		LibID:       1,
		RequestType: "return",
		RequestDate: time.Now(),
	}
	config.DB.Create(&request)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handleReturnRequest(c, 1, request.ReqID) 

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}

func TestHandleReturnRequest_DeleteError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testutils.SetupTestDB()

	request := models.RequestEvents{
		BookID:      123456,
		ReaderID:    2,
		LibID:       1,
		RequestType: "return",
		RequestDate: time.Now(),
	}
	config.DB.Create(&request)

	retRegistry := models.IssueRegistry{
		ISBN:       123456,
		LibID:      1,
		ReaderID:   2,
		Status:     "issued",
		ReturnDate: nil,
	}
	config.DB.Create(&retRegistry)

	sqlDB, _ := config.DB.DB()
	sqlDB.Close()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handleReturnRequest(c, 1, request.ReqID)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}
