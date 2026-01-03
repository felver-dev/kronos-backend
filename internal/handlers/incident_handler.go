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
// @Summary Créer un incident
// @Description Crée un nouvel incident dans le système
// @Tags incidents
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateIncidentRequest true "Données de l'incident"
// @Success 201 {object} dto.IncidentDTO
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /incidents [post]
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
// @Summary Récupérer un incident par ID
// @Description Récupère un incident par son identifiant
// @Tags incidents
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID de l'incident"
// @Success 200 {object} dto.IncidentDTO
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /incidents/{id} [get]
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
// @Summary Récupérer tous les incidents
// @Description Récupère la liste de tous les incidents
// @Tags incidents
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {array} dto.IncidentDTO
// @Failure 500 {object} utils.Response
// @Router /incidents [get]
func (h *IncidentHandler) GetAll(c *gin.Context) {
	incidents, err := h.incidentService.GetAll()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des incidents")
		return
	}

	utils.SuccessResponse(c, incidents, "Incidents récupérés avec succès")
}

// Update met à jour un incident
// @Summary Mettre à jour un incident
// @Description Met à jour les informations d'un incident
// @Tags incidents
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID de l'incident"
// @Param request body dto.UpdateIncidentRequest true "Données de mise à jour"
// @Success 200 {object} dto.IncidentDTO
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /incidents/{id} [put]
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
// @Summary Qualifier un incident
// @Description Qualifie un incident (impact et urgence)
// @Tags incidents
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID de l'incident"
// @Param request body dto.QualifyIncidentRequest true "Données de qualification"
// @Success 200 {object} dto.IncidentDTO
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /incidents/{id}/qualify [post]
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
// @Summary Résoudre un incident
// @Description Marque un incident comme résolu
// @Tags incidents
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID de l'incident"
// @Success 200 {object} dto.IncidentDTO
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /incidents/{id}/resolve [post]
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
// @Summary Supprimer un incident
// @Description Supprime un incident du système
// @Tags incidents
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID de l'incident"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /incidents/{id} [delete]
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
