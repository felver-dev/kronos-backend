package models

import (
	"time"

	"gorm.io/gorm"
)

// AssetCategory représente une catégorie d'actif IT
// Table: asset_categories
type AssetCategory struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"type:varchar(100);not null" json:"name"`
	Description string `gorm:"type:text" json:"description,omitempty"`
	ParentID    *uint  `gorm:"index" json:"parent_id,omitempty"` // Catégorie parente (optionnel)
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relations
	Parent  *AssetCategory `gorm:"foreignKey:ParentID" json:"parent,omitempty"` // Catégorie parente (optionnel)
	Assets  []Asset         `gorm:"foreignKey:CategoryID" json:"-"`             // Actifs de cette catégorie
}

// TableName spécifie le nom de la table
func (AssetCategory) TableName() string {
	return "asset_categories"
}

// Asset représente un actif IT (équipement)
// Table: assets
type Asset struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	Name           string         `gorm:"type:varchar(255);not null" json:"name"`
	SerialNumber   string         `gorm:"type:varchar(100);index" json:"serial_number,omitempty"`
	Model          string         `gorm:"type:varchar(255)" json:"model,omitempty"`
	Manufacturer   string         `gorm:"type:varchar(255)" json:"manufacturer,omitempty"`
	CategoryID     uint           `gorm:"not null;index" json:"category_id"`
	AssignedToID   *uint          `gorm:"index" json:"assigned_to_id,omitempty"` // ID utilisateur assigné (optionnel)
	Status         string         `gorm:"type:varchar(50);default:'available';index" json:"status"` // available, in_use, maintenance, retired
	PurchaseDate   *time.Time     `gorm:"type:date" json:"purchase_date,omitempty"`
	WarrantyExpiry *time.Time      `gorm:"type:date" json:"warranty_expiry,omitempty"`
	Location       string         `gorm:"type:varchar(255)" json:"location,omitempty"`
	Notes          string         `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"` // Soft delete
	CreatedByID    *uint          `gorm:"index" json:"-"`
	CreatedBy      *User          `gorm:"foreignKey:CreatedByID" json:"-"`

	// Relations
	Category    AssetCategory `gorm:"foreignKey:CategoryID" json:"category,omitempty"` // Catégorie
	AssignedTo  *User         `gorm:"foreignKey:AssignedToID" json:"assigned_to,omitempty"` // Utilisateur assigné (optionnel)
}

// TableName spécifie le nom de la table
func (Asset) TableName() string {
	return "assets"
}

