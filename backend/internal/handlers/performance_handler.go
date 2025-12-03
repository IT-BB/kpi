package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kpi-system/backend/internal/services"
)

type PerformanceHandler struct {
	service *services.PerformanceService
}

func NewPerformanceHandler() *PerformanceHandler {
	return &PerformanceHandler{
		service: services.NewPerformanceService(),
	}
}

func (h *PerformanceHandler) InitiateReview(c *gin.Context) {
	planID, err := strconv.ParseUint(c.Param("planId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plan id"})
		return
	}

	if err := h.service.InitiateReview(uint(planID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "review initiated"})
}

func (h *PerformanceHandler) SubmitSelfReview(c *gin.Context) {
	reviewID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid review id"})
		return
	}

	var req struct {
		Score   float64 `json:"score"`
		Comment string  `json:"comment"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	employeeID := c.GetUint("user_id")

	if err := h.service.SubmitSelfReview(uint(reviewID), employeeID, req.Score, req.Comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "self review submitted"})
}

func (h *PerformanceHandler) SubmitManagerReview(c *gin.Context) {
	reviewID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid review id"})
		return
	}

	var req struct {
		Score   float64 `json:"score"`
		Comment string  `json:"comment"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	managerID := c.GetUint("user_id")

	if err := h.service.SubmitManagerReview(uint(reviewID), managerID, req.Score, req.Comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "manager review submitted"})
}

func (h *PerformanceHandler) AdjustScore(c *gin.Context) {
	reviewItemID, err := strconv.ParseUint(c.Param("itemId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item id"})
		return
	}

	var req struct {
		AdjustedScore float64 `json:"adjusted_score"`
		Reason        string  `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reviewerID := c.GetUint("user_id")

	if err := h.service.AdjustScore(uint(reviewItemID), reviewerID, req.AdjustedScore, req.Reason); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "score adjusted"})
}

func (h *PerformanceHandler) FinalizeReview(c *gin.Context) {
	reviewID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid review id"})
		return
	}

	if err := h.service.FinalizeReview(uint(reviewID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "review finalized"})
}

func (h *PerformanceHandler) GetReview(c *gin.Context) {
	reviewID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid review id"})
		return
	}

	review, err := h.service.GetReview(uint(reviewID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "review not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": review})
}

func (h *PerformanceHandler) ListReviews(c *gin.Context) {
	filters := make(map[string]interface{})

	if employeeID := c.Query("employee_id"); employeeID != "" {
		if id, err := strconv.ParseUint(employeeID, 10, 32); err == nil {
			filters["employee_id"] = uint(id)
		}
	}

	if reviewerID := c.Query("reviewer_id"); reviewerID != "" {
		if id, err := strconv.ParseUint(reviewerID, 10, 32); err == nil {
			filters["reviewer_id"] = uint(id)
		}
	}

	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}

	reviews, err := h.service.ListReviews(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": reviews})
}
