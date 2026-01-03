package repositories

import (
	"time"

	"github.com/mcicare/itsm-backend/database"
	"github.com/mcicare/itsm-backend/internal/models"
)

// UserSessionRepository interface pour les opérations sur les sessions utilisateurs
type UserSessionRepository interface {
	Create(session *models.UserSession) error
	FindByID(id uint) (*models.UserSession, error)
	FindByTokenHash(tokenHash string) (*models.UserSession, error)
	FindByUserID(userID uint) ([]models.UserSession, error)
	FindActiveByUserID(userID uint) ([]models.UserSession, error)
	FindExpired() ([]models.UserSession, error)
	Update(session *models.UserSession) error
	Delete(id uint) error
	DeleteExpired() error
	DeleteByUserID(userID uint) error
	UpdateLastActivity(id uint) error
}

// userSessionRepository implémente UserSessionRepository
type userSessionRepository struct{}

// NewUserSessionRepository crée une nouvelle instance de UserSessionRepository
func NewUserSessionRepository() UserSessionRepository {
	return &userSessionRepository{}
}

// Create crée une nouvelle session
func (r *userSessionRepository) Create(session *models.UserSession) error {
	return database.DB.Create(session).Error
}

// FindByID trouve une session par son ID
func (r *userSessionRepository) FindByID(id uint) (*models.UserSession, error) {
	var session models.UserSession
	err := database.DB.Preload("User").Preload("User.Role").First(&session, id).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// FindByTokenHash trouve une session par le hash du token
func (r *userSessionRepository) FindByTokenHash(tokenHash string) (*models.UserSession, error) {
	var session models.UserSession
	err := database.DB.Preload("User").Preload("User.Role").Where("token_hash = ?", tokenHash).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// FindByUserID récupère toutes les sessions d'un utilisateur
func (r *userSessionRepository) FindByUserID(userID uint) ([]models.UserSession, error) {
	var sessions []models.UserSession
	err := database.DB.Preload("User").Where("user_id = ?", userID).Order("created_at DESC").Find(&sessions).Error
	return sessions, err
}

// FindActiveByUserID récupère les sessions actives d'un utilisateur
func (r *userSessionRepository) FindActiveByUserID(userID uint) ([]models.UserSession, error) {
	var sessions []models.UserSession
	now := time.Now()
	err := database.DB.Preload("User").Where("user_id = ? AND expires_at > ?", userID, now).Order("created_at DESC").Find(&sessions).Error
	return sessions, err
}

// FindExpired récupère toutes les sessions expirées
func (r *userSessionRepository) FindExpired() ([]models.UserSession, error) {
	var sessions []models.UserSession
	now := time.Now()
	err := database.DB.Preload("User").Where("expires_at <= ?", now).Find(&sessions).Error
	return sessions, err
}

// Update met à jour une session
func (r *userSessionRepository) Update(session *models.UserSession) error {
	return database.DB.Save(session).Error
}

// Delete supprime une session
func (r *userSessionRepository) Delete(id uint) error {
	return database.DB.Delete(&models.UserSession{}, id).Error
}

// DeleteExpired supprime toutes les sessions expirées
func (r *userSessionRepository) DeleteExpired() error {
	now := time.Now()
	return database.DB.Where("expires_at <= ?", now).Delete(&models.UserSession{}).Error
}

// DeleteByUserID supprime toutes les sessions d'un utilisateur
func (r *userSessionRepository) DeleteByUserID(userID uint) error {
	return database.DB.Where("user_id = ?", userID).Delete(&models.UserSession{}).Error
}

// UpdateLastActivity met à jour la dernière activité d'une session
func (r *userSessionRepository) UpdateLastActivity(id uint) error {
	now := time.Now()
	return database.DB.Model(&models.UserSession{}).Where("id = ?", id).Update("last_activity", now).Error
}

