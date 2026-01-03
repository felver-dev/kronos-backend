package repositories

import (
	"github.com/mcicare/itsm-backend/database"
	"github.com/mcicare/itsm-backend/internal/models"
)

// TicketAttachmentRepository interface pour les opérations sur les pièces jointes de tickets
type TicketAttachmentRepository interface {
	Create(attachment *models.TicketAttachment) error
	FindByID(id uint) (*models.TicketAttachment, error)
	FindByTicketID(ticketID uint) ([]models.TicketAttachment, error)
	FindByUserID(userID uint) ([]models.TicketAttachment, error)
	FindPrimaryByTicketID(ticketID uint) (*models.TicketAttachment, error)
	FindImagesByTicketID(ticketID uint) ([]models.TicketAttachment, error)
	Update(attachment *models.TicketAttachment) error
	Delete(id uint) error
}

// ticketAttachmentRepository implémente TicketAttachmentRepository
type ticketAttachmentRepository struct{}

// NewTicketAttachmentRepository crée une nouvelle instance de TicketAttachmentRepository
func NewTicketAttachmentRepository() TicketAttachmentRepository {
	return &ticketAttachmentRepository{}
}

// Create crée une nouvelle pièce jointe
func (r *ticketAttachmentRepository) Create(attachment *models.TicketAttachment) error {
	return database.DB.Create(attachment).Error
}

// FindByID trouve une pièce jointe par son ID
func (r *ticketAttachmentRepository) FindByID(id uint) (*models.TicketAttachment, error) {
	var attachment models.TicketAttachment
	err := database.DB.Preload("Ticket").Preload("User").Preload("User.Role").First(&attachment, id).Error
	if err != nil {
		return nil, err
	}
	return &attachment, nil
}

// FindByTicketID récupère toutes les pièces jointes d'un ticket
func (r *ticketAttachmentRepository) FindByTicketID(ticketID uint) ([]models.TicketAttachment, error) {
	var attachments []models.TicketAttachment
	err := database.DB.Preload("User").Preload("User.Role").Where("ticket_id = ?", ticketID).Order("display_order ASC, created_at ASC").Find(&attachments).Error
	return attachments, err
}

// FindByUserID récupère toutes les pièces jointes uploadées par un utilisateur
func (r *ticketAttachmentRepository) FindByUserID(userID uint) ([]models.TicketAttachment, error) {
	var attachments []models.TicketAttachment
	err := database.DB.Preload("Ticket").Preload("User").Where("user_id = ?", userID).Order("created_at DESC").Find(&attachments).Error
	return attachments, err
}

// FindPrimaryByTicketID trouve l'image principale d'un ticket
func (r *ticketAttachmentRepository) FindPrimaryByTicketID(ticketID uint) (*models.TicketAttachment, error) {
	var attachment models.TicketAttachment
	err := database.DB.Preload("User").Where("ticket_id = ? AND is_primary = ?", ticketID, true).First(&attachment).Error
	if err != nil {
		return nil, err
	}
	return &attachment, nil
}

// FindImagesByTicketID récupère toutes les images d'un ticket
func (r *ticketAttachmentRepository) FindImagesByTicketID(ticketID uint) ([]models.TicketAttachment, error) {
	var attachments []models.TicketAttachment
	err := database.DB.Preload("User").Preload("User.Role").Where("ticket_id = ? AND is_image = ?", ticketID, true).Order("display_order ASC, created_at ASC").Find(&attachments).Error
	return attachments, err
}

// Update met à jour une pièce jointe
func (r *ticketAttachmentRepository) Update(attachment *models.TicketAttachment) error {
	return database.DB.Save(attachment).Error
}

// Delete supprime une pièce jointe (soft delete)
func (r *ticketAttachmentRepository) Delete(id uint) error {
	return database.DB.Delete(&models.TicketAttachment{}, id).Error
}

