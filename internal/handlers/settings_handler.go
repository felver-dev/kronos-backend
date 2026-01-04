package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/services"
	"github.com/mcicare/itsm-backend/internal/utils"
)

// SettingsHandler gère les handlers des paramètres
type SettingsHandler struct {
	settingsService services.SettingsService
}

// NewSettingsHandler crée une nouvelle instance de SettingsHandler
func NewSettingsHandler(settingsService services.SettingsService) *SettingsHandler {
	return &SettingsHandler{
		settingsService: settingsService,
	}
}

// GetAll récupère tous les paramètres
// @Summary Récupérer les paramètres généraux
// @Description Récupère tous les paramètres généraux du système
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} utils.Response
// @Router /settings [get]
func (h *SettingsHandler) GetAll(c *gin.Context) {
	settings, err := h.settingsService.GetAll()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des paramètres")
		return
	}

	utils.SuccessResponse(c, settings, "Paramètres récupérés avec succès")
}

// Update met à jour les paramètres
// @Summary Mettre à jour les paramètres généraux
// @Description Met à jour les paramètres généraux du système
// @Tags settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.UpdateSettingsRequest true "Paramètres à mettre à jour"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} utils.Response
// @Router /settings [put]
func (h *SettingsHandler) Update(c *gin.Context) {
	var req dto.UpdateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	updatedByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	settings, err := h.settingsService.Update(req, updatedByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, settings, "Paramètres mis à jour avec succès")
}

