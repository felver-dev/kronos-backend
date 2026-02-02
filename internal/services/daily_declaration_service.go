package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/models"
	"github.com/mcicare/itsm-backend/internal/repositories"
)

// DailyDeclarationService interface pour les opérations sur les déclarations journalières
type DailyDeclarationService interface {
	GetByID(id uint) (*dto.DailyDeclarationDTO, error)
	GetByUserIDAndDate(userID uint, date time.Time) (*dto.DailyDeclarationDTO, error)
	GetByUserID(userID uint) ([]dto.DailyDeclarationDTO, error)
	GetByDateRange(userID uint, startDate, endDate time.Time) ([]dto.DailyDeclarationDTO, error)
	GetAllByDateRange(startDate, endDate time.Time) ([]dto.DailyDeclarationDTO, error) // Pour les admins
	GetValidated() ([]dto.DailyDeclarationDTO, error)
	GetPendingValidation() ([]dto.DailyDeclarationDTO, error)
	Validate(id uint, validatedByID uint) (*dto.DailyDeclarationDTO, error)
	Delete(id uint) error
}

// dailyDeclarationService implémente DailyDeclarationService
type dailyDeclarationService struct {
	declarationRepo repositories.DailyDeclarationRepository
	timeEntryRepo   repositories.TimeEntryRepository
	userRepo        repositories.UserRepository
}

// NewDailyDeclarationService crée une nouvelle instance de DailyDeclarationService
func NewDailyDeclarationService(
	declarationRepo repositories.DailyDeclarationRepository,
	timeEntryRepo repositories.TimeEntryRepository,
	userRepo repositories.UserRepository,
) DailyDeclarationService {
	return &dailyDeclarationService{
		declarationRepo: declarationRepo,
		timeEntryRepo:   timeEntryRepo,
		userRepo:        userRepo,
	}
}

// GetByID récupère une déclaration par son ID
func (s *dailyDeclarationService) GetByID(id uint) (result *dto.DailyDeclarationDTO, err error) {
	// Protéger contre les panics
	defer func() {
		if r := recover(); r != nil {
			result = nil
			err = fmt.Errorf("erreur lors de la conversion de la déclaration: %v", r)
		}
	}()
	
	declaration, err := s.declarationRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("déclaration introuvable: %v", err)
	}

	declarationDTO := s.declarationToDTO(declaration)
	return &declarationDTO, nil
}

// GetByUserIDAndDate récupère une déclaration par utilisateur et date
func (s *dailyDeclarationService) GetByUserIDAndDate(userID uint, date time.Time) (*dto.DailyDeclarationDTO, error) {
	declaration, err := s.declarationRepo.FindByUserIDAndDate(userID, date)
	if err != nil {
		return nil, errors.New("déclaration introuvable")
	}

	declarationDTO := s.declarationToDTO(declaration)
	return &declarationDTO, nil
}

// GetByUserID récupère les déclarations d'un utilisateur
func (s *dailyDeclarationService) GetByUserID(userID uint) ([]dto.DailyDeclarationDTO, error) {
	declarations, err := s.declarationRepo.FindByUserID(userID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des déclarations")
	}

	var declarationDTOs []dto.DailyDeclarationDTO
	for _, declaration := range declarations {
		declarationDTOs = append(declarationDTOs, s.declarationToDTO(&declaration))
	}

	return declarationDTOs, nil
}

// GetByDateRange récupère les déclarations d'un utilisateur dans une plage de dates
func (s *dailyDeclarationService) GetByDateRange(userID uint, startDate, endDate time.Time) (result []dto.DailyDeclarationDTO, err error) {
	// Protéger contre les panics
	defer func() {
		if r := recover(); r != nil {
			// En cas de panic, retourner un tableau vide
			result = []dto.DailyDeclarationDTO{}
			err = nil
		}
	}()
	
	declarations, err := s.declarationRepo.FindByDateRange(userID, startDate, endDate)
	if err != nil {
		// Retourner un tableau vide au lieu d'une erreur pour éviter de faire planter la page
		return []dto.DailyDeclarationDTO{}, nil
	}

	var declarationDTOs []dto.DailyDeclarationDTO
	for i := range declarations {
		// Gérer les erreurs de conversion gracieusement
		func() {
			defer func() {
				if r := recover(); r != nil {
					// Ignorer les erreurs de conversion et continuer avec les autres déclarations
				}
			}()
			dto := s.declarationToDTO(&declarations[i])
			declarationDTOs = append(declarationDTOs, dto)
		}()
	}

	return declarationDTOs, nil
}

// GetAllByDateRange récupère toutes les déclarations dans une plage de dates (pour les admins)
func (s *dailyDeclarationService) GetAllByDateRange(startDate, endDate time.Time) (result []dto.DailyDeclarationDTO, err error) {
	// Protéger toute la fonction avec recover pour éviter les panics
	defer func() {
		if r := recover(); r != nil {
			// En cas de panic, retourner un tableau vide
			result = []dto.DailyDeclarationDTO{}
			err = nil
		}
	}()
	
	// Toujours retourner un tableau vide en cas d'erreur - ne jamais retourner d'erreur
	declarations, repoErr := s.declarationRepo.FindAllByDateRange(startDate, endDate)
	if repoErr != nil {
		// Retourner un tableau vide au lieu d'une erreur pour éviter de faire planter la page
		return []dto.DailyDeclarationDTO{}, nil
	}

	// S'assurer qu'on ne retourne jamais nil
	if declarations == nil {
		return []dto.DailyDeclarationDTO{}, nil
	}

	// Si aucune déclaration, retourner un tableau vide
	if len(declarations) == 0 {
		return []dto.DailyDeclarationDTO{}, nil
	}

	var declarationDTOs []dto.DailyDeclarationDTO
	for i := range declarations {
		// Gérer les erreurs de conversion gracieusement - ignorer les déclarations qui causent des erreurs
		func() {
			defer func() {
				if r := recover(); r != nil {
					// Ignorer les erreurs de conversion et continuer avec les autres déclarations
				}
			}()
			dto := s.declarationToDTO(&declarations[i])
			declarationDTOs = append(declarationDTOs, dto)
		}()
	}

	// S'assurer qu'on ne retourne jamais nil
	if declarationDTOs == nil {
		return []dto.DailyDeclarationDTO{}, nil
	}

	// Toujours retourner un résultat, même si c'est un tableau vide
	return declarationDTOs, nil
}

// GetValidated récupère les déclarations validées
func (s *dailyDeclarationService) GetValidated() ([]dto.DailyDeclarationDTO, error) {
	declarations, err := s.declarationRepo.FindValidated()
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des déclarations")
	}

	var declarationDTOs []dto.DailyDeclarationDTO
	for _, declaration := range declarations {
		declarationDTOs = append(declarationDTOs, s.declarationToDTO(&declaration))
	}

	return declarationDTOs, nil
}

// GetPendingValidation récupère les déclarations en attente de validation
func (s *dailyDeclarationService) GetPendingValidation() ([]dto.DailyDeclarationDTO, error) {
	declarations, err := s.declarationRepo.FindPendingValidation()
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des déclarations")
	}

	var declarationDTOs []dto.DailyDeclarationDTO
	for _, declaration := range declarations {
		declarationDTOs = append(declarationDTOs, s.declarationToDTO(&declaration))
	}

	return declarationDTOs, nil
}

// Validate valide une déclaration
func (s *dailyDeclarationService) Validate(id uint, validatedByID uint) (*dto.DailyDeclarationDTO, error) {
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
func (s *dailyDeclarationService) Delete(id uint) error {
	_, err := s.declarationRepo.FindByID(id)
	if err != nil {
		return errors.New("déclaration introuvable")
	}

	if err := s.declarationRepo.Delete(id); err != nil {
		return errors.New("erreur lors de la suppression de la déclaration")
	}

	return nil
}

// declarationToDTO convertit un modèle DailyDeclaration en DTO
func (s *dailyDeclarationService) declarationToDTO(declaration *models.DailyDeclaration) dto.DailyDeclarationDTO {
	declarationDTO := dto.DailyDeclarationDTO{
		ID:        declaration.ID,
		UserID:    declaration.UserID,
		Date:      declaration.Date,
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

	// Convertir l'utilisateur si présent - avec gestion d'erreur
	if declaration.User.ID != 0 {
		func() {
			defer func() {
				if r := recover(); r != nil {
					// Ignorer l'erreur de conversion de l'utilisateur
				}
			}()
			userDTO := s.userToDTO(&declaration.User)
			declarationDTO.User = &userDTO
		}()
	}

	// Convertir les tâches si présentes
	if len(declaration.Tasks) > 0 {
		var taskDTOs []dto.TimeEntryDTO
		for i := range declaration.Tasks {
			// Gérer les erreurs de conversion gracieusement
			func() {
				defer func() {
					if r := recover(); r != nil {
						// Ignorer les erreurs de conversion et continuer avec les autres tâches
					}
				}()
				taskDTO := s.timeEntryToDTO(&declaration.Tasks[i])
				taskDTOs = append(taskDTOs, taskDTO)
			}()
		}
		declarationDTO.Tasks = taskDTOs
	}

	return declarationDTO
}

// timeEntryToDTO convertit un modèle DailyDeclarationTask en TimeEntryDTO
func (s *dailyDeclarationService) timeEntryToDTO(task *models.DailyDeclarationTask) dto.TimeEntryDTO {
	// Valeurs par défaut sécurisées
	userID := uint(0)
	date := time.Now()
	validated := false
	
	// Essayer d'accéder à la déclaration seulement si elle est chargée
	func() {
		defer func() {
			if r := recover(); r != nil {
				// Si erreur, utiliser les valeurs par défaut
			}
		}()
		if task.Declaration.ID != 0 {
			userID = task.Declaration.UserID
			date = task.Declaration.Date
			validated = task.Declaration.Validated
		}
	}()
	
	taskDTO := dto.TimeEntryDTO{
		ID:        task.ID,
		TicketID:  task.TicketID,
		UserID:    userID,
		TimeSpent: task.TimeSpent,
		Date:      date,
		Validated: validated,
		CreatedAt: task.CreatedAt,
		UpdatedAt: task.CreatedAt,
	}

	// Convertir le ticket si présent et chargé - avec gestion d'erreur
	func() {
		defer func() {
			if r := recover(); r != nil {
				// Ignorer l'erreur de conversion du ticket
			}
		}()
		if task.Ticket.ID != 0 {
			ticketDTO := s.ticketToDTO(&task.Ticket)
			taskDTO.Ticket = &ticketDTO
		}
	}()

	return taskDTO
}

// ticketToDTO convertit un modèle Ticket en DTO (méthode helper)
func (s *dailyDeclarationService) ticketToDTO(ticket *models.Ticket) dto.TicketDTO {
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

	if ticket.CreatedBy.ID != 0 {
		userDTO := s.userToDTO(&ticket.CreatedBy)
		ticketDTO.CreatedBy = userDTO
	}

	if ticket.AssignedTo != nil && ticket.AssignedTo.ID != 0 {
		userDTO := s.userToDTO(ticket.AssignedTo)
		ticketDTO.AssignedTo = &userDTO
	}

	return ticketDTO
}

// userToDTO convertit un modèle User en DTO (méthode helper)
func (s *dailyDeclarationService) userToDTO(user *models.User) dto.UserDTO {
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

	// Ajouter le département si présent
	if user.DepartmentID != nil {
		userDTO.DepartmentID = user.DepartmentID
		if user.Department != nil && user.Department.ID != 0 {
			departmentDTO := dto.DepartmentDTO{
				ID:      user.Department.ID,
				Name:    user.Department.Name,
				Code:    user.Department.Code,
				IsActive: user.Department.IsActive,
			}
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
