package models

import (
	"time"
)

// SLA représente un Service Level Agreement (délai cible)
// Table: sla
type SLA struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	Name           string    `gorm:"type:varchar(255);not null" json:"name"`
	Description    string    `gorm:"type:text" json:"description,omitempty"`
	TicketCategory string    `gorm:"type:varchar(50);not null;index" json:"ticket_category"` // incident, demande, changement, developpement
	Priority       *string   `gorm:"type:varchar(50);index" json:"priority,omitempty"`       // low, medium, high, critical (nil = tous)
	TargetTime     int       `gorm:"not null" json:"target_time"`                            // Temps cible en minutes
	Unit           string    `gorm:"type:varchar(20);default:'minutes'" json:"unit"`         // minutes, hours, days
	IsActive       bool      `gorm:"default:true;index" json:"is_active"`                    // Si le SLA est actif
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	CreatedByID    *uint     `gorm:"index" json:"-"`
	CreatedBy      *User     `gorm:"foreignKey:CreatedByID" json:"-"`

	// Relations HasMany
	TicketSLAs []TicketSLA `gorm:"foreignKey:SLAID" json:"-"`
}

// TableName spécifie le nom de la table
func (SLA) TableName() string {
	return "sla"
}

// TicketSLA représente l'association entre un ticket et un SLA avec suivi de conformité
// Table: ticket_sla
type TicketSLA struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	TicketID      uint       `gorm:"uniqueIndex:idx_ticket_sla_ticket_id;not null" json:"ticket_id"` // Relation 1:1 avec Ticket
	SLAID         uint       `gorm:"not null;index" json:"sla_id"`
	TargetTime    time.Time  `gorm:"not null;index" json:"target_time"`                      // Date/heure cible
	ActualTime    *time.Time `json:"actual_time,omitempty"`                                  // Date/heure réelle de résolution (optionnel)
	Status        string     `gorm:"type:varchar(50);default:'on_time';index" json:"status"` // on_time, at_risk, violated
	ViolationTime *int       `gorm:"type:int" json:"violation_time,omitempty"`               // Temps de violation en minutes (si violé)
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`

	// Relations - GORM utilisera automatiquement les champs existants
	Ticket Ticket `gorm:"foreignKey:TicketID;constraint:OnDelete:CASCADE" json:"ticket,omitempty"` // Ticket associé (1:1)
	SLA    SLA    `gorm:"foreignKey:SLAID" json:"sla,omitempty"`                                   // SLA associé
}

// TableName spécifie le nom de la table
func (TicketSLA) TableName() string {
	return "ticket_sla"
}
