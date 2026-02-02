package models

import (
	"time"

	"gorm.io/gorm"
)

// TicketCategory représente une catégorie de ticket
// Table: ticket_categories
type TicketCategory struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Name         string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"name"` // Nom de la catégorie (ex: incident, demande)
	Slug         string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"slug"` // Slug unique (ex: incident, demande, changement)
	Description  string         `gorm:"type:text" json:"description,omitempty"`             // Description de la catégorie
	Icon         string         `gorm:"type:varchar(100)" json:"icon,omitempty"`            // Nom de l'icône (ex: AlertTriangle, FileText)
	Color        string         `gorm:"type:varchar(50)" json:"color,omitempty"`            // Couleur associée (ex: red, blue)
	IsActive     bool           `gorm:"default:true;index" json:"is_active"`                // Catégorie active ou non
	DisplayOrder int            `gorm:"default:0;index" json:"display_order"`               // Ordre d'affichage
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"` // Soft delete

	// Relations
	Tickets []Ticket `gorm:"foreignKey:CategoryID" json:"-"` // Tickets de cette catégorie
}

// TableName spécifie le nom de la table
func (TicketCategory) TableName() string {
	return "ticket_categories"
}
