package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/scope"
	"github.com/mcicare/itsm-backend/internal/services"
	"github.com/mcicare/itsm-backend/internal/utils"
)

// TicketInternalHandler gère les handlers des tickets internes
type TicketInternalHandler struct {
	service services.TicketInternalService
}

// NewTicketInternalHandler crée une nouvelle instance
func NewTicketInternalHandler(service services.TicketInternalService) *TicketInternalHandler {
	return &TicketInternalHandler{service: service}
}

// canSeeTicketInternal retourne true si l'utilisateur peut voir le ticket interne selon le scope
func canSeeTicketInternal(s *scope.QueryScope, createdByID uint, assignedToID *uint, departmentID uint, filialeID uint) bool {
	if s == nil {
		return false
	}
	if s.HasPermission("tickets_internes.view_all") {
		return true
	}
	if s.HasPermission("tickets_internes.view_filiale") && s.FilialeID != nil && filialeID == *s.FilialeID {
		return true
	}
	if s.HasPermission("tickets_internes.view_department") && s.DepartmentID != nil && departmentID == *s.DepartmentID {
		return true
	}
	if s.HasPermission("tickets_internes.view_own") {
		if createdByID == s.UserID {
			return true
		}
		if assignedToID != nil && *assignedToID == s.UserID {
			return true
		}
	}
	return false
}

// GetMyPanier récupère le panier de l'utilisateur : tickets internes assignés à lui et non clôturés
// @Summary Récupérer mon panier (tickets internes)
// @Description Tickets internes assignés à l'utilisateur, non clôturés
// @Tags ticket-internes
// @Security BearerAuth
// @Produce json
// @Param page query int false "Numéro de page" default(1)
// @Param limit query int false "Nombre par page" default(50)
// @Success 200 {object} dto.TicketInternalListResponse
// @Failure 403 {object} utils.Response
// @Router /ticket-internes/panier [get]
func (h *TicketInternalHandler) GetMyPanier(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}
	userID := userIDVal.(uint)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 50
	}
	resp, err := h.service.GetPanier(userID, page, limit)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération du panier")
		return
	}
	utils.SuccessResponse(c, resp, "Panier récupéré avec succès")
}

// GetMyPerformance récupère la performance de l'utilisateur sur les tickets internes qu'il traite (assignés à lui)
// @Summary Ma performance sur les tickets internes
// @Description Retourne les métriques de performance sur les tickets internes assignés à l'utilisateur connecté (total, résolus, temps passé, efficacité)
// @Tags ticket-internes
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dto.TicketInternalPerformanceDTO
// @Failure 403 {object} utils.Response
// @Router /ticket-internes/performance/mine [get]
func (h *TicketInternalHandler) GetMyPerformance(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}
	userID := userIDVal.(uint)
	scope := utils.GetScopeFromContext(c)
	if scope == nil {
		utils.InternalServerErrorResponse(c, "Contexte utilisateur introuvable")
		return
	}
	if !scope.HasPermission("tickets_internes.view_own") && !scope.HasPermission("tickets_internes.view_department") &&
		!scope.HasPermission("tickets_internes.view_filiale") && !scope.HasPermission("tickets_internes.view_all") {
		utils.ErrorResponse(c, http.StatusForbidden, "Vous n'avez pas la permission de consulter les tickets internes", nil)
		return
	}
	perf, err := h.service.GetMyPerformance(userID)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors du calcul de la performance")
		return
	}
	utils.SuccessResponse(c, perf, "Performance récupérée avec succès")
}

// GetAll liste les tickets internes avec scope et pagination
// @Summary Lister les tickets internes
// @Description Liste les tickets internes avec pagination et filtres (scope selon permissions)
// @Tags ticket-internes
// @Security BearerAuth
// @Produce json
// @Param page query int false "Numéro de page" default(1)
// @Param limit query int false "Nombre par page" default(20)
// @Param status query string false "Filtrer par statut"
// @Param department_id query int false "Filtrer par département"
// @Param filiale_id query int false "Filtrer par filiale"
// @Success 200 {object} dto.TicketInternalListResponse
// @Failure 403 {object} utils.Response
// @Router /ticket-internes [get]
func (h *TicketInternalHandler) GetAll(c *gin.Context) {
	scope := utils.GetScopeFromContext(c)
	if scope == nil {
		utils.InternalServerErrorResponse(c, "Contexte utilisateur introuvable")
		return
	}
	if !scope.HasPermission("tickets_internes.view_own") && !scope.HasPermission("tickets_internes.view_department") &&
		!scope.HasPermission("tickets_internes.view_filiale") && !scope.HasPermission("tickets_internes.view_all") {
		utils.ErrorResponse(c, http.StatusForbidden, "Vous n'avez pas la permission de voir les tickets internes", nil)
		return
	}
	utils.ApplyDashboardScopeHint(c, scope)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	status := c.Query("status")
	var departmentID, filialeID *uint
	if d := c.Query("department_id"); d != "" {
		if id, err := strconv.ParseUint(d, 10, 32); err == nil {
			u := uint(id)
			departmentID = &u
		}
	}
	if f := c.Query("filiale_id"); f != "" {
		if id, err := strconv.ParseUint(f, 10, 32); err == nil {
			u := uint(id)
			filialeID = &u
		}
	}

	resp, err := h.service.GetAll(scope, page, limit, status, departmentID, filialeID)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des tickets internes")
		return
	}
	utils.SuccessResponse(c, resp, "Tickets internes récupérés avec succès")
}

// GetByID récupère un ticket interne par ID (avec vérification de visibilité)
// @Summary Récupérer un ticket interne
// @Description Récupère un ticket interne par son ID (accès selon scope)
// @Tags ticket-internes
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID du ticket interne"
// @Success 200 {object} dto.TicketInternalDTO
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /ticket-internes/{id} [get]
func (h *TicketInternalHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}
	scope := utils.GetScopeFromContext(c)
	if scope == nil {
		utils.InternalServerErrorResponse(c, "Contexte utilisateur introuvable")
		return
	}
	ticket, err := h.service.GetByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Ticket interne introuvable")
		return
	}
	var assignedID *uint
	if ticket.AssignedToID != nil {
		assignedID = ticket.AssignedToID
	}
	if !canSeeTicketInternal(scope, ticket.CreatedByID, assignedID, ticket.DepartmentID, ticket.FilialeID) {
		utils.ErrorResponse(c, http.StatusForbidden, "Vous n'avez pas accès à ce ticket interne", nil)
		return
	}
	utils.SuccessResponse(c, ticket, "Ticket interne récupéré avec succès")
}

// Create crée un ticket interne
// @Summary Créer un ticket interne
// @Description Crée un nouveau ticket interne (départements non-IT uniquement)
// @Tags ticket-internes
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateTicketInternalRequest true "Données du ticket interne"
// @Success 201 {object} dto.TicketInternalDTO
// @Failure 400 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /ticket-internes [post]
func (h *TicketInternalHandler) Create(c *gin.Context) {
	if !utils.RequirePermission(c, "tickets_internes.create") {
		utils.ForbiddenResponse(c, "Permission insuffisante: tickets_internes.create")
		return
	}
	var req dto.CreateTicketInternalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}
	createdByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}
	scope := utils.GetScopeFromContext(c)
	// Seul l'admin système (view_all) peut choisir n'importe quel département ; les autres uniquement le leur
	if scope != nil && !scope.HasPermission("tickets_internes.view_all") {
		if scope.DepartmentID == nil {
			utils.ForbiddenResponse(c, "Aucun département associé à votre compte")
			return
		}
		if req.DepartmentID != *scope.DepartmentID {
			utils.ForbiddenResponse(c, "Vous ne pouvez créer un ticket interne que dans votre département")
			return
		}
	}
	// Admin (view_all ou view_filiale) peut assigner à n'importe quel employé ; sinon uniquement membres du département
	// Seul l'admin système (view_all) peut assigner à n'importe quel employé ; les autres uniquement à un membre de leur département
	allowAssignAny := scope != nil && scope.HasPermission("tickets_internes.view_all")
	ticket, err := h.service.Create(req, createdByID.(uint), allowAssignAny)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	utils.CreatedResponse(c, ticket, "Ticket interne créé avec succès")
}

// Update met à jour un ticket interne
// @Summary Modifier un ticket interne
// @Description Met à jour les champs d'un ticket interne
// @Tags ticket-internes
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID du ticket interne"
// @Param request body dto.UpdateTicketInternalRequest true "Données à mettre à jour"
// @Success 200 {object} dto.TicketInternalDTO
// @Failure 400 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /ticket-internes/{id} [put]
func (h *TicketInternalHandler) Update(c *gin.Context) {
	if !utils.RequirePermission(c, "tickets_internes.update") {
		utils.ForbiddenResponse(c, "Permission insuffisante: tickets_internes.update")
		return
	}
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}
	scope := utils.GetScopeFromContext(c)
	existing, err := h.service.GetByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Ticket interne introuvable")
		return
	}
	if !canSeeTicketInternal(scope, existing.CreatedByID, existing.AssignedToID, existing.DepartmentID, existing.FilialeID) {
		utils.ErrorResponse(c, http.StatusForbidden, "Vous n'avez pas accès à ce ticket interne", nil)
		return
	}
	var req dto.UpdateTicketInternalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}
	updatedByID, _ := c.Get("user_id")
	ticket, err := h.service.Update(uint(id), req, updatedByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	utils.SuccessResponse(c, ticket, "Ticket interne mis à jour avec succès")
}

// Assign assigne un ticket interne
// @Summary Assigner un ticket interne
// @Description Assigne un ticket interne à un utilisateur
// @Tags ticket-internes
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID du ticket interne"
// @Param request body dto.AssignTicketInternalRequest true "ID de l'utilisateur à assigner"
// @Success 200 {object} dto.TicketInternalDTO
// @Failure 400 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /ticket-internes/{id}/assign [post]
func (h *TicketInternalHandler) Assign(c *gin.Context) {
	if !utils.RequirePermission(c, "tickets_internes.assign") {
		utils.ForbiddenResponse(c, "Permission insuffisante: tickets_internes.assign")
		return
	}
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}
	scope := utils.GetScopeFromContext(c)
	existing, err := h.service.GetByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Ticket interne introuvable")
		return
	}
	if !canSeeTicketInternal(scope, existing.CreatedByID, existing.AssignedToID, existing.DepartmentID, existing.FilialeID) {
		utils.ErrorResponse(c, http.StatusForbidden, "Vous n'avez pas accès à ce ticket interne", nil)
		return
	}
	var req dto.AssignTicketInternalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}
	assignedByID, _ := c.Get("user_id")
	allowAssignAny := scope != nil && scope.HasPermission("tickets_internes.view_all")
	ticket, err := h.service.Assign(uint(id), req, assignedByID.(uint), allowAssignAny)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	utils.SuccessResponse(c, ticket, "Ticket interne assigné avec succès")
}

// ChangeStatus change le statut d'un ticket interne
func (h *TicketInternalHandler) ChangeStatus(c *gin.Context) {
	if !utils.RequirePermission(c, "tickets_internes.update") {
		utils.ForbiddenResponse(c, "Permission insuffisante: tickets_internes.update")
		return
	}
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}
	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}
	status := req.Status
	scope := utils.GetScopeFromContext(c)
	existing, err := h.service.GetByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Ticket interne introuvable")
		return
	}
	if !canSeeTicketInternal(scope, existing.CreatedByID, existing.AssignedToID, existing.DepartmentID, existing.FilialeID) {
		utils.ErrorResponse(c, http.StatusForbidden, "Vous n'avez pas accès à ce ticket interne", nil)
		return
	}
	changedByID, _ := c.Get("user_id")
	ticket, err := h.service.ChangeStatus(uint(id), status, changedByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	utils.SuccessResponse(c, ticket, "Statut mis à jour avec succès")
}

// Validate valide un ticket interne (en_attente → resolu)
// @Summary Valider un ticket interne
// @Description Passe le ticket interne en statut résolu (validation)
// @Tags ticket-internes
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID du ticket interne"
// @Success 200 {object} dto.TicketInternalDTO
// @Failure 400 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /ticket-internes/{id}/validate [post]
func (h *TicketInternalHandler) Validate(c *gin.Context) {
	if !utils.RequirePermission(c, "tickets_internes.validate") {
		utils.ForbiddenResponse(c, "Permission insuffisante: tickets_internes.validate")
		return
	}
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}
	scope := utils.GetScopeFromContext(c)
	existing, err := h.service.GetByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Ticket interne introuvable")
		return
	}
	if !canSeeTicketInternal(scope, existing.CreatedByID, existing.AssignedToID, existing.DepartmentID, existing.FilialeID) {
		utils.ErrorResponse(c, http.StatusForbidden, "Vous n'avez pas accès à ce ticket interne", nil)
		return
	}
	validatedByID, _ := c.Get("user_id")
	ticket, err := h.service.Validate(uint(id), validatedByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	utils.SuccessResponse(c, ticket, "Ticket interne validé avec succès")
}

// Close clôture un ticket interne
func (h *TicketInternalHandler) Close(c *gin.Context) {
	if !utils.RequirePermission(c, "tickets_internes.close") {
		utils.ForbiddenResponse(c, "Permission insuffisante: tickets_internes.close")
		return
	}
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}
	scope := utils.GetScopeFromContext(c)
	existing, err := h.service.GetByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Ticket interne introuvable")
		return
	}
	if !canSeeTicketInternal(scope, existing.CreatedByID, existing.AssignedToID, existing.DepartmentID, existing.FilialeID) {
		utils.ErrorResponse(c, http.StatusForbidden, "Vous n'avez pas accès à ce ticket interne", nil)
		return
	}
	closedByID, _ := c.Get("user_id")
	ticket, err := h.service.Close(uint(id), closedByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	utils.SuccessResponse(c, ticket, "Ticket interne clôturé avec succès")
}

// Delete supprime un ticket interne
// @Summary Supprimer un ticket interne
// @Description Supprime définitivement un ticket interne
// @Tags ticket-internes
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID du ticket interne"
// @Success 200 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /ticket-internes/{id} [delete]
func (h *TicketInternalHandler) Delete(c *gin.Context) {
	if !utils.RequirePermission(c, "tickets_internes.delete") {
		utils.ForbiddenResponse(c, "Permission insuffisante: tickets_internes.delete")
		return
	}
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}
	scope := utils.GetScopeFromContext(c)
	existing, err := h.service.GetByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Ticket interne introuvable")
		return
	}
	if !canSeeTicketInternal(scope, existing.CreatedByID, existing.AssignedToID, existing.DepartmentID, existing.FilialeID) {
		utils.ErrorResponse(c, http.StatusForbidden, "Vous n'avez pas accès à ce ticket interne", nil)
		return
	}
	if err := h.service.Delete(uint(id)); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	utils.SuccessResponse(c, nil, "Ticket interne supprimé avec succès")
}
