package services

import (
	"errors"
	"strings"

	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/models"
	"github.com/mcicare/itsm-backend/internal/repositories"
)

// SearchService interface pour les opérations de recherche
type SearchService interface {
	GlobalSearch(query string, types []string, limit int) (*dto.GlobalSearchResultDTO, error)
	SearchTickets(query string, status string, limit int) ([]dto.TicketSearchResultDTO, error)
	SearchAssets(query string, category string, limit int) ([]dto.AssetSearchResultDTO, error)
	SearchKnowledgeBase(query string, category string, limit int) ([]dto.KnowledgeArticleSearchResultDTO, error)
}

// searchService implémente SearchService
type searchService struct {
	ticketRepo  repositories.TicketRepository
	assetRepo   repositories.AssetRepository
	articleRepo repositories.KnowledgeArticleRepository
}

// NewSearchService crée une nouvelle instance de SearchService
func NewSearchService(
	ticketRepo repositories.TicketRepository,
	assetRepo repositories.AssetRepository,
	articleRepo repositories.KnowledgeArticleRepository,
) SearchService {
	return &searchService{
		ticketRepo:  ticketRepo,
		assetRepo:   assetRepo,
		articleRepo: articleRepo,
	}
}

// GlobalSearch effectue une recherche globale dans tous les types
func (s *searchService) GlobalSearch(query string, types []string, limit int) (*dto.GlobalSearchResultDTO, error) {
	if query == "" {
		return nil, errors.New("requête de recherche vide")
	}

	if limit <= 0 || limit > 100 {
		limit = 20
	}

	result := &dto.GlobalSearchResultDTO{
		Query: query,
		Types: types,
	}

	// Si aucun type spécifié, rechercher dans tous
	if len(types) == 0 {
		types = []string{"tickets", "assets", "articles"}
	}

	// Rechercher dans les tickets
	if contains(types, "tickets") {
		tickets, err := s.searchTicketsInternal(query, "", limit)
		if err == nil {
			result.Tickets = tickets
		}
	}

	// Rechercher dans les actifs
	if contains(types, "assets") {
		assets, err := s.searchAssetsInternal(query, "", limit)
		if err == nil {
			result.Assets = assets
		}
	}

	// Rechercher dans la base de connaissances
	if contains(types, "articles") {
		articles, err := s.searchKnowledgeBaseInternal(query, "", limit)
		if err == nil {
			result.Articles = articles
		}
	}

	result.Total = len(result.Tickets) + len(result.Assets) + len(result.Articles)

	return result, nil
}

// SearchTickets recherche dans les tickets
func (s *searchService) SearchTickets(query string, status string, limit int) ([]dto.TicketSearchResultDTO, error) {
	return s.searchTicketsInternal(query, status, limit)
}

// SearchAssets recherche dans les actifs
func (s *searchService) SearchAssets(query string, category string, limit int) ([]dto.AssetSearchResultDTO, error) {
	return s.searchAssetsInternal(query, category, limit)
}

// SearchKnowledgeBase recherche dans la base de connaissances
func (s *searchService) SearchKnowledgeBase(query string, category string, limit int) ([]dto.KnowledgeArticleSearchResultDTO, error) {
	return s.searchKnowledgeBaseInternal(query, category, limit)
}

// searchTicketsInternal recherche interne dans les tickets
func (s *searchService) searchTicketsInternal(query string, status string, limit int) ([]dto.TicketSearchResultDTO, error) {
	tickets, err := s.ticketRepo.Search(query, status, limit)
	if err != nil {
		return nil, errors.New("erreur lors de la recherche dans les tickets")
	}

	resultDTOs := make([]dto.TicketSearchResultDTO, len(tickets))
	for i, ticket := range tickets {
		resultDTOs[i] = s.ticketToSearchResultDTO(&ticket, query)
	}

	return resultDTOs, nil
}

// searchAssetsInternal recherche interne dans les actifs
func (s *searchService) searchAssetsInternal(query string, category string, limit int) ([]dto.AssetSearchResultDTO, error) {
	assets, err := s.assetRepo.Search(query, category, limit)
	if err != nil {
		return nil, errors.New("erreur lors de la recherche dans les actifs")
	}

	resultDTOs := make([]dto.AssetSearchResultDTO, len(assets))
	for i, asset := range assets {
		resultDTOs[i] = s.assetToSearchResultDTO(&asset, query)
	}

	return resultDTOs, nil
}

// searchKnowledgeBaseInternal recherche interne dans la base de connaissances
func (s *searchService) searchKnowledgeBaseInternal(query string, category string, limit int) ([]dto.KnowledgeArticleSearchResultDTO, error) {
	articles, err := s.articleRepo.Search(query)
	if err != nil {
		return nil, errors.New("erreur lors de la recherche dans la base de connaissances")
	}

	// Filtrer par catégorie si spécifiée
	filteredArticles := []models.KnowledgeArticle{}
	for _, article := range articles {
		if category == "" || (article.Category.ID != 0 && strings.EqualFold(article.Category.Name, category)) {
			filteredArticles = append(filteredArticles, article)
		}
	}

	// Limiter les résultats
	if limit > 0 && len(filteredArticles) > limit {
		filteredArticles = filteredArticles[:limit]
	}

	// Convertir en DTOs
	resultDTOs := make([]dto.KnowledgeArticleSearchResultDTO, len(filteredArticles))
	for i, article := range filteredArticles {
		resultDTOs[i] = s.articleToSearchResultDTO(&article, query)
	}

	return resultDTOs, nil
}

// ticketToSearchResultDTO convertit un ticket en DTO de recherche
func (s *searchService) ticketToSearchResultDTO(ticket *models.Ticket, query string) dto.TicketSearchResultDTO {
	snippet := extractSnippet(ticket.Description, query, 150)
	
	result := dto.TicketSearchResultDTO{
		ID:        ticket.ID,
		Title:     ticket.Title,
		Snippet:   snippet,
		Status:    ticket.Status,
		Priority:  ticket.Priority,
		Category:  ticket.Category,
		CreatedAt: ticket.CreatedAt,
	}

	if ticket.CreatedBy.ID != 0 {
		createdByDTO := dto.UserDTO{
			ID:        ticket.CreatedBy.ID,
			Username:  ticket.CreatedBy.Username,
			Email:     ticket.CreatedBy.Email,
			FirstName: ticket.CreatedBy.FirstName,
			LastName:  ticket.CreatedBy.LastName,
			Role:      ticket.CreatedBy.Role.Name,
			IsActive:  ticket.CreatedBy.IsActive,
			CreatedAt: ticket.CreatedBy.CreatedAt,
			UpdatedAt: ticket.CreatedBy.UpdatedAt,
		}
		result.CreatedBy = &createdByDTO
	}

	if ticket.AssignedTo != nil {
		assignedToDTO := dto.UserDTO{
			ID:        ticket.AssignedTo.ID,
			Username:  ticket.AssignedTo.Username,
			Email:     ticket.AssignedTo.Email,
			FirstName: ticket.AssignedTo.FirstName,
			LastName:  ticket.AssignedTo.LastName,
			Role:      ticket.AssignedTo.Role.Name,
			IsActive:  ticket.AssignedTo.IsActive,
			CreatedAt: ticket.AssignedTo.CreatedAt,
			UpdatedAt: ticket.AssignedTo.UpdatedAt,
		}
		result.AssignedTo = &assignedToDTO
	}

	return result
}

// assetToSearchResultDTO convertit un actif en DTO de recherche
func (s *searchService) assetToSearchResultDTO(asset *models.Asset, query string) dto.AssetSearchResultDTO {
	snippet := extractSnippet(asset.Notes, query, 150)
	
	result := dto.AssetSearchResultDTO{
		ID:           asset.ID,
		Name:         asset.Name,
		Snippet:      snippet,
		SerialNumber: asset.SerialNumber,
		CategoryID:   asset.CategoryID,
		Status:       asset.Status,
		CreatedAt:    asset.CreatedAt,
	}

	if asset.Category.ID != 0 {
		categoryDTO := dto.AssetCategoryDTO{
			ID:          asset.Category.ID,
			Name:        asset.Category.Name,
			Description: asset.Category.Description,
		}
		if asset.Category.ParentID != nil {
			categoryDTO.ParentID = asset.Category.ParentID
		}
		result.Category = &categoryDTO
	}

	return result
}

// articleToSearchResultDTO convertit un article en DTO de recherche
func (s *searchService) articleToSearchResultDTO(article *models.KnowledgeArticle, query string) dto.KnowledgeArticleSearchResultDTO {
	snippet := extractSnippet(article.Content, query, 200)
	
	result := dto.KnowledgeArticleSearchResultDTO{
		ID:        article.ID,
		Title:      article.Title,
		Snippet:    snippet,
		CategoryID: article.CategoryID,
		AuthorID:   article.AuthorID,
		ViewCount:  article.ViewCount,
		CreatedAt:  article.CreatedAt,
	}

	if article.Category.ID != 0 {
		categoryDTO := dto.KnowledgeCategoryDTO{
			ID:          article.Category.ID,
			Name:        article.Category.Name,
			Description: article.Category.Description,
		}
		if article.Category.ParentID != nil {
			categoryDTO.ParentID = article.Category.ParentID
		}
		result.Category = &categoryDTO
	}

	return result
}

// extractSnippet extrait un extrait de texte autour de la requête
func extractSnippet(text, query string, maxLength int) string {
	if text == "" {
		return ""
	}

	queryLower := strings.ToLower(query)
	textLower := strings.ToLower(text)
	
	// Trouver la position de la requête
	pos := strings.Index(textLower, queryLower)
	if pos == -1 {
		// Si la requête n'est pas trouvée, retourner le début du texte
		if len(text) > maxLength {
			return text[:maxLength] + "..."
		}
		return text
	}

	// Extraire autour de la position
	start := pos - maxLength/2
	if start < 0 {
		start = 0
	}
	
	end := start + maxLength
	if end > len(text) {
		end = len(text)
	}

	snippet := text[start:end]
	if start > 0 {
		snippet = "..." + snippet
	}
	if end < len(text) {
		snippet = snippet + "..."
	}

	return snippet
}

// contains vérifie si une slice contient une valeur
func contains(slice []string, value string) bool {
	for _, item := range slice {
		if strings.EqualFold(item, value) {
			return true
		}
	}
	return false
}

