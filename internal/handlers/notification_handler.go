package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/services"
	"github.com/mcicare/itsm-backend/internal/utils"
)

// NotificationHandler gère les handlers des notifications
type NotificationHandler struct {
	notificationService services.NotificationService
}

// NewNotificationHandler crée une nouvelle instance de NotificationHandler
func NewNotificationHandler(notificationService services.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
	}
}

// GetByUserID récupère les notifications d'un utilisateur
// @Summary Récupérer les notifications d'un utilisateur
// @Description Récupère toutes les notifications de l'utilisateur connecté
// @Tags notifications
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {array} dto.NotificationDTO
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /notifications [get]
func (h *NotificationHandler) GetByUserID(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	notifications, err := h.notificationService.GetByUserID(userID.(uint))
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des notifications")
		return
	}

	utils.SuccessResponse(c, notifications, "Notifications récupérées avec succès")
}

// MarkAsRead marque une notification comme lue
// @Summary Marquer une notification comme lue
// @Description Marque une notification spécifique comme lue
// @Tags notifications
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID de la notification"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /notifications/{id}/read [post]
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	err = h.notificationService.MarkAsRead(uint(id), userID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, nil, "Notification marquée comme lue")
}

// MarkAllAsRead marque toutes les notifications comme lues
// @Summary Marquer toutes les notifications comme lues
// @Description Marque toutes les notifications de l'utilisateur connecté comme lues
// @Tags notifications
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /notifications/read-all [post]
func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	err := h.notificationService.MarkAllAsRead(userID.(uint))
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la mise à jour des notifications")
		return
	}

	utils.SuccessResponse(c, nil, "Toutes les notifications ont été marquées comme lues")
}

// GetUnreadCount récupère le nombre de notifications non lues
// @Summary Récupérer le nombre de notifications non lues
// @Description Récupère le nombre de notifications non lues de l'utilisateur connecté
// @Tags notifications
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} dto.UnreadCountDTO
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /notifications/unread/count [get]
func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	count, err := h.notificationService.GetUnreadCount(userID.(uint))
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors du comptage des notifications")
		return
	}

	utils.SuccessResponse(c, gin.H{"count": count}, "Nombre de notifications non lues récupéré")
}
