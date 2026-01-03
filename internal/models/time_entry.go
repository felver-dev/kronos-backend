package models

import (
	"time"

	"gorm.io/gorm"
)

// TimeEntry représente une entrée de temps (temps passé sur un ticket)
// Table: time_entries
type TimeEntry struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	TicketID      uint           `gorm:"not null;index" json:"ticket_id"`
	UserID        uint           `gorm:"not null;index" json:"user_id"`
	TimeSpent     int            `gorm:"not null" json:"time_spent"`           // Temps passé en minutes
	Date          time.Time      `gorm:"type:date;not null;index" json:"date"` // Date de l'entrée
	Description   string         `gorm:"type:text" json:"description,omitempty"`
	Validated     bool           `gorm:"default:false;index" json:"validated"`   // Si l'entrée a été validée
	ValidatedByID *uint          `gorm:"index" json:"validated_by_id,omitempty"` // ID du validateur (optionnel)
	ValidatedAt   *time.Time     `json:"validated_at,omitempty"`                 // Date de validation (optionnel)
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"` // Soft delete

	// Relations
	Ticket      Ticket `gorm:"foreignKey:TicketID;constraint:OnDelete:CASCADE" json:"-"` // Ticket associé
	User        User   `gorm:"foreignKey:UserID" json:"user,omitempty"`                  // Utilisateur
	ValidatedBy *User  `gorm:"foreignKey:ValidatedByID" json:"validated_by,omitempty"`   // Validateur (optionnel)
}

// TableName spécifie le nom de la table
func (TimeEntry) TableName() string {
	return "time_entries"
}
