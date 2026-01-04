package dto

import "time"

// UpdateSettingsRequest représente la requête de mise à jour des paramètres
type UpdateSettingsRequest struct {
	Settings map[string]interface{} `json:"settings"` // Map de clé-valeur des paramètres à mettre à jour
}

// RequestSourceDTO représente une source de demande
type RequestSourceDTO struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Description string    `json:"description,omitempty"`
	IsEnabled   bool      `json:"is_enabled"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateRequestSourceRequest représente la requête de création d'une source
type CreateRequestSourceRequest struct {
	Name        string `json:"name" binding:"required"`
	Code        string `json:"code" binding:"required"`
	Description string `json:"description,omitempty"`
	IsEnabled   bool   `json:"is_enabled,omitempty"`
}

// UpdateRequestSourceRequest représente la requête de mise à jour d'une source
type UpdateRequestSourceRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	IsEnabled   *bool  `json:"is_enabled,omitempty"`
}

// BackupConfigurationDTO représente la configuration de sauvegarde
type BackupConfigurationDTO struct {
	Frequency      string    `json:"frequency"`       // daily, weekly, monthly
	Time           string    `json:"time"`            // Heure de sauvegarde (format: HH:MM)
	Retention      int       `json:"retention"`       // Nombre de jours de rétention
	AutoBackup     bool      `json:"auto_backup"`     // Si la sauvegarde automatique est activée
	LastBackup     *time.Time `json:"last_backup,omitempty"`
	NextBackup     *time.Time `json:"next_backup,omitempty"`
}

// ExecuteBackupRequest représente la requête d'exécution d'une sauvegarde
type ExecuteBackupRequest struct {
	Type string `json:"type,omitempty"` // full, incremental
}

// BackupExecutionResponse représente la réponse d'exécution d'une sauvegarde
type BackupExecutionResponse struct {
	BackupID      uint   `json:"backup_id"`
	Status        string `json:"status"`         // in_progress, completed, failed
	EstimatedTime int    `json:"estimated_time"` // en secondes
	Message       string `json:"message,omitempty"`
}

