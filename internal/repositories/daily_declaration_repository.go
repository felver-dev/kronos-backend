package repositories

import (
	"time"

	"github.com/mcicare/itsm-backend/database"
	"github.com/mcicare/itsm-backend/internal/models"
)

// DailyDeclarationRepository interface pour les opérations sur les déclarations journalières
type DailyDeclarationRepository interface {
	Create(declaration *models.DailyDeclaration) error
	FindByID(id uint) (*models.DailyDeclaration, error)
	FindByUserIDAndDate(userID uint, date time.Time) (*models.DailyDeclaration, error)
	FindByUserID(userID uint) ([]models.DailyDeclaration, error)
	FindByDateRange(userID uint, startDate, endDate time.Time) ([]models.DailyDeclaration, error)
	FindValidated() ([]models.DailyDeclaration, error)
	FindPendingValidation() ([]models.DailyDeclaration, error)
	Update(declaration *models.DailyDeclaration) error
	Delete(id uint) error
}

// dailyDeclarationRepository implémente DailyDeclarationRepository
type dailyDeclarationRepository struct{}

// NewDailyDeclarationRepository crée une nouvelle instance de DailyDeclarationRepository
func NewDailyDeclarationRepository() DailyDeclarationRepository {
	return &dailyDeclarationRepository{}
}

// Create crée une nouvelle déclaration journalière
func (r *dailyDeclarationRepository) Create(declaration *models.DailyDeclaration) error {
	return database.DB.Create(declaration).Error
}

// FindByID trouve une déclaration par son ID
func (r *dailyDeclarationRepository) FindByID(id uint) (*models.DailyDeclaration, error) {
	var declaration models.DailyDeclaration
	err := database.DB.Preload("User").Preload("Validator").Preload("Tasks").Preload("Tasks.Ticket").First(&declaration, id).Error
	if err != nil {
		return nil, err
	}
	return &declaration, nil
}

// FindByUserIDAndDate trouve une déclaration par utilisateur et date
func (r *dailyDeclarationRepository) FindByUserIDAndDate(userID uint, date time.Time) (*models.DailyDeclaration, error) {
	var declaration models.DailyDeclaration
	err := database.DB.Preload("User").Preload("Validator").Preload("Tasks").Preload("Tasks.Ticket").
		Where("user_id = ? AND date = ?", userID, date.Format("2006-01-02")).
		First(&declaration).Error
	if err != nil {
		return nil, err
	}
	return &declaration, nil
}

// FindByUserID récupère toutes les déclarations d'un utilisateur
func (r *dailyDeclarationRepository) FindByUserID(userID uint) ([]models.DailyDeclaration, error) {
	var declarations []models.DailyDeclaration
	err := database.DB.Preload("User").Preload("Validator").Preload("Tasks").Where("user_id = ?", userID).Order("date DESC").Find(&declarations).Error
	return declarations, err
}

// FindByDateRange récupère les déclarations d'un utilisateur dans une plage de dates
func (r *dailyDeclarationRepository) FindByDateRange(userID uint, startDate, endDate time.Time) ([]models.DailyDeclaration, error) {
	var declarations []models.DailyDeclaration
	err := database.DB.Preload("User").Preload("Validator").Preload("Tasks").
		Where("user_id = ? AND date >= ? AND date <= ?", userID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02")).
		Order("date DESC").
		Find(&declarations).Error
	return declarations, err
}

// FindValidated récupère les déclarations validées
func (r *dailyDeclarationRepository) FindValidated() ([]models.DailyDeclaration, error) {
	var declarations []models.DailyDeclaration
	err := database.DB.Preload("User").Preload("Validator").Preload("Tasks").Where("validated = ?", true).Order("date DESC").Find(&declarations).Error
	return declarations, err
}

// FindPendingValidation récupère les déclarations en attente de validation
func (r *dailyDeclarationRepository) FindPendingValidation() ([]models.DailyDeclaration, error) {
	var declarations []models.DailyDeclaration
	err := database.DB.Preload("User").Preload("Tasks").Where("validated = ?", false).Order("date DESC").Find(&declarations).Error
	return declarations, err
}

// Update met à jour une déclaration
func (r *dailyDeclarationRepository) Update(declaration *models.DailyDeclaration) error {
	return database.DB.Save(declaration).Error
}

// Delete supprime une déclaration
func (r *dailyDeclarationRepository) Delete(id uint) error {
	return database.DB.Delete(&models.DailyDeclaration{}, id).Error
}
