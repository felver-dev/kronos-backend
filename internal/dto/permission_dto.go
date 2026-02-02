package dto

import "time"

// PermissionDTO représente une permission dans les réponses API
type PermissionDTO struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Description string    `json:"description,omitempty"`
	Module      string    `json:"module,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}
