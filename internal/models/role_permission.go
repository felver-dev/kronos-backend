package models

import (
	"time"
)

// RolePermission représente l'association entre un rôle et une permission (table de liaison)
// Table: role_permissions
type RolePermission struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	RoleID       uint      `gorm:"not null;uniqueIndex:idx_role_permission" json:"role_id"`
	PermissionID uint      `gorm:"not null;uniqueIndex:idx_role_permission" json:"permission_id"`
	CreatedAt    time.Time `json:"created_at"`

	// Relations
	Role       Role       `gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE" json:"role,omitempty"`
	Permission Permission `gorm:"foreignKey:PermissionID;constraint:OnDelete:CASCADE" json:"permission,omitempty"`
}

// TableName spécifie le nom de la table
func (RolePermission) TableName() string {
	return "role_permissions"
}
