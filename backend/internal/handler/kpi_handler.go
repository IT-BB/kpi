package handler

import (
	"errors"
	"net/http"
	"time"

	"enterprise-kpi/internal/middleware"
	"enterprise-kpi/internal/models"
	"enterprise-kpi/internal/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type KPIMetricRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Weight      int     `json:"weight" binding:"required"`
	TargetValue float64 `json:"targetValue"`
	Unit        string  `json:"unit"`
}

type KPITemplateRequest struct {
	Name        string             `json:"name" binding:"required"`
	Description string             `json:"description"`
	IsActive    bool               `json:"isActive"`
	Metrics     []KPIMetricRequest `json:"metrics" binding:"required,dive"`
}

type KPICycleRequest struct {
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description"`
	StartsAt    time.Time `json:"startsAt" binding:"required"`
	EndsAt      time.Time `json:"endsAt" binding:"required"`
	Status      string    `json:"status" binding:"required"`
}

type AssignmentRequest struct {
	CycleID    uint      `json:"cycleId" binding:"required"`
	TemplateID uint      `json:"templateId" binding:"required"`
	AssigneeID uint      `json:"assigneeId" binding:"required"`
	ManagerID  uint      `json:"managerId" binding:"required"`
	DueDate    time.Time `json:"dueDate" binding:"required"`
}

type ScoreInput struct {
	MetricID uint    `json:"metricId" binding:"required"`
	Score    float64 `json:"score"`
	Comment  string  `json:"comment"`
}

type AssignmentSubmissionRequest struct {
	Scores  []ScoreInput `json:"scores"`
	Comment string       `json:"comment"`
}

type AssignmentReviewRequest struct {
	Scores  []ScoreInput `json:"scores"`
	Comment string       `json:"comment"`
}

func (api *API) ListTemplates(c *gin.Context) {
	var templates []models.KPITemplate
	if err := api.DB.Preload("Metrics").Preload("Owner").Find(&templates).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to load templates", err.Error())
		return
	}
	response.Success(c, http.StatusOK, templates)
}

func (api *API) GetTemplate(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}
	var template models.KPITemplate
	if err := api.DB.Preload("Metrics").Preload("Owner").First(&template, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Error(c, http.StatusNotFound, "template not found", nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "failed to load template", err.Error())
		return
	}
	response.Success(c, http.StatusOK, template)
}

func (api *API) CreateTemplate(c *gin.Context) {
	actor := middleware.CurrentUser(c)
	if actor == nil {
		response.Error(c, http.StatusUnauthorized, "missing user context", nil)
		return
	}
	var req KPITemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid payload", err.Error())
		return
	}
	if len(req.Metrics) == 0 {
		response.Error(c, http.StatusBadRequest, "template requires metrics", nil)
		return
	}

	template := models.KPITemplate{
		Name:        req.Name,
		Description: req.Description,
		OwnerID:     actor.ID,
		IsActive:    req.IsActive,
	}
	for _, metric := range req.Metrics {
		template.Metrics = append(template.Metrics, models.KPIMetric{
			Name:        metric.Name,
			Description: metric.Description,
			Weight:      metric.Weight,
			TargetValue: metric.TargetValue,
			Unit:        metric.Unit,
		})
	}

	if err := api.DB.Create(&template).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to create template", err.Error())
		return
	}
	api.recordAudit(actor.ID, "kpi_template", template.ID, "create", req)
	response.Success(c, http.StatusCreated, template)
}

func (api *API) UpdateTemplate(c *gin.Context) {
	actor := middleware.CurrentUser(c)
	if actor == nil {
		response.Error(c, http.StatusUnauthorized, "missing user context", nil)
		return
	}
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}
	var template models.KPITemplate
	if err := api.DB.Preload("Metrics").First(&template, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Error(c, http.StatusNotFound, "template not found", nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "failed to load template", err.Error())
		return
	}
	var req KPITemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid payload", err.Error())
		return
	}
	template.Name = req.Name
	template.Description = req.Description
	template.IsActive = req.IsActive

	if err := api.DB.Save(&template).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to update template", err.Error())
		return
	}
	// Replace metrics for simplicity
	if err := api.DB.Where("template_id = ?", template.ID).Delete(&models.KPIMetric{}).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to refresh metrics", err.Error())
		return
	}
	batch := make([]models.KPIMetric, 0, len(req.Metrics))
	for _, metric := range req.Metrics {
		batch = append(batch, models.KPIMetric{
			TemplateID:  template.ID,
			Name:        metric.Name,
			Description: metric.Description,
			Weight:      metric.Weight,
			TargetValue: metric.TargetValue,
			Unit:        metric.Unit,
		})
	}
	if len(batch) > 0 {
		if err := api.DB.Create(&batch).Error; err != nil {
			response.Error(c, http.StatusInternalServerError, "failed to persist metrics", err.Error())
			return
		}
	}
	api.recordAudit(actor.ID, "kpi_template", template.ID, "update", req)
	response.Success(c, http.StatusOK, template)
}

func (api *API) ListCycles(c *gin.Context) {
	var cycles []models.KPICycle
	if err := api.DB.Order("starts_at DESC").Find(&cycles).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to load cycles", err.Error())
		return
	}
	response.Success(c, http.StatusOK, cycles)
}

func (api *API) CreateCycle(c *gin.Context) {
	var req KPICycleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid payload", err.Error())
		return
	}
	cycle := models.KPICycle{
		Name:        req.Name,
		Description: req.Description,
		StartsAt:    req.StartsAt,
		EndsAt:      req.EndsAt,
		Status:      req.Status,
	}
	if err := api.DB.Create(&cycle).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to create cycle", err.Error())
		return
	}
	response.Success(c, http.StatusCreated, cycle)
}

func (api *API) ListAssignments(c *gin.Context) {
	var assignments []models.KPIAssignment
	query := api.DB.Preload("Template").Preload("Template.Metrics").Preload("Assignee").Preload("Manager").Preload("Cycle").Preload("Scores.Metric").Preload("Comments")

	if assignee := c.Query("assigneeId"); assignee != "" {
		query = query.Where("assignee_id = ?", assignee)
	}
	if manager := c.Query("managerId"); manager != "" {
		query = query.Where("manager_id = ?", manager)
	}
	if cycle := c.Query("cycleId"); cycle != "" {
		query = query.Where("cycle_id = ?", cycle)
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Order("due_date ASC").Find(&assignments).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to load assignments", err.Error())
		return
	}
	response.Success(c, http.StatusOK, assignments)
}

func (api *API) CreateAssignment(c *gin.Context) {
	actor := middleware.CurrentUser(c)
	if actor == nil {
		response.Error(c, http.StatusUnauthorized, "missing user context", nil)
		return
	}
	var req AssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid payload", err.Error())
		return
	}
	var template models.KPITemplate
	if err := api.DB.Preload("Metrics").First(&template, req.TemplateID).Error; err != nil {
		response.Error(c, http.StatusBadRequest, "template not found", err.Error())
		return
	}
	assignment := models.KPIAssignment{
		CycleID:    req.CycleID,
		TemplateID: req.TemplateID,
		AssigneeID: req.AssigneeID,
		ManagerID:  req.ManagerID,
		Status:     models.AssignmentStatusDraft,
		DueDate:    req.DueDate,
	}
	for _, metric := range template.Metrics {
		assignment.Scores = append(assignment.Scores, models.KPIScore{
			MetricID: metric.ID,
		})
	}

	if err := api.DB.Create(&assignment).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to create assignment", err.Error())
		return
	}
	api.recordAudit(actor.ID, "kpi_assignment", assignment.ID, "create", req)
	response.Success(c, http.StatusCreated, assignment)
}

func (api *API) SubmitAssignment(c *gin.Context) {
	actor := middleware.CurrentUser(c)
	if actor == nil {
		response.Error(c, http.StatusUnauthorized, "missing user context", nil)
		return
	}
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}
	assignment, err := api.loadAssignment(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, http.StatusNotFound, "assignment not found", nil)
		} else {
			response.Error(c, http.StatusInternalServerError, "failed to load assignment", err.Error())
		}
		return
	}
	if actor.Role != models.RoleAdmin && assignment.AssigneeID != actor.ID {
		response.Error(c, http.StatusForbidden, "only assignee can submit", nil)
		return
	}
	var req AssignmentSubmissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid payload", err.Error())
		return
	}
	if err := api.applyScores(assignment.ID, req.Scores, false); err != nil {
		response.Error(c, http.StatusBadRequest, "failed to apply scores", err.Error())
		return
	}
	now := time.Now()
	assignment.Status = models.AssignmentStatusSubmitted
	assignment.SubmittedAt = &now
	if err := api.DB.Model(&models.KPIAssignment{}).Where("id = ?", assignment.ID).Updates(map[string]interface{}{
		"status":       assignment.Status,
		"submitted_at": assignment.SubmittedAt,
	}).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to update assignment", err.Error())
		return
	}
	if req.Comment != "" {
		api.DB.Create(&models.KPIComment{
			AssignmentID: assignment.ID,
			AuthorID:     actor.ID,
			Message:      req.Comment,
			Type:         "submission",
		})
	}
	api.recordAudit(actor.ID, "kpi_assignment", assignment.ID, "submit", req)
	response.Success(c, http.StatusOK, gin.H{"status": assignment.Status})
}

func (api *API) ReviewAssignment(c *gin.Context) {
	actor := middleware.CurrentUser(c)
	if actor == nil {
		response.Error(c, http.StatusUnauthorized, "missing user context", nil)
		return
	}
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}
	assignment, err := api.loadAssignment(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, http.StatusNotFound, "assignment not found", nil)
		} else {
			response.Error(c, http.StatusInternalServerError, "failed to load assignment", err.Error())
		}
		return
	}
	if actor.Role != models.RoleAdmin && assignment.ManagerID != actor.ID {
		response.Error(c, http.StatusForbidden, "only manager can review", nil)
		return
	}
	var req AssignmentReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid payload", err.Error())
		return
	}
	if err := api.applyScores(assignment.ID, req.Scores, true); err != nil {
		response.Error(c, http.StatusBadRequest, "failed to apply scores", err.Error())
		return
	}
	now := time.Now()
	assignment.Status = models.AssignmentStatusReviewed
	assignment.ReviewedAt = &now
	if err := api.DB.Model(&models.KPIAssignment{}).Where("id = ?", assignment.ID).Updates(map[string]interface{}{
		"status":      assignment.Status,
		"reviewed_at": assignment.ReviewedAt,
	}).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to update assignment", err.Error())
		return
	}
	if req.Comment != "" {
		api.DB.Create(&models.KPIComment{
			AssignmentID: assignment.ID,
			AuthorID:     actor.ID,
			Message:      req.Comment,
			Type:         "review",
		})
	}
	api.recordAudit(actor.ID, "kpi_assignment", assignment.ID, "review", req)
	response.Success(c, http.StatusOK, gin.H{"status": assignment.Status})
}

func (api *API) FinalizeAssignment(c *gin.Context) {
	actor := middleware.CurrentUser(c)
	if actor == nil {
		response.Error(c, http.StatusUnauthorized, "missing user context", nil)
		return
	}
	if actor.Role != models.RoleAdmin {
		response.Error(c, http.StatusForbidden, "only admins can finalize", nil)
		return
	}
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}
	assignment, err := api.loadAssignment(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, http.StatusNotFound, "assignment not found", nil)
		} else {
			response.Error(c, http.StatusInternalServerError, "failed to load assignment", err.Error())
		}
		return
	}
	now := time.Now()
	assignment.Status = models.AssignmentStatusApproved
	assignment.ApprovedAt = &now
	if err := api.DB.Model(&models.KPIAssignment{}).Where("id = ?", assignment.ID).Updates(map[string]interface{}{
		"status":      assignment.Status,
		"approved_at": assignment.ApprovedAt,
	}).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to finalize assignment", err.Error())
		return
	}
	api.recordAudit(actor.ID, "kpi_assignment", assignment.ID, "finalize", gin.H{"status": assignment.Status})
	response.Success(c, http.StatusOK, gin.H{"status": assignment.Status})
}

func (api *API) applyScores(assignmentID uint, scores []ScoreInput, manager bool) error {
	for _, score := range scores {
		updates := map[string]interface{}{
			"comment": score.Comment,
		}
		if manager {
			updates["manager_score"] = score.Score
		} else {
			updates["self_score"] = score.Score
		}
		if err := api.DB.Model(&models.KPIScore{}).
			Where("assignment_id = ? AND metric_id = ?", assignmentID, score.MetricID).
			Updates(updates).Error; err != nil {
			return err
		}
	}
	return nil
}

func (api *API) loadAssignment(id uint) (models.KPIAssignment, error) {
	var assignment models.KPIAssignment
	err := api.DB.
		Preload("Template").
		Preload("Template.Metrics").
		Preload("Assignee").
		Preload("Manager").
		Preload("Scores").
		Preload("Scores.Metric").
		Preload("Comments").
		First(&assignment, id).Error
	return assignment, err
}
