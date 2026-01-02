package dto

import "time"

// ServiceRequestDTO représente une demande de service dans les réponses API
type ServiceRequestDTO struct {
	ID                uint                   `json:"id"`
	TicketID          uint                   `json:"ticket_id"`
	Ticket            *TicketDTO             `json:"ticket,omitempty"`             // Ticket associé (optionnel)
	TypeID            uint                   `json:"type_id"`                      // ID du type de demande
	Type              *ServiceRequestTypeDTO `json:"type,omitempty"`               // Type de demande (optionnel)
	Deadline          *time.Time             `json:"deadline,omitempty"`           // Date limite (optionnel)
	Validated         bool                   `json:"validated"`                    // Si la demande a été validée
	ValidatedBy       *uint                  `json:"validated_by,omitempty"`       // ID du validateur (optionnel)
	ValidatedAt       *time.Time             `json:"validated_at,omitempty"`       // Date de validation (optionnel)
	ValidationComment string                 `json:"validation_comment,omitempty"` // Commentaire de validation (optionnel)
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
}

// ServiceRequestTypeDTO représente un type de demande de service paramétrable
type ServiceRequestTypeDTO struct {
	ID              uint   `json:"id"`
	Name            string `json:"name"` // Nom du type (ex: "Installation", "Configuration")
	Description     string `json:"description,omitempty"`
	DefaultDeadline int    `json:"default_deadline"` // Délai par défaut en heures
	IsActive        bool   `json:"is_active"`        // Si le type est actif
}

// CreateServiceRequestRequest représente la requête de création d'une demande de service
type CreateServiceRequestRequest struct {
	TicketID uint    `json:"ticket_id" binding:"required"` // ID du ticket (obligatoire)
	TypeID   uint    `json:"type_id" binding:"required"`   // ID du type (obligatoire)
	Deadline *string `json:"deadline,omitempty"`           // Date limite format "2006-01-02" (optionnel)
}

// UpdateServiceRequestRequest représente la requête de mise à jour d'une demande de service
type UpdateServiceRequestRequest struct {
	TypeID   *uint   `json:"type_id,omitempty"`  // ID du type (optionnel)
	Deadline *string `json:"deadline,omitempty"` // Date limite format "2006-01-02" (optionnel)
}

// ValidateServiceRequestRequest représente la requête de validation d'une demande de service
type ValidateServiceRequestRequest struct {
	Validated bool   `json:"validated" binding:"required"` // true pour valider, false pour invalider
	Comment   string `json:"comment,omitempty"`            // Commentaire de validation (optionnel)
}

// CreateServiceRequestTypeRequest représente la requête de création d'un type de demande
type CreateServiceRequestTypeRequest struct {
	Name            string `json:"name" binding:"required"`                   // Nom (obligatoire)
	Description     string `json:"description,omitempty"`                     // Description (optionnel)
	DefaultDeadline int    `json:"default_deadline" binding:"required,min=1"` // Délai par défaut en heures (obligatoire, min 1)
}

// UpdateServiceRequestTypeRequest représente la requête de mise à jour d'un type de demande
type UpdateServiceRequestTypeRequest struct {
	Name            string `json:"name,omitempty"`
	Description     string `json:"description,omitempty"`
	DefaultDeadline *int   `json:"default_deadline,omitempty" binding:"omitempty,min=1"`
	IsActive        *bool  `json:"is_active,omitempty"` // Statut actif (optionnel)
}

// DeadlineStatusDTO représente le statut d'une deadline
type DeadlineStatusDTO struct {
	Deadline  time.Time `json:"deadline"`   // Date limite
	Remaining int       `json:"remaining"`  // Jours restants (peut être négatif si dépassé)
	Unit      string    `json:"unit"`       // "days"
	IsOverdue bool      `json:"is_overdue"` // Si la deadline est dépassée
}
