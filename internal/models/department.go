package models

import (
	"time"

	"gorm.io/gorm"
)

// Department représente un département de l'entreprise
// Table: departments
type Department struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"type:varchar(255);not null;index" json:"name"`        // Nom du département (obligatoire)
	Code        string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"`    // Code unique du département (obligatoire)
	Description *string        `gorm:"type:text" json:"description,omitempty"`              // Description (optionnel)
	FilialeID   *uint          `gorm:"index" json:"filiale_id,omitempty"`                   // ID de la filiale (optionnel)
	OfficeID    *uint          `gorm:"index" json:"office_id,omitempty"`                      // ID du siège (optionnel)
	IsITDepartment bool        `gorm:"default:false;index" json:"is_it_department"`           // Si c'est un département IT (uniquement pour la filiale fournisseur de logiciels)
	Office      *Office        `gorm:"foreignKey:OfficeID" json:"office,omitempty"`          // Relation vers le siège
	Filiale     *Filiale       `gorm:"foreignKey:FilialeID" json:"filiale,omitempty"`        // Relation vers la filiale
	IsActive    bool           `gorm:"default:true;index" json:"is_active"`                   // Si le département est actif
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"` // Soft delete
}

// TableName spécifie le nom de la table
func (Department) TableName() string {
	return "departments"
}
