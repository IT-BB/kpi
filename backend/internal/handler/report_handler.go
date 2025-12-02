package handler

import (
	"net/http"

	"enterprise-kpi/internal/models"
	"enterprise-kpi/internal/response"

	"github.com/gin-gonic/gin"
)

type ReportSummary struct {
	TotalAssignments    int64                 `json:"totalAssignments"`
	StatusBreakdown     map[string]int64      `json:"statusBreakdown"`
	AverageSelfScore    float64               `json:"averageSelfScore"`
	AverageManagerScore float64               `json:"averageManagerScore"`
	DepartmentSummaries []DepartmentReportRow `json:"departments"`
}

type DepartmentReportRow struct {
	Department          string  `json:"department"`
	Headcount           int64   `json:"headcount"`
	AverageSelfScore    float64 `json:"averageSelfScore"`
	AverageManagerScore float64 `json:"averageManagerScore"`
}

// GetSummary aggregates key KPI insights across the organization.
func (api *API) GetSummary(c *gin.Context) {
	var total int64
	if err := api.DB.Model(&models.KPIAssignment{}).Count(&total).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to count assignments", err.Error())
		return
	}

	statuses := map[string]int64{}
	for _, status := range []string{
		models.AssignmentStatusDraft,
		models.AssignmentStatusSubmitted,
		models.AssignmentStatusReviewed,
		models.AssignmentStatusApproved,
	} {
		var count int64
		if err := api.DB.Model(&models.KPIAssignment{}).Where("status = ?", status).Count(&count).Error; err == nil {
			statuses[status] = count
		}
	}

	type avgRow struct {
		AvgSelf    *float64
		AvgManager *float64
	}
	var averages avgRow
	api.DB.Model(&models.KPIScore{}).Select("AVG(self_score) as avg_self, AVG(manager_score) as avg_manager").Scan(&averages)

	var deptRows []DepartmentReportRow
	api.DB.Model(&models.KPIAssignment{}).
		Select("COALESCE(departments.name, 'Unassigned') as department, COUNT(DISTINCT users.id) as headcount, AVG(kpi_scores.self_score) as average_self_score, AVG(kpi_scores.manager_score) as average_manager_score").
		Joins("JOIN users ON users.id = kpi_assignments.assignee_id").
		Joins("LEFT JOIN departments ON departments.id = users.department_id").
		Joins("LEFT JOIN kpi_scores ON kpi_scores.assignment_id = kpi_assignments.id").
		Group("COALESCE(departments.name, 'Unassigned')").
		Scan(&deptRows)

	summary := ReportSummary{
		TotalAssignments:    total,
		StatusBreakdown:     statuses,
		AverageSelfScore:    derefFloat(averages.AvgSelf),
		AverageManagerScore: derefFloat(averages.AvgManager),
		DepartmentSummaries: deptRows,
	}
	response.Success(c, http.StatusOK, summary)
}

func derefFloat(value *float64) float64 {
	if value == nil {
		return 0
	}
	return *value
}
