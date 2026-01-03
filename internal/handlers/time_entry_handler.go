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
// @Summary Créer une entrée de temps
// @Description Crée une nouvelle entrée de temps dans le système
// @Tags time-entries
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateTimeEntryRequest true "Données de l'entrée de temps"
// @Success 201 {object} dto.TimeEntryDTO
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /time-entries [post]
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
// @Summary Récupérer une entrée de temps par ID
// @Description Récupère une entrée de temps par son identifiant
// @Tags time-entries
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID de l'entrée de temps"
// @Success 200 {object} dto.TimeEntryDTO
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /time-entries/{id} [get]
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
// @Summary Récupérer toutes les entrées de temps
// @Description Récupère la liste de toutes les entrées de temps
// @Tags time-entries
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {array} dto.TimeEntryDTO
// @Failure 500 {object} utils.Response
// @Router /time-entries [get]
func (h *TimeEntryHandler) GetAll(c *gin.Context) {
	timeEntries, err := h.timeEntryService.GetAll()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des entrées de temps")
		return
	}

	utils.SuccessResponse(c, timeEntries, "Entrées de temps récupérées avec succès")
}

// Validate valide une entrée de temps
// @Summary Valider une entrée de temps
// @Description Valide une entrée de temps
// @Tags time-entries
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID de l'entrée de temps"
// @Param request body dto.ValidateTimeEntryRequest true "Données de validation"
// @Success 200 {object} dto.TimeEntryDTO
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /time-entries/{id}/validate [post]
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
// @Summary Supprimer une entrée de temps
// @Description Supprime une entrée de temps du système
// @Tags time-entries
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID de l'entrée de temps"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /time-entries/{id} [delete]
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
