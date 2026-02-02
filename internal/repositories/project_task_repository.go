package repositories

import (
	"fmt"
	"regexp"
	"time"

	"github.com/mcicare/itsm-backend/database"
	"github.com/mcicare/itsm-backend/internal/models"
)

// ProjectTaskRepository interface pour les tâches de projet
type ProjectTaskRepository interface {
	FindByProjectID(projectID uint) ([]models.ProjectTask, error)
	FindByPhaseID(phaseID uint) ([]models.ProjectTask, error)
	FindByID(id uint) (*models.ProjectTask, error)
	Create(t *models.ProjectTask) error
	Update(t *models.ProjectTask) error
	UpdateActualTime(taskID uint, minutes int) error
	Delete(id uint) error
	GenerateCode(projectID uint) (string, error)
	ReplaceAssignees(taskID uint, userIDs []uint) error
}

type projectTaskRepository struct{}

// NewProjectTaskRepository crée une nouvelle instance
func NewProjectTaskRepository() ProjectTaskRepository {
	return &projectTaskRepository{}
}

func (r *projectTaskRepository) FindByProjectID(projectID uint) ([]models.ProjectTask, error) {
	var list []models.ProjectTask
	err := database.DB.Where("project_id = ?", projectID).
		Preload("ProjectPhase").Preload("AssignedTo").Preload("CreatedBy").Preload("Assignees").Preload("Assignees.User").
		Order("project_phase_id ASC, display_order ASC, id ASC").
		Find(&list).Error
	return list, err
}

func (r *projectTaskRepository) FindByPhaseID(phaseID uint) ([]models.ProjectTask, error) {
	var list []models.ProjectTask
	err := database.DB.Where("project_phase_id = ?", phaseID).
		Preload("AssignedTo").Preload("CreatedBy").
		Order("display_order ASC, id ASC").
		Find(&list).Error
	return list, err
}

func (r *projectTaskRepository) FindByID(id uint) (*models.ProjectTask, error) {
	var t models.ProjectTask
	err := database.DB.Preload("ProjectPhase").Preload("AssignedTo").Preload("CreatedBy").Preload("Assignees").Preload("Assignees.User").
		First(&t, id).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *projectTaskRepository) Create(t *models.ProjectTask) error {
	return database.DB.Create(t).Error
}

func (r *projectTaskRepository) Update(t *models.ProjectTask) error {
	return database.DB.Save(t).Error
}

func (r *projectTaskRepository) UpdateActualTime(taskID uint, minutes int) error {
	return database.DB.Model(&models.ProjectTask{}).Where("id = ?", taskID).Update("actual_time", minutes).Error
}

func (r *projectTaskRepository) Delete(id uint) error {
	_ = database.DB.Where("project_task_id = ?", id).Delete(&models.ProjectTaskAssignee{}).Error
	return database.DB.Delete(&models.ProjectTask{}, id).Error
}

func (r *projectTaskRepository) ReplaceAssignees(taskID uint, userIDs []uint) error {
	if err := database.DB.Where("project_task_id = ?", taskID).Delete(&models.ProjectTaskAssignee{}).Error; err != nil {
		return err
	}
	if len(userIDs) == 0 {
		return nil
	}
	list := make([]models.ProjectTaskAssignee, 0, len(userIDs))
	for _, uid := range userIDs {
		list = append(list, models.ProjectTaskAssignee{ProjectTaskID: taskID, UserID: uid})
	}
	return database.DB.Create(&list).Error
}

var codeSuffixRE = regexp.MustCompile(`^TAP-\d{4}-(\d+)$`)

// GenerateCode génère un code TAP-YYYY-NNNN pour une nouvelle tâche du projet
func (r *projectTaskRepository) GenerateCode(projectID uint) (string, error) {
	year := time.Now().Year()
	prefix := fmt.Sprintf("TAP-%d-", year)

	var codes []string
	if err := database.DB.Model(&models.ProjectTask{}).Where("project_id = ? AND code LIKE ?", projectID, prefix+"%").Pluck("code", &codes).Error; err != nil {
		return "", err
	}
	maxN := 0
	for _, c := range codes {
		if m := codeSuffixRE.FindStringSubmatch(c); len(m) == 2 {
			var n int
			if _, err := fmt.Sscanf(m[1], "%d", &n); err == nil && n > maxN {
				maxN = n
			}
		}
	}
	return prefix + fmt.Sprintf("%04d", maxN+1), nil
}
