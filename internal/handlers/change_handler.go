package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/services"
	"github.com/mcicare/itsm-backend/internal/utils"
)

// ChangeHandler gère les handlers des changements
type ChangeHandler struct {
	changeService services.ChangeService
}

// NewChangeHandler crée une nouvelle instance de ChangeHandler
func NewChangeHandler(changeService services.ChangeService) *ChangeHandler {
	return &ChangeHandler{
		changeService: changeService,
	}
}

// Create crée un nouveau changement
// @Summary Créer un changement
// @Description Crée un nouveau changement dans le système
// @Tags changes
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateChangeRequest true "Données du changement"
// @Success 201 {object} dto.ChangeDTO
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /changes [post]
func (h *ChangeHandler) Create(c *gin.Context) {
	var req dto.CreateChangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	createdByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	change, err := h.changeService.Create(req, createdByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.CreatedResponse(c, change, "Changement créé avec succès")
}

// GetByID récupère un changement par son ID
// @Summary Récupérer un changement par ID
// @Description Récupère un changement par son identifiant
// @Tags changes
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID du changement"
// @Success 200 {object} dto.ChangeDTO
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /changes/{id} [get]
func (h *ChangeHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	change, err := h.changeService.GetByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Changement introuvable")
		return
	}

	utils.SuccessResponse(c, change, "Changement récupéré avec succès")
}

// GetAll récupère tous les changements
// @Summary Récupérer tous les changements
// @Description Récupère la liste de tous les changements
// @Tags changes
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {array} dto.ChangeDTO
// @Failure 500 {object} utils.Response
// @Router /changes [get]
func (h *ChangeHandler) GetAll(c *gin.Context) {
	changes, err := h.changeService.GetAll()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des changements")
		return
	}

	utils.SuccessResponse(c, changes, "Changements récupérés avec succès")
}

// Update met à jour un changement
// @Summary Mettre à jour un changement
// @Description Met à jour les informations d'un changement
// @Tags changes
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID du changement"
// @Param request body dto.UpdateChangeRequest true "Données de mise à jour"
// @Success 200 {object} dto.ChangeDTO
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /changes/{id} [put]
func (h *ChangeHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	var req dto.UpdateChangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	updatedByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	change, err := h.changeService.Update(uint(id), req, updatedByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, change, "Changement mis à jour avec succès")
}

// RecordResult enregistre le résultat d'un changement
// @Summary Enregistrer le résultat d'un changement
// @Description Enregistre le résultat d'un changement (succès/échec)
// @Tags changes
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID du changement"
// @Param request body dto.RecordChangeResultRequest true "Résultat du changement"
// @Success 200 {object} dto.ChangeDTO
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /changes/{id}/result [post]
func (h *ChangeHandler) RecordResult(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	var req dto.RecordChangeResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	recordedByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	change, err := h.changeService.RecordResult(uint(id), req, recordedByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, change, "Résultat enregistré avec succès")
}

// Delete supprime un changement
// @Summary Supprimer un changement
// @Description Supprime un changement du système
// @Tags changes
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID du changement"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /changes/{id} [delete]
func (h *ChangeHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	err = h.changeService.Delete(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Changement introuvable")
		return
	}

	utils.SuccessResponse(c, nil, "Changement supprimé avec succès")
}

// UpdateRisk met à jour le risque d'un changement
// @Summary Mettre à jour le risque
// @Description Met à jour le niveau de risque d'un changement
// @Tags changes
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID du changement"
// @Param request body dto.UpdateRiskRequest true "Données de risque"
// @Success 200 {object} dto.ChangeDTO
// @Failure 400 {object} utils.Response
// @Router /changes/{id}/risk [put]
func (h *ChangeHandler) UpdateRisk(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	var req dto.UpdateRiskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	updatedByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	change, err := h.changeService.UpdateRisk(uint(id), req, updatedByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, change, "Risque mis à jour avec succès")
}

// AssignResponsible assigne un responsable à un changement
// @Summary Assigner un responsable
// @Description Assigne un responsable à un changement
// @Tags changes
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID du changement"
// @Param request body dto.AssignResponsibleRequest true "ID du responsable"
// @Success 200 {object} dto.ChangeDTO
// @Failure 400 {object} utils.Response
// @Router /changes/{id}/assign-responsible [post]
func (h *ChangeHandler) AssignResponsible(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	var req dto.AssignResponsibleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	assignedByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	change, err := h.changeService.AssignResponsible(uint(id), req, assignedByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, change, "Responsable assigné avec succès")
}

// GetResult récupère le résultat d'un changement
// @Summary Récupérer le résultat d'un changement
// @Description Récupère le résultat post-changement d'un changement
// @Tags changes
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID du changement"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} utils.Response
// @Router /changes/{id}/result [get]
func (h *ChangeHandler) GetResult(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	change, err := h.changeService.GetByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Changement introuvable")
		return
	}

	response := map[string]interface{}{
		"result":            change.Result,
		"description":      change.ResultDescription,
		"date":              change.ResultDate,
		"has_result":        change.Result != "",
	}

	utils.SuccessResponse(c, response, "Résultat récupéré avec succès")
}

// GetByRisk récupère les changements par niveau de risque
// @Summary Récupérer les changements par risque
// @Description Récupère les changements filtrés par niveau de risque (low, medium, high, critical)
// @Tags changes
// @Security BearerAuth
// @Produce json
// @Param riskLevel path string true "Niveau de risque (low, medium, high, critical)"
// @Success 200 {array} dto.ChangeDTO
// @Failure 400 {object} utils.Response
// @Router /changes/by-risk/{riskLevel} [get]
func (h *ChangeHandler) GetByRisk(c *gin.Context) {
	riskLevel := c.Param("riskLevel")

	changes, err := h.changeService.GetByRisk(riskLevel)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, changes, "Changements récupérés avec succès")
}

// GetByResponsible récupère les changements par responsable
// @Summary Récupérer les changements par responsable
// @Description Récupère les changements assignés à un responsable spécifique
// @Tags changes
// @Security BearerAuth
// @Produce json
// @Param userId path int true "ID du responsable"
// @Success 200 {array} dto.ChangeDTO
// @Failure 400 {object} utils.Response
// @Router /changes/by-responsible/{userId} [get]
func (h *ChangeHandler) GetByResponsible(c *gin.Context) {
	userIDParam := c.Param("userId")
	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID utilisateur invalide")
		return
	}

	changes, err := h.changeService.GetByResponsible(uint(userID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, changes, "Changements récupérés avec succès")
}
