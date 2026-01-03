package repositories

import (
	"github.com/mcicare/itsm-backend/database"
	"github.com/mcicare/itsm-backend/internal/models"
)

// TicketAssetRepository interface pour les opérations sur les associations ticket-actif
type TicketAssetRepository interface {
	Create(ticketAsset *models.TicketAsset) error
	FindByID(id uint) (*models.TicketAsset, error)
	FindByTicketID(ticketID uint) ([]models.TicketAsset, error)
	FindByAssetID(assetID uint) ([]models.TicketAsset, error)
	Delete(id uint) error
	DeleteByTicketAndAsset(ticketID, assetID uint) error
}

// ticketAssetRepository implémente TicketAssetRepository
type ticketAssetRepository struct{}

// NewTicketAssetRepository crée une nouvelle instance de TicketAssetRepository
func NewTicketAssetRepository() TicketAssetRepository {
	return &ticketAssetRepository{}
}

// Create crée une nouvelle association ticket-actif
func (r *ticketAssetRepository) Create(ticketAsset *models.TicketAsset) error {
	return database.DB.Create(ticketAsset).Error
}

// FindByID trouve une association par son ID
func (r *ticketAssetRepository) FindByID(id uint) (*models.TicketAsset, error) {
	var ticketAsset models.TicketAsset
	err := database.DB.Preload("Ticket").Preload("Asset").First(&ticketAsset, id).Error
	if err != nil {
		return nil, err
	}
	return &ticketAsset, nil
}

// FindByTicketID récupère toutes les associations d'un ticket
func (r *ticketAssetRepository) FindByTicketID(ticketID uint) ([]models.TicketAsset, error) {
	var ticketAssets []models.TicketAsset
	err := database.DB.Preload("Asset").Where("ticket_id = ?", ticketID).Find(&ticketAssets).Error
	return ticketAssets, err
}

// FindByAssetID récupère toutes les associations d'un actif
func (r *ticketAssetRepository) FindByAssetID(assetID uint) ([]models.TicketAsset, error) {
	var ticketAssets []models.TicketAsset
	err := database.DB.Preload("Ticket").Where("asset_id = ?", assetID).Find(&ticketAssets).Error
	return ticketAssets, err
}

// Delete supprime une association
func (r *ticketAssetRepository) Delete(id uint) error {
	return database.DB.Delete(&models.TicketAsset{}, id).Error
}

// DeleteByTicketAndAsset supprime une association par ticket et actif
func (r *ticketAssetRepository) DeleteByTicketAndAsset(ticketID, assetID uint) error {
	return database.DB.Where("ticket_id = ? AND asset_id = ?", ticketID, assetID).Delete(&models.TicketAsset{}).Error
}

