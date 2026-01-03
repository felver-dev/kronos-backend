package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/services"
	"github.com/mcicare/itsm-backend/internal/utils"
)

// TicketHandler gère les handlers des tickets
type TicketHandler struct {
	ticketService services.TicketService
}

// NewTicketHandler crée une nouvelle instance de TicketHandler
func NewTicketHandler(ticketService services.TicketService) *TicketHandler {
	return &TicketHandler{
		ticketService: ticketService,
	}
}

// Create crée un nouveau ticket
// @Summary Créer un ticket
// @Description Crée un nouveau ticket dans le système
// @Tags tickets
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateTicketRequest true "Données du ticket"
// @Success 201 {object} dto.TicketDTO
// @Failure 400 {object} utils.Response
// @Router /tickets [post]
func (h *TicketHandler) Create(c *gin.Context) {
	var req dto.CreateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	createdByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	ticket, err := h.ticketService.Create(req, createdByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.CreatedResponse(c, ticket, "Ticket créé avec succès")
}

// GetByID récupère un ticket par son ID
// @Summary Récupérer un ticket
// @Description Récupère les informations d'un ticket par son ID
// @Tags tickets
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID du ticket"
// @Success 200 {object} dto.TicketDTO
// @Failure 404 {object} utils.Response
// @Router /tickets/{id} [get]
func (h *TicketHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	ticket, err := h.ticketService.GetByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Ticket introuvable")
		return
	}

	utils.SuccessResponse(c, ticket, "Ticket récupéré avec succès")
}

// GetAll récupère tous les tickets avec pagination
// @Summary Liste des tickets
// @Description Récupère la liste des tickets avec pagination
// @Tags tickets
// @Security BearerAuth
// @Produce json
// @Param page query int false "Numéro de page" default(1)
// @Param limit query int false "Nombre d'éléments par page" default(20)
// @Success 200 {object} dto.TicketListResponse
// @Failure 500 {object} utils.Response
// @Router /tickets [get]
func (h *TicketHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	response, err := h.ticketService.GetAll(page, limit)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des tickets")
		return
	}

	utils.SuccessResponse(c, response, "Tickets récupérés avec succès")
}

// Update met à jour un ticket
// @Summary Mettre à jour un ticket
// @Description Met à jour les informations d'un ticket
// @Tags tickets
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID du ticket"
// @Param request body dto.UpdateTicketRequest true "Données à mettre à jour"
// @Success 200 {object} dto.TicketDTO
// @Failure 400 {object} utils.Response
// @Router /tickets/{id} [put]
func (h *TicketHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	var req dto.UpdateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	updatedByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	ticket, err := h.ticketService.Update(uint(id), req, updatedByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, ticket, "Ticket mis à jour avec succès")
}

// Assign assigne un ticket à un utilisateur
// @Summary Assigner un ticket
// @Description Assigne un ticket à un technicien
// @Tags tickets
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID du ticket"
// @Param request body dto.AssignTicketRequest true "Données d'assignation"
// @Success 200 {object} dto.TicketDTO
// @Failure 400 {object} utils.Response
// @Router /tickets/{id}/assign [post]
func (h *TicketHandler) Assign(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	var req dto.AssignTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	assignedByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	ticket, err := h.ticketService.Assign(uint(id), req, assignedByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, ticket, "Ticket assigné avec succès")
}

// ChangeStatus change le statut d'un ticket
// @Summary Changer le statut d'un ticket
// @Description Change le statut d'un ticket
// @Tags tickets
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID du ticket"
// @Param request body map[string]string true "Nouveau statut"
// @Success 200 {object} dto.TicketDTO
// @Failure 400 {object} utils.Response
// @Router /tickets/{id}/status [put]
func (h *TicketHandler) ChangeStatus(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	changedByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	ticket, err := h.ticketService.ChangeStatus(uint(id), req.Status, changedByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, ticket, "Statut modifié avec succès")
}

// Close ferme un ticket
// @Summary Fermer un ticket
// @Description Ferme un ticket
// @Tags tickets
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID du ticket"
// @Success 200 {object} dto.TicketDTO
// @Failure 400 {object} utils.Response
// @Router /tickets/{id}/close [post]
func (h *TicketHandler) Close(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	closedByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	ticket, err := h.ticketService.Close(uint(id), closedByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, ticket, "Ticket fermé avec succès")
}

// Delete supprime un ticket
// @Summary Supprimer un ticket
// @Description Supprime un ticket du système
// @Tags tickets
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID du ticket"
// @Success 200 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /tickets/{id} [delete]
func (h *TicketHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	err = h.ticketService.Delete(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Ticket introuvable")
		return
	}

	utils.SuccessResponse(c, nil, "Ticket supprimé avec succès")
}

// AddComment ajoute un commentaire à un ticket
// @Summary Ajouter un commentaire
// @Description Ajoute un commentaire à un ticket
// @Tags tickets
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID du ticket"
// @Param request body dto.CreateTicketCommentRequest true "Commentaire"
// @Success 201 {object} dto.TicketCommentDTO
// @Failure 400 {object} utils.Response
// @Router /tickets/{id}/comments [post]
func (h *TicketHandler) AddComment(c *gin.Context) {
	idParam := c.Param("id")
	ticketID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	var req dto.CreateTicketCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	comment, err := h.ticketService.AddComment(uint(ticketID), req, userID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.CreatedResponse(c, comment, "Commentaire ajouté avec succès")
}

// GetComments récupère les commentaires d'un ticket
// @Summary Récupérer les commentaires
// @Description Récupère tous les commentaires d'un ticket
// @Tags tickets
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID du ticket"
// @Success 200 {array} dto.TicketCommentDTO
// @Failure 404 {object} utils.Response
// @Router /tickets/{id}/comments [get]
func (h *TicketHandler) GetComments(c *gin.Context) {
	idParam := c.Param("id")
	ticketID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	comments, err := h.ticketService.GetComments(uint(ticketID))
	if err != nil {
		utils.NotFoundResponse(c, "Ticket introuvable")
		return
	}

	utils.SuccessResponse(c, comments, "Commentaires récupérés avec succès")
}

