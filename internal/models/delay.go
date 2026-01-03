package models

import "time"

// Delay représente un retard détecté sur un ticket
// Table: delays
type Delay struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	TicketID        uint       `gorm:"uniqueIndex;not null;index" json:"ticket_id"` // Relation 1:1 avec Ticket
	UserID          uint       `gorm:"not null;index" json:"user_id"`                // Technicien en retard
	EstimatedTime   int        `gorm:"not null" json:"estimated_time"`              // Temps estimé en minutes
	ActualTime      int        `gorm:"not null" json:"actual_time"`                 // Temps réel en minutes
	DelayTime       int        `gorm:"not null" json:"delay_time"`                  // Retard en minutes (actual - estimated)
	DelayPercentage float64    `gorm:"type:decimal(5,2);not null" json:"delay_percentage"` // Pourcentage de retard
	Status          string     `gorm:"type:varchar(50);default:'unjustified';index" json:"status"` // unjustified, pending, justified, rejected
	DetectedAt      time.Time  `gorm:"index" json:"detected_at"`                    // Date de détection
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`

	// Relations - GORM utilisera automatiquement les champs existants
	Ticket       Ticket            `gorm:"constraint:OnDelete:CASCADE" json:"ticket,omitempty"` // Ticket associé (1:1)
	User         User              `gorm:"foreignKey:UserID" json:"user,omitempty"`                                    // Technicien
	Justification *DelayJustification `gorm:"foreignKey:DelayID;constraint:OnDelete:CASCADE" json:"justification,omitempty"` // Justification (1:1, optionnel)
}

// TableName spécifie le nom de la table
func (Delay) TableName() string {
	return "delays"
}

// DelayJustification représente une justification de retard
// Table: delay_justifications
type DelayJustification struct {
	ID                uint       `gorm:"primaryKey" json:"id"`
	DelayID           uint       `gorm:"uniqueIndex;not null;index" json:"delay_id"` // Relation 1:1 avec Delay
	UserID            uint       `gorm:"not null;index" json:"user_id"`              // Technicien qui justifie
	Justification     string     `gorm:"type:text;not null" json:"justification"`     // Texte de justification
	Status            string     `gorm:"type:varchar(50);default:'pending';index" json:"status"` // pending, validated, rejected
	ValidatedByID     *uint      `gorm:"index" json:"validated_by_id,omitempty"`     // ID du validateur (optionnel)
	ValidatedAt       *time.Time `json:"validated_at,omitempty"`                     // Date de validation (optionnel)
	ValidationComment string     `gorm:"type:text" json:"validation_comment,omitempty"` // Commentaire du validateur (optionnel)
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`

	// Relations
	Delay       Delay `gorm:"foreignKey:DelayID;constraint:OnDelete:CASCADE" json:"delay,omitempty"` // Retard associé (1:1)
	User        User  `gorm:"foreignKey:UserID" json:"user,omitempty"`                               // Technicien qui justifie
	ValidatedBy *User `gorm:"foreignKey:ValidatedByID" json:"validated_by,omitempty"`                 // Validateur (optionnel)
}

// TableName spécifie le nom de la table
func (DelayJustification) TableName() string {
	return "delay_justifications"
}

