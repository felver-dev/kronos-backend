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
	GetAll(page, limit int, userID *uint, action, entityType string) (*dto.AuditLogListResponse, error)
	GetByID(id uint) (*dto.AuditLogDTO, error)
	GetByUserID(userID uint, startDate, endDate *time.Time) ([]dto.AuditLogDTO, error)
	GetByAction(action string) ([]dto.AuditLogDTO, error)
	GetByEntity(entityType string, entityID uint) ([]dto.AuditLogDTO, error)
	GetTicketAuditTrail(ticketID uint) ([]dto.AuditLogDTO, error)
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
func (s *auditService) GetAll(page, limit int, userID *uint, action, entityType string) (*dto.AuditLogListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 50
	}

	// TODO: Implémenter la pagination et les filtres dans le repository
	// Pour l'instant, on récupère les logs récents
	logs, err := s.auditLogRepo.FindRecent(limit * page)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des logs d'audit")
	}

	// Filtrer par userID si fourni
	if userID != nil {
		filtered := []models.AuditLog{}
		for _, log := range logs {
			if log.UserID != nil && *log.UserID == *userID {
				filtered = append(filtered, log)
			}
		}
		logs = filtered
	}

	// Filtrer par action si fournie
	if action != "" {
		filtered := []models.AuditLog{}
		for _, log := range logs {
			if log.Action == action {
				filtered = append(filtered, log)
			}
		}
		logs = filtered
	}

	// Filtrer par entityType si fourni
	if entityType != "" {
		filtered := []models.AuditLog{}
		for _, log := range logs {
			if log.EntityType == entityType {
				filtered = append(filtered, log)
			}
		}
		logs = filtered
	}

	// Pagination manuelle
	total := len(logs)
	start := (page - 1) * limit
	end := start + limit
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	var paginatedLogs []models.AuditLog
	if start < total {
		paginatedLogs = logs[start:end]
	}

	logDTOs := make([]dto.AuditLogDTO, len(paginatedLogs))
	for i, log := range paginatedLogs {
		logDTOs[i] = s.auditLogToDTO(&log)
	}

	return &dto.AuditLogListResponse{
		Logs: logDTOs,
		Pagination: dto.PaginationDTO{
			Page:  page,
			Limit: limit,
			Total: int64(total),
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
func (s *auditService) GetByUserID(userID uint, startDate, endDate *time.Time) ([]dto.AuditLogDTO, error) {
	var logs []models.AuditLog
	var err error

	if startDate != nil && endDate != nil {
		logs, err = s.auditLogRepo.FindByDateRange(*startDate, *endDate)
		if err != nil {
			return nil, errors.New("erreur lors de la récupération des logs d'audit")
		}
		// Filtrer par userID
		filtered := []models.AuditLog{}
		for _, log := range logs {
			if log.UserID != nil && *log.UserID == userID {
				filtered = append(filtered, log)
			}
		}
		logs = filtered
	} else {
		logs, err = s.auditLogRepo.FindByUserID(userID)
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
func (s *auditService) GetByAction(action string) ([]dto.AuditLogDTO, error) {
	logs, err := s.auditLogRepo.FindByAction(action)
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
func (s *auditService) GetByEntity(entityType string, entityID uint) ([]dto.AuditLogDTO, error) {
	logs, err := s.auditLogRepo.FindByEntity(entityType, entityID)
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
func (s *auditService) GetTicketAuditTrail(ticketID uint) ([]dto.AuditLogDTO, error) {
	return s.GetByEntity("ticket", ticketID)
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

