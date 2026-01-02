package dto

import "time"

// IncidentDTO représente un incident dans les réponses API
type IncidentDTO struct {
	ID             uint       `json:"id"`
	TicketID       uint       `json:"ticket_id"`
	Ticket         *TicketDTO `json:"ticket,omitempty"`          // Ticket associé (optionnel)
	Impact         string     `json:"impact"`                    // low, medium, high, critical
	Urgency        string     `json:"urgency"`                   // low, medium, high, critical
	ResolutionTime *int       `json:"resolution_time,omitempty"` // Temps de résolution en minutes (optionnel)
	ResolvedAt     *time.Time `json:"resolved_at,omitempty"`     // Date de résolution (optionnel)
	LinkedAssets   []AssetDTO `json:"linked_assets,omitempty"`   // Actifs liés (optionnel)
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// CreateIncidentRequest représente la requête de création d'un incident
type CreateIncidentRequest struct {
	TicketID uint   `json:"ticket_id" binding:"required"`                              // ID du ticket (obligatoire)
	Impact   string `json:"impact" binding:"required,oneof=low medium high critical"`  // Impact (obligatoire)
	Urgency  string `json:"urgency" binding:"required,oneof=low medium high critical"` // Urgence (obligatoire)
}

// UpdateIncidentRequest représente la requête de mise à jour d'un incident
type UpdateIncidentRequest struct {
	Impact  string `json:"impact,omitempty" binding:"omitempty,oneof=low medium high critical"`  // Impact (optionnel)
	Urgency string `json:"urgency,omitempty" binding:"omitempty,oneof=low medium high critical"` // Urgence (optionnel)
}

// QualifyIncidentRequest représente la requête de qualification d'un incident (impact/urgence)
type QualifyIncidentRequest struct {
	Impact  string `json:"impact" binding:"required,oneof=low medium high critical"`  // Impact (obligatoire)
	Urgency string `json:"urgency" binding:"required,oneof=low medium high critical"` // Urgence (obligatoire)
}

// LinkAssetRequest représente la requête de liaison d'un actif à un incident
type LinkAssetRequest struct {
	AssetID uint `json:"asset_id" binding:"required"` // ID de l'actif à lier (obligatoire)
}

// ResolutionTimeDTO représente le temps de résolution d'un incident
type ResolutionTimeDTO struct {
	ResolutionTime int        `json:"resolution_time"`       // Temps de résolution en minutes
	Unit           string     `json:"unit"`                  // "minutes"
	StartedAt      time.Time  `json:"started_at"`            // Date de début
	ResolvedAt     *time.Time `json:"resolved_at,omitempty"` // Date de résolution (optionnel)
}
