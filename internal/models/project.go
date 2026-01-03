package models

import (
	"time"
)

// Project représente un projet
// Table: projects
type Project struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	Name            string     `gorm:"type:varchar(255);not null" json:"name"`
	Description     string     `gorm:"type:text" json:"description,omitempty"`
	TotalBudgetTime *int       `gorm:"type:int" json:"total_budget_time,omitempty"`           // Budget temps total en minutes (optionnel)
	ConsumedTime    int        `gorm:"default:0" json:"consumed_time"`                        // Temps consommé en minutes (calculé)
	Status          string     `gorm:"type:varchar(50);default:'active';index" json:"status"` // active, completed, cancelled
	StartDate       *time.Time `gorm:"type:date" json:"start_date,omitempty"`
	EndDate         *time.Time `gorm:"type:date" json:"end_date,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	CreatedByID     *uint      `gorm:"index" json:"-"`
	CreatedBy       *User      `gorm:"foreignKey:CreatedByID" json:"-"`

	// Relations
	Tickets []Ticket `gorm:"many2many:ticket_projects;" json:"tickets,omitempty"` // Tickets associés (many-to-many)
}

// TableName spécifie le nom de la table
func (Project) TableName() string {
	return "projects"
}

// TicketProject représente l'association entre un ticket et un projet (table de liaison)
// Table: ticket_projects
type TicketProject struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	TicketID  uint      `gorm:"not null;uniqueIndex:idx_ticket_project" json:"ticket_id"`
	ProjectID uint      `gorm:"not null;uniqueIndex:idx_ticket_project" json:"project_id"`
	CreatedAt time.Time `json:"created_at"`

	// Relations
	Ticket  Ticket  `gorm:"foreignKey:TicketID;constraint:OnDelete:CASCADE" json:"ticket,omitempty"`
	Project Project `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"project,omitempty"`
}

// TableName spécifie le nom de la table
func (TicketProject) TableName() string {
	return "ticket_projects"
}
