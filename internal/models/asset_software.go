package models

import (
	"time"

	"gorm.io/gorm"
)

// AssetSoftware représente un logiciel installé sur un actif IT
// Table: asset_software
type AssetSoftware struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	AssetID         *uint          `gorm:"index" json:"asset_id,omitempty"`                      // ID de l'actif (optionnel - permet de créer des logiciels indépendamment)
	SoftwareName    string         `gorm:"type:varchar(255);not null" json:"software_name"`    // Nom du logiciel (ex: "Windows 11", "Office 365")
	Version         string         `gorm:"type:varchar(100)" json:"version,omitempty"`          // Version (ex: "11.0", "2021")
	LicenseKey      string         `gorm:"type:varchar(255)" json:"license_key,omitempty"`      // Clé de licence (optionnel)
	InstallationDate *time.Time    `gorm:"type:date" json:"installation_date,omitempty"`        // Date d'installation
	Notes           string         `gorm:"type:text" json:"notes,omitempty"`                    // Notes optionnelles
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"` // Soft delete

	// Relations
	Asset *Asset `gorm:"foreignKey:AssetID;constraint:OnDelete:SET NULL" json:"asset,omitempty"`
}

// TableName spécifie le nom de la table
func (AssetSoftware) TableName() string {
	return "asset_software"
}
