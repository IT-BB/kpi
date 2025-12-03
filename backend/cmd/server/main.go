package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/kpi-system/backend/internal/config"
	"github.com/kpi-system/backend/internal/database"
	"github.com/kpi-system/backend/internal/handlers"
)

func main() {
	cfg := config.Load()

	if err := database.Initialize(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	if err := database.AutoMigrate(); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	if err := database.SeedInitialData(); err != nil {
		log.Fatalf("Failed to seed initial data: %v", err)
	}

	r := gin.Default()

	setupRoutes(r)

	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func setupRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")

	deptHandler := handlers.NewDepartmentHandler()
	api.POST("/departments", deptHandler.Create)
	api.PUT("/departments/:id", deptHandler.Update)
	api.DELETE("/departments/:id", deptHandler.Delete)
	api.GET("/departments/:id", deptHandler.Get)
	api.GET("/departments", deptHandler.List)

	kpiHandler := handlers.NewKPIHandler()
	api.POST("/kpi-indicators", kpiHandler.Create)
	api.PUT("/kpi-indicators/:id", kpiHandler.Update)
	api.DELETE("/kpi-indicators/:id", kpiHandler.Delete)
	api.GET("/kpi-indicators/:id", kpiHandler.Get)
	api.GET("/kpi-indicators", kpiHandler.List)
	api.POST("/kpi-indicators/:id/publish", kpiHandler.Publish)

	assessmentHandler := handlers.NewAssessmentHandler()
	api.POST("/assessment-plans", assessmentHandler.CreatePlan)
	api.PUT("/assessment-plans/:id", assessmentHandler.UpdatePlan)
	api.POST("/assessment-plans/:id/confirm", assessmentHandler.ConfirmPlan)
	api.GET("/assessment-plans/:id", assessmentHandler.GetPlan)
	api.GET("/assessment-plans", assessmentHandler.ListPlans)
	api.POST("/kpi-progress", assessmentHandler.UpdateProgress)
	api.GET("/kpi-progress/:planItemId", assessmentHandler.GetProgress)

	performanceHandler := handlers.NewPerformanceHandler()
	api.POST("/performance-reviews/initiate/:planId", performanceHandler.InitiateReview)
	api.POST("/performance-reviews/:id/self-review", performanceHandler.SubmitSelfReview)
	api.POST("/performance-reviews/:id/manager-review", performanceHandler.SubmitManagerReview)
	api.POST("/performance-reviews/items/:itemId/adjust-score", performanceHandler.AdjustScore)
	api.POST("/performance-reviews/:id/finalize", performanceHandler.FinalizeReview)
	api.GET("/performance-reviews/:id", performanceHandler.GetReview)
	api.GET("/performance-reviews", performanceHandler.ListReviews)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
		})
	})
}
