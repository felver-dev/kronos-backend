package models

import (
	"time"

	"gorm.io/gorm"
)

// Office représente un siège/bureau de l'entreprise
// Table: offices
type Office struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(255);not null;index" json:"name"`       // Nom du siège (obligatoire)
	Code      *string        `gorm:"type:varchar(50);uniqueIndex" json:"code,omitempty"` // Code du siège (optionnel pour compatibilité, mais requis via l'API)
	Country   string         `gorm:"type:varchar(100);not null;index" json:"country"`    // Pays (obligatoire)
	City      string         `gorm:"type:varchar(100);not null;index" json:"city"`       // Ville (obligatoire)
	Commune   *string        `gorm:"type:varchar(100)" json:"commune,omitempty"`         // Commune (optionnel)
	Address   *string        `gorm:"type:text" json:"address,omitempty"`                 // Adresse complète (optionnel)
	FilialeID *uint          `gorm:"index" json:"filiale_id,omitempty"`                  // ID de la filiale (optionnel)
	Longitude *float64       `gorm:"type:decimal(10,8)" json:"longitude,omitempty"`      // Longitude (optionnel)
	Latitude  *float64       `gorm:"type:decimal(10,8)" json:"latitude,omitempty"`       // Latitude (optionnel)
	IsActive  bool           `gorm:"default:true;index" json:"is_active"`                // Si le siège est actif
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // Soft delete

	// Relations
	Filiale *Filiale `gorm:"foreignKey:FilialeID" json:"filiale,omitempty"` // Filiale (optionnel)
}

// TableName spécifie le nom de la table
func (Office) TableName() string {
	return "offices"
}
