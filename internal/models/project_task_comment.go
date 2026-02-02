package models

import (
	"time"
)

// ProjectTaskComment représente un commentaire sur une tâche de projet
// Table: project_task_comments
type ProjectTaskComment struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	ProjectTaskID uint      `gorm:"not null;index" json:"project_task_id"`
	UserID        uint      `gorm:"not null;index" json:"user_id"`
	Comment       string    `gorm:"type:text;not null" json:"comment"`
	IsInternal    bool      `gorm:"default:false" json:"is_internal"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	ProjectTask *ProjectTask `gorm:"foreignKey:ProjectTaskID" json:"-"`
	User        *User        `gorm:"foreignKey:UserID" json:"-"`
}

// TableName spécifie le nom de la table
func (ProjectTaskComment) TableName() string {
	return "project_task_comments"
}
