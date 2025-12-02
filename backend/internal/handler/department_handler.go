package handler

import (
	"net/http"

	"enterprise-kpi/internal/models"
	"enterprise-kpi/internal/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DepartmentRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	ManagerID   *uint  `json:"managerId"`
}

func (api *API) ListDepartments(c *gin.Context) {
	var departments []models.Department
	if err := api.DB.Preload("Manager").Find(&departments).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to fetch departments", err.Error())
		return
	}
	response.Success(c, http.StatusOK, departments)
}

func (api *API) CreateDepartment(c *gin.Context) {
	var req DepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid payload", err.Error())
		return
	}
	dept := models.Department{
		Name:        req.Name,
		Description: req.Description,
		ManagerID:   req.ManagerID,
	}
	if err := api.DB.Create(&dept).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to create department", err.Error())
		return
	}
	response.Success(c, http.StatusCreated, dept)
}

func (api *API) UpdateDepartment(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}

	var dept models.Department
	if err := api.DB.First(&dept, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Error(c, http.StatusNotFound, "department not found", nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "failed to load department", err.Error())
		return
	}

	var req DepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid payload", err.Error())
		return
	}

	dept.Name = req.Name
	dept.Description = req.Description
	dept.ManagerID = req.ManagerID

	if err := api.DB.Save(&dept).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to update department", err.Error())
		return
	}
	response.Success(c, http.StatusOK, dept)
}

func (api *API) DeleteDepartment(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}
	if err := api.DB.Delete(&models.Department{}, id).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to delete department", err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}
