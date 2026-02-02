package services

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/models"
	"github.com/mcicare/itsm-backend/internal/repositories"
)

// AuditService interface pour les opérations sur les logs d'audit
type AuditService interface {
	GetAll(scope interface{}, page, limit int, userID *uint, action, entityType string) (*dto.AuditLogListResponse, error) // scope peut être *scope.QueryScope ou nil
	GetByID(id uint) (*dto.AuditLogDTO, error)
	GetByUserID(scope interface{}, userID uint, startDate, endDate *time.Time) ([]dto.AuditLogDTO, error)
	GetByAction(scope interface{}, action string) ([]dto.AuditLogDTO, error)
	GetByEntity(scope interface{}, entityType string, entityID uint) ([]dto.AuditLogDTO, error)
	GetTicketAuditTrail(scope interface{}, ticketID uint) ([]dto.AuditLogDTO, error)
}

// auditService implémente AuditService
type auditService struct {
	auditLogRepo repositories.AuditLogRepository
}

// NewAuditService crée une nouvelle instance de AuditService
func NewAuditService(auditLogRepo repositories.AuditLogRepository) AuditService {
	return &auditService{
		auditLogRepo: auditLogRepo,
	}
}

// GetAll récupère tous les logs d'audit avec pagination et filtres
// Le scope est utilisé pour filtrer automatiquement selon les permissions de l'utilisateur
func (s *auditService) GetAll(scopeParam interface{}, page, limit int, userID *uint, action, entityType string) (*dto.AuditLogListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 50
	}

	// Récupérer les logs paginés + le total via le repository
	logs, total, err := s.auditLogRepo.FindPaginated(scopeParam, page, limit, userID, action, entityType)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des logs d'audit")
	}

	logDTOs := make([]dto.AuditLogDTO, len(logs))
	for i, log := range logs {
		logDTOs[i] = s.auditLogToDTO(&log)
	}

	return &dto.AuditLogListResponse{
		Logs: logDTOs,
		Pagination: dto.PaginationDTO{
			Page:  page,
			Limit: limit,
			Total: total,
		},
	}, nil
}

// GetByID récupère un log d'audit par son ID
func (s *auditService) GetByID(id uint) (*dto.AuditLogDTO, error) {
	log, err := s.auditLogRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("log d'audit introuvable")
	}

	logDTO := s.auditLogToDTO(log)
	return &logDTO, nil
}

// GetByUserID récupère les logs d'audit d'un utilisateur
// Le scope est utilisé pour filtrer automatiquement selon les permissions de l'utilisateur
func (s *auditService) GetByUserID(scopeParam interface{}, userID uint, startDate, endDate *time.Time) ([]dto.AuditLogDTO, error) {
	var logs []models.AuditLog
	var err error

	if startDate != nil && endDate != nil {
		logs, err = s.auditLogRepo.FindByDateRange(scopeParam, *startDate, *endDate)
		if err != nil {
			return nil, errors.New("erreur lors de la récupération des logs d'audit")
		}
		// Filtrer par userID (le scope a déjà été appliqué)
		filtered := []models.AuditLog{}
		for _, log := range logs {
			if log.UserID != nil && *log.UserID == userID {
				filtered = append(filtered, log)
			}
		}
		logs = filtered
	} else {
		logs, err = s.auditLogRepo.FindByUserID(scopeParam, userID)
		if err != nil {
			return nil, errors.New("erreur lors de la récupération des logs d'audit")
		}
	}

	logDTOs := make([]dto.AuditLogDTO, len(logs))
	for i, log := range logs {
		logDTOs[i] = s.auditLogToDTO(&log)
	}

	return logDTOs, nil
}

// GetByAction récupère les logs d'audit d'une action
// Le scope est utilisé pour filtrer automatiquement selon les permissions de l'utilisateur
func (s *auditService) GetByAction(scopeParam interface{}, action string) ([]dto.AuditLogDTO, error) {
	logs, err := s.auditLogRepo.FindByAction(scopeParam, action)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des logs d'audit")
	}

	logDTOs := make([]dto.AuditLogDTO, len(logs))
	for i, log := range logs {
		logDTOs[i] = s.auditLogToDTO(&log)
	}

	return logDTOs, nil
}

// GetByEntity récupère les logs d'audit d'une entité
// Le scope est utilisé pour filtrer automatiquement selon les permissions de l'utilisateur
func (s *auditService) GetByEntity(scopeParam interface{}, entityType string, entityID uint) ([]dto.AuditLogDTO, error) {
	logs, err := s.auditLogRepo.FindByEntity(scopeParam, entityType, entityID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des logs d'audit")
	}

	logDTOs := make([]dto.AuditLogDTO, len(logs))
	for i, log := range logs {
		logDTOs[i] = s.auditLogToDTO(&log)
	}

	return logDTOs, nil
}

// GetTicketAuditTrail récupère la piste d'audit d'un ticket
// Le scope est utilisé pour filtrer automatiquement selon les permissions de l'utilisateur
func (s *auditService) GetTicketAuditTrail(scopeParam interface{}, ticketID uint) ([]dto.AuditLogDTO, error) {
	return s.GetByEntity(scopeParam, "ticket", ticketID)
}

// auditLogToDTO convertit un modèle AuditLog en DTO
func (s *auditService) auditLogToDTO(log *models.AuditLog) dto.AuditLogDTO {
	logDTO := dto.AuditLogDTO{
		ID:          log.ID,
		Action:      log.Action,
		EntityType:  log.EntityType,
		EntityID:    log.EntityID,
		IPAddress:   log.IPAddress,
		UserAgent:   log.UserAgent,
		Description: log.Description,
		CreatedAt:   log.CreatedAt,
	}

	if log.UserID != nil {
		logDTO.UserID = log.UserID
		if log.User != nil {
			userDTO := dto.UserDTO{
				ID:        log.User.ID,
				Username:  log.User.Username,
				Email:     log.User.Email,
				FirstName: log.User.FirstName,
				LastName:  log.User.LastName,
				Role:      log.User.Role.Name,
				IsActive:  log.User.IsActive,
				CreatedAt: log.User.CreatedAt,
				UpdatedAt: log.User.UpdatedAt,
			}
			logDTO.User = &userDTO
		}
	}

	// Convertir les valeurs JSON si présentes
	if log.OldValues != nil && len(log.OldValues) > 0 {
		var oldValues map[string]interface{}
		if err := json.Unmarshal(log.OldValues, &oldValues); err == nil {
			logDTO.OldValues = oldValues
		}
	}
	if log.NewValues != nil && len(log.NewValues) > 0 {
		var newValues map[string]interface{}
		if err := json.Unmarshal(log.NewValues, &newValues); err == nil {
			logDTO.NewValues = newValues
		}
	}

	return logDTO
}

