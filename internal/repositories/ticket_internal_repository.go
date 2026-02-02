package repositories

import (
	"fmt"

	"github.com/mcicare/itsm-backend/database"
	"github.com/mcicare/itsm-backend/internal/models"
	"github.com/mcicare/itsm-backend/internal/scope"
	"gorm.io/gorm"
)

// TicketInternalRepository interface pour les opérations sur les tickets internes
type TicketInternalRepository interface {
	Create(t *models.TicketInternal) error
	FindByID(id uint) (*models.TicketInternal, error)
	FindAll(scopeParam interface{}, page, limit int, status string, departmentID, filialeID *uint) ([]models.TicketInternal, int64, error)
	FindPanierByUser(userID uint, page, limit int) ([]models.TicketInternal, int64, error)
	GetStatsForDashboard(scopeParam interface{}) (total int, byStatus map[string]int, open int, closed int, err error)
	GetPerformanceByAssignedUser(userID uint) (totalAssigned, resolved, inProgress, open int, totalTimeSpent int, err error)
	Update(t *models.TicketInternal) error
	UpdateFields(id uint, updates map[string]interface{}) error
	Delete(id uint) error
	GetNextSequenceNumber(year int) (int, error)
	CodeExists(code string) (bool, error)
}

type ticketInternalRepository struct{}

// NewTicketInternalRepository crée une nouvelle instance
func NewTicketInternalRepository() TicketInternalRepository {
	return &ticketInternalRepository{}
}

func applyTicketInternalPreloads(query *gorm.DB) *gorm.DB {
	return query.Preload("Department").Preload("Department.Filiale").
		Preload("Filiale").
		Preload("CreatedBy").Preload("CreatedBy.Department").
		Preload("AssignedTo").Preload("AssignedTo.Department").
		Preload("ValidatedBy").Preload("Ticket")
}

// Create crée un ticket interne
func (r *ticketInternalRepository) Create(t *models.TicketInternal) error {
	return database.DB.Create(t).Error
}

// FindByID trouve un ticket interne par ID avec relations
func (r *ticketInternalRepository) FindByID(id uint) (*models.TicketInternal, error) {
	var t models.TicketInternal
	if err := applyTicketInternalPreloads(database.DB).First(&t, id).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

// FindPanierByUser récupère les tickets internes assignés à l'utilisateur et non clôturés (panier)
func (r *ticketInternalRepository) FindPanierByUser(userID uint, page, limit int) ([]models.TicketInternal, int64, error) {
	var list []models.TicketInternal
	if userID == 0 {
		return list, 0, nil
	}
	query := database.DB.Model(&models.TicketInternal{}).
		Where("assigned_to_id = ?", userID).
		Where("status != ?", "cloture")
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * limit
	err := applyTicketInternalPreloads(query).
		Order("ticket_internes.updated_at DESC").
		Offset(offset).Limit(limit).
		Find(&list).Error
	return list, total, err
}

// FindAll liste les tickets internes avec scope et pagination
func (r *ticketInternalRepository) FindAll(scopeParam interface{}, page, limit int, status string, departmentID, filialeID *uint) ([]models.TicketInternal, int64, error) {
	query := database.DB.Model(&models.TicketInternal{})
	if scopeParam != nil {
		if s, ok := scopeParam.(*scope.QueryScope); ok {
			query = scope.ApplyTicketInternalScopeToTable(query, s)
		}
	}
	if departmentID != nil {
		query = query.Where("department_id = ?", *departmentID)
	}
	if filialeID != nil {
		query = query.Where("filiale_id = ?", *filialeID)
	}
	if status != "" && status != "all" {
		query = query.Where("status = ?", status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	var list []models.TicketInternal
	err := applyTicketInternalPreloads(query).Order("ticket_internes.created_at DESC").Offset(offset).Limit(limit).Find(&list).Error
	return list, total, err
}

// GetStatsForDashboard retourne les statistiques des tickets internes pour le tableau de bord (même scope que FindAll)
func (r *ticketInternalRepository) GetStatsForDashboard(scopeParam interface{}) (total int, byStatus map[string]int, open int, closed int, err error) {
	byStatus = make(map[string]int)
	query := database.DB.Model(&models.TicketInternal{})
	if scopeParam != nil {
		if s, ok := scopeParam.(*scope.QueryScope); ok {
			query = scope.ApplyTicketInternalScopeToTable(query, s)
		}
	}
	var total64 int64
	if err := query.Count(&total64).Error; err != nil {
		return 0, nil, 0, 0, err
	}
	total = int(total64)

	type row struct {
		Status string `gorm:"column:status"`
		Count  int64  `gorm:"column:count"`
	}
	var rows []row
	if err := query.Select("status, COUNT(*) as count").Group("status").Scan(&rows).Error; err != nil {
		return total, byStatus, 0, 0, err
	}
	for _, rw := range rows {
		byStatus[rw.Status] = int(rw.Count)
		switch rw.Status {
		case "ouvert", "en_cours", "en_attente":
			open += int(rw.Count)
		case "resolu", "cloture":
			closed += int(rw.Count)
		}
	}
	return total, byStatus, open, closed, nil
}

// GetPerformanceByAssignedUser retourne les métriques de performance pour les tickets internes assignés à l'utilisateur
func (r *ticketInternalRepository) GetPerformanceByAssignedUser(userID uint) (totalAssigned, resolved, inProgress, open int, totalTimeSpent int, err error) {
	if userID == 0 {
		return 0, 0, 0, 0, 0, nil
	}
	assignedFilter := database.DB.Model(&models.TicketInternal{}).Where("assigned_to_id = ?", userID)
	var total64 int64
	if err := assignedFilter.Count(&total64).Error; err != nil {
		return 0, 0, 0, 0, 0, err
	}
	totalAssigned = int(total64)

	var resolved64, inProgress64, open64 int64
	if err := database.DB.Model(&models.TicketInternal{}).Where("assigned_to_id = ?", userID).Where("status = ?", "cloture").Count(&resolved64).Error; err != nil {
		return totalAssigned, 0, 0, 0, 0, err
	}
	resolved = int(resolved64)
	if err := database.DB.Model(&models.TicketInternal{}).Where("assigned_to_id = ?", userID).Where("status = ?", "en_cours").Count(&inProgress64).Error; err != nil {
		return totalAssigned, resolved, 0, 0, 0, err
	}
	inProgress = int(inProgress64)
	if err := database.DB.Model(&models.TicketInternal{}).Where("assigned_to_id = ?", userID).Where("status IN ?", []string{"ouvert", "en_attente", "valide"}).Count(&open64).Error; err != nil {
		return totalAssigned, resolved, inProgress, 0, 0, err
	}
	open = int(open64)

	var sumTime int64
	if err := database.DB.Model(&models.TicketInternal{}).Where("assigned_to_id = ?", userID).Select("COALESCE(SUM(COALESCE(actual_time, 0)), 0)").Scan(&sumTime).Error; err != nil {
		return totalAssigned, resolved, inProgress, open, 0, err
	}
	totalTimeSpent = int(sumTime)
	return totalAssigned, resolved, inProgress, open, totalTimeSpent, nil
}

// Update met à jour un ticket interne
func (r *ticketInternalRepository) Update(t *models.TicketInternal) error {
	return database.DB.Save(t).Error
}

// UpdateFields met à jour des champs spécifiques
func (r *ticketInternalRepository) UpdateFields(id uint, updates map[string]interface{}) error {
	return database.DB.Model(&models.TicketInternal{}).Where("id = ?", id).Updates(updates).Error
}

// Delete supprime (soft delete) un ticket interne
func (r *ticketInternalRepository) Delete(id uint) error {
	return database.DB.Delete(&models.TicketInternal{}, id).Error
}

// GetNextSequenceNumber retourne le prochain numéro de séquence pour l'année (code TKI-YYYY-NNNN)
func (r *ticketInternalRepository) GetNextSequenceNumber(year int) (int, error) {
	codePattern := fmt.Sprintf("TKI-%d-%%", year)
	var list []models.TicketInternal
	err := database.DB.Unscoped().Model(&models.TicketInternal{}).Where("code LIKE ?", codePattern).Select("code").Find(&list).Error
	if err != nil {
		return 0, err
	}
	maxSeq := 0
	for _, t := range list {
		var y, seq int
		if _, err := fmt.Sscanf(t.Code, "TKI-%d-%d", &y, &seq); err == nil && y == year && seq > maxSeq {
			maxSeq = seq
		}
	}
	return maxSeq + 1, nil
}

// CodeExists vérifie si un code existe déjà
func (r *ticketInternalRepository) CodeExists(code string) (bool, error) {
	var count int64
	err := database.DB.Model(&models.TicketInternal{}).Where("code = ?", code).Count(&count).Error
	return count > 0, err
}
