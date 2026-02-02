package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/services"
	"github.com/mcicare/itsm-backend/internal/utils"
)

// FilialeHandler gère les handlers des filiales
type FilialeHandler struct {
	filialeService services.FilialeService
}

// NewFilialeHandler crée une nouvelle instance de FilialeHandler
func NewFilialeHandler(filialeService services.FilialeService) *FilialeHandler {
	return &FilialeHandler{
		filialeService: filialeService,
	}
}

// Create crée une nouvelle filiale
// @Summary Créer une filiale
// @Description Crée une nouvelle filiale (nécessite filiales.create)
// @Tags filiales
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateFilialeRequest true "Données de la filiale"
// @Success 201 {object} dto.FilialeDTO
// @Failure 400 {object} utils.Response
// @Router /filiales [post]
func (h *FilialeHandler) Create(c *gin.Context) {
	// Vérifier la permission
	if !utils.RequirePermission(c, "filiales.create") {
		utils.ForbiddenResponse(c, "Permission insuffisante: filiales.create")
		return
	}

	var req dto.CreateFilialeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	filiale, err := h.filialeService.Create(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.CreatedResponse(c, filiale, "Filiale créée avec succès")
}

// GetAll récupère les filiales visibles par l'utilisateur
// @Summary Récupérer les filiales
// @Description Avec filiales.view_all ou filiales.manage : toutes les filiales. Avec filiales.view ou notifications.filter_by_filiale : uniquement sa filiale.
// @Tags filiales
// @Security BearerAuth
// @Produce json
// @Success 200 {array} dto.FilialeDTO
// @Failure 500 {object} utils.Response
// @Router /filiales [get]
func (h *FilialeHandler) GetAll(c *gin.Context) {
	hasView := utils.RequirePermission(c, "filiales.view") || utils.RequirePermission(c, "filiales.view_all")
	hasFilterNotifications := utils.RequirePermission(c, "notifications.filter_by_filiale")
	if !hasView && !hasFilterNotifications {
		utils.ForbiddenResponse(c, "Permission insuffisante: filiales.view, filiales.view_all ou notifications.filter_by_filiale")
		return
	}

	canViewAll := utils.RequirePermission(c, "filiales.view_all") || utils.RequirePermission(c, "filiales.manage")

	var filiales []dto.FilialeDTO
	var err error

	if canViewAll {
		filiales, err = h.filialeService.GetAll()
	} else {
		scope := utils.GetScopeFromContext(c)
		if scope != nil && scope.FilialeID != nil {
			var one *dto.FilialeDTO
			one, err = h.filialeService.GetByID(*scope.FilialeID)
			if err == nil {
				filiales = []dto.FilialeDTO{*one}
			}
		} else {
			filiales = []dto.FilialeDTO{}
			err = nil
		}
	}

	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des filiales")
		return
	}

	utils.SuccessResponse(c, filiales, "Filiales récupérées avec succès")
}

// GetActive récupère uniquement les filiales actives
// @Summary Récupérer les filiales actives
// @Description Récupère la liste des filiales actives (route publique pour l'inscription)
// @Tags filiales
// @Produce json
// @Success 200 {array} dto.FilialeDTO
// @Failure 500 {object} utils.Response
// @Router /filiales/active [get]
func (h *FilialeHandler) GetActive(c *gin.Context) {
	// Route publique - pas de vérification de permission nécessaire
	filiales, err := h.filialeService.GetActive()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des filiales actives")
		return
	}

	utils.SuccessResponse(c, filiales, "Filiales actives récupérées avec succès")
}

// GetByID récupère une filiale par son ID
// @Summary Récupérer une filiale par ID
// @Description Récupère une filiale par son identifiant (nécessite filiales.view)
// @Tags filiales
// @Security BearerAuth
// @Produce json
// @Param filiale_id path int true "ID de la filiale"
// @Success 200 {object} dto.FilialeDTO
// @Failure 404 {object} utils.Response
// @Router /filiales/{filiale_id} [get]
func (h *FilialeHandler) GetByID(c *gin.Context) {
	// Vérifier la permission
	if !utils.RequirePermission(c, "filiales.view") {
		utils.ForbiddenResponse(c, "Permission insuffisante: filiales.view")
		return
	}

	idParam := c.Param("filiale_id")
	if idParam == "" {
		idParam = c.Param("id") // Fallback pour compatibilité
	}
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	filiale, err := h.filialeService.GetByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Filiale introuvable")
		return
	}

	utils.SuccessResponse(c, filiale, "Filiale récupérée avec succès")
}

// GetByCode récupère une filiale par son code
// @Summary Récupérer une filiale par code
// @Description Récupère une filiale par son code (nécessite filiales.view)
// @Tags filiales
// @Security BearerAuth
// @Produce json
// @Param code path string true "Code de la filiale"
// @Success 200 {object} dto.FilialeDTO
// @Failure 404 {object} utils.Response
// @Router /filiales/code/{code} [get]
func (h *FilialeHandler) GetByCode(c *gin.Context) {
	// Vérifier la permission
	if !utils.RequirePermission(c, "filiales.view") {
		utils.ForbiddenResponse(c, "Permission insuffisante: filiales.view")
		return
	}

	code := c.Param("code")
	if code == "" {
		utils.BadRequestResponse(c, "Code invalide")
		return
	}

	filiale, err := h.filialeService.GetByCode(code)
	if err != nil {
		utils.NotFoundResponse(c, "Filiale introuvable")
		return
	}

	utils.SuccessResponse(c, filiale, "Filiale récupérée avec succès")
}

// GetSoftwareProvider récupère la filiale fournisseur de logiciels (is_software_provider=true)
// @Summary Récupérer la filiale fournisseur de logiciels
// @Description Récupère la filiale marquée comme fournisseur de logiciels / IT (nécessite filiales.view)
// @Tags filiales
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dto.FilialeDTO
// @Failure 404 {object} utils.Response
// @Router /filiales/software-provider [get]
func (h *FilialeHandler) GetSoftwareProvider(c *gin.Context) {
	// Vérifier la permission
	if !utils.RequirePermission(c, "filiales.view") {
		utils.ForbiddenResponse(c, "Permission insuffisante: filiales.view")
		return
	}

	filiale, err := h.filialeService.GetSoftwareProvider()
	if err != nil {
		utils.NotFoundResponse(c, "Filiale fournisseur de logiciels introuvable")
		return
	}

	utils.SuccessResponse(c, filiale, "Filiale fournisseur de logiciels récupérée avec succès")
}

// Update met à jour une filiale
// @Summary Mettre à jour une filiale
// @Description Met à jour une filiale (nécessite filiales.update)
// @Tags filiales
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param filiale_id path int true "ID de la filiale"
// @Param request body dto.UpdateFilialeRequest true "Données de mise à jour"
// @Success 200 {object} dto.FilialeDTO
// @Failure 400 {object} utils.Response
// @Router /filiales/{filiale_id} [put]
func (h *FilialeHandler) Update(c *gin.Context) {
	// Vérifier la permission
	if !utils.RequirePermission(c, "filiales.update") {
		utils.ForbiddenResponse(c, "Permission insuffisante: filiales.update")
		return
	}

	idParam := c.Param("filiale_id")
	if idParam == "" {
		idParam = c.Param("id") // Fallback pour compatibilité
	}
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	var req dto.UpdateFilialeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	filiale, err := h.filialeService.Update(uint(id), req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, filiale, "Filiale mise à jour avec succès")
}

// Delete supprime une filiale
// @Summary Supprimer une filiale
// @Description Supprime une filiale (nécessite filiales.manage)
// @Tags filiales
// @Security BearerAuth
// @Produce json
// @Param filiale_id path int true "ID de la filiale"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /filiales/{filiale_id} [delete]
func (h *FilialeHandler) Delete(c *gin.Context) {
	// Vérifier la permission
	if !utils.RequirePermission(c, "filiales.manage") {
		utils.ForbiddenResponse(c, "Permission insuffisante: filiales.manage")
		return
	}

	idParam := c.Param("filiale_id")
	if idParam == "" {
		idParam = c.Param("id") // Fallback pour compatibilité
	}
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	err = h.filialeService.Delete(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, nil, "Filiale supprimée avec succès")
}
