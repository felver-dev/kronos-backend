package dto

import "time"

// DashboardDTO représente le tableau de bord complet
type DashboardDTO struct {
	Tickets       TicketStatsDTO        `json:"tickets"`        // Statistiques des tickets (normaux)
	SLA           SLAStatsDTO           `json:"sla"`            // Statistiques SLA
	Performance   PerformanceStatsDTO   `json:"performance"`    // Statistiques de performance
	Alerts        []AlertDTO            `json:"alerts"`         // Alertes en cours
	Period        string                `json:"period"`         // Période (week, month, etc.)
	Users         UserStatsDTO          `json:"users"`           // Statistiques des utilisateurs
	Assets        AssetStatsDTO         `json:"assets"`         // Statistiques des actifs
	WorkedHours   WorkedHoursStatsDTO    `json:"worked_hours"`    // Heures travaillées (API, non affichées au board)
	Message       string                `json:"message,omitempty"` // Ex: "Aucun département associé à votre compte"
	TicketInternes *TicketInternalStatsDTO `json:"ticket_internes,omitempty"` // Stats tickets internes (rempli seulement si l'utilisateur a les permissions tickets_internes)
}

// TicketInternalStatsDTO représente les statistiques des tickets internes pour le tableau de bord
type TicketInternalStatsDTO struct {
	Total    int            `json:"total"`
	ByStatus map[string]int `json:"by_status"`
	ByPriority map[string]int `json:"by_priority,omitempty"`
	Open     int            `json:"open"`
	Closed   int            `json:"closed"`
}

// WorkedHoursStatsDTO représente les heures travaillées (données conservées, non affichées au board)
type WorkedHoursStatsDTO struct {
	TotalMinutes int     `json:"total_minutes"`
	TotalHours   float64 `json:"total_hours"`
	Period       string  `json:"period"`
}

// AssetStatsDTO représente les statistiques des actifs
type AssetStatsDTO struct {
	Total      int            `json:"total"`       // Nombre total d'actifs
	ByStatus   map[string]int `json:"by_status"`   // Par statut
	ByCategory map[string]int `json:"by_category"` // Par catégorie
}

// TicketStatsDTO représente les statistiques des tickets
type TicketStatsDTO struct {
	Total                 int            `json:"total"`                   // Nombre total de tickets
	ByCategory            map[string]int `json:"by_category"`             // Par catégorie
	ByStatus              map[string]int `json:"by_status"`               // Par statut
	ByPriority            map[string]int `json:"by_priority"`             // Par priorité
	AverageResolutionTime float64        `json:"average_resolution_time"` // Temps moyen de résolution en minutes
	Delayed               int            `json:"delayed"`                 // Tickets en retard
	Open                  int            `json:"open"`                    // Tickets ouverts
	Closed                int            `json:"closed"`                  // Tickets fermés
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
	TotalTimeSpent      int     `json:"total_time_spent"`      // Temps total passé en minutes
	AverageEfficiency   float64 `json:"average_efficiency"`    // Efficacité moyenne en %
	AverageProductivity float64 `json:"average_productivity"`  // Productivité moyenne (tickets/heure)
	TotalTicketsTreated int     `json:"total_tickets_treated"` // Nombre total de tickets traités
}

// AlertDTO représente une alerte
type AlertDTO struct {
	Type      string    `json:"type"`               // delay_alert, budget_alert, etc.
	Severity  string    `json:"severity"`           // low, medium, high, critical
	Message   string    `json:"message"`            // Message de l'alerte
	LinkURL   string    `json:"link_url,omitempty"` // URL vers la ressource (optionnel)
	CreatedAt time.Time `json:"created_at"`
}

// TicketCountReportDTO représente le rapport de nombre de tickets par période
type TicketCountReportDTO struct {
	Period    string               `json:"period"`    // week, month, etc.
	Count     int                  `json:"count"`     // Nombre total
	Breakdown []PeriodBreakdownDTO `json:"breakdown"` // Répartition par sous-période
}

// PeriodBreakdownDTO représente la répartition par sous-période (alignée sur les statuts en base)
type PeriodBreakdownDTO struct {
	Date       time.Time `json:"date"`        // Date de la période
	Count      int       `json:"count"`      // Nombre total de tickets pour cette période
	Open       int       `json:"open"`       // ouvert
	InProgress int       `json:"in_progress"` // en_cours
	Pending    int       `json:"pending"`    // en_attente
	Resolved   int       `json:"resolved"`   // resolu (validé)
	Closed     int       `json:"closed"`     // cloture (fermé)
}

// TicketTypeDistributionDTO représente la répartition des tickets par type
type TicketTypeDistributionDTO struct {
	Incidents      int `json:"incidents"`      // Nombre d'incidents
	Demandes       int `json:"demandes"`       // Nombre de demandes
	Changements    int `json:"changements"`    // Nombre de changements
	Developpements int `json:"developpements"` // Nombre de développements
	Assistance     int `json:"assistance"`     // Nombre d'assistances
	Support        int `json:"support"`         // Nombre de supports
}

// AverageResolutionTimeDTO représente le temps moyen de résolution
type AverageResolutionTimeDTO struct {
	AverageTime int                `json:"average_time"` // Temps moyen en minutes
	Unit        string             `json:"unit"`         // minutes
	Breakdown   map[string]float64 `json:"breakdown"`    // Répartition par catégorie/priorité
}

// WorkloadByAgentDTO représente la charge de travail par agent
type WorkloadByAgentDTO struct {
	UserID        uint     `json:"user_id"`
	User          *UserDTO `json:"user,omitempty"`
	TicketCount   int      `json:"ticket_count"`     // Nombre total de tickets assignés
	ResolvedCount int      `json:"resolved_count"`   // Nombre de tickets résolus
	InProgressCount int    `json:"in_progress_count"` // Nombre de tickets en cours
	PendingCount  int      `json:"pending_count"`    // Nombre de tickets en attente
	OpenCount     int      `json:"open_count"`       // Nombre de tickets ouverts
	DelayedCount  int      `json:"delayed_count"`    // Nombre de tickets en retard
	AverageTime   float64  `json:"average_time"`      // Temps moyen de résolution en minutes
	TotalTime     int      `json:"total_time"`       // Temps total passé en minutes
	Efficiency    float64  `json:"efficiency"`       // Efficacité en % (résolus / total)
}

// CustomReportRequest représente la requête pour un rapport personnalisé
type CustomReportRequest struct {
	Metrics   []string       `json:"metrics" binding:"required"` // Métriques à inclure (obligatoire)
	Period    string         `json:"period" binding:"required"`  // Période (obligatoire)
	StartDate *time.Time     `json:"start_date,omitempty"`       // Date de début (optionnel)
	EndDate   *time.Time     `json:"end_date,omitempty"`         // Date de fin (optionnel)
	Filters   map[string]any `json:"filters,omitempty"`          // Filtres personnalisés (optionnel)
	Format    string         `json:"format,omitempty"`           // Format d'export (pdf, excel, csv) (optionnel)
}

// DelayedTicketDTO représente un ticket en retard
type DelayedTicketDTO struct {
	TicketID     uint       `json:"ticket_id"`
	Ticket       *TicketDTO `json:"ticket,omitempty"`
	ExpectedDate time.Time  `json:"expected_date"` // Date attendue de résolution
	DelayedBy    int        `json:"delayed_by"`    // Nombre de jours de retard
	Priority     string     `json:"priority"`      // Priorité du ticket
	Category     string     `json:"category"`      // Catégorie du ticket
}

// IndividualPerformanceReportDTO représente le rapport de performance individuel
type IndividualPerformanceReportDTO struct {
	UserID                uint                   `json:"user_id"`
	User                  *UserDTO               `json:"user,omitempty"`
	Period                string                 `json:"period"`
	TicketsTreated        int                    `json:"tickets_treated"`         // Nombre de tickets traités
	AverageResolutionTime float64                `json:"average_resolution_time"` // Temps moyen de résolution en minutes
	Efficiency            float64                `json:"efficiency"`              // Efficacité en %
	Productivity          float64                `json:"productivity"`            // Productivité (tickets/heure)
	TotalTimeSpent        int                    `json:"total_time_spent"`        // Temps total passé en minutes
	Breakdown             map[string]interface{} `json:"breakdown,omitempty"`     // Répartition détaillée
}

// AssetReportDTO représente le résumé des actifs
type AssetReportDTO struct {
	Period     string         `json:"period"`
	Total      int            `json:"total"`
	ByStatus   map[string]int `json:"by_status"`
	ByCategory map[string]int `json:"by_category"`
}

// KnowledgeReportDTO représente le résumé de la base de connaissances
type KnowledgeReportDTO struct {
	Period     string         `json:"period"`
	Total      int            `json:"total"`
	Published  int            `json:"published"`
	Draft      int            `json:"draft"`
	ByCategory map[string]int `json:"by_category"`
}