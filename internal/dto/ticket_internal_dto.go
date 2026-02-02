package dto

import "time"

// TicketInternalDTO représente un ticket interne dans les réponses API
type TicketInternalDTO struct {
	ID                uint        `json:"id"`
	Code              string      `json:"code"`
	Title             string      `json:"title"`
	Description       string      `json:"description"`
	Category          string      `json:"category"`
	Status            string      `json:"status"`
	Priority          string      `json:"priority"`
	DepartmentID      uint        `json:"department_id"`
	Department        *DepartmentDTO `json:"department,omitempty"`
	FilialeID         uint        `json:"filiale_id"`
	Filiale           *FilialeDTO `json:"filiale,omitempty"`
	CreatedByID       uint        `json:"created_by_id"`
	CreatedBy         UserDTO     `json:"created_by"`
	AssignedToID      *uint       `json:"assigned_to_id,omitempty"`
	AssignedTo        *UserDTO    `json:"assigned_to,omitempty"`
	ValidatedByUserID *uint       `json:"validated_by_user_id,omitempty"`
	ValidatedBy       *UserDTO    `json:"validated_by,omitempty"`
	ValidatedAt       *time.Time  `json:"validated_at,omitempty"`
	EstimatedTime     *int        `json:"estimated_time,omitempty"`
	ActualTime        *int        `json:"actual_time,omitempty"`
	TicketID          *uint       `json:"ticket_id,omitempty"`
	CreatedAt         time.Time   `json:"created_at"`
	UpdatedAt         time.Time   `json:"updated_at"`
	ClosedAt          *time.Time  `json:"closed_at,omitempty"`
}

// CreateTicketInternalRequest représente la requête de création d'un ticket interne
type CreateTicketInternalRequest struct {
	Title         string `json:"title" binding:"required"`
	Description   string `json:"description" binding:"required"`
	Category      string `json:"category" binding:"required"` // slug: tache_interne, demande_interne, etc.
	Priority      string `json:"priority,omitempty" binding:"omitempty,oneof=low medium high critical"`
	DepartmentID  uint   `json:"department_id" binding:"required"` // Département propriétaire (non-IT)
	EstimatedTime *int   `json:"estimated_time,omitempty"`
	AssignedToID  *uint  `json:"assigned_to_id,omitempty"`
	TicketID      *uint  `json:"ticket_id,omitempty"` // Lien optionnel vers un ticket normal
}

// UpdateTicketInternalRequest représente la requête de mise à jour d'un ticket interne
type UpdateTicketInternalRequest struct {
	Title         string `json:"title,omitempty"`
	Description   string `json:"description,omitempty"`
	Category      string `json:"category,omitempty"`
	Status        string `json:"status,omitempty" binding:"omitempty,oneof=ouvert en_cours en_attente resolu cloture"`
	Priority      string `json:"priority,omitempty" binding:"omitempty,oneof=low medium high critical"`
	EstimatedTime *int   `json:"estimated_time,omitempty"` // en minutes
	ActualTime    *int   `json:"actual_time,omitempty"`    // temps passé en minutes (saisi par l'assigné ou son chef)
	AssignedToID  *uint  `json:"assigned_to_id,omitempty"`
}

// AssignTicketInternalRequest représente la requête d'assignation d'un ticket interne
type AssignTicketInternalRequest struct {
	AssignedToID  *uint `json:"assigned_to_id,omitempty"`
	EstimatedTime *int  `json:"estimated_time,omitempty"`
}

// TicketInternalListResponse représente la réponse de liste avec pagination
type TicketInternalListResponse struct {
	Tickets    []TicketInternalDTO `json:"tickets"`
	Pagination PaginationDTO       `json:"pagination"`
}

// TicketInternalPerformanceDTO représente la performance de l'utilisateur sur les tickets internes qu'il traite (assignés à lui)
type TicketInternalPerformanceDTO struct {
	TotalAssigned   int     `json:"total_assigned"`    // Nombre total de tickets internes assignés à l'utilisateur
	Resolved        int     `json:"resolved"`          // Nombre clôturés
	InProgress      int     `json:"in_progress"`      // En cours
	Open            int     `json:"open"`             // Ouverts / en attente
	TotalTimeSpent  int     `json:"total_time_spent"` // Temps total passé en minutes (somme actual_time)
	Efficiency      float64 `json:"efficiency"`       // Efficacité en % (resolved / total_assigned si total > 0)
}
