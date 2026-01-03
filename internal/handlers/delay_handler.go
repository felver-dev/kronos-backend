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
func (h *DelayHandler) GetAll(c *gin.Context) {
	delays, err := h.delayService.GetAll()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des retards")
		return
	}

	utils.SuccessResponse(c, delays, "Retards récupérés avec succès")
}

// CreateJustification crée une justification pour un retard
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
