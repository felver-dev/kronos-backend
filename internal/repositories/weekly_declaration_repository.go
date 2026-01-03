package repositories

import (
	"github.com/mcicare/itsm-backend/database"
	"github.com/mcicare/itsm-backend/internal/models"
)

// WeeklyDeclarationRepository interface pour les opérations sur les déclarations hebdomadaires
type WeeklyDeclarationRepository interface {
	Create(declaration *models.WeeklyDeclaration) error
	FindByID(id uint) (*models.WeeklyDeclaration, error)
	FindByUserIDAndWeek(userID uint, week string) (*models.WeeklyDeclaration, error)
	FindByUserID(userID uint) ([]models.WeeklyDeclaration, error)
	FindValidated() ([]models.WeeklyDeclaration, error)
	FindPendingValidation() ([]models.WeeklyDeclaration, error)
	Update(declaration *models.WeeklyDeclaration) error
	Delete(id uint) error
}

// weeklyDeclarationRepository implémente WeeklyDeclarationRepository
type weeklyDeclarationRepository struct{}

// NewWeeklyDeclarationRepository crée une nouvelle instance de WeeklyDeclarationRepository
func NewWeeklyDeclarationRepository() WeeklyDeclarationRepository {
	return &weeklyDeclarationRepository{}
}

// Create crée une nouvelle déclaration hebdomadaire
func (r *weeklyDeclarationRepository) Create(declaration *models.WeeklyDeclaration) error {
	return database.DB.Create(declaration).Error
}

// FindByID trouve une déclaration par son ID
func (r *weeklyDeclarationRepository) FindByID(id uint) (*models.WeeklyDeclaration, error) {
	var declaration models.WeeklyDeclaration
	err := database.DB.Preload("User").Preload("Validator").Preload("Tasks").Preload("Tasks.Ticket").First(&declaration, id).Error
	if err != nil {
		return nil, err
	}
	return &declaration, nil
}

// FindByUserIDAndWeek trouve une déclaration par utilisateur et semaine
func (r *weeklyDeclarationRepository) FindByUserIDAndWeek(userID uint, week string) (*models.WeeklyDeclaration, error) {
	var declaration models.WeeklyDeclaration
	err := database.DB.Preload("User").Preload("Validator").Preload("Tasks").Preload("Tasks.Ticket").
		Where("user_id = ? AND week = ?", userID, week).
		First(&declaration).Error
	if err != nil {
		return nil, err
	}
	return &declaration, nil
}

// FindByUserID récupère toutes les déclarations d'un utilisateur
func (r *weeklyDeclarationRepository) FindByUserID(userID uint) ([]models.WeeklyDeclaration, error) {
	var declarations []models.WeeklyDeclaration
	err := database.DB.Preload("User").Preload("Validator").Preload("Tasks").Where("user_id = ?", userID).Order("week DESC").Find(&declarations).Error
	return declarations, err
}

// FindValidated récupère les déclarations validées
func (r *weeklyDeclarationRepository) FindValidated() ([]models.WeeklyDeclaration, error) {
	var declarations []models.WeeklyDeclaration
	err := database.DB.Preload("User").Preload("Validator").Preload("Tasks").Where("validated = ?", true).Order("week DESC").Find(&declarations).Error
	return declarations, err
}

// FindPendingValidation récupère les déclarations en attente de validation
func (r *weeklyDeclarationRepository) FindPendingValidation() ([]models.WeeklyDeclaration, error) {
	var declarations []models.WeeklyDeclaration
	err := database.DB.Preload("User").Preload("Tasks").Where("validated = ?", false).Order("week DESC").Find(&declarations).Error
	return declarations, err
}

// Update met à jour une déclaration
func (r *weeklyDeclarationRepository) Update(declaration *models.WeeklyDeclaration) error {
	return database.DB.Save(declaration).Error
}

// Delete supprime une déclaration (soft delete)
func (r *weeklyDeclarationRepository) Delete(id uint) error {
	return database.DB.Delete(&models.WeeklyDeclaration{}, id).Error
}

