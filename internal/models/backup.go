package models

import (
	"time"
)

// Backup représente un historique de sauvegarde
// Table: backups
type Backup struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	ConfigurationID uint       `gorm:"not null;index" json:"configuration_id"`
	FilePath        string     `gorm:"type:varchar(500);not null" json:"file_path"` // Chemin vers le fichier de sauvegarde
	FileSize        *int64     `gorm:"type:bigint" json:"file_size,omitempty"`      // Taille en bytes (optionnel)
	Status          string     `gorm:"type:varchar(50);not null;index" json:"status"` // in_progress, completed, failed
	StartedAt       time.Time  `gorm:"index" json:"started_at"`
	CompletedAt     *time.Time `json:"completed_at,omitempty"` // Date de fin (optionnel)
	ErrorMessage    string     `gorm:"type:text" json:"error_message,omitempty"`     // Message d'erreur si échec (optionnel)
	CreatedByID     *uint       `gorm:"index" json:"-"`
	CreatedBy       *User       `gorm:"foreignKey:CreatedByID" json:"-"`

	// Relations
	Configuration BackupConfiguration `gorm:"foreignKey:ConfigurationID" json:"configuration,omitempty"`
}

// TableName spécifie le nom de la table
func (Backup) TableName() string {
	return "backups"
}

