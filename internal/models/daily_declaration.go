package models

import "time"

// DailyDeclaration représente une déclaration journalière des tâches
// Table: daily_declarations
type DailyDeclaration struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	UserID            uint           `gorm:"not null;index" json:"user_id"`
	Date              time.Time      `gorm:"type:date;not null;index" json:"date"` // Date de la déclaration
	TaskCount         int            `gorm:"default:0" json:"task_count"`            // Nombre de tâches
	TotalTime         int            `gorm:"default:0" json:"total_time"`            // Temps total en minutes
	Validated         bool           `gorm:"default:false;index" json:"validated"`   // Si la déclaration a été validée
	ValidatedByID     *uint          `gorm:"index" json:"validated_by_id,omitempty"`  // ID du validateur (optionnel)
	ValidatedAt       *time.Time     `json:"validated_at,omitempty"`                 // Date de validation (optionnel)
	ValidationComment string         `gorm:"type:text" json:"validation_comment,omitempty"` // Commentaire de validation (optionnel)
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`

	// Relations
	User        User                  `gorm:"foreignKey:UserID" json:"user,omitempty"` // Utilisateur
	ValidatedBy *User                 `gorm:"foreignKey:ValidatedByID" json:"validated_by,omitempty"` // Validateur (optionnel)
	Tasks       []DailyDeclarationTask `gorm:"foreignKey:DeclarationID;constraint:OnDelete:CASCADE" json:"tasks,omitempty"` // Tâches déclarées
}

// TableName spécifie le nom de la table
func (DailyDeclaration) TableName() string {
	return "daily_declarations"
}

// DailyDeclarationTask représente une tâche déclarée dans une déclaration journalière
// Table: daily_declaration_tasks
type DailyDeclarationTask struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	DeclarationID uint      `gorm:"not null;index" json:"declaration_id"`
	TicketID      uint      `gorm:"not null;index" json:"ticket_id"`
	TimeSpent     int       `gorm:"not null" json:"time_spent"` // Temps passé en minutes
	CreatedAt     time.Time `json:"created_at"`

	// Relations - GORM utilisera automatiquement les champs existants
	Declaration DailyDeclaration `gorm:"foreignKey:DeclarationID;constraint:OnDelete:CASCADE" json:"-"` // Déclaration associée
	Ticket      Ticket           `gorm:"foreignKey:TicketID;constraint:OnDelete:CASCADE" json:"ticket,omitempty"` // Ticket associé
}

// TableName spécifie le nom de la table
func (DailyDeclarationTask) TableName() string {
	return "daily_declaration_tasks"
}

