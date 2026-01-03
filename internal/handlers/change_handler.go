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
func (h *ChangeHandler) GetAll(c *gin.Context) {
	changes, err := h.changeService.GetAll()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des changements")
		return
	}

	utils.SuccessResponse(c, changes, "Changements récupérés avec succès")
}

// Update met à jour un changement
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

