package models

import "time"

const (
	AssignmentStatusDraft     = "draft"
	AssignmentStatusSubmitted = "submitted"
	AssignmentStatusReviewed  = "reviewed"
	AssignmentStatusApproved  = "approved"
)

// KPIAssignment links a user to a template within a KPI cycle.
type KPIAssignment struct {
	BaseModel
	CycleID     uint         `json:"cycleId"`
	Cycle       KPICycle     `json:"cycle"`
	TemplateID  uint         `json:"templateId"`
	Template    KPITemplate  `json:"template"`
	AssigneeID  uint         `json:"assigneeId"`
	Assignee    User         `json:"assignee"`
	ManagerID   uint         `json:"managerId"`
	Manager     User         `json:"manager"`
	Status      string       `gorm:"size:50" json:"status"`
	DueDate     time.Time    `json:"dueDate"`
	SubmittedAt *time.Time   `json:"submittedAt"`
	ReviewedAt  *time.Time   `json:"reviewedAt"`
	ApprovedAt  *time.Time   `json:"approvedAt"`
	Scores      []KPIScore   `json:"scores"`
	Comments    []KPIComment `json:"comments"`
}

// KPIScore captures both self and manager scores for a metric.
type KPIScore struct {
	BaseModel
	AssignmentID uint      `json:"assignmentId"`
	MetricID     uint      `json:"metricId"`
	Metric       KPIMetric `json:"metric"`
	SelfScore    float64   `json:"selfScore"`
	ManagerScore float64   `json:"managerScore"`
	Comment      string    `gorm:"size:500" json:"comment"`
}

// KPIComment records workflow comments (submission, review, approval).
type KPIComment struct {
	BaseModel
	AssignmentID uint   `json:"assignmentId"`
	AuthorID     uint   `json:"authorId"`
	Author       User   `json:"author"`
	Message      string `gorm:"size:500" json:"message"`
	Type         string `gorm:"size:50" json:"type"`
}
