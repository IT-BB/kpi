package models

import (
	"time"

	"gorm.io/gorm"
)

type Review360Campaign struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Name        string         `gorm:"not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	StartDate   time.Time      `gorm:"not null" json:"start_date"`
	EndDate     time.Time      `gorm:"not null" json:"end_date"`
	Status      string         `gorm:"not null;default:'draft'" json:"status"`
	CreatedBy   uint           `gorm:"not null" json:"created_by"`
	Creator     *User          `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}

type Review360Assignment struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	CampaignID   uint           `gorm:"not null" json:"campaign_id"`
	Campaign     *Review360Campaign `gorm:"foreignKey:CampaignID" json:"campaign,omitempty"`
	RevieweeID   uint           `gorm:"not null" json:"reviewee_id"`
	Reviewee     *User          `gorm:"foreignKey:RevieweeID" json:"reviewee,omitempty"`
	ReviewerID   uint           `gorm:"not null" json:"reviewer_id"`
	Reviewer     *User          `gorm:"foreignKey:ReviewerID" json:"reviewer,omitempty"`
	Relationship string         `gorm:"not null" json:"relationship"`
	Status       string         `gorm:"not null;default:'pending'" json:"status"`
	CompletedAt  *time.Time     `json:"completed_at"`
}

type Review360Response struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	AssignmentID uint           `gorm:"not null" json:"assignment_id"`
	Assignment   *Review360Assignment `gorm:"foreignKey:AssignmentID" json:"assignment,omitempty"`
	Dimension    string         `gorm:"not null" json:"dimension"`
	Score        int            `gorm:"not null" json:"score"`
	Comment      string         `gorm:"type:text" json:"comment"`
}

func (Review360Campaign) TableName() string {
	return "review360_campaigns"
}

func (Review360Assignment) TableName() string {
	return "review360_assignments"
}

func (Review360Response) TableName() string {
	return "review360_responses"
}

const (
	RelationshipManager     = "manager"
	RelationshipPeer        = "peer"
	RelationshipSubordinate = "subordinate"
	RelationshipSelf        = "self"

	DimensionLeadership    = "leadership"
	DimensionCommunication = "communication"
	DimensionTeamwork      = "teamwork"
	DimensionTechnical     = "technical"
	DimensionInnovation    = "innovation"
)
