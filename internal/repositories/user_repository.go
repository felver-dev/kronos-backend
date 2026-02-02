package repositories

import (
	"fmt"
	"strings"
	"time"

	"github.com/mcicare/itsm-backend/database"
	"github.com/mcicare/itsm-backend/internal/models"
	"github.com/mcicare/itsm-backend/internal/scope"
	"gorm.io/gorm"
)

// UserRepository interface pour les opérations sur les utilisateurs
type UserRepository interface {
	Create(user *models.User) error
	FindByID(id uint) (*models.User, error)
	CountByIDs(ids []uint) (int64, error)
	FindByUsername(username string) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	FindAll(scope interface{}) ([]models.User, error) // scope peut être *scope.QueryScope ou nil
	FindByRole(scope interface{}, roleID uint) ([]models.User, error)
	FindActive(scope interface{}) ([]models.User, error)
	Search(scope interface{}, query string, limit int) ([]models.User, error) // scope peut être *scope.QueryScope ou nil
	CountByRole(roleID uint, count *int64) error
	Update(user *models.User) error
	Delete(id uint) error
	UpdateLastLogin(userID uint) error
}

// userRepository implémente UserRepository
type userRepository struct{}

// NewUserRepository crée une nouvelle instance de UserRepository
func NewUserRepository() UserRepository {
	return &userRepository{}
}

// applyUserPreloads applique les Preloads standards pour les utilisateurs
func applyUserPreloads(query *gorm.DB) *gorm.DB {
	return query.Preload("Role").
		Preload("Department").Preload("Department.Office").
		Preload("Filiale")
}

// Create crée un nouvel utilisateur
func (r *userRepository) Create(user *models.User) error {
	return database.DB.Create(user).Error
}

// FindByID trouve un utilisateur par son ID avec son rôle et département
func (r *userRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := applyUserPreloads(database.DB).
		First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// CountByIDs compte les utilisateurs par IDs (requête légère)
func (r *userRepository) CountByIDs(ids []uint) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	var count int64
	if err := database.DB.Model(&models.User{}).Where("id IN ?", ids).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// FindByUsername trouve un utilisateur par son nom d'utilisateur avec son rôle et département
func (r *userRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	err := applyUserPreloads(database.DB).
		Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByEmail trouve un utilisateur par son email avec son rôle et département
func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := database.DB.Preload("Role").Preload("Department").Preload("Department.Office").Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindAll récupère tous les utilisateurs avec leurs rôles et départements
// Le scope est utilisé pour filtrer automatiquement selon les permissions de l'utilisateur
func (r *userRepository) FindAll(scopeParam interface{}) ([]models.User, error) {
	var users []models.User
	
	// Construire la requête de base
	query := applyUserPreloads(database.DB.Model(&models.User{}))
	
	// Appliquer le scope si fourni
	if scopeParam != nil {
		if queryScope, ok := scopeParam.(*scope.QueryScope); ok {
			query = scope.ApplyUserScope(query, queryScope)
		}
	}
	
	err := query.Find(&users).Error
	return users, err
}

// FindByRole récupère tous les utilisateurs d'un rôle donné avec leurs départements
// Le scope est utilisé pour filtrer automatiquement selon les permissions de l'utilisateur
func (r *userRepository) FindByRole(scopeParam interface{}, roleID uint) ([]models.User, error) {
	var users []models.User
	
	// Construire la requête de base
	query := applyUserPreloads(database.DB.Model(&models.User{})).
		Where("role_id = ?", roleID)
	
	// Appliquer le scope si fourni
	if scopeParam != nil {
		if queryScope, ok := scopeParam.(*scope.QueryScope); ok {
			query = scope.ApplyUserScope(query, queryScope)
		}
	}
	
	err := query.Find(&users).Error
	return users, err
}

// FindActive récupère tous les utilisateurs actifs avec leurs départements
// Le scope est utilisé pour filtrer automatiquement selon les permissions de l'utilisateur
func (r *userRepository) FindActive(scopeParam interface{}) ([]models.User, error) {
	var users []models.User
	
	// Construire la requête de base
	query := applyUserPreloads(database.DB.Model(&models.User{})).
		Where("is_active = ?", true)
	
	// Appliquer le scope si fourni
	if scopeParam != nil {
		if queryScope, ok := scopeParam.(*scope.QueryScope); ok {
			query = scope.ApplyUserScope(query, queryScope)
		}
	}
	
	err := query.Find(&users).Error
	return users, err
}

// Search recherche des utilisateurs par nom, username ou email
func (r *userRepository) Search(scopeParam interface{}, query string, limit int) ([]models.User, error) {
	if limit <= 0 {
		limit = 20
	}
	like := "%" + strings.ToLower(query) + "%"
	
	// Construire la requête de base
	db := applyUserPreloads(database.DB.Model(&models.User{})).
		Where(
			"LOWER(users.username) LIKE ? OR LOWER(users.email) LIKE ? OR LOWER(users.first_name) LIKE ? OR LOWER(users.last_name) LIKE ?",
			like, like, like, like,
		)
	
	// Appliquer le scope si fourni (doit être fait avant les autres filtres)
	if scopeParam != nil {
		if queryScope, ok := scopeParam.(*scope.QueryScope); ok {
			db = scope.ApplyUserScope(db, queryScope)
		}
	}
	
	var users []models.User
	err := db.Limit(limit).Find(&users).Error
	return users, err
}

// CountByRole compte le nombre d'utilisateurs actifs (non supprimés et is_active = true) pour un rôle donné
func (r *userRepository) CountByRole(roleID uint, count *int64) error {
	return database.DB.Model(&models.User{}).
		Where("role_id = ? AND deleted_at IS NULL AND is_active = ?", roleID, true).
		Count(count).Error
}

// Update met à jour un utilisateur
func (r *userRepository) Update(user *models.User) error {
	// Log pour déboguer
	fmt.Printf("Repository Update - User ID: %d, RoleID à sauvegarder: %d\n", user.ID, user.RoleID)
	
	// Utiliser Where + Updates pour forcer la mise à jour de tous les champs, y compris role_id
	// Cela évite que GORM ignore role_id si Role est préchargé
	// On utilise Omit pour exclure les champs qu'on ne veut pas mettre à jour
	err := database.DB.Model(&models.User{}).
		Where("id = ?", user.ID).
		Omit("created_at", "created_by_id", "password_hash").
		Updates(map[string]interface{}{
			"username":      user.Username,
			"email":         user.Email,
			"first_name":    user.FirstName,
			"last_name":     user.LastName,
			"phone":         user.Phone,
			"avatar":        user.Avatar,
			"department_id": user.DepartmentID,
			"role_id":       user.RoleID, // Forcer la mise à jour du role_id
			"is_active":     user.IsActive,
			"updated_by_id": user.UpdatedByID,
			"updated_at":    time.Now(),
		}).Error
	
	if err != nil {
		fmt.Printf("Repository Update - Erreur: %v\n", err)
		return err
	}
	
	// Vérifier que le role_id a bien été mis à jour
	var checkUser models.User
	database.DB.Select("id", "role_id").First(&checkUser, user.ID)
	fmt.Printf("Repository Update - Après mise à jour, RoleID dans la DB: %d\n", checkUser.RoleID)
	
	return nil
}

// Delete supprime un utilisateur (soft delete)
func (r *userRepository) Delete(id uint) error {
	return database.DB.Delete(&models.User{}, id).Error
}

// UpdateLastLogin met à jour la date de dernière connexion
func (r *userRepository) UpdateLastLogin(userID uint) error {
	now := time.Now()
	return database.DB.Model(&models.User{}).Where("id = ?", userID).Update("last_login", now).Error
}

