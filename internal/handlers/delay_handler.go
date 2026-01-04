package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/services"
	"github.com/mcicare/itsm-backend/internal/utils"
)

// DelayHandler gère les handlers des retards
type DelayHandler struct {
	delayService services.DelayService
}

// NewDelayHandler crée une nouvelle instance de DelayHandler
func NewDelayHandler(delayService services.DelayService) *DelayHandler {
	return &DelayHandler{
		delayService: delayService,
	}
}

// GetByID récupère un retard par son ID
// @Summary Récupérer un retard par ID
// @Description Récupère un retard par son identifiant
// @Tags delays
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID du retard"
// @Success 200 {object} dto.DelayDTO
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /delays/{id} [get]
func (h *DelayHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	delay, err := h.delayService.GetByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Retard introuvable")
		return
	}

	utils.SuccessResponse(c, delay, "Retard récupéré avec succès")
}

// GetAll récupère tous les retards
// @Summary Récupérer tous les retards
// @Description Récupère la liste de tous les retards
// @Tags delays
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {array} dto.DelayDTO
// @Failure 500 {object} utils.Response
// @Router /delays [get]
func (h *DelayHandler) GetAll(c *gin.Context) {
	delays, err := h.delayService.GetAll()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des retards")
		return
	}

	utils.SuccessResponse(c, delays, "Retards récupérés avec succès")
}

// CreateJustification crée une justification pour un retard
// @Summary Créer une justification de retard
// @Description Crée une justification pour un retard
// @Tags delays
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param delay_id path int true "ID du retard"
// @Param request body dto.CreateDelayJustificationRequest true "Données de la justification"
// @Success 201 {object} dto.DelayJustificationDTO
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /delays/{delay_id}/justifications [post]
func (h *DelayHandler) CreateJustification(c *gin.Context) {
	delayIDParam := c.Param("id")
	delayID, err := strconv.ParseUint(delayIDParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	var req dto.CreateDelayJustificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	justification, err := h.delayService.CreateJustification(uint(delayID), req, userID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.CreatedResponse(c, justification, "Justification créée avec succès")
}

// ValidateJustification valide une justification
// @Summary Valider une justification de retard
// @Description Valide une justification de retard
// @Tags delays
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID de la justification"
// @Param request body dto.ValidateDelayJustificationRequest true "Données de validation"
// @Success 200 {object} dto.DelayJustificationDTO
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /delays/justifications/{id}/validate [post]
func (h *DelayHandler) ValidateJustification(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	var req dto.ValidateDelayJustificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	validatedByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	justification, err := h.delayService.ValidateJustification(uint(id), req, validatedByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, justification, "Justification validée avec succès")
}

// GetJustificationByDelayID récupère la justification d'un retard
// @Summary Récupérer la justification d'un retard
// @Description Récupère la justification d'un retard par l'ID du retard
// @Tags delays
// @Security BearerAuth
// @Produce json
// @Param delayId path int true "ID du retard"
// @Success 200 {object} dto.DelayJustificationDTO
// @Failure 404 {object} utils.Response
// @Router /delays/{delayId}/justification [get]
func (h *DelayHandler) GetJustificationByDelayID(c *gin.Context) {
	delayIDParam := c.Param("id")
	delayID, err := strconv.ParseUint(delayIDParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	justification, err := h.delayService.GetJustificationByDelayID(uint(delayID))
	if err != nil {
		utils.NotFoundResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, justification, "Justification récupérée avec succès")
}

// UpdateJustification met à jour une justification
// @Summary Mettre à jour une justification de retard
// @Description Met à jour une justification de retard (avant validation)
// @Tags delays
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param delayId path int true "ID du retard"
// @Param request body dto.UpdateDelayJustificationRequest true "Nouvelle justification"
// @Success 200 {object} dto.DelayJustificationDTO
// @Failure 400 {object} utils.Response
// @Router /delays/{delayId}/justification [put]
func (h *DelayHandler) UpdateJustification(c *gin.Context) {
	delayIDParam := c.Param("id")
	delayID, err := strconv.ParseUint(delayIDParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	// Récupérer la justification pour obtenir son ID
	justification, err := h.delayService.GetJustificationByDelayID(uint(delayID))
	if err != nil {
		utils.NotFoundResponse(c, "Justification introuvable")
		return
	}

	var req dto.UpdateDelayJustificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	updatedJustification, err := h.delayService.UpdateJustification(justification.ID, req, userID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, updatedJustification, "Justification mise à jour avec succès")
}

// DeleteJustification supprime une justification
// @Summary Supprimer une justification de retard
// @Description Supprime une justification de retard (seulement si en attente)
// @Tags delays
// @Security BearerAuth
// @Produce json
// @Param delayId path int true "ID du retard"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /delays/{delayId}/justification [delete]
func (h *DelayHandler) DeleteJustification(c *gin.Context) {
	delayIDParam := c.Param("id")
	delayID, err := strconv.ParseUint(delayIDParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	err = h.delayService.DeleteJustification(uint(delayID), userID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, nil, "Justification supprimée avec succès")
}

// GetJustificationByTicketID récupère la justification d'un ticket
// @Summary Récupérer la justification d'un ticket
// @Description Récupère la justification de retard d'un ticket
// @Tags tickets
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID du ticket"
// @Success 200 {object} dto.DelayJustificationDTO
// @Failure 404 {object} utils.Response
// @Router /tickets/{ticketId}/delay-justification [get]
func (h *DelayHandler) GetJustificationByTicketID(c *gin.Context) {
	ticketIDParam := c.Param("id")
	ticketID, err := strconv.ParseUint(ticketIDParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID de ticket invalide")
		return
	}

	justification, err := h.delayService.GetJustificationByTicketID(uint(ticketID))
	if err != nil {
		utils.NotFoundResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, justification, "Justification récupérée avec succès")
}

// GetJustificationsByUserID récupère les justifications d'un utilisateur
// @Summary Récupérer les justifications d'un utilisateur
// @Description Récupère toutes les justifications de retards d'un utilisateur
// @Tags users
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID de l'utilisateur"
// @Success 200 {array} dto.DelayJustificationDTO
// @Failure 500 {object} utils.Response
// @Router /users/{userId}/delay-justifications [get]
func (h *DelayHandler) GetJustificationsByUserID(c *gin.Context) {
	userIDParam := c.Param("id")
	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID utilisateur invalide")
		return
	}

	justifications, err := h.delayService.GetJustificationsByUserID(uint(userID))
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des justifications")
		return
	}

	utils.SuccessResponse(c, justifications, "Justifications récupérées avec succès")
}

// RejectJustification rejette une justification
// @Summary Rejeter une justification de retard
// @Description Rejette une justification de retard
// @Tags delays
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param delayId path int true "ID du retard"
// @Param request body dto.ValidateDelayJustificationRequest true "Commentaire de rejet"
// @Success 200 {object} dto.DelayJustificationDTO
// @Failure 400 {object} utils.Response
// @Router /delays/{delayId}/justification/reject [post]
func (h *DelayHandler) RejectJustification(c *gin.Context) {
	delayIDParam := c.Param("id")
	delayID, err := strconv.ParseUint(delayIDParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	var req dto.ValidateDelayJustificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	rejectedByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	justification, err := h.delayService.RejectJustification(uint(delayID), req, rejectedByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, justification, "Justification rejetée avec succès")
}

// GetValidatedJustifications récupère les justifications validées
// @Summary Récupérer les justifications validées
// @Description Récupère toutes les justifications validées
// @Tags delays
// @Security BearerAuth
// @Produce json
// @Success 200 {array} dto.DelayJustificationDTO
// @Failure 500 {object} utils.Response
// @Router /delays/justifications/validated [get]
func (h *DelayHandler) GetValidatedJustifications(c *gin.Context) {
	justifications, err := h.delayService.GetValidatedJustifications()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des justifications validées")
		return
	}

	utils.SuccessResponse(c, justifications, "Justifications validées récupérées avec succès")
}

// GetRejectedJustifications récupère les justifications rejetées
// @Summary Récupérer les justifications rejetées
// @Description Récupère toutes les justifications rejetées
// @Tags delays
// @Security BearerAuth
// @Produce json
// @Success 200 {array} dto.DelayJustificationDTO
// @Failure 500 {object} utils.Response
// @Router /delays/justifications/rejected [get]
func (h *DelayHandler) GetRejectedJustifications(c *gin.Context) {
	justifications, err := h.delayService.GetRejectedJustifications()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des justifications rejetées")
		return
	}

	utils.SuccessResponse(c, justifications, "Justifications rejetées récupérées avec succès")
}

// GetJustificationsHistory récupère l'historique des justifications
// @Summary Historique des justifications
// @Description Récupère l'historique de toutes les justifications
// @Tags delays
// @Security BearerAuth
// @Produce json
// @Success 200 {array} dto.DelayJustificationDTO
// @Failure 500 {object} utils.Response
// @Router /delays/justifications/history [get]
func (h *DelayHandler) GetJustificationsHistory(c *gin.Context) {
	justifications, err := h.delayService.GetJustificationsHistory()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération de l'historique")
		return
	}

	utils.SuccessResponse(c, justifications, "Historique des justifications récupéré avec succès")
}
