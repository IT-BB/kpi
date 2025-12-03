package models

import (
	"time"

	"gorm.io/gorm"
)

type KPIProgress struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	PlanItemID    uint           `gorm:"not null" json:"plan_item_id"`
	PlanItem      *AssessmentPlanItem `gorm:"foreignKey:PlanItemID" json:"plan_item,omitempty"`
	CurrentValue  float64        `gorm:"not null" json:"current_value"`
	CompletionRate float64       `json:"completion_rate"`
	UpdatedBy     uint           `gorm:"not null" json:"updated_by"`
	Updater       *User          `gorm:"foreignKey:UpdatedBy" json:"updater,omitempty"`
	Notes         string         `gorm:"type:text" json:"notes"`
}

func (KPIProgress) TableName() string {
	return "kpi_progress"
}
