package services

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/mcicare/itsm-backend/database"
	"github.com/mcicare/itsm-backend/internal/models"
	"github.com/mcicare/itsm-backend/internal/repositories"
	"gorm.io/gorm"
)

// ProjectService interface pour les opérations sur les projets
type ProjectService interface {
	Create(name, description string, totalBudgetTime *int, startDate, endDate *string, createdByID uint) (*models.Project, error)
	GetByID(id uint) (*models.Project, error)
	GetAll(scope interface{}) ([]models.Project, error)
	GetByStatus(scope interface{}, status string) ([]models.Project, error)
	Update(id uint, name, description string, totalBudgetTime *int, status string, startDate, endDate *string, updatedByID uint) (*models.Project, error)
	Delete(id uint) error
	UpdateConsumedTime(projectID uint, consumedTime int) error
	AddBudgetExtension(projectID uint, additionalMinutes int, justification string, startDate, endDate *string, createdByID uint) (*models.ProjectBudgetExtension, error)
	GetBudgetExtensions(projectID uint) ([]models.ProjectBudgetExtension, error)
	UpdateBudgetExtension(projectID, extID uint, additionalMinutes int, justification string, startDate, endDate *string, updatedByID uint) (*models.ProjectBudgetExtension, error)
	DeleteBudgetExtension(projectID, extID uint) error
	// Phases
	GetPhases(projectID uint) ([]models.ProjectPhase, error)
	CreatePhase(projectID uint, name, description string, displayOrder int, status string) (*models.ProjectPhase, error)
	UpdatePhase(phaseID uint, name, description string, displayOrder *int, status string) (*models.ProjectPhase, error)
	DeletePhase(phaseID uint) error
	ReorderPhases(projectID uint, order []uint) error
	// Functions
	GetFunctions(projectID uint) ([]models.ProjectFunction, error)
	CreateFunction(projectID uint, name, functionType string, displayOrder int) (*models.ProjectFunction, error)
	UpdateFunction(functionID uint, name string, functionType *string, displayOrder *int) (*models.ProjectFunction, error)
	DeleteFunction(functionID uint) error
	// Members
	GetMembers(projectID uint) ([]models.ProjectMember, error)
	AddMember(projectID, userID uint, functionIDs []uint) (*models.ProjectMember, error)
	RemoveMember(projectID, userID uint) error
	SetMemberFunctions(projectID, userID uint, functionIDs []uint) error
	SetProjectManager(projectID, userID uint) error
	SetLead(projectID, userID uint) error
	// Phase members
	GetPhaseMembers(phaseID uint) ([]models.ProjectPhaseMember, error)
	AddPhaseMember(phaseID, userID uint, projectFunctionID *uint) (*models.ProjectPhaseMember, error)
	RemovePhaseMember(phaseID, userID uint) error
	SetPhaseMemberFunction(phaseID, userID uint, projectFunctionID *uint) error
	// Tasks
	GetTasks(projectID uint) ([]models.ProjectTask, error)
	GetTasksByPhase(phaseID uint) ([]models.ProjectTask, error)
	CreateTask(projectID, phaseID, createdByID uint, title, description, status, priority string, assigneeIDs []uint, estimatedTime *int, dueDate *string) (*models.ProjectTask, error)
	UpdateTask(taskID uint, title, description, status, priority string, assigneeIDs *[]uint, estimatedTime *int, actualTime *int, dueDate *string, projectPhaseID *uint) (*models.ProjectTask, error)
	DeleteTask(taskID uint) error
}

// projectService implémente ProjectService
type projectService struct {
	projectRepo        repositories.ProjectRepository
	userRepo           repositories.UserRepository
	budgetExtRepo      repositories.ProjectBudgetExtensionRepository
	phaseRepo          repositories.ProjectPhaseRepository
	functionRepo       repositories.ProjectFunctionRepository
	memberRepo         repositories.ProjectMemberRepository
	phaseMemberRepo    repositories.ProjectPhaseMemberRepository
	taskRepo           repositories.ProjectTaskRepository
	notificationService NotificationService
}

// NewProjectService crée une nouvelle instance de ProjectService
func NewProjectService(
	projectRepo repositories.ProjectRepository,
	userRepo repositories.UserRepository,
	budgetExtRepo repositories.ProjectBudgetExtensionRepository,
	phaseRepo repositories.ProjectPhaseRepository,
	functionRepo repositories.ProjectFunctionRepository,
	memberRepo repositories.ProjectMemberRepository,
	phaseMemberRepo repositories.ProjectPhaseMemberRepository,
	taskRepo repositories.ProjectTaskRepository,
	notificationService NotificationService,
) ProjectService {
	return &projectService{
		projectRepo:        projectRepo,
		userRepo:           userRepo,
		budgetExtRepo:      budgetExtRepo,
		phaseRepo:          phaseRepo,
		functionRepo:       functionRepo,
		memberRepo:         memberRepo,
		phaseMemberRepo:    phaseMemberRepo,
		taskRepo:           taskRepo,
		notificationService: notificationService,
	}
}

func parseOptionalDate(s *string) (*time.Time, error) {
	if s == nil || *s == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", *s)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// workingMinutesBetween calcule le nombre de minutes de travail (8 h/jour, jours ouvrés lun-ven)
// dans l’intervalle [start, end] inclus. Retourne 0 si end < start.
func workingMinutesBetween(start, end time.Time) int {
	s := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)
	e := time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, time.UTC)
	if e.Before(s) {
		return 0
	}
	var count int
	for d := s; !d.After(e); d = d.AddDate(0, 0, 1) {
		wd := d.Weekday()
		if wd != time.Saturday && wd != time.Sunday {
			count++
		}
	}
	return count * (8 * 60)
}

// Create crée un nouveau projet
func (s *projectService) Create(name, description string, totalBudgetTime *int, startDate, endDate *string, createdByID uint) (*models.Project, error) {
	// Vérifier que l'utilisateur existe et récupérer sa filiale pour le projet
	creator, err := s.userRepo.FindByID(createdByID)
	if err != nil {
		return nil, errors.New("utilisateur créateur introuvable")
	}
	// Filiale du projet : celle du créateur (user → département) pour que le projet apparaisse au tableau de bord département/filiale
	filialeID := creator.FilialeID
	if filialeID == nil && creator.Department != nil && creator.Department.FilialeID != nil {
		filialeID = creator.Department.FilialeID
	}

	createdByIDPtr := &createdByID
	project := &models.Project{
		Name:            name,
		Description:     description,
		TotalBudgetTime: totalBudgetTime,
		ConsumedTime:    0,
		Status:          "active",
		CreatedByID:     createdByIDPtr,
		FilialeID:       filialeID,
	}
	if t, err := parseOptionalDate(startDate); err != nil {
		return nil, errors.New("start_date invalide (attendu: AAAA-MM-JJ)")
	} else if t != nil {
		project.StartDate = t
	}
	if t, err := parseOptionalDate(endDate); err != nil {
		return nil, errors.New("end_date invalide (attendu: AAAA-MM-JJ)")
	} else if t != nil {
		project.EndDate = t
	}

	// Validation : budget temps ≤ temps de travail dans l’intervalle [début, fin] (jours ouvrés × 8 h/j)
	if project.StartDate != nil && project.EndDate != nil {
		if project.EndDate.Before(*project.StartDate) {
			return nil, errors.New("la date de fin prévue doit être postérieure ou égale à la date de début")
		}
		if project.TotalBudgetTime != nil && *project.TotalBudgetTime > 0 {
			max := workingMinutesBetween(*project.StartDate, *project.EndDate)
			if *project.TotalBudgetTime > max {
				return nil, errors.New("le budget temps ne peut pas dépasser le temps de travail disponible entre la date de début et la date de fin prévues (jours ouvrés × 8 h/j)")
			}
		}
	}

	if err := s.projectRepo.Create(project); err != nil {
		return nil, errors.New("erreur lors de la création du projet")
	}

	// Récupérer le projet créé
	createdProject, err := s.projectRepo.FindByID(project.ID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération du projet créé")
	}

	// Préremplir les deux fonctions de direction : Chef de projet et Lead
	pid := &createdProject.ID
	if err := s.functionRepo.Create(&models.ProjectFunction{ProjectID: pid, Name: "Chef de projet", Type: "direction", DisplayOrder: 0}); err != nil {
		log.Printf("[Create] project %d: création Chef de projet: %v", createdProject.ID, err)
	}
	if err := s.functionRepo.Create(&models.ProjectFunction{ProjectID: pid, Name: "Lead", Type: "direction", DisplayOrder: 1}); err != nil {
		log.Printf("[Create] project %d: création Lead: %v", createdProject.ID, err)
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
// Le scope est utilisé pour filtrer automatiquement selon les permissions de l'utilisateur
func (s *projectService) GetAll(scopeParam interface{}) ([]models.Project, error) {
	projects, err := s.projectRepo.FindAll(scopeParam)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des projets")
	}

	return projects, nil
}

// GetByStatus récupère les projets par statut
// Le scope est utilisé pour filtrer automatiquement selon les permissions de l'utilisateur
func (s *projectService) GetByStatus(scopeParam interface{}, status string) ([]models.Project, error) {
	projects, err := s.projectRepo.FindByStatus(scopeParam, status)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des projets")
	}

	return projects, nil
}

// Update met à jour un projet
func (s *projectService) Update(id uint, name, description string, totalBudgetTime *int, status string, startDate, endDate *string, updatedByID uint) (*models.Project, error) {
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
	if startDate != nil {
		if *startDate == "" {
			project.StartDate = nil
		} else if t, err := time.Parse("2006-01-02", *startDate); err != nil {
			return nil, errors.New("start_date invalide (attendu: AAAA-MM-JJ)")
		} else {
			project.StartDate = &t
		}
	}
	if endDate != nil {
		if *endDate == "" {
			project.EndDate = nil
		} else if t, err := time.Parse("2006-01-02", *endDate); err != nil {
			return nil, errors.New("end_date invalide (attendu: AAAA-MM-JJ)")
		} else {
			project.EndDate = &t
		}
	}

	// Validation : budget temps ≤ temps de travail dans l’intervalle [début, fin] (jours ouvrés × 8 h/j)
	if project.StartDate != nil && project.EndDate != nil {
		if project.EndDate.Before(*project.StartDate) {
			return nil, errors.New("la date de fin prévue doit être postérieure ou égale à la date de début")
		}
		if project.TotalBudgetTime != nil && *project.TotalBudgetTime > 0 {
			max := workingMinutesBetween(*project.StartDate, *project.EndDate)
			if *project.TotalBudgetTime > max {
				return nil, errors.New("le budget temps ne peut pas dépasser le temps de travail disponible entre la date de début et la date de fin prévues (jours ouvrés × 8 h/j)")
			}
		}
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

// Delete supprime un projet et toutes les données liées (cascade manuelle car FK ON DELETE RESTRICT).
func (s *projectService) Delete(id uint) error {
	_, err := s.projectRepo.FindByID(id)
	if err != nil {
		return errors.New("projet introuvable")
	}

	return database.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Tâches du projet : libérer time_entries puis supprimer commentaires, pièces jointes, historique, assignees, tâches
		var taskIDs []uint
		if err := tx.Model(&models.ProjectTask{}).Where("project_id = ?", id).Pluck("id", &taskIDs).Error; err != nil {
			log.Printf("Delete project: list tasks error: %v", err)
			return errors.New("erreur lors de la suppression du projet")
		}
		if len(taskIDs) > 0 {
			_ = tx.Model(&models.TimeEntry{}).Where("project_task_id IN ?", taskIDs).Update("project_task_id", nil).Error
			_ = tx.Where("project_task_id IN ?", taskIDs).Delete(&models.ProjectTaskComment{}).Error
			_ = tx.Where("project_task_id IN ?", taskIDs).Delete(&models.ProjectTaskAttachment{}).Error
			_ = tx.Where("project_task_id IN ?", taskIDs).Delete(&models.ProjectTaskHistory{}).Error
			_ = tx.Where("project_task_id IN ?", taskIDs).Delete(&models.ProjectTaskAssignee{}).Error
			if err := tx.Where("project_id = ?", id).Delete(&models.ProjectTask{}).Error; err != nil {
				log.Printf("Delete project: delete tasks error: %v", err)
				return errors.New("erreur lors de la suppression du projet")
			}
		}
		// 2. Membres des phases puis phases
		var phaseIDs []uint
		if err := tx.Model(&models.ProjectPhase{}).Where("project_id = ?", id).Pluck("id", &phaseIDs).Error; err != nil {
			log.Printf("Delete project: list phases error: %v", err)
			return errors.New("erreur lors de la suppression du projet")
		}
		if len(phaseIDs) > 0 {
			_ = tx.Where("project_phase_id IN ?", phaseIDs).Delete(&models.ProjectPhaseMember{}).Error
		}
		if err := tx.Where("project_id = ?", id).Delete(&models.ProjectPhase{}).Error; err != nil {
			log.Printf("Delete project: delete phases error: %v", err)
			return errors.New("erreur lors de la suppression du projet")
		}
		// 3. Fonctions des membres puis membres
		var memberIDs []uint
		if err := tx.Model(&models.ProjectMember{}).Where("project_id = ?", id).Pluck("id", &memberIDs).Error; err != nil {
			log.Printf("Delete project: list members error: %v", err)
			return errors.New("erreur lors de la suppression du projet")
		}
		if len(memberIDs) > 0 {
			_ = tx.Where("project_member_id IN ?", memberIDs).Delete(&models.ProjectMemberFunction{}).Error
		}
		if err := tx.Where("project_id = ?", id).Delete(&models.ProjectMember{}).Error; err != nil {
			log.Printf("Delete project: delete members error: %v", err)
			return errors.New("erreur lors de la suppression du projet")
		}
		// 4. Fonctions du projet (project_id = id, pas les globales NULL)
		if err := tx.Where("project_id = ?", id).Delete(&models.ProjectFunction{}).Error; err != nil {
			log.Printf("Delete project: delete functions error: %v", err)
			return errors.New("erreur lors de la suppression du projet")
		}
		// 5. Extensions de budget
		if err := tx.Where("project_id = ?", id).Delete(&models.ProjectBudgetExtension{}).Error; err != nil {
			log.Printf("Delete project: delete budget extensions error: %v", err)
			return errors.New("erreur lors de la suppression du projet")
		}
		// 6. Liaison tickets-projets
		if err := tx.Where("project_id = ?", id).Delete(&models.TicketProject{}).Error; err != nil {
			log.Printf("Delete project: delete ticket_projects error: %v", err)
			return errors.New("erreur lors de la suppression du projet")
		}
		// 7. Projet
		if err := tx.Delete(&models.Project{}, id).Error; err != nil {
			log.Printf("Delete project: delete project error: %v", err)
			return errors.New("erreur lors de la suppression du projet")
		}
		return nil
	})
}

// UpdateConsumedTime met à jour le temps consommé d'un projet
func (s *projectService) UpdateConsumedTime(projectID uint, consumedTime int) error {
	return s.projectRepo.UpdateConsumedTime(projectID, consumedTime)
}

// AddBudgetExtension ajoute une extension au budget temps du projet (temps + justification).
// startDate et endDate délimitent la période de l'extension (optionnel). Réservé aux projets clôturés.
func (s *projectService) AddBudgetExtension(projectID uint, additionalMinutes int, justification string, startDate, endDate *string, createdByID uint) (*models.ProjectBudgetExtension, error) {
	p, err := s.projectRepo.FindByID(projectID)
	if err != nil {
		return nil, errors.New("projet introuvable")
	}
	if p.Status != "completed" {
		return nil, errors.New("l'extension de budget n'est possible que pour un projet clôturé")
	}
	if additionalMinutes <= 0 {
		return nil, errors.New("le temps ajouté doit être strictement positif")
	}
	if justification == "" || len(justification) < 3 {
		return nil, errors.New("une justification d'au moins 3 caractères est requise")
	}
	ext := &models.ProjectBudgetExtension{
		ProjectID:         projectID,
		AdditionalMinutes: additionalMinutes,
		Justification:     justification,
		CreatedByID:       &createdByID,
	}
	if startDate != nil && *startDate != "" {
		t, err := time.Parse("2006-01-02", *startDate)
		if err != nil {
			return nil, errors.New("start_date invalide (attendu: AAAA-MM-JJ)")
		}
		ext.StartDate = &t
	}
	if endDate != nil && *endDate != "" {
		t, err := time.Parse("2006-01-02", *endDate)
		if err != nil {
			return nil, errors.New("end_date invalide (attendu: AAAA-MM-JJ)")
		}
		ext.EndDate = &t
	}
	if (ext.StartDate != nil) != (ext.EndDate != nil) {
		return nil, errors.New("renseignez les deux dates (début et fin de la période) ou aucune")
	}
	if ext.StartDate != nil && ext.EndDate != nil {
		if ext.EndDate.Before(*ext.StartDate) {
			return nil, errors.New("la date de fin de l'extension doit être postérieure ou égale à la date de début")
		}
		if additionalMinutes > workingMinutesBetween(*ext.StartDate, *ext.EndDate) {
			return nil, errors.New("le temps ajouté ne peut pas dépasser le temps de travail disponible entre les dates de l'extension (jours ouvrés × 8 h/j)")
		}
	}
	if err := s.budgetExtRepo.Create(ext); err != nil {
		return nil, errors.New("erreur lors de l'enregistrement de l'extension")
	}
	if err := s.projectRepo.IncrementTotalBudgetTime(projectID, additionalMinutes); err != nil {
		return nil, errors.New("erreur lors de la mise à jour du budget")
	}
	// Réactiver le projet (statut actif) : une extension indique que le projet reprend
	if err := s.projectRepo.UpdateStatus(projectID, "active"); err != nil {
		return nil, errors.New("erreur lors de la réactivation du projet")
	}
	// Recharger avec CreatedBy pour la réponse
	list, _ := s.budgetExtRepo.FindByProjectID(projectID)
	for i := range list {
		if list[i].ID == ext.ID {
			return &list[i], nil
		}
	}
	return ext, nil
}

// GetBudgetExtensions retourne l'historique des extensions de budget d'un projet
func (s *projectService) GetBudgetExtensions(projectID uint) ([]models.ProjectBudgetExtension, error) {
	return s.budgetExtRepo.FindByProjectID(projectID)
}

// UpdateBudgetExtension met à jour une extension de budget (temps, justification, période).
// Recalcule le delta sur total_budget_time : IncrementTotalBudgetTime(projectID, new - old).
func (s *projectService) UpdateBudgetExtension(projectID, extID uint, additionalMinutes int, justification string, startDate, endDate *string, updatedByID uint) (*models.ProjectBudgetExtension, error) {
	ext, err := s.budgetExtRepo.FindByID(extID)
	if err != nil || ext == nil {
		return nil, errors.New("extension introuvable")
	}
	if ext.ProjectID != projectID {
		return nil, errors.New("cette extension n'appartient pas à ce projet")
	}
	if additionalMinutes <= 0 {
		return nil, errors.New("le temps ajouté doit être strictement positif")
	}
	if justification == "" || len(justification) < 3 {
		return nil, errors.New("une justification d'au moins 3 caractères est requise")
	}
	oldMinutes := ext.AdditionalMinutes
	ext.AdditionalMinutes = additionalMinutes
	ext.Justification = justification
	if startDate != nil && *startDate != "" {
		t, err := time.Parse("2006-01-02", *startDate)
		if err != nil {
			return nil, errors.New("start_date invalide (attendu: AAAA-MM-JJ)")
		}
		ext.StartDate = &t
	} else {
		ext.StartDate = nil
	}
	if endDate != nil && *endDate != "" {
		t, err := time.Parse("2006-01-02", *endDate)
		if err != nil {
			return nil, errors.New("end_date invalide (attendu: AAAA-MM-JJ)")
		}
		ext.EndDate = &t
	} else {
		ext.EndDate = nil
	}
	if (ext.StartDate != nil) != (ext.EndDate != nil) {
		return nil, errors.New("renseignez les deux dates (début et fin de la période) ou aucune")
	}
	if ext.StartDate != nil && ext.EndDate != nil {
		if ext.EndDate.Before(*ext.StartDate) {
			return nil, errors.New("la date de fin de l'extension doit être postérieure ou égale à la date de début")
		}
		if additionalMinutes > workingMinutesBetween(*ext.StartDate, *ext.EndDate) {
			return nil, errors.New("le temps ajouté ne peut pas dépasser le temps de travail disponible entre les dates de l'extension (jours ouvrés × 8 h/j)")
		}
	}
	if err := s.budgetExtRepo.Update(ext); err != nil {
		return nil, errors.New("erreur lors de la mise à jour de l'extension")
	}
	delta := additionalMinutes - oldMinutes
	if delta != 0 {
		if err := s.projectRepo.IncrementTotalBudgetTime(projectID, delta); err != nil {
			return nil, errors.New("erreur lors de la mise à jour du budget")
		}
	}
	_ = updatedByID
	list, _ := s.budgetExtRepo.FindByProjectID(projectID)
	for i := range list {
		if list[i].ID == ext.ID {
			return &list[i], nil
		}
	}
	return ext, nil
}

// DeleteBudgetExtension supprime une extension et diminue le budget du projet de ce montant.
func (s *projectService) DeleteBudgetExtension(projectID, extID uint) error {
	ext, err := s.budgetExtRepo.FindByID(extID)
	if err != nil || ext == nil {
		return errors.New("extension introuvable")
	}
	if ext.ProjectID != projectID {
		return errors.New("cette extension n'appartient pas à ce projet")
	}
	if err := s.projectRepo.IncrementTotalBudgetTime(projectID, -ext.AdditionalMinutes); err != nil {
		return errors.New("erreur lors de la mise à jour du budget")
	}
	if err := s.budgetExtRepo.Delete(ext.ID); err != nil {
		return errors.New("erreur lors de la suppression de l'extension")
	}
	return nil
}

// --- Phases ---
func (s *projectService) GetPhases(projectID uint) ([]models.ProjectPhase, error) {
	if _, err := s.projectRepo.FindByID(projectID); err != nil {
		return nil, errors.New("projet introuvable")
	}
	return s.phaseRepo.FindByProjectID(projectID)
}

func (s *projectService) CreatePhase(projectID uint, name, description string, displayOrder int, status string) (*models.ProjectPhase, error) {
	if _, err := s.projectRepo.FindByID(projectID); err != nil {
		return nil, errors.New("projet introuvable")
	}
	if name == "" {
		return nil, errors.New("le nom de l'étape est requis")
	}
	if status == "" {
		status = "not_started"
	}
	p := &models.ProjectPhase{ProjectID: projectID, Name: name, Description: description, DisplayOrder: displayOrder, Status: status}
	if err := s.phaseRepo.Create(p); err != nil {
		return nil, err
	}
	return s.phaseRepo.FindByID(p.ID)
}

func (s *projectService) UpdatePhase(phaseID uint, name, description string, displayOrder *int, status string) (*models.ProjectPhase, error) {
	p, err := s.phaseRepo.FindByID(phaseID)
	if err != nil {
		return nil, errors.New("étape introuvable")
	}
	if name != "" {
		p.Name = name
	}
	if description != "" {
		p.Description = description
	}
	if displayOrder != nil {
		p.DisplayOrder = *displayOrder
	}
	if status != "" {
		p.Status = status
	}
	if err := s.phaseRepo.Update(p); err != nil {
		return nil, err
	}
	return s.phaseRepo.FindByID(phaseID)
}

func (s *projectService) DeletePhase(phaseID uint) error {
	if _, err := s.phaseRepo.FindByID(phaseID); err != nil {
		return errors.New("étape introuvable")
	}
	return s.phaseRepo.Delete(phaseID)
}

func (s *projectService) ReorderPhases(projectID uint, order []uint) error {
	if _, err := s.projectRepo.FindByID(projectID); err != nil {
		return errors.New("projet introuvable")
	}
	return s.phaseRepo.Reorder(projectID, order)
}

// --- Functions ---
func (s *projectService) GetFunctions(projectID uint) ([]models.ProjectFunction, error) {
	if _, err := s.projectRepo.FindByID(projectID); err != nil {
		return nil, errors.New("projet introuvable")
	}
	return s.functionRepo.FindByProjectID(projectID)
}

func (s *projectService) CreateFunction(projectID uint, name, functionType string, displayOrder int) (*models.ProjectFunction, error) {
	if _, err := s.projectRepo.FindByID(projectID); err != nil {
		return nil, errors.New("projet introuvable")
	}
	if name == "" {
		return nil, errors.New("le nom de la fonction est requis")
	}
	if functionType != "direction" && functionType != "execution" {
		functionType = "execution"
	}
	pid := &projectID
	f := &models.ProjectFunction{ProjectID: pid, Name: name, Type: functionType, DisplayOrder: displayOrder}
	if err := s.functionRepo.Create(f); err != nil {
		return nil, err
	}
	return s.functionRepo.FindByID(f.ID)
}

func (s *projectService) UpdateFunction(functionID uint, name string, functionType *string, displayOrder *int) (*models.ProjectFunction, error) {
	f, err := s.functionRepo.FindByID(functionID)
	if err != nil {
		return nil, errors.New("fonction introuvable")
	}
	if name != "" {
		f.Name = name
	}
	if functionType != nil && (*functionType == "direction" || *functionType == "execution") {
		f.Type = *functionType
	}
	if displayOrder != nil {
		f.DisplayOrder = *displayOrder
	}
	if err := s.functionRepo.Update(f); err != nil {
		return nil, err
	}
	return s.functionRepo.FindByID(functionID)
}

func (s *projectService) DeleteFunction(functionID uint) error {
	if _, err := s.functionRepo.FindByID(functionID); err != nil {
		return errors.New("fonction introuvable")
	}
	return s.functionRepo.Delete(functionID)
}

// --- Members ---
func (s *projectService) GetMembers(projectID uint) ([]models.ProjectMember, error) {
	if _, err := s.projectRepo.FindByID(projectID); err != nil {
		return nil, errors.New("projet introuvable")
	}
	return s.memberRepo.FindByProjectID(projectID)
}

func (s *projectService) AddMember(projectID, userID uint, functionIDs []uint) (*models.ProjectMember, error) {
	if _, err := s.projectRepo.FindByID(projectID); err != nil {
		return nil, errors.New("projet introuvable")
	}
	if _, err := s.userRepo.FindByID(userID); err != nil {
		return nil, errors.New("utilisateur introuvable")
	}
	existing, _ := s.memberRepo.FindByProjectIDAndUserID(projectID, userID)
	if existing != nil {
		return nil, errors.New("cet utilisateur est déjà membre du projet")
	}
	for _, fid := range functionIDs {
		fn, err := s.functionRepo.FindByID(fid)
		if err != nil || fn == nil {
			return nil, errors.New("une des fonctions est introuvable")
		}
		if fn.ProjectID != nil && *fn.ProjectID != projectID {
			return nil, errors.New("une des fonctions n'appartient pas à ce projet")
		}
	}
	m := &models.ProjectMember{ProjectID: projectID, UserID: userID}
	if err := s.memberRepo.Create(m); err != nil {
		return nil, err
	}
	if len(functionIDs) > 0 {
		if err := s.memberRepo.ReplaceMemberFunctions(m.ID, functionIDs); err != nil {
			return nil, err
		}
	}
	list, _ := s.memberRepo.FindByProjectID(projectID)
	for i := range list {
		if list[i].ID == m.ID {
			return &list[i], nil
		}
	}
	// Notifier le membre qu'il a été ajouté au projet
	s.notifyProjectMemberAdded(projectID, userID)
	return m, nil
}

func (s *projectService) notifyProjectMemberAdded(projectID, userID uint) {
	if s.notificationService == nil {
		return
	}
	proj, err := s.projectRepo.FindByID(projectID)
	if err != nil || proj == nil {
		return
	}
	linkURL := fmt.Sprintf("/app/projects/%d", projectID)
	title := "Vous avez été ajouté à un projet"
	message := fmt.Sprintf("Vous avez été ajouté au projet « %s ». Consultez le projet pour voir les tâches et l'équipe.", proj.Name)
	metadata := map[string]any{"project_id": projectID, "project_name": proj.Name}
	if err := s.notificationService.Create(userID, "project_member_added", title, message, linkURL, metadata); err != nil {
		log.Printf("Erreur notification projet membre ajouté (user %d): %v", userID, err)
	}
}

func (s *projectService) RemoveMember(projectID, userID uint) error {
	m, err := s.memberRepo.FindByProjectIDAndUserID(projectID, userID)
	if err != nil || m == nil {
		return errors.New("membre introuvable")
	}
	// Retirer ce membre des tâches où il est assigné
	tasks, _ := s.taskRepo.FindByProjectID(projectID)
	for _, t := range tasks {
		var keep []uint
		for _, a := range t.Assignees {
			if a.UserID != userID {
				keep = append(keep, a.UserID)
			}
		}
		if len(keep) != len(t.Assignees) {
			_ = s.taskRepo.ReplaceAssignees(t.ID, keep)
		}
	}
	// Supprimer les liens project_member_functions avant de supprimer le membre (contrainte FK)
	if err := s.memberRepo.ReplaceMemberFunctions(m.ID, []uint{}); err != nil {
		return err
	}
	return s.memberRepo.Delete(m.ID)
}

func (s *projectService) SetMemberFunctions(projectID, userID uint, functionIDs []uint) error {
	m, err := s.memberRepo.FindByProjectIDAndUserID(projectID, userID)
	if err != nil || m == nil {
		return errors.New("membre introuvable")
	}
	for _, fid := range functionIDs {
		fn, err := s.functionRepo.FindByID(fid)
		if err != nil || fn == nil {
			return errors.New("une des fonctions est introuvable")
		}
		if fn.ProjectID != nil && *fn.ProjectID != projectID {
			return errors.New("une des fonctions n'appartient pas à ce projet")
		}
	}
	return s.memberRepo.ReplaceMemberFunctions(m.ID, functionIDs)
}

func (s *projectService) SetProjectManager(projectID, userID uint) error {
	if _, err := s.projectRepo.FindByID(projectID); err != nil {
		return errors.New("projet introuvable")
	}
	if userID != 0 {
		m, _ := s.memberRepo.FindByProjectIDAndUserID(projectID, userID)
		if m == nil {
			return errors.New("l'utilisateur doit être membre du projet pour être désigné chef de projet")
		}
	}
	return s.memberRepo.SetProjectManager(projectID, userID)
}

func (s *projectService) SetLead(projectID, userID uint) error {
	if _, err := s.projectRepo.FindByID(projectID); err != nil {
		return errors.New("projet introuvable")
	}
	if userID != 0 {
		m, _ := s.memberRepo.FindByProjectIDAndUserID(projectID, userID)
		if m == nil {
			return errors.New("l'utilisateur doit être membre du projet pour être désigné lead")
		}
	}
	return s.memberRepo.SetLead(projectID, userID)
}

// --- Phase members ---
func (s *projectService) GetPhaseMembers(phaseID uint) ([]models.ProjectPhaseMember, error) {
	if _, err := s.phaseRepo.FindByID(phaseID); err != nil {
		return nil, errors.New("étape introuvable")
	}
	return s.phaseMemberRepo.FindByPhaseID(phaseID)
}

func (s *projectService) AddPhaseMember(phaseID, userID uint, projectFunctionID *uint) (*models.ProjectPhaseMember, error) {
	if _, err := s.phaseRepo.FindByID(phaseID); err != nil {
		return nil, errors.New("étape introuvable")
	}
	if _, err := s.userRepo.FindByID(userID); err != nil {
		return nil, errors.New("utilisateur introuvable")
	}
	existing, _ := s.phaseMemberRepo.FindByPhaseIDAndUserID(phaseID, userID)
	if existing != nil {
		return nil, errors.New("cet utilisateur est déjà membre de l'étape")
	}
	m := &models.ProjectPhaseMember{ProjectPhaseID: phaseID, UserID: userID, ProjectFunctionID: projectFunctionID}
	if err := s.phaseMemberRepo.Create(m); err != nil {
		return nil, err
	}
	list, _ := s.phaseMemberRepo.FindByPhaseID(phaseID)
	for i := range list {
		if list[i].ID == m.ID {
			return &list[i], nil
		}
	}
	// Notifier le membre qu'il a été ajouté à l'étape (et donc au projet)
	s.notifyPhaseMemberAdded(phaseID, userID)
	return m, nil
}

func (s *projectService) notifyPhaseMemberAdded(phaseID, userID uint) {
	if s.notificationService == nil {
		return
	}
	ph, err := s.phaseRepo.FindByID(phaseID)
	if err != nil || ph == nil {
		return
	}
	proj, err := s.projectRepo.FindByID(ph.ProjectID)
	if err != nil || proj == nil {
		return
	}
	linkURL := fmt.Sprintf("/app/projects/%d", ph.ProjectID)
	title := "Vous avez été ajouté à une étape de projet"
	message := fmt.Sprintf("Vous avez été ajouté à l'étape « %s » du projet « %s ». Consultez le projet pour voir les tâches.", ph.Name, proj.Name)
	metadata := map[string]any{"project_id": ph.ProjectID, "phase_id": phaseID, "phase_name": ph.Name, "project_name": proj.Name}
	if err := s.notificationService.Create(userID, "project_phase_member_added", title, message, linkURL, metadata); err != nil {
		log.Printf("Erreur notification étape projet membre ajouté (user %d): %v", userID, err)
	}
}

func (s *projectService) RemovePhaseMember(phaseID, userID uint) error {
	m, err := s.phaseMemberRepo.FindByPhaseIDAndUserID(phaseID, userID)
	if err != nil || m == nil {
		return errors.New("membre d'étape introuvable")
	}
	return s.phaseMemberRepo.Delete(m.ID)
}

func (s *projectService) SetPhaseMemberFunction(phaseID, userID uint, projectFunctionID *uint) error {
	m, err := s.phaseMemberRepo.FindByPhaseIDAndUserID(phaseID, userID)
	if err != nil || m == nil {
		return errors.New("membre d'étape introuvable")
	}
	m.ProjectFunctionID = projectFunctionID
	return s.phaseMemberRepo.Update(m)
}

// --- Tasks ---
func (s *projectService) GetTasks(projectID uint) ([]models.ProjectTask, error) {
	if _, err := s.projectRepo.FindByID(projectID); err != nil {
		return nil, errors.New("projet introuvable")
	}
	return s.taskRepo.FindByProjectID(projectID)
}

func (s *projectService) GetTasksByPhase(phaseID uint) ([]models.ProjectTask, error) {
	if _, err := s.phaseRepo.FindByID(phaseID); err != nil {
		return nil, errors.New("étape introuvable")
	}
	return s.taskRepo.FindByPhaseID(phaseID)
}

func parseDate(s *string) *time.Time {
	if s == nil || *s == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02", *s)
	if err != nil {
		return nil
	}
	return &t
}

// ensureAssigneesAsMembers ajoute chaque assigné comme membre du projet (sans fonction) s'il n'est pas déjà membre.
func (s *projectService) ensureAssigneesAsMembers(projectID uint, assigneeIDs []uint) {
	for _, uid := range assigneeIDs {
		m, _ := s.memberRepo.FindByProjectIDAndUserID(projectID, uid)
		if m != nil {
			continue
		}
		newM := &models.ProjectMember{ProjectID: projectID, UserID: uid}
		if err := s.memberRepo.Create(newM); err != nil {
			log.Printf("[ensureAssigneesAsMembers] project=%d user=%d: %v", projectID, uid, err)
		}
	}
}

// recalcAndUpdateProjectConsumedTime recalcule le temps consommé du projet (somme des actual_time des tâches) et met à jour la colonne consumed_time.
func (s *projectService) recalcAndUpdateProjectConsumedTime(projectID uint) error {
	tasks, err := s.taskRepo.FindByProjectID(projectID)
	if err != nil {
		return err
	}
	sum := 0
	for _, t := range tasks {
		sum += t.ActualTime
	}
	return s.projectRepo.UpdateConsumedTime(projectID, sum)
}

func (s *projectService) CreateTask(projectID, phaseID, createdByID uint, title, description, status, priority string, assigneeIDs []uint, estimatedTime *int, dueDate *string) (*models.ProjectTask, error) {
	if _, err := s.projectRepo.FindByID(projectID); err != nil {
		return nil, errors.New("projet introuvable")
	}
	ph, err := s.phaseRepo.FindByID(phaseID)
	if err != nil || ph == nil || ph.ProjectID != projectID {
		return nil, errors.New("étape introuvable ou n'appartient pas au projet")
	}
	if _, err := s.userRepo.FindByID(createdByID); err != nil {
		return nil, errors.New("créateur introuvable")
	}
	if title == "" {
		return nil, errors.New("le titre de la tâche est requis")
	}
	if status == "" {
		status = "ouvert"
	}
	if priority == "" {
		priority = "medium"
	}
	code, err := s.taskRepo.GenerateCode(projectID)
	if err != nil {
		return nil, errors.New("erreur génération du code tâche")
	}
	t := &models.ProjectTask{
		ProjectID:      projectID,
		ProjectPhaseID: phaseID,
		Code:           code,
		Title:          title,
		Description:    description,
		Status:         status,
		Priority:       priority,
		CreatedByID:    createdByID,
		EstimatedTime:  estimatedTime,
		DueDate:        parseDate(dueDate),
	}
	if err := s.taskRepo.Create(t); err != nil {
		return nil, err
	}
	if len(assigneeIDs) > 0 {
		if err := s.taskRepo.ReplaceAssignees(t.ID, assigneeIDs); err != nil {
			return nil, err
		}
		s.ensureAssigneesAsMembers(projectID, assigneeIDs)
		for _, assigneeID := range assigneeIDs {
			s.notifyTaskAssigned(projectID, t.ID, t.Code, t.Title, assigneeID)
		}
	}
	_ = s.recalcAndUpdateProjectConsumedTime(projectID)
	return s.taskRepo.FindByID(t.ID)
}

func (s *projectService) notifyTaskAssigned(projectID, taskID uint, taskCode, taskTitle string, assigneeID uint) {
	if s.notificationService == nil {
		return
	}
	linkURL := fmt.Sprintf("/app/projects/%d", projectID)
	title := "Tâche assignée"
	message := fmt.Sprintf("Une tâche vous a été assignée : %s - %s. Consultez le projet pour plus de détails.", taskCode, taskTitle)
	metadata := map[string]any{"project_id": projectID, "task_id": taskID, "task_code": taskCode}
	if err := s.notificationService.Create(assigneeID, "project_task_assigned", title, message, linkURL, metadata); err != nil {
		log.Printf("Erreur notification tâche projet assignée (user %d): %v", assigneeID, err)
	}
}

func (s *projectService) UpdateTask(taskID uint, title, description, status, priority string, assigneeIDs *[]uint, estimatedTime *int, actualTime *int, dueDate *string, projectPhaseID *uint) (*models.ProjectTask, error) {
	t, err := s.taskRepo.FindByID(taskID)
	if err != nil {
		return nil, errors.New("tâche introuvable")
	}
	if title != "" {
		t.Title = title
	}
	if description != "" {
		t.Description = description
	}
	if status != "" {
		t.Status = status
		if status == "cloture" {
			now := time.Now()
			t.ClosedAt = &now
		}
	}
	if priority != "" {
		t.Priority = priority
	}
	if projectPhaseID != nil && *projectPhaseID != 0 {
		ph, err := s.phaseRepo.FindByID(*projectPhaseID)
		if err != nil || ph == nil || ph.ProjectID != t.ProjectID {
			return nil, errors.New("étape introuvable ou n'appartient pas au projet")
		}
		t.ProjectPhaseID = *projectPhaseID
		t.ProjectPhase = nil // évite que GORM réécrive project_phase_id depuis l'association préchargée
	}
	if estimatedTime != nil {
		t.EstimatedTime = estimatedTime
	}
	if actualTime != nil {
		t.ActualTime = *actualTime
	}
	if dueDate != nil {
		t.DueDate = parseDate(dueDate)
	}
	if err := s.taskRepo.Update(t); err != nil {
		return nil, err
	}
	if actualTime != nil {
		_ = s.taskRepo.UpdateActualTime(taskID, *actualTime)
	}
	if assigneeIDs != nil {
		// Notifier uniquement les nouveaux assignés (ceux qui n'étaient pas déjà assignés)
		existingIDs := make(map[uint]bool)
		for _, a := range t.Assignees {
			existingIDs[a.UserID] = true
		}
		for _, uid := range *assigneeIDs {
			if !existingIDs[uid] {
				s.notifyTaskAssigned(t.ProjectID, taskID, t.Code, t.Title, uid)
			}
		}
		if err := s.taskRepo.ReplaceAssignees(taskID, *assigneeIDs); err != nil {
			return nil, err
		}
		s.ensureAssigneesAsMembers(t.ProjectID, *assigneeIDs)
	}
	return s.taskRepo.FindByID(taskID)
}

func (s *projectService) DeleteTask(taskID uint) error {
	t, err := s.taskRepo.FindByID(taskID)
	if err != nil {
		return errors.New("tâche introuvable")
	}
	projectID := t.ProjectID
	if err := s.taskRepo.Delete(taskID); err != nil {
		return err
	}
	_ = s.recalcAndUpdateProjectConsumedTime(projectID)
	return nil
}
