package services

import (
	"encoding/json"
	"errors"

	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/models"
	"github.com/mcicare/itsm-backend/internal/repositories"
)

// NotificationService interface pour les opérations sur les notifications
type NotificationService interface {
	GetByID(id uint) (*dto.NotificationDTO, error)
	GetByUserID(userID uint) ([]dto.NotificationDTO, error)
	GetUnreadByUserID(userID uint) ([]dto.NotificationDTO, error)
	GetByType(userID uint, notificationType string) ([]dto.NotificationDTO, error)
	MarkAsRead(id uint, userID uint) error
	MarkAllAsRead(userID uint) error
	Delete(id uint, userID uint) error
	GetUnreadCount(userID uint) (int64, error)
}

// notificationService implémente NotificationService
type notificationService struct {
	notificationRepo repositories.NotificationRepository
	userRepo         repositories.UserRepository
}

// NewNotificationService crée une nouvelle instance de NotificationService
func NewNotificationService(
	notificationRepo repositories.NotificationRepository,
	userRepo repositories.UserRepository,
) NotificationService {
	return &notificationService{
		notificationRepo: notificationRepo,
		userRepo:         userRepo,
	}
}

// GetByID récupère une notification par son ID
func (s *notificationService) GetByID(id uint) (*dto.NotificationDTO, error) {
	notification, err := s.notificationRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("notification introuvable")
	}

	notificationDTO := s.notificationToDTO(notification)
	return &notificationDTO, nil
}

// GetByUserID récupère toutes les notifications d'un utilisateur
func (s *notificationService) GetByUserID(userID uint) ([]dto.NotificationDTO, error) {
	notifications, err := s.notificationRepo.FindByUserID(userID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des notifications")
	}

	var notificationDTOs []dto.NotificationDTO
	for _, notification := range notifications {
		notificationDTOs = append(notificationDTOs, s.notificationToDTO(&notification))
	}

	return notificationDTOs, nil
}

// GetUnreadByUserID récupère les notifications non lues d'un utilisateur
func (s *notificationService) GetUnreadByUserID(userID uint) ([]dto.NotificationDTO, error) {
	notifications, err := s.notificationRepo.FindUnreadByUserID(userID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des notifications")
	}

	var notificationDTOs []dto.NotificationDTO
	for _, notification := range notifications {
		notificationDTOs = append(notificationDTOs, s.notificationToDTO(&notification))
	}

	return notificationDTOs, nil
}

// GetByType récupère les notifications d'un type spécifique pour un utilisateur
func (s *notificationService) GetByType(userID uint, notificationType string) ([]dto.NotificationDTO, error) {
	notifications, err := s.notificationRepo.FindByType(userID, notificationType)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des notifications")
	}

	var notificationDTOs []dto.NotificationDTO
	for _, notification := range notifications {
		notificationDTOs = append(notificationDTOs, s.notificationToDTO(&notification))
	}

	return notificationDTOs, nil
}

// MarkAsRead marque une notification comme lue
func (s *notificationService) MarkAsRead(id uint, userID uint) error {
	notification, err := s.notificationRepo.FindByID(id)
	if err != nil {
		return errors.New("notification introuvable")
	}

	// Vérifier que la notification appartient à l'utilisateur
	if notification.UserID != userID {
		return errors.New("vous n'êtes pas autorisé à modifier cette notification")
	}

	if err := s.notificationRepo.MarkAsRead(id); err != nil {
		return errors.New("erreur lors de la mise à jour de la notification")
	}

	return nil
}

// MarkAllAsRead marque toutes les notifications d'un utilisateur comme lues
func (s *notificationService) MarkAllAsRead(userID uint) error {
	if err := s.notificationRepo.MarkAllAsRead(userID); err != nil {
		return errors.New("erreur lors de la mise à jour des notifications")
	}

	return nil
}

// Delete supprime une notification
func (s *notificationService) Delete(id uint, userID uint) error {
	notification, err := s.notificationRepo.FindByID(id)
	if err != nil {
		return errors.New("notification introuvable")
	}

	// Vérifier que la notification appartient à l'utilisateur
	if notification.UserID != userID {
		return errors.New("vous n'êtes pas autorisé à supprimer cette notification")
	}

	if err := s.notificationRepo.Delete(id); err != nil {
		return errors.New("erreur lors de la suppression de la notification")
	}

	return nil
}

// GetUnreadCount récupère le nombre de notifications non lues d'un utilisateur
func (s *notificationService) GetUnreadCount(userID uint) (int64, error) {
	count, err := s.notificationRepo.CountUnread(userID)
	if err != nil {
		return 0, errors.New("erreur lors du comptage des notifications")
	}

	return count, nil
}

// notificationToDTO convertit un modèle Notification en DTO
func (s *notificationService) notificationToDTO(notification *models.Notification) dto.NotificationDTO {
	// Convertir Metadata de datatypes.JSON en map[string]any
	var metadata map[string]any
	if len(notification.Metadata) > 0 {
		json.Unmarshal(notification.Metadata, &metadata)
	}

	notificationDTO := dto.NotificationDTO{
		ID:        notification.ID,
		UserID:    notification.UserID,
		Type:      notification.Type,
		Title:     notification.Title,
		Message:   notification.Message,
		IsRead:    notification.IsRead,
		LinkURL:   notification.LinkURL,
		Metadata:  metadata,
		CreatedAt: notification.CreatedAt,
	}

	if notification.ReadAt != nil {
		notificationDTO.ReadAt = notification.ReadAt
	}

	// Convertir l'utilisateur si présent
	if notification.User.ID != 0 {
		userDTO := s.userToDTO(&notification.User)
		notificationDTO.User = &userDTO
	}

	return notificationDTO
}

// userToDTO convertit un modèle User en DTO (méthode helper)
func (s *notificationService) userToDTO(user *models.User) dto.UserDTO {
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
