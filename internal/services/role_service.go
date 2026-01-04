package services

import (
	"errors"

	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/models"
	"github.com/mcicare/itsm-backend/internal/repositories"
)

// RoleService interface pour les opérations sur les rôles
type RoleService interface {
	Create(req dto.CreateRoleRequest, createdByID uint) (*dto.RoleDTO, error)
	GetByID(id uint) (*dto.RoleDTO, error)
	GetAll() ([]dto.RoleDTO, error)
	Update(id uint, req dto.UpdateRoleRequest, updatedByID uint) (*dto.RoleDTO, error)
	Delete(id uint) error
}

// roleService implémente RoleService
type roleService struct {
	roleRepo repositories.RoleRepository
	userRepo repositories.UserRepository
}

// NewRoleService crée une nouvelle instance de RoleService
func NewRoleService(
	roleRepo repositories.RoleRepository,
	userRepo repositories.UserRepository,
) RoleService {
	return &roleService{
		roleRepo: roleRepo,
		userRepo: userRepo,
	}
}

// Create crée un nouveau rôle
func (s *roleService) Create(req dto.CreateRoleRequest, createdByID uint) (*dto.RoleDTO, error) {
	// Vérifier que le nom n'existe pas déjà
	existingRole, _ := s.roleRepo.FindByName(req.Name)
	if existingRole != nil {
		return nil, errors.New("un rôle avec ce nom existe déjà")
	}

	role := &models.Role{
		Name:        req.Name,
		Description: req.Description,
		IsSystem:    false, // Les rôles créés via API ne sont pas des rôles système
	}

	if err := s.roleRepo.Create(role); err != nil {
		return nil, errors.New("erreur lors de la création du rôle")
	}

	createdRole, err := s.roleRepo.FindByID(role.ID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération du rôle créé")
	}

	roleDTO := s.roleToDTO(createdRole)
	return &roleDTO, nil
}

// GetByID récupère un rôle par son ID
func (s *roleService) GetByID(id uint) (*dto.RoleDTO, error) {
	role, err := s.roleRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("rôle introuvable")
	}

	roleDTO := s.roleToDTO(role)
	return &roleDTO, nil
}

// GetAll récupère tous les rôles
func (s *roleService) GetAll() ([]dto.RoleDTO, error) {
	roles, err := s.roleRepo.FindAll()
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des rôles")
	}

	var roleDTOs []dto.RoleDTO
	for _, role := range roles {
		roleDTOs = append(roleDTOs, s.roleToDTO(&role))
	}

	return roleDTOs, nil
}

// Update met à jour un rôle
func (s *roleService) Update(id uint, req dto.UpdateRoleRequest, updatedByID uint) (*dto.RoleDTO, error) {
	role, err := s.roleRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("rôle introuvable")
	}

	// Vérifier que ce n'est pas un rôle système (ne peut pas être modifié)
	if role.IsSystem {
		return nil, errors.New("impossible de modifier un rôle système")
	}

	// Vérifier que le nouveau nom n'existe pas déjà (si fourni)
	if req.Name != "" && req.Name != role.Name {
		existingRole, _ := s.roleRepo.FindByName(req.Name)
		if existingRole != nil {
			return nil, errors.New("un rôle avec ce nom existe déjà")
		}
		role.Name = req.Name
	}

	if req.Description != "" {
		role.Description = req.Description
	}

	if err := s.roleRepo.Update(role); err != nil {
		return nil, errors.New("erreur lors de la mise à jour du rôle")
	}

	updatedRole, err := s.roleRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération du rôle mis à jour")
	}

	roleDTO := s.roleToDTO(updatedRole)
	return &roleDTO, nil
}

// Delete supprime un rôle
func (s *roleService) Delete(id uint) error {
	role, err := s.roleRepo.FindByID(id)
	if err != nil {
		return errors.New("rôle introuvable")
	}

	// Vérifier que ce n'est pas un rôle système
	if role.IsSystem {
		return errors.New("impossible de supprimer un rôle système")
	}

	// Vérifier qu'aucun utilisateur n'utilise ce rôle
	users, err := s.userRepo.FindByRole(id)
	if err == nil && len(users) > 0 {
		return errors.New("impossible de supprimer un rôle utilisé par des utilisateurs")
	}

	if err := s.roleRepo.Delete(id); err != nil {
		return errors.New("erreur lors de la suppression du rôle")
	}

	return nil
}

// roleToDTO convertit un modèle Role en DTO
func (s *roleService) roleToDTO(role *models.Role) dto.RoleDTO {
	return dto.RoleDTO{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		IsSystem:    role.IsSystem,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}
}

