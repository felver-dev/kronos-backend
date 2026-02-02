package repositories

import (
	"github.com/mcicare/itsm-backend/database"
	"github.com/mcicare/itsm-backend/internal/models"
)

// DepartmentRepository interface pour les opérations sur les départements
type DepartmentRepository interface {
	Create(department *models.Department) error
	FindByID(id uint) (*models.Department, error)
	FindByCode(code string) (*models.Department, error)
	FindAll() ([]models.Department, error)
	FindActive() ([]models.Department, error)
	FindByOfficeID(officeID uint) ([]models.Department, error)
	FindByFilialeID(filialeID uint) ([]models.Department, error)
	Update(department *models.Department) error
	Delete(id uint) error
}

// departmentRepository implémente DepartmentRepository
type departmentRepository struct{}

// NewDepartmentRepository crée une nouvelle instance de DepartmentRepository
func NewDepartmentRepository() DepartmentRepository {
	return &departmentRepository{}
}

// Create crée un nouveau département
func (r *departmentRepository) Create(department *models.Department) error {
	return database.DB.Create(department).Error
}

// FindByID trouve un département par son ID
func (r *departmentRepository) FindByID(id uint) (*models.Department, error) {
	var department models.Department
	err := database.DB.Preload("Office").Preload("Filiale").First(&department, id).Error
	if err != nil {
		return nil, err
	}
	return &department, nil
}

// FindByCode trouve un département par son code
func (r *departmentRepository) FindByCode(code string) (*models.Department, error) {
	var department models.Department
	err := database.DB.Preload("Office").Preload("Filiale").Where("code = ?", code).First(&department).Error
	if err != nil {
		return nil, err
	}
	return &department, nil
}

// FindAll récupère tous les départements
func (r *departmentRepository) FindAll() ([]models.Department, error) {
	var departments []models.Department
	err := database.DB.Preload("Office").Preload("Filiale").Order("name ASC").Find(&departments).Error
	return departments, err
}

// FindActive récupère tous les départements actifs
func (r *departmentRepository) FindActive() ([]models.Department, error) {
	var departments []models.Department
	err := database.DB.Preload("Office").Preload("Filiale").Where("is_active = ?", true).Order("name ASC").Find(&departments).Error
	return departments, err
}

// FindByOfficeID récupère les départements d'un siège
func (r *departmentRepository) FindByOfficeID(officeID uint) ([]models.Department, error) {
	var departments []models.Department
	err := database.DB.Preload("Office").Preload("Filiale").Where("office_id = ?", officeID).Order("name ASC").Find(&departments).Error
	return departments, err
}

// FindByFilialeID récupère les départements d'une filiale
func (r *departmentRepository) FindByFilialeID(filialeID uint) ([]models.Department, error) {
	var departments []models.Department
	err := database.DB.Preload("Office").Preload("Filiale").Where("filiale_id = ? AND is_active = ?", filialeID, true).Order("name ASC").Find(&departments).Error
	return departments, err
}

// Update met à jour un département
func (r *departmentRepository) Update(department *models.Department) error {
	return database.DB.Save(department).Error
}

// Delete supprime un département (soft delete)
func (r *departmentRepository) Delete(id uint) error {
	return database.DB.Delete(&models.Department{}, id).Error
}
