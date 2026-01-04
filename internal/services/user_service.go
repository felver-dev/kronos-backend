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
	GetAll() ([]dto.UserDTO, error)
	GetByRole(roleID uint) ([]dto.UserDTO, error)
	Update(id uint, req dto.UpdateUserRequest, updatedByID uint) (*dto.UserDTO, error)
	Delete(id uint) error
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
	userRepo repositories.UserRepository
	roleRepo repositories.RoleRepository
}

// NewUserService crée une nouvelle instance de UserService
func NewUserService(userRepo repositories.UserRepository, roleRepo repositories.RoleRepository) UserService {
	return &userService{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

// Create crée un nouvel utilisateur
func (s *userService) Create(req dto.CreateUserRequest, createdByID uint) (*dto.UserDTO, error) {
	// Vérifier que le rôle existe
	_, err := s.roleRepo.FindByID(req.RoleID)
	if err != nil {
		return nil, errors.New("rôle introuvable")
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
func (s *userService) GetAll() ([]dto.UserDTO, error) {
	users, err := s.userRepo.FindAll()
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des utilisateurs")
	}

	userDTOs := make([]dto.UserDTO, len(users))
	for i, user := range users {
		userDTOs[i] = s.userToDTO(&user)
	}

	return userDTOs, nil
}

// GetByRole récupère tous les utilisateurs d'un rôle
func (s *userService) GetByRole(roleID uint) ([]dto.UserDTO, error) {
	users, err := s.userRepo.FindByRole(roleID)
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

	// Vérifier que le nouveau rôle existe si fourni
	if req.RoleID != 0 {
		_, err := s.roleRepo.FindByID(req.RoleID)
		if err != nil {
			return nil, errors.New("rôle introuvable")
		}
		user.RoleID = req.RoleID
	}

	// Mettre à jour les champs fournis
	if req.Email != "" {
		// Vérifier que l'email n'est pas déjà utilisé par un autre utilisateur
		existingUser, _ := s.userRepo.FindByEmail(req.Email)
		if existingUser != nil && existingUser.ID != id {
			return nil, errors.New("cet email est déjà utilisé")
		}
		user.Email = req.Email
	}

	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}

	if req.LastName != "" {
		user.LastName = req.LastName
	}

	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	user.UpdatedByID = &updatedByID

	// Sauvegarder
	if err := s.userRepo.Update(user); err != nil {
		return nil, errors.New("erreur lors de la mise à jour de l'utilisateur")
	}

	// Récupérer l'utilisateur mis à jour avec ses relations
	updatedUser, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération de l'utilisateur mis à jour")
	}

	userDTO := s.userToDTO(updatedUser)
	return &userDTO, nil
}

// Delete supprime un utilisateur (soft delete)
func (s *userService) Delete(id uint) error {
	// Vérifier que l'utilisateur existe
	_, err := s.userRepo.FindByID(id)
	if err != nil {
		return errors.New("utilisateur introuvable")
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
	return dto.UserDTO{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Avatar:    user.Avatar,
		Role:      user.Role.Name,
		IsActive:  user.IsActive,
		LastLogin: user.LastLogin,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
