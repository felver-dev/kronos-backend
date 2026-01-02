package dto

import "time"

// DashboardDTO représente le tableau de bord complet
type DashboardDTO struct {
	Tickets     TicketStatsDTO     `json:"tickets"`      // Statistiques des tickets
	SLA         SLAStatsDTO        `json:"sla"`          // Statistiques SLA
	Performance PerformanceStatsDTO `json:"performance"` // Statistiques de performance
	Alerts      []AlertDTO         `json:"alerts"`       // Alertes en cours
	Period      string             `json:"period"`       // Période (week, month, etc.)
}

// TicketStatsDTO représente les statistiques des tickets
type TicketStatsDTO struct {
	Total              int                `json:"total"`                // Nombre total de tickets
	ByCategory         map[string]int     `json:"by_category"`         // Par catégorie
	ByStatus           map[string]int     `json:"by_status"`            // Par statut
	ByPriority         map[string]int     `json:"by_priority"`          // Par priorité
	AverageResolutionTime float64         `json:"average_resolution_time"` // Temps moyen de résolution en minutes
	Delayed            int                `json:"delayed"`              // Tickets en retard
	Open               int                `json:"open"`                 // Tickets ouverts
	Closed             int                `json:"closed"`               // Tickets fermés
}

// SLAStatsDTO représente les statistiques SLA
type SLAStatsDTO struct {
	OverallCompliance float64            `json:"overall_compliance"` // Conformité globale en %
	ByCategory        map[string]float64 `json:"by_category"`        // Conformité par catégorie
	ByPriority        map[string]float64 `json:"by_priority"`        // Conformité par priorité
	TotalViolations   int                `json:"total_violations"`   // Nombre total de violations
	AtRisk            int                `json:"at_risk"`            // Tickets à risque
}

// PerformanceStatsDTO représente les statistiques de performance
type PerformanceStatsDTO struct {
	TotalTimeSpent      int     `json:"total_time_spent"`       // Temps total passé en minutes
	AverageEfficiency   float64 `json:"average_efficiency"`     // Efficacité moyenne en %
	AverageProductivity float64 `json:"average_productivity"`   // Productivité moyenne (tickets/heure)
	TotalTicketsTreated int     `json:"total_tickets_treated"` // Nombre total de tickets traités
}

// AlertDTO représente une alerte
type AlertDTO struct {
	Type      string    `json:"type"`       // delay_alert, budget_alert, etc.
	Severity  string    `json:"severity"`   // low, medium, high, critical
	Message   string    `json:"message"`    // Message de l'alerte
	LinkURL   string    `json:"link_url,omitempty"` // URL vers la ressource (optionnel)
	CreatedAt time.Time `json:"created_at"`
}

// TicketCountReportDTO représente le rapport de nombre de tickets par période
type TicketCountReportDTO struct {
	Period    string            `json:"period"`    // week, month, etc.
	Count     int               `json:"count"`     // Nombre total
	Breakdown []PeriodBreakdownDTO `json:"breakdown"` // Répartition par sous-période
}

// PeriodBreakdownDTO représente la répartition par sous-période
type PeriodBreakdownDTO struct {
	Date  time.Time `json:"date"`  // Date de la période
	Count int       `json:"count"` // Nombre pour cette période
}

// TicketTypeDistributionDTO représente la répartition des tickets par type
type TicketTypeDistributionDTO struct {
	Incidents      int `json:"incidents"`      // Nombre d'incidents
	Demandes       int `json:"demandes"`       // Nombre de demandes
	Changements    int `json:"changements"`    // Nombre de changements
	Developpements int `json:"developpements"` // Nombre de développements
}

// AverageResolutionTimeDTO représente le temps moyen de résolution
type AverageResolutionTimeDTO struct {
	AverageTime int                `json:"average_time"` // Temps moyen en minutes
	Unit        string             `json:"unit"`         // minutes
	Breakdown   map[string]float64 `json:"breakdown"`    // Répartition par catégorie/priorité
}

// WorkloadByAgentDTO représente la charge de travail par agent
type WorkloadByAgentDTO struct {
	UserID      uint   `json:"user_id"`
	User        *UserDTO `json:"user,omitempty"`
	TicketCount int    `json:"ticket_count"`     // Nombre de tickets
	AverageTime float64 `json:"average_time"`    // Temps moyen en minutes
	TotalTime   int    `json:"total_time"`       // Temps total en minutes
}

// CustomReportRequest représente la requête pour un rapport personnalisé
type CustomReportRequest struct {
	Metrics []string    `json:"metrics" binding:"required"` // Métriques à inclure (obligatoire)
	Period  string      `json:"period" binding:"required"`  // Période (obligatoire)
	StartDate *time.Time `json:"start_date,omitempty"`     // Date de début (optionnel)
	EndDate   *time.Time `json:"end_date,omitempty"`       // Date de fin (optionnel)
	Filters   map[string]any `json:"filters,omitempty"`    // Filtres personnalisés (optionnel)
	Format    string    `json:"format,omitempty"`          // Format d'export (pdf, excel, csv) (optionnel)
}

