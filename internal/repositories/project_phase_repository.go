package repositories

import (
	"github.com/mcicare/itsm-backend/database"
	"github.com/mcicare/itsm-backend/internal/models"
)

// ProjectPhaseRepository interface pour les étapes de projet
type ProjectPhaseRepository interface {
	FindByProjectID(projectID uint) ([]models.ProjectPhase, error)
	FindByID(id uint) (*models.ProjectPhase, error)
	Create(p *models.ProjectPhase) error
	Update(p *models.ProjectPhase) error
	Delete(id uint) error
	Reorder(projectID uint, order []uint) error
}

type projectPhaseRepository struct{}

// NewProjectPhaseRepository crée une nouvelle instance
func NewProjectPhaseRepository() ProjectPhaseRepository {
	return &projectPhaseRepository{}
}

func (r *projectPhaseRepository) FindByProjectID(projectID uint) ([]models.ProjectPhase, error) {
	var list []models.ProjectPhase
	err := database.DB.Where("project_id = ?", projectID).Order("display_order ASC, id ASC").Find(&list).Error
	return list, err
}

func (r *projectPhaseRepository) FindByID(id uint) (*models.ProjectPhase, error) {
	var p models.ProjectPhase
	err := database.DB.First(&p, id).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *projectPhaseRepository) Create(p *models.ProjectPhase) error {
	return database.DB.Create(p).Error
}

func (r *projectPhaseRepository) Update(p *models.ProjectPhase) error {
	return database.DB.Save(p).Error
}

func (r *projectPhaseRepository) Delete(id uint) error {
	return database.DB.Delete(&models.ProjectPhase{}, id).Error
}

func (r *projectPhaseRepository) Reorder(projectID uint, order []uint) error {
	for i, id := range order {
		if err := database.DB.Model(&models.ProjectPhase{}).Where("id = ? AND project_id = ?", id, projectID).Update("display_order", i).Error; err != nil {
			return err
		}
	}
	return nil
}
