package repositories

import (
	"github.com/mcicare/itsm-backend/database"
	"github.com/mcicare/itsm-backend/internal/models"
)

// ProjectFunctionRepository interface pour les fonctions projet
type ProjectFunctionRepository interface {
	FindByProjectID(projectID uint) ([]models.ProjectFunction, error)
	FindByID(id uint) (*models.ProjectFunction, error)
	Create(f *models.ProjectFunction) error
	Update(f *models.ProjectFunction) error
	Delete(id uint) error
}

type projectFunctionRepository struct{}

// NewProjectFunctionRepository crée une nouvelle instance
func NewProjectFunctionRepository() ProjectFunctionRepository {
	return &projectFunctionRepository{}
}

func (r *projectFunctionRepository) FindByProjectID(projectID uint) ([]models.ProjectFunction, error) {
	var list []models.ProjectFunction
	err := database.DB.Where("(project_id = ? OR project_id IS NULL)", projectID).
		Order("display_order ASC, name ASC").Find(&list).Error
	return list, err
}

func (r *projectFunctionRepository) FindByID(id uint) (*models.ProjectFunction, error) {
	var f models.ProjectFunction
	err := database.DB.First(&f, id).Error
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *projectFunctionRepository) Create(f *models.ProjectFunction) error {
	return database.DB.Create(f).Error
}

func (r *projectFunctionRepository) Update(f *models.ProjectFunction) error {
	return database.DB.Save(f).Error
}

func (r *projectFunctionRepository) Delete(id uint) error {
	return database.DB.Delete(&models.ProjectFunction{}, id).Error
}
