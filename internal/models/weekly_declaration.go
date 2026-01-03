package models

import "time"

// WeeklyDeclaration représente une déclaration hebdomadaire des tâches
// Table: weekly_declarations
type WeeklyDeclaration struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	UserID            uint           `gorm:"not null;index" json:"user_id"`
	Week              string         `gorm:"type:varchar(10);not null;index" json:"week"` // Format ISO: "2024-W03"
	StartDate         time.Time      `gorm:"type:date;not null" json:"start_date"`         // Date de début de la semaine
	EndDate           time.Time      `gorm:"type:date;not null" json:"end_date"`           // Date de fin de la semaine
	TaskCount         int            `gorm:"default:0" json:"task_count"`                  // Nombre total de tâches
	TotalTime         int            `gorm:"default:0" json:"total_time"`                  // Temps total en minutes
	Validated         bool           `gorm:"default:false;index" json:"validated"`         // Si la déclaration a été validée
	ValidatedByID     *uint          `gorm:"index" json:"validated_by_id,omitempty"`       // ID du validateur (optionnel)
	ValidatedAt       *time.Time     `json:"validated_at,omitempty"`                       // Date de validation (optionnel)
	ValidationComment string         `gorm:"type:text" json:"validation_comment,omitempty"` // Commentaire de validation (optionnel)
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`

	// Relations
	User        User                  `gorm:"foreignKey:UserID" json:"user,omitempty"` // Utilisateur
	ValidatedBy *User                 `gorm:"foreignKey:ValidatedByID" json:"validated_by,omitempty"` // Validateur (optionnel)
	Tasks       []WeeklyDeclarationTask `gorm:"foreignKey:DeclarationID;constraint:OnDelete:CASCADE" json:"tasks,omitempty"` // Tâches déclarées
}

// TableName spécifie le nom de la table
func (WeeklyDeclaration) TableName() string {
	return "weekly_declarations"
}

// WeeklyDeclarationTask représente une tâche déclarée dans une déclaration hebdomadaire
// Table: weekly_declaration_tasks
type WeeklyDeclarationTask struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	DeclarationID uint      `gorm:"not null;index" json:"declaration_id"`
	TicketID      uint      `gorm:"not null;index" json:"ticket_id"`
	Date          time.Time `gorm:"type:date;not null;index" json:"date"` // Date de la tâche
	TimeSpent     int       `gorm:"not null" json:"time_spent"`           // Temps passé en minutes
	CreatedAt     time.Time `json:"created_at"`

	// Relations - GORM utilisera automatiquement les champs existants
	Declaration WeeklyDeclaration `gorm:"foreignKey:DeclarationID;constraint:OnDelete:CASCADE" json:"-"` // Déclaration associée
	Ticket      Ticket           `gorm:"constraint:OnDelete:CASCADE" json:"ticket,omitempty"` // Ticket associé
}

// TableName spécifie le nom de la table
func (WeeklyDeclarationTask) TableName() string {
	return "weekly_declaration_tasks"
}

