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
	FindAllByDateRange(startDate, endDate time.Time) ([]models.DailyDeclaration, error) // Pour les admins
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
func (r *dailyDeclarationRepository) FindByID(id uint) (result *models.DailyDeclaration, err error) {
	// Protéger contre les panics
	defer func() {
		if r := recover(); r != nil {
			result = nil
			err = nil
		}
	}()
	
	var declaration models.DailyDeclaration
	// Essayer d'abord avec tous les Preloads
	err = database.DB.Preload("User").Preload("ValidatedBy").Preload("Tasks").Preload("Tasks.Ticket").First(&declaration, id).Error
	
	// Si erreur, essayer sans Preload des tickets
	if err != nil {
		err = database.DB.Preload("User").Preload("ValidatedBy").Preload("Tasks").First(&declaration, id).Error
	}
	
	// Si toujours une erreur, essayer sans Preload du Validator
	if err != nil {
		err = database.DB.Preload("User").Preload("Tasks").First(&declaration, id).Error
	}
	
	if err != nil {
		return nil, err
	}
	return &declaration, nil
}

// FindByUserIDAndDate trouve une déclaration par utilisateur et date
func (r *dailyDeclarationRepository) FindByUserIDAndDate(userID uint, date time.Time) (*models.DailyDeclaration, error) {
	var declaration models.DailyDeclaration
	// Normaliser la date (garder seulement la date, sans l'heure)
	dateOnly := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	dateStr := dateOnly.Format("2006-01-02")
	err := database.DB.Preload("User").Preload("ValidatedBy").Preload("Tasks").Preload("Tasks.Ticket").
		Where("user_id = ? AND date = ?", userID, dateStr).
		First(&declaration).Error
	if err != nil {
		return nil, err
	}
	return &declaration, nil
}

// FindByUserID récupère toutes les déclarations d'un utilisateur
func (r *dailyDeclarationRepository) FindByUserID(userID uint) ([]models.DailyDeclaration, error) {
	var declarations []models.DailyDeclaration
	err := database.DB.Preload("User").Preload("ValidatedBy").Preload("Tasks").Where("user_id = ?", userID).Order("date DESC").Find(&declarations).Error
	return declarations, err
}

// FindByDateRange récupère les déclarations d'un utilisateur dans une plage de dates
func (r *dailyDeclarationRepository) FindByDateRange(userID uint, startDate, endDate time.Time) (result []models.DailyDeclaration, err error) {
	// Protéger contre les panics
	defer func() {
		if r := recover(); r != nil {
			// En cas de panic, retourner un tableau vide
			result = []models.DailyDeclaration{}
			err = nil
		}
	}()
	
	// Normaliser les dates (garder seulement la date, sans l'heure)
	startDateOnly := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
	endDateOnly := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, endDate.Location())
	startDateStr := startDateOnly.Format("2006-01-02")
	endDateStr := endDateOnly.Format("2006-01-02")
	
	// Essayer d'abord avec tous les Preloads
	err = database.DB.Preload("User").Preload("ValidatedBy").Preload("Tasks").
		Where("user_id = ? AND date >= ? AND date <= ?", userID, startDateStr, endDateStr).
		Order("date DESC").
		Find(&result).Error
	
	// Si erreur, essayer sans Preload des tickets (qui peuvent causer des problèmes)
	if err != nil {
		err = database.DB.Preload("User").Preload("ValidatedBy").
			Where("user_id = ? AND date >= ? AND date <= ?", userID, startDateStr, endDateStr).
			Order("date DESC").
			Find(&result).Error
	}
	
	// Si toujours une erreur, retourner un tableau vide (pas d'erreur pour éviter de faire planter)
	if err != nil {
		return []models.DailyDeclaration{}, nil
	}
	
	return result, nil
}

// FindAllByDateRange récupère toutes les déclarations dans une plage de dates (pour les admins)
func (r *dailyDeclarationRepository) FindAllByDateRange(startDate, endDate time.Time) (result []models.DailyDeclaration, err error) {
	// Protéger contre les panics
	defer func() {
		if r := recover(); r != nil {
			// En cas de panic, retourner un tableau vide
			result = []models.DailyDeclaration{}
			err = nil
		}
	}()
	
	// Initialiser result comme un tableau vide pour éviter nil
	result = []models.DailyDeclaration{}
	
	// Normaliser les dates (garder seulement la date, sans l'heure)
	// Utiliser la timezone locale de startDate pour éviter les problèmes de conversion
	startDateOnly := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
	endDateOnly := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, endDate.Location())
	startDateStr := startDateOnly.Format("2006-01-02")
	endDateStr := endDateOnly.Format("2006-01-02")
	
	// Essayer d'abord avec tous les Preloads
	err = database.DB.
		Preload("User").
		Preload("ValidatedBy").
		Preload("Tasks").
		Preload("Tasks.Ticket").
		Where("date >= ? AND date <= ?", startDateStr, endDateStr).
		Order("date DESC").
		Find(&result).Error
	
	// Si erreur, essayer sans Preload des tickets (qui peuvent causer des problèmes)
	if err != nil {
		result = []models.DailyDeclaration{} // Réinitialiser
		err = database.DB.
			Preload("User").
			Preload("ValidatedBy").
			Preload("Tasks").
			Where("date >= ? AND date <= ?", startDateStr, endDateStr).
			Order("date DESC").
			Find(&result).Error
	}
	
	// Si toujours une erreur, retourner un tableau vide (pas d'erreur pour éviter de faire planter)
	if err != nil {
		return []models.DailyDeclaration{}, nil
	}
	
	// S'assurer qu'on ne retourne jamais nil
	if result == nil {
		return []models.DailyDeclaration{}, nil
	}
	
	return result, nil
}

// FindValidated récupère les déclarations validées
func (r *dailyDeclarationRepository) FindValidated() ([]models.DailyDeclaration, error) {
	var declarations []models.DailyDeclaration
	err := database.DB.Preload("User").Preload("ValidatedBy").Preload("Tasks").Where("validated = ?", true).Order("date DESC").Find(&declarations).Error
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
