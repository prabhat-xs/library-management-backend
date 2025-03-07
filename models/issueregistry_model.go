package models

import "time"

type IssueRegistry struct {
	IssueID            uint      `gorm:"primaryKey" json:"issueID"`
	ISBN               uint      `gorm:"not null" binding:"required" json:"isbn"`
	ReaderID           uint      `gorm:"not null" binding:"required" json:"readerID"`
	IssueApproverID    uint      `gorm:"not null" binding:"required" json:"issueapproverID"`
	IssueStatus        string    `gorm:"not null" binding:"required" json:"status"`
	IssueDate          time.Time `gorm:"not null" binding:"required" json:"date"`
	ExpectedReturnDate time.Time `gorm:"not null" binding:"required" json:"expected_return_date"`
	ReturnDate         time.Time `gorm:"not null" binding:"required" json:"return_date"`
	ReturnApproverID   uint      `gorm:"not null" binding:"required" json:"returnapproverID"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`

	Reader User `gorm:"foreignKey:ReaderID"`
	Approver User `gorm:"foreignKey:IssueApproverID"`
}
