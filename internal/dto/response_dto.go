package dto

// PaginationDTO contient les informations de pagination pour les listes
type PaginationDTO struct {
	Page       int   `json:"page"`        // Page actuelle (commence à 1)
	Limit      int   `json:"limit"`       // Nombre d'éléments par page
	Total      int64 `json:"total"`       // Nombre total d'éléments
	TotalPages int   `json:"total_pages"` // Nombre total de pages
}

// ErrorDTO représente une erreur dans les réponses API
type ErrorDTO struct {
	Code    string `json:"code,omitempty"`    // Code d'erreur (optionnel)
	Message string `json:"message"`           // Message d'erreur
	Field   string `json:"field,omitempty"`   // Champ concerné (optionnel)
	Details any    `json:"details,omitempty"` // Détails supplémentaires (optionnel)
}

// ValidationErrorDTO représente une erreur de validation de champ
// Utilisé pour les erreurs de validation des formulaires
type ValidationErrorDTO struct {
	Field   string `json:"field"`         // Nom du champ
	Message string `json:"message"`       // Message d'erreur
	Tag     string `json:"tag,omitempty"` // Tag de validation (ex: "required", "email")
}
