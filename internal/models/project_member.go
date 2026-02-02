package models

import (
	"time"
)

// ProjectMember représente un membre du projet (affectation au projet global).
// Les fonctions (direction + exécution) sont dans project_member_functions (Functions).
// Table: project_members
type ProjectMember struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	ProjectID        uint      `gorm:"not null;uniqueIndex:idx_project_member" json:"project_id"`
	UserID           uint      `gorm:"not null;uniqueIndex:idx_project_member" json:"user_id"`
	ProjectFunctionID *uint    `gorm:"index" json:"project_function_id,omitempty"` // conservé pour rétrocompat, préférer Functions
	IsProjectManager bool      `gorm:"default:false" json:"is_project_manager"`
	IsLead           bool      `gorm:"default:false" json:"is_lead"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`

	Project  *Project          `gorm:"foreignKey:ProjectID" json:"-"`
	User     *User             `gorm:"foreignKey:UserID" json:"-"`
	Function *ProjectFunction  `gorm:"foreignKey:ProjectFunctionID" json:"-"`
	Functions []ProjectFunction `gorm:"many2many:project_member_functions;" json:"functions,omitempty"`
}

// TableName spécifie le nom de la table
func (ProjectMember) TableName() string {
	return "project_members"
}
