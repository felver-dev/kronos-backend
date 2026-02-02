package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/services"
	"github.com/mcicare/itsm-backend/internal/utils"
)

// FilialeSoftwareHandler gère les handlers des déploiements de logiciels
type FilialeSoftwareHandler struct {
	deploymentService services.FilialeSoftwareService
}

// NewFilialeSoftwareHandler crée une nouvelle instance de FilialeSoftwareHandler
func NewFilialeSoftwareHandler(deploymentService services.FilialeSoftwareService) *FilialeSoftwareHandler {
	return &FilialeSoftwareHandler{
		deploymentService: deploymentService,
	}
}

// Create crée un nouveau déploiement
// @Summary Créer un déploiement
// @Description Crée un nouveau déploiement de logiciel chez une filiale (nécessite software.deploy)
// @Tags filiale-software
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param filiale_id path int true "ID de la filiale"
// @Param request body dto.CreateFilialeSoftwareRequest true "Données du déploiement"
// @Success 201 {object} dto.FilialeSoftwareDTO
// @Failure 400 {object} utils.Response
// @Router /filiales/{filiale_id}/software [post]
func (h *FilialeSoftwareHandler) Create(c *gin.Context) {
	// Vérifier la permission
	if !utils.RequirePermission(c, "software.deploy") {
		utils.ForbiddenResponse(c, "Permission insuffisante: software.deploy")
		return
	}

	// Récupérer filiale_id depuis l'URL
	filialeIDParam := c.Param("filiale_id")
	filialeID, err := strconv.ParseUint(filialeIDParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID de la filiale invalide")
		return
	}

	var req dto.CreateFilialeSoftwareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	// Utiliser le filiale_id de l'URL (prioritaire sur celui du body)
	req.FilialeID = uint(filialeID)

	deployment, err := h.deploymentService.Create(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.CreatedResponse(c, deployment, "Déploiement créé avec succès")
}

// GetByFilialeID récupère tous les déploiements d'une filiale
// @Summary Récupérer les déploiements d'une filiale
// @Description Récupère tous les déploiements actifs d'une filiale (accessible pour créer des tickets)
// @Tags filiale-software
// @Security BearerAuth
// @Produce json
// @Param filiale_id path int true "ID de la filiale"
// @Success 200 {array} dto.FilialeSoftwareDTO
// @Failure 500 {object} utils.Response
// @Router /filiales/{filiale_id}/software [get]
func (h *FilialeSoftwareHandler) GetByFilialeID(c *gin.Context) {
	// Permettre l'accès si l'utilisateur peut voir les logiciels OU créer des tickets
	// Cela permet aux utilisateurs de voir les logiciels déployés dans leur filiale pour créer des tickets
	canViewSoftware := utils.RequirePermission(c, "software.view")
	canCreateTickets := utils.RequirePermission(c, "tickets.create")

	if !canViewSoftware && !canCreateTickets {
		utils.ForbiddenResponse(c, "Permission insuffisante")
		return
	}

	filialeIDParam := c.Param("filiale_id")
	filialeID, err := strconv.ParseUint(filialeIDParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID de la filiale invalide")
		return
	}

	// Récupérer uniquement les déploiements actifs si l'utilisateur ne peut que créer des tickets
	// Sinon, retourner tous les déploiements
	deployments, err := h.deploymentService.GetByFilialeID(uint(filialeID))
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des déploiements")
		return
	}

	// Filtrer pour ne retourner que les déploiements actifs si l'utilisateur ne peut que créer des tickets
	if !canViewSoftware && canCreateTickets {
		activeDeployments, err := h.deploymentService.GetActiveByFiliale(uint(filialeID))
		if err != nil {
			utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des déploiements actifs")
			return
		}
		utils.SuccessResponse(c, activeDeployments, "Déploiements actifs récupérés avec succès")
		return
	}

	utils.SuccessResponse(c, deployments, "Déploiements récupérés avec succès")
}

// GetBySoftwareID récupère tous les déploiements d'un logiciel
// @Summary Récupérer les déploiements d'un logiciel
// @Description Récupère tous les déploiements d'un logiciel (nécessite software.view)
// @Tags filiale-software
// @Security BearerAuth
// @Produce json
// @Param software_id path int true "ID du logiciel"
// @Success 200 {array} dto.FilialeSoftwareDTO
// @Failure 500 {object} utils.Response
// @Router /software/{software_id}/deployments [get]
func (h *FilialeSoftwareHandler) GetBySoftwareID(c *gin.Context) {
	// Vérifier la permission
	if !utils.RequirePermission(c, "software.view") {
		utils.ForbiddenResponse(c, "Permission insuffisante: software.view")
		return
	}

	softwareIDParam := c.Param("software_id")
	softwareID, err := strconv.ParseUint(softwareIDParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID du logiciel invalide")
		return
	}

	deployments, err := h.deploymentService.GetBySoftwareID(uint(softwareID))
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des déploiements")
		return
	}

	utils.SuccessResponse(c, deployments, "Déploiements récupérés avec succès")
}

// GetByID récupère un déploiement par son ID
// @Summary Récupérer un déploiement par ID
// @Description Récupère un déploiement par son identifiant (nécessite software.view)
// @Tags filiale-software
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID du déploiement"
// @Success 200 {object} dto.FilialeSoftwareDTO
// @Failure 404 {object} utils.Response
// @Router /filiales-software/{id} [get]
func (h *FilialeSoftwareHandler) GetByID(c *gin.Context) {
	// Vérifier la permission
	if !utils.RequirePermission(c, "software.view") {
		utils.ForbiddenResponse(c, "Permission insuffisante: software.view")
		return
	}

	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	deployment, err := h.deploymentService.GetByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Déploiement introuvable")
		return
	}

	utils.SuccessResponse(c, deployment, "Déploiement récupéré avec succès")
}

// GetAll récupère tous les déploiements
// @Summary Récupérer tous les déploiements
// @Description Récupère tous les déploiements (nécessite software.view)
// @Tags filiale-software
// @Security BearerAuth
// @Produce json
// @Success 200 {array} dto.FilialeSoftwareDTO
// @Failure 500 {object} utils.Response
// @Router /filiales-software [get]
func (h *FilialeSoftwareHandler) GetAll(c *gin.Context) {
	// Vérifier la permission
	if !utils.RequirePermission(c, "software.view") {
		utils.ForbiddenResponse(c, "Permission insuffisante: software.view")
		return
	}

	deployments, err := h.deploymentService.GetAll()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des déploiements")
		return
	}

	utils.SuccessResponse(c, deployments, "Déploiements récupérés avec succès")
}

// GetActive récupère tous les déploiements actifs
// @Summary Récupérer les déploiements actifs
// @Description Récupère tous les déploiements actifs (nécessite software.view)
// @Tags filiale-software
// @Security BearerAuth
// @Produce json
// @Success 200 {array} dto.FilialeSoftwareDTO
// @Failure 500 {object} utils.Response
// @Router /filiales-software/active [get]
func (h *FilialeSoftwareHandler) GetActive(c *gin.Context) {
	// Vérifier la permission
	if !utils.RequirePermission(c, "software.view") {
		utils.ForbiddenResponse(c, "Permission insuffisante: software.view")
		return
	}

	deployments, err := h.deploymentService.GetActive()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des déploiements actifs")
		return
	}

	utils.SuccessResponse(c, deployments, "Déploiements actifs récupérés avec succès")
}

// Update met à jour un déploiement
// @Summary Mettre à jour un déploiement
// @Description Met à jour un déploiement (nécessite software.manage_deployments)
// @Tags filiale-software
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID du déploiement"
// @Param request body dto.UpdateFilialeSoftwareRequest true "Données de mise à jour"
// @Success 200 {object} dto.FilialeSoftwareDTO
// @Failure 400 {object} utils.Response
// @Router /filiales-software/{id} [put]
func (h *FilialeSoftwareHandler) Update(c *gin.Context) {
	// Vérifier la permission
	if !utils.RequirePermission(c, "software.manage_deployments") {
		utils.ForbiddenResponse(c, "Permission insuffisante: software.manage_deployments")
		return
	}

	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	var req dto.UpdateFilialeSoftwareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	deployment, err := h.deploymentService.Update(uint(id), req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, deployment, "Déploiement mis à jour avec succès")
}

// Delete supprime un déploiement
// @Summary Supprimer un déploiement
// @Description Supprime un déploiement (nécessite software.manage_deployments)
// @Tags filiale-software
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID du déploiement"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /filiales-software/{id} [delete]
func (h *FilialeSoftwareHandler) Delete(c *gin.Context) {
	// Vérifier la permission
	if !utils.RequirePermission(c, "software.manage_deployments") {
		utils.ForbiddenResponse(c, "Permission insuffisante: software.manage_deployments")
		return
	}

	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	err = h.deploymentService.Delete(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, nil, "Déploiement supprimé avec succès")
}
