package models

import (
	"time"
)

// ProjectTaskHistory enregistre l'historique des changements d'une tâche (optionnel)
// Table: project_task_history
type ProjectTaskHistory struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	ProjectTaskID uint      `gorm:"not null;index" json:"project_task_id"`
	UserID        uint      `gorm:"not null;index" json:"user_id"`
	Action        string    `gorm:"type:varchar(100)" json:"action"`
	FieldName     string    `gorm:"type:varchar(100)" json:"field_name,omitempty"`
	OldValue      string    `gorm:"type:text" json:"old_value,omitempty"`
	NewValue      string    `gorm:"type:text" json:"new_value,omitempty"`
	CreatedAt     time.Time `json:"created_at"`

	ProjectTask *ProjectTask `gorm:"foreignKey:ProjectTaskID" json:"-"`
	User        *User        `gorm:"foreignKey:UserID" json:"-"`
}

// TableName spécifie le nom de la table
func (ProjectTaskHistory) TableName() string {
	return "project_task_history"
}
