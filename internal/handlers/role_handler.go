package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/services"
	"github.com/mcicare/itsm-backend/internal/utils"
)

// RoleHandler gère les handlers des rôles
type RoleHandler struct {
	roleService services.RoleService
}

// NewRoleHandler crée une nouvelle instance de RoleHandler
func NewRoleHandler(roleService services.RoleService) *RoleHandler {
	return &RoleHandler{
		roleService: roleService,
	}
}

// GetAll récupère les rôles visibles pour l'utilisateur selon ses permissions.
// @Summary Récupérer les rôles
// @Description roles.manage : tous les rôles. roles.view_department : rôles de son département. roles.view_filiale ou roles.view : rôles globaux + rôles de sa filiale.
// @Tags roles
// @Security BearerAuth
// @Produce json
// @Success 200 {array} dto.RoleDTO
// @Failure 500 {object} utils.Response
// @Router /roles [get]
func (h *RoleHandler) GetAll(c *gin.Context) {
	scope := utils.GetScopeFromContext(c)
	if scope == nil {
		utils.InternalServerErrorResponse(c, "Contexte utilisateur introuvable")
		return
	}
	canManageAll := scope.HasPermission("roles.manage")
	canViewDepartment := scope.HasPermission("roles.view_department")
	canViewFiliale := scope.HasPermission("roles.view_filiale")
	canView := scope.HasPermission("roles.view")
	canDelegateOnly := scope.HasPermission("roles.delegate_permissions")

	// Utilisateur qui n'a que delegate_permissions : peut récupérer uniquement les rôles qu'il a créés (pour les assigner à des utilisateurs)
	if canDelegateOnly && !canManageAll && !canViewDepartment && !canViewFiliale && !canView {
		roles, err := h.roleService.GetMyDelegations(scope.UserID)
		if err != nil {
			utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des rôles délégués")
			return
		}
		utils.SuccessResponse(c, roles, "Rôles récupérés avec succès")
		return
	}

	if !canManageAll && !canViewDepartment && !canViewFiliale && !canView {
		utils.ErrorResponse(c, http.StatusForbidden, "Vous n'avez pas la permission de voir les rôles", nil)
		return
	}

	viewMode := "filiale" // défaut : rôles globaux + filiale
	if canManageAll {
		viewMode = "all"
	} else if canViewDepartment && scope.DepartmentID != nil {
		viewMode = "department"
	} else if canViewFiliale || canView {
		viewMode = "filiale"
	}

	roles, err := h.roleService.GetAllForAssignment(scope.FilialeID, scope.DepartmentID, viewMode)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des rôles")
		return
	}
	// Inclure aussi les rôles délégués (créés par l'utilisateur) pour qu'ils soient sélectionnables dans le modal utilisateur
	delegated, errDeleg := h.roleService.GetMyDelegations(scope.UserID)
	if errDeleg == nil && len(delegated) > 0 {
		seen := make(map[uint]bool)
		for _, r := range roles {
			seen[r.ID] = true
		}
		for _, r := range delegated {
			if !seen[r.ID] {
				roles = append(roles, r)
				seen[r.ID] = true
			}
		}
	}
	utils.SuccessResponse(c, roles, "Rôles récupérés avec succès")
}

// GetByID récupère un rôle par son ID
// @Summary Récupérer un rôle par ID
// @Description Récupère un rôle par son identifiant
// @Tags roles
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID du rôle"
// @Success 200 {object} dto.RoleDTO
// @Failure 404 {object} utils.Response
// @Router /roles/{id} [get]
func (h *RoleHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	role, err := h.roleService.GetByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Rôle introuvable")
		return
	}

	utils.SuccessResponse(c, role, "Rôle récupéré avec succès")
}

// Create crée un nouveau rôle
// @Summary Créer un rôle
// @Description Crée un nouveau rôle dans le système
// @Tags roles
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateRoleRequest true "Données du rôle"
// @Success 201 {object} dto.RoleDTO
// @Failure 400 {object} utils.Response
// @Router /roles [post]
func (h *RoleHandler) Create(c *gin.Context) {
	var req dto.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	createdByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	role, err := h.roleService.Create(req, createdByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.CreatedResponse(c, role, "Rôle créé avec succès")
}

// Update met à jour un rôle
// @Summary Mettre à jour un rôle
// @Description Met à jour les informations d'un rôle
// @Tags roles
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID du rôle"
// @Param request body dto.UpdateRoleRequest true "Données à mettre à jour"
// @Success 200 {object} dto.RoleDTO
// @Failure 400 {object} utils.Response
// @Router /roles/{id} [put]
func (h *RoleHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	var req dto.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	updatedByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	canManageAllRoles := utils.RequirePermission(c, "roles.manage")
	role, err := h.roleService.Update(uint(id), req, updatedByID.(uint), canManageAllRoles)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, role, "Rôle mis à jour avec succès")
}

// Delete supprime un rôle
// @Summary Supprimer un rôle
// @Description Supprime un rôle du système (seul le créateur ou un utilisateur avec roles.manage peut supprimer)
// @Tags roles
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID du rôle"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /roles/{id} [delete]
func (h *RoleHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	deletedByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	canManageAllRoles := utils.RequirePermission(c, "roles.manage")
	err = h.roleService.Delete(uint(id), deletedByID.(uint), canManageAllRoles)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, nil, "Rôle supprimé avec succès")
}

// GetRolePermissions récupère les permissions d'un rôle
// @Summary Récupérer les permissions d'un rôle
// @Description Récupère la liste des permissions associées à un rôle
// @Tags roles
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID du rôle"
// @Success 200 {object} utils.Response{data=[]string}
// @Failure 404 {object} utils.Response
// @Router /roles/{id}/permissions [get]
func (h *RoleHandler) GetRolePermissions(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	permissions, err := h.roleService.GetRolePermissions(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, permissions, "Permissions récupérées avec succès")
}

// UpdateRolePermissions met à jour les permissions d'un rôle
// @Summary Mettre à jour les permissions d'un rôle
// @Description Met à jour la liste des permissions associées à un rôle
// @Tags roles
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID du rôle"
// @Param request body map[string][]string true "Liste des codes de permissions" example: {"permissions": ["tickets.view_all", "tickets.create"]}
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /roles/{id}/permissions [put]
func (h *RoleHandler) UpdateRolePermissions(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	var req struct {
		Permissions []string `json:"permissions" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	updatedByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	canManageAllRoles := utils.RequirePermission(c, "roles.manage")
	err = h.roleService.UpdateRolePermissions(uint(id), req.Permissions, updatedByID.(uint), canManageAllRoles)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, nil, "Permissions mises à jour avec succès")
}

// GetAssignablePermissions récupère les permissions que l'utilisateur actuel peut déléguer
// @Summary Récupérer les permissions assignables
// @Description Récupère la liste des permissions que l'utilisateur actuel peut assigner à d'autres rôles
// @Tags roles
// @Security BearerAuth
// @Produce json
// @Success 200 {object} utils.Response{data=[]string}
// @Failure 500 {object} utils.Response
// @Router /roles/assignable-permissions [get]
func (h *RoleHandler) GetAssignablePermissions(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	permissions, err := h.roleService.GetAssignablePermissions(userID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, permissions, "Permissions assignables récupérées avec succès")
}

// GetMyDelegations récupère les rôles créés par l'utilisateur courant
// @Summary Récupérer les rôles délégués
// @Description Récupère la liste des rôles créés par l'utilisateur courant (délégation)
// @Tags roles
// @Security BearerAuth
// @Produce json
// @Success 200 {object} utils.Response{data=[]dto.RoleDTO}
// @Failure 500 {object} utils.Response
// @Router /roles/my-delegations [get]
func (h *RoleHandler) GetMyDelegations(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	roles, err := h.roleService.GetMyDelegations(userID.(uint))
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des rôles délégués")
		return
	}

	utils.SuccessResponse(c, roles, "Rôles délégués récupérés avec succès")
}

// GetForDelegationPage récupère les rôles à afficher sur la page "Délégation des rôles"
// @Summary Rôles pour la page Délégation
// @Description Rôles créés par l'utilisateur + rôles utilisés par au moins un utilisateur de sa filiale (sans actions pour ces derniers)
// @Tags roles
// @Security BearerAuth
// @Produce json
// @Success 200 {object} utils.Response{data=[]dto.RoleDTO}
// @Failure 500 {object} utils.Response
// @Router /roles/for-delegation [get]
func (h *RoleHandler) GetForDelegationPage(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}
	scope := utils.GetScopeFromContext(c)
	var filialeID *uint
	if scope != nil {
		filialeID = scope.FilialeID
	}
	roles, err := h.roleService.GetForDelegationPage(userID.(uint), filialeID)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des rôles")
		return
	}
	utils.SuccessResponse(c, roles, "Rôles récupérés avec succès")
}
