package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/services"
	"github.com/mcicare/itsm-backend/internal/utils"
)

// BackupHandler gère les handlers des sauvegardes
type BackupHandler struct {
	backupService services.BackupService
}

// NewBackupHandler crée une nouvelle instance de BackupHandler
func NewBackupHandler(backupService services.BackupService) *BackupHandler {
	return &BackupHandler{
		backupService: backupService,
	}
}

// GetConfiguration récupère la configuration de sauvegarde
// @Summary Récupérer la configuration de sauvegarde
// @Description Récupère la configuration de sauvegarde du système
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dto.BackupConfigurationDTO
// @Failure 500 {object} utils.Response
// @Router /settings/backup [get]
func (h *BackupHandler) GetConfiguration(c *gin.Context) {
	config, err := h.backupService.GetConfiguration()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération de la configuration")
		return
	}

	utils.SuccessResponse(c, config, "Configuration de sauvegarde récupérée avec succès")
}

// UpdateConfiguration met à jour la configuration de sauvegarde
// @Summary Mettre à jour la configuration de sauvegarde
// @Description Met à jour la configuration de sauvegarde du système
// @Tags settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.BackupConfigurationDTO true "Configuration de sauvegarde"
// @Success 200 {object} dto.BackupConfigurationDTO
// @Failure 400 {object} utils.Response
// @Router /settings/backup [put]
func (h *BackupHandler) UpdateConfiguration(c *gin.Context) {
	var req dto.BackupConfigurationDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	updatedByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	config, err := h.backupService.UpdateConfiguration(req, updatedByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, config, "Configuration de sauvegarde mise à jour avec succès")
}

// ExecuteBackup exécute une sauvegarde manuelle
// @Summary Exécuter une sauvegarde manuelle
// @Description Exécute une sauvegarde manuelle du système
// @Tags settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.ExecuteBackupRequest false "Type de sauvegarde (optionnel)"
// @Success 200 {object} dto.BackupExecutionResponse
// @Failure 400 {object} utils.Response
// @Router /settings/backup/execute [post]
func (h *BackupHandler) ExecuteBackup(c *gin.Context) {
	var req dto.ExecuteBackupRequest
	_ = c.ShouldBindJSON(&req) // Optionnel

	executedByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	backupType := req.Type
	if backupType == "" {
		backupType = "full"
	}

	response, err := h.backupService.ExecuteBackup(backupType, executedByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, response, "Sauvegarde démarrée avec succès")
}

