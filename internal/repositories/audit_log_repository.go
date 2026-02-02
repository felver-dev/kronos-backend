package repositories

import (
	"time"

	"github.com/mcicare/itsm-backend/database"
	"github.com/mcicare/itsm-backend/internal/models"
	"github.com/mcicare/itsm-backend/internal/scope"
)

// AuditLogRepository interface pour les opérations sur les logs d'audit
type AuditLogRepository interface {
	Create(log *models.AuditLog) error
	FindByID(id uint) (*models.AuditLog, error)
	FindByUserID(scope interface{}, userID uint) ([]models.AuditLog, error)
	FindByEntity(scope interface{}, entityType string, entityID uint) ([]models.AuditLog, error)
	FindByAction(scope interface{}, action string) ([]models.AuditLog, error)
	FindByDateRange(scope interface{}, startDate, endDate time.Time) ([]models.AuditLog, error)
	FindRecent(scope interface{}, limit int) ([]models.AuditLog, error)
	// FindPaginated récupère les logs d'audit avec pagination et filtres, et renvoie aussi le total
	FindPaginated(scope interface{}, page, limit int, userID *uint, action, entityType string) ([]models.AuditLog, int64, error)
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
// Le scope est utilisé pour filtrer automatiquement selon les permissions de l'utilisateur
func (r *auditLogRepository) FindByUserID(scopeParam interface{}, userID uint) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	
	// Construire la requête de base
	query := database.DB.Model(&models.AuditLog{}).
		Preload("User").Preload("User.Role").
		Where("audit_logs.user_id = ?", userID)
	
	// Appliquer le scope si fourni
	if scopeParam != nil {
		if queryScope, ok := scopeParam.(*scope.QueryScope); ok {
			query = scope.ApplyAuditScope(query, queryScope)
		}
	}
	
	err := query.Order("audit_logs.created_at DESC").Find(&logs).Error
	return logs, err
}

// FindByEntity récupère tous les logs d'audit d'une entité
// Le scope est utilisé pour filtrer automatiquement selon les permissions de l'utilisateur
func (r *auditLogRepository) FindByEntity(scopeParam interface{}, entityType string, entityID uint) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	
	// Construire la requête de base
	query := database.DB.Model(&models.AuditLog{}).
		Preload("User").Preload("User.Role").
		Where("audit_logs.entity_type = ? AND audit_logs.entity_id = ?", entityType, entityID)
	
	// Appliquer le scope si fourni
	if scopeParam != nil {
		if queryScope, ok := scopeParam.(*scope.QueryScope); ok {
			query = scope.ApplyAuditScope(query, queryScope)
		}
	}
	
	err := query.Order("audit_logs.created_at DESC").Find(&logs).Error
	return logs, err
}

// FindByAction récupère tous les logs d'audit d'une action
// Le scope est utilisé pour filtrer automatiquement selon les permissions de l'utilisateur
func (r *auditLogRepository) FindByAction(scopeParam interface{}, action string) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	
	// Construire la requête de base
	query := database.DB.Model(&models.AuditLog{}).
		Preload("User").Preload("User.Role").
		Where("audit_logs.action = ?", action)
	
	// Appliquer le scope si fourni
	if scopeParam != nil {
		if queryScope, ok := scopeParam.(*scope.QueryScope); ok {
			query = scope.ApplyAuditScope(query, queryScope)
		}
	}
	
	err := query.Order("audit_logs.created_at DESC").Find(&logs).Error
	return logs, err
}

// FindByDateRange récupère les logs d'audit dans une plage de dates
// Le scope est utilisé pour filtrer automatiquement selon les permissions de l'utilisateur
func (r *auditLogRepository) FindByDateRange(scopeParam interface{}, startDate, endDate time.Time) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	
	// Construire la requête de base
	query := database.DB.Model(&models.AuditLog{}).
		Preload("User").Preload("User.Role").
		Where("audit_logs.created_at >= ? AND audit_logs.created_at <= ?", startDate, endDate)
	
	// Appliquer le scope si fourni
	if scopeParam != nil {
		if queryScope, ok := scopeParam.(*scope.QueryScope); ok {
			query = scope.ApplyAuditScope(query, queryScope)
		}
	}
	
	err := query.Order("audit_logs.created_at DESC").Find(&logs).Error
	return logs, err
}

// FindRecent récupère les logs d'audit les plus récents
// Le scope est utilisé pour filtrer automatiquement selon les permissions de l'utilisateur
func (r *auditLogRepository) FindRecent(scopeParam interface{}, limit int) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	
	// Construire la requête de base
	query := database.DB.Model(&models.AuditLog{}).
		Preload("User").Preload("User.Role")
	
	// Appliquer le scope si fourni
	if scopeParam != nil {
		if queryScope, ok := scopeParam.(*scope.QueryScope); ok {
			query = scope.ApplyAuditScope(query, queryScope)
		}
	}
	
	err := query.Order("audit_logs.created_at DESC").Limit(limit).Find(&logs).Error
	return logs, err
}

// FindPaginated récupère les logs d'audit avec pagination et filtres
// Le scope est utilisé pour filtrer automatiquement selon les permissions de l'utilisateur
func (r *auditLogRepository) FindPaginated(scopeParam interface{}, page, limit int, userID *uint, action, entityType string) ([]models.AuditLog, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 50
	}

	offset := (page - 1) * limit

	// Construire la requête de base avec les filtres
	query := database.DB.Model(&models.AuditLog{})

	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}
	if action != "" {
		query = query.Where("action = ?", action)
	}
	if entityType != "" {
		query = query.Where("entity_type = ?", entityType)
	}

	// Appliquer le scope si fourni (permissions audit)
	if scopeParam != nil {
		if queryScope, ok := scopeParam.(*scope.QueryScope); ok {
			query = scope.ApplyAuditScope(query, queryScope)
		}
	}

	// Compter le total après filtres
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Si aucun résultat, retourner tout de suite
	if total == 0 {
		return []models.AuditLog{}, 0, nil
	}

	// Récupérer la page demandée avec les préchargements nécessaires
	var logs []models.AuditLog
	if err := query.
		Preload("User").
		Preload("User.Role").
		Order("audit_logs.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// Delete supprime un log d'audit
func (r *auditLogRepository) Delete(id uint) error {
	return database.DB.Delete(&models.AuditLog{}, id).Error
}

// DeleteOld supprime les logs d'audit plus anciens qu'une date donnée (pour nettoyage)
func (r *auditLogRepository) DeleteOld(olderThan time.Time) error {
	return database.DB.Where("created_at < ?", olderThan).Delete(&models.AuditLog{}).Error
}
