package handler

import (
	"net/http"
	"strings"
	"time"

	"enterprise-kpi/internal/middleware"
	"enterprise-kpi/internal/models"
	"enterprise-kpi/internal/response"
	"enterprise-kpi/pkg/password"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserRequest struct {
	FullName     string `json:"fullName" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	Password     string `json:"password"`
	Role         string `json:"role" binding:"required"`
	Title        string `json:"title"`
	DepartmentID *uint  `json:"departmentId"`
	ManagerID    *uint  `json:"managerId"`
}

func (api *API) ListUsers(c *gin.Context) {
	role := strings.TrimSpace(c.Query("role"))
	var users []models.User
	query := api.DB.Preload("Department").Preload("Manager")
	if role != "" {
		query = query.Where("role = ?", role)
	}
	if err := query.Find(&users).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to list users", err.Error())
		return
	}
	response.Success(c, http.StatusOK, users)
}

func (api *API) CreateUser(c *gin.Context) {
	var req UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid payload", err.Error())
		return
	}
	if req.Password == "" {
		response.Error(c, http.StatusBadRequest, "password is required", nil)
		return
	}
	hash, err := password.Hash(req.Password)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to hash password", err.Error())
		return
	}
	user := models.User{
		FullName:     req.FullName,
		Email:        req.Email,
		PasswordHash: hash,
		Role:         req.Role,
		Title:        req.Title,
		DepartmentID: req.DepartmentID,
		ManagerID:    req.ManagerID,
		JoinedAt:     time.Now(),
	}
	if err := api.DB.Create(&user).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to create user", err.Error())
		return
	}
	response.Success(c, http.StatusCreated, user)
}

func (api *API) UpdateUser(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}

	var user models.User
	if err := api.DB.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Error(c, http.StatusNotFound, "user not found", nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "failed to load user", err.Error())
		return
	}

	var req UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid payload", err.Error())
		return
	}

	user.FullName = req.FullName
	user.Email = req.Email
	user.Role = req.Role
	user.Title = req.Title
	user.DepartmentID = req.DepartmentID
	user.ManagerID = req.ManagerID

	if req.Password != "" {
		hash, err := password.Hash(req.Password)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, "failed to hash password", err.Error())
			return
		}
		user.PasswordHash = hash
	}

	if err := api.DB.Save(&user).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to update user", err.Error())
		return
	}
	response.Success(c, http.StatusOK, user)
}

func (api *API) DeleteUser(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}
	if err := api.DB.Delete(&models.User{}, id).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to delete user", err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}

// Me returns the current user enriched with relationships.
func (api *API) Me(c *gin.Context) {
	actor := middleware.CurrentUser(c)
	if actor == nil {
		response.Error(c, http.StatusUnauthorized, "missing user context", nil)
		return
	}
	var user models.User
	if err := api.DB.Preload("Department").Preload("Manager").First(&user, actor.ID).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to load profile", err.Error())
		return
	}
	response.Success(c, http.StatusOK, user)
}
