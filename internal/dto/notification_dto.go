package dto

import "time"

// NotificationDTO représente une notification dans les réponses API
type NotificationDTO struct {
	ID        uint           `json:"id"`
	UserID    uint           `json:"user_id"`
	User      *UserDTO       `json:"user,omitempty"`     // Utilisateur (optionnel)
	Type      string         `json:"type"`               // delay_alert, budget_alert, validation_pending, etc.
	Title     string         `json:"title"`              // Titre de la notification
	Message   string         `json:"message"`            // Message de la notification
	IsRead    bool           `json:"is_read"`            // Si la notification a été lue
	ReadAt    *time.Time     `json:"read_at,omitempty"`  // Date de lecture (optionnel)
	LinkURL   string         `json:"link_url,omitempty"` // URL vers la ressource concernée (optionnel)
	Metadata  map[string]any `json:"metadata,omitempty"` // Données supplémentaires (optionnel)
	CreatedAt time.Time      `json:"created_at"`
}

// NotificationListResponse représente la réponse de liste de notifications
type NotificationListResponse struct {
	Notifications []NotificationDTO `json:"notifications"`
	UnreadCount   int               `json:"unread_count"` // Nombre de notifications non lues
	Total         int64             `json:"total"`
	Page          int               `json:"page"`
	Limit         int               `json:"limit"`
	TotalPages    int               `json:"total_pages"`
}

// MarkNotificationReadRequest représente la requête pour marquer une notification comme lue
// Pas besoin de body, mais on peut l'utiliser pour des actions futures
type MarkNotificationReadRequest struct {
	Read bool `json:"read,omitempty"` // true pour marquer comme lu, false pour non lu (optionnel)
}

// UnreadCountDTO représente le nombre de notifications non lues
type UnreadCountDTO struct {
	Count int `json:"count"` // Nombre de notifications non lues
}
