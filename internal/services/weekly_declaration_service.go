package services

import (
	"errors"
	"time"

	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/models"
	"github.com/mcicare/itsm-backend/internal/repositories"
)

// WeeklyDeclarationService interface pour les opérations sur les déclarations hebdomadaires
type WeeklyDeclarationService interface {
	GetByID(id uint) (*dto.WeeklyDeclarationDTO, error)
	GetByUserIDAndWeek(userID uint, week string) (*dto.WeeklyDeclarationDTO, error)
	GetByUserID(userID uint) ([]dto.WeeklyDeclarationDTO, error)
	GetValidated() ([]dto.WeeklyDeclarationDTO, error)
	GetPendingValidation() ([]dto.WeeklyDeclarationDTO, error)
	Validate(id uint, validatedByID uint) (*dto.WeeklyDeclarationDTO, error)
	Delete(id uint) error
}

// weeklyDeclarationService implémente WeeklyDeclarationService
type weeklyDeclarationService struct {
	declarationRepo repositories.WeeklyDeclarationRepository
	userRepo        repositories.UserRepository
}

// NewWeeklyDeclarationService crée une nouvelle instance de WeeklyDeclarationService
func NewWeeklyDeclarationService(
	declarationRepo repositories.WeeklyDeclarationRepository,
	userRepo repositories.UserRepository,
) WeeklyDeclarationService {
	return &weeklyDeclarationService{
		declarationRepo: declarationRepo,
		userRepo:        userRepo,
	}
}

// GetByID récupère une déclaration par son ID
func (s *weeklyDeclarationService) GetByID(id uint) (*dto.WeeklyDeclarationDTO, error) {
	declaration, err := s.declarationRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("déclaration introuvable")
	}

	declarationDTO := s.declarationToDTO(declaration)
	return &declarationDTO, nil
}

// GetByUserIDAndWeek récupère une déclaration par utilisateur et semaine
func (s *weeklyDeclarationService) GetByUserIDAndWeek(userID uint, week string) (*dto.WeeklyDeclarationDTO, error) {
	declaration, err := s.declarationRepo.FindByUserIDAndWeek(userID, week)
	if err != nil {
		return nil, errors.New("déclaration introuvable")
	}

	declarationDTO := s.declarationToDTO(declaration)
	return &declarationDTO, nil
}

// GetByUserID récupère les déclarations d'un utilisateur
func (s *weeklyDeclarationService) GetByUserID(userID uint) ([]dto.WeeklyDeclarationDTO, error) {
	declarations, err := s.declarationRepo.FindByUserID(userID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des déclarations")
	}

	var declarationDTOs []dto.WeeklyDeclarationDTO
	for _, declaration := range declarations {
		declarationDTOs = append(declarationDTOs, s.declarationToDTO(&declaration))
	}

	return declarationDTOs, nil
}

// GetValidated récupère les déclarations validées
func (s *weeklyDeclarationService) GetValidated() ([]dto.WeeklyDeclarationDTO, error) {
	declarations, err := s.declarationRepo.FindValidated()
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des déclarations")
	}

	var declarationDTOs []dto.WeeklyDeclarationDTO
	for _, declaration := range declarations {
		declarationDTOs = append(declarationDTOs, s.declarationToDTO(&declaration))
	}

	return declarationDTOs, nil
}

// GetPendingValidation récupère les déclarations en attente de validation
func (s *weeklyDeclarationService) GetPendingValidation() ([]dto.WeeklyDeclarationDTO, error) {
	declarations, err := s.declarationRepo.FindPendingValidation()
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des déclarations")
	}

	var declarationDTOs []dto.WeeklyDeclarationDTO
	for _, declaration := range declarations {
		declarationDTOs = append(declarationDTOs, s.declarationToDTO(&declaration))
	}

	return declarationDTOs, nil
}

// Validate valide une déclaration
func (s *weeklyDeclarationService) Validate(id uint, validatedByID uint) (*dto.WeeklyDeclarationDTO, error) {
	declaration, err := s.declarationRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("déclaration introuvable")
	}

	// Vérifier que le validateur existe
	_, err = s.userRepo.FindByID(validatedByID)
	if err != nil {
		return nil, errors.New("validateur introuvable")
	}

	now := time.Now()
	declaration.Validated = true
	declaration.ValidatedByID = &validatedByID
	declaration.ValidatedAt = &now

	if err := s.declarationRepo.Update(declaration); err != nil {
		return nil, errors.New("erreur lors de la validation de la déclaration")
	}

	// Récupérer la déclaration mise à jour
	updatedDeclaration, err := s.declarationRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération de la déclaration mise à jour")
	}

	declarationDTO := s.declarationToDTO(updatedDeclaration)
	return &declarationDTO, nil
}

// Delete supprime une déclaration
func (s *weeklyDeclarationService) Delete(id uint) error {
	_, err := s.declarationRepo.FindByID(id)
	if err != nil {
		return errors.New("déclaration introuvable")
	}

	if err := s.declarationRepo.Delete(id); err != nil {
		return errors.New("erreur lors de la suppression de la déclaration")
	}

	return nil
}

// declarationToDTO convertit un modèle WeeklyDeclaration en DTO
func (s *weeklyDeclarationService) declarationToDTO(declaration *models.WeeklyDeclaration) dto.WeeklyDeclarationDTO {
	declarationDTO := dto.WeeklyDeclarationDTO{
		ID:        declaration.ID,
		UserID:    declaration.UserID,
		Week:      declaration.Week,
		StartDate: declaration.StartDate,
		EndDate:   declaration.EndDate,
		TaskCount: declaration.TaskCount,
		TotalTime: declaration.TotalTime,
		Validated: declaration.Validated,
		CreatedAt: declaration.CreatedAt,
		UpdatedAt: declaration.UpdatedAt,
	}

	if declaration.ValidatedByID != nil {
		declarationDTO.ValidatedBy = declaration.ValidatedByID
	}
	if declaration.ValidatedAt != nil {
		declarationDTO.ValidatedAt = declaration.ValidatedAt
	}

	// Convertir l'utilisateur si présent
	if declaration.User.ID != 0 {
		userDTO := s.userToDTO(&declaration.User)
		declarationDTO.User = &userDTO
	}

	// Calculer la répartition par jour si les tâches sont présentes
	if len(declaration.Tasks) > 0 {
		dailyBreakdown := make(map[time.Time]dto.DailyBreakdownDTO)
		for _, task := range declaration.Tasks {
			date := task.Date
			if breakdown, exists := dailyBreakdown[date]; exists {
				breakdown.TaskCount++
				breakdown.TotalTime += task.TimeSpent
				dailyBreakdown[date] = breakdown
			} else {
				dailyBreakdown[date] = dto.DailyBreakdownDTO{
					Date:      date,
					TaskCount: 1,
					TotalTime: task.TimeSpent,
				}
			}
		}

		var breakdownDTOs []dto.DailyBreakdownDTO
		for _, breakdown := range dailyBreakdown {
			breakdownDTOs = append(breakdownDTOs, breakdown)
		}
		declarationDTO.DailyBreakdown = breakdownDTOs
	}

	return declarationDTO
}

// userToDTO convertit un modèle User en DTO (méthode helper)
func (s *weeklyDeclarationService) userToDTO(user *models.User) dto.UserDTO {
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

