package models

import (
	"time"
)

// Setting représente un paramètre système
// Table: settings
type Setting struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Key         string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"key"` // Clé unique du paramètre
	Value       string    `gorm:"type:text" json:"value"`                             // Valeur du paramètre (peut être JSON)
	Type        string    `gorm:"type:varchar(50);default:'string'" json:"type"`      // string, number, boolean, json
	Category    string    `gorm:"type:varchar(100);index" json:"category,omitempty"`   // Catégorie du paramètre (optionnel)
	Description string    `gorm:"type:text" json:"description,omitempty"`             // Description du paramètre (optionnel)
	IsPublic    bool      `gorm:"default:false;index" json:"is_public"`                // Si le paramètre est accessible publiquement (sans auth)
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	UpdatedByID *uint      `gorm:"index" json:"-"`
	UpdatedBy   *User      `gorm:"foreignKey:UpdatedByID" json:"-"`
}

// TableName spécifie le nom de la table
func (Setting) TableName() string {
	return "settings"
}

