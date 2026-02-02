package models

import "time"

// Delay représente un retard détecté sur un ticket ou un ticket interne
// Table: delays. Soit ticket_id soit ticket_internal_id (l'un des deux).
type Delay struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	TicketID         *uint     `gorm:"index" json:"ticket_id,omitempty"`         // Ticket normal (optionnel si ticket_internal_id) — index unique créé en migration
	TicketInternalID *uint     `gorm:"index" json:"ticket_internal_id,omitempty"` // Ticket interne (optionnel si ticket_id) — index unique créé en migration
	FilialeID        *uint     `gorm:"index" json:"filiale_id,omitempty"`                                            // ID de la filiale (via ticket, pour faciliter les filtres)
	UserID          uint      `gorm:"not null;index" json:"user_id"`                              // Technicien en retard
	EstimatedTime   int       `gorm:"not null" json:"estimated_time"`                             // Temps estimé en minutes
	ActualTime      int       `gorm:"not null" json:"actual_time"`                                // Temps réel en minutes
	DelayTime       int       `gorm:"not null" json:"delay_time"`                                 // Retard en minutes (actual - estimated)
	DelayPercentage float64   `gorm:"type:decimal(5,2);not null" json:"delay_percentage"`         // Pourcentage de retard
	Status          string    `gorm:"type:varchar(50);default:'unjustified';index" json:"status"` // unjustified, pending, justified, rejected
	DetectedAt      time.Time `gorm:"index" json:"detected_at"`                                   // Date de détection
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	// Relations
	Ticket         *Ticket         `gorm:"foreignKey:TicketID;constraint:OnDelete:CASCADE" json:"ticket,omitempty"`
	TicketInternal *TicketInternal `gorm:"foreignKey:TicketInternalID;constraint:OnDelete:CASCADE" json:"ticket_internal,omitempty"`
	Filiale        *Filiale        `gorm:"foreignKey:FilialeID" json:"filiale,omitempty"`
	User           User            `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Justification  *DelayJustification `gorm:"foreignKey:DelayID;constraint:OnDelete:CASCADE" json:"justification,omitempty"`
}

// TableName spécifie le nom de la table
func (Delay) TableName() string {
	return "delays"
}

// DelayJustification représente une justification de retard
// Table: delay_justifications
type DelayJustification struct {
	ID                uint       `gorm:"primaryKey" json:"id"`
	DelayID           uint       `gorm:"uniqueIndex:idx_delay_justifications_delay_id;not null" json:"delay_id"`             // Relation 1:1 avec Delay
	UserID            uint       `gorm:"not null;index" json:"user_id"`                          // Technicien qui justifie
	Justification     string     `gorm:"type:text;not null" json:"justification"`                // Texte de justification
	Status            string     `gorm:"type:varchar(50);default:'pending';index" json:"status"` // pending, validated, rejected
	ValidatedByID     *uint      `gorm:"index" json:"validated_by_id,omitempty"`                 // ID du validateur (optionnel)
	ValidatedAt       *time.Time `json:"validated_at,omitempty"`                                 // Date de validation (optionnel)
	ValidationComment string     `gorm:"type:text" json:"validation_comment,omitempty"`          // Commentaire du validateur (optionnel)
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`

	// Relations
	Delay       Delay `gorm:"foreignKey:DelayID;constraint:OnDelete:CASCADE" json:"delay,omitempty"` // Retard associé (1:1)
	User        User  `gorm:"foreignKey:UserID" json:"user,omitempty"`                               // Technicien qui justifie
	ValidatedBy *User `gorm:"foreignKey:ValidatedByID" json:"validated_by,omitempty"`                // Validateur (optionnel)
}

// TableName spécifie le nom de la table
func (DelayJustification) TableName() string {
	return "delay_justifications"
}
