package models

import (
	"time"
)

// ProjectFunction représente une fonction (rôle projet).
// Type: "direction" (Chef de projet, Lead, Tech Lead…) ou "execution" (Dev, Testeur…).
// Un membre peut avoir plusieurs fonctions (ex. Dev + Lead).
// Table: project_functions
type ProjectFunction struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	ProjectID    *uint  `gorm:"index" json:"project_id,omitempty"` // NULL = catalogue global
	Name         string `gorm:"type:varchar(100);not null" json:"name"`
	Type         string `gorm:"column:function_type;type:varchar(20);default:execution" json:"type"` // "direction" | "execution" (colonne function_type pour éviter le mot réservé MySQL)
	DisplayOrder int    `gorm:"default:0" json:"display_order"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	Project *Project `gorm:"foreignKey:ProjectID" json:"-"`
}

// TableName spécifie le nom de la table
func (ProjectFunction) TableName() string {
	return "project_functions"
}
