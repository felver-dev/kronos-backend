package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/services"
	"github.com/mcicare/itsm-backend/internal/utils"
)

// TicketHandler gère les handlers des tickets
type TicketHandler struct {
	ticketService services.TicketService
}

// NewTicketHandler crée une nouvelle instance de TicketHandler
func NewTicketHandler(ticketService services.TicketService) *TicketHandler {
	return &TicketHandler{
		ticketService: ticketService,
	}
}

// Create crée un nouveau ticket
// @Summary Créer un ticket
// @Description Crée un nouveau ticket dans le système
// @Tags tickets
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateTicketRequest true "Données du ticket"
// @Success 201 {object} dto.TicketDTO
// @Failure 400 {object} utils.Response
// @Router /tickets [post]
func (h *TicketHandler) Create(c *gin.Context) {
	// Vérifier la permission de création
	if !utils.RequirePermission(c, "tickets.create") {
		utils.ForbiddenResponse(c, "Permission insuffisante: tickets.create")
		return
	}

	var req dto.CreateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	createdByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	// Vérifier si l'utilisateur peut créer un ticket pour n'importe quelle filiale
	// (permission tickets.create_any_filiale OU résolveur = département IT de la filiale fournisseur)
	scope := utils.GetScopeFromContext(c)
	canCreateAnyFiliale := utils.RequirePermission(c, "tickets.create_any_filiale") || (scope != nil && scope.IsResolver)

	// Si l'utilisateur n'a pas la permission de créer pour n'importe quelle filiale,
	// il ne peut créer un ticket que pour sa propre filiale
	if !canCreateAnyFiliale {
		// Utiliser le scope qui contient déjà la filiale de l'utilisateur
		scope := utils.GetScopeFromContext(c)
		if scope != nil && scope.FilialeID != nil {
			// Si une filiale est spécifiée et qu'elle est différente de celle du créateur, refuser
			if req.FilialeID != nil && *req.FilialeID != *scope.FilialeID {
				utils.ForbiddenResponse(c, "Vous ne pouvez créer un ticket que pour votre propre filiale")
				return
			}

			// Forcer la filiale du créateur si aucune filiale n'est spécifiée
			if req.FilialeID == nil {
				req.FilialeID = scope.FilialeID
			}
		} else {
			// Fallback: Si le scope n'a pas de filiale, récupérer l'utilisateur pour obtenir sa filiale
			// Note: On utilise le service pour éviter une dépendance circulaire
			// Mais on peut aussi utiliser le scope qui devrait toujours avoir la filiale
			utils.ForbiddenResponse(c, "Impossible de déterminer votre filiale")
			return
		}
	}

	ticket, err := h.ticketService.Create(req, createdByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.CreatedResponse(c, ticket, "Ticket créé avec succès")
}

// GetByID récupère un ticket par son ID
// @Summary Récupérer un ticket
// @Description Récupère les informations d'un ticket par son ID
// @Tags tickets
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID du ticket"
// @Success 200 {object} dto.TicketDTO
// @Failure 404 {object} utils.Response
// @Router /tickets/{id} [get]
func (h *TicketHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	includeDepartment := c.Query("include_department") == "true" || c.Query("include_department") == "1"
	ticket, err := h.ticketService.GetByID(uint(id), includeDepartment)
	if err != nil {
		utils.NotFoundResponse(c, "Ticket introuvable")
		return
	}

	utils.SuccessResponse(c, ticket, "Ticket récupéré avec succès")
}

// GetAll récupère tous les tickets avec pagination et filtres optionnels
// @Summary Liste des tickets
// @Description Récupère la liste des tickets avec pagination. Filtres optionnels: status, filiale_id, user_id (assigné)
// @Tags tickets
// @Security BearerAuth
// @Produce json
// @Param page query int false "Numéro de page" default(1)
// @Param limit query int false "Nombre d'éléments par page" default(20)
// @Param status query string false "Filtrer par statut (ouvert, en_cours, en_attente, resolu, cloture)"
// @Param filiale_id query int false "Filtrer par ID filiale"
// @Param user_id query int false "Filtrer par ID utilisateur assigné"
// @Success 200 {object} dto.TicketListResponse
// @Failure 500 {object} utils.Response
// @Router /tickets [get]
func (h *TicketHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	status := c.Query("status")
	filialeIDStr := c.Query("filiale_id")
	userIDStr := c.Query("user_id")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	queryScope := utils.GetScopeFromContext(c)
	utils.ApplyDashboardScopeHint(c, queryScope)

	var filialeID *uint
	if filialeIDStr != "" {
		if id, err := strconv.ParseUint(filialeIDStr, 10, 32); err == nil {
			uid := uint(id)
			filialeID = &uid
		}
	}
	var assigneeUserID *uint
	if userIDStr != "" {
		if id, err := strconv.ParseUint(userIDStr, 10, 32); err == nil {
			uid := uint(id)
			assigneeUserID = &uid
		}
	}

	var response interface{}
	var err error
	if status != "" || filialeID != nil || assigneeUserID != nil {
		response, err = h.ticketService.GetAllWithFilters(queryScope, page, limit, status, filialeID, assigneeUserID)
	} else {
		response, err = h.ticketService.GetAll(queryScope, page, limit)
	}
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des tickets")
		return
	}

	utils.SuccessResponse(c, response, "Tickets récupérés avec succès")
}

// GetByDepartment récupère les tickets par département du demandeur
// @Summary Liste des tickets par département
// @Description Récupère la liste des tickets dont le demandeur appartient au département donné
// @Tags tickets
// @Security BearerAuth
// @Produce json
// @Param departmentId path int true "ID du département"
// @Param page query int false "Numéro de page" default(1)
// @Param limit query int false "Nombre d'éléments par page" default(20)
// @Success 200 {object} dto.TicketListResponse
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /tickets/by-department/{departmentId} [get]
func (h *TicketHandler) GetByDepartment(c *gin.Context) {
	deptParam := c.Param("departmentId")
	deptID, err := strconv.ParseUint(deptParam, 10, 32)
	if err != nil || deptID == 0 {
		utils.BadRequestResponse(c, "ID de département invalide")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	response, err := h.ticketService.GetByDepartment(uint(deptID), page, limit)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des tickets par département")
		return
	}

	utils.SuccessResponse(c, response, "Tickets par département récupérés avec succès")
}

// Update met à jour un ticket
// @Summary Mettre à jour un ticket
// @Description Met à jour les informations d'un ticket
// @Tags tickets
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID du ticket"
// @Param request body dto.UpdateTicketRequest true "Données à mettre à jour"
// @Success 200 {object} dto.TicketDTO
// @Failure 400 {object} utils.Response
// @Router /tickets/{id} [put]
func (h *TicketHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	var req dto.UpdateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Log l'erreur de validation pour déboguer
		fmt.Printf("DEBUG: Erreur de validation lors de la mise à jour du ticket: %v\n", err)
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	// Log la requête reçue pour déboguer
	fmt.Printf("DEBUG: Requête de mise à jour reçue - Category: %s\n", req.Category)

	updatedByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	ticket, err := h.ticketService.Update(uint(id), req, updatedByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, ticket, "Ticket mis à jour avec succès")
}

// Assign assigne un ticket à un utilisateur
// @Summary Assigner un ticket
// @Description Assigne un ticket à un technicien
// @Tags tickets
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID du ticket"
// @Param request body dto.AssignTicketRequest true "Données d'assignation"
// @Success 200 {object} dto.TicketDTO
// @Failure 400 {object} utils.Response
// @Router /tickets/{id}/assign [post]
func (h *TicketHandler) Assign(c *gin.Context) {
	start := time.Now()
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	var req dto.AssignTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}
	fmt.Printf("DEBUG: Assign handler - ticket=%s user_ids=%v lead_id=%v\n", idParam, req.UserIDs, req.LeadID)

	assignedByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	ticket, err := h.ticketService.Assign(uint(id), req, assignedByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	fmt.Printf("PERF Assign ticket=%d users=%d dur=%s\n", id, len(req.UserIDs), time.Since(start))
	utils.SuccessResponse(c, ticket, "Ticket assigné avec succès")
}

// ChangeStatus change le statut d'un ticket
// @Summary Changer le statut d'un ticket
// @Description Change le statut d'un ticket
// @Tags tickets
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID du ticket"
// @Param request body map[string]string true "Nouveau statut"
// @Success 200 {object} dto.TicketDTO
// @Failure 400 {object} utils.Response
// @Router /tickets/{id}/status [put]
func (h *TicketHandler) ChangeStatus(c *gin.Context) {
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

	changedByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	ticket, err := h.ticketService.ChangeStatus(uint(id), req.Status, changedByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, ticket, "Statut modifié avec succès")
}

// Close ferme un ticket
// @Summary Fermer un ticket
// @Description Ferme un ticket
// @Tags tickets
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID du ticket"
// @Success 200 {object} dto.TicketDTO
// @Failure 400 {object} utils.Response
// @Router /tickets/{id}/close [post]
func (h *TicketHandler) Close(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	closedByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	ticket, err := h.ticketService.Close(uint(id), closedByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, ticket, "Ticket fermé avec succès")
}

// ValidateTicket valide un ticket en attente et le passe à « résolu »
// @Summary Valider un ticket en attente
// @Description Valide un ticket en attente de validation (le passe à « résolu »). Seuls les utilisateurs avec tickets.validate ou tickets.validate_own peuvent valider, ou le créateur du ticket.
// @Tags tickets
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID du ticket"
// @Success 200 {object} dto.TicketDTO
// @Failure 400 {object} utils.Response
// @Router /tickets/{id}/validate [post]
func (h *TicketHandler) ValidateTicket(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	validatedByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	// Vérifier la permission de validation OU si l'utilisateur est le créateur du ticket
	// D'abord récupérer le ticket pour vérifier si l'utilisateur est le créateur
	ticket, err := h.ticketService.GetByID(uint(id), false)
	if err != nil {
		utils.NotFoundResponse(c, "Ticket introuvable")
		return
	}

	// Autoriser si l'utilisateur a la permission tickets.validate OU tickets.validate_own OU est le créateur
	hasPermission := utils.RequireAnyPermission(c, "tickets.validate", "tickets.validate_own")
	isCreator := ticket.CreatedBy.ID != 0 && ticket.CreatedBy.ID == validatedByID.(uint)

	if !hasPermission && !isCreator {
		utils.ForbiddenResponse(c, "Permission insuffisante: tickets.validate ou tickets.validate_own, ou vous devez être le créateur du ticket")
		return
	}

	// Valider le ticket
	validatedTicket, err := h.ticketService.ValidateTicket(uint(id), validatedByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, validatedTicket, "Ticket validé et fermé avec succès")
}

// Delete supprime un ticket
// @Summary Supprimer un ticket
// @Description Supprime un ticket du système
// @Tags tickets
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID du ticket"
// @Success 200 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /tickets/{id} [delete]
func (h *TicketHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	err = h.ticketService.Delete(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Ticket introuvable")
		return
	}

	utils.SuccessResponse(c, nil, "Ticket supprimé avec succès")
}

// AddComment ajoute un commentaire à un ticket
// @Summary Ajouter un commentaire
// @Description Ajoute un commentaire à un ticket
// @Tags tickets
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID du ticket"
// @Param request body dto.CreateTicketCommentRequest true "Commentaire"
// @Success 201 {object} dto.TicketCommentDTO
// @Failure 400 {object} utils.Response
// @Router /tickets/{id}/comments [post]
func (h *TicketHandler) AddComment(c *gin.Context) {
	start := time.Now()
	idParam := c.Param("id")
	ticketID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	var req dto.CreateTicketCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	comment, err := h.ticketService.AddComment(uint(ticketID), req, userID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	fmt.Printf("PERF AddComment ticket=%d user=%d len=%d dur=%s\n", ticketID, userID.(uint), len(req.Comment), time.Since(start))
	utils.CreatedResponse(c, comment, "Commentaire ajouté avec succès")
}

// GetComments récupère les commentaires d'un ticket
// @Summary Récupérer les commentaires
// @Description Récupère tous les commentaires d'un ticket
// @Tags tickets
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID du ticket"
// @Success 200 {array} dto.TicketCommentDTO
// @Failure 404 {object} utils.Response
// @Router /tickets/{id}/comments [get]
func (h *TicketHandler) GetComments(c *gin.Context) {
	idParam := c.Param("id")
	ticketID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	scope := utils.GetScopeFromContext(c)
	canViewInternal := scope != nil && scope.DepartmentIsIT

	comments, err := h.ticketService.GetComments(uint(ticketID), canViewInternal)
	if err != nil {
		utils.NotFoundResponse(c, "Ticket introuvable")
		return
	}

	utils.SuccessResponse(c, comments, "Commentaires récupérés avec succès")
}

// UpdateComment met à jour un commentaire (seul l'auteur peut modifier).
// @Summary Modifier un commentaire
// @Description Met à jour le texte d'un commentaire. Seul l'auteur du commentaire peut le modifier.
// @Tags tickets
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID du ticket"
// @Param commentId path int true "ID du commentaire"
// @Param request body dto.UpdateTicketCommentRequest true "Nouveau texte du commentaire"
// @Success 200 {object} dto.TicketCommentDTO
// @Failure 400 {object} utils.Response
// @Router /tickets/{id}/comments/{commentId} [put]
func (h *TicketHandler) UpdateComment(c *gin.Context) {
	idParam := c.Param("id")
	ticketID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID du ticket invalide")
		return
	}
	commentIDParam := c.Param("commentId")
	commentID, err := strconv.ParseUint(commentIDParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID du commentaire invalide")
		return
	}
	var req dto.UpdateTicketCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}
	comment, err := h.ticketService.UpdateComment(uint(ticketID), uint(commentID), req, userID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	utils.SuccessResponse(c, comment, "Commentaire modifié avec succès")
}

// DeleteComment supprime un commentaire (seul l'auteur peut supprimer).
// @Summary Supprimer un commentaire
// @Description Supprime un commentaire. Seul l'auteur du commentaire peut le supprimer.
// @Tags tickets
// @Security BearerAuth
// @Param id path int true "ID du ticket"
// @Param commentId path int true "ID du commentaire"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /tickets/{id}/comments/{commentId} [delete]
func (h *TicketHandler) DeleteComment(c *gin.Context) {
	idParam := c.Param("id")
	ticketID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID du ticket invalide")
		return
	}
	commentIDParam := c.Param("commentId")
	commentID, err := strconv.ParseUint(commentIDParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID du commentaire invalide")
		return
	}
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}
	err = h.ticketService.DeleteComment(uint(ticketID), uint(commentID), userID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	utils.SuccessResponse(c, nil, "Commentaire supprimé avec succès")
}

// Reassign réassigne un ticket à un autre utilisateur
// @Summary Réassigner un ticket
// @Description Réassigne un ticket à un autre technicien
// @Tags tickets
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID du ticket"
// @Param request body dto.AssignTicketRequest true "Données de réassignation"
// @Success 200 {object} dto.TicketDTO
// @Failure 400 {object} utils.Response
// @Router /tickets/{id}/reassign [post]
func (h *TicketHandler) Reassign(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	var req dto.AssignTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	reassignedByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	ticket, err := h.ticketService.Assign(uint(id), req, reassignedByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, ticket, "Ticket réassigné avec succès")
}

// GetHistory récupère l'historique d'un ticket
// @Summary Récupérer l'historique d'un ticket
// @Description Récupère l'historique complet des modifications d'un ticket
// @Tags tickets
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID du ticket"
// @Success 200 {array} dto.TicketHistoryDTO
// @Failure 404 {object} utils.Response
// @Router /tickets/{id}/history [get]
func (h *TicketHandler) GetHistory(c *gin.Context) {
	idParam := c.Param("id")
	ticketID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	history, err := h.ticketService.GetHistory(uint(ticketID))
	if err != nil {
		utils.NotFoundResponse(c, "Ticket introuvable")
		return
	}

	utils.SuccessResponse(c, history, "Historique récupéré avec succès")
}

// GetBySource récupère les tickets par source
// @Summary Récupérer les tickets par source
// @Description Récupère les tickets filtrés par source (mail, appel, direct)
// @Tags tickets
// @Security BearerAuth
// @Produce json
// @Param source path string true "Source (mail, appel, direct)"
// @Param page query int false "Numéro de page" default(1)
// @Param limit query int false "Nombre d'éléments par page" default(20)
// @Success 200 {object} dto.TicketListResponse
// @Failure 400 {object} utils.Response
// @Router /tickets/by-source/{source} [get]
func (h *TicketHandler) GetBySource(c *gin.Context) {
	source := c.Param("source")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	// Extraire le QueryScope du contexte (injecté par AuthMiddleware)
	queryScope := utils.GetScopeFromContext(c)

	response, err := h.ticketService.GetBySource(queryScope, source, page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, response, "Tickets récupérés avec succès")
}

// GetByCategory récupère les tickets par catégorie
// @Summary Récupérer les tickets par catégorie
// @Description Récupère les tickets filtrés par catégorie (incident, demande, changement, developpement, assistance, support)
// @Tags tickets
// @Security BearerAuth
// @Produce json
// @Param category path string true "Catégorie (incident, demande, changement, developpement, assistance, support)"
// @Param page query int false "Numéro de page" default(1)
// @Param limit query int false "Nombre d'éléments par page" default(20)
// @Param status query string false "Filtrer par statut (ouvert, en_cours, en_attente, cloture)" default(all)
// @Param priority query string false "Filtrer par priorité (low, medium, high, critical)" default(all)
// @Success 200 {object} dto.TicketListResponse
// @Failure 400 {object} utils.Response
// @Router /tickets/by-category/{category} [get]
func (h *TicketHandler) GetByCategory(c *gin.Context) {
	category := c.Param("category")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	// Récupérer les paramètres de filtrage (peuvent être absents)
	statusParam := c.Query("status")
	priorityParam := c.Query("priority")

	// Utiliser "all" comme valeur par défaut si le paramètre est absent
	status := statusParam
	if status == "" {
		status = "all"
	}
	priority := priorityParam
	if priority == "" {
		priority = "all"
	}

	fmt.Printf("DEBUG: GetByCategory - category: %s, page: %d, limit: %d, status: '%s' (param: '%s'), priority: '%s' (param: '%s')\n",
		category, page, limit, status, statusParam, priority, priorityParam)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	// Extraire le QueryScope du contexte (injecté par AuthMiddleware)
	queryScope := utils.GetScopeFromContext(c)

	response, err := h.ticketService.GetByCategory(queryScope, category, page, limit, status, priority)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	fmt.Printf("DEBUG: GetByCategory - Nombre de tickets retournés: %d\n", len(response.Tickets))
	for i, ticket := range response.Tickets {
		if i < 3 { // Afficher les 3 premiers pour debug
			fmt.Printf("DEBUG: Ticket %d - ID: %d, Priority: %s\n", i+1, ticket.ID, ticket.Priority)
		}
	}

	utils.SuccessResponse(c, response, "Tickets récupérés avec succès")
}

// GetByStatus récupère les tickets par statut
// @Summary Récupérer les tickets par statut
// @Description Récupère les tickets filtrés par statut (ouvert, en_cours, en_attente, cloture)
// @Tags tickets
// @Security BearerAuth
// @Produce json
// @Param status path string true "Statut (ouvert, en_cours, en_attente, cloture)"
// @Param page query int false "Numéro de page" default(1)
// @Param limit query int false "Nombre d'éléments par page" default(20)
// @Success 200 {object} dto.TicketListResponse
// @Failure 400 {object} utils.Response
// @Router /tickets/by-status/{status} [get]
func (h *TicketHandler) GetByStatus(c *gin.Context) {
	status := c.Param("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	// Extraire le QueryScope du contexte (injecté par AuthMiddleware)
	queryScope := utils.GetScopeFromContext(c)

	response, err := h.ticketService.GetByStatus(queryScope, status, page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, response, "Tickets récupérés avec succès")
}

// GetByAssignee récupère les tickets assignés à un utilisateur
// @Summary Récupérer les tickets par assigné
// @Description Récupère les tickets assignés à un utilisateur spécifique
// @Tags tickets
// @Security BearerAuth
// @Produce json
// @Param userId path int true "ID de l'utilisateur"
// @Param page query int false "Numéro de page" default(1)
// @Param limit query int false "Nombre d'éléments par page" default(20)
// @Success 200 {object} dto.TicketListResponse
// @Failure 400 {object} utils.Response
// @Router /tickets/by-assignee/{userId} [get]
func (h *TicketHandler) GetByAssignee(c *gin.Context) {
	userIDParam := c.Param("userId")
	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID utilisateur invalide")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	response, err := h.ticketService.GetByAssignedTo(uint(userID), page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, response, "Tickets récupérés avec succès")
}

// GetMyPanier récupère le panier de l'utilisateur: tickets qui lui sont assignés et non clôturés
// @Summary Récupérer mon panier
// @Description Tickets assignés à l'utilisateur, non clôturés. À la clôture, un ticket disparaît du panier.
// @Tags tickets
// @Security BearerAuth
// @Produce json
// @Param page query int false "Numéro de page" default(1)
// @Param limit query int false "Nombre d'éléments par page" default(20)
// @Success 200 {object} dto.TicketListResponse
// @Failure 400 {object} utils.Response
// @Router /tickets/panier [get]
func (h *TicketHandler) GetMyPanier(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	response, err := h.ticketService.GetPanier(userID.(uint), page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	utils.SuccessResponse(c, response, "Panier récupéré avec succès")
}

// GetMyTickets récupère les tickets de l'utilisateur connecté
// @Summary Récupérer mes tickets
// @Description Récupère les tickets créés par l'utilisateur connecté ou qui lui sont assignés
// @Tags tickets
// @Security BearerAuth
// @Produce json
// @Param page query int false "Numéro de page" default(1)
// @Param limit query int false "Nombre d'éléments par page" default(20)
// @Success 200 {object} dto.TicketListResponse
// @Failure 400 {object} utils.Response
// @Router /tickets/my-tickets [get]
func (h *TicketHandler) GetMyTickets(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}
	start := time.Now()

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	status := c.Query("status")
	response, err := h.ticketService.GetByUser(userID.(uint), page, limit, status)
	fmt.Printf("PERF GetMyTickets user=%d page=%d limit=%d status=%s err=%v dur=%s\n", userID.(uint), page, limit, status, err, time.Since(start))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, response, "Tickets récupérés avec succès")
}
