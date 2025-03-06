package models

import "time"

type User struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	Name           string    `gorm:"unique;not null" binding:"required" json:"name"`
	Email          string    `gorm:"unique" binding:"required" json:"email"`
	Contact_number string    `gorm:"not null" binding:"required" json:"contact_number"`
	Role           string    `gorm:"not null" binding:"required,oneof=admin reader" json:"role"`
	LibID         uint      `gorm:"not null" json:"lib_ID"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	IssueRequests []RequestEvents `gorm:"foreignKey:ReaderID"`

	ApprovedIssues []IssueRegistry `gorm:"foreignKey:IssueApproverID"`
}
