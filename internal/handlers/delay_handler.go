package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/services"
	"github.com/mcicare/itsm-backend/internal/utils"
)

// DelayHandler gère les handlers des retards
type DelayHandler struct {
	delayService services.DelayService
}

// NewDelayHandler crée une nouvelle instance de DelayHandler
func NewDelayHandler(delayService services.DelayService) *DelayHandler {
	return &DelayHandler{
		delayService: delayService,
	}
}

// GetByID récupère un retard par son ID
// @Summary Récupérer un retard par ID
// @Description Récupère un retard par son identifiant
// @Tags delays
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID du retard"
// @Success 200 {object} dto.DelayDTO
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /delays/{id} [get]
func (h *DelayHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	delay, err := h.delayService.GetByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Retard introuvable")
		return
	}

	utils.SuccessResponse(c, delay, "Retard récupéré avec succès")
}

// GetAll récupère tous les retards
// @Summary Récupérer tous les retards
// @Description Récupère la liste de tous les retards
// @Tags delays
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {array} dto.DelayDTO
// @Failure 500 {object} utils.Response
// @Router /delays [get]
func (h *DelayHandler) GetAll(c *gin.Context) {
	delays, err := h.delayService.GetAll()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des retards")
		return
	}

	utils.SuccessResponse(c, delays, "Retards récupérés avec succès")
}

// CreateJustification crée une justification pour un retard
// @Summary Créer une justification de retard
// @Description Crée une justification pour un retard
// @Tags delays
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param delay_id path int true "ID du retard"
// @Param request body dto.CreateDelayJustificationRequest true "Données de la justification"
// @Success 201 {object} dto.DelayJustificationDTO
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /delays/{delay_id}/justifications [post]
func (h *DelayHandler) CreateJustification(c *gin.Context) {
	delayIDParam := c.Param("delay_id")
	delayID, err := strconv.ParseUint(delayIDParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	var req dto.CreateDelayJustificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	justification, err := h.delayService.CreateJustification(uint(delayID), req, userID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.CreatedResponse(c, justification, "Justification créée avec succès")
}

// ValidateJustification valide une justification
// @Summary Valider une justification de retard
// @Description Valide une justification de retard
// @Tags delays
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID de la justification"
// @Param request body dto.ValidateDelayJustificationRequest true "Données de validation"
// @Success 200 {object} dto.DelayJustificationDTO
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /delays/justifications/{id}/validate [post]
func (h *DelayHandler) ValidateJustification(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	var req dto.ValidateDelayJustificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	validatedByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	justification, err := h.delayService.ValidateJustification(uint(id), req, validatedByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, justification, "Justification validée avec succès")
}
