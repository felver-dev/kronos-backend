package services

import (
	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/repositories"
)

// SettingsService interface pour les opérations sur les paramètres
type SettingsService interface {
	GetAll() (map[string]interface{}, error)
	Update(req dto.UpdateSettingsRequest, updatedByID uint) (map[string]interface{}, error)
}

// settingsService implémente SettingsService
type settingsService struct {
	settingsRepo repositories.SettingsRepository
}

// NewSettingsService crée une nouvelle instance de SettingsService
func NewSettingsService(settingsRepo repositories.SettingsRepository) SettingsService {
	return &settingsService{
		settingsRepo: settingsRepo,
	}
}

// GetAll récupère tous les paramètres
func (s *settingsService) GetAll() (map[string]interface{}, error) {
	settings, err := s.settingsRepo.FindAll()
	if err != nil {
		// Si la table n'existe pas ou est vide, retourner un map vide au lieu d'une erreur
		// Cela permet au frontend de fonctionner même si aucun paramètre n'a été configuré
		return make(map[string]interface{}), nil
	}

	// Organiser les paramètres par catégorie
	result := make(map[string]interface{})
	for _, setting := range settings {
		if setting.Category != "" {
			if result[setting.Category] == nil {
				result[setting.Category] = make(map[string]interface{})
			}
			categoryMap := result[setting.Category].(map[string]interface{})
			categoryMap[setting.Key] = setting.Value
		} else {
			result[setting.Key] = setting.Value
		}
	}

	return result, nil
}

// Update met à jour les paramètres
func (s *settingsService) Update(req dto.UpdateSettingsRequest, updatedByID uint) (map[string]interface{}, error) {
	// TODO: Implémenter la mise à jour des paramètres
	// Pour l'instant, on retourne les paramètres actuels
	return s.GetAll()
}

