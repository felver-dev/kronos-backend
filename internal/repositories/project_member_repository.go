package repositories

import (
	"github.com/mcicare/itsm-backend/database"
	"github.com/mcicare/itsm-backend/internal/models"
)

// ProjectMemberRepository interface pour les membres du projet
type ProjectMemberRepository interface {
	FindByProjectID(projectID uint) ([]models.ProjectMember, error)
	FindByProjectIDAndUserID(projectID, userID uint) (*models.ProjectMember, error)
	Create(m *models.ProjectMember) error
	Update(m *models.ProjectMember) error
	Delete(id uint) error
	SetProjectManager(projectID, userID uint) error
	SetLead(projectID, userID uint) error
	ReplaceMemberFunctions(projectMemberID uint, functionIDs []uint) error
}

type projectMemberRepository struct{}

// NewProjectMemberRepository crée une nouvelle instance
func NewProjectMemberRepository() ProjectMemberRepository {
	return &projectMemberRepository{}
}

func (r *projectMemberRepository) FindByProjectID(projectID uint) ([]models.ProjectMember, error) {
	var list []models.ProjectMember
	err := database.DB.Where("project_id = ?", projectID).
		Preload("User").Preload("Function").Preload("Functions").
		Find(&list).Error
	return list, err
}

func (r *projectMemberRepository) FindByProjectIDAndUserID(projectID, userID uint) (*models.ProjectMember, error) {
	var m models.ProjectMember
	err := database.DB.Where("project_id = ? AND user_id = ?", projectID, userID).First(&m).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *projectMemberRepository) Create(m *models.ProjectMember) error {
	return database.DB.Create(m).Error
}

func (r *projectMemberRepository) Update(m *models.ProjectMember) error {
	return database.DB.Save(m).Error
}

func (r *projectMemberRepository) Delete(id uint) error {
	return database.DB.Delete(&models.ProjectMember{}, id).Error
}

func (r *projectMemberRepository) SetProjectManager(projectID, userID uint) error {
	if err := database.DB.Model(&models.ProjectMember{}).Where("project_id = ?", projectID).Update("is_project_manager", false).Error; err != nil {
		return err
	}
	if userID != 0 {
		if err := database.DB.Model(&models.ProjectMember{}).Where("project_id = ? AND user_id = ?", projectID, userID).Update("is_project_manager", true).Error; err != nil {
			return err
		}
	}
	var ptr *uint
	if userID != 0 {
		ptr = &userID
	}
	return database.DB.Model(&models.Project{}).Where("id = ?", projectID).Update("project_manager_id", ptr).Error
}

func (r *projectMemberRepository) SetLead(projectID, userID uint) error {
	if err := database.DB.Model(&models.ProjectMember{}).Where("project_id = ?", projectID).Update("is_lead", false).Error; err != nil {
		return err
	}
	if userID != 0 {
		if err := database.DB.Model(&models.ProjectMember{}).Where("project_id = ? AND user_id = ?", projectID, userID).Update("is_lead", true).Error; err != nil {
			return err
		}
	}
	var ptr *uint
	if userID != 0 {
		ptr = &userID
	}
	return database.DB.Model(&models.Project{}).Where("id = ?", projectID).Update("lead_id", ptr).Error
}

func (r *projectMemberRepository) ReplaceMemberFunctions(projectMemberID uint, functionIDs []uint) error {
	if err := database.DB.Where("project_member_id = ?", projectMemberID).Delete(&models.ProjectMemberFunction{}).Error; err != nil {
		return err
	}
	for _, fid := range functionIDs {
		if err := database.DB.Create(&models.ProjectMemberFunction{ProjectMemberID: projectMemberID, ProjectFunctionID: fid}).Error; err != nil {
			return err
		}
	}
	return nil
}
