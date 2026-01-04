package dto

import "time"

// GlobalSearchResultDTO représente le résultat d'une recherche globale
type GlobalSearchResultDTO struct {
	Query   string                        `json:"query"`
	Types   []string                      `json:"types"`
	Tickets []TicketSearchResultDTO       `json:"tickets,omitempty"`
	Assets  []AssetSearchResultDTO        `json:"assets,omitempty"`
	Articles []KnowledgeArticleSearchResultDTO `json:"articles,omitempty"`
	Total   int                           `json:"total"`
}

// TicketSearchResultDTO représente un résultat de recherche de ticket
type TicketSearchResultDTO struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Snippet     string    `json:"snippet"`     // Extrait de la description correspondant
	Status      string    `json:"status"`
	Priority    string    `json:"priority"`
	Category    string    `json:"category"`
	CreatedBy   *UserDTO  `json:"created_by,omitempty"`
	AssignedTo  *UserDTO  `json:"assigned_to,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// AssetSearchResultDTO représente un résultat de recherche d'actif
type AssetSearchResultDTO struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Snippet     string    `json:"snippet"`     // Extrait de la description correspondant
	SerialNumber string   `json:"serial_number,omitempty"`
	CategoryID  uint      `json:"category_id"`
	Category    *AssetCategoryDTO `json:"category,omitempty"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

