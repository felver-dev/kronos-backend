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
		Priority:        req.Priority,
		TargetTime:     req.TargetTime,
		Unit:           unit,
		IsActive:       isActive,
		CreatedByID:     createdByIDPtr,
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

