package services

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mcicare/itsm-backend/config"
	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/models"
	"github.com/mcicare/itsm-backend/internal/repositories"
	"github.com/mcicare/itsm-backend/internal/utils"
)

// UserService interface pour les opérations sur les utilisateurs
type UserService interface {
	Create(req dto.CreateUserRequest, createdByID uint) (*dto.UserDTO, error)
	GetByID(id uint) (*dto.UserDTO, error)
	GetAll(scope interface{}) ([]dto.UserDTO, error)       // scope peut être *scope.QueryScope ou nil
	GetAllActive(scope interface{}) ([]dto.UserDTO, error) // scope peut être *scope.QueryScope ou nil
	GetByRole(scope interface{}, roleID uint) ([]dto.UserDTO, error)
	Update(id uint, req dto.UpdateUserRequest, updatedByID uint) (*dto.UserDTO, error)
	Delete(id uint, deletedByID uint) error
	ChangePassword(userID uint, oldPassword, newPassword string) error
	ResetPassword(userID uint, newPassword string) error
	Activate(id uint) error
	Deactivate(id uint) error
	GetPermissions(userID uint) (*dto.UserPermissionsDTO, error)
	UpdatePermissions(userID uint, req dto.UpdateUserPermissionsRequest, updatedByID uint) (*dto.UserPermissionsDTO, error)
	UploadAvatar(userID uint, filePath string, updatedByID uint) (*dto.UserDTO, error)
	GetAvatarPath(userID uint) (string, error)
	GetAvatarThumbnailPath(userID uint) (string, error)
	DeleteAvatar(userID uint, updatedByID uint) (*dto.UserDTO, error)
}

// userService implémente UserService
type userService struct {
	userRepo       repositories.UserRepository
	roleRepo       repositories.RoleRepository
	departmentRepo repositories.DepartmentRepository
	ticketRepo     repositories.TicketRepository
}

// NewUserService crée une nouvelle instance de UserService
func NewUserService(userRepo repositories.UserRepository, roleRepo repositories.RoleRepository, departmentRepo repositories.DepartmentRepository, ticketRepo repositories.TicketRepository) UserService {
	return &userService{
		userRepo:       userRepo,
		roleRepo:       roleRepo,
		departmentRepo: departmentRepo,
		ticketRepo:     ticketRepo,
	}
}

// Create crée un nouvel utilisateur
func (s *userService) Create(req dto.CreateUserRequest, createdByID uint) (*dto.UserDTO, error) {
	// Récupérer l'utilisateur créateur pour vérifier ses permissions
	creator, err := s.userRepo.FindByID(createdByID)
	if err != nil {
		return nil, errors.New("utilisateur créateur introuvable")
	}

	// Vérifier que le rôle existe (l'assignation est contrôlée par les permissions users.create / users.update)
	if req.RoleID == 0 {
		return nil, errors.New("un rôle doit être sélectionné")
	}
	_, err = s.roleRepo.FindByID(req.RoleID)
	if err != nil {
		return nil, errors.New("rôle introuvable")
	}

	// Vérifier que le département existe si fourni
	// Ignorer si DepartmentID est nil ou pointe vers 0 (valeur invalide)
	if req.DepartmentID != nil && *req.DepartmentID != 0 {
		dept, err := s.departmentRepo.FindByID(*req.DepartmentID)
		if err != nil {
			return nil, errors.New("département introuvable")
		}

		// Déterminer la filiale cible pour la validation
		// Si req.FilialeID est défini, utiliser celui-ci (le handler l'a déjà validé)
		// Sinon, utiliser la filiale du créateur
		var targetFilialeID *uint
		if req.FilialeID != nil {
			targetFilialeID = req.FilialeID
		} else if creator.FilialeID != nil {
			targetFilialeID = creator.FilialeID
		}

		// Vérifier que le département appartient à la filiale cible
		if targetFilialeID != nil && dept.FilialeID != nil {
			if *dept.FilialeID != *targetFilialeID {
				if req.FilialeID != nil && *req.FilialeID != *creator.FilialeID {
					// Le créateur a la permission users.create_any_filiale mais le département n'appartient pas à la filiale sélectionnée
					return nil, errors.New("le département n'appartient pas à la filiale sélectionnée")
				} else {
					// Le département n'appartient pas à la filiale de l'utilisateur créateur
					return nil, errors.New("le département n'appartient pas à votre filiale")
				}
			}
		} else if dept.FilialeID == nil {
			// Le département n'a pas de filiale assignée, ce qui ne devrait pas arriver normalement
			return nil, errors.New("le département n'a pas de filiale assignée")
		}
	}

	// Vérifier que la filiale existe si fournie
	if req.FilialeID != nil {
		// Note: On suppose qu'un repository filiale existe, sinon on peut ignorer cette vérification
		// ou utiliser un service filiale si disponible
	}

	// Vérifier que l'username n'existe pas déjà
	existingUser, _ := s.userRepo.FindByUsername(req.Username)
	if existingUser != nil {
		return nil, errors.New("ce nom d'utilisateur est déjà utilisé")
	}

	// Vérifier que l'email n'existe pas déjà
	existingUser, _ = s.userRepo.FindByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("cet email est déjà utilisé")
	}

	// Hasher le mot de passe
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("erreur lors du hashage du mot de passe")
	}

	// Créer l'utilisateur
	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: passwordHash,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Phone:        req.Phone,
		DepartmentID: req.DepartmentID,
		FilialeID:    req.FilialeID,
		RoleID:       req.RoleID,
		IsActive:     true, // Par défaut actif
		CreatedByID:  &createdByID,
	}

	if err := s.userRepo.Create(user); err != nil {
		// Retourner l'erreur réelle pour faciliter le débogage
		return nil, fmt.Errorf("erreur lors de la création de l'utilisateur: %w", err)
	}

	// Récupérer l'utilisateur créé avec ses relations
	createdUser, err := s.userRepo.FindByID(user.ID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération de l'utilisateur créé")
	}

	// Convertir en DTO
	userDTO := s.userToDTO(createdUser)
	return &userDTO, nil
}

// GetByID récupère un utilisateur par son ID
func (s *userService) GetByID(id uint) (*dto.UserDTO, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("utilisateur introuvable")
	}

	userDTO := s.userToDTO(user)
	return &userDTO, nil
}

// GetAll récupère tous les utilisateurs
// Le scope est utilisé pour filtrer automatiquement selon les permissions de l'utilisateur
func (s *userService) GetAll(scopeParam interface{}) ([]dto.UserDTO, error) {
	users, err := s.userRepo.FindAll(scopeParam)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des utilisateurs")
	}

	userDTOs := make([]dto.UserDTO, len(users))
	for i, user := range users {
		userDTOs[i] = s.userToDTO(&user)
	}

	return userDTOs, nil
}

// GetAllActive récupère tous les utilisateurs actifs
// Le scope est utilisé pour filtrer automatiquement selon les permissions de l'utilisateur
// Si scope est nil, retourne tous les utilisateurs actifs sans filtre de filiale
func (s *userService) GetAllActive(scopeParam interface{}) ([]dto.UserDTO, error) {
	users, err := s.userRepo.FindActive(scopeParam)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des utilisateurs actifs")
	}

	userDTOs := make([]dto.UserDTO, len(users))
	for i, user := range users {
		userDTOs[i] = s.userToDTO(&user)
	}

	return userDTOs, nil
}

// GetByRole récupère tous les utilisateurs d'un rôle
// Le scope est utilisé pour filtrer automatiquement selon les permissions de l'utilisateur
func (s *userService) GetByRole(scopeParam interface{}, roleID uint) ([]dto.UserDTO, error) {
	users, err := s.userRepo.FindByRole(scopeParam, roleID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des utilisateurs")
	}

	userDTOs := make([]dto.UserDTO, len(users))
	for i, user := range users {
		userDTOs[i] = s.userToDTO(&user)
	}

	return userDTOs, nil
}

// Update met à jour un utilisateur
func (s *userService) Update(id uint, req dto.UpdateUserRequest, updatedByID uint) (*dto.UserDTO, error) {
	// Récupérer l'utilisateur existant
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("utilisateur introuvable")
	}

	// Log pour déboguer
	fmt.Printf("Update user %d - RoleID reçu: %d, Rôle actuel: %d\n", id, req.RoleID, user.RoleID)

	// Récupérer l'utilisateur qui effectue la mise à jour pour vérifier ses permissions
	updater, err := s.userRepo.FindByID(updatedByID)
	if err != nil {
		return nil, errors.New("utilisateur modificateur introuvable")
	}

	// Vérifier que le nouveau rôle existe si fourni (l'assignation est contrôlée par les permissions users.update)
	if req.RoleID != 0 {
		_, err := s.roleRepo.FindByID(req.RoleID)
		if err != nil {
			return nil, errors.New("rôle introuvable")
		}
		fmt.Printf("Rôle trouvé, mise à jour du rôle de %d vers %d\n", user.RoleID, req.RoleID)
		user.RoleID = req.RoleID
	} else {
		fmt.Printf("RoleID est 0 ou non fourni, pas de mise à jour du rôle\n")
	}

	// Gérer la mise à jour du département
	// Si DepartmentID est fourni (non nil et non 0), vérifier qu'il existe et l'assigner
	// Si DepartmentID est nil, ne pas modifier le département existant
	// Si DepartmentID pointe vers 0, cela signifie qu'on veut retirer l'utilisateur du département
	if req.DepartmentID != nil {
		// Si DepartmentID pointe vers 0, retirer l'utilisateur du département
		if *req.DepartmentID == 0 {
			user.DepartmentID = nil
		} else {
			// DepartmentID pointe vers une valeur valide, vérifier qu'il existe
			dept, err := s.departmentRepo.FindByID(*req.DepartmentID)
			if err != nil {
				return nil, errors.New("département introuvable")
			}

			// Déterminer la filiale cible (celle modifiée, celle de l'utilisateur modifié, ou celle du modificateur)
			var targetFilialeID *uint
			if req.FilialeID != nil {
				targetFilialeID = req.FilialeID
			} else if user.FilialeID != nil {
				targetFilialeID = user.FilialeID
			} else if updater.FilialeID != nil {
				targetFilialeID = updater.FilialeID
			}

			// Vérifier que le département appartient à la filiale cible
			if targetFilialeID != nil && dept.FilialeID != nil {
				if *dept.FilialeID != *targetFilialeID {
					return nil, errors.New("le département n'appartient pas à la filiale de l'utilisateur")
				}
			}

			user.DepartmentID = req.DepartmentID
		}
	}
	// Si DepartmentID est nil dans la requête, ne pas modifier le département existant
	// Le frontend n'a pas envoyé ce champ, donc on garde la valeur actuelle

	// Gérer la mise à jour de la filiale
	// Si FilialeID est fourni (non nil), l'assigner
	if req.FilialeID != nil {
		user.FilialeID = req.FilialeID
	} else {
		// Si FilialeID est nil dans la requête, on vérifie si l'utilisateur a déjà une filiale
		// Si oui et que le champ a été explicitement envoyé comme null, on retire la filiale
		if user.FilialeID != nil {
			// Si l'utilisateur a une filiale et que req.FilialeID est nil,
			// on retire la filiale (le frontend a envoyé filiale_id: null)
			user.FilialeID = nil
		}
		// Si l'utilisateur n'a pas de filiale et que req.FilialeID est nil,
		// on ne fait rien (pas de changement)
	}

	// Mettre à jour les champs fournis
	// Username peut être modifié si fourni
	if req.Username != "" {
		// Vérifier que le username n'est pas déjà utilisé par un autre utilisateur
		existingUser, _ := s.userRepo.FindByUsername(req.Username)
		if existingUser != nil && existingUser.ID != id {
			return nil, errors.New("ce nom d'utilisateur est déjà utilisé")
		}
		user.Username = req.Username
	}

	// Email est toujours requis, donc on le met à jour si fourni
	if req.Email != "" {
		// Vérifier que l'email n'est pas déjà utilisé par un autre utilisateur
		existingUser, _ := s.userRepo.FindByEmail(req.Email)
		if existingUser != nil && existingUser.ID != id {
			return nil, errors.New("cet email est déjà utilisé")
		}
		user.Email = req.Email
	} else if req.Email == "" && user.Email == "" {
		// Si l'email n'est pas fourni et que l'utilisateur n'a pas d'email, c'est une erreur
		return nil, errors.New("l'email est requis")
	}

	// Vérifier si le nom ou prénom a changé pour mettre à jour les tickets
	oldFirstName := user.FirstName
	oldLastName := user.LastName
	oldRequesterName := ""
	if oldFirstName != "" || oldLastName != "" {
		oldRequesterName = fmt.Sprintf("%s %s", strings.TrimSpace(oldFirstName), strings.TrimSpace(oldLastName))
	} else {
		oldRequesterName = user.Username
	}

	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}

	if req.LastName != "" {
		user.LastName = req.LastName
	}

	if req.Phone != "" {
		user.Phone = req.Phone
	}

	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	user.UpdatedByID = &updatedByID

	// Sauvegarder
	if err := s.userRepo.Update(user); err != nil {
		// Retourner l'erreur réelle pour le débogage
		return nil, fmt.Errorf("erreur lors de la mise à jour de l'utilisateur: %w", err)
	}

	// Récupérer l'utilisateur mis à jour avec ses relations
	updatedUser, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération de l'utilisateur mis à jour")
	}

	// Log pour vérifier que le rôle a bien été mis à jour
	fmt.Printf("Après mise à jour - User ID: %d, RoleID dans la DB: %d, Role Name: %s\n",
		updatedUser.ID, updatedUser.RoleID, updatedUser.Role.Name)

	// Mettre à jour le nom du demandeur dans tous les tickets créés par cet utilisateur
	// si le nom ou prénom a changé
	newRequesterName := ""
	if updatedUser.FirstName != "" || updatedUser.LastName != "" {
		newRequesterName = fmt.Sprintf("%s %s", strings.TrimSpace(updatedUser.FirstName), strings.TrimSpace(updatedUser.LastName))
	} else {
		newRequesterName = updatedUser.Username
	}

	// Si le nom du demandeur a changé, mettre à jour tous les tickets
	if oldRequesterName != newRequesterName {
		// 1. Mettre à jour les tickets où requester_id correspond à cet utilisateur
		if err := s.ticketRepo.UpdateRequesterNameByRequesterID(id, newRequesterName); err != nil {
			// Log l'erreur mais ne bloque pas la mise à jour de l'utilisateur
			fmt.Printf("Erreur lors de la mise à jour du nom du demandeur dans les tickets (par requester_id): %v\n", err)
		}
		// 2. Mettre à jour les tickets créés par cet utilisateur
		if err := s.ticketRepo.UpdateRequesterNameByCreatedBy(id, newRequesterName); err != nil {
			// Log l'erreur mais ne bloque pas la mise à jour de l'utilisateur
			fmt.Printf("Erreur lors de la mise à jour du nom du demandeur dans les tickets (par created_by_id): %v\n", err)
		}
		// 3. Mettre à jour aussi les tickets où le requester_name correspond à l'ancien nom
		// (au cas où le ticket a été créé par quelqu'un d'autre pour cet utilisateur)
		if err := s.ticketRepo.UpdateRequesterNameByName(oldRequesterName, newRequesterName); err != nil {
			// Log l'erreur mais ne bloque pas la mise à jour de l'utilisateur
			fmt.Printf("Erreur lors de la mise à jour du nom du demandeur dans les tickets (par requester_name): %v\n", err)
		}
	}

	userDTO := s.userToDTO(updatedUser)
	return &userDTO, nil
}

// Delete supprime un utilisateur (soft delete)
func (s *userService) Delete(id uint, deletedByID uint) error {
	// Vérifier que l'utilisateur existe et récupérer ses informations
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return errors.New("utilisateur introuvable")
	}

	// Empêcher la suppression du compte admin par défaut (point d'entrée de l'application)
	if user.Username == "admin" && user.Email == "admin@kronos.com" {
		return errors.New("impossible de supprimer le compte administrateur par défaut")
	}

	// Vérifier si l'utilisateur à supprimer est un admin
	// Si c'est le cas, vérifier qu'il n'est pas le dernier admin
	if user.Role.Name == "ADMIN" {
		// Compter le nombre d'admins actifs (non supprimés)
		var adminCount int64
		adminRole, err := s.roleRepo.FindByName("ADMIN")
		if err != nil {
			return errors.New("erreur lors de la vérification du rôle admin")
		}

		err = s.userRepo.CountByRole(adminRole.ID, &adminCount)
		if err != nil {
			return errors.New("erreur lors du comptage des administrateurs")
		}

		// Si c'est le dernier admin, empêcher la suppression
		if adminCount <= 1 {
			return errors.New("impossible de supprimer le dernier administrateur du système")
		}
	}

	// Supprimer (soft delete)
	if err := s.userRepo.Delete(id); err != nil {
		return errors.New("erreur lors de la suppression de l'utilisateur")
	}

	return nil
}

// ChangePassword change le mot de passe d'un utilisateur
func (s *userService) ChangePassword(userID uint, oldPassword, newPassword string) error {
	// Récupérer l'utilisateur
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("utilisateur introuvable")
	}

	// Vérifier l'ancien mot de passe
	if !utils.CheckPasswordHash(oldPassword, user.PasswordHash) {
		return errors.New("ancien mot de passe incorrect")
	}

	// Hasher le nouveau mot de passe
	passwordHash, err := utils.HashPassword(newPassword)
	if err != nil {
		return errors.New("erreur lors du hashage du mot de passe")
	}

	// Mettre à jour
	user.PasswordHash = passwordHash
	if err := s.userRepo.Update(user); err != nil {
		return errors.New("erreur lors de la mise à jour du mot de passe")
	}

	return nil
}

// ResetPassword réinitialise le mot de passe d'un utilisateur (admin uniquement)
func (s *userService) ResetPassword(userID uint, newPassword string) error {
	// Récupérer l'utilisateur
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("utilisateur introuvable")
	}

	// Hasher le nouveau mot de passe
	passwordHash, err := utils.HashPassword(newPassword)
	if err != nil {
		return errors.New("erreur lors du hashage du mot de passe")
	}

	// Mettre à jour
	user.PasswordHash = passwordHash
	if err := s.userRepo.Update(user); err != nil {
		return errors.New("erreur lors de la réinitialisation du mot de passe")
	}

	return nil
}

// Activate active un utilisateur
func (s *userService) Activate(id uint) error {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return errors.New("utilisateur introuvable")
	}

	user.IsActive = true
	if err := s.userRepo.Update(user); err != nil {
		return errors.New("erreur lors de l'activation de l'utilisateur")
	}

	return nil
}

// Deactivate désactive un utilisateur
func (s *userService) Deactivate(id uint) error {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return errors.New("utilisateur introuvable")
	}

	user.IsActive = false
	if err := s.userRepo.Update(user); err != nil {
		return errors.New("erreur lors de la désactivation de l'utilisateur")
	}

	return nil
}

// GetPermissions récupère les permissions d'un utilisateur
func (s *userService) GetPermissions(userID uint) (*dto.UserPermissionsDTO, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("utilisateur introuvable")
	}

	userDTO := s.userToDTO(user)

	// Pour l'instant, on retourne les permissions du rôle
	// TODO: Implémenter la gestion des permissions personnalisées par utilisateur
	permissions := []string{} // Les permissions seront récupérées depuis le rôle

	permissionsDTO := &dto.UserPermissionsDTO{
		UserID:      userID,
		User:        &userDTO,
		Permissions: permissions,
	}

	return permissionsDTO, nil
}

// UpdatePermissions met à jour les permissions d'un utilisateur
func (s *userService) UpdatePermissions(userID uint, req dto.UpdateUserPermissionsRequest, updatedByID uint) (*dto.UserPermissionsDTO, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("utilisateur introuvable")
	}

	// TODO: Implémenter la sauvegarde des permissions personnalisées
	// Pour l'instant, on retourne simplement les permissions demandées
	userDTO := s.userToDTO(user)

	permissionsDTO := &dto.UserPermissionsDTO{
		UserID:      userID,
		User:        &userDTO,
		Permissions: req.Permissions,
	}

	return permissionsDTO, nil
}

// UploadAvatar upload un avatar pour un utilisateur
func (s *userService) UploadAvatar(userID uint, fileName string, updatedByID uint) (*dto.UserDTO, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("utilisateur introuvable")
	}

	// Supprimer l'ancien avatar s'il existe
	if user.Avatar != "" {
		oldPath := filepath.Join(config.AppConfig.AvatarDir, user.Avatar)
		if _, err := os.Stat(oldPath); err == nil {
			os.Remove(oldPath)
		}
		// Supprimer aussi la miniature
		thumbnailPath := strings.Replace(oldPath, filepath.Ext(oldPath), "_thumb"+filepath.Ext(oldPath), 1)
		if _, err := os.Stat(thumbnailPath); err == nil {
			os.Remove(thumbnailPath)
		}
	}

	// Mettre à jour l'avatar dans la base de données (on stocke juste le nom du fichier)
	user.Avatar = fileName
	user.UpdatedByID = &updatedByID

	if err := s.userRepo.Update(user); err != nil {
		return nil, errors.New("erreur lors de la mise à jour de l'avatar")
	}

	updatedUser, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération de l'utilisateur mis à jour")
	}

	userDTO := s.userToDTO(updatedUser)
	return &userDTO, nil
}

// GetAvatarPath récupère le chemin de l'avatar d'un utilisateur
func (s *userService) GetAvatarPath(userID uint) (string, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return "", errors.New("utilisateur introuvable")
	}

	if user.Avatar == "" {
		return "", errors.New("aucun avatar trouvé")
	}

	fullPath := filepath.Join(config.AppConfig.AvatarDir, user.Avatar)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", errors.New("fichier avatar introuvable")
	}

	return fullPath, nil
}

// GetAvatarThumbnailPath récupère le chemin de la miniature de l'avatar
func (s *userService) GetAvatarThumbnailPath(userID uint) (string, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return "", errors.New("utilisateur introuvable")
	}

	if user.Avatar == "" {
		return "", errors.New("aucun avatar trouvé")
	}

	// Générer le chemin de la miniature
	avatarPath := filepath.Join(config.AppConfig.AvatarDir, user.Avatar)
	ext := filepath.Ext(avatarPath)
	thumbnailPath := strings.Replace(avatarPath, ext, "_thumb"+ext, 1)

	if _, err := os.Stat(thumbnailPath); os.IsNotExist(err) {
		// Si la miniature n'existe pas, retourner l'avatar original
		if _, err := os.Stat(avatarPath); os.IsNotExist(err) {
			return "", errors.New("fichier avatar introuvable")
		}
		return avatarPath, nil
	}

	return thumbnailPath, nil
}

// DeleteAvatar supprime l'avatar d'un utilisateur
func (s *userService) DeleteAvatar(userID uint, updatedByID uint) (*dto.UserDTO, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("utilisateur introuvable")
	}

	if user.Avatar == "" {
		return nil, errors.New("aucun avatar à supprimer")
	}

	// Supprimer le fichier
	avatarPath := filepath.Join(config.AppConfig.AvatarDir, user.Avatar)
	if _, err := os.Stat(avatarPath); err == nil {
		os.Remove(avatarPath)
	}

	// Supprimer la miniature
	ext := filepath.Ext(avatarPath)
	thumbnailPath := strings.Replace(avatarPath, ext, "_thumb"+ext, 1)
	if _, err := os.Stat(thumbnailPath); err == nil {
		os.Remove(thumbnailPath)
	}

	// Mettre à jour dans la base de données
	user.Avatar = ""
	user.UpdatedByID = &updatedByID

	if err := s.userRepo.Update(user); err != nil {
		return nil, errors.New("erreur lors de la suppression de l'avatar")
	}

	updatedUser, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération de l'utilisateur mis à jour")
	}

	userDTO := s.userToDTO(updatedUser)
	return &userDTO, nil
}

// userToDTO convertit un modèle User en DTO UserDTO
func (s *userService) userToDTO(user *models.User) dto.UserDTO {
	userDTO := dto.UserDTO{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Phone:        user.Phone,
		DepartmentID: user.DepartmentID,
		FilialeID:    user.FilialeID,
		Avatar:       user.Avatar,
		Role:         user.Role.Name,
		Permissions:  s.getPermissionsForRole(user.Role.Name),
		IsActive:     user.IsActive,
		LastLogin:    user.LastLogin,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}

	// Inclure la filiale si présente
	if user.Filiale != nil {
		userDTO.Filiale = &dto.FilialeDTO{
			ID:          user.Filiale.ID,
			Code:        user.Filiale.Code,
			Name:        user.Filiale.Name,
			Country:     user.Filiale.Country,
			City:        user.Filiale.City,
			Address:     user.Filiale.Address,
			Phone:       user.Filiale.Phone,
			Email:       user.Filiale.Email,
			IsActive:    user.Filiale.IsActive,
			IsSoftwareProvider: user.Filiale.IsSoftwareProvider,
			CreatedAt:   user.Filiale.CreatedAt,
			UpdatedAt:   user.Filiale.UpdatedAt,
		}
	}

	// Inclure le département complet si présent
	if user.Department != nil {
		departmentDTO := dto.DepartmentDTO{
			ID:          user.Department.ID,
			Name:        user.Department.Name,
			Code:        user.Department.Code,
			Description: user.Department.Description,
			OfficeID:    user.Department.OfficeID,
			IsActive:    user.Department.IsActive,
			CreatedAt:   user.Department.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   user.Department.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
		// Inclure le siège si présent
		if user.Department.Office != nil {
			departmentDTO.Office = &dto.OfficeDTO{
				ID:        user.Department.Office.ID,
				Name:      user.Department.Office.Name,
				Country:   user.Department.Office.Country,
				City:      user.Department.Office.City,
				Commune:   user.Department.Office.Commune,
				Address:   user.Department.Office.Address,
				Longitude: user.Department.Office.Longitude,
				Latitude:  user.Department.Office.Latitude,
				IsActive:  user.Department.Office.IsActive,
				CreatedAt: user.Department.Office.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
				UpdatedAt: user.Department.Office.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
			}
		}
		userDTO.Department = &departmentDTO
	}

	return userDTO
}

// getPermissionsForRole retourne la liste des permissions associées à un rôle donné.
// Les permissions sont récupérées depuis la base de données via la table role_permissions
func (s *userService) getPermissionsForRole(roleName string) []string {
	// Récupérer le rôle par son nom
	role, err := s.roleRepo.FindByName(roleName)
	if err != nil {
		// Retourner des permissions minimales par défaut si le rôle n'existe pas
		return []string{"tickets.view_own"}
	}

	// Récupérer les permissions du rôle depuis la base de données
	permissions, err := s.roleRepo.GetPermissionsByRoleID(role.ID)
	if err != nil {
		// Retourner des permissions minimales par défaut en cas d'erreur
		return []string{"tickets.view_own"}
	}

	// Si aucune permission n'est assignée, retourner des permissions minimales
	if len(permissions) == 0 {
		return []string{"tickets.view_own"}
	}

	return permissions
}
