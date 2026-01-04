package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/services"
	"github.com/mcicare/itsm-backend/internal/utils"
)

// AssetHandler gère les handlers des actifs IT
type AssetHandler struct {
	assetService services.AssetService
}

// NewAssetHandler crée une nouvelle instance de AssetHandler
func NewAssetHandler(assetService services.AssetService) *AssetHandler {
	return &AssetHandler{
		assetService: assetService,
	}
}

// Create crée un nouvel actif
// @Summary Créer un actif IT
// @Description Crée un nouvel actif IT dans le système
// @Tags assets
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateAssetRequest true "Données de l'actif"
// @Success 201 {object} dto.AssetDTO
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /assets [post]
func (h *AssetHandler) Create(c *gin.Context) {
	var req dto.CreateAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	createdByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	asset, err := h.assetService.Create(req, createdByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.CreatedResponse(c, asset, "Actif créé avec succès")
}

// GetByID récupère un actif par son ID
// @Summary Récupérer un actif par ID
// @Description Récupère un actif IT par son identifiant
// @Tags assets
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID de l'actif"
// @Success 200 {object} dto.AssetDTO
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /assets/{id} [get]
func (h *AssetHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	asset, err := h.assetService.GetByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Actif introuvable")
		return
	}

	utils.SuccessResponse(c, asset, "Actif récupéré avec succès")
}

// GetAll récupère tous les actifs
// @Summary Récupérer tous les actifs
// @Description Récupère la liste de tous les actifs IT
// @Tags assets
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {array} dto.AssetDTO
// @Failure 500 {object} utils.Response
// @Router /assets [get]
func (h *AssetHandler) GetAll(c *gin.Context) {
	assets, err := h.assetService.GetAll()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des actifs")
		return
	}

	utils.SuccessResponse(c, assets, "Actifs récupérés avec succès")
}

// Assign assigne un actif à un utilisateur
// @Summary Assigner un actif à un utilisateur
// @Description Assigne un actif IT à un utilisateur
// @Tags assets
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID de l'actif"
// @Param request body dto.AssignAssetRequest true "Données d'assignation"
// @Success 200 {object} dto.AssetDTO
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /assets/{id}/assign [post]
func (h *AssetHandler) Assign(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	var req dto.AssignAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	assignedByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	asset, err := h.assetService.Assign(uint(id), req, assignedByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, asset, "Actif assigné avec succès")
}

// Delete supprime un actif
// @Summary Supprimer un actif
// @Description Supprime un actif IT du système
// @Tags assets
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID de l'actif"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /assets/{id} [delete]
func (h *AssetHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	err = h.assetService.Delete(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Actif introuvable")
		return
	}

	utils.SuccessResponse(c, nil, "Actif supprimé avec succès")
}

// Update met à jour un actif
// @Summary Mettre à jour un actif
// @Description Met à jour les informations d'un actif IT
// @Tags assets
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID de l'actif"
// @Param request body dto.UpdateAssetRequest true "Données à mettre à jour"
// @Success 200 {object} dto.AssetDTO
// @Failure 400 {object} utils.Response
// @Router /assets/{id} [put]
func (h *AssetHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	var req dto.UpdateAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	updatedByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	asset, err := h.assetService.Update(uint(id), req, updatedByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, asset, "Actif mis à jour avec succès")
}

// Unassign retire l'assignation d'un actif
// @Summary Retirer l'assignation d'un actif
// @Description Retire l'assignation d'un actif IT à un utilisateur
// @Tags assets
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID de l'actif"
// @Success 200 {object} dto.AssetDTO
// @Failure 400 {object} utils.Response
// @Router /assets/{id}/unassign-user [delete]
func (h *AssetHandler) Unassign(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	unassignedByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	// Créer une requête avec UserID = 0 pour indiquer la désassignation
	req := dto.AssignAssetRequest{UserID: 0}
	asset, err := h.assetService.Unassign(uint(id), req, unassignedByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, asset, "Assignation retirée avec succès")
}

// GetAssignedUser récupère l'utilisateur assigné à un actif
// @Summary Récupérer l'utilisateur assigné
// @Description Récupère l'utilisateur assigné à un actif IT
// @Tags assets
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID de l'actif"
// @Success 200 {object} dto.UserDTO
// @Failure 404 {object} utils.Response
// @Router /assets/{id}/assigned-user [get]
func (h *AssetHandler) GetAssignedUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	asset, err := h.assetService.GetByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Actif introuvable")
		return
	}

	if asset.AssignedUser == nil {
		utils.NotFoundResponse(c, "Aucun utilisateur assigné à cet actif")
		return
	}

	utils.SuccessResponse(c, asset.AssignedUser, "Utilisateur assigné récupéré avec succès")
}

// GetByCategory récupère les actifs d'une catégorie
// @Summary Récupérer les actifs par catégorie
// @Description Récupère les actifs IT d'une catégorie spécifique
// @Tags assets
// @Security BearerAuth
// @Produce json
// @Param categoryId path int true "ID de la catégorie"
// @Success 200 {array} dto.AssetDTO
// @Failure 400 {object} utils.Response
// @Router /assets/by-category/{categoryId} [get]
func (h *AssetHandler) GetByCategory(c *gin.Context) {
	categoryIDParam := c.Param("categoryId")
	categoryID, err := strconv.ParseUint(categoryIDParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID catégorie invalide")
		return
	}

	assets, err := h.assetService.GetByCategory(uint(categoryID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, assets, "Actifs récupérés avec succès")
}

// GetByUser récupère les actifs assignés à un utilisateur
// @Summary Récupérer les actifs par utilisateur
// @Description Récupère les actifs IT assignés à un utilisateur spécifique
// @Tags assets
// @Security BearerAuth
// @Produce json
// @Param userId path int true "ID de l'utilisateur"
// @Success 200 {array} dto.AssetDTO
// @Failure 400 {object} utils.Response
// @Router /assets/by-user/{userId} [get]
func (h *AssetHandler) GetByUser(c *gin.Context) {
	userIDParam := c.Param("userId")
	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID utilisateur invalide")
		return
	}

	assets, err := h.assetService.GetByAssignedTo(uint(userID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, assets, "Actifs récupérés avec succès")
}

// GetInventory récupère l'inventaire des actifs
// @Summary Récupérer l'inventaire des actifs
// @Description Récupère l'inventaire complet des actifs IT avec statistiques
// @Tags assets
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dto.AssetInventoryDTO
// @Failure 500 {object} utils.Response
// @Router /assets/inventory [get]
func (h *AssetHandler) GetInventory(c *gin.Context) {
	inventory, err := h.assetService.GetInventory()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération de l'inventaire")
		return
	}

	utils.SuccessResponse(c, inventory, "Inventaire récupéré avec succès")
}

// GetLinkedTickets récupère les tickets liés à un actif
// @Summary Récupérer les tickets liés
// @Description Récupère la liste des tickets liés à un actif IT
// @Tags assets
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID de l'actif"
// @Success 200 {array} dto.TicketDTO
// @Failure 404 {object} utils.Response
// @Router /assets/{id}/tickets [get]
func (h *AssetHandler) GetLinkedTickets(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	tickets, err := h.assetService.GetLinkedTickets(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, tickets, "Tickets liés récupérés avec succès")
}

// LinkTicket lie un ticket à un actif
// @Summary Lier un ticket à un actif
// @Description Lie un ticket à un actif IT
// @Tags assets
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID de l'actif"
// @Param ticketId path int true "ID du ticket"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /assets/{id}/link-ticket/{ticketId} [post]
func (h *AssetHandler) LinkTicket(c *gin.Context) {
	idParam := c.Param("id")
	assetID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID actif invalide")
		return
	}

	ticketIDParam := c.Param("ticketId")
	ticketID, err := strconv.ParseUint(ticketIDParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID ticket invalide")
		return
	}

	linkedByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	err = h.assetService.LinkTicket(uint(assetID), uint(ticketID), linkedByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, nil, "Ticket lié avec succès")
}

// UnlinkTicket supprime la liaison entre un ticket et un actif
// @Summary Délier un ticket d'un actif
// @Description Supprime la liaison entre un ticket et un actif IT
// @Tags assets
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID de l'actif"
// @Param ticketId path int true "ID du ticket"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /assets/{id}/unlink-ticket/{ticketId} [delete]
func (h *AssetHandler) UnlinkTicket(c *gin.Context) {
	idParam := c.Param("id")
	assetID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID actif invalide")
		return
	}

	ticketIDParam := c.Param("ticketId")
	ticketID, err := strconv.ParseUint(ticketIDParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID ticket invalide")
		return
	}

	err = h.assetService.UnlinkTicket(uint(assetID), uint(ticketID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, nil, "Liaison supprimée avec succès")
}
