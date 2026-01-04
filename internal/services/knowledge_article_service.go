package services

import (
	"errors"

	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/models"
	"github.com/mcicare/itsm-backend/internal/repositories"
)

// KnowledgeArticleService interface pour les opérations sur les articles de la base de connaissances
type KnowledgeArticleService interface {
	Create(req dto.CreateKnowledgeArticleRequest, authorID uint) (*dto.KnowledgeArticleDTO, error)
	GetByID(id uint) (*dto.KnowledgeArticleDTO, error)
	GetAll() ([]dto.KnowledgeArticleDTO, error)
	GetPublished() ([]dto.KnowledgeArticleDTO, error)
	GetByCategory(categoryID uint) ([]dto.KnowledgeArticleDTO, error)
	GetByAuthor(authorID uint) ([]dto.KnowledgeArticleDTO, error)
	Search(query string) ([]dto.KnowledgeArticleSearchResultDTO, error)
	Update(id uint, req dto.UpdateKnowledgeArticleRequest, updatedByID uint) (*dto.KnowledgeArticleDTO, error)
	Publish(id uint, published bool, updatedByID uint) (*dto.KnowledgeArticleDTO, error)
	Delete(id uint) error
	IncrementViewCount(id uint) error
}

// KnowledgeCategoryService interface pour les opérations sur les catégories de la base de connaissances
type KnowledgeCategoryService interface {
	Create(req dto.CreateKnowledgeCategoryRequest, createdByID uint) (*dto.KnowledgeCategoryDTO, error)
	GetByID(id uint) (*dto.KnowledgeCategoryDTO, error)
	GetAll() ([]dto.KnowledgeCategoryDTO, error)
	GetByParentID(parentID uint) ([]dto.KnowledgeCategoryDTO, error)
	GetActive() ([]dto.KnowledgeCategoryDTO, error)
	Update(id uint, req dto.UpdateKnowledgeCategoryRequest, updatedByID uint) (*dto.KnowledgeCategoryDTO, error)
	Delete(id uint) error
}

// knowledgeArticleService implémente KnowledgeArticleService
type knowledgeArticleService struct {
	articleRepo  repositories.KnowledgeArticleRepository
	categoryRepo repositories.KnowledgeCategoryRepository
	userRepo     repositories.UserRepository
}

// NewKnowledgeArticleService crée une nouvelle instance de KnowledgeArticleService
func NewKnowledgeArticleService(
	articleRepo repositories.KnowledgeArticleRepository,
	categoryRepo repositories.KnowledgeCategoryRepository,
	userRepo repositories.UserRepository,
) KnowledgeArticleService {
	return &knowledgeArticleService{
		articleRepo:  articleRepo,
		categoryRepo: categoryRepo,
		userRepo:     userRepo,
	}
}

// Create crée un nouvel article
func (s *knowledgeArticleService) Create(req dto.CreateKnowledgeArticleRequest, authorID uint) (*dto.KnowledgeArticleDTO, error) {
	// Vérifier que la catégorie existe
	_, err := s.categoryRepo.FindByID(req.CategoryID)
	if err != nil {
		return nil, errors.New("catégorie introuvable")
	}

	// Créer l'article
	article := &models.KnowledgeArticle{
		Title:       req.Title,
		Content:     req.Content,
		CategoryID:  req.CategoryID,
		AuthorID:    authorID,
		IsPublished: req.IsPublished,
		ViewCount:   0,
	}

	if err := s.articleRepo.Create(article); err != nil {
		return nil, errors.New("erreur lors de la création de l'article")
	}

	// Récupérer l'article créé avec ses relations
	createdArticle, err := s.articleRepo.FindByID(article.ID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération de l'article créé")
	}

	articleDTO := s.articleToDTO(createdArticle)
	return &articleDTO, nil
}

// GetByID récupère un article par son ID
func (s *knowledgeArticleService) GetByID(id uint) (*dto.KnowledgeArticleDTO, error) {
	article, err := s.articleRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("article introuvable")
	}

	// Incrémenter le compteur de vues si l'article est publié
	if article.IsPublished {
		s.articleRepo.IncrementViewCount(id)
	}

	articleDTO := s.articleToDTO(article)
	return &articleDTO, nil
}

// GetAll récupère tous les articles
func (s *knowledgeArticleService) GetAll() ([]dto.KnowledgeArticleDTO, error) {
	articles, err := s.articleRepo.FindAll()
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des articles")
	}

	var articleDTOs []dto.KnowledgeArticleDTO
	for _, article := range articles {
		articleDTOs = append(articleDTOs, s.articleToDTO(&article))
	}

	return articleDTOs, nil
}

// GetPublished récupère les articles publiés
func (s *knowledgeArticleService) GetPublished() ([]dto.KnowledgeArticleDTO, error) {
	articles, err := s.articleRepo.FindPublished()
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des articles")
	}

	var articleDTOs []dto.KnowledgeArticleDTO
	for _, article := range articles {
		articleDTOs = append(articleDTOs, s.articleToDTO(&article))
	}

	return articleDTOs, nil
}

// GetByCategory récupère les articles d'une catégorie
func (s *knowledgeArticleService) GetByCategory(categoryID uint) ([]dto.KnowledgeArticleDTO, error) {
	articles, err := s.articleRepo.FindByCategory(categoryID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des articles")
	}

	var articleDTOs []dto.KnowledgeArticleDTO
	for _, article := range articles {
		articleDTOs = append(articleDTOs, s.articleToDTO(&article))
	}

	return articleDTOs, nil
}

// GetByAuthor récupère les articles d'un auteur
func (s *knowledgeArticleService) GetByAuthor(authorID uint) ([]dto.KnowledgeArticleDTO, error) {
	articles, err := s.articleRepo.FindByAuthor(authorID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des articles")
	}

	var articleDTOs []dto.KnowledgeArticleDTO
	for _, article := range articles {
		articleDTOs = append(articleDTOs, s.articleToDTO(&article))
	}

	return articleDTOs, nil
}

// Search recherche des articles
func (s *knowledgeArticleService) Search(query string) ([]dto.KnowledgeArticleSearchResultDTO, error) {
	articles, err := s.articleRepo.Search(query)
	if err != nil {
		return nil, errors.New("erreur lors de la recherche des articles")
	}

	var resultDTOs []dto.KnowledgeArticleSearchResultDTO
	for _, article := range articles {
		resultDTOs = append(resultDTOs, s.articleToSearchResultDTO(&article))
	}

	return resultDTOs, nil
}

// Update met à jour un article
func (s *knowledgeArticleService) Update(id uint, req dto.UpdateKnowledgeArticleRequest, updatedByID uint) (*dto.KnowledgeArticleDTO, error) {
	article, err := s.articleRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("article introuvable")
	}

	// Vérifier que l'utilisateur est l'auteur ou a les droits
	if article.AuthorID != updatedByID {
		// TODO: Vérifier les permissions (admin, etc.)
	}

	// Mettre à jour les champs fournis
	if req.Title != "" {
		article.Title = req.Title
	}
	if req.Content != "" {
		article.Content = req.Content
	}
	if req.CategoryID != nil {
		// Vérifier que la catégorie existe
		_, err = s.categoryRepo.FindByID(*req.CategoryID)
		if err != nil {
			return nil, errors.New("catégorie introuvable")
		}
		article.CategoryID = *req.CategoryID
	}
	if req.IsPublished != nil {
		article.IsPublished = *req.IsPublished
	}

	if err := s.articleRepo.Update(article); err != nil {
		return nil, errors.New("erreur lors de la mise à jour de l'article")
	}

	// Récupérer l'article mis à jour
	updatedArticle, err := s.articleRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération de l'article mis à jour")
	}

	articleDTO := s.articleToDTO(updatedArticle)
	return &articleDTO, nil
}

// Publish publie ou dépublie un article
func (s *knowledgeArticleService) Publish(id uint, published bool, updatedByID uint) (*dto.KnowledgeArticleDTO, error) {
	article, err := s.articleRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("article introuvable")
	}

	article.IsPublished = published

	if err := s.articleRepo.Update(article); err != nil {
		return nil, errors.New("erreur lors de la mise à jour de l'article")
	}

	// Récupérer l'article mis à jour
	updatedArticle, err := s.articleRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération de l'article mis à jour")
	}

	articleDTO := s.articleToDTO(updatedArticle)
	return &articleDTO, nil
}

// Delete supprime un article
func (s *knowledgeArticleService) Delete(id uint) error {
	_, err := s.articleRepo.FindByID(id)
	if err != nil {
		return errors.New("article introuvable")
	}

	if err := s.articleRepo.Delete(id); err != nil {
		return errors.New("erreur lors de la suppression de l'article")
	}

	return nil
}

// IncrementViewCount incrémente le compteur de vues d'un article
func (s *knowledgeArticleService) IncrementViewCount(id uint) error {
	return s.articleRepo.IncrementViewCount(id)
}

// articleToDTO convertit un modèle KnowledgeArticle en DTO
func (s *knowledgeArticleService) articleToDTO(article *models.KnowledgeArticle) dto.KnowledgeArticleDTO {
	articleDTO := dto.KnowledgeArticleDTO{
		ID:          article.ID,
		Title:       article.Title,
		Content:     article.Content,
		CategoryID:  article.CategoryID,
		AuthorID:    article.AuthorID,
		IsPublished: article.IsPublished,
		ViewCount:   article.ViewCount,
		CreatedAt:   article.CreatedAt,
		UpdatedAt:   article.UpdatedAt,
	}

	// Convertir la catégorie si présente
	if article.Category.ID != 0 {
		categoryDTO := s.categoryToDTO(&article.Category)
		articleDTO.Category = &categoryDTO
	}

	// Convertir l'auteur si présent
	if article.Author.ID != 0 {
		userDTO := s.userToDTO(&article.Author)
		articleDTO.Author = &userDTO
	}

	return articleDTO
}

// articleToSearchResultDTO convertit un modèle KnowledgeArticle en DTO de recherche
func (s *knowledgeArticleService) articleToSearchResultDTO(article *models.KnowledgeArticle) dto.KnowledgeArticleSearchResultDTO {
	// Créer un extrait du contenu (premiers 200 caractères)
	snippet := article.Content
	if len(snippet) > 200 {
		snippet = snippet[:200] + "..."
	}

	resultDTO := dto.KnowledgeArticleSearchResultDTO{
		ID:         article.ID,
		Title:      article.Title,
		Snippet:    snippet,
		CategoryID: article.CategoryID,
		AuthorID:   article.AuthorID,
		ViewCount:  article.ViewCount,
		CreatedAt:  article.CreatedAt,
	}

	// Convertir la catégorie si présente
	if article.Category.ID != 0 {
		categoryDTO := s.categoryToDTO(&article.Category)
		resultDTO.Category = &categoryDTO
	}

	return resultDTO
}

// categoryToDTO convertit un modèle KnowledgeCategory en DTO
func (s *knowledgeArticleService) categoryToDTO(category *models.KnowledgeCategory) dto.KnowledgeCategoryDTO {
	categoryDTO := dto.KnowledgeCategoryDTO{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}

	if category.ParentID != nil {
		categoryDTO.ParentID = category.ParentID
	}

	return categoryDTO
}

// userToDTO convertit un modèle User en DTO (méthode helper)
func (s *knowledgeArticleService) userToDTO(user *models.User) dto.UserDTO {
	userDTO := dto.UserDTO{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Avatar:    user.Avatar,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	if user.RoleID != 0 {
		userDTO.Role = user.Role.Name
	}

	if user.LastLogin != nil {
		userDTO.LastLogin = user.LastLogin
	}

	return userDTO
}

// knowledgeCategoryService implémente KnowledgeCategoryService
type knowledgeCategoryService struct {
	categoryRepo repositories.KnowledgeCategoryRepository
	userRepo     repositories.UserRepository
}

// NewKnowledgeCategoryService crée une nouvelle instance de KnowledgeCategoryService
func NewKnowledgeCategoryService(
	categoryRepo repositories.KnowledgeCategoryRepository,
	userRepo repositories.UserRepository,
) KnowledgeCategoryService {
	return &knowledgeCategoryService{
		categoryRepo: categoryRepo,
		userRepo:     userRepo,
	}
}

// Create crée une nouvelle catégorie
func (s *knowledgeCategoryService) Create(req dto.CreateKnowledgeCategoryRequest, createdByID uint) (*dto.KnowledgeCategoryDTO, error) {
	category := &models.KnowledgeCategory{
		Name:        req.Name,
		Description: req.Description,
		ParentID:    req.ParentID,
		IsActive:    true,
	}

	if err := s.categoryRepo.Create(category); err != nil {
		return nil, errors.New("erreur lors de la création de la catégorie")
	}

	createdCategory, err := s.categoryRepo.FindByID(category.ID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération de la catégorie créée")
	}

	categoryDTO := s.categoryToDTO(createdCategory)
	return &categoryDTO, nil
}

// GetByID récupère une catégorie par son ID
func (s *knowledgeCategoryService) GetByID(id uint) (*dto.KnowledgeCategoryDTO, error) {
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("catégorie introuvable")
	}

	categoryDTO := s.categoryToDTO(category)
	return &categoryDTO, nil
}

// GetAll récupère toutes les catégories
func (s *knowledgeCategoryService) GetAll() ([]dto.KnowledgeCategoryDTO, error) {
	categories, err := s.categoryRepo.FindAll()
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des catégories")
	}

	var categoryDTOs []dto.KnowledgeCategoryDTO
	for _, category := range categories {
		categoryDTOs = append(categoryDTOs, s.categoryToDTO(&category))
	}

	return categoryDTOs, nil
}

// GetByParentID récupère les catégories enfants d'une catégorie parente
func (s *knowledgeCategoryService) GetByParentID(parentID uint) ([]dto.KnowledgeCategoryDTO, error) {
	categories, err := s.categoryRepo.FindByParentID(parentID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des catégories")
	}

	var categoryDTOs []dto.KnowledgeCategoryDTO
	for _, category := range categories {
		categoryDTOs = append(categoryDTOs, s.categoryToDTO(&category))
	}

	return categoryDTOs, nil
}

// GetActive récupère toutes les catégories actives
func (s *knowledgeCategoryService) GetActive() ([]dto.KnowledgeCategoryDTO, error) {
	categories, err := s.categoryRepo.FindActive()
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des catégories")
	}

	var categoryDTOs []dto.KnowledgeCategoryDTO
	for _, category := range categories {
		categoryDTOs = append(categoryDTOs, s.categoryToDTO(&category))
	}

	return categoryDTOs, nil
}

// Update met à jour une catégorie
func (s *knowledgeCategoryService) Update(id uint, req dto.UpdateKnowledgeCategoryRequest, updatedByID uint) (*dto.KnowledgeCategoryDTO, error) {
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("catégorie introuvable")
	}

	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Description != "" {
		category.Description = req.Description
	}
	if req.ParentID != nil {
		category.ParentID = req.ParentID
	}

	if err := s.categoryRepo.Update(category); err != nil {
		return nil, errors.New("erreur lors de la mise à jour de la catégorie")
	}

	updatedCategory, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération de la catégorie mise à jour")
	}

	categoryDTO := s.categoryToDTO(updatedCategory)
	return &categoryDTO, nil
}

// Delete supprime une catégorie
func (s *knowledgeCategoryService) Delete(id uint) error {
	_, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return errors.New("catégorie introuvable")
	}

	if err := s.categoryRepo.Delete(id); err != nil {
		return errors.New("erreur lors de la suppression de la catégorie")
	}

	return nil
}

// categoryToDTO convertit un modèle KnowledgeCategory en DTO (pour knowledgeCategoryService)
func (s *knowledgeCategoryService) categoryToDTO(category *models.KnowledgeCategory) dto.KnowledgeCategoryDTO {
	categoryDTO := dto.KnowledgeCategoryDTO{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}

	if category.ParentID != nil {
		categoryDTO.ParentID = category.ParentID
	}

	return categoryDTO
}
