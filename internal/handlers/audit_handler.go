package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/services"
	"github.com/mcicare/itsm-backend/internal/utils"
)

// AuditHandler gère les handlers des logs d'audit
type AuditHandler struct {
	auditService services.AuditService
}

// NewAuditHandler crée une nouvelle instance de AuditHandler
func NewAuditHandler(auditService services.AuditService) *AuditHandler {
	return &AuditHandler{
		auditService: auditService,
	}
}

// Types alias pour la doc Swagger (évite "cannot find type definition: dto.AuditLogListResponse")
type (
	auditLogListResponse = dto.AuditLogListResponse
	auditLogDTO          = dto.AuditLogDTO
)

// GetAll récupère tous les logs d'audit
// @Summary Liste des logs d'audit
// @Description Récupère la liste des logs d'audit avec pagination et filtres
// @Tags audit
// @Security BearerAuth
// @Produce json
// @Param page query int false "Numéro de page (défaut: 1)"
// @Param limit query int false "Nombre d'éléments par page (défaut: 50, max: 100)"
// @Param userId query int false "Filtrer par ID utilisateur"
// @Param action query string false "Filtrer par action"
// @Param entityType query string false "Filtrer par type d'entité"
// @Success 200 {object} auditLogListResponse
// @Failure 500 {object} utils.Response
// @Router /audit-logs [get]
func (h *AuditHandler) GetAll(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "50")
	userIDStr := c.Query("userId")
	action := c.Query("action")
	entityType := c.Query("entityType")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 50
	}

	var userID *uint
	if userIDStr != "" {
		id, err := strconv.ParseUint(userIDStr, 10, 32)
		if err == nil {
			uid := uint(id)
			userID = &uid
		}
	}

	// Extraire le QueryScope du contexte (injecté par AuthMiddleware)
	queryScope := utils.GetScopeFromContext(c)
	
	logs, err := h.auditService.GetAll(queryScope, page, limit, userID, action, entityType)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des logs d'audit")
		return
	}

	utils.SuccessResponse(c, logs, "Logs d'audit récupérés avec succès")
}

// GetByID récupère un log d'audit par son ID
// @Summary Détails d'un log d'audit
// @Description Récupère les détails d'un log d'audit par son ID
// @Tags audit
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID du log d'audit"
// @Success 200 {object} auditLogDTO
// @Failure 404 {object} utils.Response
// @Router /audit-logs/{id} [get]
func (h *AuditHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	log, err := h.auditService.GetByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Log d'audit introuvable")
		return
	}

	utils.SuccessResponse(c, log, "Log d'audit récupéré avec succès")
}

// GetByUserID récupère les logs d'audit d'un utilisateur
// @Summary Logs d'audit par utilisateur
// @Description Récupère tous les logs d'audit d'un utilisateur
// @Tags audit
// @Security BearerAuth
// @Produce json
// @Param userId path int true "ID de l'utilisateur"
// @Param startDate query string false "Date de début (format: 2006-01-02)"
// @Param endDate query string false "Date de fin (format: 2006-01-02)"
// @Success 200 {array} auditLogDTO
// @Failure 400 {object} utils.Response
// @Router /audit-logs/by-user/{userId} [get]
func (h *AuditHandler) GetByUserID(c *gin.Context) {
	userIDParam := c.Param("userId")
	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID utilisateur invalide")
		return
	}

	startDateStr := c.Query("startDate")
	endDateStr := c.Query("endDate")

	var startDate, endDate *time.Time
	if startDateStr != "" {
		if t, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = &t
		}
	}
	if endDateStr != "" {
		if t, err := time.Parse("2006-01-02", endDateStr); err == nil {
			// Ajouter 23h59 pour inclure toute la journée
			t = t.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
			endDate = &t
		}
	}

	// Extraire le QueryScope du contexte (injecté par AuthMiddleware)
	queryScope := utils.GetScopeFromContext(c)
	
	logs, err := h.auditService.GetByUserID(queryScope, uint(userID), startDate, endDate)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, logs, "Logs d'audit récupérés avec succès")
}

// GetByAction récupère les logs d'audit d'une action
// @Summary Logs d'audit par action
// @Description Récupère tous les logs d'audit d'un type d'action
// @Tags audit
// @Security BearerAuth
// @Produce json
// @Param action path string true "Type d'action (create, update, delete, etc.)"
// @Success 200 {array} auditLogDTO
// @Failure 500 {object} utils.Response
// @Router /audit-logs/by-action/{action} [get]
func (h *AuditHandler) GetByAction(c *gin.Context) {
	action := c.Param("action")
	if action == "" {
		utils.BadRequestResponse(c, "Action manquante")
		return
	}

	// Extraire le QueryScope du contexte (injecté par AuthMiddleware)
	queryScope := utils.GetScopeFromContext(c)
	
	logs, err := h.auditService.GetByAction(queryScope, action)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des logs d'audit")
		return
	}

	utils.SuccessResponse(c, logs, "Logs d'audit récupérés avec succès")
}

// GetByEntity récupère les logs d'audit d'une entité
// @Summary Logs d'audit par entité
// @Description Récupère tous les logs d'audit d'une entité spécifique
// @Tags audit
// @Security BearerAuth
// @Produce json
// @Param entityType path string true "Type d'entité (ticket, user, asset, etc.)"
// @Param entityId path int true "ID de l'entité"
// @Success 200 {array} auditLogDTO
// @Failure 400 {object} utils.Response
// @Router /audit-logs/by-entity/{entityType}/{entityId} [get]
func (h *AuditHandler) GetByEntity(c *gin.Context) {
	entityType := c.Param("entityType")
	entityIDParam := c.Param("entityId")

	entityID, err := strconv.ParseUint(entityIDParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID d'entité invalide")
		return
	}

	// Extraire le QueryScope du contexte (injecté par AuthMiddleware)
	queryScope := utils.GetScopeFromContext(c)
	
	logs, err := h.auditService.GetByEntity(queryScope, entityType, uint(entityID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, logs, "Logs d'audit récupérés avec succès")
}

// GetTicketAuditTrail récupère la piste d'audit d'un ticket
// @Summary Piste d'audit d'un ticket
// @Description Récupère la piste d'audit complète d'un ticket
// @Tags tickets
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID du ticket"
// @Success 200 {array} auditLogDTO
// @Failure 400 {object} utils.Response
// @Router /tickets/{id}/audit-trail [get]
func (h *AuditHandler) GetTicketAuditTrail(c *gin.Context) {
	ticketIDParam := c.Param("id")
	ticketID, err := strconv.ParseUint(ticketIDParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID de ticket invalide")
		return
	}

	// Extraire le QueryScope du contexte (injecté par AuthMiddleware)
	queryScope := utils.GetScopeFromContext(c)
	
	logs, err := h.auditService.GetTicketAuditTrail(queryScope, uint(ticketID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, logs, "Piste d'audit récupérée avec succès")
}

