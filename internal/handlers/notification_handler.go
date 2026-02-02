package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/dto"
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

	// S'assurer qu'on retourne toujours un tableau vide [] au lieu de null
	if notifications == nil {
		notifications = []dto.NotificationDTO{}
	}

	utils.SuccessResponse(c, notifications, "Notifications récupérées avec succès")
}

// GetUnread récupère uniquement les notifications non lues (pour la cloche)
// @Summary Récupérer les notifications non lues
// @Description Récupère les notifications non lues de l'utilisateur connecté (affichage cloche)
// @Tags notifications
// @Security BearerAuth
// @Produce json
// @Success 200 {array} dto.NotificationDTO
// @Router /notifications/unread [get]
func (h *NotificationHandler) GetUnread(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	notifications, err := h.notificationService.GetUnreadByUserID(userID.(uint))
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des notifications")
		return
	}
	if notifications == nil {
		notifications = []dto.NotificationDTO{}
	}
	utils.SuccessResponse(c, notifications, "Notifications non lues récupérées avec succès")
}

// List récupère l'historique des notifications avec filtres et pagination (page historique)
// @Summary Historique des notifications
// @Description Liste des notifications avec filtres (utilisateur, filiale, date, lu/non lu) et pagination
// @Tags notifications
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page"
// @Param limit query int false "Limite par page"
// @Param is_read query bool false "Filtrer par lu (true) / non lu (false)"
// @Param date_from query string false "Date début (ISO)"
// @Param date_to query string false "Date fin (ISO)"
// @Param user_id query int false "Filtrer par utilisateur (admin)"
// @Param filiale_id query int false "Filtrer par filiale (admin)"
// @Success 200 {object} dto.NotificationListResponse
// @Router /notifications/history [get]
func (h *NotificationHandler) List(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}
	uid := userID.(uint)

	opts := services.NotificationListOpts{Page: 1, Limit: 20}
	if v := c.Query("page"); v != "" {
		if p, err := strconv.Atoi(v); err == nil && p >= 1 {
			opts.Page = p
		}
	}
	if v := c.Query("limit"); v != "" {
		if l, err := strconv.Atoi(v); err == nil && l >= 1 && l <= 100 {
			opts.Limit = l
		}
	}
	if v := c.Query("is_read"); v == "true" {
		t := true
		opts.IsRead = &t
	} else if v == "false" {
		f := false
		opts.IsRead = &f
	}
	if v := c.Query("date_from"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			opts.DateFrom = &t
		}
	}
	if v := c.Query("date_to"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			opts.DateTo = &t
		}
	}
	if v := c.Query("search"); v != "" {
		opts.Search = strings.TrimSpace(v)
	}
	// Filiale : tout le monde peut filtrer par filiale (mes notifications pour les non-admin, ou vue admin)
	if v := c.Query("filiale_id"); v != "" {
		if id, err := strconv.ParseUint(v, 10, 32); err == nil {
			f := uint(id)
			opts.FilterFilialeID = &f
		}
	}

	scope := utils.GetScopeFromContext(c)
	if scope != nil && scope.HasPermission("users.view_all") {
		if v := c.Query("user_id"); v != "" {
			if id, err := strconv.ParseUint(v, 10, 32); err == nil {
				u := uint(id)
				opts.FilterUserID = &u
			}
		}
	}

	resp, err := h.notificationService.List(uid, opts)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}
	utils.SuccessResponse(c, resp, "Historique récupéré avec succès")
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
