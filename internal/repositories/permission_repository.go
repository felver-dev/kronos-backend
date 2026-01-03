package repositories

import (
	"github.com/mcicare/itsm-backend/database"
	"github.com/mcicare/itsm-backend/internal/models"
)

// PermissionRepository interface pour les opérations sur les permissions
type PermissionRepository interface {
	Create(permission *models.Permission) error
	FindByID(id uint) (*models.Permission, error)
	FindByCode(code string) (*models.Permission, error)
	FindAll() ([]models.Permission, error)
	FindByModule(module string) ([]models.Permission, error)
	Update(permission *models.Permission) error
	Delete(id uint) error
}

// permissionRepository implémente PermissionRepository
type permissionRepository struct{}

// NewPermissionRepository crée une nouvelle instance de PermissionRepository
func NewPermissionRepository() PermissionRepository {
	return &permissionRepository{}
}

// Create crée une nouvelle permission
func (r *permissionRepository) Create(permission *models.Permission) error {
	return database.DB.Create(permission).Error
}

// FindByID trouve une permission par son ID
func (r *permissionRepository) FindByID(id uint) (*models.Permission, error) {
	var permission models.Permission
	err := database.DB.First(&permission, id).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

// FindByCode trouve une permission par son code
func (r *permissionRepository) FindByCode(code string) (*models.Permission, error) {
	var permission models.Permission
	err := database.DB.Where("code = ?", code).First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

// FindAll récupère toutes les permissions
func (r *permissionRepository) FindAll() ([]models.Permission, error) {
	var permissions []models.Permission
	err := database.DB.Find(&permissions).Error
	return permissions, err
}

// FindByModule récupère toutes les permissions d'un module
func (r *permissionRepository) FindByModule(module string) ([]models.Permission, error) {
	var permissions []models.Permission
	err := database.DB.Where("module = ?", module).Find(&permissions).Error
	return permissions, err
}

// Update met à jour une permission
func (r *permissionRepository) Update(permission *models.Permission) error {
	return database.DB.Save(permission).Error
}

// Delete supprime une permission
func (r *permissionRepository) Delete(id uint) error {
	return database.DB.Delete(&models.Permission{}, id).Error
}

