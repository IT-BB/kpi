package models

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Name        string         `gorm:"uniqueIndex;not null" json:"name"`
	Description string         `json:"description"`
	Permissions string         `gorm:"type:text" json:"permissions"`
}

func (Role) TableName() string {
	return "roles"
}

const (
	RoleAdmin    = "admin"
	RoleHR       = "hr"
	RoleManager  = "manager"
	RoleEmployee = "employee"
)
