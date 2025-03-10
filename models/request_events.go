package models

import "time"

type RequestEvents struct {
	ReqID        uint      `gorm:"primaryKey"  json:"reqID"`
	BookID       uint      `gorm:"not null" binding:"required" json:"bookID"`
	ReaderID     uint      `gorm:"not null" binding:"required" json:"readerID"`
	RequestDate  time.Time `gorm:"not null"`
	ApprovalDate *time.Time
	ApproverID   *uint
	RequestType  string     `gorm:"default:'issue';check:request_type IN ('issue','return')"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	Book     Books `gorm:"foreignKey:BookID"`
	Reader   User  `gorm:"foreignKey:ReaderID"`
	Approver User  `gorm:"foreignKey:ApproverID"`
}
