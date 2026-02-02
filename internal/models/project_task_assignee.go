package models

import (
	"time"
)

// ProjectTaskAssignee permet l'assignation multiple sur une tâche (optionnel, style ticket_assignees)
// Table: project_task_assignees
type ProjectTaskAssignee struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	ProjectTaskID  uint      `gorm:"not null;uniqueIndex:idx_task_assignee" json:"project_task_id"`
	UserID         uint      `gorm:"not null;uniqueIndex:idx_task_assignee" json:"user_id"`
	IsLead         bool      `gorm:"default:false" json:"is_lead"`
	CreatedAt      time.Time `json:"created_at"`

	ProjectTask *ProjectTask `gorm:"foreignKey:ProjectTaskID" json:"-"`
	User        *User        `gorm:"foreignKey:UserID" json:"-"`
}

// TableName spécifie le nom de la table
func (ProjectTaskAssignee) TableName() string {
	return "project_task_assignees"
}
