package handler

import (
	"enterprise-kpi/internal/config"

	"gorm.io/gorm"
)

// API aggregates all HTTP handlers with shared dependencies.
type API struct {
	DB  *gorm.DB
	Cfg config.Config
}

// New wires the dependencies for HTTP handlers.
func New(db *gorm.DB, cfg config.Config) *API {
	return &API{DB: db, Cfg: cfg}
}
