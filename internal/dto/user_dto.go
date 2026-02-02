package dto

import "time"

// UserDTO représente un utilisateur dans les réponses API
// C'est la version "publique" du modèle User, sans les informations sensibles
type UserDTO struct {
	ID           uint           `json:"id"`
	Username     string         `json:"username"`
	Email        string         `json:"email"`
	Phone        string         `json:"phone,omitempty"`
	FirstName    string         `json:"first_name,omitempty"`
	LastName     string         `json:"last_name,omitempty"`
	DepartmentID *uint          `json:"department_id,omitempty"` // ID du département (optionnel)
	Department   *DepartmentDTO `json:"department,omitempty"`    // Département complet (optionnel)
	FilialeID    *uint          `json:"filiale_id,omitempty"`    // ID de la filiale (optionnel)
	Filiale      *FilialeDTO    `json:"filiale,omitempty"`       // Filiale complète (optionnel)
	Avatar       string         `json:"avatar,omitempty"`        // Chemin vers l'avatar
	Role         string         `json:"role"`                    // Nom du rôle (ex: "DSI", "TECHNICIEN_IT")
	Permissions  []string       `json:"permissions,omitempty"`   // Liste des permissions (optionnelle)
	IsActive     bool           `json:"is_active"`
	LastLogin    *time.Time     `json:"last_login,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

// CreateUserRequest représente la requête de création d'un utilisateur
type CreateUserRequest struct {
	Username     string `json:"username" binding:"required"`       // Nom d'utilisateur (obligatoire)
	Email        string `json:"email" binding:"required,email"`    // Email (obligatoire, format email)
	Password     string `json:"password" binding:"required,min=6"` // Mot de passe (obligatoire, min 6 caractères)
	FirstName    string `json:"first_name,omitempty"`              // Prénom (optionnel)
	LastName     string `json:"last_name,omitempty"`               // Nom (optionnel)
	Phone        string `json:"phone,omitempty"`                   // Téléphone (optionnel)
	DepartmentID *uint  `json:"department_id,omitempty"`           // ID du département (optionnel)
	FilialeID    *uint  `json:"filiale_id,omitempty"`              // ID de la filiale (optionnel)
	RoleID       uint   `json:"role_id"`                           // ID du rôle (optionnel, USER par défaut)
}

// UpdateUserRequest représente la requête de mise à jour d'un utilisateur
type UpdateUserRequest struct {
	Username     string `json:"username,omitempty"`                        // Nom d'utilisateur (optionnel)
	Email        string `json:"email,omitempty" binding:"omitempty,email"` // Email (optionnel, format email si fourni)
	FirstName    string `json:"first_name,omitempty"`                      // Prénom (optionnel)
	LastName     string `json:"last_name,omitempty"`                       // Nom (optionnel)
	Phone        string `json:"phone,omitempty"`                           // Téléphone (optionnel)
	DepartmentID *uint  `json:"department_id,omitempty"`                   // ID du département (optionnel, nil pour supprimer)
	FilialeID    *uint  `json:"filiale_id,omitempty"`                      // ID de la filiale (optionnel, nil pour supprimer)
	RoleID       uint   `json:"role_id,omitempty"`                         // ID du rôle (optionnel)
	IsActive     *bool  `json:"is_active,omitempty"`                       // Statut actif (optionnel, pointeur pour distinguer false de non fourni)
}

// RoleDTO représente un rôle dans les réponses API
type RoleDTO struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"` // Nom du rôle (ex: "DSI", "TECHNICIEN_IT")
	Description string    `json:"description,omitempty"`
	IsSystem    bool      `json:"is_system"`               // Si c'est un rôle système (ne peut pas être modifié/supprimé)
	CreatedByID *uint     `json:"created_by_id,omitempty"` // ID de l'utilisateur créateur (nil pour les rôles système)
	FilialeID   *uint     `json:"filiale_id,omitempty"`    // ID de la filiale (nil pour les rôles globaux)
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateRoleRequest représente la requête de création d'un rôle
type CreateRoleRequest struct {
	Name        string   `json:"name" binding:"required"` // Nom du rôle (obligatoire)
	Description string   `json:"description,omitempty"`   // Description (optionnel)
	Permissions []string `json:"permissions,omitempty"`   // Liste des codes de permissions à assigner (optionnel)
	FilialeID   *uint    `json:"filiale_id,omitempty"`    // ID de la filiale (optionnel, sera défini automatiquement si non fourni)
}

// UpdateRoleRequest représente la requête de mise à jour d'un rôle
type UpdateRoleRequest struct {
	Name        string `json:"name,omitempty"`        // Nom (optionnel)
	Description string `json:"description,omitempty"` // Description (optionnel)
}

// UserPermissionsDTO représente les permissions d'un utilisateur
type UserPermissionsDTO struct {
	UserID      uint     `json:"user_id"`
	User        *UserDTO `json:"user,omitempty"`
	Permissions []string `json:"permissions"` // Liste des permissions (ex: ["tickets.create", "tickets.update"])
}

// UpdateUserPermissionsRequest représente la requête de mise à jour des permissions d'un utilisateur
type UpdateUserPermissionsRequest struct {
	Permissions []string `json:"permissions" binding:"required"` // Liste des permissions (obligatoire)
}
