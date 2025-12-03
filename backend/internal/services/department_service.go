package services

import (
    "errors"
    "fmt"

    "github.com/kpi-system/backend/internal/database"
    "github.com/kpi-system/backend/internal/models"
)

type DepartmentService struct{}

func NewDepartmentService() *DepartmentService {
    return &DepartmentService{}
}

func (s *DepartmentService) Create(dept *models.Department) error {
    var existing models.Department
    if err := database.DB.Where("name = ?", dept.Name).First(&existing).Error; err == nil {
        return errors.New("department name already exists")
    }

    if dept.ParentID != nil {
        var parent models.Department
        if err := database.DB.First(&parent, *dept.ParentID).Error; err != nil {
            return errors.New("parent department not found")
        }
        dept.Level = parent.Level + 1
    } else {
        dept.Level = 1
    }

    return database.DB.Create(dept).Error
}

func (s *DepartmentService) Update(id uint, updates map[string]interface{}) error {
    var dept models.Department
    if err := database.DB.First(&dept, id).Error; err != nil {
        return errors.New("department not found")
    }

    if name, ok := updates["name"].(string); ok && name != dept.Name {
        var existing models.Department
        if err := database.DB.Where("name = ? AND id != ?", name, id).First(&existing).Error; err == nil {
            return errors.New("department name already exists")
        }
    }

    return database.DB.Model(&dept).Updates(updates).Error
}

func (s *DepartmentService) Delete(id uint) error {
    var userCount int64
    if err := database.DB.Model(&models.User{}).Where("department_id = ?", id).Count(&userCount).Error; err != nil {
        return err
    }

    if userCount > 0 {
        return fmt.Errorf("cannot delete department with %d users. Please reassign users first", userCount)
    }

    var childCount int64
    if err := database.DB.Model(&models.Department{}).Where("parent_id = ?", id).Count(&childCount).Error; err != nil {
        return err
    }

    if childCount > 0 {
        return errors.New("cannot delete department with child departments")
    }

    return database.DB.Delete(&models.Department{}, id).Error
}

func (s *DepartmentService) GetByID(id uint) (*models.Department, error) {
    var dept models.Department
    if err := database.DB.Preload("Parent").First(&dept, id).Error; err != nil {
        return nil, err
    }
    return &dept, nil
}

func (s *DepartmentService) List() ([]models.Department, error) {
    var depts []models.Department
    if err := database.DB.Preload("Parent").Find(&depts).Error; err != nil {
        return nil, err
    }
    return depts, nil
}
