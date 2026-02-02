package dto

// DepartmentDTO représente un département
type DepartmentDTO struct {
	ID             uint        `json:"id"`
	Name           string      `json:"name"`                  // Nom du département (obligatoire)
	Code           string      `json:"code"`                  // Code unique du département (obligatoire)
	Description    *string     `json:"description,omitempty"` // Description (optionnel)
	FilialeID      *uint       `json:"filiale_id,omitempty"`  // ID de la filiale (optionnel)
	Filiale        *FilialeDTO `json:"filiale,omitempty"`     // Filiale associée (optionnel)
	OfficeID       *uint       `json:"office_id,omitempty"`   // ID du siège (optionnel)
	Office         *OfficeDTO  `json:"office,omitempty"`      // Siège associé (optionnel)
	IsActive       bool        `json:"is_active"`             // Si le département est actif
	IsITDepartment bool        `json:"is_it_department"`      // Si c'est un département IT (uniquement pour la filiale fournisseur de logiciels)
	CreatedAt      string      `json:"created_at"`
	UpdatedAt      string      `json:"updated_at"`
}

// CreateDepartmentRequest représente la requête de création d'un département
type CreateDepartmentRequest struct {
	Name           string  `json:"name" binding:"required"`       // Nom du département (obligatoire)
	Code           string  `json:"code" binding:"required"`       // Code unique du département (obligatoire)
	Description    *string `json:"description,omitempty"`         // Description (optionnel)
	FilialeID      *uint   `json:"filiale_id" binding:"required"` // ID de la filiale (obligatoire)
	OfficeID       *uint   `json:"office_id,omitempty"`           // ID du siège (optionnel)
	IsActive       *bool   `json:"is_active,omitempty"`           // Si le département est actif (optionnel, défaut: true)
	IsITDepartment *bool   `json:"is_it_department,omitempty"`    // Si c'est un département IT (optionnel, défaut: false, uniquement pour MCI CARE CI)
}

// UpdateDepartmentRequest représente la requête de mise à jour d'un département
type UpdateDepartmentRequest struct {
	Name           *string `json:"name,omitempty"`             // Nom du département (optionnel)
	Code           *string `json:"code,omitempty"`             // Code unique du département (optionnel)
	Description    *string `json:"description,omitempty"`      // Description (optionnel)
	FilialeID      *uint   `json:"filiale_id,omitempty"`       // ID de la filiale (optionnel)
	OfficeID       *uint   `json:"office_id,omitempty"`        // ID du siège (optionnel)
	IsActive       *bool   `json:"is_active,omitempty"`        // Si le département est actif (optionnel)
	IsITDepartment *bool   `json:"is_it_department,omitempty"` // Si c'est un département IT (optionnel, uniquement pour la filiale fournisseur de logiciels)
}
