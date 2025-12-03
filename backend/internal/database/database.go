package database

import (
	"fmt"
	"log"

	"github.com/kpi-system/backend/internal/config"
	"github.com/kpi-system/backend/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Initialize(cfg *config.Config) error {
	var err error
	var dialector gorm.Dialector

	switch cfg.Database.Driver {
	case "mysql":
		dsn := cfg.Database.DSN
		if dsn == "" {
			dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
				cfg.Database.Username,
				cfg.Database.Password,
				cfg.Database.Host,
				cfg.Database.Port,
				cfg.Database.Database,
			)
		}
		dialector = mysql.Open(dsn)
	case "sqlite":
		dialector = sqlite.Open(cfg.Database.Database)
	default:
		return fmt.Errorf("unsupported database driver: %s", cfg.Database.Driver)
	}

	DB, err = gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connected successfully")
	return nil
}

func AutoMigrate() error {
	err := DB.AutoMigrate(
		&models.Role{},
		&models.Department{},
		&models.User{},
		&models.KPIIndicator{},
		&models.AssessmentPlan{},
		&models.AssessmentPlanItem{},
		&models.KPIProgress{},
		&models.PerformanceReview{},
		&models.PerformanceReviewItem{},
		&models.Review360Campaign{},
		&models.Review360Assignment{},
		&models.Review360Response{},
		&models.AuditLog{},
	)
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database migration completed")
	return nil
}

func SeedInitialData() error {
	var count int64
	DB.Model(&models.Role{}).Count(&count)
	if count > 0 {
		return nil
	}

	roles := []models.Role{
		{Name: models.RoleAdmin, Description: "System Administrator", Permissions: "all"},
		{Name: models.RoleHR, Description: "HR Manager", Permissions: "hr,reports"},
		{Name: models.RoleManager, Description: "Department Manager", Permissions: "team,kpi"},
		{Name: models.RoleEmployee, Description: "Employee", Permissions: "self"},
	}

	for _, role := range roles {
		if err := DB.Create(&role).Error; err != nil {
			return fmt.Errorf("failed to seed role: %w", err)
		}
	}

	log.Println("Initial data seeded successfully")
	return nil
}
