package models

import (
	"time"
)

// ProjectPhaseMember représente un membre affecté à une seule étape (sans être membre du projet global)
// Table: project_phase_members
type ProjectPhaseMember struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	ProjectPhaseID    uint      `gorm:"not null;uniqueIndex:idx_phase_member" json:"project_phase_id"`
	UserID            uint      `gorm:"not null;uniqueIndex:idx_phase_member" json:"user_id"`
	ProjectFunctionID *uint     `gorm:"index" json:"project_function_id,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`

	ProjectPhase *ProjectPhase   `gorm:"foreignKey:ProjectPhaseID" json:"-"`
	User         *User           `gorm:"foreignKey:UserID" json:"-"`
	Function     *ProjectFunction `gorm:"foreignKey:ProjectFunctionID" json:"-"`
}

// TableName spécifie le nom de la table
func (ProjectPhaseMember) TableName() string {
	return "project_phase_members"
}
