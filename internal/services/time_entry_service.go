package services

import (
	"errors"
	"time"

	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/models"
	"github.com/mcicare/itsm-backend/internal/repositories"
)

// TimeEntryService interface pour les opérations sur les entrées de temps
type TimeEntryService interface {
	Create(req dto.CreateTimeEntryRequest, userID uint) (*dto.TimeEntryDTO, error)
	GetByID(id uint) (*dto.TimeEntryDTO, error)
	GetAll(scope interface{}) ([]dto.TimeEntryDTO, error) // scope peut être *scope.QueryScope ou nil
	GetByTicketID(ticketID uint) ([]dto.TimeEntryDTO, error)
	GetByUserID(userID uint) ([]dto.TimeEntryDTO, error)
	GetByDateRange(userID uint, startDate, endDate time.Time) ([]dto.TimeEntryDTO, error)
	GetValidated(scope interface{}) ([]dto.TimeEntryDTO, error)
	GetPendingValidation(scope interface{}) ([]dto.TimeEntryDTO, error)
	Update(id uint, req dto.UpdateTimeEntryRequest, updatedByID uint) (*dto.TimeEntryDTO, error)
	Validate(id uint, req dto.ValidateTimeEntryRequest, validatedByID uint) (*dto.TimeEntryDTO, error)
	Delete(id uint) error
	GetTotalByTicketID(ticketID uint) (int, error)
	GetTotalByUserID(userID uint) (int, error)
}

// timeEntryService implémente TimeEntryService
type timeEntryService struct {
	timeEntryRepo repositories.TimeEntryRepository
	ticketRepo    repositories.TicketRepository
	userRepo      repositories.UserRepository
	delayRepo     repositories.DelayRepository
}

// NewTimeEntryService crée une nouvelle instance de TimeEntryService
func NewTimeEntryService(
	timeEntryRepo repositories.TimeEntryRepository,
	ticketRepo repositories.TicketRepository,
	userRepo repositories.UserRepository,
	delayRepo repositories.DelayRepository,
) TimeEntryService {
	return &timeEntryService{
		timeEntryRepo: timeEntryRepo,
		ticketRepo:    ticketRepo,
		userRepo:      userRepo,
		delayRepo:     delayRepo,
	}
}

// Create crée une nouvelle entrée de temps
func (s *timeEntryService) Create(req dto.CreateTimeEntryRequest, userID uint) (*dto.TimeEntryDTO, error) {
	// Vérifier que le ticket existe
	_, err := s.ticketRepo.FindByID(req.TicketID)
	if err != nil {
		return nil, errors.New("ticket introuvable")
	}

	// Parser la date
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, errors.New("format de date invalide, attendu: YYYY-MM-DD")
	}

	ticketID := req.TicketID
	timeEntry := &models.TimeEntry{
		TicketID:    &ticketID,
		UserID:      userID,
		TimeSpent:   req.TimeSpent,
		Date:        date,
		Description: req.Description,
		Validated:   false,
	}

	if err := s.timeEntryRepo.Create(timeEntry); err != nil {
		return nil, errors.New("erreur lors de la création de l'entrée de temps")
	}

	// Mettre à jour le temps réel du ticket + retards
	s.updateTicketActualTime(req.TicketID, userID)

	// Récupérer l'entrée créée avec ses relations
	createdEntry, err := s.timeEntryRepo.FindByID(timeEntry.ID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération de l'entrée créée")
	}

	// Convertir en DTO
	entryDTO := s.timeEntryToDTO(createdEntry)
	return &entryDTO, nil
}

// GetByID récupère une entrée de temps par son ID
func (s *timeEntryService) GetByID(id uint) (*dto.TimeEntryDTO, error) {
	timeEntry, err := s.timeEntryRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("entrée de temps introuvable")
	}

	entryDTO := s.timeEntryToDTO(timeEntry)
	return &entryDTO, nil
}

// GetAll récupère toutes les entrées de temps
// Le scope est utilisé pour filtrer automatiquement selon les permissions de l'utilisateur
func (s *timeEntryService) GetAll(scopeParam interface{}) ([]dto.TimeEntryDTO, error) {
	timeEntries, err := s.timeEntryRepo.FindAll(scopeParam)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des entrées de temps: " + err.Error())
	}

	var entryDTOs []dto.TimeEntryDTO
	for _, entry := range timeEntries {
		// Gérer les erreurs de conversion gracieusement
		entryDTO := s.timeEntryToDTOSafe(&entry)
		entryDTOs = append(entryDTOs, entryDTO)
	}

	return entryDTOs, nil
}

// GetByTicketID récupère les entrées de temps d'un ticket
func (s *timeEntryService) GetByTicketID(ticketID uint) ([]dto.TimeEntryDTO, error) {
	timeEntries, err := s.timeEntryRepo.FindByTicketID(ticketID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des entrées de temps")
	}

	var entryDTOs []dto.TimeEntryDTO
	for _, entry := range timeEntries {
		entryDTOs = append(entryDTOs, s.timeEntryToDTO(&entry))
	}

	return entryDTOs, nil
}

// GetByUserID récupère les entrées de temps d'un utilisateur
func (s *timeEntryService) GetByUserID(userID uint) ([]dto.TimeEntryDTO, error) {
	timeEntries, err := s.timeEntryRepo.FindByUserID(userID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des entrées de temps")
	}

	var entryDTOs []dto.TimeEntryDTO
	for _, entry := range timeEntries {
		entryDTOs = append(entryDTOs, s.timeEntryToDTO(&entry))
	}

	return entryDTOs, nil
}

// GetByDateRange récupère les entrées de temps d'un utilisateur dans une plage de dates
func (s *timeEntryService) GetByDateRange(userID uint, startDate, endDate time.Time) ([]dto.TimeEntryDTO, error) {
	timeEntries, err := s.timeEntryRepo.FindByDateRange(userID, startDate, endDate)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des entrées de temps")
	}

	var entryDTOs []dto.TimeEntryDTO
	for _, entry := range timeEntries {
		entryDTOs = append(entryDTOs, s.timeEntryToDTO(&entry))
	}

	return entryDTOs, nil
}

// GetValidated récupère les entrées de temps validées
func (s *timeEntryService) GetValidated(scopeParam interface{}) ([]dto.TimeEntryDTO, error) {
	timeEntries, err := s.timeEntryRepo.FindValidated(scopeParam)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des entrées de temps")
	}

	var entryDTOs []dto.TimeEntryDTO
	for _, entry := range timeEntries {
		entryDTOs = append(entryDTOs, s.timeEntryToDTO(&entry))
	}

	return entryDTOs, nil
}

// GetPendingValidation récupère les entrées de temps en attente de validation
func (s *timeEntryService) GetPendingValidation(scopeParam interface{}) ([]dto.TimeEntryDTO, error) {
	timeEntries, err := s.timeEntryRepo.FindPendingValidation(scopeParam)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des entrées de temps")
	}

	var entryDTOs []dto.TimeEntryDTO
	for _, entry := range timeEntries {
		entryDTOs = append(entryDTOs, s.timeEntryToDTO(&entry))
	}

	return entryDTOs, nil
}

// Update met à jour une entrée de temps
func (s *timeEntryService) Update(id uint, req dto.UpdateTimeEntryRequest, updatedByID uint) (*dto.TimeEntryDTO, error) {
	timeEntry, err := s.timeEntryRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("entrée de temps introuvable")
	}

	// Vérifier que l'entrée n'est pas validée (on ne peut pas modifier une entrée validée)
	if timeEntry.Validated {
		return nil, errors.New("impossible de modifier une entrée de temps validée")
	}

	// Mettre à jour les champs fournis
	if req.TimeSpent > 0 {
		timeEntry.TimeSpent = req.TimeSpent
	}
	if req.Description != "" {
		timeEntry.Description = req.Description
	}

	if err := s.timeEntryRepo.Update(timeEntry); err != nil {
		return nil, errors.New("erreur lors de la mise à jour de l'entrée de temps")
	}

	// Mettre à jour le temps réel du ticket (uniquement si lié à un ticket normal)
	if timeEntry.TicketID != nil {
		s.updateTicketActualTime(*timeEntry.TicketID, updatedByID)
	}

	// Récupérer l'entrée mise à jour
	updatedEntry, err := s.timeEntryRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération de l'entrée mise à jour")
	}

	entryDTO := s.timeEntryToDTO(updatedEntry)
	return &entryDTO, nil
}

// Validate valide ou invalide une entrée de temps
func (s *timeEntryService) Validate(id uint, req dto.ValidateTimeEntryRequest, validatedByID uint) (*dto.TimeEntryDTO, error) {
	timeEntry, err := s.timeEntryRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("entrée de temps introuvable")
	}

	// Vérifier que le validateur existe
	_, err = s.userRepo.FindByID(validatedByID)
	if err != nil {
		return nil, errors.New("validateur introuvable")
	}

	now := time.Now()
	timeEntry.Validated = *req.Validated
	timeEntry.ValidatedByID = &validatedByID
	timeEntry.ValidatedAt = &now

	if err := s.timeEntryRepo.Update(timeEntry); err != nil {
		return nil, errors.New("erreur lors de la validation de l'entrée de temps")
	}

	// Récupérer l'entrée mise à jour
	updatedEntry, err := s.timeEntryRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération de l'entrée mise à jour")
	}

	entryDTO := s.timeEntryToDTO(updatedEntry)
	return &entryDTO, nil
}

// Delete supprime une entrée de temps
func (s *timeEntryService) Delete(id uint) error {
	timeEntry, err := s.timeEntryRepo.FindByID(id)
	if err != nil {
		return errors.New("entrée de temps introuvable")
	}

	// Vérifier que l'entrée n'est pas validée
	if timeEntry.Validated {
		return errors.New("impossible de supprimer une entrée de temps validée")
	}

	if err := s.timeEntryRepo.Delete(id); err != nil {
		return errors.New("erreur lors de la suppression de l'entrée de temps")
	}

	if timeEntry.TicketID != nil {
		s.updateTicketActualTime(*timeEntry.TicketID, timeEntry.UserID)
	}

	return nil
}

// GetTotalByTicketID calcule le temps total passé sur un ticket
func (s *timeEntryService) GetTotalByTicketID(ticketID uint) (int, error) {
	return s.timeEntryRepo.SumByTicketID(ticketID)
}

// GetTotalByUserID calcule le temps total passé par un utilisateur
func (s *timeEntryService) GetTotalByUserID(userID uint) (int, error) {
	return s.timeEntryRepo.SumByUserID(userID)
}

// updateTicketActualTime met à jour le temps réel d'un ticket
func (s *timeEntryService) updateTicketActualTime(ticketID uint, userID uint) {
	total, err := s.timeEntryRepo.SumByTicketID(ticketID)
	if err != nil {
		return
	}

	ticket, err := s.ticketRepo.FindByID(ticketID)
	if err != nil {
		return
	}

	ticket.ActualTime = &total
	s.ticketRepo.Update(ticket)

	// Détection du retard
	if ticket.EstimatedTime == nil || *ticket.EstimatedTime <= 0 {
		return
	}
	estimated := *ticket.EstimatedTime
	delayTime := total - estimated

	if delayTime <= 0 {
		// Si plus de retard, supprimer uniquement si encore "unjustified"
		if existing, err := s.delayRepo.FindByTicketID(ticketID); err == nil {
			if existing.Status == "unjustified" {
				_ = s.delayRepo.Delete(existing.ID)
			}
		}
		return
	}

	percentage := float64(delayTime) / float64(estimated) * 100
	if percentage > 999.99 {
		percentage = 999.99
	}
	existing, err := s.delayRepo.FindByTicketID(ticketID)
	if err != nil || existing == nil {
		ownerID := userID
		if ownerID == 0 {
			if ticket.AssignedToID != nil {
				ownerID = *ticket.AssignedToID
			} else {
				ownerID = ticket.CreatedByID
			}
		}
		tid := ticketID
		delay := &models.Delay{
			TicketID:        &tid,
			UserID:          ownerID,
			EstimatedTime:   estimated,
			ActualTime:      total,
			DelayTime:       delayTime,
			DelayPercentage: percentage,
			Status:          "unjustified",
			DetectedAt:      time.Now(),
		}
		_ = s.delayRepo.Create(delay)
		return
	}

	existing.ActualTime = total
	existing.EstimatedTime = estimated
	existing.DelayTime = delayTime
	existing.DelayPercentage = percentage
	if userID != 0 {
		existing.UserID = userID
	}
	if existing.Status == "rejected" {
		existing.Status = "unjustified"
	}
	_ = s.delayRepo.Update(existing)
}

// timeEntryToDTO convertit un modèle TimeEntry en DTO
func (s *timeEntryService) timeEntryToDTO(timeEntry *models.TimeEntry) dto.TimeEntryDTO {
	ticketID := uint(0)
	if timeEntry.TicketID != nil {
		ticketID = *timeEntry.TicketID
	}
	entryDTO := dto.TimeEntryDTO{
		ID:            timeEntry.ID,
		TicketID:      ticketID,
		ProjectTaskID: timeEntry.ProjectTaskID,
		UserID:        timeEntry.UserID,
		TimeSpent:     timeEntry.TimeSpent,
		Date:          timeEntry.Date,
		Description:   timeEntry.Description,
		Validated:     timeEntry.Validated,
		CreatedAt:     timeEntry.CreatedAt,
		UpdatedAt:     timeEntry.UpdatedAt,
	}

	// Convertir le ticket si présent (vérifier que le ticket a été chargé)
	if timeEntry.Ticket != nil && timeEntry.Ticket.ID != 0 {
		ticketDTO := s.ticketToDTO(timeEntry.Ticket)
		entryDTO.Ticket = &ticketDTO
	}

	// Convertir l'utilisateur si présent (vérifier que l'utilisateur a été chargé)
	if timeEntry.User.ID != 0 {
		userDTO := s.userToDTO(&timeEntry.User)
		entryDTO.User = &userDTO
	}

	// Convertir le validateur si présent
	if timeEntry.ValidatedByID != nil {
		entryDTO.ValidatedBy = timeEntry.ValidatedByID
	}
	if timeEntry.ValidatedAt != nil {
		entryDTO.ValidatedAt = timeEntry.ValidatedAt
	}

	return entryDTO
}

// timeEntryToDTOSafe convertit un modèle TimeEntry en DTO avec gestion d'erreur gracieuse
func (s *timeEntryService) timeEntryToDTOSafe(timeEntry *models.TimeEntry) dto.TimeEntryDTO {
	var ticketID uint
	if timeEntry.TicketID != nil {
		ticketID = *timeEntry.TicketID
	}
	entryDTO := dto.TimeEntryDTO{
		ID:            timeEntry.ID,
		TicketID:      ticketID,
		ProjectTaskID: timeEntry.ProjectTaskID,
		UserID:        timeEntry.UserID,
		TimeSpent:     timeEntry.TimeSpent,
		Date:          timeEntry.Date,
		Description:   timeEntry.Description,
		Validated:     timeEntry.Validated,
		CreatedAt:     timeEntry.CreatedAt,
		UpdatedAt:     timeEntry.UpdatedAt,
	}

	if timeEntry.Ticket != nil && timeEntry.Ticket.ID != 0 {
		ticketDTO := s.ticketToDTOSafe(timeEntry.Ticket)
		entryDTO.Ticket = &ticketDTO
	}

	// Convertir l'utilisateur si présent (vérifier que l'utilisateur a été chargé)
	if timeEntry.User.ID != 0 {
		userDTO := s.userToDTOSafe(&timeEntry.User)
		entryDTO.User = &userDTO
	}

	// Convertir le validateur si présent
	if timeEntry.ValidatedByID != nil {
		entryDTO.ValidatedBy = timeEntry.ValidatedByID
	}
	if timeEntry.ValidatedAt != nil {
		entryDTO.ValidatedAt = timeEntry.ValidatedAt
	}

	return entryDTO
}

// ticketToDTOSafe convertit un modèle Ticket en DTO avec gestion d'erreur gracieuse
func (s *timeEntryService) ticketToDTOSafe(ticket *models.Ticket) dto.TicketDTO {
	ticketDTO := dto.TicketDTO{
		ID:          ticket.ID,
		Code:        ticket.Code,
		Title:       ticket.Title,
		Description: ticket.Description,
		Category:    ticket.Category,
		Source:      ticket.Source,
		Status:      ticket.Status,
		Priority:    ticket.Priority,
		CreatedAt:   ticket.CreatedAt,
		UpdatedAt:   ticket.UpdatedAt,
	}

	if ticket.EstimatedTime != nil {
		ticketDTO.EstimatedTime = ticket.EstimatedTime
	}
	if ticket.ActualTime != nil {
		ticketDTO.ActualTime = ticket.ActualTime
	}
	if ticket.ClosedAt != nil {
		ticketDTO.ClosedAt = ticket.ClosedAt
	}

	// Convertir CreatedBy si présent
	if ticket.CreatedBy.ID != 0 {
		userDTO := s.userToDTOSafe(&ticket.CreatedBy)
		ticketDTO.CreatedBy = userDTO
	}

	// Convertir AssignedTo si présent
	if ticket.AssignedTo != nil && ticket.AssignedTo.ID != 0 {
		userDTO := s.userToDTOSafe(ticket.AssignedTo)
		ticketDTO.AssignedTo = &userDTO
	}

	return ticketDTO
}

// userToDTOSafe convertit un modèle User en DTO avec gestion d'erreur gracieuse
func (s *timeEntryService) userToDTOSafe(user *models.User) dto.UserDTO {
	userDTO := dto.UserDTO{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Avatar:    user.Avatar,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	// Vérifier que le rôle a été chargé
	if user.RoleID != 0 && user.Role.ID != 0 {
		userDTO.Role = user.Role.Name
	} else if user.RoleID != 0 {
		// Si RoleID existe mais Role n'est pas chargé, on laisse vide
		userDTO.Role = ""
	}

	// Ajouter le département si présent (gérer gracieusement les erreurs)
	if user.DepartmentID != nil {
		userDTO.DepartmentID = user.DepartmentID
		if user.Department != nil && user.Department.ID != 0 {
			departmentDTO := dto.DepartmentDTO{
				ID:      user.Department.ID,
				Name:    user.Department.Name,
				Code:    user.Department.Code,
				IsActive: user.Department.IsActive,
			}
			// Ajouter le siège si présent
			if user.Department.Office != nil && user.Department.Office.ID != 0 {
				departmentDTO.Office = &dto.OfficeDTO{
					ID:      user.Department.Office.ID,
					Name:    user.Department.Office.Name,
					Country: user.Department.Office.Country,
					City:    user.Department.Office.City,
				}
			}
			userDTO.Department = &departmentDTO
		}
	}

	if user.LastLogin != nil {
		userDTO.LastLogin = user.LastLogin
	}

	return userDTO
}

// ticketToDTO convertit un modèle Ticket en DTO (méthode helper)
func (s *timeEntryService) ticketToDTO(ticket *models.Ticket) dto.TicketDTO {
	ticketDTO := dto.TicketDTO{
		ID:          ticket.ID,
		Code:        ticket.Code,
		Title:       ticket.Title,
		Description: ticket.Description,
		Category:    ticket.Category,
		Source:      ticket.Source,
		Status:      ticket.Status,
		Priority:    ticket.Priority,
		CreatedAt:   ticket.CreatedAt,
		UpdatedAt:   ticket.UpdatedAt,
	}

	if ticket.EstimatedTime != nil {
		ticketDTO.EstimatedTime = ticket.EstimatedTime
	}
	if ticket.ActualTime != nil {
		ticketDTO.ActualTime = ticket.ActualTime
	}
	if ticket.ClosedAt != nil {
		ticketDTO.ClosedAt = ticket.ClosedAt
	}

	// Convertir CreatedBy
	if ticket.CreatedBy.ID != 0 {
		userDTO := s.userToDTO(&ticket.CreatedBy)
		ticketDTO.CreatedBy = userDTO
	}

	// Convertir AssignedTo
	if ticket.AssignedTo != nil && ticket.AssignedTo.ID != 0 {
		userDTO := s.userToDTO(ticket.AssignedTo)
		ticketDTO.AssignedTo = &userDTO
	}

	return ticketDTO
}

// userToDTO convertit un modèle User en DTO (méthode helper)
func (s *timeEntryService) userToDTO(user *models.User) dto.UserDTO {
	userDTO := dto.UserDTO{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Avatar:    user.Avatar,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	// Vérifier que le rôle a été chargé
	if user.RoleID != 0 && user.Role.ID != 0 {
		userDTO.Role = user.Role.Name
	}

	if user.LastLogin != nil {
		userDTO.LastLogin = user.LastLogin
	}

	return userDTO
}
