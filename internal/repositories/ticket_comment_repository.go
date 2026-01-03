package repositories

import (
	"github.com/mcicare/itsm-backend/database"
	"github.com/mcicare/itsm-backend/internal/models"
)

// TicketCommentRepository interface pour les opérations sur les commentaires de tickets
type TicketCommentRepository interface {
	Create(comment *models.TicketComment) error
	FindByID(id uint) (*models.TicketComment, error)
	FindByTicketID(ticketID uint) ([]models.TicketComment, error)
	FindByUserID(userID uint) ([]models.TicketComment, error)
	FindInternalByTicketID(ticketID uint) ([]models.TicketComment, error)
	FindPublicByTicketID(ticketID uint) ([]models.TicketComment, error)
	Update(comment *models.TicketComment) error
	Delete(id uint) error
}

// ticketCommentRepository implémente TicketCommentRepository
type ticketCommentRepository struct{}

// NewTicketCommentRepository crée une nouvelle instance de TicketCommentRepository
func NewTicketCommentRepository() TicketCommentRepository {
	return &ticketCommentRepository{}
}

// Create crée un nouveau commentaire
func (r *ticketCommentRepository) Create(comment *models.TicketComment) error {
	return database.DB.Create(comment).Error
}

// FindByID trouve un commentaire par son ID
func (r *ticketCommentRepository) FindByID(id uint) (*models.TicketComment, error) {
	var comment models.TicketComment
	err := database.DB.Preload("Ticket").Preload("User").Preload("User.Role").First(&comment, id).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// FindByTicketID récupère tous les commentaires d'un ticket
func (r *ticketCommentRepository) FindByTicketID(ticketID uint) ([]models.TicketComment, error) {
	var comments []models.TicketComment
	err := database.DB.Preload("User").Preload("User.Role").Where("ticket_id = ?", ticketID).Order("created_at ASC").Find(&comments).Error
	return comments, err
}

// FindByUserID récupère tous les commentaires d'un utilisateur
func (r *ticketCommentRepository) FindByUserID(userID uint) ([]models.TicketComment, error) {
	var comments []models.TicketComment
	err := database.DB.Preload("Ticket").Preload("User").Where("user_id = ?", userID).Order("created_at DESC").Find(&comments).Error
	return comments, err
}

// FindInternalByTicketID récupère les commentaires internes d'un ticket
func (r *ticketCommentRepository) FindInternalByTicketID(ticketID uint) ([]models.TicketComment, error) {
	var comments []models.TicketComment
	err := database.DB.Preload("User").Preload("User.Role").Where("ticket_id = ? AND is_internal = ?", ticketID, true).Order("created_at ASC").Find(&comments).Error
	return comments, err
}

// FindPublicByTicketID récupère les commentaires publics d'un ticket
func (r *ticketCommentRepository) FindPublicByTicketID(ticketID uint) ([]models.TicketComment, error) {
	var comments []models.TicketComment
	err := database.DB.Preload("User").Preload("User.Role").Where("ticket_id = ? AND is_internal = ?", ticketID, false).Order("created_at ASC").Find(&comments).Error
	return comments, err
}

// Update met à jour un commentaire
func (r *ticketCommentRepository) Update(comment *models.TicketComment) error {
	return database.DB.Save(comment).Error
}

// Delete supprime un commentaire (soft delete)
func (r *ticketCommentRepository) Delete(id uint) error {
	return database.DB.Delete(&models.TicketComment{}, id).Error
}

