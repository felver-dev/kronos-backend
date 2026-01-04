package dto

import "time"

// StatisticsOverviewDTO représente la vue d'ensemble des statistiques
type StatisticsOverviewDTO struct {
	Period      string              `json:"period"`
	Tickets     TicketStatsDTO      `json:"tickets"`
	SLA         SLAStatsDTO         `json:"sla"`
	Performance PerformanceStatsDTO `json:"performance"`
	Users       UserStatsDTO        `json:"users"`
}

// UserStatsDTO représente les statistiques des utilisateurs
type UserStatsDTO struct {
	Total        int            `json:"total"`
	Active       int            `json:"active"`
	ByRole       map[string]int `json:"by_role"`
	AverageTicketsPerUser float64 `json:"average_tickets_per_user"`
}

// WorkloadStatisticsDTO représente les statistiques de charge de travail
type WorkloadStatisticsDTO struct {
	Period        string    `json:"period"`
	UserID        *uint     `json:"user_id,omitempty"`
	TotalTickets  int       `json:"total_tickets"`
	AveragePerDay float64  `json:"average_per_day"`
	PeakDay       time.Time `json:"peak_day"`
	PeakDayCount  int       `json:"peak_day_count"`
	Distribution  []WorkloadDayDTO `json:"distribution"`
}

// WorkloadDayDTO représente la charge de travail pour un jour
type WorkloadDayDTO struct {
	Date  time.Time `json:"date"`
	Count int       `json:"count"`
}

// PerformanceStatisticsDTO représente les statistiques de performance
type PerformanceStatisticsDTO struct {
	Period                string  `json:"period"`
	AverageResolutionTime float64 `json:"average_resolution_time"` // en minutes
	SLACompliance         float64 `json:"sla_compliance"`          // en pourcentage
	Efficiency            float64 `json:"efficiency"`              // en pourcentage
	Productivity          float64 `json:"productivity"`            // en pourcentage
	FirstResponseTime     float64 `json:"first_response_time"`    // en minutes
}

// TrendsStatisticsDTO représente les statistiques de tendances
type TrendsStatisticsDTO struct {
	Metric   string          `json:"metric"`
	Period   string          `json:"period"`
	Trend    string          `json:"trend"`    // increasing, decreasing, stable
	Data     []TrendDataDTO  `json:"data"`
	Forecast []TrendDataDTO  `json:"forecast,omitempty"`
}

// TrendDataDTO représente une donnée de tendance
type TrendDataDTO struct {
	Date  time.Time `json:"date"`
	Value float64   `json:"value"`
}

// KPIStatisticsDTO représente les indicateurs de succès (KPI)
type KPIStatisticsDTO struct {
	Period              string  `json:"period"`
	TicketsRegistered   int     `json:"tickets_registered"`
	UtilizationRate     float64 `json:"utilization_rate"`      // en pourcentage
	ReportProduction    int     `json:"report_production"`     // nombre de rapports produits
	AverageSatisfaction float64 `json:"average_satisfaction"`  // en pourcentage (si disponible)
	ResponseTime        float64 `json:"response_time"`         // temps de réponse moyen en minutes
	ResolutionRate      float64 `json:"resolution_rate"`       // taux de résolution en pourcentage
}

