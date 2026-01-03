package models

import (
	"time"

	"gorm.io/gorm"
)

// Ticket représente un ticket dans le système
// Table: tickets
type Ticket struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	Title          string         `gorm:"type:varchar(255);not null" json:"title"`
	Description    string         `gorm:"type:text" json:"description"`
	Category       string         `gorm:"type:varchar(50);not null;index" json:"category"`                // incident, demande, changement, developpement
	Source         string         `gorm:"type:varchar(50);not null" json:"source"`                        // mail, appel, direct
	Status         string         `gorm:"type:varchar(50);not null;default:'ouvert';index" json:"status"` // ouvert, en_cours, en_attente, cloture
	Priority       string         `gorm:"type:varchar(50);default:'medium'" json:"priority"`              // low, medium, high, critical
	AssignedToID   *uint          `gorm:"index" json:"assigned_to_id,omitempty"`                          // ID utilisateur assigné (optionnel)
	CreatedByID    uint           `gorm:"not null;index" json:"created_by_id"`
	PrimaryImageID *uint          `json:"primary_image_id,omitempty"`               // ID de l'image principale (optionnel)
	EstimatedTime  *int           `gorm:"type:int" json:"estimated_time,omitempty"` // Temps estimé en minutes (optionnel)
	ActualTime     *int           `gorm:"type:int" json:"actual_time,omitempty"`    // Temps réel en minutes (calculé)
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	ClosedAt       *time.Time     `json:"closed_at,omitempty"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"` // Soft delete

	// Relations
	AssignedTo *User `gorm:"foreignKey:AssignedToID" json:"assigned_to,omitempty"` // Utilisateur assigné
	CreatedBy  User  `gorm:"foreignKey:CreatedByID" json:"created_by"`             // Créateur du ticket
	// PrimaryImage  *TicketAttachment `gorm:"foreignKey:PrimaryImageID" json:"-"`        // Image principale (à créer plus tard)

	// Relations HasMany (définies dans les autres modèles)
	// Comments    []TicketComment `gorm:"foreignKey:TicketID" json:"-"`
	// History     []TicketHistory `gorm:"foreignKey:TicketID" json:"-"`
	// Attachments []TicketAttachment `gorm:"foreignKey:TicketID" json:"-"`
	// TimeEntries []TimeEntry `gorm:"foreignKey:TicketID" json:"-"`
}

// TableName spécifie le nom de la table
func (Ticket) TableName() string {
	return "tickets"
}
