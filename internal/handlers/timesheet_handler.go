package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/services"
	"github.com/mcicare/itsm-backend/internal/utils"
)

// TimesheetHandler gère les handlers des timesheets
type TimesheetHandler struct {
	timesheetService services.TimesheetService
}

// NewTimesheetHandler crée une nouvelle instance de TimesheetHandler
func NewTimesheetHandler(timesheetService services.TimesheetService) *TimesheetHandler {
	return &TimesheetHandler{
		timesheetService: timesheetService,
	}
}

// ========== Saisie du temps par ticket ==========

// CreateTimeEntry crée une nouvelle entrée de temps
// @Summary Créer une entrée de temps
// @Description Crée une nouvelle entrée de temps dans le timesheet
// @Tags timesheet
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateTimeEntryRequest true "Données de l'entrée de temps"
// @Success 201 {object} dto.TimeEntryDTO
// @Failure 400 {object} utils.Response
// @Router /timesheet/entries [post]
func (h *TimesheetHandler) CreateTimeEntry(c *gin.Context) {
	var req dto.CreateTimeEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	timeEntry, err := h.timesheetService.CreateTimeEntry(req, userID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.CreatedResponse(c, timeEntry, "Entrée de temps créée avec succès")
}

// GetTimeEntries récupère toutes les entrées de temps
// @Summary Récupérer toutes les entrées de temps
// @Description Récupère la liste de toutes les entrées de temps
// @Tags timesheet
// @Security BearerAuth
// @Produce json
// @Success 200 {array} dto.TimeEntryDTO
// @Failure 500 {object} utils.Response
// @Router /timesheet/entries [get]
func (h *TimesheetHandler) GetTimeEntries(c *gin.Context) {
	entries, err := h.timesheetService.GetTimeEntries()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des entrées de temps")
		return
	}

	utils.SuccessResponse(c, entries, "Entrées de temps récupérées avec succès")
}

// GetTimeEntryByID récupère une entrée de temps par son ID
// @Summary Récupérer une entrée de temps par ID
// @Description Récupère une entrée de temps par son identifiant
// @Tags timesheet
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID de l'entrée de temps"
// @Success 200 {object} dto.TimeEntryDTO
// @Failure 404 {object} utils.Response
// @Router /timesheet/entries/{id} [get]
func (h *TimesheetHandler) GetTimeEntryByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	entry, err := h.timesheetService.GetTimeEntryByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Entrée de temps introuvable")
		return
	}

	utils.SuccessResponse(c, entry, "Entrée de temps récupérée avec succès")
}

// UpdateTimeEntry met à jour une entrée de temps
// @Summary Mettre à jour une entrée de temps
// @Description Met à jour une entrée de temps
// @Tags timesheet
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID de l'entrée de temps"
// @Param request body dto.UpdateTimeEntryRequest true "Données de mise à jour"
// @Success 200 {object} dto.TimeEntryDTO
// @Failure 400 {object} utils.Response
// @Router /timesheet/entries/{id} [put]
func (h *TimesheetHandler) UpdateTimeEntry(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	var req dto.UpdateTimeEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	entry, err := h.timesheetService.UpdateTimeEntry(uint(id), req, userID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, entry, "Entrée de temps mise à jour avec succès")
}

// GetTimeEntriesByTicketID récupère les entrées de temps d'un ticket
// @Summary Récupérer les entrées de temps d'un ticket
// @Description Récupère toutes les entrées de temps associées à un ticket
// @Tags timesheet
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID du ticket"
// @Success 200 {array} dto.TimeEntryDTO
// @Failure 500 {object} utils.Response
// @Router /tickets/{id}/time-entries [get]
func (h *TimesheetHandler) GetTimeEntriesByTicketID(c *gin.Context) {
	ticketIDParam := c.Param("id")
	ticketID, err := strconv.ParseUint(ticketIDParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID de ticket invalide")
		return
	}

	entries, err := h.timesheetService.GetTimeEntriesByTicketID(uint(ticketID))
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des entrées de temps")
		return
	}

	utils.SuccessResponse(c, entries, "Entrées de temps récupérées avec succès")
}

// GetTimeEntriesByUserID récupère les entrées de temps d'un utilisateur
// @Summary Récupérer les entrées de temps d'un utilisateur
// @Description Récupère toutes les entrées de temps d'un utilisateur
// @Tags timesheet
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID de l'utilisateur"
// @Success 200 {array} dto.TimeEntryDTO
// @Failure 500 {object} utils.Response
// @Router /users/{id}/time-entries [get]
func (h *TimesheetHandler) GetTimeEntriesByUserID(c *gin.Context) {
	userIDParam := c.Param("id")
	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID utilisateur invalide")
		return
	}

	entries, err := h.timesheetService.GetTimeEntriesByUserID(uint(userID))
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des entrées de temps")
		return
	}

	utils.SuccessResponse(c, entries, "Entrées de temps récupérées avec succès")
}

// GetTimeEntriesByDate récupère les entrées de temps d'une date
// @Summary Récupérer les entrées de temps d'une date
// @Description Récupère toutes les entrées de temps d'une date spécifique
// @Tags timesheet
// @Security BearerAuth
// @Produce json
// @Param date path string true "Date (format: YYYY-MM-DD)"
// @Success 200 {array} dto.TimeEntryDTO
// @Failure 400 {object} utils.Response
// @Router /timesheet/entries/by-date/{date} [get]
func (h *TimesheetHandler) GetTimeEntriesByDate(c *gin.Context) {
	dateParam := c.Param("date")
	date, err := time.Parse("2006-01-02", dateParam)
	if err != nil {
		utils.BadRequestResponse(c, "Format de date invalide, attendu: YYYY-MM-DD")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	entries, err := h.timesheetService.GetTimeEntriesByDate(date, userID.(uint))
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des entrées de temps")
		return
	}

	utils.SuccessResponse(c, entries, "Entrées de temps récupérées avec succès")
}

// ========== Déclaration par jour ==========

// GetDailyDeclaration récupère une déclaration journalière
// @Summary Récupérer une déclaration journalière
// @Description Récupère la déclaration journalière d'une date
// @Tags timesheet
// @Security BearerAuth
// @Produce json
// @Param date path string true "Date (format: YYYY-MM-DD)"
// @Success 200 {object} dto.DailyDeclarationDTO
// @Failure 404 {object} utils.Response
// @Router /timesheet/daily/{date} [get]
func (h *TimesheetHandler) GetDailyDeclaration(c *gin.Context) {
	dateParam := c.Param("date")
	date, err := time.Parse("2006-01-02", dateParam)
	if err != nil {
		utils.BadRequestResponse(c, "Format de date invalide, attendu: YYYY-MM-DD")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	declaration, err := h.timesheetService.GetDailyDeclaration(date, userID.(uint))
	if err != nil {
		utils.NotFoundResponse(c, "Déclaration introuvable")
		return
	}

	utils.SuccessResponse(c, declaration, "Déclaration récupérée avec succès")
}

// CreateOrUpdateDailyDeclaration crée ou met à jour une déclaration journalière
// @Summary Créer ou mettre à jour une déclaration journalière
// @Description Crée ou met à jour une déclaration journalière
// @Tags timesheet
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param date path string true "Date (format: YYYY-MM-DD)"
// @Param request body []dto.DailyTaskRequest true "Liste des tâches"
// @Success 200 {object} dto.DailyDeclarationDTO
// @Failure 400 {object} utils.Response
// @Router /timesheet/daily/{date} [post]
func (h *TimesheetHandler) CreateOrUpdateDailyDeclaration(c *gin.Context) {
	dateParam := c.Param("date")
	date, err := time.Parse("2006-01-02", dateParam)
	if err != nil {
		utils.BadRequestResponse(c, "Format de date invalide, attendu: YYYY-MM-DD")
		return
	}

	var tasks []dto.DailyTaskRequest
	if err := c.ShouldBindJSON(&tasks); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	declaration, err := h.timesheetService.CreateOrUpdateDailyDeclaration(date, userID.(uint), tasks)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, declaration, "Déclaration créée/mise à jour avec succès")
}

// GetDailyTasks récupère les tâches d'une déclaration journalière
// @Summary Récupérer les tâches d'une déclaration journalière
// @Description Récupère toutes les tâches d'une déclaration journalière
// @Tags timesheet
// @Security BearerAuth
// @Produce json
// @Param date path string true "Date (format: YYYY-MM-DD)"
// @Success 200 {array} dto.DailyTaskDTO
// @Failure 404 {object} utils.Response
// @Router /timesheet/daily/{date}/tasks [get]
func (h *TimesheetHandler) GetDailyTasks(c *gin.Context) {
	dateParam := c.Param("date")
	date, err := time.Parse("2006-01-02", dateParam)
	if err != nil {
		utils.BadRequestResponse(c, "Format de date invalide, attendu: YYYY-MM-DD")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	tasks, err := h.timesheetService.GetDailyTasks(date, userID.(uint))
	if err != nil {
		utils.NotFoundResponse(c, "Tâches introuvables")
		return
	}

	utils.SuccessResponse(c, tasks, "Tâches récupérées avec succès")
}

// CreateDailyTask crée une tâche dans une déclaration journalière
// @Summary Créer une tâche journalière
// @Description Crée une nouvelle tâche dans une déclaration journalière
// @Tags timesheet
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param date path string true "Date (format: YYYY-MM-DD)"
// @Param request body dto.DailyTaskRequest true "Données de la tâche"
// @Success 201 {object} dto.DailyTaskDTO
// @Failure 400 {object} utils.Response
// @Router /timesheet/daily/{date}/tasks [post]
func (h *TimesheetHandler) CreateDailyTask(c *gin.Context) {
	dateParam := c.Param("date")
	date, err := time.Parse("2006-01-02", dateParam)
	if err != nil {
		utils.BadRequestResponse(c, "Format de date invalide, attendu: YYYY-MM-DD")
		return
	}

	var task dto.DailyTaskRequest
	if err := c.ShouldBindJSON(&task); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	createdTask, err := h.timesheetService.CreateDailyTask(date, userID.(uint), task)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.CreatedResponse(c, createdTask, "Tâche créée avec succès")
}

// DeleteDailyTask supprime une tâche d'une déclaration journalière
// @Summary Supprimer une tâche journalière
// @Description Supprime une tâche d'une déclaration journalière
// @Tags timesheet
// @Security BearerAuth
// @Produce json
// @Param date path string true "Date (format: YYYY-MM-DD)"
// @Param taskId path int true "ID de la tâche"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /timesheet/daily/{date}/tasks/{taskId} [delete]
func (h *TimesheetHandler) DeleteDailyTask(c *gin.Context) {
	dateParam := c.Param("date")
	date, err := time.Parse("2006-01-02", dateParam)
	if err != nil {
		utils.BadRequestResponse(c, "Format de date invalide, attendu: YYYY-MM-DD")
		return
	}

	taskIDParam := c.Param("taskId")
	taskID, err := strconv.ParseUint(taskIDParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID de tâche invalide")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	err = h.timesheetService.DeleteDailyTask(date, userID.(uint), uint(taskID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, nil, "Tâche supprimée avec succès")
}

// GetDailySummary récupère le résumé d'une déclaration journalière
// @Summary Récupérer le résumé d'une déclaration journalière
// @Description Récupère le résumé (temps total, nombre de tâches) d'une déclaration journalière
// @Tags timesheet
// @Security BearerAuth
// @Produce json
// @Param date path string true "Date (format: YYYY-MM-DD)"
// @Success 200 {object} dto.DailySummaryDTO
// @Failure 404 {object} utils.Response
// @Router /timesheet/daily/{date}/summary [get]
func (h *TimesheetHandler) GetDailySummary(c *gin.Context) {
	dateParam := c.Param("date")
	date, err := time.Parse("2006-01-02", dateParam)
	if err != nil {
		utils.BadRequestResponse(c, "Format de date invalide, attendu: YYYY-MM-DD")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	summary, err := h.timesheetService.GetDailySummary(date, userID.(uint))
	if err != nil {
		utils.NotFoundResponse(c, "Résumé introuvable")
		return
	}

	utils.SuccessResponse(c, summary, "Résumé récupéré avec succès")
}

// GetDailyCalendar récupère le calendrier des déclarations journalières
// @Summary Récupérer le calendrier journalier
// @Description Récupère le calendrier des déclarations journalières dans une plage de dates
// @Tags timesheet
// @Security BearerAuth
// @Produce json
// @Param startDate query string false "Date de début (format: YYYY-MM-DD)"
// @Param endDate query string false "Date de fin (format: YYYY-MM-DD)"
// @Success 200 {array} dto.DailyCalendarDTO
// @Failure 400 {object} utils.Response
// @Router /timesheet/daily/calendar [get]
func (h *TimesheetHandler) GetDailyCalendar(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	startDateStr := c.DefaultQuery("startDate", time.Now().AddDate(0, 0, -30).Format("2006-01-02"))
	endDateStr := c.DefaultQuery("endDate", time.Now().Format("2006-01-02"))

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		utils.BadRequestResponse(c, "Format de date de début invalide")
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		utils.BadRequestResponse(c, "Format de date de fin invalide")
		return
	}

	calendar, err := h.timesheetService.GetDailyCalendar(userID.(uint), startDate, endDate)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération du calendrier")
		return
	}

	utils.SuccessResponse(c, calendar, "Calendrier récupéré avec succès")
}

// GetDailyRange récupère les déclarations journalières dans une plage de dates
// @Summary Récupérer les déclarations journalières dans une plage
// @Description Récupère toutes les déclarations journalières dans une plage de dates
// @Tags timesheet
// @Security BearerAuth
// @Produce json
// @Param startDate query string false "Date de début (format: YYYY-MM-DD)"
// @Param endDate query string false "Date de fin (format: YYYY-MM-DD)"
// @Success 200 {array} dto.DailyDeclarationDTO
// @Failure 400 {object} utils.Response
// @Router /timesheet/daily/range [get]
func (h *TimesheetHandler) GetDailyRange(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	startDateStr := c.DefaultQuery("startDate", time.Now().AddDate(0, 0, -30).Format("2006-01-02"))
	endDateStr := c.DefaultQuery("endDate", time.Now().Format("2006-01-02"))

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		utils.BadRequestResponse(c, "Format de date de début invalide")
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		utils.BadRequestResponse(c, "Format de date de fin invalide")
		return
	}

	declarations, err := h.timesheetService.GetDailyRange(userID.(uint), startDate, endDate)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des déclarations")
		return
	}

	utils.SuccessResponse(c, declarations, "Déclarations récupérées avec succès")
}

// ========== Déclaration par semaine ==========

// GetWeeklyDeclaration récupère une déclaration hebdomadaire
// @Summary Récupérer une déclaration hebdomadaire
// @Description Récupère la déclaration hebdomadaire d'une semaine
// @Tags timesheet
// @Security BearerAuth
// @Produce json
// @Param week path string true "Semaine (format: YYYY-Www)"
// @Success 200 {object} dto.WeeklyDeclarationDTO
// @Failure 404 {object} utils.Response
// @Router /timesheet/weekly/{week} [get]
func (h *TimesheetHandler) GetWeeklyDeclaration(c *gin.Context) {
	week := c.Param("week")

	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	declaration, err := h.timesheetService.GetWeeklyDeclaration(week, userID.(uint))
	if err != nil {
		utils.NotFoundResponse(c, "Déclaration introuvable")
		return
	}

	utils.SuccessResponse(c, declaration, "Déclaration récupérée avec succès")
}

// CreateOrUpdateWeeklyDeclaration crée ou met à jour une déclaration hebdomadaire
// @Summary Créer ou mettre à jour une déclaration hebdomadaire
// @Description Crée ou met à jour une déclaration hebdomadaire
// @Tags timesheet
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param week path string true "Semaine (format: YYYY-Www)"
// @Param request body []dto.WeeklyTaskRequest true "Liste des tâches"
// @Success 200 {object} dto.WeeklyDeclarationDTO
// @Failure 400 {object} utils.Response
// @Router /timesheet/weekly/{week} [post]
func (h *TimesheetHandler) CreateOrUpdateWeeklyDeclaration(c *gin.Context) {
	week := c.Param("week")

	var tasks []dto.WeeklyTaskRequest
	if err := c.ShouldBindJSON(&tasks); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	declaration, err := h.timesheetService.CreateOrUpdateWeeklyDeclaration(week, userID.(uint), tasks)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, declaration, "Déclaration créée/mise à jour avec succès")
}

// GetWeeklyTasks récupère les tâches d'une déclaration hebdomadaire
// @Summary Récupérer les tâches d'une déclaration hebdomadaire
// @Description Récupère toutes les tâches d'une déclaration hebdomadaire
// @Tags timesheet
// @Security BearerAuth
// @Produce json
// @Param week path string true "Semaine (format: YYYY-Www)"
// @Success 200 {array} dto.WeeklyTaskDTO
// @Failure 404 {object} utils.Response
// @Router /timesheet/weekly/{week}/tasks [get]
func (h *TimesheetHandler) GetWeeklyTasks(c *gin.Context) {
	week := c.Param("week")

	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	tasks, err := h.timesheetService.GetWeeklyTasks(week, userID.(uint))
	if err != nil {
		utils.NotFoundResponse(c, "Tâches introuvables")
		return
	}

	utils.SuccessResponse(c, tasks, "Tâches récupérées avec succès")
}

// GetWeeklySummary récupère le résumé d'une déclaration hebdomadaire
// @Summary Récupérer le résumé d'une déclaration hebdomadaire
// @Description Récupère le résumé (temps total, nombre de tâches) d'une déclaration hebdomadaire
// @Tags timesheet
// @Security BearerAuth
// @Produce json
// @Param week path string true "Semaine (format: YYYY-Www)"
// @Success 200 {object} dto.WeeklySummaryDTO
// @Failure 404 {object} utils.Response
// @Router /timesheet/weekly/{week}/summary [get]
func (h *TimesheetHandler) GetWeeklySummary(c *gin.Context) {
	week := c.Param("week")

	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	summary, err := h.timesheetService.GetWeeklySummary(week, userID.(uint))
	if err != nil {
		utils.NotFoundResponse(c, "Résumé introuvable")
		return
	}

	utils.SuccessResponse(c, summary, "Résumé récupéré avec succès")
}

// GetWeeklyDailyBreakdown récupère la répartition quotidienne d'une déclaration hebdomadaire
// @Summary Récupérer la répartition quotidienne
// @Description Récupère la répartition quotidienne d'une déclaration hebdomadaire
// @Tags timesheet
// @Security BearerAuth
// @Produce json
// @Param week path string true "Semaine (format: YYYY-Www)"
// @Success 200 {array} dto.DailyBreakdownDTO
// @Failure 404 {object} utils.Response
// @Router /timesheet/weekly/{week}/daily-breakdown [get]
func (h *TimesheetHandler) GetWeeklyDailyBreakdown(c *gin.Context) {
	week := c.Param("week")

	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	breakdown, err := h.timesheetService.GetWeeklyDailyBreakdown(week, userID.(uint))
	if err != nil {
		utils.NotFoundResponse(c, "Répartition introuvable")
		return
	}

	utils.SuccessResponse(c, breakdown, "Répartition récupérée avec succès")
}

// ValidateWeeklyDeclaration valide une déclaration hebdomadaire
// @Summary Valider une déclaration hebdomadaire
// @Description Valide une déclaration hebdomadaire
// @Tags timesheet
// @Security BearerAuth
// @Produce json
// @Param week path string true "Semaine (format: YYYY-Www)"
// @Success 200 {object} dto.WeeklyDeclarationDTO
// @Failure 400 {object} utils.Response
// @Router /timesheet/weekly/{week}/validate [post]
func (h *TimesheetHandler) ValidateWeeklyDeclaration(c *gin.Context) {
	week := c.Param("week")

	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	validatedByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	declaration, err := h.timesheetService.ValidateWeeklyDeclaration(week, userID.(uint), validatedByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, declaration, "Déclaration validée avec succès")
}

// GetWeeklyValidationStatus récupère le statut de validation d'une déclaration hebdomadaire
// @Summary Récupérer le statut de validation
// @Description Récupère le statut de validation d'une déclaration hebdomadaire
// @Tags timesheet
// @Security BearerAuth
// @Produce json
// @Param week path string true "Semaine (format: YYYY-Www)"
// @Success 200 {object} dto.ValidationStatusDTO
// @Failure 404 {object} utils.Response
// @Router /timesheet/weekly/{week}/validation-status [get]
func (h *TimesheetHandler) GetWeeklyValidationStatus(c *gin.Context) {
	week := c.Param("week")

	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	status, err := h.timesheetService.GetWeeklyValidationStatus(week, userID.(uint))
	if err != nil {
		utils.NotFoundResponse(c, "Statut introuvable")
		return
	}

	utils.SuccessResponse(c, status, "Statut récupéré avec succès")
}

// ========== Budget temps ==========

// SetTicketEstimatedTime définit le temps estimé d'un ticket
// @Summary Définir le temps estimé d'un ticket
// @Description Définit le temps estimé d'un ticket
// @Tags timesheet
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID du ticket"
// @Param estimatedTime body int true "Temps estimé en minutes"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /tickets/{id}/estimated-time [post]
func (h *TimesheetHandler) SetTicketEstimatedTime(c *gin.Context) {
	ticketIDParam := c.Param("id")
	ticketID, err := strconv.ParseUint(ticketIDParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID de ticket invalide")
		return
	}

	var estimatedTime int
	if err := c.ShouldBindJSON(&estimatedTime); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	err = h.timesheetService.SetTicketEstimatedTime(uint(ticketID), estimatedTime, userID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, nil, "Temps estimé défini avec succès")
}

// GetTicketEstimatedTime récupère le temps estimé d'un ticket
// @Summary Récupérer le temps estimé d'un ticket
// @Description Récupère le temps estimé d'un ticket
// @Tags timesheet
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID du ticket"
// @Success 200 {object} dto.EstimatedTimeDTO
// @Failure 404 {object} utils.Response
// @Router /tickets/{id}/estimated-time [get]
func (h *TimesheetHandler) GetTicketEstimatedTime(c *gin.Context) {
	ticketIDParam := c.Param("id")
	ticketID, err := strconv.ParseUint(ticketIDParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID de ticket invalide")
		return
	}

	estimatedTime, err := h.timesheetService.GetTicketEstimatedTime(uint(ticketID))
	if err != nil {
		utils.NotFoundResponse(c, "Temps estimé introuvable")
		return
	}

	utils.SuccessResponse(c, estimatedTime, "Temps estimé récupéré avec succès")
}

// UpdateTicketEstimatedTime met à jour le temps estimé d'un ticket
// @Summary Mettre à jour le temps estimé d'un ticket
// @Description Met à jour le temps estimé d'un ticket
// @Tags timesheet
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID du ticket"
// @Param estimatedTime body int true "Nouveau temps estimé en minutes"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /tickets/{id}/estimated-time [put]
func (h *TimesheetHandler) UpdateTicketEstimatedTime(c *gin.Context) {
	ticketIDParam := c.Param("id")
	ticketID, err := strconv.ParseUint(ticketIDParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID de ticket invalide")
		return
	}

	var estimatedTime int
	if err := c.ShouldBindJSON(&estimatedTime); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	err = h.timesheetService.UpdateTicketEstimatedTime(uint(ticketID), estimatedTime, userID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, nil, "Temps estimé mis à jour avec succès")
}

// GetTicketTimeComparison récupère la comparaison temps estimé vs réel d'un ticket
// @Summary Comparer temps estimé vs réel
// @Description Récupère la comparaison entre le temps estimé et le temps réel d'un ticket
// @Tags timesheet
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID du ticket"
// @Success 200 {object} dto.TimeComparisonDTO
// @Failure 404 {object} utils.Response
// @Router /tickets/{id}/time-comparison [get]
func (h *TimesheetHandler) GetTicketTimeComparison(c *gin.Context) {
	ticketIDParam := c.Param("id")
	ticketID, err := strconv.ParseUint(ticketIDParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID de ticket invalide")
		return
	}

	comparison, err := h.timesheetService.GetTicketTimeComparison(uint(ticketID))
	if err != nil {
		utils.NotFoundResponse(c, "Comparaison introuvable")
		return
	}

	utils.SuccessResponse(c, comparison, "Comparaison récupérée avec succès")
}

// GetProjectTimeBudget récupère le budget temps d'un projet
// @Summary Récupérer le budget temps d'un projet
// @Description Récupère le budget temps d'un projet
// @Tags timesheet
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID du projet"
// @Success 200 {object} dto.ProjectTimeBudgetDTO
// @Failure 404 {object} utils.Response
// @Router /projects/{id}/time-budget [get]
func (h *TimesheetHandler) GetProjectTimeBudget(c *gin.Context) {
	projectIDParam := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID de projet invalide")
		return
	}

	budget, err := h.timesheetService.GetProjectTimeBudget(uint(projectID))
	if err != nil {
		utils.NotFoundResponse(c, "Budget introuvable")
		return
	}

	utils.SuccessResponse(c, budget, "Budget récupéré avec succès")
}

// SetProjectTimeBudget définit le budget temps d'un projet
// @Summary Définir le budget temps d'un projet
// @Description Définit le budget temps d'un projet
// @Tags timesheet
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID du projet"
// @Param request body dto.SetProjectTimeBudgetRequest true "Données du budget"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /projects/{id}/time-budget [post]
func (h *TimesheetHandler) SetProjectTimeBudget(c *gin.Context) {
	projectIDParam := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID de projet invalide")
		return
	}

	var budget dto.SetProjectTimeBudgetRequest
	if err := c.ShouldBindJSON(&budget); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	err = h.timesheetService.SetProjectTimeBudget(uint(projectID), budget, userID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, nil, "Budget défini avec succès")
}

// GetBudgetAlerts récupère les alertes de budget
// @Summary Récupérer les alertes de budget
// @Description Récupère toutes les alertes de budget
// @Tags timesheet
// @Security BearerAuth
// @Produce json
// @Success 200 {array} dto.BudgetAlertDTO
// @Failure 500 {object} utils.Response
// @Router /timesheet/budget-alerts [get]
func (h *TimesheetHandler) GetBudgetAlerts(c *gin.Context) {
	alerts, err := h.timesheetService.GetBudgetAlerts()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des alertes")
		return
	}

	utils.SuccessResponse(c, alerts, "Alertes récupérées avec succès")
}

// GetTicketBudgetStatus récupère le statut du budget d'un ticket
// @Summary Récupérer le statut du budget d'un ticket
// @Description Récupère le statut du budget (on_budget, over_budget, under_budget) d'un ticket
// @Tags timesheet
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID du ticket"
// @Success 200 {object} dto.BudgetStatusDTO
// @Failure 404 {object} utils.Response
// @Router /timesheet/budget-status/{id} [get]
func (h *TimesheetHandler) GetTicketBudgetStatus(c *gin.Context) {
	ticketIDParam := c.Param("id")
	ticketID, err := strconv.ParseUint(ticketIDParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID de ticket invalide")
		return
	}

	status, err := h.timesheetService.GetTicketBudgetStatus(uint(ticketID))
	if err != nil {
		utils.NotFoundResponse(c, "Statut introuvable")
		return
	}

	utils.SuccessResponse(c, status, "Statut récupéré avec succès")
}

// ========== Validation ==========

// ValidateTimeEntry valide une entrée de temps
// @Summary Valider une entrée de temps
// @Description Valide une entrée de temps
// @Tags timesheet
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID de l'entrée de temps"
// @Param request body dto.ValidateTimeEntryRequest true "Données de validation"
// @Success 200 {object} dto.TimeEntryDTO
// @Failure 400 {object} utils.Response
// @Router /timesheet/entries/{id}/validate [post]
func (h *TimesheetHandler) ValidateTimeEntry(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	var req dto.ValidateTimeEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	validatedByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	entry, err := h.timesheetService.ValidateTimeEntry(uint(id), req, validatedByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, entry, "Entrée validée avec succès")
}

// GetPendingValidationEntries récupère les entrées en attente de validation
// @Summary Récupérer les entrées en attente de validation
// @Description Récupère toutes les entrées de temps en attente de validation
// @Tags timesheet
// @Security BearerAuth
// @Produce json
// @Success 200 {array} dto.TimeEntryDTO
// @Failure 500 {object} utils.Response
// @Router /timesheet/entries/pending-validation [get]
func (h *TimesheetHandler) GetPendingValidationEntries(c *gin.Context) {
	entries, err := h.timesheetService.GetPendingValidationEntries()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des entrées")
		return
	}

	utils.SuccessResponse(c, entries, "Entrées récupérées avec succès")
}

// GetValidationHistory récupère l'historique de validation
// @Summary Récupérer l'historique de validation
// @Description Récupère l'historique de toutes les validations
// @Tags timesheet
// @Security BearerAuth
// @Produce json
// @Success 200 {array} dto.ValidationHistoryDTO
// @Failure 500 {object} utils.Response
// @Router /timesheet/validation-history [get]
func (h *TimesheetHandler) GetValidationHistory(c *gin.Context) {
	history, err := h.timesheetService.GetValidationHistory()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération de l'historique")
		return
	}

	utils.SuccessResponse(c, history, "Historique récupéré avec succès")
}

// ========== Alertes ==========

// GetDelayAlerts récupère les alertes de retard
// @Summary Récupérer les alertes de retard
// @Description Récupère toutes les alertes de retard
// @Tags timesheet
// @Security BearerAuth
// @Produce json
// @Success 200 {array} dto.DelayAlertDTO
// @Failure 500 {object} utils.Response
// @Router /timesheet/alerts/delays [get]
func (h *TimesheetHandler) GetDelayAlerts(c *gin.Context) {
	alerts, err := h.timesheetService.GetDelayAlerts()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des alertes")
		return
	}

	utils.SuccessResponse(c, alerts, "Alertes récupérées avec succès")
}

// GetBudgetAlertsForTimesheet récupère les alertes de budget pour le timesheet
// @Summary Récupérer les alertes de budget
// @Description Récupère toutes les alertes de budget
// @Tags timesheet
// @Security BearerAuth
// @Produce json
// @Success 200 {array} dto.BudgetAlertDTO
// @Failure 500 {object} utils.Response
// @Router /timesheet/alerts/budget [get]
func (h *TimesheetHandler) GetBudgetAlertsForTimesheet(c *gin.Context) {
	alerts, err := h.timesheetService.GetBudgetAlertsForTimesheet()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des alertes")
		return
	}

	utils.SuccessResponse(c, alerts, "Alertes récupérées avec succès")
}

// GetOverloadAlerts récupère les alertes de surcharge
// @Summary Récupérer les alertes de surcharge
// @Description Récupère toutes les alertes de surcharge
// @Tags timesheet
// @Security BearerAuth
// @Produce json
// @Success 200 {array} dto.OverloadAlertDTO
// @Failure 500 {object} utils.Response
// @Router /timesheet/alerts/overload [get]
func (h *TimesheetHandler) GetOverloadAlerts(c *gin.Context) {
	alerts, err := h.timesheetService.GetOverloadAlerts()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des alertes")
		return
	}

	utils.SuccessResponse(c, alerts, "Alertes récupérées avec succès")
}

// GetUnderloadAlerts récupère les alertes de sous-charge
// @Summary Récupérer les alertes de sous-charge
// @Description Récupère toutes les alertes de sous-charge
// @Tags timesheet
// @Security BearerAuth
// @Produce json
// @Success 200 {array} dto.UnderloadAlertDTO
// @Failure 500 {object} utils.Response
// @Router /timesheet/alerts/underload [get]
func (h *TimesheetHandler) GetUnderloadAlerts(c *gin.Context) {
	alerts, err := h.timesheetService.GetUnderloadAlerts()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des alertes")
		return
	}

	utils.SuccessResponse(c, alerts, "Alertes récupérées avec succès")
}

// SendReminderAlerts envoie des rappels
// @Summary Envoyer des rappels
// @Description Envoie des rappels aux utilisateurs spécifiés
// @Tags timesheet
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body []uint true "Liste des IDs utilisateurs"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /timesheet/alerts/reminders [post]
func (h *TimesheetHandler) SendReminderAlerts(c *gin.Context) {
	var userIDs []uint
	if err := c.ShouldBindJSON(&userIDs); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	err := h.timesheetService.SendReminderAlerts(userIDs)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, nil, "Rappels envoyés avec succès")
}

// GetPendingJustificationAlerts récupère les alertes de justifications en attente
// @Summary Récupérer les alertes de justifications en attente
// @Description Récupère toutes les alertes de justifications en attente
// @Tags timesheet
// @Security BearerAuth
// @Produce json
// @Success 200 {array} dto.PendingJustificationAlertDTO
// @Failure 500 {object} utils.Response
// @Router /timesheet/alerts/justifications-pending [get]
func (h *TimesheetHandler) GetPendingJustificationAlerts(c *gin.Context) {
	alerts, err := h.timesheetService.GetPendingJustificationAlerts()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des alertes")
		return
	}

	utils.SuccessResponse(c, alerts, "Alertes récupérées avec succès")
}

// ========== Historique ==========

// GetTimesheetHistory récupère l'historique du timesheet
// @Summary Récupérer l'historique du timesheet
// @Description Récupère l'historique du timesheet dans une plage de dates
// @Tags timesheet
// @Security BearerAuth
// @Produce json
// @Param startDate query string false "Date de début (format: YYYY-MM-DD)"
// @Param endDate query string false "Date de fin (format: YYYY-MM-DD)"
// @Success 200 {array} dto.TimesheetHistoryDTO
// @Failure 400 {object} utils.Response
// @Router /timesheet/history [get]
func (h *TimesheetHandler) GetTimesheetHistory(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	startDateStr := c.DefaultQuery("startDate", time.Now().AddDate(0, 0, -30).Format("2006-01-02"))
	endDateStr := c.DefaultQuery("endDate", time.Now().Format("2006-01-02"))

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		utils.BadRequestResponse(c, "Format de date de début invalide")
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		utils.BadRequestResponse(c, "Format de date de fin invalide")
		return
	}

	history, err := h.timesheetService.GetTimesheetHistory(userID.(uint), startDate, endDate)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération de l'historique")
		return
	}

	utils.SuccessResponse(c, history, "Historique récupéré avec succès")
}

// GetTimesheetHistoryEntry récupère une entrée de l'historique
// @Summary Récupérer une entrée de l'historique
// @Description Récupère une entrée détaillée de l'historique
// @Tags timesheet
// @Security BearerAuth
// @Produce json
// @Param entryId path int true "ID de l'entrée"
// @Success 200 {object} dto.TimesheetHistoryEntryDTO
// @Failure 404 {object} utils.Response
// @Router /timesheet/history/{entryId} [get]
func (h *TimesheetHandler) GetTimesheetHistoryEntry(c *gin.Context) {
	entryIDParam := c.Param("entryId")
	entryID, err := strconv.ParseUint(entryIDParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	entry, err := h.timesheetService.GetTimesheetHistoryEntry(uint(entryID))
	if err != nil {
		utils.NotFoundResponse(c, "Entrée introuvable")
		return
	}

	utils.SuccessResponse(c, entry, "Entrée récupérée avec succès")
}

// GetTimesheetAuditTrail récupère la piste d'audit du timesheet
// @Summary Récupérer la piste d'audit
// @Description Récupère la piste d'audit du timesheet dans une plage de dates
// @Tags timesheet
// @Security BearerAuth
// @Produce json
// @Param startDate query string false "Date de début (format: YYYY-MM-DD)"
// @Param endDate query string false "Date de fin (format: YYYY-MM-DD)"
// @Success 200 {array} dto.AuditTrailDTO
// @Failure 400 {object} utils.Response
// @Router /timesheet/audit-trail [get]
func (h *TimesheetHandler) GetTimesheetAuditTrail(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	startDateStr := c.DefaultQuery("startDate", time.Now().AddDate(0, 0, -30).Format("2006-01-02"))
	endDateStr := c.DefaultQuery("endDate", time.Now().Format("2006-01-02"))

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		utils.BadRequestResponse(c, "Format de date de début invalide")
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		utils.BadRequestResponse(c, "Format de date de fin invalide")
		return
	}

	trail, err := h.timesheetService.GetTimesheetAuditTrail(userID.(uint), startDate, endDate)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération de la piste d'audit")
		return
	}

	utils.SuccessResponse(c, trail, "Piste d'audit récupérée avec succès")
}

// GetTimesheetModifications récupère les modifications du timesheet
// @Summary Récupérer les modifications
// @Description Récupère toutes les modifications du timesheet dans une plage de dates
// @Tags timesheet
// @Security BearerAuth
// @Produce json
// @Param startDate query string false "Date de début (format: YYYY-MM-DD)"
// @Param endDate query string false "Date de fin (format: YYYY-MM-DD)"
// @Success 200 {array} dto.ModificationDTO
// @Failure 400 {object} utils.Response
// @Router /timesheet/modifications [get]
func (h *TimesheetHandler) GetTimesheetModifications(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	startDateStr := c.DefaultQuery("startDate", time.Now().AddDate(0, 0, -30).Format("2006-01-02"))
	endDateStr := c.DefaultQuery("endDate", time.Now().Format("2006-01-02"))

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		utils.BadRequestResponse(c, "Format de date de début invalide")
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		utils.BadRequestResponse(c, "Format de date de fin invalide")
		return
	}

	modifications, err := h.timesheetService.GetTimesheetModifications(userID.(uint), startDate, endDate)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des modifications")
		return
	}

	utils.SuccessResponse(c, modifications, "Modifications récupérées avec succès")
}

