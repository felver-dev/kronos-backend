package models

import (
	"time"

	"gorm.io/datatypes"
)

// Notification représente une notification pour un utilisateur
// Table: notifications
type Notification struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"not null;index" json:"user_id"`
	Type      string         `gorm:"type:varchar(100);not null;index" json:"type"` // delay_alert, budget_alert, validation_pending, etc.
	Title     string         `gorm:"type:varchar(255);not null" json:"title"`      // Titre de la notification
	Message   string         `gorm:"type:text;not null" json:"message"`            // Message de la notification
	IsRead    bool           `gorm:"default:false;index" json:"is_read"`           // Si la notification a été lue
	ReadAt    *time.Time     `json:"read_at,omitempty"`                            // Date de lecture (optionnel)
	LinkURL   string         `gorm:"type:varchar(500)" json:"link_url,omitempty"`  // URL vers la ressource concernée (optionnel)
	Metadata  datatypes.JSON `gorm:"type:json" json:"metadata,omitempty"`          // Données supplémentaires en JSON (optionnel)
	CreatedAt time.Time      `gorm:"index" json:"created_at"`

	// Relations
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"` // Utilisateur destinataire
}

// TableName spécifie le nom de la table
func (Notification) TableName() string {
	return "notifications"
}
