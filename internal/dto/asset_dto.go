package dto

import "time"

// AssetDTO représente un actif IT dans les réponses API
type AssetDTO struct {
	ID             uint              `json:"id"`
	Name           string            `json:"name"`
	SerialNumber   string            `json:"serial_number,omitempty"`
	Model          string            `json:"model,omitempty"`
	Manufacturer   string            `json:"manufacturer,omitempty"`
	CategoryID     uint              `json:"category_id"`
	Category       *AssetCategoryDTO `json:"category,omitempty"`      // Catégorie (optionnel)
	AssignedTo     *uint             `json:"assigned_to,omitempty"`   // ID utilisateur assigné (optionnel)
	AssignedUser   *UserDTO          `json:"assigned_user,omitempty"` // Utilisateur assigné (optionnel)
	Status         string            `json:"status"`                  // available, in_use, maintenance, retired
	PurchaseDate   *time.Time        `json:"purchase_date,omitempty"`
	WarrantyExpiry *time.Time        `json:"warranty_expiry,omitempty"`
	Location       string            `json:"location,omitempty"`
	Notes          string            `json:"notes,omitempty"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
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
	Name           string  `json:"name" binding:"required"`                                                         // Nom (obligatoire)
	SerialNumber   string  `json:"serial_number,omitempty"`                                                         // Numéro de série (optionnel)
	Model          string  `json:"model,omitempty"`                                                                 // Modèle (optionnel)
	Manufacturer   string  `json:"manufacturer,omitempty"`                                                          // Fabricant (optionnel)
	CategoryID     uint    `json:"category_id" binding:"required"`                                                  // ID catégorie (obligatoire)
	AssignedTo     *uint   `json:"assigned_to,omitempty"`                                                           // ID utilisateur (optionnel)
	Status         string  `json:"status,omitempty" binding:"omitempty,oneof=available in_use maintenance retired"` // Statut (optionnel)
	PurchaseDate   *string `json:"purchase_date,omitempty"`                                                         // Date d'achat format "2006-01-02" (optionnel)
	WarrantyExpiry *string `json:"warranty_expiry,omitempty"`                                                       // Date expiration garantie format "2006-01-02" (optionnel)
	Location       string  `json:"location,omitempty"`                                                              // Localisation (optionnel)
	Notes          string  `json:"notes,omitempty"`                                                                 // Notes (optionnel)
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

// AssetInventoryDTO représente l'inventaire des actifs
type AssetInventoryDTO struct {
	Total      int            `json:"total"`       // Nombre total d'actifs
	ByStatus   map[string]int `json:"by_status"`   // Répartition par statut
	ByCategory map[string]int `json:"by_category"` // Répartition par catégorie
	Assigned   int            `json:"assigned"`    // Nombre d'actifs assignés
	Available  int            `json:"available"`   // Nombre d'actifs disponibles
}

// CreateAssetCategoryRequest représente la requête de création d'une catégorie d'actif
type CreateAssetCategoryRequest struct {
	Name        string `json:"name" binding:"required"` // Nom (obligatoire)
	Description string `json:"description,omitempty"`   // Description (optionnel)
	ParentID    *uint  `json:"parent_id,omitempty"`    // ID catégorie parente (optionnel)
}

// UpdateAssetCategoryRequest représente la requête de mise à jour d'une catégorie d'actif
type UpdateAssetCategoryRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	ParentID    *uint  `json:"parent_id,omitempty"` // nil pour retirer la catégorie parente
}

// DeleteAssetCategoryRequest représente la requête de suppression d'une catégorie d'actif
type DeleteAssetCategoryRequest struct {
	ConfirmName string `json:"confirm_name,omitempty"` // Nom de confirmation pour suppression en cascade
}

// AssetCategoryListResponse représente la réponse de liste de catégories avec pagination
type AssetCategoryListResponse struct {
	Categories []AssetCategoryDTO `json:"categories"`
	Pagination PaginationDTO      `json:"pagination"`
}

// AssetSoftwareDTO représente un logiciel installé sur un actif
type AssetSoftwareDTO struct {
	ID              uint       `json:"id"`
	AssetID         *uint      `json:"asset_id,omitempty"`                                      // ID de l'actif (optionnel)
	SoftwareName    string     `json:"software_name"`
	Version         string     `json:"version,omitempty"`
	LicenseKey      string     `json:"license_key,omitempty"`
	InstallationDate *time.Time `json:"installation_date,omitempty"`
	Notes           string     `json:"notes,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	Asset           *AssetDTO  `json:"asset,omitempty"`
}

// CreateAssetSoftwareRequest représente la requête de création d'un logiciel installé
type CreateAssetSoftwareRequest struct {
	AssetID         *uint   `json:"asset_id,omitempty"`                                          // ID de l'actif (optionnel - permet de créer des logiciels indépendamment)
	SoftwareName    string  `json:"software_name" binding:"required"`
	Version         string  `json:"version,omitempty"`
	LicenseKey      string  `json:"license_key,omitempty"`
	InstallationDate *string `json:"installation_date,omitempty"` // Format "2006-01-02"
	Notes           string  `json:"notes,omitempty"`
}

// UpdateAssetSoftwareRequest représente la requête de mise à jour d'un logiciel installé
type UpdateAssetSoftwareRequest struct {
	SoftwareName    string  `json:"software_name,omitempty"`
	Version         string  `json:"version,omitempty"`
	LicenseKey      string  `json:"license_key,omitempty"`
	InstallationDate *string `json:"installation_date,omitempty"` // Format "2006-01-02"
	Notes           string  `json:"notes,omitempty"`
}

// AssetSoftwareStatisticsDTO représente les statistiques sur les logiciels installés
type AssetSoftwareStatisticsDTO struct {
	BySoftware     []SoftwareCountDTO `json:"by_software"`
	BySoftwareName []SoftwareNameCountDTO `json:"by_software_name"`
}

// SoftwareCountDTO représente le nombre d'actifs ayant un logiciel spécifique
type SoftwareCountDTO struct {
	SoftwareName string `json:"software_name"`
	Version      string `json:"version"`
	Count        int64  `json:"count"`
	CategoryName string `json:"category_name"`
}

// SoftwareNameCountDTO représente le nombre d'actifs ayant un logiciel (toutes versions)
type SoftwareNameCountDTO struct {
	SoftwareName string `json:"software_name"`
	Count        int64  `json:"count"`
}