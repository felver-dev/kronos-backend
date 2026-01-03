package services

import (
	"errors"

	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/models"
	"github.com/mcicare/itsm-backend/internal/repositories"
)

// PerformanceService interface pour les opérations sur les performances
type PerformanceService interface {
	GetPerformanceByUserID(userID uint) (*dto.PerformanceDTO, error)
	GetEfficiencyByUserID(userID uint) (*dto.EfficiencyDTO, error)
	GetProductivityByUserID(userID uint) (*dto.ProductivityDTO, error)
	GetBudgetComplianceByUserID(userID uint) (*dto.BudgetComplianceDTO, error)
	GetWorkloadByUserID(userID uint) (*dto.WorkloadDTO, error)
	GetPerformanceRanking(limit int) ([]dto.PerformanceRankingDTO, error)
	GetTeamComparison() (*dto.TeamComparisonDTO, error)
}

// performanceService implémente PerformanceService
type performanceService struct {
	ticketRepo    repositories.TicketRepository
	timeEntryRepo repositories.TimeEntryRepository
	delayRepo     repositories.DelayRepository
	userRepo      repositories.UserRepository
}

// NewPerformanceService crée une nouvelle instance de PerformanceService
func NewPerformanceService(
	ticketRepo repositories.TicketRepository,
	timeEntryRepo repositories.TimeEntryRepository,
	delayRepo repositories.DelayRepository,
	userRepo repositories.UserRepository,
) PerformanceService {
	return &performanceService{
		ticketRepo:    ticketRepo,
		timeEntryRepo: timeEntryRepo,
		delayRepo:     delayRepo,
		userRepo:      userRepo,
	}
}

// GetPerformanceByUserID récupère les métriques de performance d'un utilisateur
func (s *performanceService) GetPerformanceByUserID(userID uint) (*dto.PerformanceDTO, error) {
	// Vérifier que l'utilisateur existe
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("utilisateur introuvable")
	}

	// TODO: Implémenter les calculs de performance
	// Pour l'instant, retourner une structure vide
	userDTO := s.userToDTO(user)
	return &dto.PerformanceDTO{
		UserID: userID,
		User:   &userDTO,
	}, nil
}

// GetEfficiencyByUserID récupère les métriques d'efficacité d'un utilisateur
func (s *performanceService) GetEfficiencyByUserID(userID uint) (*dto.EfficiencyDTO, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("utilisateur introuvable")
	}

	userDTO := s.userToDTO(user)
	return &dto.EfficiencyDTO{
		UserID: userID,
		User:   &userDTO,
	}, nil
}

// GetProductivityByUserID récupère les métriques de productivité d'un utilisateur
func (s *performanceService) GetProductivityByUserID(userID uint) (*dto.ProductivityDTO, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("utilisateur introuvable")
	}

	userDTO := s.userToDTO(user)
	return &dto.ProductivityDTO{
		UserID: userID,
		User:   &userDTO,
	}, nil
}

// GetBudgetComplianceByUserID récupère les métriques de respect des budgets d'un utilisateur
func (s *performanceService) GetBudgetComplianceByUserID(userID uint) (*dto.BudgetComplianceDTO, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("utilisateur introuvable")
	}

	userDTO := s.userToDTO(user)
	return &dto.BudgetComplianceDTO{
		UserID: userID,
		User:   &userDTO,
	}, nil
}

// GetWorkloadByUserID récupère la charge de travail d'un utilisateur
func (s *performanceService) GetWorkloadByUserID(userID uint) (*dto.WorkloadDTO, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("utilisateur introuvable")
	}

	userDTO := s.userToDTO(user)
	return &dto.WorkloadDTO{
		UserID: userID,
		User:   &userDTO,
	}, nil
}

// GetPerformanceRanking récupère le classement des performances
func (s *performanceService) GetPerformanceRanking(limit int) ([]dto.PerformanceRankingDTO, error) {
	// TODO: Implémenter le calcul du classement
	return []dto.PerformanceRankingDTO{}, nil
}

// GetTeamComparison récupère la comparaison de l'équipe
func (s *performanceService) GetTeamComparison() (*dto.TeamComparisonDTO, error) {
	// TODO: Implémenter la comparaison d'équipe
	return &dto.TeamComparisonDTO{}, nil
}

// userToDTO convertit un modèle User en DTO (méthode helper)
func (s *performanceService) userToDTO(user *models.User) dto.UserDTO {
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

	if user.RoleID != 0 {
		userDTO.Role = user.Role.Name
	}

	if user.LastLogin != nil {
		userDTO.LastLogin = user.LastLogin
	}

	return userDTO
}

