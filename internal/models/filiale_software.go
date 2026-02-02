package models

import (
	"time"

	"gorm.io/gorm"
)

// FilialeSoftware représente un déploiement d'un logiciel chez une filiale
// Table: filiale_software
type FilialeSoftware struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	FilialeID   uint           `gorm:"not null;index" json:"filiale_id"`                       // ID de la filiale
	SoftwareID  uint           `gorm:"not null;index" json:"software_id"`                     // ID du logiciel
	Version     string         `gorm:"type:varchar(50)" json:"version,omitempty"`               // Version déployée chez cette filiale
	DeployedAt  *time.Time     `json:"deployed_at,omitempty"`                                  // Date de déploiement
	IsActive    bool           `gorm:"default:true;index" json:"is_active"`                     // Si le déploiement est actif
	Notes       *string        `gorm:"type:text" json:"notes,omitempty"`                        // Notes sur le déploiement
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"` // Soft delete

	// Relations
	Filiale  Filiale  `gorm:"foreignKey:FilialeID" json:"filiale,omitempty"`
	Software Software `gorm:"foreignKey:SoftwareID" json:"software,omitempty"`
}

// TableName spécifie le nom de la table
func (FilialeSoftware) TableName() string {
	return "filiale_software"
}
