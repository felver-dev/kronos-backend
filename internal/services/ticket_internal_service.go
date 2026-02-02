package services

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/models"
	"github.com/mcicare/itsm-backend/internal/repositories"
)

// TicketInternalService interface pour les opérations sur les tickets internes
type TicketInternalService interface {
	Create(req dto.CreateTicketInternalRequest, createdByID uint, allowAssignAny bool) (*dto.TicketInternalDTO, error)
	GetByID(id uint) (*dto.TicketInternalDTO, error)
	GetAll(scopeParam interface{}, page, limit int, status string, departmentID, filialeID *uint) (*dto.TicketInternalListResponse, error)
	GetPanier(userID uint, page, limit int) (*dto.TicketInternalListResponse, error)
	GetMyPerformance(userID uint) (*dto.TicketInternalPerformanceDTO, error)
	Update(id uint, req dto.UpdateTicketInternalRequest, updatedByID uint) (*dto.TicketInternalDTO, error)
	Assign(id uint, req dto.AssignTicketInternalRequest, assignedByID uint, allowAssignAny bool) (*dto.TicketInternalDTO, error)
	ChangeStatus(id uint, status string, changedByID uint) (*dto.TicketInternalDTO, error)
	Validate(id uint, validatedByID uint) (*dto.TicketInternalDTO, error)
	Close(id uint, closedByID uint) (*dto.TicketInternalDTO, error)
	Delete(id uint) error
}

type ticketInternalService struct {
	repo                repositories.TicketInternalRepository
	userRepo            repositories.UserRepository
	departmentRepo      repositories.DepartmentRepository
	notificationService NotificationService
}

// NewTicketInternalService crée une nouvelle instance
func NewTicketInternalService(
	repo repositories.TicketInternalRepository,
	userRepo repositories.UserRepository,
	departmentRepo repositories.DepartmentRepository,
	notificationService NotificationService,
) TicketInternalService {
	return &ticketInternalService{
		repo:                repo,
		userRepo:            userRepo,
		departmentRepo:      departmentRepo,
		notificationService: notificationService,
	}
}

func (s *ticketInternalService) Create(req dto.CreateTicketInternalRequest, createdByID uint, allowAssignAny bool) (*dto.TicketInternalDTO, error) {
	creator, err := s.userRepo.FindByID(createdByID)
	if err != nil {
		return nil, errors.New("utilisateur créateur introuvable")
	}
	// Si un assigné est fourni et que l'utilisateur n'a pas le droit d'assigner à n'importe qui, vérifier que l'assigné est du même département
	if req.AssignedToID != nil && *req.AssignedToID != 0 && !allowAssignAny {
		if creator.DepartmentID == nil {
			return nil, errors.New("vous ne pouvez assigner qu'à un membre de votre département")
		}
		assignee, errAssign := s.userRepo.FindByID(*req.AssignedToID)
		if errAssign != nil {
			return nil, errors.New("utilisateur assigné introuvable")
		}
		if assignee.DepartmentID == nil || *assignee.DepartmentID != *creator.DepartmentID {
			return nil, errors.New("vous ne pouvez assigner qu'à un membre de votre département")
		}
	}
	dept, err := s.departmentRepo.FindByID(req.DepartmentID)
	if err != nil || dept == nil {
		return nil, errors.New("département introuvable")
	}
	if dept.IsITDepartment {
		return nil, errors.New("les tickets internes ne concernent que les départements non-IT")
	}
	if dept.FilialeID == nil {
		return nil, errors.New("le département doit être rattaché à une filiale")
	}

	year := time.Now().Year()
	seq, err := s.repo.GetNextSequenceNumber(year)
	if err != nil {
		return nil, fmt.Errorf("génération du code: %w", err)
	}
	code := fmt.Sprintf("TKI-%d-%04d", year, seq)

	t := &models.TicketInternal{
		Code:          code,
		Title:         req.Title,
		Description:   req.Description,
		Category:      req.Category,
		Status:        "ouvert",
		Priority:      req.Priority,
		DepartmentID:  req.DepartmentID,
		FilialeID:     *dept.FilialeID,
		CreatedByID:   createdByID,
		AssignedToID:  req.AssignedToID,
		EstimatedTime: req.EstimatedTime,
		TicketID:      req.TicketID,
	}
	if t.Priority == "" {
		t.Priority = "medium"
	}
	if err := s.repo.Create(t); err != nil {
		return nil, err
	}
	loaded, _ := s.repo.FindByID(t.ID)
	dtoOut := s.toDTO(loaded)
	// Notification à l'assigné si le ticket a été créé avec un assigné
	if dtoOut != nil && req.AssignedToID != nil && *req.AssignedToID != 0 {
		s.notifyTicketInternalAssigned(*req.AssignedToID, dtoOut.Code, dtoOut.Title, dtoOut.ID, createdByID)
	}
	return dtoOut, nil
}

func (s *ticketInternalService) GetByID(id uint) (*dto.TicketInternalDTO, error) {
	t, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("ticket interne introuvable")
	}
	return s.toDTO(t), nil
}

func (s *ticketInternalService) GetAll(scopeParam interface{}, page, limit int, status string, departmentID, filialeID *uint) (*dto.TicketInternalListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	list, total, err := s.repo.FindAll(scopeParam, page, limit, status, departmentID, filialeID)
	if err != nil {
		return nil, err
	}
	dtos := make([]dto.TicketInternalDTO, 0, len(list))
	for i := range list {
		dtos = append(dtos, *s.toDTO(&list[i]))
	}
	return &dto.TicketInternalListResponse{
		Tickets: dtos,
		Pagination: dto.PaginationDTO{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: int((total + int64(limit) - 1) / int64(limit)),
		},
	}, nil
}

func (s *ticketInternalService) Update(id uint, req dto.UpdateTicketInternalRequest, updatedByID uint) (*dto.TicketInternalDTO, error) {
	t, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("ticket interne introuvable")
	}
	if req.Title != "" {
		t.Title = req.Title
	}
	if req.Description != "" {
		t.Description = req.Description
	}
	if req.Category != "" {
		t.Category = req.Category
	}
	if req.Priority != "" {
		t.Priority = req.Priority
	}
	if req.EstimatedTime != nil {
		t.EstimatedTime = req.EstimatedTime
	}
	if req.ActualTime != nil {
		t.ActualTime = req.ActualTime
	}
	if req.AssignedToID != nil {
		t.AssignedToID = req.AssignedToID
	}
	if req.Status != "" {
		t.Status = req.Status
		if req.Status == "cloture" {
			now := time.Now()
			t.ClosedAt = &now
		}
	}
	if err := s.repo.Update(t); err != nil {
		return nil, err
	}
	loaded, _ := s.repo.FindByID(id)
	return s.toDTO(loaded), nil
}

func (s *ticketInternalService) Assign(id uint, req dto.AssignTicketInternalRequest, assignedByID uint, allowAssignAny bool) (*dto.TicketInternalDTO, error) {
	t, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("ticket interne introuvable")
	}
	if t.Status == "cloture" || t.Status == "resolu" {
		return nil, errors.New("impossible d'assigner un ticket clôturé ou résolu")
	}
	// Sans view_all, on ne peut assigner qu'à un membre de son département
	if req.AssignedToID != nil && *req.AssignedToID != 0 && !allowAssignAny {
		assigner, errAssigner := s.userRepo.FindByID(assignedByID)
		if errAssigner != nil {
			return nil, errors.New("utilisateur introuvable")
		}
		if assigner.DepartmentID == nil {
			return nil, errors.New("vous ne pouvez assigner qu'à un membre de votre département")
		}
		assignee, errAssignee := s.userRepo.FindByID(*req.AssignedToID)
		if errAssignee != nil {
			return nil, errors.New("utilisateur assigné introuvable")
		}
		if assignee.DepartmentID == nil || *assignee.DepartmentID != *assigner.DepartmentID {
			return nil, errors.New("vous ne pouvez assigner qu'à un membre de votre département")
		}
	}
	updates := map[string]interface{}{}
	if req.AssignedToID != nil {
		updates["assigned_to_id"] = *req.AssignedToID
	}
	if req.EstimatedTime != nil {
		updates["estimated_time"] = *req.EstimatedTime
		if t.Status == "ouvert" {
			updates["status"] = "en_cours"
		}
	}
	if len(updates) == 0 {
		loaded, _ := s.repo.FindByID(id)
		return s.toDTO(loaded), nil
	}
	if err := s.repo.UpdateFields(id, updates); err != nil {
		return nil, err
	}
	loaded, _ := s.repo.FindByID(id)
	dtoOut := s.toDTO(loaded)
	// Notification à l'assigné lorsqu'un ticket est assigné
	if dtoOut != nil && req.AssignedToID != nil && *req.AssignedToID != 0 {
		s.notifyTicketInternalAssigned(*req.AssignedToID, dtoOut.Code, dtoOut.Title, dtoOut.ID, assignedByID)
	}
	return dtoOut, nil
}

func (s *ticketInternalService) GetPanier(userID uint, page, limit int) (*dto.TicketInternalListResponse, error) {
	if userID == 0 {
		return &dto.TicketInternalListResponse{Tickets: []dto.TicketInternalDTO{}, Pagination: dto.PaginationDTO{Page: page, Limit: limit, Total: 0, TotalPages: 0}}, nil
	}
	list, total, err := s.repo.FindPanierByUser(userID, page, limit)
	if err != nil {
		return nil, err
	}
	tickets := make([]dto.TicketInternalDTO, 0, len(list))
	for i := range list {
		tickets = append(tickets, *s.toDTO(&list[i]))
	}
	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}
	if totalPages < 1 {
		totalPages = 1
	}
	return &dto.TicketInternalListResponse{
		Tickets: tickets,
		Pagination: dto.PaginationDTO{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *ticketInternalService) GetMyPerformance(userID uint) (*dto.TicketInternalPerformanceDTO, error) {
	if userID == 0 {
		return &dto.TicketInternalPerformanceDTO{}, nil
	}
	totalAssigned, resolved, inProgress, open, totalTimeSpent, err := s.repo.GetPerformanceByAssignedUser(userID)
	if err != nil {
		return nil, err
	}
	efficiency := 0.0
	if totalAssigned > 0 {
		efficiency = float64(resolved) / float64(totalAssigned) * 100
	}
	return &dto.TicketInternalPerformanceDTO{
		TotalAssigned:  totalAssigned,
		Resolved:       resolved,
		InProgress:     inProgress,
		Open:           open,
		TotalTimeSpent: totalTimeSpent,
		Efficiency:     efficiency,
	}, nil
}

// notifyTicketInternalAssigned envoie une notification à l'utilisateur assigné
func (s *ticketInternalService) notifyTicketInternalAssigned(assigneeID uint, code, title string, ticketID uint, assignedByID uint) {
	if s.notificationService == nil {
		return
	}
	linkURL := fmt.Sprintf("/app/ticket-internes/%d", ticketID)
	notificationTitle := fmt.Sprintf("Ticket interne assigné : %s", code)
	notificationMessage := fmt.Sprintf("Un ticket interne vous a été assigné : %s - %s. Consultez-le dans votre panier.", code, title)
	metadata := map[string]any{"ticket_internal_id": ticketID, "code": code, "assigned_by_id": assignedByID}
	if err := s.notificationService.Create(assigneeID, "ticket_internal_assigned", notificationTitle, notificationMessage, linkURL, metadata); err != nil {
		log.Printf("Erreur notification ticket interne assigné (user %d): %v", assigneeID, err)
	}
}

func (s *ticketInternalService) ChangeStatus(id uint, status string, changedByID uint) (*dto.TicketInternalDTO, error) {
	valid := map[string]bool{"ouvert": true, "en_cours": true, "en_attente": true, "resolu": true, "cloture": true}
	if !valid[status] {
		return nil, errors.New("statut invalide")
	}
	_, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("ticket interne introuvable")
	}
	updates := map[string]interface{}{"status": status}
	if status == "cloture" {
		now := time.Now()
		updates["closed_at"] = now
	}
	if err := s.repo.UpdateFields(id, updates); err != nil {
		return nil, err
	}
	loaded, _ := s.repo.FindByID(id)
	return s.toDTO(loaded), nil
}

func (s *ticketInternalService) Validate(id uint, validatedByID uint) (*dto.TicketInternalDTO, error) {
	t, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("ticket interne introuvable")
	}
	if t.Status != "en_attente" {
		return nil, errors.New("seuls les tickets en attente de validation peuvent être validés")
	}
	now := time.Now()
	updates := map[string]interface{}{
		"status":              "resolu",
		"validated_by_user_id": validatedByID,
		"validated_at":        now,
	}
	if err := s.repo.UpdateFields(id, updates); err != nil {
		return nil, err
	}
	loaded, _ := s.repo.FindByID(id)
	return s.toDTO(loaded), nil
}

func (s *ticketInternalService) Close(id uint, closedByID uint) (*dto.TicketInternalDTO, error) {
	return s.ChangeStatus(id, "cloture", closedByID)
}

func (s *ticketInternalService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *ticketInternalService) toDTO(t *models.TicketInternal) *dto.TicketInternalDTO {
	if t == nil {
		return nil
	}
	d := &dto.TicketInternalDTO{
		ID:                t.ID,
		Code:              t.Code,
		Title:             t.Title,
		Description:       t.Description,
		Category:          t.Category,
		Status:             t.Status,
		Priority:          t.Priority,
		DepartmentID:      t.DepartmentID,
		FilialeID:         t.FilialeID,
		CreatedByID:       t.CreatedByID,
		EstimatedTime:     t.EstimatedTime,
		ActualTime:        t.ActualTime,
		TicketID:          t.TicketID,
		CreatedAt:         t.CreatedAt,
		UpdatedAt:         t.UpdatedAt,
		ClosedAt:          t.ClosedAt,
		ValidatedByUserID:  t.ValidatedByUserID,
		ValidatedAt:       t.ValidatedAt,
	}
	if t.CreatedBy.ID != 0 {
		d.CreatedBy = userToDTO(&t.CreatedBy)
	}
	if t.AssignedTo != nil {
		d.AssignedToID = t.AssignedToID
		d.AssignedTo = ptrUserDTO(t.AssignedTo)
	}
	if t.ValidatedBy != nil {
		d.ValidatedBy = ptrUserDTO(t.ValidatedBy)
	}
	if t.Department.ID != 0 {
		d.Department = departmentToDTO(&t.Department)
	}
	if t.Filiale.ID != 0 {
		d.Filiale = filialeToDTO(&t.Filiale)
	}
	return d
}

func ptrUserDTO(u *models.User) *dto.UserDTO {
	if u == nil {
		return nil
	}
	d := userToDTO(u)
	return &d
}

func departmentToDTO(d *models.Department) *dto.DepartmentDTO {
	if d == nil {
		return nil
	}
	out := &dto.DepartmentDTO{
		ID:             d.ID,
		Name:           d.Name,
		Code:           d.Code,
		IsActive:       d.IsActive,
		IsITDepartment: d.IsITDepartment,
		FilialeID:      d.FilialeID,
	}
	return out
}

func filialeToDTO(f *models.Filiale) *dto.FilialeDTO {
	if f == nil {
		return nil
	}
	return &dto.FilialeDTO{
		ID:                 f.ID,
		Name:               f.Name,
		Code:               f.Code,
		IsSoftwareProvider: f.IsSoftwareProvider,
	}
}
