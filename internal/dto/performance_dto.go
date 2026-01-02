package dto

// PerformanceDTO représente les métriques de performance d'un technicien
type PerformanceDTO struct {
	UserID               uint     `json:"user_id"`
	User                 *UserDTO `json:"user,omitempty"`
	TicketsTreated       int      `json:"tickets_treated"`         // Nombre de tickets traités
	TotalTimeSpent       int      `json:"total_time_spent"`        // Temps total passé en minutes
	AverageTimePerTicket float64  `json:"average_time_per_ticket"` // Temps moyen par ticket en minutes
	TotalEstimatedTime   int      `json:"total_estimated_time"`    // Temps estimé total en minutes
	Efficiency           float64  `json:"efficiency"`              // Efficacité en % (100 = parfait, <100 = en retard, >100 = en avance)
	Productivity         float64  `json:"productivity"`            // Productivité (tickets/heure)
	BudgetComplianceRate float64  `json:"budget_compliance_rate"`  // Taux de respect des budgets en %
	TicketsWithinBudget  int      `json:"tickets_within_budget"`   // Nombre de tickets dans le budget
	TicketsOverBudget    int      `json:"tickets_over_budget"`     // Nombre de tickets dépassant le budget
	TotalDelays          int      `json:"total_delays"`            // Nombre total de retards
	JustifiedDelays      int      `json:"justified_delays"`        // Nombre de retards justifiés
	UnjustifiedDelays    int      `json:"unjustified_delays"`      // Nombre de retards non justifiés
}

// EfficiencyDTO représente les métriques d'efficacité (estimé vs réel)
type EfficiencyDTO struct {
	UserID         uint     `json:"user_id"`
	User           *UserDTO `json:"user,omitempty"`
	Efficiency     float64  `json:"efficiency"`      // Efficacité en % (100 = parfait)
	EstimatedTotal int      `json:"estimated_total"` // Temps estimé total en minutes
	ActualTotal    int      `json:"actual_total"`    // Temps réel total en minutes
	Savings        int      `json:"savings"`         // Temps économisé en minutes (peut être négatif si retard)
	Overrun        int      `json:"overrun"`         // Temps de dépassement en minutes (si > 0)
}

// ProductivityDTO représente les métriques de productivité
type ProductivityDTO struct {
	UserID         uint     `json:"user_id"`
	User           *UserDTO `json:"user,omitempty"`
	Productivity   float64  `json:"productivity"`     // Tickets par heure
	TicketsPerHour float64  `json:"tickets_per_hour"` // Alias de productivity
	TicketsPerDay  float64  `json:"tickets_per_day"`  // Tickets par jour
	TicketsPerWeek float64  `json:"tickets_per_week"` // Tickets par semaine
}

// BudgetComplianceDTO représente les métriques de respect des budgets
type BudgetComplianceDTO struct {
	UserID                uint     `json:"user_id"`
	User                  *UserDTO `json:"user,omitempty"`
	ComplianceRate        float64  `json:"compliance_rate"`         // Taux de respect en %
	TotalTickets          int      `json:"total_tickets"`           // Nombre total de tickets
	WithinBudget          int      `json:"within_budget"`           // Tickets dans le budget
	OverBudget            int      `json:"over_budget"`             // Tickets dépassant le budget
	AverageOverrun        float64  `json:"average_overrun"`         // Dépassement moyen en minutes
	AverageOverrunPercent float64  `json:"average_overrun_percent"` // Dépassement moyen en %
}

// WorkloadDTO représente la charge de travail (réel vs théorique)
type WorkloadDTO struct {
	UserID              uint     `json:"user_id"`
	User                *UserDTO `json:"user,omitempty"`
	RealWorkload        int      `json:"real_workload"`        // Charge réelle en minutes
	TheoreticalWorkload int      `json:"theoretical_workload"` // Charge théorique en minutes
	Difference          int      `json:"difference"`           // Différence en minutes (peut être négative)
	Percentage          float64  `json:"percentage"`           // Pourcentage (100 = charge normale)
}

// PerformanceRankingDTO représente le classement d'un technicien
type PerformanceRankingDTO struct {
	Rank   int      `json:"rank"` // Position dans le classement (1 = meilleur)
	UserID uint     `json:"user_id"`
	User   *UserDTO `json:"user,omitempty"`
	Score  float64  `json:"score"`  // Score de performance
	Metric string   `json:"metric"` // Métrique utilisée (efficiency, productivity, etc.)
}

// TeamComparisonDTO représente la comparaison entre techniciens
type TeamComparisonDTO struct {
	Users   []PerformanceDTO `json:"users"`           // Liste des performances
	Average PerformanceDTO   `json:"average"`         // Moyenne de l'équipe
	Best    *PerformanceDTO  `json:"best,omitempty"`  // Meilleur technicien
	Worst   *PerformanceDTO  `json:"worst,omitempty"` // Technicien le plus en retard
	Metric  string           `json:"metric"`          // Métrique comparée
}
