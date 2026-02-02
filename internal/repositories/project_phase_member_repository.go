package repositories

import (
	"github.com/mcicare/itsm-backend/database"
	"github.com/mcicare/itsm-backend/internal/models"
)

// ProjectPhaseMemberRepository interface pour les membres d'étape
type ProjectPhaseMemberRepository interface {
	FindByPhaseID(phaseID uint) ([]models.ProjectPhaseMember, error)
	FindByPhaseIDAndUserID(phaseID, userID uint) (*models.ProjectPhaseMember, error)
	Create(m *models.ProjectPhaseMember) error
	Update(m *models.ProjectPhaseMember) error
	Delete(id uint) error
}

type projectPhaseMemberRepository struct{}

// NewProjectPhaseMemberRepository crée une nouvelle instance
func NewProjectPhaseMemberRepository() ProjectPhaseMemberRepository {
	return &projectPhaseMemberRepository{}
}

func (r *projectPhaseMemberRepository) FindByPhaseID(phaseID uint) ([]models.ProjectPhaseMember, error) {
	var list []models.ProjectPhaseMember
	err := database.DB.Where("project_phase_id = ?", phaseID).
		Preload("User").Preload("Function").
		Find(&list).Error
	return list, err
}

func (r *projectPhaseMemberRepository) FindByPhaseIDAndUserID(phaseID, userID uint) (*models.ProjectPhaseMember, error) {
	var m models.ProjectPhaseMember
	err := database.DB.Where("project_phase_id = ? AND user_id = ?", phaseID, userID).First(&m).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *projectPhaseMemberRepository) Create(m *models.ProjectPhaseMember) error {
	return database.DB.Create(m).Error
}

func (r *projectPhaseMemberRepository) Update(m *models.ProjectPhaseMember) error {
	return database.DB.Save(m).Error
}

func (r *projectPhaseMemberRepository) Delete(id uint) error {
	return database.DB.Delete(&models.ProjectPhaseMember{}, id).Error
}
