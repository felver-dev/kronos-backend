package models

import (
	"time"

	"gorm.io/gorm"
)

// Software représente un logiciel géré par la filiale fournisseur IT
// Table: software
type Software struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Code        string         `gorm:"type:varchar(50);uniqueIndex:idx_software_code_version;not null" json:"code"`        // Code du logiciel (ex: ISA) — unique avec Version
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`                  // Nom du logiciel
	Description *string        `gorm:"type:text" json:"description,omitempty"`                   // Description
	Version     string         `gorm:"type:varchar(50);uniqueIndex:idx_software_code_version" json:"version,omitempty"`    // Version (ex: 33, 35) — unique avec Code
	IsActive    bool           `gorm:"default:true;index" json:"is_active"`                     // Si le logiciel est actif
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"` // Soft delete

	// Relations HasMany
	Tickets     []Ticket         `gorm:"foreignKey:SoftwareID" json:"tickets,omitempty"`
	Deployments []FilialeSoftware `gorm:"foreignKey:SoftwareID" json:"deployments,omitempty"`
}

// TableName spécifie le nom de la table
func (Software) TableName() string {
	return "software"
}
