package handler

import (
	"net/http"
	"strconv"

	"enterprise-kpi/internal/response"

	"github.com/gin-gonic/gin"
)

func parseUintParam(c *gin.Context, key string) (uint, bool) {
	value := c.Param(key)
	id64, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid identifier", err.Error())
		return 0, false
	}
	return uint(id64), true
}
