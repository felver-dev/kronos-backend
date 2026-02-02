package services

import (
	"encoding/json"
	"errors"
	"log"
	"math"
	"strings"
	"time"

	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/models"
	"github.com/mcicare/itsm-backend/internal/repositories"
	"github.com/mcicare/itsm-backend/internal/websocket"
)

// NotificationListOpts options pour la liste d'historique des notifications
type NotificationListOpts struct {
	Page            int
	Limit           int
	IsRead          *bool
	DateFrom        *time.Time
	DateTo          *time.Time
	Search          string // recherche texte (titre, message)
	FilterUserID    *uint  // admin: filtrer par utilisateur
	FilterFilialeID *uint  // admin: filtrer par filiale
}

// NotificationService interface pour les opérations sur les notifications
type NotificationService interface {
	Create(userID uint, notificationType string, title string, message string, linkURL string, metadata map[string]any) error
	GetByID(id uint) (*dto.NotificationDTO, error)
	GetByUserID(userID uint) ([]dto.NotificationDTO, error)
	GetUnreadByUserID(userID uint) ([]dto.NotificationDTO, error)
	GetByType(userID uint, notificationType string) ([]dto.NotificationDTO, error)
	List(userID uint, opts NotificationListOpts) (*dto.NotificationListResponse, error)
	MarkAsRead(id uint, userID uint) error
	MarkAllAsRead(userID uint) error
	Delete(id uint, userID uint) error
	GetUnreadCount(userID uint) (int64, error)
}

// notificationService implémente NotificationService
type notificationService struct {
	notificationRepo repositories.NotificationRepository
	userRepo         repositories.UserRepository
	hub              *websocket.Hub // Hub WebSocket pour les notifications en temps réel
}

// NewNotificationService crée une nouvelle instance de NotificationService
func NewNotificationService(
	notificationRepo repositories.NotificationRepository,
	userRepo repositories.UserRepository,
	hub *websocket.Hub,
) NotificationService {
	return &notificationService{
		notificationRepo: notificationRepo,
		userRepo:         userRepo,
		hub:              hub,
	}
}

// Create crée une nouvelle notification
func (s *notificationService) Create(userID uint, notificationType string, title string, message string, linkURL string, metadata map[string]any) error {
	// Vérifier que l'utilisateur existe
	_, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("utilisateur destinataire introuvable")
	}

	// Convertir metadata en JSON si fourni
	var metadataJSON []byte
	if metadata != nil {
		metadataJSON, err = json.Marshal(metadata)
		if err != nil {
			return errors.New("erreur lors de la sérialisation des métadonnées")
		}
	}

	notification := &models.Notification{
		UserID:   userID,
		Type:     notificationType,
		Title:    title,
		Message:  message,
		LinkURL:  linkURL,
		Metadata: metadataJSON,
		IsRead:   false,
	}

	if err := s.notificationRepo.Create(notification); err != nil {
		return errors.New("erreur lors de la création de la notification")
	}

	// Envoyer la notification via WebSocket en temps réel
	if s.hub != nil {
		// Créer le DTO manuellement pour éviter de charger User depuis la DB
		metadataMap := make(map[string]any)
		if len(notification.Metadata) > 0 {
			json.Unmarshal(notification.Metadata, &metadataMap)
		}

		notificationDTO := dto.NotificationDTO{
			ID:        notification.ID,
			UserID:    notification.UserID,
			Type:      notification.Type,
			Title:     notification.Title,
			Message:   notification.Message,
			IsRead:    notification.IsRead,
			LinkURL:   notification.LinkURL,
			Metadata:  metadataMap,
			CreatedAt: notification.CreatedAt,
		}

		wsMessage := map[string]interface{}{
			"type":    "notification",
			"payload": notificationDTO,
		}
		s.hub.SendToUser(userID, wsMessage)
		log.Printf("📤 Notification WebSocket envoyée à l'utilisateur %d: %s", userID, notification.Title)
	}

	return nil
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
		return []dto.NotificationDTO{}, errors.New("erreur lors de la récupération des notifications")
	}

	var notificationDTOs []dto.NotificationDTO
	for _, notification := range notifications {
		notificationDTOs = append(notificationDTOs, s.notificationToDTO(&notification))
	}

	// S'assurer qu'on retourne toujours un slice vide [] au lieu de nil
	if notificationDTOs == nil {
		notificationDTOs = []dto.NotificationDTO{}
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

// List récupère l'historique des notifications avec filtres et pagination (pour la page historique)
func (s *notificationService) List(userID uint, opts NotificationListOpts) (*dto.NotificationListResponse, error) {
	page := opts.Page
	if page < 1 {
		page = 1
	}
	limit := opts.Limit
	if limit < 1 || limit > 100 {
		limit = 20
	}

	var notifications []models.Notification
	var total int64
	var err error

	search := strings.TrimSpace(opts.Search)
	if opts.FilterUserID != nil {
		// Admin: filtrer par utilisateur (et optionnellement filiale)
		notifications, total, err = s.notificationRepo.FindAllWithFilters(opts.FilterUserID, opts.FilterFilialeID, opts.IsRead, opts.DateFrom, opts.DateTo, search, page, limit)
	} else {
		// Utilisateur: ses propres notifications, avec filtre filiale optionnel (ma filiale = filiale sélectionnée)
		notifications, total, err = s.notificationRepo.FindByUserIDWithFilters(userID, opts.IsRead, opts.DateFrom, opts.DateTo, search, opts.FilterFilialeID, page, limit)
	}
	if err != nil {
		return nil, errors.New("erreur lors de la récupération de l'historique des notifications")
	}

	var unreadCount int64
	if opts.FilterUserID != nil {
		unreadCount, _ = s.notificationRepo.CountUnread(*opts.FilterUserID)
	} else if opts.FilterFilialeID == nil {
		unreadCount, _ = s.notificationRepo.CountUnread(userID)
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	dtos := make([]dto.NotificationDTO, 0, len(notifications))
	for i := range notifications {
		dtos = append(dtos, s.notificationToDTO(&notifications[i]))
	}

	return &dto.NotificationListResponse{
		Notifications: dtos,
		UnreadCount:   int(unreadCount),
		Total:         total,
		Page:          page,
		Limit:         limit,
		TotalPages:    totalPages,
	}, nil
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
