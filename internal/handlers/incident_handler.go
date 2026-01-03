package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/services"
	"github.com/mcicare/itsm-backend/internal/utils"
)

// IncidentHandler gère les handlers des incidents
type IncidentHandler struct {
	incidentService services.IncidentService
}

// NewIncidentHandler crée une nouvelle instance de IncidentHandler
func NewIncidentHandler(incidentService services.IncidentService) *IncidentHandler {
	return &IncidentHandler{
		incidentService: incidentService,
	}
}

// Create crée un nouvel incident
func (h *IncidentHandler) Create(c *gin.Context) {
	var req dto.CreateIncidentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	createdByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	incident, err := h.incidentService.Create(req, createdByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.CreatedResponse(c, incident, "Incident créé avec succès")
}

// GetByID récupère un incident par son ID
func (h *IncidentHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	incident, err := h.incidentService.GetByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Incident introuvable")
		return
	}

	utils.SuccessResponse(c, incident, "Incident récupéré avec succès")
}

// GetAll récupère tous les incidents
func (h *IncidentHandler) GetAll(c *gin.Context) {
	incidents, err := h.incidentService.GetAll()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des incidents")
		return
	}

	utils.SuccessResponse(c, incidents, "Incidents récupérés avec succès")
}

// Update met à jour un incident
func (h *IncidentHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	var req dto.UpdateIncidentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	updatedByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	incident, err := h.incidentService.Update(uint(id), req, updatedByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, incident, "Incident mis à jour avec succès")
}

// Qualify qualifie un incident
func (h *IncidentHandler) Qualify(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	var req dto.QualifyIncidentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	qualifiedByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	incident, err := h.incidentService.Qualify(uint(id), req, qualifiedByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, incident, "Incident qualifié avec succès")
}

// Resolve résout un incident
func (h *IncidentHandler) Resolve(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	resolvedByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	incident, err := h.incidentService.Resolve(uint(id), resolvedByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, incident, "Incident résolu avec succès")
}

// Delete supprime un incident
func (h *IncidentHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	err = h.incidentService.Delete(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Incident introuvable")
		return
	}

	utils.SuccessResponse(c, nil, "Incident supprimé avec succès")
}

