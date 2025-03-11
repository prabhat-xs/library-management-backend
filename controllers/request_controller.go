package controllers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prabhat-xs/library-management-backend/config"
	"github.com/prabhat-xs/library-management-backend/models"
	"gorm.io/gorm"
)

func RaiseIssueRequest(c *gin.Context) {
	var input struct {
		ISBN        uint   `binding:"required"`
		RequestType string `binding:"required"` 
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !(input.RequestType == "issue" || input.RequestType == "return") {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request type!",
		})
		return
	}

	var book models.Books
	if err := config.DB.Where("isbn = ?", input.ISBN).First(&book).Error; err != nil || book.Available_copies <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Book unavailable"})
		return
	}

	floatId, _ := c.Get("id")
	id := floatId.(uint)

	// TODO This needs to be validated using issue registry table
	// var issueReq models.RequestEvents
	// if input.RequestType == "return"{
	// 	if err:= config.DB.Where("bookID = ?",input.ISBN).Where("readerID = ?",id).Where("request_type = ?","issue").First(&issueReq); err != nil {
	// 		c.JSON(http.StatusBadRequest,gin.H{
	// 			"error": "No issue request exists corresponding to this return request",
	// 		})
	// 		return
	// 	}
	// }

	// TODO temporary check, better implementation using registry table to be done
	var issueReq models.RequestEvents
	if err := config.DB.Where("reader_Id= ? AND book_ID = ? ", id, input.ISBN).Take(&issueReq).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Duplicate request!",
		})
		return

	}

	reqEvent := models.RequestEvents{
		BookID:      input.ISBN,
		ReaderID:    id,
		RequestType: input.RequestType,
		RequestDate: time.Now(),
	}
	config.DB.Create(&reqEvent)
	c.JSON(http.StatusOK, gin.H{"message": "Issue request raised successfully"})
}

func ListRequests(c *gin.Context) {
	var requests []models.RequestEvents

	if err := config.DB.Find(&requests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"requests": requests,
	})
}

func ProcessIssueRequest(c *gin.Context) {
	var input struct {
		Action string `binding:"required" json:"action"`
		ReqID  uint   `binding:"required" json:"reqid"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Action != "approve" && input.Action != "reject" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Allowed actions are approve or reject",
		})
		return
	}

	if input.Action == "reject" {
		c.JSON(http.StatusOK, gin.H{
			"message": "Request processed succesfully! Issue req rejected!",
		})
		return
	}

	// GETTING ADMIN ID from JWT
	id, _ := c.Get("id")
	ApproverID := id.(uint)

	txErr := config.DB.Transaction(func(tx *gorm.DB) error {
		var req models.RequestEvents
		if tx.Where("req_id = ?", input.ReqID).First(&req).Error != nil {
			return errors.New("request not found")
		}

		var book models.Books
		if tx.Where("isbn = ?", req.BookID).First(&book).Error != nil || book.Available_copies <= 0 {
			return errors.New("no copies available")
		}

		now := time.Now()
		req.ApprovalDate = &now
		req.ApproverID = &ApproverID
		tx.Save(&req)

		issueReg := models.IssueRegistry{
			ISBN:               book.ISBN,
			ReaderID:           req.ReaderID,
			IssueApproverID:    ApproverID,
			IssueStatus:        input.Action,
			IssueDate:          now,
			ExpectedReturnDate: now.AddDate(0, 0, 14),
		}
		tx.Create(&issueReg)

		book.Available_copies -= 1
		tx.Save(&book)

		return nil
	})

	if txErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": txErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Issue request approved successfully"})
}
