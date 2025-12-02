package middleware

import (
	"net/http"
	"strings"

	"enterprise-kpi/internal/config"
	"enterprise-kpi/internal/models"
	jwtpkg "enterprise-kpi/pkg/auth"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const currentUserKey = "currentUser"

// AuthMiddleware validates JWT tokens and loads the actor from the database.
func AuthMiddleware(db *gorm.DB, cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "missing bearer token"})
			return
		}

		token := strings.TrimPrefix(header, "Bearer ")
		claims, err := jwtpkg.Parse(cfg.JWTSecret, token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid token"})
			return
		}

		var user models.User
		if err := db.Preload("Department").First(&user, claims.UserID).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "user not found"})
			return
		}

		c.Set(currentUserKey, &user)
		c.Next()
	}
}

// MustRoles ensures the current user has one of the provided roles.
func MustRoles(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		allowed[role] = struct{}{}
	}

	return func(c *gin.Context) {
		raw, exists := c.Get(currentUserKey)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "auth context missing"})
			return
		}
		user := raw.(*models.User)
		if _, ok := allowed[user.Role]; !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "insufficient role"})
			return
		}
		c.Next()
	}
}

// CurrentUser extracts the authenticated user from context.
func CurrentUser(c *gin.Context) *models.User {
	raw, exists := c.Get(currentUserKey)
	if !exists {
		return nil
	}
	if user, ok := raw.(*models.User); ok {
		return user
	}
	return nil
}

// CORSMiddleware enables cross origin requests for the SPA frontend.
func CORSMiddleware(allowOrigin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if allowOrigin == "*" && origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
		} else {
			c.Header("Access-Control-Allow-Origin", allowOrigin)
		}
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
