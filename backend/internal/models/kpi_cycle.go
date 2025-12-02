package models

import "time"

// KPICycle represents a KPI assessment window (monthly/quarterly/etc).
type KPICycle struct {
	BaseModel
	Name        string    `gorm:"size:150" json:"name"`
	Description string    `gorm:"size:300" json:"description"`
	StartsAt    time.Time `json:"startsAt"`
	EndsAt      time.Time `json:"endsAt"`
	Status      string    `gorm:"size:50" json:"status"`
}
