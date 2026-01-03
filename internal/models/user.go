package models

import (
	"time"

	"gorm.io/gorm"
)

// User représente un utilisateur du système
// Table: users
type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Username     string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"username"`
	Email        string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	PasswordHash string         `gorm:"type:varchar(255);not null" json:"-"` // Mot de passe hashé (non exposé dans JSON)
	FirstName    string         `gorm:"type:varchar(100)" json:"first_name,omitempty"`
	LastName     string         `gorm:"type:varchar(100)" json:"last_name,omitempty"`
	Avatar       string         `gorm:"type:varchar(500)" json:"avatar,omitempty"` // Chemin vers la photo de profil
	RoleID       uint           `gorm:"not null;index" json:"role_id"`
	IsActive     bool           `gorm:"default:true;index" json:"is_active"`
	LastLogin    *time.Time     `json:"last_login,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"` // Soft delete

	// Relations
	Role        Role  `gorm:"foreignKey:RoleID" json:"role,omitempty"` // Rôle de l'utilisateur
	CreatedBy   *User `gorm:"foreignKey:CreatedByID" json:"-"`         // Utilisateur créateur (auto-référence)
	CreatedByID *uint `gorm:"index" json:"-"`
	UpdatedBy   *User `gorm:"foreignKey:UpdatedByID" json:"-"` // Utilisateur modificateur (auto-référence)
	UpdatedByID *uint `json:"-"`

	// Relations HasMany (définies dans les autres modèles)
	// TicketsCreated []Ticket `gorm:"foreignKey:CreatedByID" json:"-"`
	// TicketsAssigned []Ticket `gorm:"foreignKey:AssignedToID" json:"-"`
}

// TableName spécifie le nom de la table
func (User) TableName() string {
	return "users"
}
