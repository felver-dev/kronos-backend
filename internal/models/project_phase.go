package models

import (
	"time"
)

// ProjectPhase représente une étape (phase) d'un projet
// Table: project_phases
type ProjectPhase struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	ProjectID    uint       `gorm:"not null;index" json:"project_id"`
	Name         string     `gorm:"type:varchar(255);not null" json:"name"`
	Description  string     `gorm:"type:text" json:"description,omitempty"`
	DisplayOrder int        `gorm:"default:0" json:"display_order"`
	StartDate    *time.Time `gorm:"type:date" json:"start_date,omitempty"`
	EndDate      *time.Time `gorm:"type:date" json:"end_date,omitempty"`
	Status       string     `gorm:"type:varchar(50);default:'not_started';index" json:"status"` // not_started, in_progress, done, cancelled
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`

	Project *Project `gorm:"foreignKey:ProjectID" json:"-"`
}

// TableName spécifie le nom de la table
func (ProjectPhase) TableName() string {
	return "project_phases"
}
