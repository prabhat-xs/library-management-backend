package models

type Books struct {
	ISBN             uint   `gorm:"primaryKey" binding:"required" json:"isbn"`
	Title            string `gorm:"not null" binding:"required" json:"title"`
	Authors          string `gorm:"not null" binding:"required" json:"authors"`
	Publisher        string `gorm:"not null" binding:"required" json:"publisher"`
	Version          string `gorm:"not null" binding:"required" json:"version"`
	Total_copies     uint   `gorm:"not null" binding:"required,min=1" json:"total_copies"`
	Available_copies uint   `gorm:"not null" bindiing:"required" json:"available_copies"`
	LibID            uint   `gorm:"not null" binding:"required" json:"lib_ID"`
}
