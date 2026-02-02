package models

import "time"

// TicketAssignee représente l'association entre un ticket et un utilisateur assigné
// Table: ticket_assignees
type TicketAssignee struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	TicketID  uint      `gorm:"not null;index;uniqueIndex:idx_ticket_assignee" json:"ticket_id"`
	UserID    uint      `gorm:"not null;index;uniqueIndex:idx_ticket_assignee" json:"user_id"`
	IsLead    bool      `gorm:"default:false;index" json:"is_lead"`
	CreatedAt time.Time `json:"created_at"`

	// Relations
	Ticket Ticket `gorm:"foreignKey:TicketID;constraint:OnDelete:CASCADE" json:"ticket,omitempty"`
	User   User   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
}

// TableName spécifie le nom de la table
func (TicketAssignee) TableName() string {
	return "ticket_assignees"
}
