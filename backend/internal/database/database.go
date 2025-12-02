package database

import (
	"log"
	"time"

	"enterprise-kpi/internal/config"
	"enterprise-kpi/internal/models"
	"enterprise-kpi/pkg/password"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect establishes a MySQL backed gorm instance.
func Connect(cfg config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

// AutoMigrate synchronizes the schema.
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Department{},
		&models.User{},
		&models.KPITemplate{},
		&models.KPIMetric{},
		&models.KPICycle{},
		&models.KPIAssignment{},
		&models.KPIScore{},
		&models.KPIComment{},
		&models.AuditLog{},
	)
}

// Seed populates reference data for a fresh environment.
func Seed(db *gorm.DB, cfg config.Config) error {
	var admin models.User
	if err := db.Where("email = ?", cfg.DefaultAdmin).First(&admin).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			hash, _ := password.Hash(cfg.DefaultPass)
			admin = models.User{
				FullName:     "Platform Administrator",
				Email:        cfg.DefaultAdmin,
				PasswordHash: hash,
				Role:         models.RoleAdmin,
				Title:        "Head of Performance",
				JoinedAt:     time.Now(),
			}
			if err := db.Create(&admin).Error; err != nil {
				return err
			}
			log.Printf("Seeded default admin %s", cfg.DefaultAdmin)
		} else {
			return err
		}
	}

	// Departments
	deptNames := []string{"Product", "Sales", "People Ops"}
	for _, name := range deptNames {
		var dept models.Department
		if err := db.Where("name = ?", name).First(&dept).Error; err == gorm.ErrRecordNotFound {
			dept = models.Department{Name: name, Description: name + " department"}
			if err := db.Create(&dept).Error; err != nil {
				return err
			}
		}
	}

	// Default KPI Template
	var template models.KPITemplate
	if err := db.Where("name = ?", "Quarterly Leadership").Preload("Metrics").First(&template).Error; err == gorm.ErrRecordNotFound {
		template = models.KPITemplate{
			Name:        "Quarterly Leadership",
			Description: "Standard leadership KPIs",
			OwnerID:     admin.ID,
			IsActive:    true,
			Metrics: []models.KPIMetric{
				{Name: "Revenue Impact", Description: "Contribution to revenue targets", Weight: 40, TargetValue: 100, Unit: "%"},
				{Name: "Team Health", Description: "Retention and engagement", Weight: 30, TargetValue: 90, Unit: "%"},
				{Name: "Operational Excellence", Description: "Process and delivery", Weight: 30, TargetValue: 95, Unit: "%"},
			},
		}
		if err := db.Create(&template).Error; err != nil {
			return err
		}
	}

	// Default KPI Cycle
	var cycle models.KPICycle
	if err := db.Where("name = ?", "FY24-Q4").First(&cycle).Error; err == gorm.ErrRecordNotFound {
		cycle = models.KPICycle{
			Name:        "FY24-Q4",
			Description: "Final quarter of FY24",
			StartsAt:    time.Now().AddDate(0, -1, 0),
			EndsAt:      time.Now().AddDate(0, 2, 0),
			Status:      "in-progress",
		}
		if err := db.Create(&cycle).Error; err != nil {
			return err
		}
	}

	return nil
}
