package models

import (
	"time"

	"gorm.io/gorm"
)

// TicketComment représente un commentaire sur un ticket
// Table: ticket_comments
type TicketComment struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	TicketID   uint           `gorm:"not null;index" json:"ticket_id"`
	UserID     uint           `gorm:"not null;index" json:"user_id"`
	Comment    string         `gorm:"type:text;not null" json:"comment"`
	IsInternal bool           `gorm:"default:false" json:"is_internal"` // Commentaire interne (visible uniquement par l'IT)
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"` // Soft delete

	// Relations
	Ticket Ticket `gorm:"foreignKey:TicketID" json:"-"`  // Ticket associé
	User   User   `gorm:"foreignKey:UserID" json:"user"` // Utilisateur auteur
}

// TableName spécifie le nom de la table
func (TicketComment) TableName() string {
	return "ticket_comments"
}
