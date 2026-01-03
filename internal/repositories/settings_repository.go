package repositories

import (
	"github.com/mcicare/itsm-backend/database"
	"github.com/mcicare/itsm-backend/internal/models"
)

// SettingsRepository interface pour les opérations sur les paramètres système
type SettingsRepository interface {
	Create(setting *models.Setting) error
	FindByID(id uint) (*models.Setting, error)
	FindByKey(key string) (*models.Setting, error)
	FindAll() ([]models.Setting, error)
	FindByCategory(category string) ([]models.Setting, error)
	FindPublic() ([]models.Setting, error)
	Update(setting *models.Setting) error
	Delete(id uint) error
	GetValue(key string) (string, error)
	SetValue(key, value string) error
}

// settingsRepository implémente SettingsRepository
type settingsRepository struct{}

// NewSettingsRepository crée une nouvelle instance de SettingsRepository
func NewSettingsRepository() SettingsRepository {
	return &settingsRepository{}
}

// Create crée un nouveau paramètre
func (r *settingsRepository) Create(setting *models.Setting) error {
	return database.DB.Create(setting).Error
}

// FindByID trouve un paramètre par son ID
func (r *settingsRepository) FindByID(id uint) (*models.Setting, error) {
	var setting models.Setting
	err := database.DB.Preload("UpdatedBy").First(&setting, id).Error
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

// FindByKey trouve un paramètre par sa clé
func (r *settingsRepository) FindByKey(key string) (*models.Setting, error) {
	var setting models.Setting
	err := database.DB.Preload("UpdatedBy").Where("key = ?", key).First(&setting).Error
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

// FindAll récupère tous les paramètres
func (r *settingsRepository) FindAll() ([]models.Setting, error) {
	var settings []models.Setting
	err := database.DB.Preload("UpdatedBy").Order("category ASC, key ASC").Find(&settings).Error
	return settings, err
}

// FindByCategory récupère les paramètres d'une catégorie
func (r *settingsRepository) FindByCategory(category string) ([]models.Setting, error) {
	var settings []models.Setting
	err := database.DB.Preload("UpdatedBy").Where("category = ?", category).Order("key ASC").Find(&settings).Error
	return settings, err
}

// FindPublic récupère les paramètres publics (accessibles sans authentification)
func (r *settingsRepository) FindPublic() ([]models.Setting, error) {
	var settings []models.Setting
	err := database.DB.Preload("UpdatedBy").Where("is_public = ?", true).Order("category ASC, key ASC").Find(&settings).Error
	return settings, err
}

// Update met à jour un paramètre
func (r *settingsRepository) Update(setting *models.Setting) error {
	return database.DB.Save(setting).Error
}

// Delete supprime un paramètre
func (r *settingsRepository) Delete(id uint) error {
	return database.DB.Delete(&models.Setting{}, id).Error
}

// GetValue récupère la valeur d'un paramètre par sa clé (méthode utilitaire)
func (r *settingsRepository) GetValue(key string) (string, error) {
	setting, err := r.FindByKey(key)
	if err != nil {
		return "", err
	}
	return setting.Value, nil
}

// SetValue met à jour la valeur d'un paramètre par sa clé (méthode utilitaire)
func (r *settingsRepository) SetValue(key, value string) error {
	setting, err := r.FindByKey(key)
	if err != nil {
		// Si le paramètre n'existe pas, on le crée
		setting = &models.Setting{
			Key:   key,
			Value: value,
			Type:  "string",
		}
		return r.Create(setting)
	}
	// Sinon, on met à jour
	setting.Value = value
	return r.Update(setting)
}
