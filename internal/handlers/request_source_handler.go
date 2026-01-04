package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/services"
	"github.com/mcicare/itsm-backend/internal/utils"
)

// RequestSourceHandler gère les handlers des sources de demande
type RequestSourceHandler struct {
	requestSourceService services.RequestSourceService
}

// NewRequestSourceHandler crée une nouvelle instance de RequestSourceHandler
func NewRequestSourceHandler(requestSourceService services.RequestSourceService) *RequestSourceHandler {
	return &RequestSourceHandler{
		requestSourceService: requestSourceService,
	}
}

// GetAll récupère toutes les sources de demande
// @Summary Récupérer les sources de demande
// @Description Récupère la liste de toutes les sources de demande configurées
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Success 200 {array} dto.RequestSourceDTO
// @Failure 500 {object} utils.Response
// @Router /settings/sources [get]
func (h *RequestSourceHandler) GetAll(c *gin.Context) {
	sources, err := h.requestSourceService.GetAll()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des sources")
		return
	}

	utils.SuccessResponse(c, sources, "Sources récupérées avec succès")
}

// GetByID récupère une source par son ID
// @Summary Récupérer une source par ID
// @Description Récupère une source de demande par son identifiant
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID de la source"
// @Success 200 {object} dto.RequestSourceDTO
// @Failure 404 {object} utils.Response
// @Router /settings/sources/{id} [get]
func (h *RequestSourceHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	source, err := h.requestSourceService.GetByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Source introuvable")
		return
	}

	utils.SuccessResponse(c, source, "Source récupérée avec succès")
}

// Create crée une nouvelle source de demande
// @Summary Créer une source de demande
// @Description Crée une nouvelle source de demande
// @Tags settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateRequestSourceRequest true "Données de la source"
// @Success 201 {object} dto.RequestSourceDTO
// @Failure 400 {object} utils.Response
// @Router /settings/sources [post]
func (h *RequestSourceHandler) Create(c *gin.Context) {
	var req dto.CreateRequestSourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	createdByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	source, err := h.requestSourceService.Create(req, createdByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.CreatedResponse(c, source, "Source créée avec succès")
}

// Update met à jour une source de demande
// @Summary Mettre à jour une source de demande
// @Description Met à jour une source de demande
// @Tags settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID de la source"
// @Param request body dto.UpdateRequestSourceRequest true "Données à mettre à jour"
// @Success 200 {object} dto.RequestSourceDTO
// @Failure 400 {object} utils.Response
// @Router /settings/sources/{id} [put]
func (h *RequestSourceHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	var req dto.UpdateRequestSourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	updatedByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	source, err := h.requestSourceService.Update(uint(id), req, updatedByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, source, "Source mise à jour avec succès")
}

// Delete supprime une source de demande
// @Summary Supprimer une source de demande
// @Description Supprime une source de demande
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID de la source"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /settings/sources/{id} [delete]
func (h *RequestSourceHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	err = h.requestSourceService.Delete(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, nil, "Source supprimée avec succès")
}

