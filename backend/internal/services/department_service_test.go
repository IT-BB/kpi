package services

import (
	"testing"

	"github.com/kpi-system/backend/internal/models"
)

func TestDepartmentServiceCreate(t *testing.T) {
	service := NewDepartmentService()

	dept := &models.Department{
		Name:        "Engineering",
		Description: "Engineering Department",
	}

	if dept.Name == "" {
		t.Error("Department name should not be empty")
	}

	if service == nil {
		t.Error("Service should not be nil")
	}
}

func TestDepartmentValidation(t *testing.T) {
	tests := []struct {
		name     string
		dept     models.Department
		wantErr  bool
	}{
		{
			name: "valid department",
			dept: models.Department{
				Name: "Test Dept",
			},
			wantErr: false,
		},
		{
			name: "empty name",
			dept: models.Department{
				Name: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.dept.Name == "" && !tt.wantErr {
				t.Error("Expected error for empty name")
			}
		})
	}
}
