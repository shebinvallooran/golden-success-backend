package models

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	NameEn        string         `json:"name_en" gorm:"not null" validate:"required,min=1,max=100"`
	NameAr        string         `json:"name_ar" gorm:"type:varchar(100)"`
	DescriptionEn string         `json:"description_en" gorm:"type:text"`
	DescriptionAr string         `json:"description_ar" gorm:"type:text"`

	// Home screen description fields
	HomeDescriptionEn string `json:"home_description_en" gorm:"type:text"`
	HomeDescriptionAr string `json:"home_description_ar" gorm:"type:text"`

	// Three key points fields
	Point1En string `json:"point1_en" gorm:"type:varchar(255)"`
	Point1Ar string `json:"point1_ar" gorm:"type:varchar(255)"`
	Point2En string `json:"point2_en" gorm:"type:varchar(255)"`
	Point2Ar string `json:"point2_ar" gorm:"type:varchar(255)"`
	Point3En string `json:"point3_en" gorm:"type:varchar(255)"`
	Point3Ar string `json:"point3_ar" gorm:"type:varchar(255)"`

	// Image fields
	ImageURL string `json:"image_url" gorm:"type:varchar(500)"`

	Priority      int            `json:"priority" gorm:"not null;default:0" validate:"min=0"`
	IsActive      bool           `json:"is_active" gorm:"default:true"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`

	// Virtual fields for backward compatibility
	Name        string `json:"name" gorm:"-"`
	Description string `json:"description" gorm:"-"`

	// Relationship with products
	Products []Product `json:"products,omitempty" gorm:"foreignKey:CategoryID"`
}

type CategoryCreateRequest struct {
	NameEn        string `json:"name_en" validate:"required,min=1,max=100"`
	NameAr        string `json:"name_ar" validate:"omitempty,max=100"`
	DescriptionEn string `json:"description_en"`
	DescriptionAr string `json:"description_ar"`

	// Home screen description fields
	HomeDescriptionEn string `json:"home_description_en"`
	HomeDescriptionAr string `json:"home_description_ar"`

	// Three key points fields
	Point1En string `json:"point1_en" validate:"omitempty,max=255"`
	Point1Ar string `json:"point1_ar" validate:"omitempty,max=255"`
	Point2En string `json:"point2_en" validate:"omitempty,max=255"`
	Point2Ar string `json:"point2_ar" validate:"omitempty,max=255"`
	Point3En string `json:"point3_en" validate:"omitempty,max=255"`
	Point3Ar string `json:"point3_ar" validate:"omitempty,max=255"`

	// Image field
	ImageURL string `json:"image_url" validate:"omitempty,max=500"`

	Priority      *int   `json:"priority" validate:"omitempty,min=0"`
	IsActive      *bool  `json:"is_active"`
}

type CategoryUpdateRequest struct {
	NameEn        *string `json:"name_en" validate:"omitempty,min=1,max=100"`
	NameAr        *string `json:"name_ar" validate:"omitempty,max=100"`
	DescriptionEn *string `json:"description_en"`
	DescriptionAr *string `json:"description_ar"`

	// Home screen description fields
	HomeDescriptionEn *string `json:"home_description_en"`
	HomeDescriptionAr *string `json:"home_description_ar"`

	// Three key points fields
	Point1En *string `json:"point1_en" validate:"omitempty,max=255"`
	Point1Ar *string `json:"point1_ar" validate:"omitempty,max=255"`
	Point2En *string `json:"point2_en" validate:"omitempty,max=255"`
	Point2Ar *string `json:"point2_ar" validate:"omitempty,max=255"`
	Point3En *string `json:"point3_en" validate:"omitempty,max=255"`
	Point3Ar *string `json:"point3_ar" validate:"omitempty,max=255"`

	// Image field
	ImageURL *string `json:"image_url" validate:"omitempty,max=500"`

	Priority      *int    `json:"priority" validate:"omitempty,min=0"`
	IsActive      *bool   `json:"is_active"`
}

type CategoryWithProductsResponse struct {
	ID            uint      `json:"id"`
	NameEn        string    `json:"name_en"`
	NameAr        string    `json:"name_ar"`
	DescriptionEn string    `json:"description_en"`
	DescriptionAr string    `json:"description_ar"`

	// Home screen description fields
	HomeDescriptionEn string `json:"home_description_en"`
	HomeDescriptionAr string `json:"home_description_ar"`

	// Three key points fields
	Point1En string `json:"point1_en"`
	Point1Ar string `json:"point1_ar"`
	Point2En string `json:"point2_en"`
	Point2Ar string `json:"point2_ar"`
	Point3En string `json:"point3_en"`
	Point3Ar string `json:"point3_ar"`

	// Image field
	ImageURL string `json:"image_url"`

	Priority      int       `json:"priority"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	ProductCount  int       `json:"product_count"`
	Products      []Product `json:"products,omitempty"`

	// Virtual fields for backward compatibility
	Name        string `json:"name"`
	Description string `json:"description"`
}
