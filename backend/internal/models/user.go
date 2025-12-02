package models

import "time"

const (
	RoleAdmin    = "admin"
	RoleManager  = "manager"
	RoleEmployee = "employee"
)

// User represents an employee or manager in the organization.
type User struct {
	BaseModel
	FullName     string      `gorm:"size:200" json:"fullName"`
	Email        string      `gorm:"size:200;uniqueIndex" json:"email"`
	PasswordHash string      `json:"-"`
	Role         string      `gorm:"size:32" json:"role"`
	Title        string      `gorm:"size:120" json:"title"`
	DepartmentID *uint       `json:"departmentId"`
	Department   *Department `json:"department"`
	ManagerID    *uint       `json:"managerId"`
	Manager      *User       `gorm:"foreignKey:ManagerID" json:"manager"`
	JoinedAt     time.Time   `json:"joinedAt"`
}
