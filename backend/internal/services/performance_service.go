package services

import (
	"errors"
	"time"

	"github.com/kpi-system/backend/internal/database"
	"github.com/kpi-system/backend/internal/models"
	"gorm.io/gorm"
)

type PerformanceService struct{}

func NewPerformanceService() *PerformanceService {
	return &PerformanceService{}
}

func (s *PerformanceService) InitiateReview(planID uint) error {
	var plan models.AssessmentPlan
	if err := database.DB.Preload("Items").First(&plan, planID).Error; err != nil {
		return errors.New("assessment plan not found")
	}

	if plan.Status != models.StatusConfirmed {
		return errors.New("can only initiate review for confirmed plans")
	}

	return database.DB.Transaction(func(tx *gorm.DB) error {
		review := &models.PerformanceReview{
			PlanID:     planID,
			EmployeeID: plan.EmployeeID,
			ReviewerID: plan.ManagerID,
			Status:     "pending",
		}

		if err := tx.Create(review).Error; err != nil {
			return err
		}

		for _, item := range plan.Items {
			var latestProgress models.KPIProgress
			tx.Where("plan_item_id = ?", item.ID).Order("created_at DESC").First(&latestProgress)

			reviewItem := &models.PerformanceReviewItem{
				ReviewID:       review.ID,
				PlanItemID:     item.ID,
				ActualValue:    latestProgress.CurrentValue,
				CompletionRate: latestProgress.CompletionRate,
			}

			if item.TargetValue > 0 {
				reviewItem.CompletionRate = (latestProgress.CurrentValue / item.TargetValue) * 100
			}

			reviewItem.AutoScore = (reviewItem.CompletionRate / 100) * item.Weight

			if err := tx.Create(reviewItem).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *PerformanceService) SubmitSelfReview(reviewID, employeeID uint, score float64, comment string) error {
	var review models.PerformanceReview
	if err := database.DB.First(&review, reviewID).Error; err != nil {
		return errors.New("review not found")
	}

	if review.EmployeeID != employeeID {
		return errors.New("only the employee can submit self review")
	}

	return database.DB.Model(&review).Updates(map[string]interface{}{
		"self_score":   score,
		"self_comment": comment,
	}).Error
}

func (s *PerformanceService) SubmitManagerReview(reviewID, managerID uint, score float64, comment string) error {
	var review models.PerformanceReview
	if err := database.DB.First(&review, reviewID).Error; err != nil {
		return errors.New("review not found")
	}

	if review.ReviewerID != managerID {
		return errors.New("only the assigned manager can submit review")
	}

	return database.DB.Model(&review).Updates(map[string]interface{}{
		"manager_score":   score,
		"manager_comment": comment,
	}).Error
}

func (s *PerformanceService) AdjustScore(reviewItemID, reviewerID uint, adjustedScore float64, reason string) error {
	var reviewItem models.PerformanceReviewItem
	if err := database.DB.Preload("Review").First(&reviewItem, reviewItemID).Error; err != nil {
		return errors.New("review item not found")
	}

	if reviewItem.Review.ReviewerID != reviewerID {
		return errors.New("unauthorized to adjust score")
	}

	if adjustedScore < 0 || adjustedScore > 100 {
		return errors.New("adjusted score must be between 0 and 100")
	}

	return database.DB.Model(&reviewItem).Updates(map[string]interface{}{
		"adjusted_score":    adjustedScore,
		"adjustment_reason": reason,
	}).Error
}

func (s *PerformanceService) FinalizeReview(reviewID uint) error {
	var review models.PerformanceReview
	if err := database.DB.Preload("Items").First(&review, reviewID).Error; err != nil {
		return errors.New("review not found")
	}

	totalScore := 0.0
	for _, item := range review.Items {
		if item.AdjustedScore != nil {
			totalScore += *item.AdjustedScore
		} else {
			totalScore += item.AutoScore
		}
	}

	performanceLevel := s.calculatePerformanceLevel(totalScore)
	now := time.Now()

	return database.DB.Model(&review).Updates(map[string]interface{}{
		"final_score":       totalScore,
		"performance_level": performanceLevel,
		"status":            "completed",
		"submitted_at":      now,
	}).Error
}

func (s *PerformanceService) calculatePerformanceLevel(score float64) string {
	if score >= 90 {
		return models.PerformanceLevelOutstanding
	} else if score >= 80 {
		return models.PerformanceLevelExcellent
	} else if score >= 70 {
		return models.PerformanceLevelGood
	} else if score >= 60 {
		return models.PerformanceLevelFair
	}
	return models.PerformanceLevelPoor
}

func (s *PerformanceService) GetReview(id uint) (*models.PerformanceReview, error) {
	var review models.PerformanceReview
	if err := database.DB.
		Preload("Plan").
		Preload("Employee").
		Preload("Reviewer").
		Preload("Items.PlanItem.Indicator").
		First(&review, id).Error; err != nil {
		return nil, err
	}
	return &review, nil
}

func (s *PerformanceService) ListReviews(filters map[string]interface{}) ([]models.PerformanceReview, error) {
	var reviews []models.PerformanceReview
	query := database.DB.
		Preload("Employee").
		Preload("Reviewer").
		Preload("Plan")

	if employeeID, ok := filters["employee_id"].(uint); ok {
		query = query.Where("employee_id = ?", employeeID)
	}

	if reviewerID, ok := filters["reviewer_id"].(uint); ok {
		query = query.Where("reviewer_id = ?", reviewerID)
	}

	if status, ok := filters["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Find(&reviews).Error; err != nil {
		return nil, err
	}

	return reviews, nil
}
