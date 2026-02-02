package models

import (
	"time"

	"gorm.io/gorm"
)

// TicketSolution représente une solution documentée pour un ticket
// Table: ticket_solutions
type TicketSolution struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	TicketID    uint           `gorm:"not null;index" json:"ticket_id"`
	Solution    string         `gorm:"type:text;not null" json:"solution"` // Solution documentée (Markdown)
	CreatedByID uint           `gorm:"not null;index" json:"created_by_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"` // Soft delete

	// Relations
	Ticket   Ticket `gorm:"foreignKey:TicketID" json:"ticket,omitempty"`
	CreatedBy User  `gorm:"foreignKey:CreatedByID" json:"created_by,omitempty"`
}

// TableName spécifie le nom de la table
func (TicketSolution) TableName() string {
	return "ticket_solutions"
}
