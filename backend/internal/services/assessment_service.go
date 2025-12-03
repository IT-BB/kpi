package services

import (
    "errors"
    "time"

    "github.com/kpi-system/backend/internal/database"
    "github.com/kpi-system/backend/internal/models"
    "gorm.io/gorm"
)

type AssessmentService struct{}

func NewAssessmentService() *AssessmentService {
    return &AssessmentService{}
}

func (s *AssessmentService) CreatePlan(plan *models.AssessmentPlan) error {
    totalWeight := 0.0
    for _, item := range plan.Items {
        totalWeight += item.Weight
    }

    if totalWeight != 100.0 {
        return errors.New("total weight of all indicators must equal 100%")
    }

    return database.DB.Transaction(func(tx *gorm.DB) error {
        if err := tx.Create(plan).Error; err != nil {
            return err
        }

        for i := range plan.Items {
            plan.Items[i].PlanID = plan.ID
            if err := tx.Create(&plan.Items[i]).Error; err != nil {
                return err
            }
        }

        return nil
    })
}

func (s *AssessmentService) UpdatePlan(id uint, updates map[string]interface{}) error {
    var plan models.AssessmentPlan
    if err := database.DB.First(&plan, id).Error; err != nil {
        return errors.New("assessment plan not found")
    }

    if plan.Status == models.StatusCompleted || plan.Status == models.StatusCancelled {
        return errors.New("cannot update completed or cancelled plan")
    }

    return database.DB.Model(&plan).Updates(updates).Error
}

func (s *AssessmentService) ConfirmPlan(planID, employeeID uint) error {
    var plan models.AssessmentPlan
    if err := database.DB.First(&plan, planID).Error; err != nil {
        return errors.New("assessment plan not found")
    }

    if plan.EmployeeID != employeeID {
        return errors.New("only the assigned employee can confirm the plan")
    }

    now := time.Now()
    return database.DB.Model(&plan).Updates(map[string]interface{}{
        "status":       models.StatusConfirmed,
        "confirmed_at": now,
    }).Error
}

func (s *AssessmentService) UpdateProgress(progress *models.KPIProgress) error {
    var planItem models.AssessmentPlanItem
    if err := database.DB.First(&planItem, progress.PlanItemID).Error; err != nil {
        return errors.New("plan item not found")
    }

    if planItem.TargetValue > 0 {
        progress.CompletionRate = (progress.CurrentValue / planItem.TargetValue) * 100
    }

    return database.DB.Create(progress).Error
}

func (s *AssessmentService) GetPlanByID(id uint) (*models.AssessmentPlan, error) {
    var plan models.AssessmentPlan
    if err := database.DB.
        Preload("Employee").
        Preload("Manager").
        Preload("Items.Indicator").
        First(&plan, id).Error; err != nil {
        return nil, err
    }
    return &plan, nil
}

func (s *AssessmentService) ListPlans(filters map[string]interface{}) ([]models.AssessmentPlan, error) {
    var plans []models.AssessmentPlan
    query := database.DB.
        Preload("Employee").
        Preload("Manager").
        Preload("Items")

    if employeeID, ok := filters["employee_id"].(uint); ok {
        query = query.Where("employee_id = ?", employeeID)
    }

    if managerID, ok := filters["manager_id"].(uint); ok {
        query = query.Where("manager_id = ?", managerID)
    }

    if status, ok := filters["status"].(string); ok && status != "" {
        query = query.Where("status = ?", status)
    }

    if err := query.Find(&plans).Error; err != nil {
        return nil, err
    }

    return plans, nil
}

func (s *AssessmentService) GetProgress(planItemID uint) ([]models.KPIProgress, error) {
    var progress []models.KPIProgress
    if err := database.DB.
        Where("plan_item_id = ?", planItemID).
        Preload("Updater").
        Order("created_at DESC").
        Find(&progress).Error; err != nil {
        return nil, err
    }
    return progress, nil
}

func (s *AssessmentService) CalculateScore(planItemID uint) (float64, error) {
    var planItem models.AssessmentPlanItem
    if err := database.DB.Preload("Indicator").First(&planItem, planItemID).Error; err != nil {
        return 0, err
    }

    var latestProgress models.KPIProgress
    if err := database.DB.
        Where("plan_item_id = ?", planItemID).
        Order("created_at DESC").
        First(&latestProgress).Error; err != nil {
        return 0, nil
    }

    completionRate := latestProgress.CompletionRate
    score := 0.0

    switch planItem.Indicator.CalculationType {
    case models.CalculationTypeLinear:
        score = (completionRate / 100) * 100
    case models.CalculationTypeStepped:
        if completionRate >= 100 {
            score = 100
        } else if completionRate >= 80 {
            score = 80
        } else if completionRate >= 60 {
            score = 60
        } else {
            score = completionRate
        }
    default:
        score = completionRate
    }

    return score * planItem.Weight / 100, nil
}
