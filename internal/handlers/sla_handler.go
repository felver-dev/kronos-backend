package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/services"
	"github.com/mcicare/itsm-backend/internal/utils"
)

// SLAHandler gère les handlers des SLA
type SLAHandler struct {
	slaService services.SLAService
}

// NewSLAHandler crée une nouvelle instance de SLAHandler
func NewSLAHandler(slaService services.SLAService) *SLAHandler {
	return &SLAHandler{
		slaService: slaService,
	}
}

// Create crée un nouveau SLA
func (h *SLAHandler) Create(c *gin.Context) {
	var req dto.CreateSLARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	createdByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	sla, err := h.slaService.Create(req, createdByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.CreatedResponse(c, sla, "SLA créé avec succès")
}

// GetByID récupère un SLA par son ID
func (h *SLAHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	sla, err := h.slaService.GetByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "SLA introuvable")
		return
	}

	utils.SuccessResponse(c, sla, "SLA récupéré avec succès")
}

// GetAll récupère tous les SLA
func (h *SLAHandler) GetAll(c *gin.Context) {
	slas, err := h.slaService.GetAll()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des SLA")
		return
	}

	utils.SuccessResponse(c, slas, "SLA récupérés avec succès")
}

// GetTicketSLAStatus récupère le statut SLA d'un ticket
func (h *SLAHandler) GetTicketSLAStatus(c *gin.Context) {
	ticketIDParam := c.Param("ticket_id")
	ticketID, err := strconv.ParseUint(ticketIDParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	status, err := h.slaService.GetTicketSLAStatus(uint(ticketID))
	if err != nil {
		utils.NotFoundResponse(c, "SLA introuvable pour ce ticket")
		return
	}

	utils.SuccessResponse(c, status, "Statut SLA récupéré avec succès")
}
