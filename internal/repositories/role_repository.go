package repositories

import (
	"errors"

	"github.com/mcicare/itsm-backend/database"
	"github.com/mcicare/itsm-backend/internal/models"
)

// RoleRepository interface pour les opérations sur les rôles
type RoleRepository interface {
	Create(role *models.Role) error
	FindByID(id uint) (*models.Role, error)
	FindByName(name string) (*models.Role, error)
	FindAll() ([]models.Role, error)
	FindByFilialeOrGlobal(filialeID *uint) ([]models.Role, error)
	FindByDepartmentID(departmentID uint) ([]models.Role, error)
	FindByCreatedByOrUsedInFiliale(createdByID uint, filialeID uint) ([]models.Role, error)
	FindByCreatedByID(createdByID uint) ([]models.Role, error)
	Update(role *models.Role) error
	Delete(id uint) error
	GetPermissionsByRoleID(roleID uint) ([]string, error)              // Récupère les codes des permissions d'un rôle
	UpdateRolePermissions(roleID uint, permissionCodes []string) error // Met à jour les permissions d'un rôle
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
	err := database.DB.
		Preload("CreatedBy").
		Preload("Filiale").
		First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// FindByName trouve un rôle par son nom
func (r *roleRepository) FindByName(name string) (*models.Role, error) {
	var role models.Role
	err := database.DB.
		Preload("CreatedBy").
		Preload("Filiale").
		Where("name = ?", name).First(&role).Error
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

// FindByFilialeOrGlobal récupère les rôles globaux (filiale_id IS NULL) ou de la filiale donnée
func (r *roleRepository) FindByFilialeOrGlobal(filialeID *uint) ([]models.Role, error) {
	var roles []models.Role
	q := database.DB.Model(&models.Role{})
	if filialeID == nil {
		q = q.Where("filiale_id IS NULL")
	} else {
		q = q.Where("filiale_id IS NULL OR filiale_id = ?", *filialeID)
	}
	err := q.Find(&roles).Error
	return roles, err
}

// FindByDepartmentID récupère les rôles assignés à au moins un utilisateur du département donné
func (r *roleRepository) FindByDepartmentID(departmentID uint) ([]models.Role, error) {
	var roleIDs []uint
	if err := database.DB.Model(&models.User{}).Where("department_id = ?", departmentID).Distinct("role_id").Pluck("role_id", &roleIDs).Error; err != nil {
		return nil, err
	}
	if len(roleIDs) == 0 {
		return []models.Role{}, nil
	}
	var roles []models.Role
	if err := database.DB.Where("id IN ?", roleIDs).Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

// FindByCreatedByOrUsedInFiliale récupère les rôles créés par l'utilisateur OU utilisés par au moins un utilisateur de la filiale
func (r *roleRepository) FindByCreatedByOrUsedInFiliale(createdByID uint, filialeID uint) ([]models.Role, error) {
	var roleIDs []uint
	if err := database.DB.Model(&models.User{}).Where("filiale_id = ?", filialeID).Distinct("role_id").Pluck("role_id", &roleIDs).Error; err != nil {
		return nil, err
	}
	var roles []models.Role
	q := database.DB.Model(&models.Role{}).Where("created_by_id = ?", createdByID)
	if len(roleIDs) > 0 {
		q = database.DB.Model(&models.Role{}).Where("created_by_id = ? OR id IN ?", createdByID, roleIDs)
	}
	if err := q.Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

// FindByCreatedByID récupère les rôles créés par un utilisateur
func (r *roleRepository) FindByCreatedByID(createdByID uint) ([]models.Role, error) {
	var roles []models.Role
	err := database.DB.
		Where("created_by_id = ?", createdByID).
		Find(&roles).Error
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

// GetPermissionsByRoleID récupère les codes des permissions associées à un rôle
func (r *roleRepository) GetPermissionsByRoleID(roleID uint) ([]string, error) {
	var rolePermissions []models.RolePermission
	err := database.DB.
		Where("role_id = ?", roleID).
		Preload("Permission").
		Find(&rolePermissions).Error

	if err != nil {
		return nil, err
	}

	permissions := make([]string, len(rolePermissions))
	for i, rp := range rolePermissions {
		permissions[i] = rp.Permission.Code
	}

	return permissions, nil
}

// UpdateRolePermissions met à jour les permissions d'un rôle
// Supprime toutes les permissions existantes et ajoute les nouvelles
func (r *roleRepository) UpdateRolePermissions(roleID uint, permissionCodes []string) error {
	// Vérifier que le rôle existe
	_, err := r.FindByID(roleID)
	if err != nil {
		return err
	}

	// Supprimer toutes les permissions existantes pour ce rôle
	if err := database.DB.Where("role_id = ?", roleID).Delete(&models.RolePermission{}).Error; err != nil {
		return err
	}

	// Si aucune permission n'est fournie, on a terminé
	if len(permissionCodes) == 0 {
		return nil
	}

	// Récupérer les IDs des permissions à partir de leurs codes
	var permissions []models.Permission
	if err := database.DB.Where("code IN ?", permissionCodes).Find(&permissions).Error; err != nil {
		return err
	}

	// Vérifier que toutes les permissions existent
	if len(permissions) != len(permissionCodes) {
		return errors.New("certaines permissions sont introuvables")
	}

	// Créer les associations rôle-permission
	rolePermissions := make([]models.RolePermission, len(permissions))
	for i, perm := range permissions {
		rolePermissions[i] = models.RolePermission{
			RoleID:       roleID,
			PermissionID: perm.ID,
		}
	}

	// Insérer toutes les associations en une seule transaction
	if err := database.DB.Create(&rolePermissions).Error; err != nil {
		return err
	}

	return nil
}
