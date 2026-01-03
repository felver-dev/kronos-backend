package repositories

import (
	"github.com/mcicare/itsm-backend/database"
	"github.com/mcicare/itsm-backend/internal/models"
)

// RoleRepository interface pour les opérations sur les rôles
type RoleRepository interface {
	Create(role *models.Role) error
	FindByID(id uint) (*models.Role, error)
	FindByName(name string) (*models.Role, error)
	FindAll() ([]models.Role, error)
	Update(role *models.Role) error
	Delete(id uint) error
}

// roleRepository implémente RoleRepository
type roleRepository struct{}

// NewRoleRepository crée une nouvelle instance de RoleRepository
func NewRoleRepository() RoleRepository {
	return &roleRepository{}
}

// Create crée un nouveau rôle
func (r *roleRepository) Create(role *models.Role) error {
	return database.DB.Create(role).Error
}

// FindByID trouve un rôle par son ID
func (r *roleRepository) FindByID(id uint) (*models.Role, error) {
	var role models.Role
	err := database.DB.First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// FindByName trouve un rôle par son nom
func (r *roleRepository) FindByName(name string) (*models.Role, error) {
	var role models.Role
	err := database.DB.Where("name = ?", name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// FindAll récupère tous les rôles
func (r *roleRepository) FindAll() ([]models.Role, error) {
	var roles []models.Role
	err := database.DB.Find(&roles).Error
	return roles, err
}

// Update met à jour un rôle
func (r *roleRepository) Update(role *models.Role) error {
	return database.DB.Save(role).Error
}

// Delete supprime un rôle (soft delete)
func (r *roleRepository) Delete(id uint) error {
	return database.DB.Delete(&models.Role{}, id).Error
}
