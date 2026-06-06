package models

import "time"

type Product struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CategoryID  uint      `gorm:"not null;index" json:"categoryId"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description"`
	Price       float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	ImageURL    string    `json:"imageUrl"`
	IsAvailable bool      `gorm:"not null;default:true" json:"isAvailable"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`

	Category Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}
