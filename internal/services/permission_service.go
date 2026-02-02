package services

import (
	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/models"
	"github.com/mcicare/itsm-backend/internal/repositories"
)

// PermissionService interface pour les opérations sur les permissions
type PermissionService interface {
	GetAll() ([]dto.PermissionDTO, error)
	GetByModule(module string) ([]dto.PermissionDTO, error)
	GetByCode(code string) (*dto.PermissionDTO, error)
}

// permissionService implémente PermissionService
type permissionService struct {
	permissionRepo repositories.PermissionRepository
}

// NewPermissionService crée une nouvelle instance de PermissionService
func NewPermissionService(permissionRepo repositories.PermissionRepository) PermissionService {
	return &permissionService{
		permissionRepo: permissionRepo,
	}
}

// GetAll récupère toutes les permissions
func (s *permissionService) GetAll() ([]dto.PermissionDTO, error) {
	permissions, err := s.permissionRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var permissionDTOs []dto.PermissionDTO
	for _, perm := range permissions {
		permissionDTOs = append(permissionDTOs, s.permissionToDTO(&perm))
	}

	return permissionDTOs, nil
}

// GetByModule récupère les permissions d'un module
func (s *permissionService) GetByModule(module string) ([]dto.PermissionDTO, error) {
	permissions, err := s.permissionRepo.FindByModule(module)
	if err != nil {
		return nil, err
	}

	var permissionDTOs []dto.PermissionDTO
	for _, perm := range permissions {
		permissionDTOs = append(permissionDTOs, s.permissionToDTO(&perm))
	}

	return permissionDTOs, nil
}

// GetByCode récupère une permission par son code
func (s *permissionService) GetByCode(code string) (*dto.PermissionDTO, error) {
	permission, err := s.permissionRepo.FindByCode(code)
	if err != nil {
		return nil, err
	}

	dto := s.permissionToDTO(permission)
	return &dto, nil
}

// permissionToDTO convertit un modèle Permission en DTO
func (s *permissionService) permissionToDTO(permission *models.Permission) dto.PermissionDTO {
	return dto.PermissionDTO{
		ID:          permission.ID,
		Name:        permission.Name,
		Code:        permission.Code,
		Description: permission.Description,
		Module:      permission.Module,
		CreatedAt:   permission.CreatedAt,
	}
}
