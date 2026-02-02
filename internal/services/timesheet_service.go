package services

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/models"
	"github.com/mcicare/itsm-backend/internal/repositories"
	"github.com/mcicare/itsm-backend/database"
)

// TimesheetService interface pour les opérations sur les timesheets
type TimesheetService interface {
	// Saisie du temps par ticket
	CreateTimeEntry(req dto.CreateTimeEntryRequest, userID uint) (*dto.TimeEntryDTO, error)
	GetTimeEntries(scope interface{}) ([]dto.TimeEntryDTO, error) // scope peut être *scope.QueryScope ou nil
	GetTimeEntryByID(id uint) (*dto.TimeEntryDTO, error)
	UpdateTimeEntry(id uint, req dto.UpdateTimeEntryRequest, userID uint) (*dto.TimeEntryDTO, error)
	GetTimeEntriesByTicketID(ticketID uint) ([]dto.TimeEntryDTO, error)
	GetTimeEntriesByUserID(userID uint) ([]dto.TimeEntryDTO, error)
	GetTimeEntriesByDate(date time.Time, userID uint) ([]dto.TimeEntryDTO, error)

	// Déclaration par jour
	GetDailyDeclaration(date time.Time, userID uint) (*dto.DailyDeclarationDTO, error)
	CreateOrUpdateDailyDeclaration(date time.Time, userID uint, tasks []dto.DailyTaskRequest) (*dto.DailyDeclarationDTO, error)
	GetDailyTasks(date time.Time, userID uint) ([]dto.DailyTaskDTO, error)
	CreateDailyTask(date time.Time, userID uint, task dto.DailyTaskRequest) (*dto.DailyTaskDTO, error)
	DeleteDailyTask(date time.Time, userID uint, taskID uint) error
	GetDailySummary(date time.Time, userID uint) (*dto.DailySummaryDTO, error)
	GetDailyCalendar(userID uint, startDate, endDate time.Time) ([]dto.DailyCalendarDTO, error)
	GetDailyRange(userID uint, startDate, endDate time.Time) ([]dto.DailyDeclarationDTO, error)
	GetAllDailyRange(startDate, endDate time.Time) ([]dto.DailyDeclarationDTO, error) // Pour les admins

	// Déclaration par semaine
	GetWeeklyDeclaration(week string, userID uint) (*dto.WeeklyDeclarationDTO, error)
	CreateOrUpdateWeeklyDeclaration(week string, userID uint, tasks []dto.WeeklyTaskRequest) (*dto.WeeklyDeclarationDTO, error)
	GetWeeklyTasks(week string, userID uint) ([]dto.WeeklyTaskDTO, error)
	GetWeeklySummary(week string, userID uint) (*dto.WeeklySummaryDTO, error)
	GetWeeklyDailyBreakdown(week string, userID uint) ([]dto.DailyBreakdownDTO, error)
	ValidateWeeklyDeclaration(week string, userID uint, validatedByID uint) (*dto.WeeklyDeclarationDTO, error)
	GetWeeklyValidationStatus(week string, userID uint) (*dto.ValidationStatusDTO, error)

	// Budget temps
	SetTicketEstimatedTime(ticketID uint, estimatedTime int, userID uint) error
	GetTicketEstimatedTime(ticketID uint) (*dto.EstimatedTimeDTO, error)
	UpdateTicketEstimatedTime(ticketID uint, estimatedTime int, userID uint) error
	GetTicketTimeComparison(ticketID uint) (*dto.TimeComparisonDTO, error)
	GetProjectTimeBudget(projectID uint) (*dto.ProjectTimeBudgetDTO, error)
	SetProjectTimeBudget(projectID uint, budget dto.SetProjectTimeBudgetRequest, userID uint) error
	GetBudgetAlerts() ([]dto.BudgetAlertDTO, error)
	GetTicketBudgetStatus(ticketID uint) (*dto.BudgetStatusDTO, error)

	// Validation
	ValidateTimeEntry(id uint, req dto.ValidateTimeEntryRequest, validatedByID uint) (*dto.TimeEntryDTO, error)
	GetPendingValidationEntries(scope interface{}) ([]dto.TimeEntryDTO, error) // scope peut être *scope.QueryScope ou nil
	GetValidationHistory() ([]dto.ValidationHistoryDTO, error)

	// Alertes
	GetDelayAlerts() ([]dto.DelayAlertDTO, error)
	GetBudgetAlertsForTimesheet() ([]dto.BudgetAlertDTO, error)
	GetOverloadAlerts() ([]dto.OverloadAlertDTO, error)
	GetUnderloadAlerts() ([]dto.UnderloadAlertDTO, error)
	SendReminderAlerts(userIDs []uint) error
	GetPendingJustificationAlerts() ([]dto.PendingJustificationAlertDTO, error)

	// Historique
	GetTimesheetHistory(userID uint, startDate, endDate time.Time) ([]dto.TimesheetHistoryDTO, error)
	GetTimesheetHistoryEntry(entryID uint) (*dto.TimesheetHistoryEntryDTO, error)
	GetTimesheetAuditTrail(userID uint, startDate, endDate time.Time) ([]dto.AuditTrailDTO, error)
	GetTimesheetModifications(userID uint, startDate, endDate time.Time) ([]dto.ModificationDTO, error)
}

// timesheetService implémente TimesheetService
type timesheetService struct {
	timeEntryService         TimeEntryService
	dailyDeclarationService  DailyDeclarationService
	weeklyDeclarationService WeeklyDeclarationService
	ticketRepo               repositories.TicketRepository
	projectRepo              repositories.ProjectRepository
	delayRepo                repositories.DelayRepository
	delayJustificationRepo   repositories.DelayJustificationRepository
	userRepo                 repositories.UserRepository
}

// NewTimesheetService crée une nouvelle instance de TimesheetService
func NewTimesheetService(
	timeEntryService TimeEntryService,
	dailyDeclarationService DailyDeclarationService,
	weeklyDeclarationService WeeklyDeclarationService,
	ticketRepo repositories.TicketRepository,
	projectRepo repositories.ProjectRepository,
	delayRepo repositories.DelayRepository,
	delayJustificationRepo repositories.DelayJustificationRepository,
	userRepo repositories.UserRepository,
) TimesheetService {
	return &timesheetService{
		timeEntryService:         timeEntryService,
		dailyDeclarationService:  dailyDeclarationService,
		weeklyDeclarationService: weeklyDeclarationService,
		ticketRepo:               ticketRepo,
		projectRepo:              projectRepo,
		delayRepo:                delayRepo,
		delayJustificationRepo:   delayJustificationRepo,
		userRepo:                 userRepo,
	}
}

// CreateTimeEntry crée une nouvelle entrée de temps
func (s *timesheetService) CreateTimeEntry(req dto.CreateTimeEntryRequest, userID uint) (*dto.TimeEntryDTO, error) {
	return s.timeEntryService.Create(req, userID)
}

// GetTimeEntries récupère toutes les entrées de temps
func (s *timesheetService) GetTimeEntries(scope interface{}) ([]dto.TimeEntryDTO, error) {
	return s.timeEntryService.GetAll(scope)
}

// GetTimeEntryByID récupère une entrée de temps par son ID
func (s *timesheetService) GetTimeEntryByID(id uint) (*dto.TimeEntryDTO, error) {
	return s.timeEntryService.GetByID(id)
}

// UpdateTimeEntry met à jour une entrée de temps
func (s *timesheetService) UpdateTimeEntry(id uint, req dto.UpdateTimeEntryRequest, userID uint) (*dto.TimeEntryDTO, error) {
	return s.timeEntryService.Update(id, req, userID)
}

// GetTimeEntriesByTicketID récupère les entrées de temps d'un ticket
func (s *timesheetService) GetTimeEntriesByTicketID(ticketID uint) ([]dto.TimeEntryDTO, error) {
	return s.timeEntryService.GetByTicketID(ticketID)
}

// GetTimeEntriesByUserID récupère les entrées de temps d'un utilisateur
func (s *timesheetService) GetTimeEntriesByUserID(userID uint) ([]dto.TimeEntryDTO, error) {
	return s.timeEntryService.GetByUserID(userID)
}

// GetTimeEntriesByDate récupère les entrées de temps d'une date
func (s *timesheetService) GetTimeEntriesByDate(date time.Time, userID uint) ([]dto.TimeEntryDTO, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)
	return s.timeEntryService.GetByDateRange(userID, startOfDay, endOfDay)
}

// GetDailyDeclaration récupère une déclaration journalière
func (s *timesheetService) GetDailyDeclaration(date time.Time, userID uint) (*dto.DailyDeclarationDTO, error) {
	return s.dailyDeclarationService.GetByUserIDAndDate(userID, date)
}

// CreateOrUpdateDailyDeclaration crée ou met à jour une déclaration journalière
func (s *timesheetService) CreateOrUpdateDailyDeclaration(date time.Time, userID uint, tasks []dto.DailyTaskRequest) (*dto.DailyDeclarationDTO, error) {
	if len(tasks) == 0 {
		return nil, errors.New("au moins une tâche est requise")
	}

	// Normaliser la date (garder seulement la date, sans l'heure)
	dateOnly := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	// Vérifier si une déclaration existe déjà pour cette date et cet utilisateur
	existingDeclaration, err := s.dailyDeclarationService.GetByUserIDAndDate(userID, dateOnly)
	
	var declaration *models.DailyDeclaration
	totalTime := 0

	// Calculer le temps total
	for _, task := range tasks {
		totalTime += task.TimeSpent
	}

	if err != nil || existingDeclaration == nil {
		// Créer une nouvelle déclaration
		declaration = &models.DailyDeclaration{
			UserID:    userID,
			Date:      dateOnly,
			TaskCount: len(tasks),
			TotalTime: totalTime,
			Validated: false,
		}

		// Créer les tâches
		declaration.Tasks = make([]models.DailyDeclarationTask, len(tasks))
		for i, task := range tasks {
			declaration.Tasks[i] = models.DailyDeclarationTask{
				TicketID:  task.TicketID,
				TimeSpent: task.TimeSpent,
			}
		}

		// Créer la déclaration avec les tâches
		if err := database.DB.Create(declaration).Error; err != nil {
			return nil, fmt.Errorf("erreur lors de la création de la déclaration: %v", err)
		}
		
		// Vérifier que l'ID a été généré
		if declaration.ID == 0 {
			return nil, errors.New("l'ID de la déclaration n'a pas été généré après la création")
		}
		
		// Recharger la déclaration avec les relations pour s'assurer que tout est à jour
		// Utiliser une requête simple sans Preloads complexes pour éviter les erreurs
		var reloadedDeclaration models.DailyDeclaration
		if err := database.DB.Preload("Tasks").First(&reloadedDeclaration, declaration.ID).Error; err != nil {
			return nil, fmt.Errorf("erreur lors du rechargement de la déclaration: %v", err)
		}
		
		// Construire le DTO manuellement pour éviter les problèmes de conversion
		declarationDTO := dto.DailyDeclarationDTO{
			ID:        reloadedDeclaration.ID,
			UserID:    reloadedDeclaration.UserID,
			Date:      reloadedDeclaration.Date,
			TaskCount: reloadedDeclaration.TaskCount,
			TotalTime: reloadedDeclaration.TotalTime,
			Validated: reloadedDeclaration.Validated,
			CreatedAt: reloadedDeclaration.CreatedAt,
			UpdatedAt: reloadedDeclaration.UpdatedAt,
		}
		
		if reloadedDeclaration.ValidatedByID != nil {
			declarationDTO.ValidatedBy = reloadedDeclaration.ValidatedByID
		}
		if reloadedDeclaration.ValidatedAt != nil {
			declarationDTO.ValidatedAt = reloadedDeclaration.ValidatedAt
		}
		
		// Convertir les tâches en TimeEntryDTO
		if len(reloadedDeclaration.Tasks) > 0 {
			taskDTOs := make([]dto.TimeEntryDTO, len(reloadedDeclaration.Tasks))
			for i, task := range reloadedDeclaration.Tasks {
				taskDTOs[i] = dto.TimeEntryDTO{
					ID:        task.ID,
					TicketID:  task.TicketID,
					UserID:    reloadedDeclaration.UserID,
					TimeSpent: task.TimeSpent,
					Date:      reloadedDeclaration.Date,
					Validated: reloadedDeclaration.Validated,
					CreatedAt: task.CreatedAt,
					// Le modèle DailyDeclarationTask n'a pas de champ UpdatedAt, on utilise UpdatedAt de la déclaration
					UpdatedAt: reloadedDeclaration.UpdatedAt,
				}
			}
			declarationDTO.Tasks = taskDTOs
		}
		
		return &declarationDTO, nil
	} else {
		// Mettre à jour la déclaration existante
		declaration = &models.DailyDeclaration{
			ID:        existingDeclaration.ID,
			UserID:    userID,
			Date:      dateOnly,
			TaskCount: len(tasks),
			TotalTime: totalTime,
			Validated: existingDeclaration.Validated,
		}

		// Supprimer les anciennes tâches
		if err := database.DB.Where("declaration_id = ?", existingDeclaration.ID).Delete(&models.DailyDeclarationTask{}).Error; err != nil {
			return nil, errors.New("erreur lors de la suppression des anciennes tâches")
		}

		// Créer les nouvelles tâches
		declaration.Tasks = make([]models.DailyDeclarationTask, len(tasks))
		for i, task := range tasks {
			declaration.Tasks[i] = models.DailyDeclarationTask{
				DeclarationID: existingDeclaration.ID,
				TicketID:      task.TicketID,
				TimeSpent:     task.TimeSpent,
			}
		}

		// Mettre à jour la déclaration
		if err := database.DB.Model(&models.DailyDeclaration{}).Where("id = ?", existingDeclaration.ID).
			Updates(map[string]interface{}{
				"task_count": len(tasks),
				"total_time": totalTime,
			}).Error; err != nil {
			return nil, errors.New("erreur lors de la mise à jour de la déclaration")
		}

		// Créer les nouvelles tâches
		if err := database.DB.Create(&declaration.Tasks).Error; err != nil {
			return nil, fmt.Errorf("erreur lors de la création des tâches: %v", err)
		}
		
		// Recharger la déclaration mise à jour avec les relations
		var reloadedDeclaration models.DailyDeclaration
		if err := database.DB.Preload("Tasks").First(&reloadedDeclaration, existingDeclaration.ID).Error; err != nil {
			return nil, fmt.Errorf("erreur lors du rechargement de la déclaration: %v", err)
		}
		
		// Construire le DTO manuellement
		declarationDTO := dto.DailyDeclarationDTO{
			ID:        reloadedDeclaration.ID,
			UserID:    reloadedDeclaration.UserID,
			Date:      reloadedDeclaration.Date,
			TaskCount: reloadedDeclaration.TaskCount,
			TotalTime: reloadedDeclaration.TotalTime,
			Validated: reloadedDeclaration.Validated,
			CreatedAt: reloadedDeclaration.CreatedAt,
			UpdatedAt: reloadedDeclaration.UpdatedAt,
		}
		
		if reloadedDeclaration.ValidatedByID != nil {
			declarationDTO.ValidatedBy = reloadedDeclaration.ValidatedByID
		}
		if reloadedDeclaration.ValidatedAt != nil {
			declarationDTO.ValidatedAt = reloadedDeclaration.ValidatedAt
		}
		
		// Convertir les tâches en TimeEntryDTO
		if len(reloadedDeclaration.Tasks) > 0 {
			taskDTOs := make([]dto.TimeEntryDTO, len(reloadedDeclaration.Tasks))
			for i, task := range reloadedDeclaration.Tasks {
				taskDTOs[i] = dto.TimeEntryDTO{
					ID:        task.ID,
					TicketID:  task.TicketID,
					UserID:    reloadedDeclaration.UserID,
					TimeSpent: task.TimeSpent,
					Date:      reloadedDeclaration.Date,
					Validated: reloadedDeclaration.Validated,
					CreatedAt: task.CreatedAt,
					// Pas de UpdatedAt sur la tâche, on utilise UpdatedAt de la déclaration
					UpdatedAt: reloadedDeclaration.UpdatedAt,
				}
			}
			declarationDTO.Tasks = taskDTOs
		}
		
		return &declarationDTO, nil
	}
}

// GetDailyTasks récupère les tâches d'une déclaration journalière
func (s *timesheetService) GetDailyTasks(date time.Time, userID uint) ([]dto.DailyTaskDTO, error) {
	declaration, err := s.dailyDeclarationService.GetByUserIDAndDate(userID, date)
	if err != nil || declaration == nil {
		// Pas de déclaration => pas de tâches
		return []dto.DailyTaskDTO{}, nil
	}

	if len(declaration.Tasks) == 0 {
		return []dto.DailyTaskDTO{}, nil
	}

	taskDTOs := make([]dto.DailyTaskDTO, 0, len(declaration.Tasks))
	for _, task := range declaration.Tasks {
		taskDTOs = append(taskDTOs, dto.DailyTaskDTO{
			ID:        task.ID,
			TicketID:  task.TicketID,
			Ticket:    task.Ticket,
			TimeSpent: task.TimeSpent,
			CreatedAt: task.CreatedAt,
		})
	}

	return taskDTOs, nil
}

// CreateDailyTask crée une tâche dans une déclaration journalière
func (s *timesheetService) CreateDailyTask(date time.Time, userID uint, task dto.DailyTaskRequest) (*dto.DailyTaskDTO, error) {
	existingTasks, _ := s.GetDailyTasks(date, userID)
	tasks := make([]dto.DailyTaskRequest, 0, len(existingTasks)+1)
	for _, existing := range existingTasks {
		tasks = append(tasks, dto.DailyTaskRequest{
			TicketID:  existing.TicketID,
			TimeSpent: existing.TimeSpent,
		})
	}
	tasks = append(tasks, task)

	updated, err := s.CreateOrUpdateDailyDeclaration(date, userID, tasks)
	if err != nil {
		return nil, err
	}
	if len(updated.Tasks) == 0 {
		return nil, errors.New("tâche introuvable après création")
	}

	created := updated.Tasks[len(updated.Tasks)-1]
	return &dto.DailyTaskDTO{
		ID:        created.ID,
		TicketID:  created.TicketID,
		Ticket:    created.Ticket,
		TimeSpent: created.TimeSpent,
		CreatedAt: created.CreatedAt,
	}, nil
}

// DeleteDailyTask supprime une tâche d'une déclaration journalière
func (s *timesheetService) DeleteDailyTask(date time.Time, userID uint, taskID uint) error {
	existingTasks, _ := s.GetDailyTasks(date, userID)
	tasks := make([]dto.DailyTaskRequest, 0, len(existingTasks))
	for _, existing := range existingTasks {
		if existing.ID == taskID {
			continue
		}
		tasks = append(tasks, dto.DailyTaskRequest{
			TicketID:  existing.TicketID,
			TimeSpent: existing.TimeSpent,
		})
	}
	// Si toutes les tâches sont supprimées, on refuse pour éviter une déclaration vide
	if len(tasks) == 0 {
		return errors.New("au moins une tâche est requise")
	}

	_, err := s.CreateOrUpdateDailyDeclaration(date, userID, tasks)
	return err
}

// GetDailySummary récupère le résumé d'une déclaration journalière
func (s *timesheetService) GetDailySummary(date time.Time, userID uint) (*dto.DailySummaryDTO, error) {
	declaration, err := s.dailyDeclarationService.GetByUserIDAndDate(userID, date)
	if err != nil || declaration == nil {
		return nil, errors.New("déclaration introuvable")
	}

	summary := &dto.DailySummaryDTO{
		Date:      declaration.Date,
		TaskCount: declaration.TaskCount,
		TotalTime: declaration.TotalTime,
		Validated: declaration.Validated,
	}
	return summary, nil
}

// GetDailyCalendar récupère le calendrier des déclarations journalières
func (s *timesheetService) GetDailyCalendar(userID uint, startDate, endDate time.Time) ([]dto.DailyCalendarDTO, error) {
	declarations, err := s.dailyDeclarationService.GetByDateRange(userID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	calendar := make([]dto.DailyCalendarDTO, 0, len(declarations))
	for _, declaration := range declarations {
		calendar = append(calendar, dto.DailyCalendarDTO{
			Date:      declaration.Date,
			HasEntry:  true,
			TotalTime: declaration.TotalTime,
			Validated: declaration.Validated,
		})
	}

	return calendar, nil
}

// GetDailyRange récupère les déclarations journalières dans une plage de dates
func (s *timesheetService) GetDailyRange(userID uint, startDate, endDate time.Time) ([]dto.DailyDeclarationDTO, error) {
	return s.dailyDeclarationService.GetByDateRange(userID, startDate, endDate)
}

// GetAllDailyRange récupère toutes les déclarations journalières dans une plage de dates (pour les admins)
func (s *timesheetService) GetAllDailyRange(startDate, endDate time.Time) ([]dto.DailyDeclarationDTO, error) {
	return s.dailyDeclarationService.GetAllByDateRange(startDate, endDate)
}

// GetWeeklyDeclaration récupère une déclaration hebdomadaire
func (s *timesheetService) GetWeeklyDeclaration(week string, userID uint) (*dto.WeeklyDeclarationDTO, error) {
	return s.weeklyDeclarationService.GetByUserIDAndWeek(userID, week)
}

// parseWeekString parse le format YYYY-MM-Wn et retourne l'année, le mois et le numéro de semaine
func parseWeekString(week string) (year int, month int, weekNum int, err error) {
	// Format attendu: YYYY-MM-Wn (ex: 2024-01-W1)
	parts := strings.Split(week, "-")
	if len(parts) != 3 {
		return 0, 0, 0, fmt.Errorf("format de semaine invalide, attendu: YYYY-MM-Wn")
	}

	year, err = strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("année invalide")
	}

	month, err = strconv.Atoi(parts[1])
	if err != nil || month < 1 || month > 12 {
		return 0, 0, 0, fmt.Errorf("mois invalide")
	}

	// Extraire le numéro de semaine (W1, W2, etc.)
	if !strings.HasPrefix(parts[2], "W") {
		return 0, 0, 0, fmt.Errorf("format de semaine invalide, attendu: Wn")
	}
	weekNum, err = strconv.Atoi(parts[2][1:])
	if err != nil || weekNum < 1 || weekNum > 5 {
		return 0, 0, 0, fmt.Errorf("numéro de semaine invalide (doit être entre 1 et 5)")
	}

	return year, month, weekNum, nil
}

// calculateWeekDates calcule les dates de début et fin de semaine à partir du format YYYY-MM-Wn
func calculateWeekDates(year, month, weekNum int) (startDate, endDate time.Time, err error) {
	// Trouver le premier jour du mois
	firstDay := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	
	// Trouver le premier lundi du mois (ou le 1er si c'est un lundi)
	firstDayOfWeek := int(firstDay.Weekday())
	daysToFirstMonday := 0
	if firstDayOfWeek == 0 {
		// Dimanche -> lundi suivant
		daysToFirstMonday = 1
	} else if firstDayOfWeek > 1 {
		// Mardi à samedi -> lundi suivant
		daysToFirstMonday = 8 - firstDayOfWeek
	}
	
	firstMonday := firstDay.AddDate(0, 0, daysToFirstMonday)
	
	// Calculer le lundi de la semaine demandée
	weekStart := firstMonday.AddDate(0, 0, (weekNum-1)*7)
	
	// Si la semaine demandée commence avant le premier jour du mois, utiliser le 1er
	if weekStart.Before(firstDay) {
		weekStart = firstDay
	}
	
	// Le dimanche de fin de semaine
	weekEnd := weekStart.AddDate(0, 0, 6)
	
	// S'assurer que la fin de semaine ne dépasse pas la fin du mois
	lastDayOfMonth := firstDay.AddDate(0, 1, -1)
	if weekEnd.After(lastDayOfMonth) {
		weekEnd = lastDayOfMonth
	}
	
	return weekStart, weekEnd, nil
}

// CreateOrUpdateWeeklyDeclaration crée ou met à jour une déclaration hebdomadaire
func (s *timesheetService) CreateOrUpdateWeeklyDeclaration(week string, userID uint, tasks []dto.WeeklyTaskRequest) (*dto.WeeklyDeclarationDTO, error) {
	if len(tasks) == 0 {
		return nil, errors.New("au moins une tâche est requise")
	}

	// Parser le format de semaine YYYY-MM-Wn
	year, month, weekNum, err := parseWeekString(week)
	if err != nil {
		return nil, err
	}

	// Calculer les dates de début et fin de semaine
	startDate, endDate, err := calculateWeekDates(year, month, weekNum)
	if err != nil {
		return nil, err
	}

	// Vérifier si une déclaration existe déjà pour cette semaine et cet utilisateur
	existingDeclaration, err := s.weeklyDeclarationService.GetByUserIDAndWeek(userID, week)
	
	var declaration *models.WeeklyDeclaration
	totalTime := 0

	// Calculer le temps total et parser les dates des tâches
	weeklyTasks := make([]models.WeeklyDeclarationTask, len(tasks))
	for i, task := range tasks {
		// Parser la date de la tâche
		taskDate, err := time.Parse("2006-01-02", task.Date)
		if err != nil {
			return nil, fmt.Errorf("format de date invalide pour la tâche %d: %v", i+1, err)
		}
		
		totalTime += task.TimeSpent
		weeklyTasks[i] = models.WeeklyDeclarationTask{
			TicketID:  task.TicketID,
			Date:      taskDate,
			TimeSpent: task.TimeSpent,
		}
	}

	if err != nil || existingDeclaration == nil {
		// Créer une nouvelle déclaration
		declaration = &models.WeeklyDeclaration{
			UserID:    userID,
			Week:      week,
			StartDate: startDate,
			EndDate:   endDate,
			TaskCount: len(tasks),
			TotalTime: totalTime,
			Validated: false,
			Tasks:     weeklyTasks,
		}

		// Créer la déclaration avec les tâches
		if err := database.DB.Create(declaration).Error; err != nil {
			return nil, errors.New("erreur lors de la création de la déclaration")
		}
	} else {
		// Mettre à jour la déclaration existante
		// Supprimer les anciennes tâches
		if err := database.DB.Where("declaration_id = ?", existingDeclaration.ID).Delete(&models.WeeklyDeclarationTask{}).Error; err != nil {
			return nil, errors.New("erreur lors de la suppression des anciennes tâches")
		}

		// Mettre à jour les IDs des tâches
		for i := range weeklyTasks {
			weeklyTasks[i].DeclarationID = existingDeclaration.ID
		}

		// Mettre à jour la déclaration
		if err := database.DB.Model(&models.WeeklyDeclaration{}).Where("id = ?", existingDeclaration.ID).
			Updates(map[string]interface{}{
				"start_date": startDate,
				"end_date":   endDate,
				"task_count": len(tasks),
				"total_time": totalTime,
			}).Error; err != nil {
			return nil, errors.New("erreur lors de la mise à jour de la déclaration")
		}

		// Créer les nouvelles tâches
		if err := database.DB.Create(&weeklyTasks).Error; err != nil {
			return nil, errors.New("erreur lors de la création des tâches")
		}
	}

	// Récupérer la déclaration complète avec les relations
	updatedDeclaration, err := s.weeklyDeclarationService.GetByUserIDAndWeek(userID, week)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération de la déclaration")
	}

	return updatedDeclaration, nil
}

// GetWeeklyTasks récupère les tâches d'une déclaration hebdomadaire
func (s *timesheetService) GetWeeklyTasks(week string, userID uint) ([]dto.WeeklyTaskDTO, error) {
	// TODO: Implémenter
	return nil, errors.New("non implémenté")
}

// GetWeeklySummary récupère le résumé d'une déclaration hebdomadaire
func (s *timesheetService) GetWeeklySummary(week string, userID uint) (*dto.WeeklySummaryDTO, error) {
	// TODO: Implémenter
	return nil, errors.New("non implémenté")
}

// GetWeeklyDailyBreakdown récupère la répartition quotidienne d'une déclaration hebdomadaire
func (s *timesheetService) GetWeeklyDailyBreakdown(week string, userID uint) ([]dto.DailyBreakdownDTO, error) {
	// TODO: Implémenter
	return nil, errors.New("non implémenté")
}

// ValidateWeeklyDeclaration valide une déclaration hebdomadaire
func (s *timesheetService) ValidateWeeklyDeclaration(week string, userID uint, validatedByID uint) (*dto.WeeklyDeclarationDTO, error) {
	declaration, err := s.weeklyDeclarationService.GetByUserIDAndWeek(userID, week)
	if err != nil {
		return nil, err
	}
	return s.weeklyDeclarationService.Validate(declaration.ID, validatedByID)
}

// GetWeeklyValidationStatus récupère le statut de validation d'une déclaration hebdomadaire
func (s *timesheetService) GetWeeklyValidationStatus(week string, userID uint) (*dto.ValidationStatusDTO, error) {
	declaration, err := s.weeklyDeclarationService.GetByUserIDAndWeek(userID, week)
	if err != nil {
		return nil, err
	}
	return &dto.ValidationStatusDTO{
		Validated:   declaration.Validated,
		ValidatedBy: declaration.ValidatedBy,
		ValidatedAt: declaration.ValidatedAt,
	}, nil
}

// SetTicketEstimatedTime définit le temps estimé d'un ticket.
// Si le ticket est en "ouvert", il passe automatiquement à "en_cours".
func (s *timesheetService) SetTicketEstimatedTime(ticketID uint, estimatedTime int, userID uint) error {
	ticket, err := s.ticketRepo.FindByID(ticketID)
	if err != nil {
		return errors.New("ticket introuvable")
	}
	ticket.EstimatedTime = &estimatedTime
	if ticket.Status == "ouvert" {
		ticket.Status = "en_cours"
	}
	return s.ticketRepo.Update(ticket)
}

// GetTicketEstimatedTime récupère le temps estimé d'un ticket
func (s *timesheetService) GetTicketEstimatedTime(ticketID uint) (*dto.EstimatedTimeDTO, error) {
	ticket, err := s.ticketRepo.FindByID(ticketID)
	if err != nil {
		return nil, errors.New("ticket introuvable")
	}
	
	// Vérifier que le ticket a un temps estimé défini (0 est une valeur valide)
	if ticket.EstimatedTime == nil {
		return nil, errors.New("temps estimé introuvable")
	}
	
	return &dto.EstimatedTimeDTO{
		TicketID:      ticketID,
		EstimatedTime: *ticket.EstimatedTime,
	}, nil
}

// UpdateTicketEstimatedTime met à jour le temps estimé d'un ticket
func (s *timesheetService) UpdateTicketEstimatedTime(ticketID uint, estimatedTime int, userID uint) error {
	return s.SetTicketEstimatedTime(ticketID, estimatedTime, userID)
}

// GetTicketTimeComparison récupère la comparaison temps estimé vs réel d'un ticket
func (s *timesheetService) GetTicketTimeComparison(ticketID uint) (*dto.TimeComparisonDTO, error) {
	ticket, err := s.ticketRepo.FindByID(ticketID)
	if err != nil {
		return nil, errors.New("ticket introuvable")
	}
	estimatedTime := 0
	if ticket.EstimatedTime != nil {
		estimatedTime = *ticket.EstimatedTime
	}
	actualTime := 0
	if ticket.ActualTime != nil {
		actualTime = *ticket.ActualTime
	}
	diff := actualTime - estimatedTime
	return &dto.TimeComparisonDTO{
		TicketID:      ticketID,
		EstimatedTime: estimatedTime,
		ActualTime:    actualTime,
		Difference:    diff,
		Percentage:    float64(diff) / float64(estimatedTime) * 100,
	}, nil
}

// GetProjectTimeBudget récupère le budget temps d'un projet
func (s *timesheetService) GetProjectTimeBudget(projectID uint) (*dto.ProjectTimeBudgetDTO, error) {
	// TODO: Implémenter
	return nil, errors.New("non implémenté")
}

// SetProjectTimeBudget définit le budget temps d'un projet
func (s *timesheetService) SetProjectTimeBudget(projectID uint, budget dto.SetProjectTimeBudgetRequest, userID uint) error {
	// TODO: Implémenter
	return errors.New("non implémenté")
}

// GetBudgetAlerts récupère les alertes de budget
func (s *timesheetService) GetBudgetAlerts() ([]dto.BudgetAlertDTO, error) {
	// Récupérer tous les tickets avec temps estimé (utiliser une grande limite pour récupérer tous les tickets)
	tickets, _, err := s.ticketRepo.FindAll(nil, 1, 10000, nil) // nil scope = pas de filtre (utilisé en interne)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des tickets")
	}

	var alerts []dto.BudgetAlertDTO
	for _, ticket := range tickets {
		// Vérifier si le ticket a un temps estimé
		if ticket.EstimatedTime == nil || *ticket.EstimatedTime == 0 {
			continue
		}

		// Récupérer le temps réel passé sur ce ticket
		timeEntries, err := s.timeEntryService.GetByTicketID(ticket.ID)
		if err != nil {
			// Ignorer les erreurs pour un ticket spécifique
			continue
		}

		actualTime := 0
		for _, entry := range timeEntries {
			actualTime += entry.TimeSpent
		}

		estimatedTime := *ticket.EstimatedTime
		// Créer une alerte si le temps réel dépasse 80% du temps estimé
		if actualTime > 0 {
			percentage := float64(actualTime) / float64(estimatedTime) * 100
			if percentage >= 80 {
				ticketID := ticket.ID
				alerts = append(alerts, dto.BudgetAlertDTO{
					TicketID:   &ticketID,
					AlertType: "budget_exceeded",
					Message:   "Le temps réel dépasse le temps estimé",
					Budget:    estimatedTime,
					Spent:     actualTime,
					Percentage: percentage,
					CreatedAt: time.Now(),
				})
			}
		}
	}

	return alerts, nil
}

// GetTicketBudgetStatus récupère le statut du budget d'un ticket
func (s *timesheetService) GetTicketBudgetStatus(ticketID uint) (*dto.BudgetStatusDTO, error) {
	comparison, err := s.GetTicketTimeComparison(ticketID)
	if err != nil {
		return nil, err
	}
	status := "on_budget"
	if comparison.Difference > 0 {
		status = "over_budget"
	} else if comparison.Difference < 0 {
		status = "under_budget"
	}
	return &dto.BudgetStatusDTO{
		TicketID: ticketID,
		Status:   status,
		Comparison: *comparison,
	}, nil
}

// ValidateTimeEntry valide une entrée de temps
func (s *timesheetService) ValidateTimeEntry(id uint, req dto.ValidateTimeEntryRequest, validatedByID uint) (*dto.TimeEntryDTO, error) {
	return s.timeEntryService.Validate(id, req, validatedByID)
}

// GetPendingValidationEntries récupère les entrées en attente de validation
func (s *timesheetService) GetPendingValidationEntries(scope interface{}) ([]dto.TimeEntryDTO, error) {
	return s.timeEntryService.GetPendingValidation(scope)
}

// GetValidationHistory récupère l'historique de validation
func (s *timesheetService) GetValidationHistory() ([]dto.ValidationHistoryDTO, error) {
	// TODO: Implémenter
	return nil, errors.New("non implémenté")
}

// GetDelayAlerts récupère les alertes de retard
func (s *timesheetService) GetDelayAlerts() ([]dto.DelayAlertDTO, error) {
	// Passer nil comme scope car c'est une méthode interne
	delays, err := s.delayRepo.FindUnjustified(nil)
	if err != nil {
		return nil, err
	}
	alerts := make([]dto.DelayAlertDTO, len(delays))
	for i, delay := range delays {
		var ticketID uint
		if delay.TicketID != nil {
			ticketID = *delay.TicketID
		}
		alerts[i] = dto.DelayAlertDTO{
			DelayID:    delay.ID,
			TicketID:   ticketID,
			UserID:     delay.UserID,
			DelayTime:  delay.DelayTime,
			DetectedAt: delay.DetectedAt,
		}
	}
	return alerts, nil
}

// GetBudgetAlertsForTimesheet récupère les alertes de budget pour le timesheet
func (s *timesheetService) GetBudgetAlertsForTimesheet() ([]dto.BudgetAlertDTO, error) {
	return s.GetBudgetAlerts()
}

// GetOverloadAlerts récupère les alertes de surcharge
func (s *timesheetService) GetOverloadAlerts() ([]dto.OverloadAlertDTO, error) {
	// TODO: Implémenter
	return nil, errors.New("non implémenté")
}

// GetUnderloadAlerts récupère les alertes de sous-charge
func (s *timesheetService) GetUnderloadAlerts() ([]dto.UnderloadAlertDTO, error) {
	// TODO: Implémenter
	return nil, errors.New("non implémenté")
}

// SendReminderAlerts envoie des rappels
func (s *timesheetService) SendReminderAlerts(userIDs []uint) error {
	// TODO: Implémenter
	return errors.New("non implémenté")
}

// GetPendingJustificationAlerts récupère les alertes de justifications en attente
func (s *timesheetService) GetPendingJustificationAlerts() ([]dto.PendingJustificationAlertDTO, error) {
	justifications, err := s.delayJustificationRepo.FindPending()
	if err != nil {
		return nil, err
	}
	alerts := make([]dto.PendingJustificationAlertDTO, len(justifications))
	for i, justification := range justifications {
		alerts[i] = dto.PendingJustificationAlertDTO{
			JustificationID: justification.ID,
			DelayID:          justification.DelayID,
			UserID:           justification.UserID,
			CreatedAt:        justification.CreatedAt,
		}
	}
	return alerts, nil
}

// GetTimesheetHistory récupère l'historique du timesheet
func (s *timesheetService) GetTimesheetHistory(userID uint, startDate, endDate time.Time) ([]dto.TimesheetHistoryDTO, error) {
	// TODO: Implémenter
	return nil, errors.New("non implémenté")
}

// GetTimesheetHistoryEntry récupère une entrée de l'historique
func (s *timesheetService) GetTimesheetHistoryEntry(entryID uint) (*dto.TimesheetHistoryEntryDTO, error) {
	// TODO: Implémenter
	return nil, errors.New("non implémenté")
}

// GetTimesheetAuditTrail récupère la piste d'audit du timesheet
func (s *timesheetService) GetTimesheetAuditTrail(userID uint, startDate, endDate time.Time) ([]dto.AuditTrailDTO, error) {
	// TODO: Implémenter
	return nil, errors.New("non implémenté")
}

// GetTimesheetModifications récupère les modifications du timesheet
func (s *timesheetService) GetTimesheetModifications(userID uint, startDate, endDate time.Time) ([]dto.ModificationDTO, error) {
	// TODO: Implémenter
	return nil, errors.New("non implémenté")
}

