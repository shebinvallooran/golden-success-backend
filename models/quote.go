package models

import (
	"time"

	"gorm.io/gorm"
)

type Quote struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"type:varchar(255);not null" validate:"required"`
	Email       string         `json:"email" gorm:"type:varchar(255);not null" validate:"required,email"`
	Phone       string         `json:"phone" gorm:"type:varchar(50);not null" validate:"required"`
	Company     string         `json:"company" gorm:"type:varchar(255)"`
	ProductName string         `json:"product_name" gorm:"type:varchar(255);not null" validate:"required"`
	Quantity    int            `json:"quantity" gorm:"default:1" validate:"required,min=1"`
	Message     string         `json:"message" gorm:"type:text"`
	Status      string         `json:"status" gorm:"type:varchar(50);default:'pending'"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type QuoteStatusUpdateRequest struct {
	Status string `json:"status" validate:"required,oneof=pending contacted resolved"`
}
