package repositories

import (
	"github.com/mcicare/itsm-backend/database"
	"github.com/mcicare/itsm-backend/internal/models"
	"github.com/mcicare/itsm-backend/internal/scope"
	"gorm.io/gorm"
)

// ProjectRepository interface pour les opérations sur les projets
type ProjectRepository interface {
	Create(project *models.Project) error
	FindByID(id uint) (*models.Project, error)
	FindAll(scope interface{}) ([]models.Project, error) // scope peut être *scope.QueryScope ou nil
	FindByStatus(scope interface{}, status string) ([]models.Project, error)
	Update(project *models.Project) error
	UpdateStatus(projectID uint, status string) error
	Delete(id uint) error
	UpdateConsumedTime(projectID uint, consumedTime int) error
	IncrementTotalBudgetTime(projectID uint, additionalMinutes int) error
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

// FindByID trouve un projet par son ID (sans Preload Tickets/CreatedBy pour alléger la fiche détail)
func (r *projectRepository) FindByID(id uint) (*models.Project, error) {
	var project models.Project
	err := database.DB.First(&project, id).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

// FindAll récupère tous les projets
// Le scope est utilisé pour filtrer automatiquement selon les permissions de l'utilisateur
func (r *projectRepository) FindAll(scopeParam interface{}) ([]models.Project, error) {
	var projects []models.Project
	
	// Construire la requête de base
	query := database.DB.Model(&models.Project{}).
		Preload("CreatedBy")
	
	// Appliquer le scope si fourni
	if scopeParam != nil {
		if queryScope, ok := scopeParam.(*scope.QueryScope); ok {
			query = scope.ApplyProjectScope(query, queryScope)
		}
	}
	
	err := query.Find(&projects).Error
	return projects, err
}

// FindByStatus récupère les projets par statut
// Le scope est utilisé pour filtrer automatiquement selon les permissions de l'utilisateur
func (r *projectRepository) FindByStatus(scopeParam interface{}, status string) ([]models.Project, error) {
	var projects []models.Project
	
	// Construire la requête de base
	query := database.DB.Model(&models.Project{}).
		Preload("CreatedBy").
		Where("projects.status = ?", status)
	
	// Appliquer le scope si fourni
	if scopeParam != nil {
		if queryScope, ok := scopeParam.(*scope.QueryScope); ok {
			query = scope.ApplyProjectScope(query, queryScope)
		}
	}
	
	err := query.Find(&projects).Error
	return projects, err
}

// Update met à jour un projet
func (r *projectRepository) Update(project *models.Project) error {
	updates := map[string]interface{}{
		"name":              project.Name,
		"description":       project.Description,
		"total_budget_time": project.TotalBudgetTime,
		"status":            project.Status,
		"start_date":        project.StartDate,
		"end_date":          project.EndDate,
	}
	return database.DB.Model(project).Updates(updates).Error
}

// UpdateStatus met à jour uniquement le statut d'un projet
func (r *projectRepository) UpdateStatus(projectID uint, status string) error {
	return database.DB.Model(&models.Project{}).Where("id = ?", projectID).Update("status", status).Error
}

// Delete supprime un projet
func (r *projectRepository) Delete(id uint) error {
	return database.DB.Delete(&models.Project{}, id).Error
}

// UpdateConsumedTime met à jour le temps consommé d'un projet
func (r *projectRepository) UpdateConsumedTime(projectID uint, consumedTime int) error {
	return database.DB.Model(&models.Project{}).Where("id = ?", projectID).Update("consumed_time", consumedTime).Error
}

// IncrementTotalBudgetTime ajoute des minutes au budget total du projet (pour les extensions)
func (r *projectRepository) IncrementTotalBudgetTime(projectID uint, additionalMinutes int) error {
	return database.DB.Model(&models.Project{}).Where("id = ?", projectID).
		Update("total_budget_time", gorm.Expr("COALESCE(total_budget_time, 0) + ?", additionalMinutes)).Error
}
