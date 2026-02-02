package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/services"
	"github.com/mcicare/itsm-backend/internal/utils"
)

// PermissionHandler gère les handlers des permissions
type PermissionHandler struct {
	permissionService services.PermissionService
}

// NewPermissionHandler crée une nouvelle instance de PermissionHandler
func NewPermissionHandler(permissionService services.PermissionService) *PermissionHandler {
	return &PermissionHandler{
		permissionService: permissionService,
	}
}

// GetAll récupère toutes les permissions
// @Summary Récupérer toutes les permissions
// @Description Récupère la liste de toutes les permissions disponibles dans le système
// @Tags permissions
// @Security BearerAuth
// @Produce json
// @Param module query string false "Filtrer par module"
// @Success 200 {array} dto.PermissionDTO
// @Failure 500 {object} utils.Response
// @Router /permissions [get]
func (h *PermissionHandler) GetAll(c *gin.Context) {
	module := c.Query("module")

	var permissions []dto.PermissionDTO
	var err error

	if module != "" {
		permissions, err = h.permissionService.GetByModule(module)
	} else {
		permissions, err = h.permissionService.GetAll()
	}

	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des permissions")
		return
	}

	utils.SuccessResponse(c, permissions, "Permissions récupérées avec succès")
}

// GetByCode récupère une permission par son code
// @Summary Récupérer une permission par code
// @Description Récupère une permission par son code unique
// @Tags permissions
// @Security BearerAuth
// @Produce json
// @Param code path string true "Code de la permission"
// @Success 200 {object} dto.PermissionDTO
// @Failure 404 {object} utils.Response
// @Router /permissions/code/{code} [get]
func (h *PermissionHandler) GetByCode(c *gin.Context) {
	code := c.Param("code")

	permission, err := h.permissionService.GetByCode(code)
	if err != nil {
		utils.NotFoundResponse(c, "Permission introuvable")
		return
	}

	utils.SuccessResponse(c, permission, "Permission récupérée avec succès")
}
