package dto

import "time"

// ChangeDTO représente un changement dans les réponses API
type ChangeDTO struct {
	ID                uint       `json:"id"`
	TicketID          uint       `json:"ticket_id"`
	Ticket            *TicketDTO `json:"ticket,omitempty"`             // Ticket associé (optionnel)
	Risk              string     `json:"risk"`                         // low, medium, high, critical
	RiskDescription   string     `json:"risk_description,omitempty"`   // Description du risque (optionnel)
	ResponsibleID     *uint      `json:"responsible_id,omitempty"`     // ID du responsable (optionnel)
	Responsible       *UserDTO   `json:"responsible,omitempty"`        // Responsable (optionnel)
	Result            string     `json:"result,omitempty"`             // success, partial, failed, rolled_back
	ResultDescription string     `json:"result_description,omitempty"` // Description du résultat (optionnel)
	ResultDate        *time.Time `json:"result_date,omitempty"`        // Date du résultat (optionnel)
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

// CreateChangeRequest représente la requête de création d'un changement
type CreateChangeRequest struct {
	TicketID        uint   `json:"ticket_id" binding:"required"`                           // ID du ticket (obligatoire)
	Risk            string `json:"risk" binding:"required,oneof=low medium high critical"` // Risque (obligatoire)
	RiskDescription string `json:"risk_description,omitempty"`                             // Description du risque (optionnel)
}

// UpdateChangeRequest représente la requête de mise à jour d'un changement
type UpdateChangeRequest struct {
	Risk            string `json:"risk,omitempty" binding:"omitempty,oneof=low medium high critical"` // Risque (optionnel)
	RiskDescription string `json:"risk_description,omitempty"`                                        // Description du risque (optionnel)
}

// AssignResponsibleRequest représente la requête d'assignation d'un responsable
type AssignResponsibleRequest struct {
	UserID uint `json:"user_id" binding:"required"` // ID de l'utilisateur responsable (obligatoire)
}

// UpdateRiskRequest représente la requête de mise à jour du risque
type UpdateRiskRequest struct {
	Risk            string `json:"risk" binding:"required,oneof=low medium high critical"` // Risque (obligatoire)
	RiskDescription string `json:"risk_description,omitempty"`                             // Description du risque (optionnel)
}

// RecordChangeResultRequest représente la requête d'enregistrement du résultat post-changement
type RecordChangeResultRequest struct {
	Result      string `json:"result" binding:"required,oneof=success partial failed rolled_back"` // Résultat (obligatoire)
	Description string `json:"description" binding:"required"`                                     // Description (obligatoire)
	Issues      string `json:"issues,omitempty"`                                                   // Problèmes rencontrés (optionnel)
}

// ChangeResultDTO représente le résultat d'un changement
type ChangeResultDTO struct {
	Result         string    `json:"result"`                     // success, partial, failed, rolled_back
	Description    string    `json:"description"`                // Description du résultat
	Issues         string    `json:"issues,omitempty"`           // Problèmes rencontrés (optionnel)
	Date           time.Time `json:"date"`                       // Date d'enregistrement
	RecordedBy     uint      `json:"recorded_by"`                // ID de l'utilisateur qui a enregistré
	RecordedByUser *UserDTO  `json:"recorded_by_user,omitempty"` // Utilisateur (optionnel)
}
