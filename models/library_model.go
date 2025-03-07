package models

import "time"

type Library struct {
	ID        uint      `gorm:"primaryKey"  json:"id"`
	Name      string    `binding:"required" gorm:"unique;not null" json:"lib_name"`
	CreatedAt time.Time `json:"created_at"`

	Users []User `gorm:"foreignKey:LibID"`
	Books []Books `gorm:"foreignKey:LibID"`
}
