package repositories

import (
	"strings"

	"github.com/mcicare/itsm-backend/database"
	"github.com/mcicare/itsm-backend/internal/models"
)

// ProjectBudgetExtensionRepository interface pour les extensions de budget
type ProjectBudgetExtensionRepository interface {
	Create(ext *models.ProjectBudgetExtension) error
	FindByID(id uint) (*models.ProjectBudgetExtension, error)
	FindByProjectID(projectID uint) ([]models.ProjectBudgetExtension, error)
	Update(ext *models.ProjectBudgetExtension) error
	Delete(id uint) error
}

type projectBudgetExtensionRepository struct{}

// NewProjectBudgetExtensionRepository crée une nouvelle instance
func NewProjectBudgetExtensionRepository() ProjectBudgetExtensionRepository {
	return &projectBudgetExtensionRepository{}
}

// Create crée une extension de budget
func (r *projectBudgetExtensionRepository) Create(ext *models.ProjectBudgetExtension) error {
	return database.DB.Create(ext).Error
}

// FindByID retourne une extension par son ID
func (r *projectBudgetExtensionRepository) FindByID(id uint) (*models.ProjectBudgetExtension, error) {
	var ext models.ProjectBudgetExtension
	if err := database.DB.First(&ext, id).Error; err != nil {
		return nil, err
	}
	return &ext, nil
}

// FindByProjectID retourne les extensions d'un projet, par date décroissante
func (r *projectBudgetExtensionRepository) FindByProjectID(projectID uint) ([]models.ProjectBudgetExtension, error) {
	var list []models.ProjectBudgetExtension
	err := database.DB.Where("project_id = ?", projectID).
		Order("created_at DESC").
		Find(&list).Error
	if err != nil {
		// Si la table n'existe pas encore (migration non exécutée), retourner une liste vide au lieu d'échouer
		if strings.Contains(err.Error(), "doesn't exist") || strings.Contains(err.Error(), "1146") {
			return []models.ProjectBudgetExtension{}, nil
		}
		return nil, err
	}
	// Préchargement manuel de CreatedBy (ID, Username)
	for i := range list {
		if list[i].CreatedByID != nil {
			var u struct {
				ID       uint   `json:"id"`
				Username string `json:"username"`
			}
			if database.DB.Table("users").Select("id", "username").Where("id = ?", *list[i].CreatedByID).First(&u).Error == nil {
				list[i].CreatedBy = &models.User{ID: u.ID, Username: u.Username}
			}
		}
	}
	return list, nil
}

// Update met à jour une extension de budget
func (r *projectBudgetExtensionRepository) Update(ext *models.ProjectBudgetExtension) error {
	return database.DB.Save(ext).Error
}

// Delete supprime une extension de budget
func (r *projectBudgetExtensionRepository) Delete(id uint) error {
	return database.DB.Delete(&models.ProjectBudgetExtension{}, id).Error
}
