package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/mcicare/itsm-backend/config"
	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/models"
	"github.com/mcicare/itsm-backend/internal/repositories"
	"github.com/mcicare/itsm-backend/internal/utils"
)

// AuthService interface pour les opérations d'authentification
type AuthService interface {
	Register(req dto.RegisterRequest) (*dto.RegisterResponse, error)
	Login(req dto.LoginRequest) (*dto.LoginResponse, error)
	RefreshToken(refreshToken string) (string, error)
	Logout(userID uint, tokenHash string) error
}

// authService implémente AuthService
type authService struct {
	userRepo    repositories.UserRepository
	sessionRepo repositories.UserSessionRepository
	roleRepo    repositories.RoleRepository
}

// NewAuthService crée une nouvelle instance de AuthService
func NewAuthService(userRepo repositories.UserRepository, sessionRepo repositories.UserSessionRepository, roleRepo repositories.RoleRepository) AuthService {
	return &authService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		roleRepo:    roleRepo,
	}
}

// Register crée un nouveau compte utilisateur et connecte automatiquement l'utilisateur
func (s *authService) Register(req dto.RegisterRequest) (*dto.RegisterResponse, error) {
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

	// Trouver le rôle par défaut pour les nouveaux utilisateurs
	// Priorité: USER > CLIENT > premier rôle non-système > premier rôle disponible
	var defaultRole *models.Role
	var err error

	// 1. Essayer de trouver le rôle USER (rôle par défaut pour les nouveaux utilisateurs)
	defaultRole, err = s.roleRepo.FindByName("USER")
	if err != nil {
		// 2. Si "USER" n'existe pas, essayer "CLIENT"
		defaultRole, err = s.roleRepo.FindByName("CLIENT")
		if err != nil {
			// 3. Récupérer tous les rôles et chercher un rôle non-système
			roles, err := s.roleRepo.FindAll()
			if err != nil || len(roles) == 0 {
				return nil, errors.New("aucun rôle disponible dans le système. Veuillez contacter l'administrateur")
			}
			
			// Chercher un rôle non-système (pour éviter d'assigner ADMIN par erreur)
			for i := range roles {
				if !roles[i].IsSystem {
					defaultRole = &roles[i]
					break
				}
			}
			
			// Si aucun rôle non-système n'existe, prendre le premier disponible
			// (mais loguer un avertissement)
			if defaultRole == nil {
				defaultRole = &roles[0]
				// Note: Ce cas ne devrait normalement pas arriver car USER est créé par défaut
			}
		}
	}

	// Hasher le mot de passe
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("erreur lors du hashage du mot de passe")
	}

	// Vérifier que la filiale existe si fournie
	if req.FilialeID == nil {
		return nil, errors.New("la filiale est obligatoire")
	}

	// Créer l'utilisateur (sans createdByID pour l'inscription publique)
	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: passwordHash,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		FilialeID:    req.FilialeID,
		RoleID:       defaultRole.ID,
		IsActive:     true, // Par défaut actif
		CreatedByID:  nil, // Pas de créateur pour l'inscription publique
	}

	if err := s.userRepo.Create(user); err != nil {
		// Vérifier si c'est une erreur de contrainte de clé étrangère (filiale invalide)
		if strings.Contains(err.Error(), "foreign key constraint") {
			if strings.Contains(err.Error(), "filiale") {
				return nil, errors.New("filiale invalide ou introuvable")
			}
			return nil, errors.New("données invalides")
		}
		return nil, fmt.Errorf("erreur lors de la création du compte: %v", err)
	}

	// Récupérer l'utilisateur créé avec ses relations
	createdUser, err := s.userRepo.FindByID(user.ID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération de l'utilisateur créé")
	}

	// Générer le token JWT
	token, err := utils.GenerateToken(createdUser.ID, createdUser.Username, createdUser.Role.Name)
	if err != nil {
		return nil, errors.New("erreur lors de la génération du token")
	}

	// Générer le refresh token
	refreshToken, err := utils.GenerateRefreshToken(createdUser.ID)
	if err != nil {
		return nil, errors.New("erreur lors de la génération du refresh token")
	}

	// Créer une session utilisateur
	expiresAt := time.Now().Add(time.Duration(config.AppConfig.JWTExpirationHours) * time.Hour)
	session := &models.UserSession{
		UserID:           createdUser.ID,
		TokenHash:        utils.HashString(token),
		RefreshTokenHash: utils.HashString(refreshToken),
		ExpiresAt:        expiresAt,
		LastActivity:     time.Now(),
	}
	if err := s.sessionRepo.Create(session); err != nil {
		return nil, errors.New("erreur lors de la création de la session")
	}

	// Convertir l'utilisateur en DTO
	userDTO := s.userToDTO(createdUser)

	// Retourner la réponse
	return &dto.RegisterResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User:         userDTO,
	}, nil
}

// Login authentifie un utilisateur et retourne un token JWT
func (s *authService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	// Trouver l'utilisateur par email
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("email ou mot de passe incorrect")
	}

	// Vérifier si l'utilisateur est actif
	if !user.IsActive {
		return nil, errors.New("Ce compte est inactif. Vous n'avez plus accès au système. Contactez le service IT si vous pensez qu'il s'agit d'une erreur.")
	}

	// Vérifier le mot de passe
	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, errors.New("email ou mot de passe incorrect")
	}

	// Générer le token JWT
	token, err := utils.GenerateToken(user.ID, user.Username, user.Role.Name)
	if err != nil {
		return nil, errors.New("erreur lors de la génération du token")
	}

	// Générer le refresh token
	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, errors.New("erreur lors de la génération du refresh token")
	}

	// Créer une session utilisateur
	expiresAt := time.Now().Add(time.Duration(config.AppConfig.JWTExpirationHours) * time.Hour)
	session := &models.UserSession{
		UserID:           user.ID,
		TokenHash:        utils.HashString(token), // Hash du token pour sécurité
		RefreshTokenHash: utils.HashString(refreshToken),
		ExpiresAt:        expiresAt,
		LastActivity:     time.Now(),
	}
	if err := s.sessionRepo.Create(session); err != nil {
		return nil, errors.New("erreur lors de la création de la session")
	}

	// Mettre à jour la date de dernière connexion
	if err := s.userRepo.UpdateLastLogin(user.ID); err != nil {
		// Log l'erreur mais ne bloque pas la connexion
		// On pourrait utiliser un logger ici
	}

	// Convertir l'utilisateur en DTO
	userDTO := s.userToDTO(user)

	// Retourner la réponse
	return &dto.LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User:         userDTO,
	}, nil
}

// RefreshToken génère un nouveau token à partir d'un refresh token
func (s *authService) RefreshToken(refreshToken string) (string, error) {
	// Valider le refresh token
	claims, err := utils.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", errors.New("refresh token invalide ou expiré")
	}

	// Trouver l'utilisateur
	user, err := s.userRepo.FindByID(claims.UserID)
	if err != nil {
		return "", errors.New("utilisateur introuvable")
	}

	// Vérifier si l'utilisateur est actif
	if !user.IsActive {
		return "", errors.New("compte utilisateur désactivé")
	}

	// Vérifier que la session existe
	refreshTokenHash := utils.HashString(refreshToken)
	sessions, err := s.sessionRepo.FindByUserID(user.ID)
	if err != nil {
		return "", errors.New("session introuvable")
	}

	// Trouver la session correspondante
	var foundSession *models.UserSession
	for _, session := range sessions {
		if session.RefreshTokenHash == refreshTokenHash {
			foundSession = &session
			break
		}
	}

	if foundSession == nil {
		return "", errors.New("session introuvable")
	}

	// Vérifier que la session n'est pas expirée
	if foundSession.ExpiresAt.Before(time.Now()) {
		return "", errors.New("session expirée")
	}

	// Générer un nouveau token
	newToken, err := utils.GenerateToken(user.ID, user.Username, user.Role.Name)
	if err != nil {
		return "", errors.New("erreur lors de la génération du token")
	}

	// Mettre à jour le hash du token dans la session
	foundSession.TokenHash = utils.HashString(newToken)
	foundSession.LastActivity = time.Now()
	if err := s.sessionRepo.Update(foundSession); err != nil {
		return "", errors.New("erreur lors de la mise à jour de la session")
	}

	return newToken, nil
}

// Logout déconnecte un utilisateur en supprimant sa session
func (s *authService) Logout(userID uint, tokenHash string) error {
	// Trouver la session par token hash
	session, err := s.sessionRepo.FindByTokenHash(tokenHash)
	if err != nil {
		return errors.New("session introuvable")
	}

	// Vérifier que la session appartient à l'utilisateur
	if session.UserID != userID {
		return errors.New("session non autorisée")
	}

	// Supprimer la session
	return s.sessionRepo.Delete(session.ID)
}

// userToDTO convertit un modèle User en DTO UserDTO
func (s *authService) userToDTO(user *models.User) dto.UserDTO {
	userDTO := dto.UserDTO{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		DepartmentID: user.DepartmentID,
		FilialeID:    user.FilialeID,
		Avatar:       user.Avatar,
		Role:         user.Role.Name,                           // Nom du rôle brut (ex: "DSI")
		Permissions:  s.getPermissionsForRole(user.Role.Name),   // Permissions dérivées du rôle
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

	// Inclure le département si présent
	if user.Department != nil {
		userDTO.Department = &dto.DepartmentDTO{
			ID:          user.Department.ID,
			Name:        user.Department.Name,
			Code:        user.Department.Code,
			Description: user.Department.Description,
			OfficeID:    user.Department.OfficeID,
			IsActive:    user.Department.IsActive,
			CreatedAt:   user.Department.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   user.Department.UpdatedAt.Format(time.RFC3339),
		}
		// Inclure le siège si présent
		if user.Department.Office != nil {
			userDTO.Department.Office = &dto.OfficeDTO{
				ID:        user.Department.Office.ID,
				Name:      user.Department.Office.Name,
				Country:   user.Department.Office.Country,
				City:      user.Department.Office.City,
				Commune:   user.Department.Office.Commune,
				Address:   user.Department.Office.Address,
				Longitude: user.Department.Office.Longitude,
				Latitude:  user.Department.Office.Latitude,
				IsActive:  user.Department.Office.IsActive,
				CreatedAt: user.Department.Office.CreatedAt.Format(time.RFC3339),
				UpdatedAt: user.Department.Office.UpdatedAt.Format(time.RFC3339),
			}
		}
	}

	return userDTO
}

// getPermissionsForRole retourne la liste des permissions associées à un rôle donné.
// Les permissions sont récupérées depuis la base de données via la table role_permissions
// Cette fonction doit être identique à celle dans scope.go pour la cohérence
func (s *authService) getPermissionsForRole(roleName string) []string {
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
