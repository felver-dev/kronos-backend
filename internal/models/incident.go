package models

import (
	"time"
)

// Incident représente un incident (extension d'un ticket)
// Table: incidents
type Incident struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	TicketID       uint       `gorm:"uniqueIndex;not null;index" json:"ticket_id"`    // Relation 1:1 avec Ticket
	Impact         string     `gorm:"type:varchar(50);not null;index" json:"impact"`  // low, medium, high, critical
	Urgency        string     `gorm:"type:varchar(50);not null;index" json:"urgency"` // low, medium, high, critical
	ResolutionTime *int       `gorm:"type:int" json:"resolution_time,omitempty"`      // Temps de résolution en minutes (calculé)
	ResolvedAt     *time.Time `json:"resolved_at,omitempty"`                          // Date de résolution (optionnel)
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`

	// Relations - GORM utilisera automatiquement le champ TicketID comme clé étrangère
	// Ne pas spécifier foreignKey pour éviter la duplication de colonne
	Ticket Ticket `gorm:"constraint:OnDelete:CASCADE" json:"ticket,omitempty"` // Ticket associé (1:1)
}

// TableName spécifie le nom de la table
func (Incident) TableName() string {
	return "incidents"
}
