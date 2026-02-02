package dto

// TicketCategoryDTO représente une catégorie de ticket
type TicketCategoryDTO struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	Description  string `json:"description,omitempty"`
	Icon         string `json:"icon,omitempty"`
	Color        string `json:"color,omitempty"`
	IsActive     bool   `json:"is_active"`
	DisplayOrder int    `json:"display_order"`
}

// CreateTicketCategoryRequest représente la requête de création d'une catégorie
type CreateTicketCategoryRequest struct {
	Name         string `json:"name" binding:"required"`         // Nom (obligatoire)
	Slug         string `json:"slug" binding:"required"`        // Slug unique (obligatoire)
	Description  string `json:"description,omitempty"`           // Description (optionnel)
	Icon         string `json:"icon,omitempty"`                  // Nom de l'icône (optionnel)
	Color        string `json:"color,omitempty"`                 // Couleur (optionnel)
	IsActive     bool   `json:"is_active,omitempty"`             // Actif (optionnel, défaut: true)
	DisplayOrder int    `json:"display_order,omitempty"`         // Ordre d'affichage (optionnel, défaut: 0)
}

// UpdateTicketCategoryRequest représente la requête de mise à jour d'une catégorie
type UpdateTicketCategoryRequest struct {
	Name         string `json:"name,omitempty"`         // Nom (optionnel)
	Slug         string `json:"slug,omitempty"`        // Slug (optionnel)
	Description  string `json:"description,omitempty"` // Description (optionnel)
	Icon         string `json:"icon,omitempty"`       // Nom de l'icône (optionnel)
	Color        string `json:"color,omitempty"`       // Couleur (optionnel)
	IsActive     *bool  `json:"is_active,omitempty"`    // Actif (optionnel)
	DisplayOrder *int   `json:"display_order,omitempty"` // Ordre d'affichage (optionnel)
}
