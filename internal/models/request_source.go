package models

import (
	"time"
)

// RequestSource représente une source de demande configurée
// Table: request_sources
type RequestSource struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	Code        string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"` // Code unique: 'mail', 'appel', 'direct'
	Description string    `gorm:"type:text" json:"description,omitempty"`
	IsEnabled   bool      `gorm:"default:true;index" json:"is_enabled"` // Si la source est activée
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName spécifie le nom de la table
func (RequestSource) TableName() string {
	return "request_sources"
}
