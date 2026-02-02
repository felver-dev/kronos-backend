package models

import (
	"time"
)

// ProjectBudgetExtension enregistre une extension du budget temps d'un projet (temps ajouté + justification).
// StartDate et EndDate délimitent la période de cette extension (phase).
// Table: project_budget_extensions
type ProjectBudgetExtension struct {
	ID                uint       `gorm:"primaryKey" json:"id"`
	ProjectID         uint       `gorm:"not null;index" json:"project_id"`
	AdditionalMinutes int        `gorm:"not null" json:"additional_minutes"` // Minutes ajoutées au budget
	Justification     string     `gorm:"type:text;not null" json:"justification"`
	StartDate         *time.Time `gorm:"type:date" json:"start_date,omitempty"` // Début de la période de l'extension
	EndDate           *time.Time `gorm:"type:date" json:"end_date,omitempty"`   // Fin de la période de l'extension
	CreatedByID       *uint      `gorm:"index" json:"created_by_id,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`

	// Relations
	Project   Project `gorm:"foreignKey:ProjectID" json:"-"`
	CreatedBy *User   `gorm:"foreignKey:CreatedByID" json:"created_by,omitempty"`
}

// TableName spécifie le nom de la table
func (ProjectBudgetExtension) TableName() string {
	return "project_budget_extensions"
}
