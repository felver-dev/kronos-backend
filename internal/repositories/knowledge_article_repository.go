package repositories

import (
	"github.com/mcicare/itsm-backend/database"
	"github.com/mcicare/itsm-backend/internal/models"
	"github.com/mcicare/itsm-backend/internal/scope"
)

// KnowledgeArticleRepository interface pour les opérations sur les articles de la base de connaissances
type KnowledgeArticleRepository interface {
	Create(article *models.KnowledgeArticle) error
	FindByID(id uint) (*models.KnowledgeArticle, error)
	FindAll(scope interface{}) ([]models.KnowledgeArticle, error) // scope peut être *scope.QueryScope ou nil
	FindPublished(scope interface{}) ([]models.KnowledgeArticle, error)
	FindByCategory(scope interface{}, categoryID uint) ([]models.KnowledgeArticle, error)
	FindByAuthor(scope interface{}, authorID uint) ([]models.KnowledgeArticle, error) // scope peut être *scope.QueryScope ou nil
	Search(scope interface{}, query string) ([]models.KnowledgeArticle, error)
	Update(article *models.KnowledgeArticle) error
	Delete(id uint) error
	IncrementViewCount(id uint) error
}

// KnowledgeCategoryRepository interface pour les opérations sur les catégories de la base de connaissances
type KnowledgeCategoryRepository interface {
	Create(category *models.KnowledgeCategory) error
	FindByID(id uint) (*models.KnowledgeCategory, error)
	FindAll() ([]models.KnowledgeCategory, error)
	FindByParentID(parentID uint) ([]models.KnowledgeCategory, error)
	FindActive() ([]models.KnowledgeCategory, error)
	Update(category *models.KnowledgeCategory) error
	Delete(id uint) error
}

// knowledgeArticleRepository implémente KnowledgeArticleRepository
type knowledgeArticleRepository struct{}

// knowledgeCategoryRepository implémente KnowledgeCategoryRepository
type knowledgeCategoryRepository struct{}

// NewKnowledgeArticleRepository crée une nouvelle instance de KnowledgeArticleRepository
func NewKnowledgeArticleRepository() KnowledgeArticleRepository {
	return &knowledgeArticleRepository{}
}

// NewKnowledgeCategoryRepository crée une nouvelle instance de KnowledgeCategoryRepository
func NewKnowledgeCategoryRepository() KnowledgeCategoryRepository {
	return &knowledgeCategoryRepository{}
}

// Create crée un nouvel article
func (r *knowledgeArticleRepository) Create(article *models.KnowledgeArticle) error {
	return database.DB.Create(article).Error
}

// FindByID trouve un article par son ID
func (r *knowledgeArticleRepository) FindByID(id uint) (*models.KnowledgeArticle, error) {
	var article models.KnowledgeArticle
	err := database.DB.Preload("Category").Preload("Author").Preload("Attachments").First(&article, id).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

// FindAll récupère tous les articles
// Le scope est utilisé pour filtrer automatiquement selon les permissions de l'utilisateur
func (r *knowledgeArticleRepository) FindAll(scopeParam interface{}) ([]models.KnowledgeArticle, error) {
	var articles []models.KnowledgeArticle
	
	// Construire la requête de base
	query := database.DB.Model(&models.KnowledgeArticle{}).
		Preload("Category").Preload("Author")
	
	// Appliquer le scope si fourni
	if scopeParam != nil {
		if queryScope, ok := scopeParam.(*scope.QueryScope); ok {
			query = scope.ApplyKnowledgeScope(query, queryScope)
		}
	}
	
	err := query.Find(&articles).Error
	return articles, err
}

// FindPublished récupère tous les articles publiés
// Le scope est utilisé pour filtrer automatiquement selon les permissions de l'utilisateur
func (r *knowledgeArticleRepository) FindPublished(scopeParam interface{}) ([]models.KnowledgeArticle, error) {
	var articles []models.KnowledgeArticle
	
	// Construire la requête de base
	query := database.DB.Model(&models.KnowledgeArticle{}).
		Preload("Category").Preload("Author").
		Where("knowledge_articles.is_published = ?", true)
	
	// Appliquer le scope si fourni
	if scopeParam != nil {
		if queryScope, ok := scopeParam.(*scope.QueryScope); ok {
			query = scope.ApplyKnowledgeScope(query, queryScope)
		}
	}
	
	err := query.Order("knowledge_articles.created_at DESC").Find(&articles).Error
	return articles, err
}

// FindByCategory récupère les articles d'une catégorie
// Le scope est utilisé pour filtrer automatiquement selon les permissions de l'utilisateur
func (r *knowledgeArticleRepository) FindByCategory(scopeParam interface{}, categoryID uint) ([]models.KnowledgeArticle, error) {
	var articles []models.KnowledgeArticle
	
	// Construire la requête de base
	query := database.DB.Model(&models.KnowledgeArticle{}).
		Preload("Category").Preload("Author").
		Where("knowledge_articles.category_id = ?", categoryID)
	
	// Appliquer le scope si fourni
	if scopeParam != nil {
		if queryScope, ok := scopeParam.(*scope.QueryScope); ok {
			query = scope.ApplyKnowledgeScope(query, queryScope)
		}
	}
	
	err := query.Find(&articles).Error
	return articles, err
}

// FindByAuthor récupère les articles d'un auteur
// Le scope est utilisé pour filtrer automatiquement selon les permissions de l'utilisateur
func (r *knowledgeArticleRepository) FindByAuthor(scopeParam interface{}, authorID uint) ([]models.KnowledgeArticle, error) {
	var articles []models.KnowledgeArticle
	
	// Construire la requête de base
	query := database.DB.Model(&models.KnowledgeArticle{}).
		Preload("Category").Preload("Author").
		Where("knowledge_articles.author_id = ?", authorID)
	
	// Appliquer le scope si fourni
	if scopeParam != nil {
		if queryScope, ok := scopeParam.(*scope.QueryScope); ok {
			query = scope.ApplyKnowledgeScope(query, queryScope)
		}
	}
	
	err := query.Find(&articles).Error
	return articles, err
}

// Search recherche des articles par titre ou contenu
// Le scope est utilisé pour filtrer automatiquement selon les permissions de l'utilisateur
func (r *knowledgeArticleRepository) Search(scopeParam interface{}, searchQuery string) ([]models.KnowledgeArticle, error) {
	var articles []models.KnowledgeArticle
	
	// Construire la requête de base
	query := database.DB.Model(&models.KnowledgeArticle{}).
		Preload("Category").Preload("Author").
		Where("(knowledge_articles.title LIKE ? OR knowledge_articles.content LIKE ?)", "%"+searchQuery+"%", "%"+searchQuery+"%")
	
	// Appliquer le scope si fourni
	if scopeParam != nil {
		if queryScope, ok := scopeParam.(*scope.QueryScope); ok {
			query = scope.ApplyKnowledgeScope(query, queryScope)
		}
	}
	
	err := query.Order("knowledge_articles.created_at DESC").Find(&articles).Error
	return articles, err
}

// Update met à jour un article
func (r *knowledgeArticleRepository) Update(article *models.KnowledgeArticle) error {
	return database.DB.Save(article).Error
}

// Delete supprime un article (soft delete)
func (r *knowledgeArticleRepository) Delete(id uint) error {
	return database.DB.Delete(&models.KnowledgeArticle{}, id).Error
}

// IncrementViewCount incrémente le compteur de vues d'un article
func (r *knowledgeArticleRepository) IncrementViewCount(id uint) error {
	return database.DB.Model(&models.KnowledgeArticle{}).Where("id = ?", id).Update("view_count", database.DB.Raw("view_count + 1")).Error
}

// Create crée une nouvelle catégorie
func (r *knowledgeCategoryRepository) Create(category *models.KnowledgeCategory) error {
	return database.DB.Create(category).Error
}

// FindByID trouve une catégorie par son ID
func (r *knowledgeCategoryRepository) FindByID(id uint) (*models.KnowledgeCategory, error) {
	var category models.KnowledgeCategory
	err := database.DB.Preload("Parent").Preload("Children").First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// FindAll récupère toutes les catégories
func (r *knowledgeCategoryRepository) FindAll() ([]models.KnowledgeCategory, error) {
	var categories []models.KnowledgeCategory
	err := database.DB.Preload("Parent").Find(&categories).Error
	return categories, err
}

// FindByParentID récupère les catégories enfants d'une catégorie parente
func (r *knowledgeCategoryRepository) FindByParentID(parentID uint) ([]models.KnowledgeCategory, error) {
	var categories []models.KnowledgeCategory
	err := database.DB.Where("parent_id = ?", parentID).Find(&categories).Error
	return categories, err
}

// FindActive récupère toutes les catégories actives
func (r *knowledgeCategoryRepository) FindActive() ([]models.KnowledgeCategory, error) {
	var categories []models.KnowledgeCategory
	err := database.DB.Preload("Parent").Where("is_active = ?", true).Find(&categories).Error
	return categories, err
}

// Update met à jour une catégorie
func (r *knowledgeCategoryRepository) Update(category *models.KnowledgeCategory) error {
	return database.DB.Save(category).Error
}

// Delete supprime une catégorie
func (r *knowledgeCategoryRepository) Delete(id uint) error {
	return database.DB.Delete(&models.KnowledgeCategory{}, id).Error
}
