package models

import "time"

type Library struct {
	ID        uint      `gorm:"primaryKey" binding:"required" json:"id"`
	Name      string    `binding:"required" gorm:"unique;not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`
}
