package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/services"
	"github.com/mcicare/itsm-backend/internal/utils"
)

// DepartmentHandler gère les handlers des départements
type DepartmentHandler struct {
	departmentService services.DepartmentService
}

// NewDepartmentHandler crée une nouvelle instance de DepartmentHandler
func NewDepartmentHandler(departmentService services.DepartmentService) *DepartmentHandler {
	return &DepartmentHandler{
		departmentService: departmentService,
	}
}

// Create crée un nouveau département
// @Summary Créer un département
// @Description Crée un nouveau département
// @Tags departments
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateDepartmentRequest true "Données du département"
// @Success 201 {object} dto.DepartmentDTO
// @Failure 400 {object} utils.Response
// @Router /departments [post]
func (h *DepartmentHandler) Create(c *gin.Context) {
	var req dto.CreateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	// Vérifier si l'utilisateur peut créer un département dans n'importe quelle filiale
	canSelectAnyFiliale := utils.RequirePermission(c, "departments.create_any_filiale")

	// Si l'utilisateur n'a pas la permission de gérer les filiales,
	// il ne peut créer un département que dans sa propre filiale
	if !canSelectAnyFiliale {
		// Récupérer le scope qui contient déjà la filiale de l'utilisateur
		scope := utils.GetScopeFromContext(c)
		if scope != nil && scope.FilialeID != nil {
			// Si une filiale est spécifiée et qu'elle est différente de celle du créateur, refuser
			if req.FilialeID != nil && *req.FilialeID != *scope.FilialeID {
				utils.ForbiddenResponse(c, "Vous ne pouvez créer un département que dans votre propre filiale")
				return
			}

			// Forcer la filiale du créateur si aucune filiale n'est spécifiée
			if req.FilialeID == nil {
				req.FilialeID = scope.FilialeID
			}
		} else {
			// Si le scope n'a pas de filiale, refuser la création
			utils.ForbiddenResponse(c, "Impossible de déterminer votre filiale")
			return
		}
	}

	department, err := h.departmentService.Create(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.CreatedResponse(c, department, "Département créé avec succès")
}

// GetAll récupère tous les départements
// @Summary Récupérer tous les départements
// @Description Récupère la liste des départements selon le scope : departments.view_all = toutes les filiales ; departments.view_filiale ou departments.view = sa filiale uniquement.
// @Tags departments
// @Security BearerAuth
// @Produce json
// @Param active query bool false "Récupérer uniquement les départements actifs"
// @Success 200 {array} dto.DepartmentDTO
// @Failure 403 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /departments [get]
func (h *DepartmentHandler) GetAll(c *gin.Context) {
	activeOnly := c.DefaultQuery("active", "false") == "true"
	scope := utils.GetScopeFromContext(c)
	if scope == nil {
		utils.InternalServerErrorResponse(c, "Contexte utilisateur introuvable")
		return
	}

	// Accès : au moins une permission de vue requise
	hasViewAll := scope.HasPermission("departments.view_all")
	hasViewFiliale := scope.HasPermission("departments.view_filiale") || scope.HasPermission("departments.view")
	if !hasViewAll && !hasViewFiliale {
		utils.ErrorResponse(c, http.StatusForbidden, "Vous n'avez pas la permission de voir les départements", nil)
		return
	}

	// Vue globale : view_all ou rétrocompat (filiales.manage / create_any_filiale / update_any_filiale)
	canViewAll := hasViewAll ||
		scope.HasPermission("filiales.manage") ||
		scope.HasPermission("departments.create_any_filiale") ||
		scope.HasPermission("departments.update_any_filiale")

	var departments []dto.DepartmentDTO
	var err error

	if canViewAll {
		departments, err = h.departmentService.GetAll(activeOnly)
	} else {
		if scope.FilialeID != nil {
			departments, err = h.departmentService.GetByFilialeID(*scope.FilialeID)
		} else {
			departments = []dto.DepartmentDTO{}
		}
	}

	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des départements")
		return
	}

	utils.SuccessResponse(c, departments, "Départements récupérés avec succès")
}

// GetActive récupère uniquement les départements actifs (route publique)
// @Summary Récupérer les départements actifs
// @Description Récupère la liste des départements actifs (route publique pour l'inscription)
// @Tags departments
// @Produce json
// @Success 200 {array} dto.DepartmentDTO
// @Failure 500 {object} utils.Response
// @Router /departments/active [get]
func (h *DepartmentHandler) GetActive(c *gin.Context) {
	departments, err := h.departmentService.GetAll(true)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des départements")
		return
	}

	utils.SuccessResponse(c, departments, "Départements actifs récupérés avec succès")
}

// GetByID récupère un département par son ID
// @Summary Récupérer un département par ID
// @Description Récupère un département par son identifiant (scope : view_all = tout, view_filiale/view = même filiale uniquement)
// @Tags departments
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID du département"
// @Success 200 {object} dto.DepartmentDTO
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /departments/{id} [get]
func (h *DepartmentHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	department, err := h.departmentService.GetByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Département introuvable")
		return
	}

	scope := utils.GetScopeFromContext(c)
	if scope != nil && !scope.HasPermission("departments.view_all") &&
		!scope.HasPermission("filiales.manage") &&
		!scope.HasPermission("departments.create_any_filiale") &&
		!scope.HasPermission("departments.update_any_filiale") {
		if scope.FilialeID != nil && (department.FilialeID == nil || *department.FilialeID != *scope.FilialeID) {
			utils.ErrorResponse(c, http.StatusForbidden, "Vous n'avez pas accès à ce département", nil)
			return
		}
	}

	utils.SuccessResponse(c, department, "Département récupéré avec succès")
}

// Update met à jour un département
// @Summary Mettre à jour un département
// @Description Met à jour un département
// @Tags departments
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID du département"
// @Param request body dto.UpdateDepartmentRequest true "Données de mise à jour"
// @Success 200 {object} dto.DepartmentDTO
// @Failure 400 {object} utils.Response
// @Router /departments/{id} [put]
func (h *DepartmentHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	var req dto.UpdateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	// Vérifier si l'utilisateur peut modifier un département dans n'importe quelle filiale
	canSelectAnyFiliale := utils.RequirePermission(c, "departments.update_any_filiale")

	// Si l'utilisateur n'a pas la permission de modifier dans n'importe quelle filiale,
	// il ne peut modifier la filiale que pour utiliser la sienne
	if !canSelectAnyFiliale && req.FilialeID != nil {
		// Récupérer le scope qui contient la filiale de l'utilisateur
		scope := utils.GetScopeFromContext(c)
		if scope != nil && scope.FilialeID != nil {
			// Si une filiale différente de celle du modificateur est spécifiée, refuser
			if *req.FilialeID != *scope.FilialeID {
				utils.ForbiddenResponse(c, "Vous ne pouvez modifier la filiale que pour utiliser votre propre filiale")
				return
			}
		} else {
			// Si le scope n'a pas de filiale, refuser la modification
			utils.ForbiddenResponse(c, "Impossible de déterminer votre filiale")
			return
		}
	}

	department, err := h.departmentService.Update(uint(id), req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, department, "Département mis à jour avec succès")
}

// Delete supprime un département
// @Summary Supprimer un département
// @Description Supprime un département
// @Tags departments
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID du département"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /departments/{id} [delete]
func (h *DepartmentHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	err = h.departmentService.Delete(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, nil, "Département supprimé avec succès")
}

// GetByOfficeID récupère les départements d'un siège
// @Summary Récupérer les départements d'un siège
// @Description Récupère tous les départements d'un siège spécifique
// @Tags departments
// @Security BearerAuth
// @Produce json
// @Param office_id path int true "ID du siège"
// @Success 200 {array} dto.DepartmentDTO
// @Failure 500 {object} utils.Response
// @Router /departments/office/{office_id} [get]
func (h *DepartmentHandler) GetByOfficeID(c *gin.Context) {
	officeIDParam := c.Param("office_id")
	officeID, err := strconv.ParseUint(officeIDParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID du siège invalide")
		return
	}

	departments, err := h.departmentService.GetByOfficeID(uint(officeID))
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des départements")
		return
	}

	utils.SuccessResponse(c, departments, "Départements récupérés avec succès")
}

// GetByFilialeID récupère les départements d'une filiale
// @Summary Récupérer les départements d'une filiale
// @Description Récupère tous les départements actifs d'une filiale spécifique. Les filiales non fournisseur ne peuvent voir que leurs propres départements.
// @Tags departments
// @Security BearerAuth
// @Produce json
// @Param filiale_id path int true "ID de la filiale"
// @Success 200 {array} dto.DepartmentDTO
// @Failure 403 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /departments/filiale/{filiale_id} [get]
func (h *DepartmentHandler) GetByFilialeID(c *gin.Context) {
	filialeIDParam := c.Param("filiale_id")
	filialeID, err := strconv.ParseUint(filialeIDParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID de la filiale invalide")
		return
	}

	// Vérifier si l'utilisateur peut voir les départements de n'importe quelle filiale
	canViewAnyFiliale := utils.RequirePermission(c, "departments.view_all") || utils.RequirePermission(c, "filiales.manage")

	// Si l'utilisateur n'a pas la permission de voir toutes les filiales,
	// vérifier qu'il demande les départements de sa propre filiale
	if !canViewAnyFiliale {
		scope := utils.GetScopeFromContext(c)
		if scope == nil || scope.FilialeID == nil {
			utils.ForbiddenResponse(c, "Impossible de déterminer votre filiale")
			return
		}

		if uint(filialeID) != *scope.FilialeID {
			utils.ForbiddenResponse(c, "Vous ne pouvez voir que les départements de votre propre filiale")
			return
		}
	}

	departments, err := h.departmentService.GetByFilialeID(uint(filialeID))
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des départements")
		return
	}

	utils.SuccessResponse(c, departments, "Départements récupérés avec succès")
}
