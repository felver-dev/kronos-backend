package dto

import "time"

// AssetDTO représente un actif IT dans les réponses API
type AssetDTO struct {
	ID             uint       `json:"id"`
	Name           string     `json:"name"`
	SerialNumber   string     `json:"serial_number,omitempty"`
	Model          string     `json:"model,omitempty"`
	Manufacturer   string     `json:"manufacturer,omitempty"`
	CategoryID     uint       `json:"category_id"`
	Category       *AssetCategoryDTO `json:"category,omitempty"` // Catégorie (optionnel)
	AssignedTo     *uint      `json:"assigned_to,omitempty"`      // ID utilisateur assigné (optionnel)
	AssignedUser   *UserDTO   `json:"assigned_user,omitempty"`   // Utilisateur assigné (optionnel)
	Status         string     `json:"status"`                    // available, in_use, maintenance, retired
	PurchaseDate   *time.Time `json:"purchase_date,omitempty"`
	WarrantyExpiry *time.Time `json:"warranty_expiry,omitempty"`
	Location       string     `json:"location,omitempty"`
	Notes          string     `json:"notes,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// AssetCategoryDTO représente une catégorie d'actif
type AssetCategoryDTO struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	ParentID    *uint  `json:"parent_id,omitempty"` // Catégorie parente (optionnel)
}

// CreateAssetRequest représente la requête de création d'un actif
type CreateAssetRequest struct {
	Name           string    `json:"name" binding:"required"`                    // Nom (obligatoire)
	SerialNumber   string    `json:"serial_number,omitempty"`                   // Numéro de série (optionnel)
	Model          string    `json:"model,omitempty"`                            // Modèle (optionnel)
	Manufacturer   string    `json:"manufacturer,omitempty"`                    // Fabricant (optionnel)
	CategoryID     uint      `json:"category_id" binding:"required"`             // ID catégorie (obligatoire)
	AssignedTo     *uint     `json:"assigned_to,omitempty"`                      // ID utilisateur (optionnel)
	Status         string    `json:"status,omitempty" binding:"omitempty,oneof=available in_use maintenance retired"` // Statut (optionnel)
	PurchaseDate   *string   `json:"purchase_date,omitempty"`                   // Date d'achat format "2006-01-02" (optionnel)
	WarrantyExpiry *string   `json:"warranty_expiry,omitempty"`                 // Date expiration garantie format "2006-01-02" (optionnel)
	Location       string    `json:"location,omitempty"`                        // Localisation (optionnel)
	Notes          string    `json:"notes,omitempty"`                            // Notes (optionnel)
}

// UpdateAssetRequest représente la requête de mise à jour d'un actif
type UpdateAssetRequest struct {
	Name           string  `json:"name,omitempty"`
	SerialNumber   string  `json:"serial_number,omitempty"`
	Model          string  `json:"model,omitempty"`
	Manufacturer   string  `json:"manufacturer,omitempty"`
	CategoryID     *uint   `json:"category_id,omitempty"`
	AssignedTo     *uint   `json:"assigned_to,omitempty"` // nil pour retirer l'assignation
	Status         string  `json:"status,omitempty" binding:"omitempty,oneof=available in_use maintenance retired"`
	PurchaseDate   *string `json:"purchase_date,omitempty"`
	WarrantyExpiry *string `json:"warranty_expiry,omitempty"`
	Location       string  `json:"location,omitempty"`
	Notes          string  `json:"notes,omitempty"`
}

// AssignAssetRequest représente la requête d'assignation d'un actif à un utilisateur
type AssignAssetRequest struct {
	UserID uint `json:"user_id" binding:"required"` // ID de l'utilisateur (obligatoire)
}

