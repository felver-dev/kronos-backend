package dto

import "time"

// FilialeDTO représente une filiale dans les réponses API
type FilialeDTO struct {
	ID          uint      `json:"id"`
	Code        string    `json:"code"`        // Code unique de la filiale
	Name        string    `json:"name"`        // Nom de la filiale
	Country     string    `json:"country,omitempty"`     // Pays
	City        string    `json:"city,omitempty"`        // Ville
	Address     *string   `json:"address,omitempty"`     // Adresse complète
	Phone       string    `json:"phone,omitempty"`       // Téléphone
	Email       string    `json:"email,omitempty"`       // Email de contact
	IsActive    bool      `json:"is_active"`             // Si la filiale est active
	IsSoftwareProvider bool    `json:"is_software_provider"`         // Filiale fournisseur de logiciels / IT
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateFilialeRequest représente la requête de création d'une filiale
type CreateFilialeRequest struct {
	Code        string  `json:"code" binding:"required"`        // Code unique (obligatoire)
	Name        string  `json:"name" binding:"required"`         // Nom (obligatoire)
	Country     string  `json:"country,omitempty"`               // Pays (optionnel)
	City        string  `json:"city,omitempty"`                 // Ville (optionnel)
	Address     *string `json:"address,omitempty"`              // Adresse (optionnel)
	Phone       string  `json:"phone,omitempty"`                 // Téléphone (optionnel)
	Email       string  `json:"email,omitempty"`                 // Email (optionnel)
	IsSoftwareProvider bool   `json:"is_software_provider,omitempty"`       // Filiale fournisseur de logiciels (optionnel)
}

// UpdateFilialeRequest représente la requête de mise à jour d'une filiale
type UpdateFilialeRequest struct {
	Name        string  `json:"name,omitempty"`                 // Nom (optionnel)
	Country     string  `json:"country,omitempty"`              // Pays (optionnel)
	City        string  `json:"city,omitempty"`                 // Ville (optionnel)
	Address     *string `json:"address,omitempty"`               // Adresse (optionnel)
	Phone       string  `json:"phone,omitempty"`                 // Téléphone (optionnel)
	Email       string  `json:"email,omitempty"`                 // Email (optionnel)
	IsActive    *bool   `json:"is_active,omitempty"`            // Si la filiale est active (optionnel)
	IsSoftwareProvider *bool  `json:"is_software_provider,omitempty"`       // Filiale fournisseur de logiciels (optionnel)
}
