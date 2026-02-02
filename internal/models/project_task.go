package models

import (
	"time"
)

// ProjectTask représente une tâche de projet (unité de travail dans une étape)
// Code type TAP-YYYY-NNNN. Aucun lien avec tickets.
// Table: project_tasks
type ProjectTask struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	ProjectID       uint       `gorm:"not null;index;uniqueIndex:idx_project_tasks_project_code,priority:1" json:"project_id"`
	ProjectPhaseID  uint       `gorm:"not null;index" json:"project_phase_id"`
	Code            string     `gorm:"type:varchar(50);not null;uniqueIndex:idx_project_tasks_project_code,priority:2" json:"code"` // ex. TAP-2025-0001, unique par (project_id, code)
	Title           string     `gorm:"type:varchar(255);not null" json:"title"`
	Description     string     `gorm:"type:text" json:"description,omitempty"`
	Status          string     `gorm:"type:varchar(50);default:'ouvert';index" json:"status"`   // ouvert, en_cours, en_attente, cloture
	Priority        string     `gorm:"type:varchar(50);default:'medium';index" json:"priority"` // low, medium, high, critical
	AssignedToID    *uint      `gorm:"index" json:"assigned_to_id,omitempty"`
	CreatedByID     uint       `gorm:"not null;index" json:"created_by_id"`
	EstimatedTime   *int       `gorm:"type:int" json:"estimated_time,omitempty"` // minutes
	ActualTime      int        `gorm:"column:actual_time;default:0" json:"actual_time"` // minutes (calculé ou saisi)
	DueDate         *time.Time `gorm:"type:date" json:"due_date,omitempty"`
	DisplayOrder    int        `gorm:"default:0" json:"display_order"`
	ClosedAt        *time.Time `json:"closed_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`

	Project      *Project                `gorm:"foreignKey:ProjectID" json:"-"`
	ProjectPhase *ProjectPhase           `gorm:"foreignKey:ProjectPhaseID" json:"-"`
	AssignedTo   *User                   `gorm:"foreignKey:AssignedToID" json:"-"`
	CreatedBy    *User                   `gorm:"foreignKey:CreatedByID" json:"-"`
	Assignees    []ProjectTaskAssignee   `gorm:"foreignKey:ProjectTaskID" json:"assignees,omitempty"`
}

// TableName spécifie le nom de la table
func (ProjectTask) TableName() string {
	return "project_tasks"
}
