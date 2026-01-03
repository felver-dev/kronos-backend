package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/services"
	"github.com/mcicare/itsm-backend/internal/utils"
)

// TimeEntryHandler gère les handlers des entrées de temps
type TimeEntryHandler struct {
	timeEntryService services.TimeEntryService
}

// NewTimeEntryHandler crée une nouvelle instance de TimeEntryHandler
func NewTimeEntryHandler(timeEntryService services.TimeEntryService) *TimeEntryHandler {
	return &TimeEntryHandler{
		timeEntryService: timeEntryService,
	}
}

// Create crée une nouvelle entrée de temps
func (h *TimeEntryHandler) Create(c *gin.Context) {
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

	timeEntry, err := h.timeEntryService.Create(req, userID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.CreatedResponse(c, timeEntry, "Entrée de temps créée avec succès")
}

// GetByID récupère une entrée de temps par son ID
func (h *TimeEntryHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	timeEntry, err := h.timeEntryService.GetByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Entrée de temps introuvable")
		return
	}

	utils.SuccessResponse(c, timeEntry, "Entrée de temps récupérée avec succès")
}

// GetAll récupère toutes les entrées de temps
func (h *TimeEntryHandler) GetAll(c *gin.Context) {
	timeEntries, err := h.timeEntryService.GetAll()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des entrées de temps")
		return
	}

	utils.SuccessResponse(c, timeEntries, "Entrées de temps récupérées avec succès")
}

// Validate valide une entrée de temps
func (h *TimeEntryHandler) Validate(c *gin.Context) {
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

	timeEntry, err := h.timeEntryService.Validate(uint(id), req, validatedByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, timeEntry, "Entrée de temps validée avec succès")
}

// Delete supprime une entrée de temps
func (h *TimeEntryHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	err = h.timeEntryService.Delete(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Entrée de temps introuvable")
		return
	}

	utils.SuccessResponse(c, nil, "Entrée de temps supprimée avec succès")
}

