package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kpi-system/backend/internal/models"
	"github.com/kpi-system/backend/internal/services"
)

type AssessmentHandler struct {
	service *services.AssessmentService
}

func NewAssessmentHandler() *AssessmentHandler {
	return &AssessmentHandler{
		service: services.NewAssessmentService(),
	}
}

func (h *AssessmentHandler) CreatePlan(c *gin.Context) {
	var req models.AssessmentPlan
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreatePlan(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": req})
}

func (h *AssessmentHandler) UpdatePlan(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdatePlan(uint(id), updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "plan updated"})
}

func (h *AssessmentHandler) ConfirmPlan(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	employeeID := c.GetUint("user_id")

	if err := h.service.ConfirmPlan(uint(id), employeeID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "plan confirmed"})
}

func (h *AssessmentHandler) GetPlan(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	plan, err := h.service.GetPlanByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "plan not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": plan})
}

func (h *AssessmentHandler) ListPlans(c *gin.Context) {
	filters := make(map[string]interface{})

	if employeeID := c.Query("employee_id"); employeeID != "" {
		if id, err := strconv.ParseUint(employeeID, 10, 32); err == nil {
			filters["employee_id"] = uint(id)
		}
	}

	if managerID := c.Query("manager_id"); managerID != "" {
		if id, err := strconv.ParseUint(managerID, 10, 32); err == nil {
			filters["manager_id"] = uint(id)
		}
	}

	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}

	plans, err := h.service.ListPlans(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": plans})
}

func (h *AssessmentHandler) UpdateProgress(c *gin.Context) {
	var req models.KPIProgress
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.UpdatedBy = c.GetUint("user_id")

	if err := h.service.UpdateProgress(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": req})
}

func (h *AssessmentHandler) GetProgress(c *gin.Context) {
	planItemID, err := strconv.ParseUint(c.Param("planItemId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plan item id"})
		return
	}

	progress, err := h.service.GetProgress(uint(planItemID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": progress})
}
