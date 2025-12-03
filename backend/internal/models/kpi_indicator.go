package models

import (
	"time"

	"gorm.io/gorm"
)

type KPIIndicator struct {
	ID              uint           `gorm:"primarykey" json:"id"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
	Name            string         `gorm:"not null" json:"name"`
	Description     string         `gorm:"type:text" json:"description"`
	Category        string         `gorm:"not null" json:"category"`
	DataType        string         `gorm:"not null" json:"data_type"`
	Unit            string         `json:"unit"`
	CalculationType string         `gorm:"not null" json:"calculation_type"`
	Formula         string         `gorm:"type:text" json:"formula"`
	ApplicableRoles string         `gorm:"type:text" json:"applicable_roles"`
	IsPublished     bool           `gorm:"default:false" json:"is_published"`
	CreatedBy       uint           `json:"created_by"`
	Creator         *User          `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}

func (KPIIndicator) TableName() string {
	return "kpi_indicators"
}

const (
	DataTypeNumeric    = "numeric"
	DataTypePercentage = "percentage"
	DataTypeBoolean    = "boolean"

	CalculationTypeLinear     = "linear"
	CalculationTypeStepped    = "stepped"
	CalculationTypeCustom     = "custom"

	CategoryBusiness   = "business"
	CategoryTechnical  = "technical"
	CategoryManagement = "management"
	CategoryQuality    = "quality"
)
