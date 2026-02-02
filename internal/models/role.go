package models

import (
	"time"

	"gorm.io/gorm"
)

// Role représente un rôle utilisateur dans le système
// Table: roles
type Role struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"name"` // DSI, RESPONSABLE_IT, TECHNICIEN_IT
	Description string         `gorm:"type:text" json:"description,omitempty"`
	IsSystem    bool           `gorm:"default:false" json:"is_system"`       // Rôle système (ne peut pas être supprimé)
	CreatedByID *uint          `gorm:"index" json:"created_by_id,omitempty"` // ID de l'utilisateur qui a créé le rôle (nil pour les rôles système)
	FilialeID   *uint          `gorm:"index" json:"filiale_id,omitempty"`    // ID de la filiale à laquelle le rôle appartient (nil pour les rôles globaux)
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"` // Soft delete

	// Relations
	Users     []User   `gorm:"foreignKey:RoleID" json:"-"`      // Utilisateurs ayant ce rôle
	CreatedBy *User    `gorm:"foreignKey:CreatedByID" json:"-"` // Utilisateur créateur (optionnel)
	Filiale   *Filiale `gorm:"foreignKey:FilialeID" json:"-"`   // Filiale à laquelle appartient le rôle (optionnel)
}

// TableName spécifie le nom de la table (optionnel, GORM utilise le nom du modèle par défaut)
func (Role) TableName() string {
	return "roles"
}
