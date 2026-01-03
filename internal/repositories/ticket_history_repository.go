package repositories

import (
	"github.com/mcicare/itsm-backend/database"
	"github.com/mcicare/itsm-backend/internal/models"
)

// TicketHistoryRepository interface pour les opérations sur l'historique des tickets
type TicketHistoryRepository interface {
	Create(history *models.TicketHistory) error
	FindByID(id uint) (*models.TicketHistory, error)
	FindByTicketID(ticketID uint) ([]models.TicketHistory, error)
	FindByUserID(userID uint) ([]models.TicketHistory, error)
	FindByAction(ticketID uint, action string) ([]models.TicketHistory, error)
	Update(history *models.TicketHistory) error
	Delete(id uint) error
}

// ticketHistoryRepository implémente TicketHistoryRepository
type ticketHistoryRepository struct{}

// NewTicketHistoryRepository crée une nouvelle instance de TicketHistoryRepository
func NewTicketHistoryRepository() TicketHistoryRepository {
	return &ticketHistoryRepository{}
}

// Create crée une nouvelle entrée d'historique
func (r *ticketHistoryRepository) Create(history *models.TicketHistory) error {
	return database.DB.Create(history).Error
}

// FindByID trouve une entrée d'historique par son ID
func (r *ticketHistoryRepository) FindByID(id uint) (*models.TicketHistory, error) {
	var history models.TicketHistory
	err := database.DB.Preload("Ticket").Preload("User").Preload("User.Role").First(&history, id).Error
	if err != nil {
		return nil, err
	}
	return &history, nil
}

// FindByTicketID récupère tout l'historique d'un ticket
func (r *ticketHistoryRepository) FindByTicketID(ticketID uint) ([]models.TicketHistory, error) {
	var histories []models.TicketHistory
	err := database.DB.Preload("User").Preload("User.Role").Where("ticket_id = ?", ticketID).Order("created_at ASC").Find(&histories).Error
	return histories, err
}

// FindByUserID récupère l'historique des actions d'un utilisateur
func (r *ticketHistoryRepository) FindByUserID(userID uint) ([]models.TicketHistory, error) {
	var histories []models.TicketHistory
	err := database.DB.Preload("Ticket").Preload("User").Where("user_id = ?", userID).Order("created_at DESC").Find(&histories).Error
	return histories, err
}

// FindByAction récupère les entrées d'historique d'un ticket par action
func (r *ticketHistoryRepository) FindByAction(ticketID uint, action string) ([]models.TicketHistory, error) {
	var histories []models.TicketHistory
	err := database.DB.Preload("User").Preload("User.Role").Where("ticket_id = ? AND action = ?", ticketID, action).Order("created_at ASC").Find(&histories).Error
	return histories, err
}

// Update met à jour une entrée d'historique (rarement utilisé, l'historique est généralement immuable)
func (r *ticketHistoryRepository) Update(history *models.TicketHistory) error {
	return database.DB.Save(history).Error
}

// Delete supprime une entrée d'historique (rarement utilisé)
func (r *ticketHistoryRepository) Delete(id uint) error {
	return database.DB.Delete(&models.TicketHistory{}, id).Error
}
