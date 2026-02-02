package models

import (
	"time"

	"gorm.io/gorm"
)

// TimeEntry représente une entrée de temps (temps passé sur un ticket, un ticket interne ou une tâche de projet)
// Table: time_entries. Soit ticket_id, soit ticket_internal_id, soit project_task_id (l'un des trois).
type TimeEntry struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	TicketID         *uint          `gorm:"index" json:"ticket_id,omitempty"`             // Ticket normal (optionnel si ticket_internal_id)
	TicketInternalID *uint          `gorm:"index" json:"ticket_internal_id,omitempty"`     // Ticket interne (optionnel si ticket_id)
	ProjectTaskID    *uint          `gorm:"index" json:"project_task_id,omitempty"`       // Tâche de projet (optionnel)
	UserID         uint           `gorm:"not null;index" json:"user_id"`
	TimeSpent      int            `gorm:"not null" json:"time_spent"`
	Date           time.Time      `gorm:"type:date;not null;index" json:"date"`
	Description    string         `gorm:"type:text" json:"description,omitempty"`
	Validated      bool           `gorm:"default:false;index" json:"validated"`
	ValidatedByID  *uint          `gorm:"index" json:"validated_by_id,omitempty"`
	ValidatedAt    *time.Time     `json:"validated_at,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Ticket        *Ticket        `gorm:"foreignKey:TicketID;constraint:OnDelete:CASCADE" json:"-"`
	TicketInternal *TicketInternal `gorm:"foreignKey:TicketInternalID;constraint:OnDelete:CASCADE" json:"-"`
	ProjectTask   *ProjectTask   `gorm:"foreignKey:ProjectTaskID;constraint:OnDelete:CASCADE" json:"-"`
	User        User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ValidatedBy *User        `gorm:"foreignKey:ValidatedByID" json:"validated_by,omitempty"`
}

// TableName spécifie le nom de la table
func (TimeEntry) TableName() string {
	return "time_entries"
}
