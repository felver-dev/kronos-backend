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
// @Summary Créer un article de base de connaissances
// @Description Crée un nouvel article dans la base de connaissances
// @Tags knowledge-base
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateKnowledgeArticleRequest true "Données de l'article"
// @Success 201 {object} dto.KnowledgeArticleDTO
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /knowledge-base/articles [post]
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
// @Summary Récupérer un article par ID
// @Description Récupère un article de la base de connaissances par son identifiant
// @Tags knowledge-base
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID de l'article"
// @Success 200 {object} dto.KnowledgeArticleDTO
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /knowledge-base/articles/{id} [get]
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
// @Summary Récupérer les articles publiés
// @Description Récupère la liste des articles publiés (route publique)
// @Tags knowledge-base
// @Accept json
// @Produce json
// @Success 200 {array} dto.KnowledgeArticleDTO
// @Failure 500 {object} utils.Response
// @Router /knowledge-base/articles/published [get]
func (h *KnowledgeArticleHandler) GetPublished(c *gin.Context) {
	articles, err := h.knowledgeArticleService.GetPublished()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des articles")
		return
	}

	utils.SuccessResponse(c, articles, "Articles publiés récupérés avec succès")
}

// Search recherche des articles
// @Summary Rechercher des articles
// @Description Recherche des articles dans la base de connaissances (route publique)
// @Tags knowledge-base
// @Accept json
// @Produce json
// @Param q query string true "Terme de recherche"
// @Success 200 {array} dto.KnowledgeArticleSearchResultDTO
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /knowledge-base/articles/search [get]
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
// @Summary Supprimer un article
// @Description Supprime un article de la base de connaissances
// @Tags knowledge-base
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID de l'article"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /knowledge-base/articles/{id} [delete]
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
