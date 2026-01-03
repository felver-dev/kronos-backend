package models

import (
	"time"
)

// Permission représente une permission disponible dans le système
// Table: permissions
type Permission struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"name"`
	Code        string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"code"` // Code unique de la permission
	Description string    `gorm:"type:text" json:"description,omitempty"`
	Module      string    `gorm:"type:varchar(50)" json:"module,omitempty"` // Module auquel appartient la permission (tickets, users, etc.)
	CreatedAt   time.Time `json:"created_at"`

	// Relations HasMany (définies dans les autres modèles)
	// RolePermissions []RolePermission `gorm:"foreignKey:PermissionID" json:"-"`
}

// TableName spécifie le nom de la table
func (Permission) TableName() string {
	return "permissions"
}
