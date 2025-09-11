package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	NameEn        string         `json:"name_en" gorm:"not null" validate:"required,min=1,max=255"`
	NameAr        string         `json:"name_ar" gorm:"type:varchar(255)"`
	DescriptionEn string         `json:"description_en" gorm:"type:text"`
	DescriptionAr string         `json:"description_ar" gorm:"type:text"`
	SKU           string         `json:"sku" gorm:"type:varchar(100)"`
	CategoryID    uint           `json:"category_id" gorm:"default:1" validate:"required"`
	ImageURL      string         `json:"image_url" gorm:"type:text"`
	IsActive      bool           `json:"is_active" gorm:"default:true"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`

	// Virtual fields for backward compatibility
	Name        string    `json:"name" gorm:"-"`
	Description string    `json:"description" gorm:"-"`
	Category    string    `json:"category" gorm:"-"`

	// Relationship
	CategoryInfo *Category `json:"category_info,omitempty" gorm:"foreignKey:CategoryID"`
}

type ProductCreateRequest struct {
	NameEn        string `json:"name_en" validate:"required,min=1,max=255"`
	NameAr        string `json:"name_ar" validate:"omitempty,max=255"`
	DescriptionEn string `json:"description_en"`
	DescriptionAr string `json:"description_ar"`
	SKU           string `json:"sku" validate:"omitempty,max=100"`
	CategoryID    uint   `json:"category_id" validate:"required"`
	ImageURL      string `json:"image_url"`
	IsActive      *bool  `json:"is_active"`
}

type ProductUpdateRequest struct {
	NameEn        *string `json:"name_en" validate:"omitempty,min=1,max=255"`
	NameAr        *string `json:"name_ar" validate:"omitempty,max=255"`
	DescriptionEn *string `json:"description_en"`
	DescriptionAr *string `json:"description_ar"`
	SKU           *string `json:"sku" validate:"omitempty,max=100"`
	CategoryID    *uint   `json:"category_id" validate:"omitempty"`
	ImageURL      *string `json:"image_url"`
	IsActive      *bool   `json:"is_active"`
}
