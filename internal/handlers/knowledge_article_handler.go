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

// GetAll récupère tous les articles
// @Summary Récupérer tous les articles
// @Description Récupère la liste de tous les articles (publiés et non publiés)
// @Tags knowledge-base
// @Security BearerAuth
// @Produce json
// @Success 200 {array} dto.KnowledgeArticleDTO
// @Failure 500 {object} utils.Response
// @Router /knowledge-base/articles [get]
func (h *KnowledgeArticleHandler) GetAll(c *gin.Context) {
	articles, err := h.knowledgeArticleService.GetAll()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des articles")
		return
	}

	utils.SuccessResponse(c, articles, "Articles récupérés avec succès")
}

// Update met à jour un article
// @Summary Mettre à jour un article
// @Description Met à jour les informations d'un article de la base de connaissances
// @Tags knowledge-base
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID de l'article"
// @Param request body dto.UpdateKnowledgeArticleRequest true "Données à mettre à jour"
// @Success 200 {object} dto.KnowledgeArticleDTO
// @Failure 400 {object} utils.Response
// @Router /knowledge-base/articles/{id} [put]
func (h *KnowledgeArticleHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	var req dto.UpdateKnowledgeArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	updatedByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	article, err := h.knowledgeArticleService.Update(uint(id), req, updatedByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, article, "Article mis à jour avec succès")
}

// Publish publie ou dépublie un article
// @Summary Publier/Dépublier un article
// @Description Change le statut de publication d'un article
// @Tags knowledge-base
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID de l'article"
// @Param request body map[string]bool true "Statut de publication"
// @Success 200 {object} dto.KnowledgeArticleDTO
// @Failure 400 {object} utils.Response
// @Router /knowledge-base/articles/{id}/publish [post]
func (h *KnowledgeArticleHandler) Publish(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	var req struct {
		Published bool `json:"published" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	updatedByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	article, err := h.knowledgeArticleService.Publish(uint(id), req.Published, updatedByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, article, "Statut de publication mis à jour avec succès")
}

// GetByCategory récupère les articles d'une catégorie
// @Summary Récupérer les articles par catégorie
// @Description Récupère les articles d'une catégorie spécifique
// @Tags knowledge-base
// @Security BearerAuth
// @Produce json
// @Param categoryId path int true "ID de la catégorie"
// @Success 200 {array} dto.KnowledgeArticleDTO
// @Failure 400 {object} utils.Response
// @Router /knowledge-base/articles/by-category/{categoryId} [get]
func (h *KnowledgeArticleHandler) GetByCategory(c *gin.Context) {
	categoryIDParam := c.Param("categoryId")
	categoryID, err := strconv.ParseUint(categoryIDParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID catégorie invalide")
		return
	}

	articles, err := h.knowledgeArticleService.GetByCategory(uint(categoryID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, articles, "Articles récupérés avec succès")
}

// GetByAuthor récupère les articles d'un auteur
// @Summary Récupérer les articles par auteur
// @Description Récupère les articles créés par un auteur spécifique
// @Tags knowledge-base
// @Security BearerAuth
// @Produce json
// @Param authorId path int true "ID de l'auteur"
// @Success 200 {array} dto.KnowledgeArticleDTO
// @Failure 400 {object} utils.Response
// @Router /knowledge-base/articles/by-author/{authorId} [get]
func (h *KnowledgeArticleHandler) GetByAuthor(c *gin.Context) {
	authorIDParam := c.Param("authorId")
	authorID, err := strconv.ParseUint(authorIDParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID auteur invalide")
		return
	}

	articles, err := h.knowledgeArticleService.GetByAuthor(uint(authorID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, articles, "Articles récupérés avec succès")
}

// IncrementViewCount incrémente le compteur de vues d'un article
// @Summary Incrémenter le compteur de vues
// @Description Incrémente le compteur de vues d'un article (appelé automatiquement lors de la consultation)
// @Tags knowledge-base
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID de l'article"
// @Success 200 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /knowledge-base/articles/{id}/view [post]
func (h *KnowledgeArticleHandler) IncrementViewCount(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	err = h.knowledgeArticleService.IncrementViewCount(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Article introuvable")
		return
	}

	utils.SuccessResponse(c, nil, "Compteur de vues incrémenté avec succès")
}
