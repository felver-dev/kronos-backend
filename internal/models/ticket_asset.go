package models

import (
	"time"
)

// TicketAsset représente l'association entre un ticket et un actif IT (table de liaison)
// Table: ticket_assets
type TicketAsset struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	TicketID  uint      `gorm:"not null;uniqueIndex:idx_ticket_asset" json:"ticket_id"`
	AssetID   uint      `gorm:"not null;uniqueIndex:idx_ticket_asset" json:"asset_id"`
	CreatedAt time.Time `json:"created_at"`

	// Relations
	Ticket Ticket `gorm:"foreignKey:TicketID;constraint:OnDelete:CASCADE" json:"ticket,omitempty"`
	Asset  Asset  `gorm:"foreignKey:AssetID;constraint:OnDelete:CASCADE" json:"asset,omitempty"`
}

// TableName spécifie le nom de la table
func (TicketAsset) TableName() string {
	return "ticket_assets"
}
