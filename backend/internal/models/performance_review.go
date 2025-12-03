package models

import (
	"time"

	"gorm.io/gorm"
)

type PerformanceReview struct {
	ID               uint           `gorm:"primarykey" json:"id"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
	PlanID           uint           `gorm:"not null" json:"plan_id"`
	Plan             *AssessmentPlan `gorm:"foreignKey:PlanID" json:"plan,omitempty"`
	EmployeeID       uint           `gorm:"not null" json:"employee_id"`
	Employee         *User          `gorm:"foreignKey:EmployeeID" json:"employee,omitempty"`
	ReviewerID       uint           `gorm:"not null" json:"reviewer_id"`
	Reviewer         *User          `gorm:"foreignKey:ReviewerID" json:"reviewer,omitempty"`
	SelfScore        *float64       `json:"self_score"`
	SelfComment      string         `gorm:"type:text" json:"self_comment"`
	ManagerScore     *float64       `json:"manager_score"`
	ManagerComment   string         `gorm:"type:text" json:"manager_comment"`
	FinalScore       *float64       `json:"final_score"`
	PerformanceLevel string         `json:"performance_level"`
	Status           string         `gorm:"not null;default:'pending'" json:"status"`
	SubmittedAt      *time.Time     `json:"submitted_at"`
	Items            []PerformanceReviewItem `gorm:"foreignKey:ReviewID" json:"items,omitempty"`
}

type PerformanceReviewItem struct {
	ID              uint           `gorm:"primarykey" json:"id"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
	ReviewID        uint           `gorm:"not null" json:"review_id"`
	PlanItemID      uint           `gorm:"not null" json:"plan_item_id"`
	PlanItem        *AssessmentPlanItem `gorm:"foreignKey:PlanItemID" json:"plan_item,omitempty"`
	ActualValue     float64        `json:"actual_value"`
	CompletionRate  float64        `json:"completion_rate"`
	AutoScore       float64        `json:"auto_score"`
	AdjustedScore   *float64       `json:"adjusted_score"`
	AdjustmentReason string        `gorm:"type:text" json:"adjustment_reason"`
}

func (PerformanceReview) TableName() string {
	return "performance_reviews"
}

func (PerformanceReviewItem) TableName() string {
	return "performance_review_items"
}

const (
	PerformanceLevelOutstanding = "outstanding"
	PerformanceLevelExcellent   = "excellent"
	PerformanceLevelGood        = "good"
	PerformanceLevelFair        = "fair"
	PerformanceLevelPoor        = "poor"
)
