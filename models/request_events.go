package models

import "time"

type RequestEvents struct {
	ReqID        uint      `gorm:"primaryKey"  json:"reqID"`
	BookID       uint      `gorm:"not null" binding:"required" json:"bookID"`
	ReaderID     uint      `gorm:"not null" binding:"required" json:"readerID"`
	RequestDate  time.Time `gorm:"not null" binding:"required" json:"request_date"`
	ApprovalDate time.Time `gorm:"not null" binding:"required" json:"approval_date"`
	ApproverID   uint      `gorm:"not null" binding:"required" json:"approverID"`
	RequestType  string    `gorm:"not null" binding:"required,oneof=issue return" json:"request_type" `
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	Book     Books `gorm:"foreignKey:BookID"`
	Reader   User  `gorm:"foreignKey:ReaderID"`
	Approver User  `gorm:"foreignKey:ApproverID"`
}
