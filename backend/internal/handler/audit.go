package handler

import (
	"encoding/json"
	"log"

	"enterprise-kpi/internal/models"
)

func (api *API) recordAudit(actorID uint, entityType string, entityID uint, action string, metadata interface{}) {
	raw, _ := json.Marshal(metadata)
	entry := models.AuditLog{
		ActorID:    actorID,
		EntityType: entityType,
		EntityID:   entityID,
		Action:     action,
		Metadata:   string(raw),
	}
	if err := api.DB.Create(&entry).Error; err != nil {
		log.Printf("failed to persist audit log: %v", err)
	}
}
