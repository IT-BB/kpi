package models

import (
	"time"

	"gorm.io/gorm"
)

type Department struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Name        string         `gorm:"uniqueIndex;not null" json:"name"`
	Description string         `json:"description"`
	ParentID    *uint          `json:"parent_id"`
	Parent      *Department    `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Level       int            `gorm:"not null;default:1" json:"level"`
}

func (Department) TableName() string {
	return "departments"
}
