package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/services"
	"github.com/mcicare/itsm-backend/internal/utils"
)

// TicketCategoryHandler gère les handlers des catégories de tickets
type TicketCategoryHandler struct {
	ticketCategoryService services.TicketCategoryService
}

// NewTicketCategoryHandler crée une nouvelle instance de TicketCategoryHandler
func NewTicketCategoryHandler(ticketCategoryService services.TicketCategoryService) *TicketCategoryHandler {
	return &TicketCategoryHandler{
		ticketCategoryService: ticketCategoryService,
	}
}

// GetAll récupère toutes les catégories
// @Summary Récupérer toutes les catégories de tickets
// @Description Récupère la liste de toutes les catégories de tickets
// @Tags ticket-categories
// @Security BearerAuth
// @Produce json
// @Param active query bool false "Récupérer uniquement les catégories actives"
// @Success 200 {array} dto.TicketCategoryDTO
// @Failure 500 {object} utils.Response
// @Router /tickets/categories [get]
func (h *TicketCategoryHandler) GetAll(c *gin.Context) {
	activeOnly := c.Query("active") == "true"

	var categories []dto.TicketCategoryDTO
	var err error

	if activeOnly {
		categories, err = h.ticketCategoryService.GetActive()
	} else {
		categories, err = h.ticketCategoryService.GetAll()
	}

	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des catégories")
		return
	}

	utils.SuccessResponse(c, categories, "Catégories récupérées avec succès")
}

// GetByID récupère une catégorie par son ID
// @Summary Récupérer une catégorie par ID
// @Description Récupère une catégorie de ticket par son identifiant
// @Tags ticket-categories
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID de la catégorie"
// @Success 200 {object} dto.TicketCategoryDTO
// @Failure 404 {object} utils.Response
// @Router /tickets/categories/{id} [get]
func (h *TicketCategoryHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	category, err := h.ticketCategoryService.GetByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Catégorie introuvable")
		return
	}

	utils.SuccessResponse(c, category, "Catégorie récupérée avec succès")
}

// GetBySlug récupère une catégorie par son slug
// @Summary Récupérer une catégorie par slug
// @Description Récupère une catégorie de ticket par son slug
// @Tags ticket-categories
// @Security BearerAuth
// @Produce json
// @Param slug path string true "Slug de la catégorie"
// @Success 200 {object} dto.TicketCategoryDTO
// @Failure 404 {object} utils.Response
// @Router /tickets/categories/slug/{slug} [get]
func (h *TicketCategoryHandler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")

	category, err := h.ticketCategoryService.GetBySlug(slug)
	if err != nil {
		utils.NotFoundResponse(c, "Catégorie introuvable")
		return
	}

	utils.SuccessResponse(c, category, "Catégorie récupérée avec succès")
}

// Create crée une nouvelle catégorie
// @Summary Créer une catégorie
// @Description Crée une nouvelle catégorie de ticket
// @Tags ticket-categories
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateTicketCategoryRequest true "Données de la catégorie"
// @Success 201 {object} dto.TicketCategoryDTO
// @Failure 400 {object} utils.Response
// @Router /tickets/categories [post]
func (h *TicketCategoryHandler) Create(c *gin.Context) {
	var req dto.CreateTicketCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	category, err := h.ticketCategoryService.Create(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.CreatedResponse(c, category, "Catégorie créée avec succès")
}

// Update met à jour une catégorie
// @Summary Mettre à jour une catégorie
// @Description Met à jour les informations d'une catégorie de ticket
// @Tags ticket-categories
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID de la catégorie"
// @Param request body dto.UpdateTicketCategoryRequest true "Données à mettre à jour"
// @Success 200 {object} dto.TicketCategoryDTO
// @Failure 400 {object} utils.Response
// @Router /tickets/categories/{id} [put]
func (h *TicketCategoryHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	var req dto.UpdateTicketCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	category, err := h.ticketCategoryService.Update(uint(id), req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, category, "Catégorie mise à jour avec succès")
}

// Delete supprime une catégorie
// @Summary Supprimer une catégorie
// @Description Supprime une catégorie de ticket
// @Tags ticket-categories
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID de la catégorie"
// @Success 200 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /tickets/categories/{id} [delete]
func (h *TicketCategoryHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	err = h.ticketCategoryService.Delete(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, nil, "Catégorie supprimée avec succès")
}
