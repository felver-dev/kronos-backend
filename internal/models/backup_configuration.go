package models

import (
	"time"
)

// BackupConfiguration représente une configuration de sauvegarde
// Table: backup_configurations
type BackupConfiguration struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	Frequency     string     `gorm:"type:varchar(50);not null" json:"frequency"` // daily, weekly, monthly
	Time          time.Time  `gorm:"type:time;not null" json:"time"`               // Heure de la sauvegarde
	RetentionDays int        `gorm:"default:30" json:"retention_days"`            // Nombre de jours de rétention
	IsActive      bool       `gorm:"default:true;index" json:"is_active"`        // Si la configuration est active
	LastBackupAt  *time.Time `json:"last_backup_at,omitempty"`                    // Date de la dernière sauvegarde
	NextBackupAt  *time.Time `gorm:"index" json:"next_backup_at,omitempty"`      // Date de la prochaine sauvegarde
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	UpdatedByID   *uint      `gorm:"index" json:"-"`
	UpdatedBy     *User      `gorm:"foreignKey:UpdatedByID" json:"-"`

	// Relations HasMany (définies dans les autres modèles)
	// Backups []Backup `gorm:"foreignKey:ConfigurationID" json:"-"`
}

// TableName spécifie le nom de la table
func (BackupConfiguration) TableName() string {
	return "backup_configurations"
}

