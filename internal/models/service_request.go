package models

import (
	"time"
)

// ServiceRequestType représente un type de demande de service paramétrable
// Table: service_request_types
type ServiceRequestType struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	Name           string     `gorm:"type:varchar(100);not null" json:"name"`
	Description    string     `gorm:"type:text" json:"description,omitempty"`
	DefaultDeadline int       `gorm:"type:int" json:"default_deadline"` // Délai par défaut en heures
	IsActive       bool       `gorm:"default:true;index" json:"is_active"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	CreatedByID    *uint      `gorm:"index" json:"-"`
	CreatedBy      *User      `gorm:"foreignKey:CreatedByID" json:"-"`

	// Relations HasMany
	ServiceRequests []ServiceRequest `gorm:"foreignKey:TypeID" json:"-"`
}

// TableName spécifie le nom de la table
func (ServiceRequestType) TableName() string {
	return "service_request_types"
}

// ServiceRequest représente une demande de service (extension d'un ticket)
// Table: service_requests
type ServiceRequest struct {
	ID                uint       `gorm:"primaryKey" json:"id"`
	TicketID          uint       `gorm:"uniqueIndex;not null;index" json:"ticket_id"` // Relation 1:1 avec Ticket
	TypeID            uint       `gorm:"not null;index" json:"type_id"`
	Deadline          *time.Time `gorm:"index" json:"deadline,omitempty"`              // Date limite (optionnel)
	Validated         bool       `gorm:"default:false" json:"validated"`                // Si la demande a été validée
	ValidatedByID     *uint      `gorm:"index" json:"validated_by_id,omitempty"`        // ID du validateur (optionnel)
	ValidatedAt       *time.Time `json:"validated_at,omitempty"`                       // Date de validation (optionnel)
	ValidationComment string     `gorm:"type:text" json:"validation_comment,omitempty"` // Commentaire de validation (optionnel)
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`

	// Relations - GORM utilisera automatiquement les champs existants
	Ticket      Ticket            `gorm:"constraint:OnDelete:CASCADE" json:"ticket,omitempty"` // Ticket associé (1:1)
	Type        ServiceRequestType `gorm:"foreignKey:TypeID" json:"type,omitempty"`                              // Type de demande
	ValidatedBy  *User            `gorm:"foreignKey:ValidatedByID" json:"validated_by,omitempty"`                 // Validateur (optionnel)
}

// TableName spécifie le nom de la table
func (ServiceRequest) TableName() string {
	return "service_requests"
}

