package models

// AuditLog captures key workflow events for traceability.
type AuditLog struct {
	BaseModel
	ActorID    uint   `json:"actorId"`
	Actor      User   `json:"actor"`
	EntityType string `gorm:"size:100" json:"entityType"`
	EntityID   uint   `json:"entityId"`
	Action     string `gorm:"size:100" json:"action"`
	Metadata   string `gorm:"type:text" json:"metadata"`
}
