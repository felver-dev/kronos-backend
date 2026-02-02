package dto

import "time"

// TimeEntryDTO représente une entrée de temps
type TimeEntryDTO struct {
	ID             uint       `json:"id"`
	TicketID       uint       `json:"ticket_id"`
	ProjectTaskID  *uint      `json:"project_task_id,omitempty"`
	Ticket         *TicketDTO `json:"ticket,omitempty"`
	UserID         uint       `json:"user_id"`
	User        *UserDTO   `json:"user,omitempty"`
	TimeSpent   int        `json:"time_spent"`
	Date        time.Time  `json:"date"`
	Description string     `json:"description,omitempty"`
	Validated   bool       `json:"validated"`
	ValidatedBy *uint      `json:"validated_by,omitempty"`
	ValidatedAt *time.Time `json:"validated_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// CreateTimeEntryRequest représente la requête de création d'une entrée de temps
type CreateTimeEntryRequest struct {
	TicketID    uint   `json:"ticket_id" binding:"required"`
	TimeSpent   int    `json:"time_spent" binding:"required"`
	Date        string `json:"date" binding:"required"` // Format: YYYY-MM-DD
	Description string `json:"description,omitempty"`
}

// UpdateTimeEntryRequest représente la requête de mise à jour d'une entrée de temps
type UpdateTimeEntryRequest struct {
	TimeSpent   int    `json:"time_spent,omitempty"`
	Date        string `json:"date,omitempty"` // Format: YYYY-MM-DD
	Description string `json:"description,omitempty"`
}

// ValidateTimeEntryRequest représente la requête de validation d'une entrée de temps
type ValidateTimeEntryRequest struct {
	Validated *bool `json:"validated" binding:"required"`
}

// DailyDeclarationDTO représente une déclaration journalière
type DailyDeclarationDTO struct {
	ID                uint           `json:"id"`
	UserID            uint           `json:"user_id"`
	User              *UserDTO       `json:"user,omitempty"`
	Date              time.Time      `json:"date"`
	TaskCount         int            `json:"task_count"`
	TotalTime         int            `json:"total_time"`
	Validated         bool           `json:"validated"`
	ValidatedBy       *uint          `json:"validated_by,omitempty"`
	ValidatedAt       *time.Time     `json:"validated_at,omitempty"`
	ValidationComment string         `json:"validation_comment,omitempty"`
	Tasks             []TimeEntryDTO `json:"tasks,omitempty"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
}

// WeeklyDeclarationDTO représente une déclaration hebdomadaire
type WeeklyDeclarationDTO struct {
	ID            uint               `json:"id"`
	UserID        uint               `json:"user_id"`
	User          *UserDTO           `json:"user,omitempty"`
	Week          string             `json:"week"`
	StartDate     time.Time          `json:"start_date"`
	EndDate       time.Time          `json:"end_date"`
	TaskCount     int                `json:"task_count"`
	TotalTime     int                `json:"total_time"`
	Validated     bool               `json:"validated"`
	ValidatedBy   *uint              `json:"validated_by,omitempty"`
	ValidatedAt   *time.Time         `json:"validated_at,omitempty"`
	ValidationComment string          `json:"validation_comment,omitempty"`
	DailyBreakdown []DailyBreakdownDTO `json:"daily_breakdown,omitempty"`
	CreatedAt     time.Time          `json:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at"`
}

// DailyTaskRequest représente une requête pour créer/mettre à jour une tâche journalière
type DailyTaskRequest struct {
	TicketID  uint `json:"ticket_id" binding:"required"`
	TimeSpent int  `json:"time_spent" binding:"required"`
}

// DailyTaskDTO représente une tâche journalière
type DailyTaskDTO struct {
	ID        uint      `json:"id"`
	TicketID  uint      `json:"ticket_id"`
	Ticket    *TicketDTO `json:"ticket,omitempty"`
	TimeSpent int       `json:"time_spent"`
	CreatedAt time.Time `json:"created_at"`
}

// DailySummaryDTO représente le résumé d'une déclaration journalière
type DailySummaryDTO struct {
	Date      time.Time `json:"date"`
	TaskCount int       `json:"task_count"`
	TotalTime int       `json:"total_time"`
	Validated bool      `json:"validated"`
}

// DailyCalendarDTO représente une entrée du calendrier journalier
type DailyCalendarDTO struct {
	Date      time.Time `json:"date"`
	HasEntry  bool      `json:"has_entry"`
	TotalTime int       `json:"total_time"`
	Validated bool      `json:"validated"`
}

// WeeklyTaskRequest représente une requête pour créer/mettre à jour une tâche hebdomadaire
type WeeklyTaskRequest struct {
	TicketID  uint   `json:"ticket_id" binding:"required"`
	Date      string `json:"date" binding:"required"` // Format: YYYY-MM-DD
	TimeSpent int    `json:"time_spent" binding:"required"`
}

// WeeklyTaskDTO représente une tâche hebdomadaire
type WeeklyTaskDTO struct {
	ID        uint       `json:"id"`
	TicketID  uint       `json:"ticket_id"`
	Ticket    *TicketDTO `json:"ticket,omitempty"`
	Date      time.Time  `json:"date"`
	TimeSpent int        `json:"time_spent"`
	CreatedAt time.Time `json:"created_at"`
}

// WeeklySummaryDTO représente le résumé d'une déclaration hebdomadaire
type WeeklySummaryDTO struct {
	Week      string    `json:"week"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	TaskCount int       `json:"task_count"`
	TotalTime int       `json:"total_time"`
	Validated bool      `json:"validated"`
}

// DailyBreakdownDTO représente la répartition quotidienne d'une semaine
type DailyBreakdownDTO struct {
	Date      time.Time `json:"date"`
	TaskCount int       `json:"task_count"`
	TotalTime int       `json:"total_time"`
}

// ValidationStatusDTO représente le statut de validation
type ValidationStatusDTO struct {
	Validated   bool       `json:"validated"`
	ValidatedBy *uint      `json:"validated_by,omitempty"`
	ValidatedAt *time.Time `json:"validated_at,omitempty"`
}

// EstimatedTimeDTO représente le temps estimé d'un ticket
type EstimatedTimeDTO struct {
	TicketID      uint `json:"ticket_id"`
	EstimatedTime int  `json:"estimated_time"`
}

// SetEstimatedTimeRequest représente une requête pour définir/mettre à jour le temps estimé
// @Description Requête pour définir ou mettre à jour le temps estimé d'un ticket en minutes
type SetEstimatedTimeRequest struct {
	// Temps estimé en minutes
	// @Example 120
	EstimatedTime int `json:"estimated_time" binding:"required"`
}

// TimeComparisonDTO représente la comparaison temps estimé vs réel
type TimeComparisonDTO struct {
	TicketID      uint    `json:"ticket_id"`
	EstimatedTime int    `json:"estimated_time"`
	ActualTime    int    `json:"actual_time"`
	Difference    int    `json:"difference"`
	Percentage    float64 `json:"percentage"`
}

// ProjectTimeBudgetDTO représente le budget temps d'un projet
type ProjectTimeBudgetDTO struct {
	ProjectID     uint    `json:"project_id"`
	Budget        int     `json:"budget"`
	Spent         int     `json:"spent"`
	Remaining     int     `json:"remaining"`
	Percentage    float64 `json:"percentage"`
}

// SetProjectTimeBudgetRequest représente une requête pour définir le budget temps d'un projet
type SetProjectTimeBudgetRequest struct {
	Budget int `json:"budget" binding:"required"`
}

// BudgetAlertDTO représente une alerte de budget
type BudgetAlertDTO struct {
	TicketID      *uint    `json:"ticket_id,omitempty"`
	ProjectID     *uint    `json:"project_id,omitempty"`
	AlertType     string  `json:"alert_type"`
	Message       string  `json:"message"`
	Budget        int     `json:"budget"`
	Spent         int     `json:"spent"`
	Percentage    float64 `json:"percentage"`
	CreatedAt     time.Time `json:"created_at"`
}

// BudgetStatusDTO représente le statut du budget d'un ticket
type BudgetStatusDTO struct {
	TicketID   uint              `json:"ticket_id"`
	Status     string            `json:"status"` // on_budget, over_budget, under_budget
	Comparison TimeComparisonDTO `json:"comparison"`
}

// DelayAlertDTO représente une alerte de retard
type DelayAlertDTO struct {
	DelayID    uint      `json:"delay_id"`
	TicketID   uint      `json:"ticket_id"`
	UserID     uint      `json:"user_id"`
	DelayTime  int       `json:"delay_time"`
	DetectedAt time.Time `json:"detected_at"`
}

// OverloadAlertDTO représente une alerte de surcharge
type OverloadAlertDTO struct {
	UserID     uint      `json:"user_id"`
	Date       time.Time `json:"date"`
	ActualTime int       `json:"actual_time"`
	MaxTime    int       `json:"max_time"`
	Message    string    `json:"message"`
}

// UnderloadAlertDTO représente une alerte de sous-charge
type UnderloadAlertDTO struct {
	UserID     uint      `json:"user_id"`
	Date       time.Time `json:"date"`
	ActualTime int       `json:"actual_time"`
	MinTime    int       `json:"min_time"`
	Message    string    `json:"message"`
}

// PendingJustificationAlertDTO représente une alerte de justification en attente
type PendingJustificationAlertDTO struct {
	JustificationID uint      `json:"justification_id"`
	DelayID         uint      `json:"delay_id"`
	UserID          uint      `json:"user_id"`
	CreatedAt       time.Time `json:"created_at"`
}

// ValidationHistoryDTO représente une entrée de l'historique de validation
type ValidationHistoryDTO struct {
	EntryID     uint       `json:"entry_id"`
	EntryType   string     `json:"entry_type"` // time_entry, daily, weekly
	UserID      uint       `json:"user_id"`
	ValidatedBy uint       `json:"validated_by"`
	ValidatedAt time.Time  `json:"validated_at"`
	Status      string     `json:"status"` // validated, rejected
}

// TimesheetHistoryDTO représente une entrée de l'historique du timesheet
type TimesheetHistoryDTO struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Date      time.Time `json:"date"`
	Action    string    `json:"action"`
	Details   string    `json:"details"`
	CreatedAt time.Time `json:"created_at"`
}

// TimesheetHistoryEntryDTO représente une entrée détaillée de l'historique
type TimesheetHistoryEntryDTO struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Date      time.Time `json:"date"`
	Action    string    `json:"action"`
	Details   string    `json:"details"`
	Changes   map[string]interface{} `json:"changes,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// AuditTrailDTO représente une entrée de la piste d'audit
type AuditTrailDTO struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Action    string    `json:"action"`
	EntityType string   `json:"entity_type"`
	EntityID   uint     `json:"entity_id"`
	Details   string    `json:"details"`
	CreatedAt time.Time `json:"created_at"`
}

// ModificationDTO représente une modification du timesheet
type ModificationDTO struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Date      time.Time `json:"date"`
	Field     string    `json:"field"`
	OldValue  interface{} `json:"old_value"`
	NewValue  interface{} `json:"new_value"`
	CreatedAt time.Time `json:"created_at"`
}
