package models

import (
	"time"

	"gorm.io/gorm"
)

type AssessmentPlan struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Name        string         `gorm:"not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	CycleType   string         `gorm:"not null" json:"cycle_type"`
	StartDate   time.Time      `gorm:"not null" json:"start_date"`
	EndDate     time.Time      `gorm:"not null" json:"end_date"`
	EmployeeID  uint           `gorm:"not null" json:"employee_id"`
	Employee    *User          `gorm:"foreignKey:EmployeeID" json:"employee,omitempty"`
	ManagerID   uint           `gorm:"not null" json:"manager_id"`
	Manager     *User          `gorm:"foreignKey:ManagerID" json:"manager,omitempty"`
	Status      string         `gorm:"not null;default:'draft'" json:"status"`
	ConfirmedAt *time.Time     `json:"confirmed_at"`
	Items       []AssessmentPlanItem `gorm:"foreignKey:PlanID" json:"items,omitempty"`
}

type AssessmentPlanItem struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	PlanID       uint           `gorm:"not null" json:"plan_id"`
	IndicatorID  uint           `gorm:"not null" json:"indicator_id"`
	Indicator    *KPIIndicator  `gorm:"foreignKey:IndicatorID" json:"indicator,omitempty"`
	Weight       float64        `gorm:"not null" json:"weight"`
	TargetValue  float64        `gorm:"not null" json:"target_value"`
	ChallengeValue float64      `json:"challenge_value"`
	BaselineValue  float64      `json:"baseline_value"`
}

func (AssessmentPlan) TableName() string {
	return "assessment_plans"
}

func (AssessmentPlanItem) TableName() string {
	return "assessment_plan_items"
}

const (
	CycleTypeMonthly    = "monthly"
	CycleTypeQuarterly  = "quarterly"
	CycleTypeYearly     = "yearly"

	StatusDraft       = "draft"
	StatusPending     = "pending"
	StatusConfirmed   = "confirmed"
	StatusInProgress  = "in_progress"
	StatusCompleted   = "completed"
	StatusCancelled   = "cancelled"
)
