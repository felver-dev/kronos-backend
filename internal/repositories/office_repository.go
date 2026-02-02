package repositories

import (
	"github.com/mcicare/itsm-backend/database"
	"github.com/mcicare/itsm-backend/internal/models"
)

// OfficeRepository interface pour les opérations sur les sièges
type OfficeRepository interface {
	Create(office *models.Office) error
	FindByID(id uint) (*models.Office, error)
	FindByCode(code string) (*models.Office, error)
	FindAll() ([]models.Office, error)
	FindActive() ([]models.Office, error)
	FindByFilialeID(filialeID uint) ([]models.Office, error)
	FindByCountry(country string) ([]models.Office, error)
	FindByCity(city string) ([]models.Office, error)
	Update(office *models.Office) error
	Delete(id uint) error
}

// officeRepository implémente OfficeRepository
type officeRepository struct{}

// NewOfficeRepository crée une nouvelle instance de OfficeRepository
func NewOfficeRepository() OfficeRepository {
	return &officeRepository{}
}

// Create crée un nouveau siège
func (r *officeRepository) Create(office *models.Office) error {
	return database.DB.Create(office).Error
}

// FindByID trouve un siège par son ID
func (r *officeRepository) FindByID(id uint) (*models.Office, error) {
	var office models.Office
	err := database.DB.Preload("Filiale").First(&office, id).Error
	if err != nil {
		return nil, err
	}
	return &office, nil
}

// FindByCode trouve un siège par son code
func (r *officeRepository) FindByCode(code string) (*models.Office, error) {
	var office models.Office
	err := database.DB.Preload("Filiale").Where("code = ?", code).First(&office).Error
	if err != nil {
		return nil, err
	}
	return &office, nil
}

// FindAll récupère tous les sièges
func (r *officeRepository) FindAll() ([]models.Office, error) {
	var offices []models.Office
	err := database.DB.Preload("Filiale").Order("country ASC, city ASC, name ASC").Find(&offices).Error
	return offices, err
}

// FindActive récupère tous les sièges actifs
func (r *officeRepository) FindActive() ([]models.Office, error) {
	var offices []models.Office
	err := database.DB.Preload("Filiale").Where("is_active = ?", true).Order("country ASC, city ASC, name ASC").Find(&offices).Error
	return offices, err
}

// FindByFilialeID récupère les sièges d'une filiale
func (r *officeRepository) FindByFilialeID(filialeID uint) ([]models.Office, error) {
	var offices []models.Office
	err := database.DB.Preload("Filiale").Where("filiale_id = ? AND is_active = ?", filialeID, true).Order("country ASC, city ASC, name ASC").Find(&offices).Error
	return offices, err
}

// FindByCountry récupère les sièges d'un pays
func (r *officeRepository) FindByCountry(country string) ([]models.Office, error) {
	var offices []models.Office
	err := database.DB.Where("country = ?", country).Order("city ASC, name ASC").Find(&offices).Error
	return offices, err
}

// FindByCity récupère les sièges d'une ville
func (r *officeRepository) FindByCity(city string) ([]models.Office, error) {
	var offices []models.Office
	err := database.DB.Preload("Filiale").Where("city = ?", city).Order("name ASC").Find(&offices).Error
	return offices, err
}

// Update met à jour un siège
func (r *officeRepository) Update(office *models.Office) error {
	return database.DB.Save(office).Error
}

// Delete supprime un siège (soft delete)
func (r *officeRepository) Delete(id uint) error {
	return database.DB.Delete(&models.Office{}, id).Error
}
