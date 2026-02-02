package models

import (
	"time"

	"gorm.io/gorm"
)

// Filiale représente une filiale du groupe
// Table: filiales
type Filiale struct {
	ID       uint    `gorm:"primaryKey" json:"id"`
	Code     string  `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"` // Code unique de la filiale (ex: MCI, SEN, TOG)
	Name     string  `gorm:"type:varchar(255);not null" json:"name"`            // Nom de la filiale
	Country  string  `gorm:"type:varchar(100)" json:"country,omitempty"`        // Pays de la filiale
	City     string  `gorm:"type:varchar(100)" json:"city,omitempty"`           // Ville
	Address  *string `gorm:"type:text" json:"address,omitempty"`                // Adresse complète
	Phone    string  `gorm:"type:varchar(20)" json:"phone,omitempty"`           // Téléphone
	Email    string  `gorm:"type:varchar(255)" json:"email,omitempty"`          // Email de contact
	IsActive bool    `gorm:"default:true;index" json:"is_active"`               // Si la filiale est active
	// IsSoftwareProvider : filiale fournisseur de logiciels/IT. Lu depuis la colonne is_mci_care_ci en BDD (rétrocompatibilité).
	IsSoftwareProvider bool           `gorm:"column:is_mci_care_ci;default:false;index" json:"is_software_provider"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"` // Soft delete

	// Relations HasMany
	Users       []User            `gorm:"foreignKey:FilialeID" json:"users,omitempty"`
	Departments []Department      `gorm:"foreignKey:FilialeID" json:"departments,omitempty"`
	Tickets     []Ticket          `gorm:"foreignKey:FilialeID" json:"tickets,omitempty"`
	Projects    []Project         `gorm:"foreignKey:FilialeID" json:"projects,omitempty"`
	Assets      []Asset           `gorm:"foreignKey:FilialeID" json:"assets,omitempty"`
	Offices     []Office          `gorm:"foreignKey:FilialeID" json:"offices,omitempty"`
	Deployments []FilialeSoftware `gorm:"foreignKey:FilialeID" json:"deployments,omitempty"`
}

// TableName spécifie le nom de la table
func (Filiale) TableName() string {
	return "filiales"
}
