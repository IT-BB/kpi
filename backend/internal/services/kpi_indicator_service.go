package services

import (
	"errors"

	"github.com/kpi-system/backend/internal/database"
	"github.com/kpi-system/backend/internal/models"
)

type KPIIndicatorService struct{}

func NewKPIIndicatorService() *KPIIndicatorService {
	return &KPIIndicatorService{}
}

func (s *KPIIndicatorService) Create(indicator *models.KPIIndicator) error {
	return database.DB.Create(indicator).Error
}

func (s *KPIIndicatorService) Update(id uint, updates map[string]interface{}) error {
	var indicator models.KPIIndicator
	if err := database.DB.First(&indicator, id).Error; err != nil {
		return errors.New("indicator not found")
	}

	return database.DB.Model(&indicator).Updates(updates).Error
}

func (s *KPIIndicatorService) Delete(id uint) error {
	var planItemCount int64
	if err := database.DB.Model(&models.AssessmentPlanItem{}).Where("indicator_id = ?", id).Count(&planItemCount).Error; err != nil {
		return err
	}

	if planItemCount > 0 {
		return errors.New("cannot delete indicator that is referenced by assessment plans")
	}

	return database.DB.Delete(&models.KPIIndicator{}, id).Error
}

func (s *KPIIndicatorService) GetByID(id uint) (*models.KPIIndicator, error) {
	var indicator models.KPIIndicator
	if err := database.DB.Preload("Creator").First(&indicator, id).Error; err != nil {
		return nil, err
	}
	return &indicator, nil
}

func (s *KPIIndicatorService) List(filters map[string]interface{}) ([]models.KPIIndicator, error) {
	var indicators []models.KPIIndicator
	query := database.DB.Preload("Creator")

	if category, ok := filters["category"].(string); ok && category != "" {
		query = query.Where("category = ?", category)
	}

	if isPublished, ok := filters["is_published"].(bool); ok {
		query = query.Where("is_published = ?", isPublished)
	}

	if err := query.Find(&indicators).Error; err != nil {
		return nil, err
	}

	return indicators, nil
}

func (s *KPIIndicatorService) Publish(id uint) error {
	var indicator models.KPIIndicator
	if err := database.DB.First(&indicator, id).Error; err != nil {
		return errors.New("indicator not found")
	}

	if indicator.Name == "" || indicator.DataType == "" || indicator.CalculationType == "" {
		return errors.New("indicator is incomplete and cannot be published")
	}

	return database.DB.Model(&indicator).Update("is_published", true).Error
}
