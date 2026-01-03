package models

import "time"

// TicketHistory représente une entrée dans l'historique d'un ticket
// Table: ticket_history
type TicketHistory struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	TicketID    uint      `gorm:"not null;index" json:"ticket_id"`
	UserID      uint      `gorm:"not null;index" json:"user_id"`
	Action      string    `gorm:"type:varchar(100);not null;index" json:"action"` // created, updated, status_changed, assigned, etc.
	FieldName   string    `gorm:"type:varchar(100)" json:"field_name,omitempty"`  // Nom du champ modifié (optionnel)
	OldValue    string    `gorm:"type:text" json:"old_value,omitempty"`           // Ancienne valeur (optionnel)
	NewValue    string    `gorm:"type:text" json:"new_value,omitempty"`           // Nouvelle valeur (optionnel)
	Description string    `gorm:"type:text" json:"description,omitempty"`         // Description de l'action (optionnel)
	CreatedAt   time.Time `gorm:"index" json:"created_at"`                        // Date de l'action

	// Relations
	Ticket Ticket `gorm:"foreignKey:TicketID" json:"-"`  // Ticket associé
	User   User   `gorm:"foreignKey:UserID" json:"user"` // Utilisateur qui a effectué l'action
}

// TableName spécifie le nom de la table
func (TicketHistory) TableName() string {
	return "ticket_history"
}
