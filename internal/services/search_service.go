package services

import (
	"errors"
	"strings"
	"time"

	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/models"
	"github.com/mcicare/itsm-backend/internal/repositories"
)

// SearchService interface pour les opérations de recherche
type SearchService interface {
	GlobalSearch(scope interface{}, query string, types []string, limit int) (*dto.GlobalSearchResultDTO, error) // scope peut être *scope.QueryScope ou nil
	SearchTickets(scope interface{}, query string, status string, limit int) ([]dto.TicketSearchResultDTO, error) // scope peut être *scope.QueryScope ou nil
	SearchAssets(scope interface{}, query string, category string, limit int) ([]dto.AssetSearchResultDTO, error) // scope peut être *scope.QueryScope ou nil
	SearchKnowledgeBase(scope interface{}, query string, category string, limit int) ([]dto.KnowledgeArticleSearchResultDTO, error) // scope peut être *scope.QueryScope ou nil
	SearchUsers(scope interface{}, query string, limit int) ([]dto.UserSearchResultDTO, error) // scope peut être *scope.QueryScope ou nil
	SearchTimeEntries(scope interface{}, query string, limit int) ([]dto.TimeEntrySearchResultDTO, error) // scope peut être *scope.QueryScope ou nil
}

// searchService implémente SearchService
type searchService struct {
	ticketRepo  repositories.TicketRepository
	assetRepo   repositories.AssetRepository
	articleRepo repositories.KnowledgeArticleRepository
	userRepo    repositories.UserRepository
	timeEntryRepo repositories.TimeEntryRepository
}

// NewSearchService crée une nouvelle instance de SearchService
func NewSearchService(
	ticketRepo repositories.TicketRepository,
	assetRepo repositories.AssetRepository,
	articleRepo repositories.KnowledgeArticleRepository,
	userRepo repositories.UserRepository,
	timeEntryRepo repositories.TimeEntryRepository,
) SearchService {
	return &searchService{
		ticketRepo:  ticketRepo,
		assetRepo:   assetRepo,
		articleRepo: articleRepo,
		userRepo:    userRepo,
		timeEntryRepo: timeEntryRepo,
	}
}

// GlobalSearch effectue une recherche globale dans tous les types
// Le scope est utilisé pour filtrer automatiquement selon les permissions de l'utilisateur
func (s *searchService) GlobalSearch(scopeParam interface{}, query string, types []string, limit int) (*dto.GlobalSearchResultDTO, error) {
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
		types = []string{"tickets", "assets", "articles", "users", "time_entries"}
	}

	// Rechercher dans les tickets
	if contains(types, "tickets") {
		tickets, err := s.searchTicketsInternal(scopeParam, query, "", limit)
		if err == nil {
			result.Tickets = tickets
		}
	}

	// Rechercher dans les actifs
	if contains(types, "assets") {
		assets, err := s.searchAssetsInternal(scopeParam, query, "", limit)
		if err == nil {
			result.Assets = assets
		}
	}

	// Rechercher dans la base de connaissances
	if contains(types, "articles") {
		articles, err := s.searchKnowledgeBaseInternal(scopeParam, query, "", limit)
		if err == nil {
			result.Articles = articles
		}
	}

	// Rechercher dans les utilisateurs
	if contains(types, "users") {
		users, err := s.searchUsersInternal(scopeParam, query, limit)
		if err == nil {
			result.Users = users
		}
	}

	// Rechercher dans les entrées de temps
	if contains(types, "time_entries") {
		entries, err := s.searchTimeEntriesInternal(scopeParam, query, limit)
		if err == nil {
			result.TimeEntries = entries
		}
	}

	result.Total = len(result.Tickets) + len(result.Assets) + len(result.Articles) + len(result.Users) + len(result.TimeEntries)

	return result, nil
}

// SearchTickets recherche dans les tickets
func (s *searchService) SearchTickets(scopeParam interface{}, query string, status string, limit int) ([]dto.TicketSearchResultDTO, error) {
	return s.searchTicketsInternal(scopeParam, query, status, limit)
}

// SearchAssets recherche dans les actifs
func (s *searchService) SearchAssets(scopeParam interface{}, query string, category string, limit int) ([]dto.AssetSearchResultDTO, error) {
	return s.searchAssetsInternal(scopeParam, query, category, limit)
}

// SearchKnowledgeBase recherche dans la base de connaissances
func (s *searchService) SearchKnowledgeBase(scopeParam interface{}, query string, category string, limit int) ([]dto.KnowledgeArticleSearchResultDTO, error) {
	return s.searchKnowledgeBaseInternal(scopeParam, query, category, limit)
}

// SearchUsers recherche dans les utilisateurs
func (s *searchService) SearchUsers(scopeParam interface{}, query string, limit int) ([]dto.UserSearchResultDTO, error) {
	return s.searchUsersInternal(scopeParam, query, limit)
}

// SearchTimeEntries recherche dans les entrées de temps
func (s *searchService) SearchTimeEntries(scopeParam interface{}, query string, limit int) ([]dto.TimeEntrySearchResultDTO, error) {
	return s.searchTimeEntriesInternal(scopeParam, query, limit)
}

// searchTicketsInternal recherche interne dans les tickets
func (s *searchService) searchTicketsInternal(scopeParam interface{}, query string, status string, limit int) ([]dto.TicketSearchResultDTO, error) {
	tickets, err := s.ticketRepo.Search(scopeParam, query, status, limit)
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
func (s *searchService) searchAssetsInternal(scopeParam interface{}, query string, category string, limit int) ([]dto.AssetSearchResultDTO, error) {
	assets, err := s.assetRepo.Search(scopeParam, query, category, limit)
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
func (s *searchService) searchKnowledgeBaseInternal(scopeParam interface{}, query string, category string, limit int) ([]dto.KnowledgeArticleSearchResultDTO, error) {
	articles, err := s.articleRepo.Search(scopeParam, query)
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

// searchUsersInternal recherche interne dans les utilisateurs
func (s *searchService) searchUsersInternal(scopeParam interface{}, query string, limit int) ([]dto.UserSearchResultDTO, error) {
	users, err := s.userRepo.Search(scopeParam, query, limit)
	if err != nil {
		return nil, errors.New("erreur lors de la recherche dans les utilisateurs")
	}

	resultDTOs := make([]dto.UserSearchResultDTO, len(users))
	for i, user := range users {
		resultDTOs[i] = s.userToSearchResultDTO(&user, query)
	}

	return resultDTOs, nil
}

// searchTimeEntriesInternal recherche interne dans les entrées de temps
func (s *searchService) searchTimeEntriesInternal(scopeParam interface{}, query string, limit int) ([]dto.TimeEntrySearchResultDTO, error) {
	entries, err := s.timeEntryRepo.Search(scopeParam, query, limit)
	if err != nil {
		return nil, errors.New("erreur lors de la recherche dans les entrées de temps")
	}

	resultDTOs := make([]dto.TimeEntrySearchResultDTO, len(entries))
	for i, entry := range entries {
		resultDTOs[i] = s.timeEntryToSearchResultDTO(&entry, query)
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

// userToSearchResultDTO convertit un utilisateur en DTO de recherche
func (s *searchService) userToSearchResultDTO(user *models.User, query string) dto.UserSearchResultDTO {
	fullText := strings.TrimSpace(strings.Join([]string{user.Username, user.Email, user.FirstName, user.LastName}, " "))
	snippet := extractSnippet(fullText, query, 120)

	result := dto.UserSearchResultDTO{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role.Name,
		IsActive:  user.IsActive,
		Snippet:   snippet,
		CreatedAt: user.CreatedAt,
	}

	if user.Department.ID != 0 {
		departmentDTO := dto.DepartmentDTO{
			ID:          user.Department.ID,
			Name:        user.Department.Name,
			Code:        user.Department.Code,
			Description: user.Department.Description,
			OfficeID:    user.Department.OfficeID,
			IsActive:    user.Department.IsActive,
			CreatedAt:   user.Department.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   user.Department.UpdatedAt.Format(time.RFC3339),
		}
		result.Department = &departmentDTO
	}

	return result
}

// timeEntryToSearchResultDTO convertit une entrée de temps en DTO de recherche
func (s *searchService) timeEntryToSearchResultDTO(entry *models.TimeEntry, query string) dto.TimeEntrySearchResultDTO {
	snippetSource := entry.Description
	if snippetSource == "" {
		snippetSource = entry.Ticket.Title
	}
	snippet := extractSnippet(snippetSource, query, 120)

	ticketID := uint(0)
	if entry.TicketID != nil {
		ticketID = *entry.TicketID
	}
	result := dto.TimeEntrySearchResultDTO{
		ID:          entry.ID,
		TicketID:    ticketID,
		UserID:      entry.UserID,
		TimeSpent:   entry.TimeSpent,
		Date:        entry.Date,
		Description: entry.Description,
		Snippet:     snippet,
		Validated:   entry.Validated,
		CreatedAt:   entry.CreatedAt,
	}

	if entry.Ticket != nil && entry.Ticket.ID != 0 {
		ticketDTO := dto.TicketDTO{
			ID:       entry.Ticket.ID,
			Code:     entry.Ticket.Code,
			Title:    entry.Ticket.Title,
			Status:   entry.Ticket.Status,
			Priority: entry.Ticket.Priority,
			Category: entry.Ticket.Category,
			CreatedAt: entry.Ticket.CreatedAt,
			UpdatedAt: entry.Ticket.UpdatedAt,
		}
		result.Ticket = &ticketDTO
	}

	if entry.User.ID != 0 {
		userDTO := dto.UserDTO{
			ID:        entry.User.ID,
			Username:  entry.User.Username,
			Email:     entry.User.Email,
			FirstName: entry.User.FirstName,
			LastName:  entry.User.LastName,
			Role:      entry.User.Role.Name,
			IsActive:  entry.User.IsActive,
			CreatedAt: entry.User.CreatedAt,
			UpdatedAt: entry.User.UpdatedAt,
		}
		if entry.User.Department.ID != 0 {
			departmentDTO := dto.DepartmentDTO{
				ID:          entry.User.Department.ID,
				Name:        entry.User.Department.Name,
				Code:        entry.User.Department.Code,
				Description: entry.User.Department.Description,
				OfficeID:    entry.User.Department.OfficeID,
				IsActive:    entry.User.Department.IsActive,
				CreatedAt:   entry.User.Department.CreatedAt.Format(time.RFC3339),
				UpdatedAt:   entry.User.Department.UpdatedAt.Format(time.RFC3339),
			}
			userDTO.Department = &departmentDTO
			userDTO.DepartmentID = &departmentDTO.ID
		}
		result.User = &userDTO
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

