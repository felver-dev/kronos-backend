package services

import (
	"errors"
	"strings"

	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/models"
	"github.com/mcicare/itsm-backend/internal/repositories"
)

// RoleService interface pour les opérations sur les rôles
type RoleService interface {
	Create(req dto.CreateRoleRequest, createdByID uint) (*dto.RoleDTO, error)
	GetByID(id uint) (*dto.RoleDTO, error)
	GetAll() ([]dto.RoleDTO, error)
	GetAllForAssignment(filialeID *uint, departmentID *uint, viewMode string) ([]dto.RoleDTO, error)
	Update(id uint, req dto.UpdateRoleRequest, updatedByID uint, canManageAllRoles bool) (*dto.RoleDTO, error)
	Delete(id uint, deletedByID uint, canManageAllRoles bool) error
	GetRolePermissions(roleID uint) ([]string, error)                                                            // Récupère les permissions d'un rôle
	UpdateRolePermissions(roleID uint, permissionCodes []string, updatedByID uint, canManageAllRoles bool) error // Met à jour les permissions d'un rôle
	GetAssignablePermissions(userID uint) ([]string, error)                                                      // Récupère les permissions que l'utilisateur peut déléguer
	GetMyDelegations(userID uint) ([]dto.RoleDTO, error)                                                         // Rôles créés par l'utilisateur (délégation)
	GetForDelegationPage(userID uint, filialeID *uint) ([]dto.RoleDTO, error)                                    // Rôles créés par l'utilisateur + rôles utilisés par au moins un user de la filiale
}

// roleService implémente RoleService
type roleService struct {
	roleRepo       repositories.RoleRepository
	userRepo       repositories.UserRepository
	permissionRepo repositories.PermissionRepository
	filialeRepo    repositories.FilialeRepository
}

// NewRoleService crée une nouvelle instance de RoleService
func NewRoleService(
	roleRepo repositories.RoleRepository,
	userRepo repositories.UserRepository,
	permissionRepo repositories.PermissionRepository,
	filialeRepo repositories.FilialeRepository,
) RoleService {
	return &roleService{
		roleRepo:       roleRepo,
		userRepo:       userRepo,
		permissionRepo: permissionRepo,
		filialeRepo:    filialeRepo,
	}
}

// Create crée un nouveau rôle
func (s *roleService) Create(req dto.CreateRoleRequest, createdByID uint) (*dto.RoleDTO, error) {
	// Empêcher la création d'un rôle avec le nom "ADMIN" (rôle système réservé)
	if req.Name == "ADMIN" {
		return nil, errors.New("impossible de créer un rôle avec le nom ADMIN (rôle système réservé)")
	}

	// Récupérer l'utilisateur créateur pour obtenir sa filiale et ses permissions
	creator, err := s.userRepo.FindByID(createdByID)
	if err != nil {
		return nil, errors.New("utilisateur créateur introuvable")
	}

	// Filiale optionnelle : si fournie, le nom du rôle sera "code filiale" + "-" + nom
	filialeID := req.FilialeID
	finalName := strings.TrimSpace(req.Name)
	if filialeID != nil {
		filiale, err := s.filialeRepo.FindByID(*filialeID)
		if err != nil || filiale == nil {
			return nil, errors.New("filiale introuvable")
		}
		code := strings.TrimSpace(filiale.Code)
		if code != "" {
			finalName = code + "-" + finalName
		}
	}

	// Vérifier que le nom final n'existe pas déjà
	existingRole, _ := s.roleRepo.FindByName(finalName)
	if existingRole != nil {
		return nil, errors.New("un rôle avec ce nom existe déjà")
	}

	// Si des permissions sont fournies, valider qu'elles sont un sous-ensemble des permissions du créateur
	if len(req.Permissions) > 0 {
		creatorPermissions, err := s.roleRepo.GetPermissionsByRoleID(creator.RoleID)
		if err != nil {
			return nil, errors.New("erreur lors de la récupération des permissions du créateur")
		}

		// Créer un map pour vérification rapide
		creatorPermMap := make(map[string]bool)
		for _, perm := range creatorPermissions {
			creatorPermMap[perm] = true
		}

		// Vérifier que toutes les permissions demandées appartiennent au créateur
		for _, perm := range req.Permissions {
			if !creatorPermMap[perm] {
				return nil, errors.New("vous ne pouvez assigner que les permissions que vous possédez vous-même")
			}
		}
	}

	role := &models.Role{
		Name:        finalName,
		Description: req.Description,
		IsSystem:    false, // Les rôles créés via API ne sont pas des rôles système
		CreatedByID: &createdByID,
		FilialeID:   filialeID,
	}

	if err := s.roleRepo.Create(role); err != nil {
		return nil, errors.New("erreur lors de la création du rôle")
	}

	// Assigner les permissions si fournies
	if len(req.Permissions) > 0 {
		if err := s.roleRepo.UpdateRolePermissions(role.ID, req.Permissions); err != nil {
			// Si l'assignation des permissions échoue, supprimer le rôle créé
			s.roleRepo.Delete(role.ID)
			return nil, errors.New("erreur lors de l'assignation des permissions")
		}
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

// GetAll récupère tous les rôles "réels" (exclut les rôles délégués, réservés à la page Délégation).
func (s *roleService) GetAll() ([]dto.RoleDTO, error) {
	roles, err := s.roleRepo.FindAll()
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des rôles")
	}

	var roleDTOs []dto.RoleDTO
	for i := range roles {
		if roles[i].CreatedByID == nil {
			roleDTOs = append(roleDTOs, s.roleToDTO(&roles[i]))
		}
	}
	return roleDTOs, nil
}

// GetAllForAssignment retourne les rôles "réels" visibles selon le périmètre (exclut les rôles délégués).
// viewMode: "all" = tous les rôles, "department" = rôles du département (departmentID requis), "filiale" = rôles globaux + filiale (filialeID optionnel).
func (s *roleService) GetAllForAssignment(filialeID *uint, departmentID *uint, viewMode string) ([]dto.RoleDTO, error) {
	switch viewMode {
	case "all":
		return s.GetAll()
	case "department":
		if departmentID == nil {
			return []dto.RoleDTO{}, nil
		}
		roles, err := s.roleRepo.FindByDepartmentID(*departmentID)
		if err != nil {
			return nil, errors.New("erreur lors de la récupération des rôles du département")
		}
		var roleDTOs []dto.RoleDTO
		for i := range roles {
			if roles[i].CreatedByID == nil {
				roleDTOs = append(roleDTOs, s.roleToDTO(&roles[i]))
			}
		}
		return roleDTOs, nil
	case "filiale":
		fallthrough
	default:
		roles, err := s.roleRepo.FindByFilialeOrGlobal(filialeID)
		if err != nil {
			return nil, errors.New("erreur lors de la récupération des rôles")
		}
		var roleDTOs []dto.RoleDTO
		for i := range roles {
			if roles[i].CreatedByID == nil {
				roleDTOs = append(roleDTOs, s.roleToDTO(&roles[i]))
			}
		}
		return roleDTOs, nil
	}
}

// GetMyDelegations retourne les rôles créés par l'utilisateur courant
func (s *roleService) GetMyDelegations(userID uint) ([]dto.RoleDTO, error) {
	roles, err := s.roleRepo.FindByCreatedByID(userID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des rôles délégués")
	}

	var roleDTOs []dto.RoleDTO
	for _, role := range roles {
		roleDTOs = append(roleDTOs, s.roleToDTO(&role))
	}

	return roleDTOs, nil
}

// GetForDelegationPage retourne les rôles à afficher sur la page "Délégation des rôles" :
// rôles créés par l'utilisateur + rôles utilisés par au moins un utilisateur de sa filiale.
func (s *roleService) GetForDelegationPage(userID uint, filialeID *uint) ([]dto.RoleDTO, error) {
	if filialeID == nil {
		return s.GetMyDelegations(userID)
	}
	roles, err := s.roleRepo.FindByCreatedByOrUsedInFiliale(userID, *filialeID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des rôles")
	}
	var roleDTOs []dto.RoleDTO
	for i := range roles {
		roleDTOs = append(roleDTOs, s.roleToDTO(&roles[i]))
	}
	return roleDTOs, nil
}

// Update met à jour un rôle
func (s *roleService) Update(id uint, req dto.UpdateRoleRequest, updatedByID uint, canManageAllRoles bool) (*dto.RoleDTO, error) {
	role, err := s.roleRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("rôle introuvable")
	}

	// Vérifier que ce n'est pas un rôle système (ne peut pas être modifié)
	if role.IsSystem {
		return nil, errors.New("impossible de modifier un rôle système")
	}

	// Empêcher la modification du nom en "ADMIN"
	if req.Name == "ADMIN" {
		return nil, errors.New("impossible d'utiliser le nom ADMIN (rôle système réservé)")
	}

	// Vérifier que l'utilisateur qui modifie existe
	if _, err := s.userRepo.FindByID(updatedByID); err != nil {
		return nil, errors.New("utilisateur modificateur introuvable")
	}

	// Si le rôle a un créateur, seul le créateur peut le modifier (sauf si l'utilisateur a roles.manage)
	if role.CreatedByID != nil {
		if !canManageAllRoles && *role.CreatedByID != updatedByID {
			return nil, errors.New("vous ne pouvez modifier que les rôles que vous avez créés")
		}
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

// Delete supprime un rôle (seul le créateur ou un utilisateur avec roles.manage peut supprimer)
func (s *roleService) Delete(id uint, deletedByID uint, canManageAllRoles bool) error {
	role, err := s.roleRepo.FindByID(id)
	if err != nil {
		return errors.New("rôle introuvable")
	}

	// Vérifier que ce n'est pas un rôle système
	if role.IsSystem {
		return errors.New("impossible de supprimer un rôle système")
	}

	// Si le rôle a un créateur, seul le créateur peut le supprimer (sauf si l'utilisateur a roles.manage)
	if role.CreatedByID != nil {
		if !canManageAllRoles && *role.CreatedByID != deletedByID {
			return errors.New("vous ne pouvez supprimer que les rôles que vous avez créés")
		}
	}

	// Vérifier qu'aucun utilisateur n'utilise ce rôle
	// nil scope car c'est une vérification interne (doit voir tous les utilisateurs)
	users, err := s.userRepo.FindByRole(nil, id)
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
		CreatedByID: role.CreatedByID,
		FilialeID:   role.FilialeID,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}
}

// GetRolePermissions récupère les permissions d'un rôle
func (s *roleService) GetRolePermissions(roleID uint) ([]string, error) {
	// Vérifier que le rôle existe
	_, err := s.roleRepo.FindByID(roleID)
	if err != nil {
		return nil, errors.New("rôle introuvable")
	}

	// Récupérer les permissions
	permissions, err := s.roleRepo.GetPermissionsByRoleID(roleID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des permissions")
	}

	return permissions, nil
}

// UpdateRolePermissions met à jour les permissions d'un rôle
func (s *roleService) UpdateRolePermissions(roleID uint, permissionCodes []string, updatedByID uint, canManageAllRoles bool) error {
	// Vérifier que le rôle existe
	role, err := s.roleRepo.FindByID(roleID)
	if err != nil {
		return errors.New("rôle introuvable")
	}

	// Vérifier que ce n'est pas un rôle système (ne peut pas être modifié)
	if role.IsSystem {
		return errors.New("impossible de modifier les permissions d'un rôle système")
	}

	// Vérifier que l'utilisateur qui modifie est le créateur du rôle (sauf si l'utilisateur a roles.manage)
	updater, err := s.userRepo.FindByID(updatedByID)
	if err != nil {
		return errors.New("utilisateur modificateur introuvable")
	}

	// Si le rôle a un créateur, seul le créateur peut modifier ses permissions (sauf si l'utilisateur a roles.manage)
	if role.CreatedByID != nil {
		if !canManageAllRoles && *role.CreatedByID != updatedByID {
			return errors.New("vous ne pouvez modifier que les permissions des rôles que vous avez créés")
		}
	}

	// Si l'utilisateur n'a pas roles.manage et que le rôle a un créateur,
	// valider que les permissions assignées sont un sous-ensemble des permissions du modificateur
	if role.CreatedByID != nil && !canManageAllRoles {
		updaterPermissions, err := s.roleRepo.GetPermissionsByRoleID(updater.RoleID)
		if err != nil {
			return errors.New("erreur lors de la récupération des permissions du modificateur")
		}

		// Créer un map pour vérification rapide
		updaterPermMap := make(map[string]bool)
		for _, perm := range updaterPermissions {
			updaterPermMap[perm] = true
		}

		// Vérifier que toutes les permissions demandées appartiennent au modificateur
		for _, perm := range permissionCodes {
			if !updaterPermMap[perm] {
				return errors.New("vous ne pouvez assigner que les permissions que vous possédez vous-même")
			}
		}
	}

	// Mettre à jour les permissions
	if err := s.roleRepo.UpdateRolePermissions(roleID, permissionCodes); err != nil {
		return errors.New("erreur lors de la mise à jour des permissions")
	}

	return nil
}

// GetAssignablePermissions récupère les permissions que l'utilisateur peut déléguer
func (s *roleService) GetAssignablePermissions(userID uint) ([]string, error) {
	// Récupérer l'utilisateur
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("utilisateur introuvable")
	}

	// Récupérer les permissions du rôle de l'utilisateur
	permissions, err := s.roleRepo.GetPermissionsByRoleID(user.RoleID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des permissions")
	}

	// Toutes les permissions sont déléguables : l'admin système est celui qui attribue les permissions
	// aux rôles dans la liste des rôles. S'il donne délibérément roles.manage à un rôle, c'est voulu.
	// Aucune permission n'est exclue de la liste des permissions assignables.
	assignablePermissions := make([]string, len(permissions))
	copy(assignablePermissions, permissions)

	return assignablePermissions, nil
}
