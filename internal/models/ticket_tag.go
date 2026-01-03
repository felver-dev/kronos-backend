package models

import (
	"time"
)

// TicketTag représente un tag pour les tickets
// Table: ticket_tags
type TicketTag struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"name"`
	Color     string    `gorm:"type:varchar(7)" json:"color,omitempty"` // Code couleur hexadécimal (ex: #FF5733)
	CreatedAt time.Time `json:"created_at"`

	// Relations HasMany (définies dans les autres modèles)
	// Assignments []TicketTagAssignment `gorm:"foreignKey:TagID" json:"-"`
}

// TableName spécifie le nom de la table
func (TicketTag) TableName() string {
	return "ticket_tags"
}

// TicketTagAssignment représente l'association entre un ticket et un tag (table de liaison)
// Table: ticket_tag_assignments
type TicketTagAssignment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	TicketID  uint      `gorm:"not null;uniqueIndex:idx_ticket_tag" json:"ticket_id"`
	TagID     uint      `gorm:"not null;uniqueIndex:idx_ticket_tag" json:"tag_id"`
	CreatedAt time.Time `json:"created_at"`

	// Relations
	Ticket Ticket    `gorm:"foreignKey:TicketID;constraint:OnDelete:CASCADE" json:"ticket,omitempty"`
	Tag    TicketTag `gorm:"foreignKey:TagID;constraint:OnDelete:CASCADE" json:"tag,omitempty"`
}

// TableName spécifie le nom de la table
func (TicketTagAssignment) TableName() string {
	return "ticket_tag_assignments"
}
