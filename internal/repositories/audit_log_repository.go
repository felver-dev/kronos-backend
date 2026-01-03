package repositories

import (
	"time"

	"github.com/mcicare/itsm-backend/database"
	"github.com/mcicare/itsm-backend/internal/models"
)

// AuditLogRepository interface pour les opérations sur les logs d'audit
type AuditLogRepository interface {
	Create(log *models.AuditLog) error
	FindByID(id uint) (*models.AuditLog, error)
	FindByUserID(userID uint) ([]models.AuditLog, error)
	FindByEntity(entityType string, entityID uint) ([]models.AuditLog, error)
	FindByAction(action string) ([]models.AuditLog, error)
	FindByDateRange(startDate, endDate time.Time) ([]models.AuditLog, error)
	FindRecent(limit int) ([]models.AuditLog, error)
	Delete(id uint) error
	DeleteOld(olderThan time.Time) error
}

// auditLogRepository implémente AuditLogRepository
type auditLogRepository struct{}

// NewAuditLogRepository crée une nouvelle instance de AuditLogRepository
func NewAuditLogRepository() AuditLogRepository {
	return &auditLogRepository{}
}

// Create crée un nouveau log d'audit
func (r *auditLogRepository) Create(log *models.AuditLog) error {
	return database.DB.Create(log).Error
}

// FindByID trouve un log d'audit par son ID
func (r *auditLogRepository) FindByID(id uint) (*models.AuditLog, error) {
	var auditLog models.AuditLog
	err := database.DB.Preload("User").Preload("User.Role").First(&auditLog, id).Error
	if err != nil {
		return nil, err
	}
	return &auditLog, nil
}

// FindByUserID récupère tous les logs d'audit d'un utilisateur
func (r *auditLogRepository) FindByUserID(userID uint) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	err := database.DB.Preload("User").Preload("User.Role").Where("user_id = ?", userID).Order("created_at DESC").Find(&logs).Error
	return logs, err
}

// FindByEntity récupère tous les logs d'audit d'une entité
func (r *auditLogRepository) FindByEntity(entityType string, entityID uint) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	err := database.DB.Preload("User").Preload("User.Role").Where("entity_type = ? AND entity_id = ?", entityType, entityID).Order("created_at DESC").Find(&logs).Error
	return logs, err
}

// FindByAction récupère tous les logs d'audit d'une action
func (r *auditLogRepository) FindByAction(action string) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	err := database.DB.Preload("User").Preload("User.Role").Where("action = ?", action).Order("created_at DESC").Find(&logs).Error
	return logs, err
}

// FindByDateRange récupère les logs d'audit dans une plage de dates
func (r *auditLogRepository) FindByDateRange(startDate, endDate time.Time) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	err := database.DB.Preload("User").Preload("User.Role").
		Where("created_at >= ? AND created_at <= ?", startDate, endDate).
		Order("created_at DESC").
		Find(&logs).Error
	return logs, err
}

// FindRecent récupère les logs d'audit les plus récents
func (r *auditLogRepository) FindRecent(limit int) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	err := database.DB.Preload("User").Preload("User.Role").Order("created_at DESC").Limit(limit).Find(&logs).Error
	return logs, err
}

// Delete supprime un log d'audit
func (r *auditLogRepository) Delete(id uint) error {
	return database.DB.Delete(&models.AuditLog{}, id).Error
}

// DeleteOld supprime les logs d'audit plus anciens qu'une date donnée (pour nettoyage)
func (r *auditLogRepository) DeleteOld(olderThan time.Time) error {
	return database.DB.Where("created_at < ?", olderThan).Delete(&models.AuditLog{}).Error
}
