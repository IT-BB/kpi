package handler

import (
	"net/http"

	"enterprise-kpi/internal/middleware"
	"enterprise-kpi/internal/models"
	"enterprise-kpi/internal/response"
	jwtpkg "enterprise-kpi/pkg/auth"
	"enterprise-kpi/pkg/password"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// LoginRequest captures credentials from the SPA.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Login authenticates a user returning a JWT token.
func (api *API) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid payload", err.Error())
		return
	}

	var user models.User
	if err := api.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Error(c, http.StatusUnauthorized, "invalid credentials", nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "failed to load user", err.Error())
		return
	}

	if !password.Verify(user.PasswordHash, req.Password) {
		response.Error(c, http.StatusUnauthorized, "invalid credentials", nil)
		return
	}

	token, err := jwtpkg.Generate(api.Cfg.JWTSecret, user.ID, user.Email, user.Role)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to issue token", err.Error())
		return
	}

	response.Success(c, http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":           user.ID,
			"uuid":         user.UUID,
			"fullName":     user.FullName,
			"email":        user.Email,
			"role":         user.Role,
			"departmentId": user.DepartmentID,
		},
	})
}

// Profile returns the authenticated user's details.
func (api *API) Profile(c *gin.Context) {
	user := middleware.CurrentUser(c)
	if user == nil {
		response.Error(c, http.StatusUnauthorized, "missing user context", nil)
		return
	}
	response.Success(c, http.StatusOK, user)
}
