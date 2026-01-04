package dto

import (
	"time"
)

// AuditLogDTO représente un log d'audit
type AuditLogDTO struct {
	ID          uint                   `json:"id"`
	UserID      *uint                  `json:"user_id,omitempty"`
	User        *UserDTO               `json:"user,omitempty"`
	Action      string                 `json:"action"`
	EntityType  string                 `json:"entity_type"`
	EntityID    *uint                  `json:"entity_id,omitempty"`
	OldValues   map[string]interface{} `json:"old_values,omitempty" swaggertype:"object"`
	NewValues   map[string]interface{} `json:"new_values,omitempty" swaggertype:"object"`
	IPAddress   string                 `json:"ip_address,omitempty"`
	UserAgent   string                 `json:"user_agent,omitempty"`
	Description string                 `json:"description,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
}

// AuditLogListResponse représente la réponse de liste de logs d'audit avec pagination
type AuditLogListResponse struct {
	Logs       []AuditLogDTO  `json:"logs"`
	Pagination PaginationDTO  `json:"pagination"`
}

