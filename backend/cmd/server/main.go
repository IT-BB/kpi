package main

import (
	"log"

	"enterprise-kpi/internal/app"
	"enterprise-kpi/internal/config"
	"enterprise-kpi/internal/database"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	if err := database.Seed(db, cfg); err != nil {
		log.Fatalf("failed to seed reference data: %v", err)
	}

	router := app.NewRouter(db, cfg)
	log.Printf("Enterprise KPI service listening on :%s", cfg.AppPort)
	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
