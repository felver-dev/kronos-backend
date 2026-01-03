package services

import (
	"errors"
	"time"

	"github.com/mcicare/itsm-backend/config"
	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/models"
	"github.com/mcicare/itsm-backend/internal/repositories"
	"github.com/mcicare/itsm-backend/internal/utils"
)

// AuthService interface pour les opérations d'authentification
type AuthService interface {
	Login(req dto.LoginRequest) (*dto.LoginResponse, error)
	RefreshToken(refreshToken string) (string, error)
	Logout(userID uint, tokenHash string) error
}

// authService implémente AuthService
type authService struct {
	userRepo    repositories.UserRepository
	sessionRepo repositories.UserSessionRepository
}

// NewAuthService crée une nouvelle instance de AuthService
func NewAuthService(userRepo repositories.UserRepository, sessionRepo repositories.UserSessionRepository) AuthService {
	return &authService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

// Login authentifie un utilisateur et retourne un token JWT
func (s *authService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	// Trouver l'utilisateur par username
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, errors.New("nom d'utilisateur ou mot de passe incorrect")
	}

	// Vérifier si l'utilisateur est actif
	if !user.IsActive {
		return nil, errors.New("compte utilisateur désactivé")
	}

	// Vérifier le mot de passe
	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, errors.New("nom d'utilisateur ou mot de passe incorrect")
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
	return dto.UserDTO{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Avatar:    user.Avatar,
		Role:      user.Role.Name, // Le DTO utilise Role (string) au lieu de RoleID et RoleName
		IsActive:  user.IsActive,
		LastLogin: user.LastLogin,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
