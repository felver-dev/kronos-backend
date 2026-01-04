package dto

import "time"

// SLADTO représente un SLA (Service Level Agreement) dans les réponses API
type SLADTO struct {
	ID             uint      `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description,omitempty"`
	TicketCategory string    `json:"ticket_category"`    // incident, demande, changement, developpement
	Priority       *string   `json:"priority,omitempty"` // low, medium, high, critical (nil = tous)
	TargetTime     int       `json:"target_time"`        // Temps cible en minutes
	Unit           string    `json:"unit"`               // minutes, hours, days
	IsActive       bool      `json:"is_active"`          // Si le SLA est actif
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// CreateSLARequest représente la requête de création d'un SLA
type CreateSLARequest struct {
	Name           string  `json:"name" binding:"required"`                                                            // Nom (obligatoire)
	Description    string  `json:"description,omitempty"`                                                              // Description (optionnel)
	TicketCategory string  `json:"ticket_category" binding:"required,oneof=incident demande changement developpement"` // Catégorie (obligatoire)
	Priority       *string `json:"priority,omitempty" binding:"omitempty,oneof=low medium high critical"`              // Priorité (optionnel)
	TargetTime     int     `json:"target_time" binding:"required,min=1"`                                               // Temps cible en minutes (obligatoire, min 1)
	Unit           string  `json:"unit,omitempty" binding:"omitempty,oneof=minutes hours days"`                        // Unité (optionnel, défaut: minutes)
	IsActive       bool    `json:"is_active,omitempty"`                                                                // Statut actif (optionnel, défaut: true)
}

// UpdateSLARequest représente la requête de mise à jour d'un SLA
type UpdateSLARequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	TargetTime  *int   `json:"target_time,omitempty" binding:"omitempty,min=1"`
	Unit        string `json:"unit,omitempty" binding:"omitempty,oneof=minutes hours days"`
	IsActive    *bool  `json:"is_active,omitempty"`
}

// TicketSLAStatusDTO représente le statut SLA d'un ticket
type TicketSLAStatusDTO struct {
	SLAID       uint       `json:"sla_id"`                // ID du SLA appliqué
	SLA         *SLADTO    `json:"sla,omitempty"`         // SLA (optionnel)
	TargetTime  time.Time  `json:"target_time"`           // Date/heure cible
	ElapsedTime int        `json:"elapsed_time"`          // Temps écoulé en minutes
	Remaining   int        `json:"remaining"`             // Temps restant en minutes (peut être négatif)
	Status      string     `json:"status"`                // on_time, at_risk, violated
	ViolatedAt  *time.Time `json:"violated_at,omitempty"` // Date de violation (optionnel)
}

// SLAComplianceDTO représente les métriques de conformité d'un SLA
type SLAComplianceDTO struct {
	SLAID          uint    `json:"sla_id"`
	SLA            *SLADTO `json:"sla,omitempty"`
	ComplianceRate float64 `json:"compliance_rate"` // Taux de conformité en %
	TotalTickets   int     `json:"total_tickets"`   // Nombre total de tickets
	Compliant      int     `json:"compliant"`       // Nombre de tickets conformes
	Violations     int     `json:"violations"`      // Nombre de violations
}

// SLAViolationDTO représente une violation de SLA
type SLAViolationDTO struct {
	ID            uint       `json:"id"`
	TicketID      uint       `json:"ticket_id"`
	Ticket        *TicketDTO `json:"ticket,omitempty"`
	SLAID         uint       `json:"sla_id"`
	SLA           *SLADTO    `json:"sla,omitempty"`
	ViolationTime int        `json:"violation_time"` // Temps de violation en minutes
	Unit          string     `json:"unit"`           // minutes
	ViolatedAt    time.Time  `json:"violated_at"`    // Date de violation
}

// OverallSLAComplianceDTO représente la conformité globale des SLA
type OverallSLAComplianceDTO struct {
	OverallCompliance float64            `json:"overall_compliance"` // Conformité globale en %
	ByCategory        map[string]float64 `json:"by_category"`        // Conformité par catégorie
	ByPriority        map[string]float64 `json:"by_priority"`        // Conformité par priorité
	TotalTickets      int                `json:"total_tickets"`
	TotalViolations   int                `json:"total_violations"`
}

// SLAComplianceReportDTO représente un rapport de conformité des SLA
type SLAComplianceReportDTO struct {
	OverallCompliance float64            `json:"overall_compliance"` // Conformité globale en %
	ByCategory        map[string]float64 `json:"by_category"`        // Conformité par catégorie
	ByPriority        map[string]float64 `json:"by_priority"`        // Conformité par priorité
	TotalTickets      int                `json:"total_tickets"`
	TotalViolations   int                `json:"total_violations"`
	Period            string             `json:"period"`      // Période analysée
	GeneratedAt       time.Time          `json:"generated_at"` // Date de génération
}
