package services

import (
	"errors"
	"time"

	"github.com/kpi-system/backend/internal/database"
	"github.com/kpi-system/backend/internal/models"
	"gorm.io/gorm"
)

type Review360Service struct{}

func NewReview360Service() *Review360Service {
	return &Review360Service{}
}

func (s *Review360Service) CreateCampaign(campaign *models.Review360Campaign) error {
	return database.DB.Create(campaign).Error
}

func (s *Review360Service) CreateAssignment(assignment *models.Review360Assignment) error {
	var campaign models.Review360Campaign
	if err := database.DB.First(&campaign, assignment.CampaignID).Error; err != nil {
		return errors.New("campaign not found")
	}

	if campaign.Status != "active" {
		return errors.New("can only create assignments for active campaigns")
	}

	return database.DB.Create(assignment).Error
}

func (s *Review360Service) AutoGenerateAssignments(campaignID uint, revieweeIDs []uint) error {
	var campaign models.Review360Campaign
	if err := database.DB.First(&campaign, campaignID).Error; err != nil {
		return errors.New("campaign not found")
	}

	return database.DB.Transaction(func(tx *gorm.DB) error {
		for _, revieweeID := range revieweeIDs {
			var reviewee models.User
			if err := tx.Preload("Manager").First(&reviewee, revieweeID).Error; err != nil {
				continue
			}

			if reviewee.ManagerID != nil {
				assignment := &models.Review360Assignment{
					CampaignID:   campaignID,
					RevieweeID:   revieweeID,
					ReviewerID:   *reviewee.ManagerID,
					Relationship: models.RelationshipManager,
					Status:       "pending",
				}
				if err := tx.Create(assignment).Error; err != nil {
					return err
				}
			}

			selfAssignment := &models.Review360Assignment{
				CampaignID:   campaignID,
				RevieweeID:   revieweeID,
				ReviewerID:   revieweeID,
				Relationship: models.RelationshipSelf,
				Status:       "pending",
			}
			if err := tx.Create(selfAssignment).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *Review360Service) SubmitResponse(assignmentID, reviewerID uint, responses []struct {
	Dimension string
	Score     int
	Comment   string
}) error {
	var assignment models.Review360Assignment
	if err := database.DB.First(&assignment, assignmentID).Error; err != nil {
		return errors.New("assignment not found")
	}

	if assignment.ReviewerID != reviewerID {
		return errors.New("unauthorized to submit response")
	}

	if assignment.Status == "completed" {
		return errors.New("assignment already completed")
	}

	return database.DB.Transaction(func(tx *gorm.DB) error {
		for _, resp := range responses {
			response := &models.Review360Response{
				AssignmentID: assignmentID,
				Dimension:    resp.Dimension,
				Score:        resp.Score,
				Comment:      resp.Comment,
			}

			if err := tx.Create(response).Error; err != nil {
				return err
			}
		}

		now := time.Now()
		return tx.Model(&assignment).Updates(map[string]interface{}{
			"status":       "completed",
			"completed_at": now,
		}).Error
	})
}

func (s *Review360Service) GetReport(campaignID, revieweeID uint) (map[string]interface{}, error) {
	var assignments []models.Review360Assignment
	if err := database.DB.
		Where("campaign_id = ? AND reviewee_id = ? AND status = ?", campaignID, revieweeID, "completed").
		Preload("Reviewer").
		Find(&assignments).Error; err != nil {
		return nil, err
	}

	dimensionScores := make(map[string][]int)
	dimensionComments := make(map[string][]string)

	for _, assignment := range assignments {
		var responses []models.Review360Response
		database.DB.Where("assignment_id = ?", assignment.ID).Find(&responses)

		for _, resp := range responses {
			dimensionScores[resp.Dimension] = append(dimensionScores[resp.Dimension], resp.Score)
			if resp.Comment != "" {
				dimensionComments[resp.Dimension] = append(dimensionComments[resp.Dimension], resp.Comment)
			}
		}
	}

	averages := make(map[string]float64)
	for dimension, scores := range dimensionScores {
		sum := 0
		for _, score := range scores {
			sum += score
		}
		averages[dimension] = float64(sum) / float64(len(scores))
	}

	return map[string]interface{}{
		"reviewee_id":       revieweeID,
		"campaign_id":       campaignID,
		"dimension_scores":  averages,
		"dimension_comments": dimensionComments,
		"total_responses":   len(assignments),
	}, nil
}

func (s *Review360Service) GetAssignmentsByReviewer(reviewerID uint) ([]models.Review360Assignment, error) {
	var assignments []models.Review360Assignment
	if err := database.DB.
		Where("reviewer_id = ?", reviewerID).
		Preload("Campaign").
		Preload("Reviewee").
		Find(&assignments).Error; err != nil {
		return nil, err
	}
	return assignments, nil
}

func (s *Review360Service) GetCampaignProgress(campaignID uint) (map[string]interface{}, error) {
	var totalAssignments int64
	var completedAssignments int64

	database.DB.Model(&models.Review360Assignment{}).Where("campaign_id = ?", campaignID).Count(&totalAssignments)
	database.DB.Model(&models.Review360Assignment{}).Where("campaign_id = ? AND status = ?", campaignID, "completed").Count(&completedAssignments)

	completionRate := 0.0
	if totalAssignments > 0 {
		completionRate = float64(completedAssignments) / float64(totalAssignments) * 100
	}

	return map[string]interface{}{
		"campaign_id":          campaignID,
		"total_assignments":    totalAssignments,
		"completed_assignments": completedAssignments,
		"completion_rate":      completionRate,
	}, nil
}
