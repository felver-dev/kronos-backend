package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/services"
	"github.com/mcicare/itsm-backend/internal/utils"
)

// SLAHandler gère les handlers des SLA
type SLAHandler struct {
	slaService services.SLAService
}

// NewSLAHandler crée une nouvelle instance de SLAHandler
func NewSLAHandler(slaService services.SLAService) *SLAHandler {
	return &SLAHandler{
		slaService: slaService,
	}
}

// Create crée un nouveau SLA
// @Summary Créer un SLA
// @Description Crée un nouveau SLA (Service Level Agreement) dans le système
// @Tags sla
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateSLARequest true "Données du SLA"
// @Success 201 {object} dto.SLADTO
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /sla [post]
func (h *SLAHandler) Create(c *gin.Context) {
	// Vérifier la permission
	if !utils.RequirePermission(c, "sla.create") {
		utils.ForbiddenResponse(c, "Permission insuffisante: sla.create")
		return
	}

	var req dto.CreateSLARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	createdByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	sla, err := h.slaService.Create(req, createdByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.CreatedResponse(c, sla, "SLA créé avec succès")
}

// GetByID récupère un SLA par son ID
// @Summary Récupérer un SLA par ID
// @Description Récupère un SLA par son identifiant
// @Tags sla
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID du SLA"
// @Success 200 {object} dto.SLADTO
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /sla/{id} [get]
func (h *SLAHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	sla, err := h.slaService.GetByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "SLA introuvable")
		return
	}

	utils.SuccessResponse(c, sla, "SLA récupéré avec succès")
}

// GetAll récupère tous les SLA
// @Summary Récupérer tous les SLA
// @Description Récupère la liste de tous les SLA
// @Tags sla
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {array} dto.SLADTO
// @Failure 500 {object} utils.Response
// @Router /sla [get]
func (h *SLAHandler) GetAll(c *gin.Context) {
	slas, err := h.slaService.GetAll()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des SLA")
		return
	}

	utils.SuccessResponse(c, slas, "SLA récupérés avec succès")
}

// GetTicketSLAStatus récupère le statut SLA d'un ticket
// @Summary Récupérer le statut SLA d'un ticket
// @Description Récupère le statut SLA d'un ticket spécifique
// @Tags sla
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param ticket_id path int true "ID du ticket"
// @Success 200 {object} dto.TicketSLAStatusDTO
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /sla/tickets/{ticket_id}/status [get]
func (h *SLAHandler) GetTicketSLAStatus(c *gin.Context) {
	ticketIDParam := c.Param("ticket_id")
	ticketID, err := strconv.ParseUint(ticketIDParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	status, err := h.slaService.GetTicketSLAStatus(uint(ticketID))
	if err != nil {
		utils.NotFoundResponse(c, "SLA introuvable pour ce ticket")
		return
	}

	utils.SuccessResponse(c, status, "Statut SLA récupéré avec succès")
}

// Update met à jour un SLA
// @Summary Mettre à jour un SLA
// @Description Met à jour les informations d'un SLA
// @Tags sla
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID du SLA"
// @Param request body dto.UpdateSLARequest true "Données à mettre à jour"
// @Success 200 {object} dto.SLADTO
// @Failure 400 {object} utils.Response
// @Router /sla/{id} [put]
func (h *SLAHandler) Update(c *gin.Context) {
	// Vérifier la permission
	if !utils.RequirePermission(c, "sla.update") {
		utils.ForbiddenResponse(c, "Permission insuffisante: sla.update")
		return
	}

	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	var req dto.UpdateSLARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	updatedByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	sla, err := h.slaService.Update(uint(id), req, updatedByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, sla, "SLA mis à jour avec succès")
}

// Delete supprime un SLA
// @Summary Supprimer un SLA
// @Description Supprime un SLA du système
// @Tags sla
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID du SLA"
// @Success 200 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /sla/{id} [delete]
func (h *SLAHandler) Delete(c *gin.Context) {
	// Vérifier la permission
	if !utils.RequirePermission(c, "sla.delete") {
		utils.ForbiddenResponse(c, "Permission insuffisante: sla.delete")
		return
	}

	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	err = h.slaService.Delete(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "SLA introuvable")
		return
	}

	utils.SuccessResponse(c, nil, "SLA supprimé avec succès")
}

// GetCompliance récupère le taux de conformité d'un SLA
// @Summary Récupérer le taux de conformité
// @Description Récupère le taux de conformité d'un SLA (pourcentage de tickets respectant le SLA)
// @Tags sla
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID du SLA"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} utils.Response
// @Router /sla/{id}/compliance [get]
func (h *SLAHandler) GetCompliance(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	compliance, err := h.slaService.GetCompliance(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "SLA introuvable")
		return
	}

	utils.SuccessResponse(c, compliance, "Taux de conformité récupéré avec succès")
}

// GetViolations récupère les violations d'un SLA
// @Summary Récupérer les violations d'un SLA
// @Description Récupère la liste des violations d'un SLA spécifique
// @Tags sla
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID du SLA"
// @Success 200 {array} dto.SLAViolationDTO
// @Failure 404 {object} utils.Response
// @Router /sla/{id}/violations [get]
func (h *SLAHandler) GetViolations(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	queryScope := utils.GetScopeFromContext(c)
	utils.ApplyDashboardScopeHint(c, queryScope)

	violations, err := h.slaService.GetViolations(queryScope, uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "SLA introuvable")
		return
	}

	utils.SuccessResponse(c, violations, "Violations récupérées avec succès")
}

// GetAllViolations récupère toutes les violations de SLA
// @Summary Récupérer toutes les violations
// @Description Récupère toutes les violations de SLA avec filtres optionnels
// @Tags sla
// @Security BearerAuth
// @Produce json
// @Param period query string false "Période (week, month)"
// @Param category query string false "Catégorie de ticket"
// @Success 200 {array} dto.SLAViolationDTO
// @Failure 500 {object} utils.Response
// @Router /sla/violations [get]
func (h *SLAHandler) GetAllViolations(c *gin.Context) {
	period := c.Query("period")
	category := c.Query("category")

	queryScope := utils.GetScopeFromContext(c)
	utils.ApplyDashboardScopeHint(c, queryScope)

	violations, err := h.slaService.GetAllViolations(queryScope, period, category)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des violations")
		return
	}

	utils.SuccessResponse(c, violations, "Violations récupérées avec succès")
}

// GetComplianceReport génère un rapport de conformité
// @Summary Générer un rapport de conformité
// @Description Génère un rapport de conformité des SLA au format PDF ou Excel
// @Tags sla
// @Security BearerAuth
// @Produce json
// @Param period query string false "Période (week, month)"
// @Param format query string false "Format (pdf, excel)"
// @Success 200 {file} file
// @Failure 500 {object} utils.Response
// @Router /sla/compliance-report [get]
func (h *SLAHandler) GetComplianceReport(c *gin.Context) {
	period := c.DefaultQuery("period", "month")
	format := c.DefaultQuery("format", "pdf")

	report, err := h.slaService.GetComplianceReport(period, format)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la génération du rapport")
		return
	}

	// Pour l'instant, on retourne les données JSON
	// TODO: Implémenter la génération de PDF/Excel
	utils.SuccessResponse(c, report, "Rapport généré avec succès")
}

// RecalculateSLAStatuses recalcule les statuts SLA pour tous les tickets ouverts
// @Summary Recalculer les statuts SLA
// @Description Recalcule les statuts SLA pour tous les tickets ouverts qui ont un SLA associé
// @Tags sla
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{} "Nombre de statuts mis à jour"
// @Failure 500 {object} utils.Response
// @Router /sla/recalculate [post]
func (h *SLAHandler) RecalculateSLAStatuses(c *gin.Context) {
	updatedCount, err := h.slaService.RecalculateSLAStatuses()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Erreur lors du recalcul des statuts SLA", err.Error())
		return
	}

	utils.SuccessResponse(c, map[string]interface{}{
		"updated_count": updatedCount,
		"message":       "Statuts SLA recalculés avec succès",
	}, "Statuts SLA recalculés avec succès")
}