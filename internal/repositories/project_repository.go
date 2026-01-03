package repositories

import (
	"github.com/mcicare/itsm-backend/database"
	"github.com/mcicare/itsm-backend/internal/models"
)

// ProjectRepository interface pour les opérations sur les projets
type ProjectRepository interface {
	Create(project *models.Project) error
	FindByID(id uint) (*models.Project, error)
	FindAll() ([]models.Project, error)
	FindByStatus(status string) ([]models.Project, error)
	Update(project *models.Project) error
	Delete(id uint) error
	UpdateConsumedTime(projectID uint, consumedTime int) error
}

// projectRepository implémente ProjectRepository
type projectRepository struct{}

// NewProjectRepository crée une nouvelle instance de ProjectRepository
func NewProjectRepository() ProjectRepository {
	return &projectRepository{}
}

// Create crée un nouveau projet
func (r *projectRepository) Create(project *models.Project) error {
	return database.DB.Create(project).Error
}

// FindByID trouve un projet par son ID avec ses tickets
func (r *projectRepository) FindByID(id uint) (*models.Project, error) {
	var project models.Project
	err := database.DB.Preload("Tickets").Preload("CreatedBy").First(&project, id).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

// FindAll récupère tous les projets
func (r *projectRepository) FindAll() ([]models.Project, error) {
	var projects []models.Project
	err := database.DB.Preload("CreatedBy").Find(&projects).Error
	return projects, err
}

// FindByStatus récupère les projets par statut
func (r *projectRepository) FindByStatus(status string) ([]models.Project, error) {
	var projects []models.Project
	err := database.DB.Preload("CreatedBy").Where("status = ?", status).Find(&projects).Error
	return projects, err
}

// Update met à jour un projet
func (r *projectRepository) Update(project *models.Project) error {
	return database.DB.Save(project).Error
}

// Delete supprime un projet
func (r *projectRepository) Delete(id uint) error {
	return database.DB.Delete(&models.Project{}, id).Error
}

// UpdateConsumedTime met à jour le temps consommé d'un projet
func (r *projectRepository) UpdateConsumedTime(projectID uint, consumedTime int) error {
	return database.DB.Model(&models.Project{}).Where("id = ?", projectID).Update("consumed_time", consumedTime).Error
}
