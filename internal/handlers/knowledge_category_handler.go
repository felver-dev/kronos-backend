package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/services"
	"github.com/mcicare/itsm-backend/internal/utils"
)

// KnowledgeCategoryHandler gère les handlers des catégories de la base de connaissances
type KnowledgeCategoryHandler struct {
	knowledgeCategoryService services.KnowledgeCategoryService
}

// NewKnowledgeCategoryHandler crée une nouvelle instance de KnowledgeCategoryHandler
func NewKnowledgeCategoryHandler(knowledgeCategoryService services.KnowledgeCategoryService) *KnowledgeCategoryHandler {
	return &KnowledgeCategoryHandler{
		knowledgeCategoryService: knowledgeCategoryService,
	}
}

// GetAll récupère toutes les catégories
// @Summary Récupérer toutes les catégories
// @Description Récupère la liste de toutes les catégories de la base de connaissances
// @Tags knowledge-base
// @Security BearerAuth
// @Produce json
// @Success 200 {array} dto.KnowledgeCategoryDTO
// @Failure 500 {object} utils.Response
// @Router /knowledge-base/categories [get]
func (h *KnowledgeCategoryHandler) GetAll(c *gin.Context) {
	categories, err := h.knowledgeCategoryService.GetAll()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des catégories")
		return
	}

	utils.SuccessResponse(c, categories, "Catégories récupérées avec succès")
}

// GetByID récupère une catégorie par son ID
// @Summary Récupérer une catégorie par ID
// @Description Récupère une catégorie de la base de connaissances par son identifiant
// @Tags knowledge-base
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID de la catégorie"
// @Success 200 {object} dto.KnowledgeCategoryDTO
// @Failure 404 {object} utils.Response
// @Router /knowledge-base/categories/{id} [get]
func (h *KnowledgeCategoryHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	category, err := h.knowledgeCategoryService.GetByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Catégorie introuvable")
		return
	}

	utils.SuccessResponse(c, category, "Catégorie récupérée avec succès")
}

// Create crée une nouvelle catégorie
// @Summary Créer une catégorie
// @Description Crée une nouvelle catégorie dans la base de connaissances
// @Tags knowledge-base
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateKnowledgeCategoryRequest true "Données de la catégorie"
// @Success 201 {object} dto.KnowledgeCategoryDTO
// @Failure 400 {object} utils.Response
// @Router /knowledge-base/categories [post]
func (h *KnowledgeCategoryHandler) Create(c *gin.Context) {
	var req dto.CreateKnowledgeCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	createdByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	category, err := h.knowledgeCategoryService.Create(req, createdByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.CreatedResponse(c, category, "Catégorie créée avec succès")
}

// Update met à jour une catégorie
// @Summary Mettre à jour une catégorie
// @Description Met à jour les informations d'une catégorie
// @Tags knowledge-base
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID de la catégorie"
// @Param request body dto.UpdateKnowledgeCategoryRequest true "Données à mettre à jour"
// @Success 200 {object} dto.KnowledgeCategoryDTO
// @Failure 400 {object} utils.Response
// @Router /knowledge-base/categories/{id} [put]
func (h *KnowledgeCategoryHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	var req dto.UpdateKnowledgeCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	updatedByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	category, err := h.knowledgeCategoryService.Update(uint(id), req, updatedByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, category, "Catégorie mise à jour avec succès")
}

// Delete supprime une catégorie
// @Summary Supprimer une catégorie
// @Description Supprime une catégorie de la base de connaissances
// @Tags knowledge-base
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID de la catégorie"
// @Success 200 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /knowledge-base/categories/{id} [delete]
func (h *KnowledgeCategoryHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	err = h.knowledgeCategoryService.Delete(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Catégorie introuvable")
		return
	}

	utils.SuccessResponse(c, nil, "Catégorie supprimée avec succès")
}

