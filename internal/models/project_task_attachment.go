package models

import (
	"time"
)

// ProjectTaskAttachment représente une pièce jointe sur une tâche de projet
// Table: project_task_attachments
type ProjectTaskAttachment struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	ProjectTaskID uint      `gorm:"not null;index" json:"project_task_id"`
	UserID        uint      `gorm:"not null;index" json:"user_id"`
	FileName      string    `gorm:"type:varchar(255);not null" json:"file_name"`
	FilePath      string    `gorm:"type:varchar(500);not null" json:"file_path"`
	FileSize      int64     `gorm:"default:0" json:"file_size"`
	MimeType      string    `gorm:"type:varchar(100)" json:"mime_type,omitempty"`
	CreatedAt     time.Time `json:"created_at"`

	ProjectTask *ProjectTask `gorm:"foreignKey:ProjectTaskID" json:"-"`
	User        *User        `gorm:"foreignKey:UserID" json:"-"`
}

// TableName spécifie le nom de la table
func (ProjectTaskAttachment) TableName() string {
	return "project_task_attachments"
}
