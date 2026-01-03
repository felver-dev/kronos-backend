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
func (h *AssetHandler) GetAll(c *gin.Context) {
	assets, err := h.assetService.GetAll()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des actifs")
		return
	}

	utils.SuccessResponse(c, assets, "Actifs récupérés avec succès")
}

// Assign assigne un actif à un utilisateur
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
