package models

import (
	"time"
)
type Hotel struct{
	ID uint `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"size:255;not null;uniqueIndex:idx_name_city"`
	City string `json:"city" gorm:"size:100;not null;index;uniqueIndex:idx_name_city"`
	Location string `json:"location" gorm:"size:255"`
	Price float64 `json:"price" gorm:"type:decimal(10,2);not null;index"`
	Rating float64 `json:"rating" gorm:"type:decimal(3,2)"`
	ImageURL string `json:"image_url" gorm:"type:text"`
	Source string `json:"source" gorm:"size:100"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type PriceHistory struct{
	ID uint `json:"id" gorm:"primaryKey"`
	HotelID uint `json:"hotel_id" gorm:"not null;index"`
	Price float64 `json:"price" gorm:"type:decimal(10,2);not null"`
	Timestamp time.Time `json:"timestamp" gorm:"index"`

}
type ScrapingLog struct {
	ID uint `json:"id" gorm:"primaryKey"`
	City string `json:"city" gorm:"size:100;not null"`
	Status string `json:"status" gorm:"size:50;not null"`
	HotelsCount int `json:"hotels_count" gorm:"default:0"`
	ErrorMessage string `json:"error_message,omitempty" gorm:"type:text"`
	StartedAt time.Time `json:"started_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}
type User struct {
	ID uint `json:"id" gorm:"primaryKey"`
	Email string `json:"email" gorm:"uniqueIndex;not null;size:255"`
	Password string `json:"-" gorm:"not null"`
	Name string `json:"name" gorm:"size:255"`
	Role string `json:"role" gorm:"size:50;default:'user'"` // user, admin
	IsVerified bool `json:"is_verified" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}