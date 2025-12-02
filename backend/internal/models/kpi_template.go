package models

// KPITemplate describes the structure and weights for a KPI assessment.
type KPITemplate struct {
	BaseModel
	Name        string      `gorm:"size:150" json:"name"`
	Description string      `gorm:"size:500" json:"description"`
	OwnerID     uint        `json:"ownerId"`
	Owner       User        `json:"owner"`
	Metrics     []KPIMetric `json:"metrics"`
	IsActive    bool        `json:"isActive"`
}

// KPIMetric represents a measurable dimension of a KPI template.
type KPIMetric struct {
	BaseModel
	TemplateID  uint    `json:"templateId"`
	Name        string  `gorm:"size:150" json:"name"`
	Description string  `gorm:"size:300" json:"description"`
	Weight      int     `json:"weight"`
	TargetValue float64 `json:"targetValue"`
	Unit        string  `gorm:"size:50" json:"unit"`
}
