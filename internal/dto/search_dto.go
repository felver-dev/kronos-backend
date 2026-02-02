package dto

import "time"

// GlobalSearchResultDTO représente le résultat d'une recherche globale
type GlobalSearchResultDTO struct {
	Query   string                        `json:"query"`
	Types   []string                      `json:"types"`
	Tickets []TicketSearchResultDTO       `json:"tickets,omitempty"`
	Assets  []AssetSearchResultDTO        `json:"assets,omitempty"`
	Articles []KnowledgeArticleSearchResultDTO `json:"articles,omitempty"`
	Users   []UserSearchResultDTO         `json:"users,omitempty"`
	TimeEntries []TimeEntrySearchResultDTO `json:"time_entries,omitempty"`
	Total   int                           `json:"total"`
}

// UserSearchResultDTO représente un résultat de recherche d'utilisateur
type UserSearchResultDTO struct {
	ID        uint       `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	FirstName string     `json:"first_name,omitempty"`
	LastName  string     `json:"last_name,omitempty"`
	Department *DepartmentDTO `json:"department,omitempty"`
	Role      string     `json:"role"`
	IsActive  bool       `json:"is_active"`
	Snippet   string     `json:"snippet,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

// TimeEntrySearchResultDTO représente un résultat de recherche d'entrée de temps
type TimeEntrySearchResultDTO struct {
	ID          uint       `json:"id"`
	TicketID    uint       `json:"ticket_id"`
	Ticket      *TicketDTO `json:"ticket,omitempty"`
	UserID      uint       `json:"user_id"`
	User        *UserDTO   `json:"user,omitempty"`
	TimeSpent   int        `json:"time_spent"`
	Date        time.Time  `json:"date"`
	Description string     `json:"description,omitempty"`
	Snippet     string     `json:"snippet,omitempty"`
	Validated   bool       `json:"validated"`
	CreatedAt   time.Time  `json:"created_at"`
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

