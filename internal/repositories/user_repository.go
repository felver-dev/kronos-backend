package repositories

import (
	"time"

	"github.com/mcicare/itsm-backend/database"
	"github.com/mcicare/itsm-backend/internal/models"
)

// UserRepository interface pour les opérations sur les utilisateurs
type UserRepository interface {
	Create(user *models.User) error
	FindByID(id uint) (*models.User, error)
	FindByUsername(username string) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	FindAll() ([]models.User, error)
	FindByRole(roleID uint) ([]models.User, error)
	FindActive() ([]models.User, error)
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

// Create crée un nouvel utilisateur
func (r *userRepository) Create(user *models.User) error {
	return database.DB.Create(user).Error
}

// FindByID trouve un utilisateur par son ID avec son rôle
func (r *userRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := database.DB.Preload("Role").First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByUsername trouve un utilisateur par son nom d'utilisateur avec son rôle
func (r *userRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	err := database.DB.Preload("Role").Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByEmail trouve un utilisateur par son email avec son rôle
func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := database.DB.Preload("Role").Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindAll récupère tous les utilisateurs avec leurs rôles
func (r *userRepository) FindAll() ([]models.User, error) {
	var users []models.User
	err := database.DB.Preload("Role").Find(&users).Error
	return users, err
}

// FindByRole récupère tous les utilisateurs d'un rôle donné
func (r *userRepository) FindByRole(roleID uint) ([]models.User, error) {
	var users []models.User
	err := database.DB.Preload("Role").Where("role_id = ?", roleID).Find(&users).Error
	return users, err
}

// FindActive récupère tous les utilisateurs actifs
func (r *userRepository) FindActive() ([]models.User, error) {
	var users []models.User
	err := database.DB.Preload("Role").Where("is_active = ?", true).Find(&users).Error
	return users, err
}

// Update met à jour un utilisateur
func (r *userRepository) Update(user *models.User) error {
	return database.DB.Save(user).Error
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

