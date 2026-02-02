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
	Phone        string         `gorm:"type:varchar(20)" json:"phone,omitempty"` // Numéro de téléphone
	PasswordHash string         `gorm:"type:varchar(255);not null" json:"-"`     // Mot de passe hashé (non exposé dans JSON)
	FirstName    string         `gorm:"type:varchar(100)" json:"first_name,omitempty"`
	LastName     string         `gorm:"type:varchar(100)" json:"last_name,omitempty"`
	DepartmentID *uint          `gorm:"index" json:"department_id,omitempty"`      // ID du département (optionnel)
	FilialeID    *uint          `gorm:"index" json:"filiale_id,omitempty"`         // ID de la filiale (optionnel)
	Avatar       string         `gorm:"type:varchar(500)" json:"avatar,omitempty"` // Chemin vers la photo de profil
	RoleID       uint           `gorm:"not null;index" json:"role_id"`
	IsActive     bool           `gorm:"default:true;index" json:"is_active"`
	LastLogin    *time.Time     `json:"last_login,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"` // Soft delete

	// Relations
	Role       Role        `gorm:"foreignKey:RoleID" json:"role,omitempty"`             // Rôle de l'utilisateur
	Department *Department `gorm:"foreignKey:DepartmentID" json:"department,omitempty"` // Département de l'utilisateur (optionnel)
	Filiale    *Filiale    `gorm:"foreignKey:FilialeID" json:"filiale,omitempty"`       // Filiale de l'utilisateur (optionnel)
	// CreatedBy et UpdatedBy sont des auto-références
	// IMPORTANT: Ne pas utiliser gorm:"foreignKey" ici pour éviter que GORM crée des contraintes incorrectes
	// Les contraintes seront créées manuellement dans les migrations via FixUserForeignKeys()
	CreatedBy   *User `gorm:"-" json:"-"`     // Utilisateur créateur (auto-référence, chargé manuellement si nécessaire)
	CreatedByID *uint `gorm:"index" json:"-"` // Index seulement, contrainte créée manuellement
	UpdatedBy   *User `gorm:"-" json:"-"`     // Utilisateur modificateur (auto-référence, chargé manuellement si nécessaire)
	UpdatedByID *uint `gorm:"index" json:"-"` // Index seulement, contrainte créée manuellement

	// Relations HasMany (définies dans les autres modèles)
	// TicketsCreated []Ticket `gorm:"foreignKey:CreatedByID" json:"-"`
	// TicketsAssigned []Ticket `gorm:"foreignKey:AssignedToID" json:"-"`
}

// TableName spécifie le nom de la table
func (User) TableName() string {
	return "users"
}
