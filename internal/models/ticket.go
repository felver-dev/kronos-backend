package models

import (
	"time"

	"gorm.io/gorm"
)

// Ticket représente un ticket dans le système
// Table: tickets
type Ticket struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	Code           string         `gorm:"type:varchar(50);uniqueIndex" json:"code"` // Code unique: TKT-YYYY-NNNN (nullable pour migration)
	Title          string         `gorm:"type:varchar(255);not null" json:"title"`
	Description    string         `gorm:"type:text" json:"description"`
	Category       string         `gorm:"type:varchar(50);not null;index" json:"category"`                // incident, demande, changement, developpement (slug pour compatibilité)
	CategoryID     *uint          `gorm:"index" json:"category_id,omitempty"`                              // ID de la catégorie (relation optionnelle)
	Source         string         `gorm:"type:varchar(50);not null" json:"source"`                        // mail, appel, direct
	Status         string         `gorm:"type:varchar(50);not null;default:'ouvert';index" json:"status"` // ouvert, en_cours, en_attente, cloture
	Priority       string         `gorm:"type:varchar(50);default:'medium'" json:"priority"`              // low, medium, high, critical
	AssignedToID       *uint          `gorm:"index" json:"assigned_to_id,omitempty"`                          // ID utilisateur assigné (optionnel)
	CreatedByID        uint           `gorm:"not null;index" json:"created_by_id"`
	RequesterID        *uint          `gorm:"index" json:"requester_id,omitempty"`                            // ID du demandeur (relation vers users)
	RequesterName      string         `gorm:"type:varchar(255)" json:"requester_name,omitempty"`              // Nom de la personne qui a fait la demande (fallback pour demandeurs externes)
	RequesterDepartment string        `gorm:"type:varchar(100)" json:"requester_department,omitempty"`         // Département du demandeur (ex: DAF)
	FilialeID           *uint         `gorm:"index" json:"filiale_id,omitempty"`                              // ID de la filiale (optionnel)
	SoftwareID          *uint         `gorm:"index" json:"software_id,omitempty"`                             // ID du logiciel concerné (optionnel)
	ValidatedByUserID   *uint         `gorm:"index" json:"validated_by_user_id,omitempty"`                     // ID de l'utilisateur qui a validé (optionnel)
	ValidatedAt         *time.Time    `json:"validated_at,omitempty"`                                          // Date de validation (optionnel)
	PrimaryImageID     *uint          `gorm:"index" json:"primary_image_id,omitempty"`                        // ID de l'image principale (optionnel)
	EstimatedTime  *int           `gorm:"type:int" json:"estimated_time,omitempty"` // Temps estimé en minutes (optionnel)
	ActualTime     *int           `gorm:"type:int" json:"actual_time,omitempty"`    // Temps réel en minutes (calculé)
	ParentID       *uint          `gorm:"index" json:"parent_id,omitempty"`          // Ticket parent (sous-ticket)
	CreatedAt      time.Time      `gorm:"index" json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	ClosedAt       *time.Time     `json:"closed_at,omitempty"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"` // Soft delete

	// Relations
	AssignedTo   *User             `gorm:"foreignKey:AssignedToID;references:ID" json:"assigned_to,omitempty"` // Utilisateur assigné
	CreatedBy    User              `gorm:"foreignKey:CreatedByID;references:ID" json:"created_by"`             // Créateur du ticket
	Requester    *User             `gorm:"foreignKey:RequesterID;references:ID" json:"requester,omitempty"`    // Demandeur (relation vers users)
	ValidatedBy  *User             `gorm:"foreignKey:ValidatedByUserID;references:ID" json:"validated_by,omitempty"` // Utilisateur qui a validé
	Filiale      *Filiale          `gorm:"foreignKey:FilialeID" json:"filiale,omitempty"`                      // Filiale (relation optionnelle)
	Software     *Software         `gorm:"foreignKey:SoftwareID" json:"software,omitempty"`                     // Logiciel concerné (relation optionnelle)
	CategoryObj  *TicketCategory   `gorm:"foreignKey:CategoryID" json:"category_obj,omitempty"`     // Catégorie (relation optionnelle)
	PrimaryImage *TicketAttachment `gorm:"foreignKey:PrimaryImageID" json:"primary_image,omitempty"` // Image principale (optionnel)
	Parent       *Ticket           `gorm:"foreignKey:ParentID" json:"parent,omitempty"`              // Ticket parent (optionnel)

	// Relations HasMany
	Comments    []TicketComment    `gorm:"foreignKey:TicketID" json:"comments,omitempty"`
	History     []TicketHistory    `gorm:"foreignKey:TicketID" json:"history,omitempty"`
	Attachments []TicketAttachment `gorm:"foreignKey:TicketID" json:"attachments,omitempty"`
	Assignees   []TicketAssignee   `gorm:"foreignKey:TicketID" json:"assignees,omitempty"`
	Solutions   []TicketSolution   `gorm:"foreignKey:TicketID" json:"solutions,omitempty"`
	SubTickets  []Ticket           `gorm:"foreignKey:ParentID" json:"sub_tickets,omitempty"`
	// TimeEntries []TimeEntry `gorm:"foreignKey:TicketID" json:"-"`
}

// TableName spécifie le nom de la table
func (Ticket) TableName() string {
	return "tickets"
}
