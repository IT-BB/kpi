package app

import (
	"enterprise-kpi/internal/config"
	"enterprise-kpi/internal/handler"
	"enterprise-kpi/internal/middleware"
	"enterprise-kpi/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// NewRouter wires gin routes, middleware, and handlers.
func NewRouter(db *gorm.DB, cfg config.Config) *gin.Engine {
	router := gin.Default()
	router.Use(middleware.CORSMiddleware(cfg.AllowOrigins))

	apiHandler := handler.New(db, cfg)

	v1 := router.Group("/api/v1")
	{
		v1.POST("/auth/login", apiHandler.Login)
	}

	secured := v1.Group("")
	secured.Use(middleware.AuthMiddleware(db, cfg))

	secured.GET("/auth/profile", apiHandler.Profile)
	secured.GET("/me", apiHandler.Me)

	// User & Org management
	secured.GET("/users", middleware.MustRoles(models.RoleAdmin, models.RoleManager), apiHandler.ListUsers)
	secured.POST("/users", middleware.MustRoles(models.RoleAdmin), apiHandler.CreateUser)
	secured.PUT("/users/:id", middleware.MustRoles(models.RoleAdmin), apiHandler.UpdateUser)
	secured.DELETE("/users/:id", middleware.MustRoles(models.RoleAdmin), apiHandler.DeleteUser)

	secured.GET("/departments", middleware.MustRoles(models.RoleAdmin, models.RoleManager), apiHandler.ListDepartments)
	secured.POST("/departments", middleware.MustRoles(models.RoleAdmin), apiHandler.CreateDepartment)
	secured.PUT("/departments/:id", middleware.MustRoles(models.RoleAdmin), apiHandler.UpdateDepartment)
	secured.DELETE("/departments/:id", middleware.MustRoles(models.RoleAdmin), apiHandler.DeleteDepartment)

	// KPI templates & cycles
	secured.GET("/kpi/templates", apiHandler.ListTemplates)
	secured.GET("/kpi/templates/:id", apiHandler.GetTemplate)
	secured.POST("/kpi/templates", middleware.MustRoles(models.RoleAdmin, models.RoleManager), apiHandler.CreateTemplate)
	secured.PUT("/kpi/templates/:id", middleware.MustRoles(models.RoleAdmin, models.RoleManager), apiHandler.UpdateTemplate)

	secured.GET("/kpi/cycles", apiHandler.ListCycles)
	secured.POST("/kpi/cycles", middleware.MustRoles(models.RoleAdmin), apiHandler.CreateCycle)

	// KPI assignments workflow
	secured.GET("/kpi/assignments", apiHandler.ListAssignments)
	secured.POST("/kpi/assignments", middleware.MustRoles(models.RoleAdmin, models.RoleManager), apiHandler.CreateAssignment)
	secured.POST("/kpi/assignments/:id/submit", apiHandler.SubmitAssignment)
	secured.POST("/kpi/assignments/:id/review", middleware.MustRoles(models.RoleAdmin, models.RoleManager), apiHandler.ReviewAssignment)
	secured.POST("/kpi/assignments/:id/finalize", middleware.MustRoles(models.RoleAdmin), apiHandler.FinalizeAssignment)

	secured.GET("/kpi/reports/summary", middleware.MustRoles(models.RoleAdmin, models.RoleManager), apiHandler.GetSummary)

	return router
}
