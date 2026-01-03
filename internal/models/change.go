package models

import (
	"time"
)

// Change représente un changement (extension d'un ticket)
// Table: changes
type Change struct {
	ID               uint       `gorm:"primaryKey" json:"id"`
	TicketID         uint       `gorm:"uniqueIndex;not null;index" json:"ticket_id"` // Relation 1:1 avec Ticket
	Risk             string     `gorm:"type:varchar(50);not null;index" json:"risk"` // low, medium, high, critical
	RiskDescription  string     `gorm:"type:text" json:"risk_description,omitempty"` // Description du risque (optionnel)
	ResponsibleID    *uint      `gorm:"index" json:"responsible_id,omitempty"`        // ID du responsable (optionnel)
	Result           string     `gorm:"type:varchar(50)" json:"result,omitempty"`      // success, partial, failed, rolled_back
	ResultDescription string    `gorm:"type:text" json:"result_description,omitempty"` // Description du résultat (optionnel)
	ResultDate       *time.Time `json:"result_date,omitempty"`                         // Date du résultat (optionnel)
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`

	// Relations
	Ticket     Ticket  `gorm:"foreignKey:TicketID;constraint:OnDelete:CASCADE" json:"ticket,omitempty"` // Ticket associé (1:1)
	Responsible *User  `gorm:"foreignKey:ResponsibleID" json:"responsible,omitempty"`                     // Responsable du changement (optionnel)
}

// TableName spécifie le nom de la table
func (Change) TableName() string {
	return "changes"
}

