package repositories

import (
	"time"

	"github.com/mcicare/itsm-backend/database"
	"github.com/mcicare/itsm-backend/internal/models"
)

// NotificationRepository interface pour les opérations sur les notifications
type NotificationRepository interface {
	Create(notification *models.Notification) error
	FindByID(id uint) (*models.Notification, error)
	FindByUserID(userID uint) ([]models.Notification, error)
	FindUnreadByUserID(userID uint) ([]models.Notification, error)
	FindByType(userID uint, notificationType string) ([]models.Notification, error)
	Update(notification *models.Notification) error
	MarkAsRead(id uint) error
	MarkAllAsRead(userID uint) error
	Delete(id uint) error
	CountUnread(userID uint) (int64, error)
}

// notificationRepository implémente NotificationRepository
type notificationRepository struct{}

// NewNotificationRepository crée une nouvelle instance de NotificationRepository
func NewNotificationRepository() NotificationRepository {
	return &notificationRepository{}
}

// Create crée une nouvelle notification
func (r *notificationRepository) Create(notification *models.Notification) error {
	return database.DB.Create(notification).Error
}

// FindByID trouve une notification par son ID
func (r *notificationRepository) FindByID(id uint) (*models.Notification, error) {
	var notification models.Notification
	err := database.DB.Preload("User").First(&notification, id).Error
	if err != nil {
		return nil, err
	}
	return &notification, nil
}

// FindByUserID récupère toutes les notifications d'un utilisateur
func (r *notificationRepository) FindByUserID(userID uint) ([]models.Notification, error) {
	var notifications []models.Notification
	err := database.DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&notifications).Error
	return notifications, err
}

// FindUnreadByUserID récupère les notifications non lues d'un utilisateur
func (r *notificationRepository) FindUnreadByUserID(userID uint) ([]models.Notification, error) {
	var notifications []models.Notification
	err := database.DB.Where("user_id = ? AND is_read = ?", userID, false).Order("created_at DESC").Find(&notifications).Error
	return notifications, err
}

// FindByType récupère les notifications d'un utilisateur par type
func (r *notificationRepository) FindByType(userID uint, notificationType string) ([]models.Notification, error) {
	var notifications []models.Notification
	err := database.DB.Where("user_id = ? AND type = ?", userID, notificationType).Order("created_at DESC").Find(&notifications).Error
	return notifications, err
}

// Update met à jour une notification
func (r *notificationRepository) Update(notification *models.Notification) error {
	return database.DB.Save(notification).Error
}

// MarkAsRead marque une notification comme lue
func (r *notificationRepository) MarkAsRead(id uint) error {
	now := time.Now()
	return database.DB.Model(&models.Notification{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_read": true,
		"read_at": now,
	}).Error
}

// MarkAllAsRead marque toutes les notifications d'un utilisateur comme lues
func (r *notificationRepository) MarkAllAsRead(userID uint) error {
	now := time.Now()
	return database.DB.Model(&models.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Updates(map[string]interface{}{
		"is_read": true,
		"read_at": now,
	}).Error
}

// Delete supprime une notification
func (r *notificationRepository) Delete(id uint) error {
	return database.DB.Delete(&models.Notification{}, id).Error
}

// CountUnread compte les notifications non lues d'un utilisateur
func (r *notificationRepository) CountUnread(userID uint) (int64, error) {
	var count int64
	err := database.DB.Model(&models.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Count(&count).Error
	return count, err
}

