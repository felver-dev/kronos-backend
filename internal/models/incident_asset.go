package models

import (
	"time"
)

// IncidentAsset représente l'association entre un incident et un actif IT (table de liaison)
// Table: incident_assets
type IncidentAsset struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	IncidentID uint      `gorm:"not null;uniqueIndex:idx_incident_asset" json:"incident_id"`
	AssetID    uint      `gorm:"not null;uniqueIndex:idx_incident_asset" json:"asset_id"`
	CreatedAt  time.Time `json:"created_at"`

	// Relations
	Incident Incident `gorm:"foreignKey:IncidentID;constraint:OnDelete:CASCADE" json:"incident,omitempty"`
	Asset    Asset    `gorm:"foreignKey:AssetID;constraint:OnDelete:CASCADE" json:"asset,omitempty"`
}

// TableName spécifie le nom de la table
func (IncidentAsset) TableName() string {
	return "incident_assets"
}
