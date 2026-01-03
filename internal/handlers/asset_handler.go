package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/services"
	"github.com/mcicare/itsm-backend/internal/utils"
)

// AssetHandler gère les handlers des actifs IT
type AssetHandler struct {
	assetService services.AssetService
}

// NewAssetHandler crée une nouvelle instance de AssetHandler
func NewAssetHandler(assetService services.AssetService) *AssetHandler {
	return &AssetHandler{
		assetService: assetService,
	}
}

// Create crée un nouvel actif
// @Summary Créer un actif IT
// @Description Crée un nouvel actif IT dans le système
// @Tags assets
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateAssetRequest true "Données de l'actif"
// @Success 201 {object} dto.AssetDTO
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /assets [post]
func (h *AssetHandler) Create(c *gin.Context) {
	var req dto.CreateAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	createdByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	asset, err := h.assetService.Create(req, createdByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.CreatedResponse(c, asset, "Actif créé avec succès")
}

// GetByID récupère un actif par son ID
// @Summary Récupérer un actif par ID
// @Description Récupère un actif IT par son identifiant
// @Tags assets
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID de l'actif"
// @Success 200 {object} dto.AssetDTO
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /assets/{id} [get]
func (h *AssetHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	asset, err := h.assetService.GetByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Actif introuvable")
		return
	}

	utils.SuccessResponse(c, asset, "Actif récupéré avec succès")
}

// GetAll récupère tous les actifs
// @Summary Récupérer tous les actifs
// @Description Récupère la liste de tous les actifs IT
// @Tags assets
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {array} dto.AssetDTO
// @Failure 500 {object} utils.Response
// @Router /assets [get]
func (h *AssetHandler) GetAll(c *gin.Context) {
	assets, err := h.assetService.GetAll()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des actifs")
		return
	}

	utils.SuccessResponse(c, assets, "Actifs récupérés avec succès")
}

// Assign assigne un actif à un utilisateur
// @Summary Assigner un actif à un utilisateur
// @Description Assigne un actif IT à un utilisateur
// @Tags assets
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID de l'actif"
// @Param request body dto.AssignAssetRequest true "Données d'assignation"
// @Success 200 {object} dto.AssetDTO
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /assets/{id}/assign [post]
func (h *AssetHandler) Assign(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	var req dto.AssignAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	assignedByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	asset, err := h.assetService.Assign(uint(id), req, assignedByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, asset, "Actif assigné avec succès")
}

// Delete supprime un actif
// @Summary Supprimer un actif
// @Description Supprime un actif IT du système
// @Tags assets
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID de l'actif"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /assets/{id} [delete]
func (h *AssetHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	err = h.assetService.Delete(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Actif introuvable")
		return
	}

	utils.SuccessResponse(c, nil, "Actif supprimé avec succès")
}
