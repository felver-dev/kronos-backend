package dto

import "time"

// KnowledgeArticleDTO représente un article de la base de connaissances
type KnowledgeArticleDTO struct {
	ID          uint                    `json:"id"`
	Title       string                  `json:"title"`
	Content     string                  `json:"content"`
	CategoryID  uint                    `json:"category_id"`
	Category    *KnowledgeCategoryDTO   `json:"category,omitempty"` // Catégorie (optionnel)
	AuthorID    uint                    `json:"author_id"`
	Author      *UserDTO                `json:"author,omitempty"` // Auteur (optionnel)
	IsPublished bool                    `json:"is_published"`    // Si l'article est publié
	ViewCount   int                     `json:"view_count"`      // Nombre de vues
	CreatedAt   time.Time               `json:"created_at"`
	UpdatedAt   time.Time               `json:"updated_at"`
}

// KnowledgeCategoryDTO représente une catégorie d'articles
type KnowledgeCategoryDTO struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	ParentID    *uint  `json:"parent_id,omitempty"` // Catégorie parente (optionnel)
}

// CreateKnowledgeArticleRequest représente la requête de création d'un article
type CreateKnowledgeArticleRequest struct {
	Title       string `json:"title" binding:"required"`        // Titre (obligatoire)
	Content     string `json:"content" binding:"required"`      // Contenu (obligatoire)
	CategoryID  uint   `json:"category_id" binding:"required"` // ID catégorie (obligatoire)
	IsPublished bool   `json:"is_published,omitempty"`          // Si l'article est publié (optionnel, défaut: false)
}

// UpdateKnowledgeArticleRequest représente la requête de mise à jour d'un article
type UpdateKnowledgeArticleRequest struct {
	Title       string `json:"title,omitempty"`
	Content     string `json:"content,omitempty"`
	CategoryID  *uint  `json:"category_id,omitempty"`
	IsPublished *bool  `json:"is_published,omitempty"` // Statut de publication (optionnel)
}

// CreateKnowledgeCategoryRequest représente la requête de création d'une catégorie
type CreateKnowledgeCategoryRequest struct {
	Name        string `json:"name" binding:"required"` // Nom (obligatoire)
	Description string `json:"description,omitempty"`   // Description (optionnel)
	ParentID    *uint  `json:"parent_id,omitempty"`    // ID catégorie parente (optionnel)
}

// UpdateKnowledgeCategoryRequest représente la requête de mise à jour d'une catégorie
type UpdateKnowledgeCategoryRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	ParentID    *uint  `json:"parent_id,omitempty"` // nil pour retirer la catégorie parente
}

// KnowledgeArticleSearchResultDTO représente un résultat de recherche d'article
type KnowledgeArticleSearchResultDTO struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Snippet     string `json:"snippet"`     // Extrait du contenu correspondant
	CategoryID  uint   `json:"category_id"`
	Category    *KnowledgeCategoryDTO `json:"category,omitempty"`
	AuthorID    uint   `json:"author_id"`
	ViewCount   int    `json:"view_count"`
	CreatedAt   time.Time `json:"created_at"`
}

