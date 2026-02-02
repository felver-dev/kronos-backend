package dto

import "time"

// SoftwareDTO représente un logiciel dans les réponses API
type SoftwareDTO struct {
	ID          uint      `json:"id"`
	Code        string    `json:"code"`        // Code unique du logiciel
	Name        string    `json:"name"`        // Nom du logiciel
	Description *string   `json:"description,omitempty"` // Description
	Version     string    `json:"version,omitempty"`      // Version actuelle
	IsActive    bool      `json:"is_active"`              // Si le logiciel est actif
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateSoftwareRequest représente la requête de création d'un logiciel
type CreateSoftwareRequest struct {
	Code        string  `json:"code" binding:"required"`        // Code unique (obligatoire)
	Name        string  `json:"name" binding:"required"`         // Nom (obligatoire)
	Description *string `json:"description,omitempty"`          // Description (optionnel)
	Version     string  `json:"version,omitempty"`              // Version (optionnel)
}

// UpdateSoftwareRequest représente la requête de mise à jour d'un logiciel
type UpdateSoftwareRequest struct {
	Name        string  `json:"name,omitempty"`                 // Nom (optionnel)
	Description *string `json:"description,omitempty"`          // Description (optionnel)
	Version     string  `json:"version,omitempty"`              // Version (optionnel)
	IsActive    *bool   `json:"is_active,omitempty"`           // Si le logiciel est actif (optionnel)
}

// FilialeSoftwareDTO représente un déploiement de logiciel chez une filiale
type FilialeSoftwareDTO struct {
	ID         uint       `json:"id"`
	FilialeID  uint       `json:"filiale_id"`
	Filiale    FilialeDTO `json:"filiale"`
	SoftwareID uint       `json:"software_id"`
	Software   SoftwareDTO `json:"software"`
	Version    string     `json:"version,omitempty"`            // Version déployée
	DeployedAt *time.Time `json:"deployed_at,omitempty"`       // Date de déploiement
	IsActive   bool       `json:"is_active"`                    // Si le déploiement est actif
	Notes      *string    `json:"notes,omitempty"`              // Notes
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// CreateFilialeSoftwareRequest représente la requête de création d'un déploiement
type CreateFilialeSoftwareRequest struct {
	FilialeID  uint      `json:"filiale_id"`                      // ID de la filiale (peut venir de l'URL)
	SoftwareID uint      `json:"software_id" binding:"required"` // ID du logiciel (obligatoire)
	Version    string    `json:"version,omitempty"`               // Version déployée (optionnel)
	DeployedAt *time.Time `json:"deployed_at,omitempty"`         // Date de déploiement (optionnel)
	Notes      *string   `json:"notes,omitempty"`                // Notes (optionnel)
}

// UpdateFilialeSoftwareRequest représente la requête de mise à jour d'un déploiement
type UpdateFilialeSoftwareRequest struct {
	Version    string    `json:"version,omitempty"`              // Version déployée (optionnel)
	DeployedAt *time.Time `json:"deployed_at,omitempty"`        // Date de déploiement (optionnel)
	IsActive   *bool     `json:"is_active,omitempty"`            // Si le déploiement est actif (optionnel)
	Notes      *string   `json:"notes,omitempty"`                // Notes (optionnel)
}
