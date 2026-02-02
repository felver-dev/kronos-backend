package models

import (
	"time"

	"gorm.io/gorm"
)

// TicketInternal représente un ticket interne (départements non-IT, par département)
// Table: ticket_internes
type TicketInternal struct {
	ID                 uint           `gorm:"primaryKey" json:"id"`
	Code               string         `gorm:"type:varchar(50);uniqueIndex" json:"code"` // Ex: TKI-YYYY-NNNN
	Title              string         `gorm:"type:varchar(255);not null" json:"title"`
	Description        string         `gorm:"type:text" json:"description"`
	Category           string         `gorm:"type:varchar(50);not null;index" json:"category"`                // slug catégorie interne
	Status             string         `gorm:"type:varchar(50);not null;default:'ouvert';index" json:"status"` // ouvert, en_cours, en_attente, resolu, cloture
	Priority           string         `gorm:"type:varchar(50);default:'medium'" json:"priority"`             // low, medium, high, critical
	DepartmentID       uint           `gorm:"not null;index" json:"department_id"`                             // Département propriétaire
	FilialeID          uint           `gorm:"not null;index" json:"filiale_id"`                                // Filiale (déduit du département)
	CreatedByID        uint           `gorm:"not null;index" json:"created_by_id"`
	AssignedToID       *uint          `gorm:"index" json:"assigned_to_id,omitempty"`
	ValidatedByUserID  *uint          `gorm:"index" json:"validated_by_user_id,omitempty"`
	ValidatedAt        *time.Time     `json:"validated_at,omitempty"`
	EstimatedTime      *int           `gorm:"type:int" json:"estimated_time,omitempty"` // minutes
	ActualTime         *int           `gorm:"type:int" json:"actual_time,omitempty"`    // minutes
	TicketID           *uint          `gorm:"index" json:"ticket_id,omitempty"`        // Lien optionnel vers un ticket normal
	CreatedAt          time.Time      `gorm:"index" json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	ClosedAt           *time.Time     `json:"closed_at,omitempty"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Department  Department `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
	Filiale     Filiale    `gorm:"foreignKey:FilialeID" json:"filiale,omitempty"`
	CreatedBy   User       `gorm:"foreignKey:CreatedByID" json:"created_by,omitempty"`
	AssignedTo  *User      `gorm:"foreignKey:AssignedToID" json:"assigned_to,omitempty"`
	ValidatedBy *User      `gorm:"foreignKey:ValidatedByUserID" json:"validated_by,omitempty"`
	Ticket      *Ticket    `gorm:"foreignKey:TicketID" json:"ticket,omitempty"`
}

// TableName spécifie le nom de la table
func (TicketInternal) TableName() string {
	return "ticket_internes"
}
