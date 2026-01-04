package services

import (
	"errors"
	"time"

	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/models"
	"github.com/mcicare/itsm-backend/internal/repositories"
)

// SLAService interface pour les opérations sur les SLA
type SLAService interface {
	Create(req dto.CreateSLARequest, createdByID uint) (*dto.SLADTO, error)
	GetByID(id uint) (*dto.SLADTO, error)
	GetAll() ([]dto.SLADTO, error)
	GetActive() ([]dto.SLADTO, error)
	GetByCategory(category string) ([]dto.SLADTO, error)
	Update(id uint, req dto.UpdateSLARequest, updatedByID uint) (*dto.SLADTO, error)
	Delete(id uint) error
	GetTicketSLAStatus(ticketID uint) (*dto.TicketSLAStatusDTO, error)
	GetCompliance(slaID uint) (*dto.SLAComplianceDTO, error)
	GetViolations(slaID uint) ([]dto.SLAViolationDTO, error)
	GetAllViolations(period, category string) ([]dto.SLAViolationDTO, error)
	GetComplianceReport(period, format string) (*dto.SLAComplianceReportDTO, error)
}

// slaService implémente SLAService
type slaService struct {
	slaRepo       repositories.SLARepository
	ticketSLARepo repositories.TicketSLARepository
	ticketRepo    repositories.TicketRepository
}

// NewSLAService crée une nouvelle instance de SLAService
func NewSLAService(
	slaRepo repositories.SLARepository,
	ticketSLARepo repositories.TicketSLARepository,
	ticketRepo repositories.TicketRepository,
) SLAService {
	return &slaService{
		slaRepo:       slaRepo,
		ticketSLARepo: ticketSLARepo,
		ticketRepo:    ticketRepo,
	}
}

// Create crée un nouveau SLA
func (s *slaService) Create(req dto.CreateSLARequest, createdByID uint) (*dto.SLADTO, error) {
	// Définir l'unité par défaut
	unit := req.Unit
	if unit == "" {
		unit = "minutes"
	}

	// Définir le statut actif par défaut
	isActive := req.IsActive
	if !req.IsActive && req.IsActive == false {
		isActive = true
	}

	// Créer le SLA
	createdByIDPtr := &createdByID
	sla := &models.SLA{
		Name:           req.Name,
		Description:    req.Description,
		TicketCategory: req.TicketCategory,
		Priority:       req.Priority,
		TargetTime:     req.TargetTime,
		Unit:           unit,
		IsActive:       isActive,
		CreatedByID:    createdByIDPtr,
	}

	if err := s.slaRepo.Create(sla); err != nil {
		return nil, errors.New("erreur lors de la création du SLA")
	}

	// Récupérer le SLA créé
	createdSLA, err := s.slaRepo.FindByID(sla.ID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération du SLA créé")
	}

	slaDTO := s.slaToDTO(createdSLA)
	return &slaDTO, nil
}

// GetByID récupère un SLA par son ID
func (s *slaService) GetByID(id uint) (*dto.SLADTO, error) {
	sla, err := s.slaRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("SLA introuvable")
	}

	slaDTO := s.slaToDTO(sla)
	return &slaDTO, nil
}

// GetAll récupère tous les SLA
func (s *slaService) GetAll() ([]dto.SLADTO, error) {
	slas, err := s.slaRepo.FindAll()
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des SLA")
	}

	var slaDTOs []dto.SLADTO
	for _, sla := range slas {
		slaDTOs = append(slaDTOs, s.slaToDTO(&sla))
	}

	return slaDTOs, nil
}

// GetActive récupère les SLA actifs
func (s *slaService) GetActive() ([]dto.SLADTO, error) {
	slas, err := s.slaRepo.FindActive()
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des SLA")
	}

	var slaDTOs []dto.SLADTO
	for _, sla := range slas {
		slaDTOs = append(slaDTOs, s.slaToDTO(&sla))
	}

	return slaDTOs, nil
}

// GetByCategory récupère les SLA d'une catégorie
func (s *slaService) GetByCategory(category string) ([]dto.SLADTO, error) {
	slas, err := s.slaRepo.FindByCategory(category)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des SLA")
	}

	var slaDTOs []dto.SLADTO
	for _, sla := range slas {
		slaDTOs = append(slaDTOs, s.slaToDTO(&sla))
	}

	return slaDTOs, nil
}

// Update met à jour un SLA
func (s *slaService) Update(id uint, req dto.UpdateSLARequest, updatedByID uint) (*dto.SLADTO, error) {
	sla, err := s.slaRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("SLA introuvable")
	}

	// Mettre à jour les champs fournis
	if req.Name != "" {
		sla.Name = req.Name
	}
	if req.Description != "" {
		sla.Description = req.Description
	}
	if req.TargetTime != nil {
		sla.TargetTime = *req.TargetTime
	}
	if req.Unit != "" {
		sla.Unit = req.Unit
	}
	if req.IsActive != nil {
		sla.IsActive = *req.IsActive
	}

	if err := s.slaRepo.Update(sla); err != nil {
		return nil, errors.New("erreur lors de la mise à jour du SLA")
	}

	// Récupérer le SLA mis à jour
	updatedSLA, err := s.slaRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération du SLA mis à jour")
	}

	slaDTO := s.slaToDTO(updatedSLA)
	return &slaDTO, nil
}

// Delete supprime un SLA
func (s *slaService) Delete(id uint) error {
	_, err := s.slaRepo.FindByID(id)
	if err != nil {
		return errors.New("SLA introuvable")
	}

	if err := s.slaRepo.Delete(id); err != nil {
		return errors.New("erreur lors de la suppression du SLA")
	}

	return nil
}

// GetTicketSLAStatus récupère le statut SLA d'un ticket
func (s *slaService) GetTicketSLAStatus(ticketID uint) (*dto.TicketSLAStatusDTO, error) {
	// Récupérer le ticket
	ticket, err := s.ticketRepo.FindByID(ticketID)
	if err != nil {
		return nil, errors.New("ticket introuvable")
	}

	// Récupérer l'association ticket-SLA
	ticketSLA, err := s.ticketSLARepo.FindByTicketID(ticketID)
	if err != nil {
		return nil, errors.New("aucun SLA associé à ce ticket")
	}

	// Calculer le temps écoulé
	elapsedTime := int(time.Since(ticket.CreatedAt).Minutes())
	remaining := int(time.Until(ticketSLA.TargetTime).Minutes())

	// Déterminer le statut
	status := "on_time"
	var violatedAt *time.Time
	if remaining < 0 {
		status = "violated"
		violatedAt = &ticketSLA.TargetTime
	} else if remaining < (elapsedTime / 4) {
		status = "at_risk"
	}

	slaDTO := s.slaToDTO(&ticketSLA.SLA)

	return &dto.TicketSLAStatusDTO{
		SLAID:       ticketSLA.SLAID,
		SLA:         &slaDTO,
		TargetTime:  ticketSLA.TargetTime,
		ElapsedTime: elapsedTime,
		Remaining:   remaining,
		Status:      status,
		ViolatedAt:  violatedAt,
	}, nil
}

// GetCompliance récupère le taux de conformité d'un SLA
func (s *slaService) GetCompliance(slaID uint) (*dto.SLAComplianceDTO, error) {
	sla, err := s.slaRepo.FindByID(slaID)
	if err != nil {
		return nil, errors.New("SLA introuvable")
	}

	// Récupérer tous les tickets SLA associés à ce SLA
	ticketSLAs, err := s.ticketSLARepo.FindBySLAID(slaID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des tickets SLA")
	}

	totalTickets := len(ticketSLAs)
	compliant := 0
	violations := 0

	for _, tsla := range ticketSLAs {
		if tsla.Status == "on_time" {
			compliant++
		} else if tsla.Status == "violated" {
			violations++
		}
	}

	complianceRate := 0.0
	if totalTickets > 0 {
		complianceRate = (float64(compliant) / float64(totalTickets)) * 100
	}

	slaDTO := s.slaToDTO(sla)
	return &dto.SLAComplianceDTO{
		SLAID:          slaID,
		SLA:            &slaDTO,
		ComplianceRate: complianceRate,
		TotalTickets:   totalTickets,
		Compliant:      compliant,
		Violations:     violations,
	}, nil
}

// GetViolations récupère les violations d'un SLA
func (s *slaService) GetViolations(slaID uint) ([]dto.SLAViolationDTO, error) {
	// Récupérer tous les tickets SLA violés pour ce SLA
	ticketSLAs, err := s.ticketSLARepo.FindBySLAID(slaID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des violations")
	}

	var violations []dto.SLAViolationDTO
	for _, tsla := range ticketSLAs {
		if tsla.Status == "violated" && tsla.ViolationTime != nil {
			ticket, _ := s.ticketRepo.FindByID(tsla.TicketID)
			var ticketDTO *dto.TicketDTO
			if ticket != nil {
				// Convertir le ticket en DTO (simplifié)
				ticketDTO = &dto.TicketDTO{
					ID:          ticket.ID,
					Title:       ticket.Title,
					Category:    ticket.Category,
					Status:      ticket.Status,
					CreatedAt:   ticket.CreatedAt,
				}
			}

			slaDTO := s.slaToDTO(&tsla.SLA)
			violations = append(violations, dto.SLAViolationDTO{
				ID:            tsla.ID,
				TicketID:      tsla.TicketID,
				Ticket:        ticketDTO,
				SLAID:         slaID,
				SLA:           &slaDTO,
				ViolationTime: *tsla.ViolationTime,
				Unit:          "minutes",
				ViolatedAt:    tsla.TargetTime,
			})
		}
	}

	return violations, nil
}

// GetAllViolations récupère toutes les violations de SLA
func (s *slaService) GetAllViolations(period, category string) ([]dto.SLAViolationDTO, error) {
	// Récupérer tous les tickets SLA violés
	allTicketSLAs, err := s.ticketSLARepo.FindAll()
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des violations")
	}

	var violations []dto.SLAViolationDTO
	now := time.Now()
	var startDate time.Time

	// Calculer la date de début selon la période
	switch period {
	case "week":
		startDate = now.AddDate(0, 0, -7)
	case "month":
		startDate = now.AddDate(0, -1, 0)
	default:
		startDate = now.AddDate(0, -1, 0) // Par défaut: 1 mois
	}

	for _, tsla := range allTicketSLAs {
		// Filtrer par statut violé
		if tsla.Status == "violated" && tsla.ViolationTime != nil {
			// Filtrer par période
			if tsla.TargetTime.Before(startDate) {
				continue
			}

			// Filtrer par catégorie si fournie
			if category != "" && tsla.SLA.TicketCategory != category {
				continue
			}

			ticket, _ := s.ticketRepo.FindByID(tsla.TicketID)
			var ticketDTO *dto.TicketDTO
			if ticket != nil {
				ticketDTO = &dto.TicketDTO{
					ID:          ticket.ID,
					Title:       ticket.Title,
					Category:    ticket.Category,
					Status:      ticket.Status,
					CreatedAt:   ticket.CreatedAt,
				}
			}

			slaDTO := s.slaToDTO(&tsla.SLA)
			violations = append(violations, dto.SLAViolationDTO{
				ID:            tsla.ID,
				TicketID:      tsla.TicketID,
				Ticket:        ticketDTO,
				SLAID:         tsla.SLAID,
				SLA:           &slaDTO,
				ViolationTime: *tsla.ViolationTime,
				Unit:          "minutes",
				ViolatedAt:    tsla.TargetTime,
			})
		}
	}

	return violations, nil
}

// GetComplianceReport génère un rapport de conformité
func (s *slaService) GetComplianceReport(period, format string) (*dto.SLAComplianceReportDTO, error) {
	// Récupérer tous les SLA actifs
	slas, err := s.slaRepo.FindActive()
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des SLA")
	}

	overallCompliance := 0.0
	totalTickets := 0
	totalViolations := 0
	byCategory := make(map[string]float64)
	byPriority := make(map[string]float64)

	for _, sla := range slas {
		compliance, err := s.GetCompliance(sla.ID)
		if err != nil {
			continue
		}

		totalTickets += compliance.TotalTickets
		totalViolations += compliance.Violations

		// Calculer la conformité par catégorie
		if compliance.TotalTickets > 0 {
			categoryCompliance := (float64(compliance.Compliant) / float64(compliance.TotalTickets)) * 100
			byCategory[sla.TicketCategory] = categoryCompliance

			// Calculer par priorité si fournie
			if sla.Priority != nil {
				byPriority[*sla.Priority] = categoryCompliance
			}
		}
	}

	if totalTickets > 0 {
		overallCompliance = (float64(totalTickets-totalViolations) / float64(totalTickets)) * 100
	}

	return &dto.SLAComplianceReportDTO{
		OverallCompliance: overallCompliance,
		ByCategory:        byCategory,
		ByPriority:         byPriority,
		TotalTickets:      totalTickets,
		TotalViolations:   totalViolations,
		Period:            period,
		GeneratedAt:       time.Now(),
	}, nil
}

// slaToDTO convertit un modèle SLA en DTO
func (s *slaService) slaToDTO(sla *models.SLA) dto.SLADTO {
	return dto.SLADTO{
		ID:             sla.ID,
		Name:           sla.Name,
		Description:    sla.Description,
		TicketCategory: sla.TicketCategory,
		Priority:       sla.Priority,
		TargetTime:     sla.TargetTime,
		Unit:           sla.Unit,
		IsActive:       sla.IsActive,
		CreatedAt:      sla.CreatedAt,
		UpdatedAt:      sla.UpdatedAt,
	}
}
