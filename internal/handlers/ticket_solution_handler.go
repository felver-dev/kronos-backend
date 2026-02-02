package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/services"
	"github.com/mcicare/itsm-backend/internal/utils"
)

// TicketSolutionHandler gère les handlers des solutions de tickets
type TicketSolutionHandler struct {
	solutionService services.TicketSolutionService
}

// NewTicketSolutionHandler crée une nouvelle instance de TicketSolutionHandler
func NewTicketSolutionHandler(solutionService services.TicketSolutionService) *TicketSolutionHandler {
	return &TicketSolutionHandler{
		solutionService: solutionService,
	}
}

// Create crée une nouvelle solution pour un ticket
// @Summary Créer une solution pour un ticket
// @Description Crée une nouvelle solution documentée pour un ticket cloturé
// @Tags tickets
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID du ticket"
// @Param request body dto.CreateTicketSolutionRequest true "Données de la solution"
// @Success 201 {object} dto.TicketSolutionDTO
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /tickets/{id}/solutions [post]
func (h *TicketSolutionHandler) Create(c *gin.Context) {
	ticketIDParam := c.Param("id")
	ticketID, err := strconv.ParseUint(ticketIDParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID ticket invalide")
		return
	}

	var req dto.CreateTicketSolutionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	createdByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	solution, err := h.solutionService.Create(uint(ticketID), req, createdByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.CreatedResponse(c, solution, "Solution créée avec succès")
}

// GetByID récupère une solution par son ID
// @Summary Récupérer une solution par ID
// @Description Récupère une solution de ticket par son identifiant
// @Tags tickets
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID de la solution"
// @Success 200 {object} dto.TicketSolutionDTO
// @Failure 404 {object} utils.Response
// @Router /tickets/solutions/{id} [get]
func (h *TicketSolutionHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	solution, err := h.solutionService.GetByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Solution introuvable")
		return
	}

	utils.SuccessResponse(c, solution, "Solution récupérée avec succès")
}

// GetByTicketID récupère toutes les solutions d'un ticket
// @Summary Récupérer les solutions d'un ticket
// @Description Récupère toutes les solutions documentées pour un ticket
// @Tags tickets
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID du ticket"
// @Success 200 {array} dto.TicketSolutionDTO
// @Failure 400 {object} utils.Response
// @Router /tickets/{id}/solutions [get]
func (h *TicketSolutionHandler) GetByTicketID(c *gin.Context) {
	ticketIDParam := c.Param("id")
	ticketID, err := strconv.ParseUint(ticketIDParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID ticket invalide")
		return
	}

	solutions, err := h.solutionService.GetByTicketID(uint(ticketID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, solutions, "Solutions récupérées avec succès")
}

// Update met à jour une solution
// @Summary Mettre à jour une solution
// @Description Met à jour une solution documentée
// @Tags tickets
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID de la solution"
// @Param request body dto.UpdateTicketSolutionRequest true "Données à mettre à jour"
// @Success 200 {object} dto.TicketSolutionDTO
// @Failure 400 {object} utils.Response
// @Router /tickets/solutions/{id} [put]
func (h *TicketSolutionHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	var req dto.UpdateTicketSolutionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	updatedByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	solution, err := h.solutionService.Update(uint(id), req, updatedByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, solution, "Solution mise à jour avec succès")
}

// Delete supprime une solution
// @Summary Supprimer une solution
// @Description Supprime une solution documentée
// @Tags tickets
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID de la solution"
// @Success 200 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /tickets/solutions/{id} [delete]
func (h *TicketSolutionHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	err = h.solutionService.Delete(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Solution introuvable")
		return
	}

	utils.SuccessResponse(c, nil, "Solution supprimée avec succès")
}

// PublishToKB publie une solution dans la base de connaissances
// @Summary Publier une solution dans la base de connaissances
// @Description Publie une solution documentée dans la base de connaissances
// @Tags tickets
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID de la solution"
// @Param request body dto.PublishSolutionToKBRequest true "Données de publication"
// @Success 201 {object} dto.KnowledgeArticleDTO
// @Failure 400 {object} utils.Response
// @Router /tickets/solutions/{id}/publish-to-kb [post]
func (h *TicketSolutionHandler) PublishToKB(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	var req dto.PublishSolutionToKBRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	publishedByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	article, err := h.solutionService.PublishToKB(uint(id), req, publishedByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.CreatedResponse(c, article, "Solution publiée dans la base de connaissances avec succès")
}
