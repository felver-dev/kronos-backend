package models

import (
	"time"

	"gorm.io/gorm"
)

// TicketAttachment représente une pièce jointe (image, document) d'un ticket
// Table: ticket_attachments
type TicketAttachment struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	TicketID      uint           `gorm:"not null;index" json:"ticket_id"`
	UserID        uint           `gorm:"not null;index" json:"user_id"`
	FileName      string         `gorm:"type:varchar(255);not null" json:"file_name"`
	FilePath      string         `gorm:"type:varchar(500);not null" json:"file_path"`
	ThumbnailPath string         `gorm:"type:varchar(500)" json:"thumbnail_path,omitempty"` // Chemin vers la miniature (pour les images)
	FileSize      *int           `gorm:"type:int" json:"file_size,omitempty"`               // Taille en bytes (optionnel)
	MimeType      string         `gorm:"type:varchar(100)" json:"mime_type,omitempty"`
	IsImage       bool           `gorm:"default:false;index" json:"is_image"`    // TRUE si c'est une image
	DisplayOrder  int            `gorm:"default:0" json:"display_order"`         // Ordre d'affichage (pour les galeries)
	Description   string         `gorm:"type:text" json:"description,omitempty"` // Description optionnelle
	CreatedAt     time.Time      `json:"created_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"` // Soft delete

	// Relations
	Ticket Ticket `gorm:"foreignKey:TicketID" json:"-"`  // Ticket associé
	User   User   `gorm:"foreignKey:UserID" json:"user"` // Utilisateur qui a uploadé
}

// TableName spécifie le nom de la table
func (TicketAttachment) TableName() string {
	return "ticket_attachments"
}
