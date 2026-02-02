package repositories

import (
	"github.com/mcicare/itsm-backend/database"
	"github.com/mcicare/itsm-backend/internal/models"
)

// TicketCategoryRepository interface pour les opérations sur les catégories de tickets
type TicketCategoryRepository interface {
	Create(category *models.TicketCategory) error
	FindByID(id uint) (*models.TicketCategory, error)
	FindBySlug(slug string) (*models.TicketCategory, error)
	FindAll() ([]models.TicketCategory, error)
	FindActive() ([]models.TicketCategory, error)
	Update(category *models.TicketCategory) error
	Delete(id uint) error
}

// ticketCategoryRepository implémente TicketCategoryRepository
type ticketCategoryRepository struct{}

// NewTicketCategoryRepository crée une nouvelle instance de TicketCategoryRepository
func NewTicketCategoryRepository() TicketCategoryRepository {
	return &ticketCategoryRepository{}
}

// Create crée une nouvelle catégorie
func (r *ticketCategoryRepository) Create(category *models.TicketCategory) error {
	return database.DB.Create(category).Error
}

// FindByID trouve une catégorie par son ID
func (r *ticketCategoryRepository) FindByID(id uint) (*models.TicketCategory, error) {
	var category models.TicketCategory
	err := database.DB.First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// FindBySlug trouve une catégorie par son slug
func (r *ticketCategoryRepository) FindBySlug(slug string) (*models.TicketCategory, error) {
	var category models.TicketCategory
	err := database.DB.Where("slug = ?", slug).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// FindAll récupère toutes les catégories
func (r *ticketCategoryRepository) FindAll() ([]models.TicketCategory, error) {
	var categories []models.TicketCategory
	err := database.DB.Order("display_order ASC, name ASC").Find(&categories).Error
	return categories, err
}

// FindActive récupère uniquement les catégories actives
func (r *ticketCategoryRepository) FindActive() ([]models.TicketCategory, error) {
	var categories []models.TicketCategory
	err := database.DB.Where("is_active = ?", true).Order("display_order ASC, name ASC").Find(&categories).Error
	return categories, err
}

// Update met à jour une catégorie
func (r *ticketCategoryRepository) Update(category *models.TicketCategory) error {
	return database.DB.Save(category).Error
}

// Delete supprime une catégorie (soft delete)
func (r *ticketCategoryRepository) Delete(id uint) error {
	return database.DB.Delete(&models.TicketCategory{}, id).Error
}
