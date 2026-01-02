package dto

import "time"

// TicketDTO représente un ticket dans les réponses API
type TicketDTO struct {
	ID            uint       `json:"id"`
	Title         string     `json:"title"`
	Description   string     `json:"description"`
	Category      string     `json:"category"`                 // incident, demande, changement, developpement
	Source        string     `json:"source"`                   // mail, appel, direct
	Status        string     `json:"status"`                   // ouvert, en_cours, en_attente, cloture
	Priority      string     `json:"priority"`                 // low, medium, high, critical
	AssignedTo    *UserDTO   `json:"assigned_to,omitempty"`    // Utilisateur assigné (optionnel)
	CreatedBy     UserDTO    `json:"created_by"`               // Créateur du ticket
	EstimatedTime *int       `json:"estimated_time,omitempty"` // Temps estimé en minutes (optionnel)
	ActualTime    *int       `json:"actual_time,omitempty"`    // Temps réel en minutes (optionnel)
	PrimaryImage  *string    `json:"primary_image,omitempty"`  // Image principale (optionnel)
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	ClosedAt      *time.Time `json:"closed_at,omitempty"`
}

// CreateTicketRequest représente la requête de création d'un ticket
type CreateTicketRequest struct {
	Title         string `json:"title" binding:"required"`                                                    // Titre (obligatoire)
	Description   string `json:"description" binding:"required"`                                              // Description (obligatoire)
	Category      string `json:"category" binding:"required,oneof=incident demande changement developpement"` // Catégorie (obligatoire)
	Source        string `json:"source" binding:"required,oneof=mail appel direct"`                           // Source (obligatoire)
	Priority      string `json:"priority,omitempty" binding:"omitempty,oneof=low medium high critical"`       // Priorité (optionnel)
	EstimatedTime *int   `json:"estimated_time,omitempty"`                                                    // Temps estimé en minutes (optionnel)
}

// UpdateTicketRequest représente la requête de mise à jour d'un ticket
type UpdateTicketRequest struct {
	Title       string `json:"title,omitempty"`                                                               // Titre (optionnel)
	Description string `json:"description,omitempty"`                                                         // Description (optionnel)
	Status      string `json:"status,omitempty" binding:"omitempty,oneof=ouvert en_cours en_attente cloture"` // Statut (optionnel)
	Priority    string `json:"priority,omitempty" binding:"omitempty,oneof=low medium high critical"`         // Priorité (optionnel)
}

// AssignTicketRequest représente la requête d'assignation d'un ticket
type AssignTicketRequest struct {
	UserID        uint `json:"user_id" binding:"required"` // ID de l'utilisateur à assigner (obligatoire)
	EstimatedTime *int `json:"estimated_time,omitempty"`   // Temps estimé en minutes (optionnel)
}

// TicketListResponse représente la réponse de liste de tickets avec pagination
type TicketListResponse struct {
	Tickets    []TicketDTO   `json:"tickets"`
	Pagination PaginationDTO `json:"pagination"`
}

// TicketCommentDTO représente un commentaire sur un ticket
type TicketCommentDTO struct {
	ID         uint      `json:"id"`
	TicketID   uint      `json:"ticket_id"`
	User       UserDTO   `json:"user"`
	Comment    string    `json:"comment"`
	IsInternal bool      `json:"is_internal"` // Commentaire interne (visible uniquement par l'IT)
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// CreateTicketCommentRequest représente la requête de création d'un commentaire
type CreateTicketCommentRequest struct {
	Comment    string `json:"comment" binding:"required"` // Commentaire (obligatoire)
	IsInternal bool   `json:"is_internal,omitempty"`      // Commentaire interne (optionnel, défaut: false)
}
