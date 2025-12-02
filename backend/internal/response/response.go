package response

import "github.com/gin-gonic/gin"

// Success wraps successful payloads in a consistent shape.
func Success(c *gin.Context, status int, data interface{}) {
	c.JSON(status, gin.H{
		"status": "success",
		"data":   data,
	})
}

// Error sends structured error responses.
func Error(c *gin.Context, status int, message string, details interface{}) {
	payload := gin.H{
		"status":  "error",
		"message": message,
	}
	if details != nil {
		payload["details"] = details
	}
	c.JSON(status, payload)
}
