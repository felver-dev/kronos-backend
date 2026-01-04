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

// GetAll récupère tous les rôles
// @Summary Récupérer tous les rôles
// @Description Récupère la liste de tous les rôles du système
// @Tags roles
// @Security BearerAuth
// @Produce json
// @Success 200 {array} dto.RoleDTO
// @Failure 500 {object} utils.Response
// @Router /roles [get]
func (h *RoleHandler) GetAll(c *gin.Context) {
	roles, err := h.roleService.GetAll()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des rôles")
		return
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

	role, err := h.roleService.Update(uint(id), req, updatedByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, role, "Rôle mis à jour avec succès")
}

// Delete supprime un rôle
// @Summary Supprimer un rôle
// @Description Supprime un rôle du système
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

	err = h.roleService.Delete(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, nil, "Rôle supprimé avec succès")
}
