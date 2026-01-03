package models

import (
	"time"

	"gorm.io/gorm"
)

// KnowledgeCategory représente une catégorie d'article de la base de connaissances
// Table: knowledge_categories
type KnowledgeCategory struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Description string    `gorm:"type:text" json:"description,omitempty"`
	ParentID    *uint     `gorm:"index" json:"parent_id,omitempty"` // Catégorie parente (optionnel)
	IsActive    bool      `gorm:"default:true;index" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relations
	Parent   *KnowledgeCategory  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`   // Catégorie parente (optionnel)
	Children []KnowledgeCategory `gorm:"foreignKey:ParentID" json:"children,omitempty"` // Catégories enfants
	Articles []KnowledgeArticle  `gorm:"foreignKey:CategoryID" json:"-"`                // Articles de cette catégorie
}

// TableName spécifie le nom de la table
func (KnowledgeCategory) TableName() string {
	return "knowledge_categories"
}

// KnowledgeArticle représente un article de la base de connaissances
// Table: knowledge_articles
type KnowledgeArticle struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"type:varchar(255);not null" json:"title"`
	Content     string         `gorm:"type:text;not null" json:"content"`
	CategoryID  uint           `gorm:"not null;index" json:"category_id"`
	AuthorID    uint           `gorm:"not null;index" json:"author_id"`
	IsPublished bool           `gorm:"default:false;index" json:"is_published"` // Si l'article est publié
	ViewCount   int            `gorm:"default:0" json:"view_count"`             // Nombre de vues
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"` // Soft delete

	// Relations
	Category    KnowledgeCategory            `gorm:"foreignKey:CategoryID" json:"category,omitempty"`                               // Catégorie
	Author      User                         `gorm:"foreignKey:AuthorID" json:"author,omitempty"`                                   // Auteur
	Attachments []KnowledgeArticleAttachment `gorm:"foreignKey:ArticleID;constraint:OnDelete:CASCADE" json:"attachments,omitempty"` // Pièces jointes
}

// TableName spécifie le nom de la table
func (KnowledgeArticle) TableName() string {
	return "knowledge_articles"
}

// KnowledgeArticleAttachment représente une pièce jointe d'un article
// Table: knowledge_article_attachments
type KnowledgeArticleAttachment struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	ArticleID uint           `gorm:"not null;index" json:"article_id"`
	FileName  string         `gorm:"type:varchar(255);not null" json:"file_name"`
	FilePath  string         `gorm:"type:varchar(500);not null" json:"file_path"`
	FileSize  int            `gorm:"type:int" json:"file_size,omitempty"` // Taille en bytes
	MimeType  string         `gorm:"type:varchar(100)" json:"mime_type,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // Soft delete

	// Relations
	Article KnowledgeArticle `gorm:"foreignKey:ArticleID;constraint:OnDelete:CASCADE" json:"-"`
}

// TableName spécifie le nom de la table
func (KnowledgeArticleAttachment) TableName() string {
	return "knowledge_article_attachments"
}
