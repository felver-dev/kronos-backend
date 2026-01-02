package dto

import "time"

// UserDTO représente un utilisateur dans les réponses API
// C'est la version "publique" du modèle User, sans les informations sensibles
type UserDTO struct {
	ID        uint       `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	FirstName string     `json:"first_name,omitempty"`
	LastName  string     `json:"last_name,omitempty"`
	Avatar    string     `json:"avatar,omitempty"` // Chemin vers l'avatar
	Role      string     `json:"role"`             // Nom du rôle (ex: "DSI", "TECHNICIEN_IT")
	IsActive  bool       `json:"is_active"`
	LastLogin *time.Time `json:"last_login,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// CreateUserRequest représente la requête de création d'un utilisateur
type CreateUserRequest struct {
	Username  string `json:"username" binding:"required"`       // Nom d'utilisateur (obligatoire)
	Email     string `json:"email" binding:"required,email"`    // Email (obligatoire, format email)
	Password  string `json:"password" binding:"required,min=6"` // Mot de passe (obligatoire, min 6 caractères)
	FirstName string `json:"first_name,omitempty"`              // Prénom (optionnel)
	LastName  string `json:"last_name,omitempty"`               // Nom (optionnel)
	RoleID    uint   `json:"role_id" binding:"required"`        // ID du rôle (obligatoire)
}

// UpdateUserRequest représente la requête de mise à jour d'un utilisateur
type UpdateUserRequest struct {
	Email     string `json:"email,omitempty" binding:"omitempty,email"` // Email (optionnel, format email si fourni)
	FirstName string `json:"first_name,omitempty"`                      // Prénom (optionnel)
	LastName  string `json:"last_name,omitempty"`                       // Nom (optionnel)
	RoleID    uint   `json:"role_id,omitempty"`                         // ID du rôle (optionnel)
	IsActive  *bool  `json:"is_active,omitempty"`                       // Statut actif (optionnel, pointeur pour distinguer false de non fourni)
}
