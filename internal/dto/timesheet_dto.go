package dto

import "time"

// TimeEntryDTO représente une entrée de temps (temps passé sur un ticket)
type TimeEntryDTO struct {
	ID          uint       `json:"id"`
	TicketID    uint       `json:"ticket_id"`
	Ticket      *TicketDTO `json:"ticket,omitempty"` // Ticket associé (optionnel)
	UserID      uint       `json:"user_id"`
	User        *UserDTO   `json:"user,omitempty"` // Utilisateur (optionnel)
	TimeSpent   int        `json:"time_spent"`     // Temps passé en minutes
	Date        time.Time  `json:"date"`           // Date de l'entrée
	Description string     `json:"description,omitempty"`
	Validated   bool       `json:"validated"` // Si l'entrée a été validée
	ValidatedBy *uint      `json:"validated_by,omitempty"`
	ValidatedAt *time.Time `json:"validated_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// CreateTimeEntryRequest représente la requête de création d'une entrée de temps
type CreateTimeEntryRequest struct {
	TicketID    uint   `json:"ticket_id" binding:"required"`        // ID du ticket (obligatoire)
	TimeSpent   int    `json:"time_spent" binding:"required,min=1"` // Temps passé en minutes (obligatoire, min 1)
	Date        string `json:"date" binding:"required"`             // Date au format "2006-01-02" (obligatoire)
	Description string `json:"description,omitempty"`               // Description (optionnel)
}

// UpdateTimeEntryRequest représente la requête de mise à jour d'une entrée de temps
type UpdateTimeEntryRequest struct {
	TimeSpent   int    `json:"time_spent,omitempty" binding:"omitempty,min=1"` // Temps passé en minutes (optionnel)
	Description string `json:"description,omitempty"`                          // Description (optionnel)
}

// DailyDeclarationDTO représente une déclaration journalière des tâches
type DailyDeclarationDTO struct {
	ID          uint           `json:"id"`
	UserID      uint           `json:"user_id"`
	User        *UserDTO       `json:"user,omitempty"`
	Date        time.Time      `json:"date"`
	TaskCount   int            `json:"task_count"` // Nombre de tâches
	TotalTime   int            `json:"total_time"` // Temps total en minutes
	Validated   bool           `json:"validated"`  // Si la déclaration a été validée
	ValidatedBy *uint          `json:"validated_by,omitempty"`
	ValidatedAt *time.Time     `json:"validated_at,omitempty"`
	Tasks       []TimeEntryDTO `json:"tasks,omitempty"` // Tâches déclarées (optionnel)
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// WeeklyDeclarationDTO représente une déclaration hebdomadaire des tâches
type WeeklyDeclarationDTO struct {
	ID             uint                `json:"id"`
	UserID         uint                `json:"user_id"`
	User           *UserDTO            `json:"user,omitempty"`
	Week           string              `json:"week"` // Format ISO: "2024-W03"
	StartDate      time.Time           `json:"start_date"`
	EndDate        time.Time           `json:"end_date"`
	TaskCount      int                 `json:"task_count"` // Nombre total de tâches
	TotalTime      int                 `json:"total_time"` // Temps total en minutes
	Validated      bool                `json:"validated"`  // Si la déclaration a été validée
	ValidatedBy    *uint               `json:"validated_by,omitempty"`
	ValidatedAt    *time.Time          `json:"validated_at,omitempty"`
	DailyBreakdown []DailyBreakdownDTO `json:"daily_breakdown,omitempty"` // Répartition par jour (optionnel)
	CreatedAt      time.Time           `json:"created_at"`
	UpdatedAt      time.Time           `json:"updated_at"`
}

// DailyBreakdownDTO représente la répartition des tâches par jour
type DailyBreakdownDTO struct {
	Date      time.Time `json:"date"`
	TaskCount int       `json:"task_count"` // Nombre de tâches du jour
	TotalTime int       `json:"total_time"` // Temps total du jour en minutes
}

// ValidateTimeEntryRequest représente la requête de validation d'une entrée de temps
type ValidateTimeEntryRequest struct {
	Validated bool   `json:"validated" binding:"required"` // true pour valider, false pour invalider
	Comment   string `json:"comment,omitempty"`            // Commentaire de validation (optionnel)
}

// TimeComparisonDTO représente la comparaison entre temps estimé et temps réel
type TimeComparisonDTO struct {
	EstimatedTime int     `json:"estimated_time"` // Temps estimé en minutes
	ActualTime    int     `json:"actual_time"`    // Temps réel en minutes
	Difference    int     `json:"difference"`     // Différence en minutes (actual - estimated)
	Percentage    float64 `json:"percentage"`     // Pourcentage (100 = estimé, >100 = dépassement)
	Status        string  `json:"status"`         // "under", "on_time", "over"
}
