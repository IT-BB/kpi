package models

import (
	"time"

	"gorm.io/gorm"
)

type AuditLog struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	UserID      uint           `gorm:"not null" json:"user_id"`
	User        *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Action      string         `gorm:"not null" json:"action"`
	Resource    string         `gorm:"not null" json:"resource"`
	ResourceID  uint           `json:"resource_id"`
	IPAddress   string         `json:"ip_address"`
	UserAgent   string         `gorm:"type:text" json:"user_agent"`
	Details     string         `gorm:"type:text" json:"details"`
	Status      string         `gorm:"not null" json:"status"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}

const (
	ActionCreate = "create"
	ActionRead   = "read"
	ActionUpdate = "update"
	ActionDelete = "delete"
	ActionLogin  = "login"
	ActionLogout = "logout"

	AuditStatusSuccess = "success"
	AuditStatusFailure = "failure"
)
