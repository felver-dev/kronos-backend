package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/services"
	"github.com/mcicare/itsm-backend/internal/utils"
)

// SearchHandler gère les handlers de recherche
type SearchHandler struct {
	searchService services.SearchService
}

// NewSearchHandler crée une nouvelle instance de SearchHandler
func NewSearchHandler(searchService services.SearchService) *SearchHandler {
	return &SearchHandler{
		searchService: searchService,
	}
}

// GlobalSearch effectue une recherche globale
// @Summary Recherche globale
// @Description Effectue une recherche dans tous les types (tickets, actifs, articles)
// @Tags search
// @Security BearerAuth
// @Produce json
// @Param q query string true "Terme de recherche"
// @Param types query string false "Types à rechercher (tickets,assets,articles) - séparés par des virgules"
// @Param limit query int false "Limite de résultats (défaut: 20, max: 100)"
// @Success 200 {object} dto.GlobalSearchResultDTO
// @Failure 400 {object} utils.Response
// @Router /search [get]
func (h *SearchHandler) GlobalSearch(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		utils.BadRequestResponse(c, "Paramètre de recherche 'q' manquant")
		return
	}

	typesStr := c.Query("types")
	var types []string
	if typesStr != "" {
		types = strings.Split(typesStr, ",")
		// Nettoyer les espaces
		for i, t := range types {
			types[i] = strings.TrimSpace(t)
		}
	}

	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	result, err := h.searchService.GlobalSearch(query, types, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, result, "Recherche effectuée avec succès")
}

// SearchTickets recherche dans les tickets
// @Summary Rechercher dans les tickets
// @Description Effectue une recherche dans les tickets
// @Tags search
// @Security BearerAuth
// @Produce json
// @Param q query string true "Terme de recherche"
// @Param status query string false "Filtrer par statut"
// @Param limit query int false "Limite de résultats (défaut: 20, max: 100)"
// @Success 200 {array} dto.TicketSearchResultDTO
// @Failure 400 {object} utils.Response
// @Router /search/tickets [get]
func (h *SearchHandler) SearchTickets(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		utils.BadRequestResponse(c, "Paramètre de recherche 'q' manquant")
		return
	}

	status := c.Query("status")

	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	results, err := h.searchService.SearchTickets(query, status, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, results, "Recherche dans les tickets effectuée avec succès")
}

// SearchAssets recherche dans les actifs
// @Summary Rechercher dans les actifs
// @Description Effectue une recherche dans les actifs IT
// @Tags search
// @Security BearerAuth
// @Produce json
// @Param q query string true "Terme de recherche"
// @Param category query string false "Filtrer par catégorie"
// @Param limit query int false "Limite de résultats (défaut: 20, max: 100)"
// @Success 200 {array} dto.AssetSearchResultDTO
// @Failure 400 {object} utils.Response
// @Router /search/assets [get]
func (h *SearchHandler) SearchAssets(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		utils.BadRequestResponse(c, "Paramètre de recherche 'q' manquant")
		return
	}

	category := c.Query("category")

	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	results, err := h.searchService.SearchAssets(query, category, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, results, "Recherche dans les actifs effectuée avec succès")
}

// SearchKnowledgeBase recherche dans la base de connaissances
// @Summary Rechercher dans la base de connaissances
// @Description Effectue une recherche dans les articles de la base de connaissances
// @Tags search
// @Security BearerAuth
// @Produce json
// @Param q query string true "Terme de recherche"
// @Param category query string false "Filtrer par catégorie"
// @Param limit query int false "Limite de résultats (défaut: 20, max: 100)"
// @Success 200 {array} dto.KnowledgeArticleSearchResultDTO
// @Failure 400 {object} utils.Response
// @Router /search/knowledge-base [get]
func (h *SearchHandler) SearchKnowledgeBase(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		utils.BadRequestResponse(c, "Paramètre de recherche 'q' manquant")
		return
	}

	category := c.Query("category")

	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	results, err := h.searchService.SearchKnowledgeBase(query, category, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, results, "Recherche dans la base de connaissances effectuée avec succès")
}

