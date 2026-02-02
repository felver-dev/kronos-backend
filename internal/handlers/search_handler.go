package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/services"
	"github.com/mcicare/itsm-backend/internal/utils"
)

// Alias pour la doc Swagger (évite "cannot find type definition: dto.GlobalSearchResultDTO")
type globalSearchResultDTO = dto.GlobalSearchResultDTO

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
// @Description Effectue une recherche dans tous les types (tickets, actifs, articles, utilisateurs, entrées de temps)
// @Tags search
// @Security BearerAuth
// @Produce json
// @Param q query string true "Terme de recherche"
// @Param types query string false "Types à rechercher (tickets,assets,articles,users,time_entries) - séparés par des virgules"
// @Param limit query int false "Limite de résultats (défaut: 20, max: 100)"
// @Success 200 {object} globalSearchResultDTO
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
			cleaned := strings.TrimSpace(strings.ToLower(t))
			cleaned = strings.ReplaceAll(cleaned, "-", "_")
			types[i] = normalizeSearchType(cleaned)
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

	// Extraire le QueryScope du contexte (injecté par AuthMiddleware)
	queryScope := utils.GetScopeFromContext(c)
	
	result, err := h.searchService.GlobalSearch(queryScope, query, types, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, result, "Recherche effectuée avec succès")
}

func normalizeSearchType(value string) string {
	switch value {
	case "ticket":
		return "tickets"
	case "asset":
		return "assets"
	case "article":
		return "articles"
	case "user":
		return "users"
	case "time_entry":
		return "time_entries"
	default:
		return value
	}
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

	// Extraire le QueryScope du contexte (injecté par AuthMiddleware)
	queryScope := utils.GetScopeFromContext(c)
	
	results, err := h.searchService.SearchTickets(queryScope, query, status, limit)
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

	// Extraire le QueryScope du contexte (injecté par AuthMiddleware)
	queryScope := utils.GetScopeFromContext(c)
	
	results, err := h.searchService.SearchAssets(queryScope, query, category, limit)
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

	// Extraire le QueryScope du contexte (injecté par AuthMiddleware)
	queryScope := utils.GetScopeFromContext(c)
	
	results, err := h.searchService.SearchKnowledgeBase(queryScope, query, category, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, results, "Recherche dans la base de connaissances effectuée avec succès")
}


// SearchUsers recherche dans les utilisateurs
// @Summary Rechercher dans les utilisateurs
// @Description Effectue une recherche dans les utilisateurs
// @Tags search
// @Security BearerAuth
// @Produce json
// @Param q query string true "Terme de recherche"
// @Param limit query int false "Limite de résultats (défaut: 20, max: 100)"
// @Success 200 {array} dto.UserSearchResultDTO
// @Failure 400 {object} utils.Response
// @Router /search/users [get]
func (h *SearchHandler) SearchUsers(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		utils.BadRequestResponse(c, "Paramètre de recherche 'q' manquant")
		return
	}

	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	// Extraire le QueryScope du contexte (injecté par AuthMiddleware)
	queryScope := utils.GetScopeFromContext(c)
	
	results, err := h.searchService.SearchUsers(queryScope, query, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, results, "Recherche dans les utilisateurs effectuée avec succès")
}

// SearchTimeEntries recherche dans les entrées de temps
// @Summary Rechercher dans les entrées de temps
// @Description Effectue une recherche dans les entrées de temps
// @Tags search
// @Security BearerAuth
// @Produce json
// @Param q query string true "Terme de recherche"
// @Param limit query int false "Limite de résultats (défaut: 20, max: 100)"
// @Success 200 {array} dto.TimeEntrySearchResultDTO
// @Failure 400 {object} utils.Response
// @Router /search/time-entries [get]
func (h *SearchHandler) SearchTimeEntries(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		utils.BadRequestResponse(c, "Paramètre de recherche 'q' manquant")
		return
	}

	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	// Extraire le QueryScope du contexte (injecté par AuthMiddleware)
	queryScope := utils.GetScopeFromContext(c)
	
	results, err := h.searchService.SearchTimeEntries(queryScope, query, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, results, "Recherche dans les entrées de temps effectuée avec succès")
}
