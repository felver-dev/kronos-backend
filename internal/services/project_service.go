package services

import (
	"errors"

	"github.com/mcicare/itsm-backend/internal/models"
	"github.com/mcicare/itsm-backend/internal/repositories"
)

// ProjectService interface pour les opérations sur les projets
type ProjectService interface {
	Create(name, description string, totalBudgetTime *int, createdByID uint) (*models.Project, error)
	GetByID(id uint) (*models.Project, error)
	GetAll() ([]models.Project, error)
	GetByStatus(status string) ([]models.Project, error)
	Update(id uint, name, description string, totalBudgetTime *int, status string, updatedByID uint) (*models.Project, error)
	Delete(id uint) error
	UpdateConsumedTime(projectID uint, consumedTime int) error
}

// projectService implémente ProjectService
type projectService struct {
	projectRepo repositories.ProjectRepository
	userRepo    repositories.UserRepository
}

// NewProjectService crée une nouvelle instance de ProjectService
func NewProjectService(
	projectRepo repositories.ProjectRepository,
	userRepo repositories.UserRepository,
) ProjectService {
	return &projectService{
		projectRepo: projectRepo,
		userRepo:    userRepo,
	}
}

// Create crée un nouveau projet
func (s *projectService) Create(name, description string, totalBudgetTime *int, createdByID uint) (*models.Project, error) {
	// Vérifier que l'utilisateur existe
	_, err := s.userRepo.FindByID(createdByID)
	if err != nil {
		return nil, errors.New("utilisateur créateur introuvable")
	}

	createdByIDPtr := &createdByID
	project := &models.Project{
		Name:            name,
		Description:     description,
		TotalBudgetTime: totalBudgetTime,
		ConsumedTime:    0,
		Status:          "active",
		CreatedByID:     createdByIDPtr,
	}

	if err := s.projectRepo.Create(project); err != nil {
		return nil, errors.New("erreur lors de la création du projet")
	}

	// Récupérer le projet créé
	createdProject, err := s.projectRepo.FindByID(project.ID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération du projet créé")
	}

	return createdProject, nil
}

// GetByID récupère un projet par son ID
func (s *projectService) GetByID(id uint) (*models.Project, error) {
	project, err := s.projectRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("projet introuvable")
	}

	return project, nil
}

// GetAll récupère tous les projets
func (s *projectService) GetAll() ([]models.Project, error) {
	projects, err := s.projectRepo.FindAll()
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des projets")
	}

	return projects, nil
}

// GetByStatus récupère les projets par statut
func (s *projectService) GetByStatus(status string) ([]models.Project, error) {
	projects, err := s.projectRepo.FindByStatus(status)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des projets")
	}

	return projects, nil
}

// Update met à jour un projet
func (s *projectService) Update(id uint, name, description string, totalBudgetTime *int, status string, updatedByID uint) (*models.Project, error) {
	project, err := s.projectRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("projet introuvable")
	}

	// Mettre à jour les champs fournis
	if name != "" {
		project.Name = name
	}
	if description != "" {
		project.Description = description
	}
	if totalBudgetTime != nil {
		project.TotalBudgetTime = totalBudgetTime
	}
	if status != "" {
		project.Status = status
	}

	if err := s.projectRepo.Update(project); err != nil {
		return nil, errors.New("erreur lors de la mise à jour du projet")
	}

	// Récupérer le projet mis à jour
	updatedProject, err := s.projectRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération du projet mis à jour")
	}

	return updatedProject, nil
}

// Delete supprime un projet
func (s *projectService) Delete(id uint) error {
	_, err := s.projectRepo.FindByID(id)
	if err != nil {
		return errors.New("projet introuvable")
	}

	if err := s.projectRepo.Delete(id); err != nil {
		return errors.New("erreur lors de la suppression du projet")
	}

	return nil
}

// UpdateConsumedTime met à jour le temps consommé d'un projet
func (s *projectService) UpdateConsumedTime(projectID uint, consumedTime int) error {
	return s.projectRepo.UpdateConsumedTime(projectID, consumedTime)
}
