package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/services"
	"github.com/mcicare/itsm-backend/internal/utils"
)

// KnowledgeArticleHandler gère les handlers de la base de connaissances
type KnowledgeArticleHandler struct {
	knowledgeArticleService services.KnowledgeArticleService
}

// NewKnowledgeArticleHandler crée une nouvelle instance de KnowledgeArticleHandler
func NewKnowledgeArticleHandler(knowledgeArticleService services.KnowledgeArticleService) *KnowledgeArticleHandler {
	return &KnowledgeArticleHandler{
		knowledgeArticleService: knowledgeArticleService,
	}
}

// Create crée un nouvel article
func (h *KnowledgeArticleHandler) Create(c *gin.Context) {
	var req dto.CreateKnowledgeArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	authorID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	article, err := h.knowledgeArticleService.Create(req, authorID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.CreatedResponse(c, article, "Article créé avec succès")
}

// GetByID récupère un article par son ID
func (h *KnowledgeArticleHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	article, err := h.knowledgeArticleService.GetByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Article introuvable")
		return
	}

	utils.SuccessResponse(c, article, "Article récupéré avec succès")
}

// GetPublished récupère les articles publiés
func (h *KnowledgeArticleHandler) GetPublished(c *gin.Context) {
	articles, err := h.knowledgeArticleService.GetPublished()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des articles")
		return
	}

	utils.SuccessResponse(c, articles, "Articles publiés récupérés avec succès")
}

// Search recherche des articles
func (h *KnowledgeArticleHandler) Search(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		utils.BadRequestResponse(c, "Paramètre de recherche manquant")
		return
	}

	results, err := h.knowledgeArticleService.Search(query)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la recherche")
		return
	}

	utils.SuccessResponse(c, results, "Résultats de recherche récupérés avec succès")
}

// Delete supprime un article
func (h *KnowledgeArticleHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	err = h.knowledgeArticleService.Delete(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Article introuvable")
		return
	}

	utils.SuccessResponse(c, nil, "Article supprimé avec succès")
}
