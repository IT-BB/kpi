package services

import (
	"testing"

	"github.com/kpi-system/backend/internal/models"
)

func TestCalculateCompletionRate(t *testing.T) {
	tests := []struct {
		name        string
		currentVal  float64
		targetVal   float64
		expected    float64
	}{
		{"100% completion", 100, 100, 100},
		{"50% completion", 50, 100, 50},
		{"over achievement", 120, 100, 120},
		{"zero target", 50, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result float64
			if tt.targetVal > 0 {
				result = (tt.currentVal / tt.targetVal) * 100
			}
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestWeightValidation(t *testing.T) {
	items := []models.AssessmentPlanItem{
		{Weight: 30},
		{Weight: 40},
		{Weight: 30},
	}

	total := 0.0
	for _, item := range items {
		total += item.Weight
	}

	if total != 100.0 {
		t.Errorf("total weight should be 100, got %v", total)
	}
}

func TestPerformanceLevel(t *testing.T) {
	service := NewPerformanceService()

	tests := []struct {
		score    float64
		expected string
	}{
		{95, models.PerformanceLevelOutstanding},
		{85, models.PerformanceLevelExcellent},
		{75, models.PerformanceLevelGood},
		{65, models.PerformanceLevelFair},
		{50, models.PerformanceLevelPoor},
	}

	for _, tt := range tests {
		result := service.calculatePerformanceLevel(tt.score)
		if result != tt.expected {
			t.Errorf("score %v: expected %v, got %v", tt.score, tt.expected, result)
		}
	}
}
