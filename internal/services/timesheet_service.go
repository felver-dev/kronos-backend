package services

import (
	"errors"
	"time"

	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/repositories"
)

// TimesheetService interface pour les opérations sur les timesheets
type TimesheetService interface {
	// Saisie du temps par ticket
	CreateTimeEntry(req dto.CreateTimeEntryRequest, userID uint) (*dto.TimeEntryDTO, error)
	GetTimeEntries() ([]dto.TimeEntryDTO, error)
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
	GetPendingValidationEntries() ([]dto.TimeEntryDTO, error)
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
func (s *timesheetService) GetTimeEntries() ([]dto.TimeEntryDTO, error) {
	return s.timeEntryService.GetAll()
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
	// TODO: Implémenter la logique de création/mise à jour
	return nil, errors.New("non implémenté")
}

// GetDailyTasks récupère les tâches d'une déclaration journalière
func (s *timesheetService) GetDailyTasks(date time.Time, userID uint) ([]dto.DailyTaskDTO, error) {
	// TODO: Implémenter
	return nil, errors.New("non implémenté")
}

// CreateDailyTask crée une tâche dans une déclaration journalière
func (s *timesheetService) CreateDailyTask(date time.Time, userID uint, task dto.DailyTaskRequest) (*dto.DailyTaskDTO, error) {
	// TODO: Implémenter
	return nil, errors.New("non implémenté")
}

// DeleteDailyTask supprime une tâche d'une déclaration journalière
func (s *timesheetService) DeleteDailyTask(date time.Time, userID uint, taskID uint) error {
	// TODO: Implémenter
	return errors.New("non implémenté")
}

// GetDailySummary récupère le résumé d'une déclaration journalière
func (s *timesheetService) GetDailySummary(date time.Time, userID uint) (*dto.DailySummaryDTO, error) {
	// TODO: Implémenter
	return nil, errors.New("non implémenté")
}

// GetDailyCalendar récupère le calendrier des déclarations journalières
func (s *timesheetService) GetDailyCalendar(userID uint, startDate, endDate time.Time) ([]dto.DailyCalendarDTO, error) {
	// TODO: Implémenter
	return nil, errors.New("non implémenté")
}

// GetDailyRange récupère les déclarations journalières dans une plage de dates
func (s *timesheetService) GetDailyRange(userID uint, startDate, endDate time.Time) ([]dto.DailyDeclarationDTO, error) {
	return s.dailyDeclarationService.GetByDateRange(userID, startDate, endDate)
}

// GetWeeklyDeclaration récupère une déclaration hebdomadaire
func (s *timesheetService) GetWeeklyDeclaration(week string, userID uint) (*dto.WeeklyDeclarationDTO, error) {
	return s.weeklyDeclarationService.GetByUserIDAndWeek(userID, week)
}

// CreateOrUpdateWeeklyDeclaration crée ou met à jour une déclaration hebdomadaire
func (s *timesheetService) CreateOrUpdateWeeklyDeclaration(week string, userID uint, tasks []dto.WeeklyTaskRequest) (*dto.WeeklyDeclarationDTO, error) {
	// TODO: Implémenter
	return nil, errors.New("non implémenté")
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

// SetTicketEstimatedTime définit le temps estimé d'un ticket
func (s *timesheetService) SetTicketEstimatedTime(ticketID uint, estimatedTime int, userID uint) error {
	ticket, err := s.ticketRepo.FindByID(ticketID)
	if err != nil {
		return errors.New("ticket introuvable")
	}
	ticket.EstimatedTime = &estimatedTime
	return s.ticketRepo.Update(ticket)
}

// GetTicketEstimatedTime récupère le temps estimé d'un ticket
func (s *timesheetService) GetTicketEstimatedTime(ticketID uint) (*dto.EstimatedTimeDTO, error) {
	ticket, err := s.ticketRepo.FindByID(ticketID)
	if err != nil {
		return nil, errors.New("ticket introuvable")
	}
	estimatedTime := 0
	if ticket.EstimatedTime != nil {
		estimatedTime = *ticket.EstimatedTime
	}
	return &dto.EstimatedTimeDTO{
		TicketID:      ticketID,
		EstimatedTime: estimatedTime,
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
	// TODO: Implémenter
	return nil, errors.New("non implémenté")
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
func (s *timesheetService) GetPendingValidationEntries() ([]dto.TimeEntryDTO, error) {
	return s.timeEntryService.GetPendingValidation()
}

// GetValidationHistory récupère l'historique de validation
func (s *timesheetService) GetValidationHistory() ([]dto.ValidationHistoryDTO, error) {
	// TODO: Implémenter
	return nil, errors.New("non implémenté")
}

// GetDelayAlerts récupère les alertes de retard
func (s *timesheetService) GetDelayAlerts() ([]dto.DelayAlertDTO, error) {
	delays, err := s.delayRepo.FindUnjustified()
	if err != nil {
		return nil, err
	}
	alerts := make([]dto.DelayAlertDTO, len(delays))
	for i, delay := range delays {
		alerts[i] = dto.DelayAlertDTO{
			DelayID:    delay.ID,
			TicketID:   delay.TicketID,
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

