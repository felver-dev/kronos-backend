package dto

// OfficeDTO représente un siège/bureau
type OfficeDTO struct {
	ID        uint        `json:"id"`
	Name      string      `json:"name"`                 // Nom du siège (obligatoire)
	Code      *string     `json:"code,omitempty"`       // Code du siège (optionnel pour compatibilité)
	Country   string      `json:"country"`              // Pays (obligatoire)
	City      string      `json:"city"`                 // Ville (obligatoire)
	Commune   *string     `json:"commune,omitempty"`    // Commune (optionnel)
	Address   *string     `json:"address,omitempty"`    // Adresse complète (optionnel)
	FilialeID *uint       `json:"filiale_id,omitempty"` // ID de la filiale (optionnel)
	Filiale   *FilialeDTO `json:"filiale,omitempty"`    // Filiale associée (optionnel)
	Longitude *float64    `json:"longitude,omitempty"`  // Longitude (optionnel)
	Latitude  *float64    `json:"latitude,omitempty"`   // Latitude (optionnel)
	IsActive  bool        `json:"is_active"`            // Si le siège est actif
	CreatedAt string      `json:"created_at"`
	UpdatedAt string      `json:"updated_at"`
}

// CreateOfficeRequest représente la requête de création d'un siège
type CreateOfficeRequest struct {
	Name      string   `json:"name" binding:"required"`    // Nom du siège (obligatoire)
	Code      string   `json:"code,omitempty"`            // Code du siège (optionnel) - généré automatiquement si vide, sinon préfixé par le code filiale côté backend
	Country   string   `json:"country" binding:"required"` // Pays (obligatoire)
	City      string   `json:"city" binding:"required"`    // Ville (obligatoire)
	Commune   *string  `json:"commune,omitempty"`          // Commune (optionnel)
	Address   *string  `json:"address,omitempty"`          // Adresse complète (optionnel)
	FilialeID *uint    `json:"filiale_id,omitempty"`       // ID de la filiale (optionnel)
	Longitude *float64 `json:"longitude,omitempty"`        // Longitude (optionnel)
	Latitude  *float64 `json:"latitude,omitempty"`         // Latitude (optionnel)
	IsActive  *bool    `json:"is_active,omitempty"`        // Si le siège est actif (optionnel, défaut: true)
}

// UpdateOfficeRequest représente la requête de mise à jour d'un siège
type UpdateOfficeRequest struct {
	Name      *string  `json:"name,omitempty"`       // Nom du siège (optionnel)
	Code      *string  `json:"code,omitempty"`       // Code du siège (optionnel)
	Country   *string  `json:"country,omitempty"`    // Pays (optionnel)
	City      *string  `json:"city,omitempty"`       // Ville (optionnel)
	Commune   *string  `json:"commune,omitempty"`    // Commune (optionnel)
	Address   *string  `json:"address,omitempty"`    // Adresse complète (optionnel)
	FilialeID *uint    `json:"filiale_id,omitempty"` // ID de la filiale (optionnel)
	Longitude *float64 `json:"longitude,omitempty"`  // Longitude (optionnel)
	Latitude  *float64 `json:"latitude,omitempty"`   // Latitude (optionnel)
	IsActive  *bool    `json:"is_active,omitempty"`  // Si le siège est actif (optionnel)
}
