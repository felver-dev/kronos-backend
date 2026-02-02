package dto

import "time"

// DelayDTO représente un retard détecté sur un ticket
type DelayDTO struct {
	ID              uint                   `json:"id"`
	TicketID        uint                   `json:"ticket_id"`
	Ticket          *TicketDTO             `json:"ticket,omitempty"` // Ticket concerné (optionnel)
	UserID          uint                   `json:"user_id"`          // Technicien en retard
	User            *UserDTO               `json:"user,omitempty"`
	EstimatedTime   int                    `json:"estimated_time"`          // Temps estimé en minutes
	ActualTime      int                    `json:"actual_time"`             // Temps réel en minutes
	DelayTime       int                    `json:"delay_time"`              // Retard en minutes (actual - estimated)
	DelayPercentage float64                `json:"delay_percentage"`        // Pourcentage de retard
	Status          string                 `json:"status"`                  // unjustified, pending, justified, rejected
	Justification   *DelayJustificationDTO `json:"justification,omitempty"` // Justification (optionnel)
	DetectedAt      time.Time              `json:"detected_at"`             // Date de détection
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}

// DelayJustificationDTO représente une justification de retard
type DelayJustificationDTO struct {
	ID                uint       `json:"id"`
	DelayID           uint       `json:"delay_id"`
	TicketID          *uint      `json:"ticket_id,omitempty"`
	TicketCode        string     `json:"ticket_code,omitempty"`
	TicketTitle       string     `json:"ticket_title,omitempty"`
	UserID            uint       `json:"user_id"` // Technicien qui justifie
	User              *UserDTO   `json:"user,omitempty"`
	Justification     string     `json:"justification"`          // Texte de justification
	Status            string     `json:"status"`                 // pending, validated, rejected
	ValidatedBy       *uint      `json:"validated_by,omitempty"` // ID du validateur
	ValidatedAt       *time.Time `json:"validated_at,omitempty"`
	ValidationComment string     `json:"validation_comment,omitempty"` // Commentaire du validateur
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

// CreateDelayJustificationRequest représente la requête de création d'une justification
type CreateDelayJustificationRequest struct {
	Justification string `json:"justification" binding:"required"` // Texte de justification (obligatoire)
}

// UpdateDelayJustificationRequest représente la requête de mise à jour d'une justification (avant validation)
type UpdateDelayJustificationRequest struct {
	Justification string `json:"justification" binding:"required"` // Nouveau texte de justification
}

// ValidateDelayJustificationRequest représente la requête de validation/rejet d'une justification
type ValidateDelayJustificationRequest struct {
	Validated *bool  `json:"validated,omitempty"`          // true pour valider, false pour rejeter
	Comment   string `json:"comment,omitempty"`            // Commentaire du validateur (optionnel)
}

// DelayStatusStatsDTO représente les statistiques de retards par statut
type DelayStatusStatsDTO struct {
	Unjustified int `json:"unjustified"` // Retards non justifiés
	Pending     int `json:"pending"`     // Justifications en attente
	Justified   int `json:"justified"`   // Retards justifiés validés
	Rejected    int `json:"rejected"`    // Justifications rejetées
}
