package models

import (
	"time"

	"gorm.io/datatypes"
)

// AuditLog représente un log d'audit pour la traçabilité
// Table: audit_logs
type AuditLog struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      *uint          `gorm:"index" json:"user_id,omitempty"` // Utilisateur qui a effectué l'action (optionnel, peut être NULL pour actions système)
	Action      string         `gorm:"type:varchar(100);not null;index" json:"action"` // create, update, delete, login, logout, etc.
	EntityType  string         `gorm:"type:varchar(100);not null;index" json:"entity_type"` // Type d'entité (users, tickets, etc.)
	EntityID    *uint          `gorm:"index" json:"entity_id,omitempty"` // ID de l'entité concernée (optionnel)
	OldValues   datatypes.JSON `gorm:"type:json" json:"old_values,omitempty"` // Anciennes valeurs (optionnel)
	NewValues   datatypes.JSON `gorm:"type:json" json:"new_values,omitempty"` // Nouvelles valeurs (optionnel)
	IPAddress   string         `gorm:"type:varchar(45)" json:"ip_address,omitempty"` // Adresse IP (IPv4 ou IPv6)
	UserAgent   string         `gorm:"type:varchar(500)" json:"user_agent,omitempty"` // User-Agent du navigateur
	Description string         `gorm:"type:text" json:"description,omitempty"` // Description de l'action (optionnel)
	CreatedAt   time.Time      `gorm:"index" json:"created_at"`

	// Relations
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"` // Utilisateur (optionnel)
}

// TableName spécifie le nom de la table
func (AuditLog) TableName() string {
	return "audit_logs"
}

