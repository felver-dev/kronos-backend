package dto

import "time"

// TicketDTO représente un ticket dans les réponses API
type TicketDTO struct {
	ID                  uint                `json:"id"`
	Code                string              `json:"code"` // Code unique: TKT-YYYY-NNNN
	Title               string              `json:"title"`
	Description         string              `json:"description"`
	Category            string              `json:"category"`                       // incident, demande, changement, developpement
	Source              string              `json:"source"`                         // mail, appel, direct
	Status              string              `json:"status"`                         // ouvert, en_cours, en_attente, cloture
	Priority            string              `json:"priority"`                       // low, medium, high, critical
	AssignedTo          *UserDTO            `json:"assigned_to,omitempty"`          // Utilisateur assigné (optionnel)
	Assignees           []TicketAssigneeDTO `json:"assignees,omitempty"`            // Utilisateurs assignés
	Lead                *UserDTO            `json:"lead,omitempty"`                 // Responsable (lead)
	CreatedBy           UserDTO             `json:"created_by"`                     // Créateur du ticket (informaticien)
	RequesterID         *uint               `json:"requester_id,omitempty"`         // ID du demandeur (relation vers users)
	Requester           *UserDTO            `json:"requester,omitempty"`            // Demandeur (relation vers users)
	RequesterName       string              `json:"requester_name,omitempty"`       // Nom de la personne qui a fait la demande (fallback pour demandeurs externes)
	RequesterDepartment string              `json:"requester_department,omitempty"` // Département du demandeur
	FilialeID           *uint               `json:"filiale_id,omitempty"`           // ID de la filiale
	Filiale             *FilialeDTO         `json:"filiale,omitempty"`              // Filiale (optionnel)
	SoftwareID          *uint               `json:"software_id,omitempty"`          // ID du logiciel concerné
	Software            *SoftwareDTO        `json:"software,omitempty"`             // Logiciel (optionnel)
	ValidatedByUserID   *uint               `json:"validated_by_user_id,omitempty"` // ID de l'utilisateur qui a validé
	ValidatedBy         *UserDTO            `json:"validated_by,omitempty"`         // Utilisateur qui a validé (optionnel)
	ValidatedAt         *time.Time          `json:"validated_at,omitempty"`         // Date de validation
	EstimatedTime       *int                `json:"estimated_time,omitempty"`       // Temps estimé en minutes (optionnel)
	ActualTime          *int                `json:"actual_time,omitempty"`          // Temps réel en minutes (optionnel)
	PrimaryImage        *string             `json:"primary_image,omitempty"`        // Image principale (optionnel)
	ParentID            *uint               `json:"parent_id,omitempty"`            // Ticket parent (optionnel)
	SubTickets          []TicketDTO         `json:"sub_tickets,omitempty"`          // Sous-tickets (optionnel)
	CreatedAt           time.Time           `json:"created_at"`
	UpdatedAt           time.Time           `json:"updated_at"`
	ClosedAt            *time.Time          `json:"closed_at,omitempty"`
}

// TicketAssigneeDTO représente une assignation d'un utilisateur à un ticket
type TicketAssigneeDTO struct {
	User   UserDTO `json:"user"`
	IsLead bool    `json:"is_lead"`
}

// CreateTicketRequest représente la requête de création d'un ticket
type CreateTicketRequest struct {
	Title               string `json:"title" binding:"required"`                                              // Titre (obligatoire)
	Description         string `json:"description" binding:"required"`                                        // Description (obligatoire)
	Category            string `json:"category" binding:"required"`                                           // Slug de la catégorie (doit exister dans ticket_categories et être active)
	Source              string `json:"source" binding:"required,oneof=mail appel direct whatsapp kronos"`     // Source (obligatoire)
	Priority            string `json:"priority,omitempty" binding:"omitempty,oneof=low medium high critical"` // Priorité (optionnel)
	EstimatedTime       *int   `json:"estimated_time,omitempty"`                                              // Temps estimé en minutes (optionnel)
	RequesterID         *uint  `json:"requester_id,omitempty"`                                                // ID du demandeur (optionnel, prioritaire sur requester_name)
	RequesterName       string `json:"requester_name,omitempty"`                                              // Nom de la personne qui a fait la demande (obligatoire si requester_id non fourni)
	RequesterDepartment string `json:"requester_department" binding:"required"`                               // Département du demandeur (obligatoire)
	FilialeID           *uint  `json:"filiale_id,omitempty"`                                                  // ID de la filiale (optionnel, défini automatiquement depuis l'utilisateur créateur)
	SoftwareID          *uint  `json:"software_id,omitempty"`                                                 // ID du logiciel concerné (optionnel)
	ParentID            *uint  `json:"parent_id,omitempty"`                                                   // Ticket parent (optionnel)
	AssigneeIDs         []uint `json:"assignee_ids,omitempty"`                                                // Assignés (optionnel)
	LeadID              *uint  `json:"lead_id,omitempty"`                                                     // Responsable (optionnel)
}

// UpdateTicketRequest représente la requête de mise à jour d'un ticket
type UpdateTicketRequest struct {
	Title               string `json:"title,omitempty"`                                                                      // Titre (optionnel)
	Description         string `json:"description,omitempty"`                                                                // Description (optionnel)
	Category            string `json:"category,omitempty" binding:"omitempty"`                                               // Slug de la catégorie (optionnel ; si fourni, doit exister et être active)
	Status              string `json:"status,omitempty" binding:"omitempty,oneof=ouvert en_cours en_attente resolu cloture"` // Statut (optionnel, ajout de "resolu")
	Priority            string `json:"priority,omitempty" binding:"omitempty,oneof=low medium high critical"`                // Priorité (optionnel)
	RequesterID         *uint  `json:"requester_id,omitempty"`                                                               // ID du demandeur (optionnel)
	RequesterName       string `json:"requester_name,omitempty"`                                                             // Nom du demandeur (optionnel, fallback)
	RequesterDepartment string `json:"requester_department,omitempty"`                                                       // Département du demandeur (optionnel)
	SoftwareID          *uint  `json:"software_id,omitempty"`                                                                // ID du logiciel concerné (optionnel)
	ParentID            *uint  `json:"parent_id,omitempty"`                                                                  // Ticket parent (optionnel)
	AssigneeIDs         []uint `json:"assignee_ids,omitempty"`                                                               // Assignés (optionnel)
	LeadID              *uint  `json:"lead_id,omitempty"`                                                                    // Responsable (optionnel)
	EstimatedTime       *int   `json:"estimated_time,omitempty"`                                                             // Temps estimé en minutes (optionnel, résolveurs IT)
}

// AssignTicketRequest représente la requête d'assignation d'un ticket
type AssignTicketRequest struct {
	UserID        uint   `json:"user_id,omitempty"`        // ID utilisateur (ancien mode)
	UserIDs       []uint `json:"user_ids,omitempty"`       // Liste d'utilisateurs assignés
	LeadID        *uint  `json:"lead_id,omitempty"`        // Responsable (optionnel)
	EstimatedTime *int   `json:"estimated_time,omitempty"` // Temps estimé en minutes (optionnel)
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

// UpdateTicketCommentRequest représente la requête de mise à jour d'un commentaire (texte uniquement)
type UpdateTicketCommentRequest struct {
	Comment string `json:"comment" binding:"required"` // Nouveau texte du commentaire
}

// TicketHistoryDTO représente une entrée d'historique d'un ticket
type TicketHistoryDTO struct {
	ID          uint      `json:"id"`
	TicketID    uint      `json:"ticket_id"`
	User        UserDTO   `json:"user"`
	Action      string    `json:"action"`      // created, updated, status_changed, assigned, etc.
	FieldName   string    `json:"field_name"`  // Nom du champ modifié (optionnel)
	OldValue    string    `json:"old_value"`   // Ancienne valeur (optionnel)
	NewValue    string    `json:"new_value"`   // Nouvelle valeur (optionnel)
	Description string    `json:"description"` // Description de l'action (optionnel)
	CreatedAt   time.Time `json:"created_at"`
}

// TicketAttachmentDTO représente une pièce jointe d'un ticket
type TicketAttachmentDTO struct {
	ID            uint      `json:"id"`
	TicketID      uint      `json:"ticket_id"`
	User          UserDTO   `json:"user"`
	FileName      string    `json:"file_name"`
	FilePath      string    `json:"file_path"`
	ThumbnailPath string    `json:"thumbnail_path,omitempty"`
	FileSize      *int      `json:"file_size,omitempty"`
	MimeType      string    `json:"mime_type,omitempty"`
	IsImage       bool      `json:"is_image"`
	DisplayOrder  int       `json:"display_order"`
	Description   string    `json:"description,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

// CreateTicketAttachmentRequest représente la requête de création d'une pièce jointe
type CreateTicketAttachmentRequest struct {
	Description  string `json:"description,omitempty"`   // Description (optionnel)
	DisplayOrder int    `json:"display_order,omitempty"` // Ordre d'affichage (optionnel)
}

// UpdateTicketAttachmentRequest représente la requête de mise à jour d'une pièce jointe
type UpdateTicketAttachmentRequest struct {
	Description  string `json:"description,omitempty"`   // Description (optionnel)
	DisplayOrder int    `json:"display_order,omitempty"` // Ordre d'affichage (optionnel)
}

// ReorderTicketAttachmentsRequest représente la requête de réorganisation des pièces jointes
type ReorderTicketAttachmentsRequest struct {
	AttachmentIDs []uint `json:"attachment_ids" binding:"required"` // Liste des IDs dans le nouvel ordre
}

// TicketSolutionDTO représente une solution documentée pour un ticket
type TicketSolutionDTO struct {
	ID        uint      `json:"id"`
	TicketID  uint      `json:"ticket_id"`
	Solution  string    `json:"solution"` // Solution documentée (Markdown)
	CreatedBy UserDTO   `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateTicketSolutionRequest représente la requête de création d'une solution
type CreateTicketSolutionRequest struct {
	Solution string `json:"solution" binding:"required"` // Solution documentée (obligatoire)
}

// UpdateTicketSolutionRequest représente la requête de mise à jour d'une solution
type UpdateTicketSolutionRequest struct {
	Solution string `json:"solution" binding:"required"` // Solution documentée (obligatoire)
}

// PublishSolutionToKBRequest représente la requête de publication d'une solution dans la base de connaissances
type PublishSolutionToKBRequest struct {
	Title      string `json:"title" binding:"required"`       // Titre de l'article (obligatoire)
	CategoryID uint   `json:"category_id" binding:"required"` // ID de la catégorie KB (obligatoire)
}
