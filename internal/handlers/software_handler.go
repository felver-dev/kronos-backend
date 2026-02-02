package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/services"
	"github.com/mcicare/itsm-backend/internal/utils"
)

// SoftwareHandler gère les handlers des logiciels
type SoftwareHandler struct {
	softwareService services.SoftwareService
}

// NewSoftwareHandler crée une nouvelle instance de SoftwareHandler
func NewSoftwareHandler(softwareService services.SoftwareService) *SoftwareHandler {
	return &SoftwareHandler{
		softwareService: softwareService,
	}
}

// Create crée un nouveau logiciel
// @Summary Créer un logiciel
// @Description Crée un nouveau logiciel (nécessite software.create)
// @Tags software
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateSoftwareRequest true "Données du logiciel"
// @Success 201 {object} dto.SoftwareDTO
// @Failure 400 {object} utils.Response
// @Router /software [post]
func (h *SoftwareHandler) Create(c *gin.Context) {
	// Vérifier la permission
	if !utils.RequirePermission(c, "software.create") {
		utils.ForbiddenResponse(c, "Permission insuffisante: software.create")
		return
	}

	var req dto.CreateSoftwareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	software, err := h.softwareService.Create(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.CreatedResponse(c, software, "Logiciel créé avec succès")
}

// GetAll récupère tous les logiciels
// @Summary Récupérer tous les logiciels
// @Description Récupère la liste de tous les logiciels (nécessite software.view)
// @Tags software
// @Security BearerAuth
// @Produce json
// @Success 200 {array} dto.SoftwareDTO
// @Failure 500 {object} utils.Response
// @Router /software [get]
func (h *SoftwareHandler) GetAll(c *gin.Context) {
	// Vérifier la permission
	if !utils.RequirePermission(c, "software.view") {
		utils.ForbiddenResponse(c, "Permission insuffisante: software.view")
		return
	}

	software, err := h.softwareService.GetAll()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des logiciels")
		return
	}

	utils.SuccessResponse(c, software, "Logiciels récupérés avec succès")
}

// GetActive récupère uniquement les logiciels actifs
// @Summary Récupérer les logiciels actifs
// @Description Récupère la liste des logiciels actifs (route publique pour la création de tickets)
// @Tags software
// @Produce json
// @Success 200 {array} dto.SoftwareDTO
// @Failure 500 {object} utils.Response
// @Router /software/active [get]
func (h *SoftwareHandler) GetActive(c *gin.Context) {
	// Route publique - pas de vérification de permission nécessaire
	software, err := h.softwareService.GetActive()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des logiciels actifs")
		return
	}

	utils.SuccessResponse(c, software, "Logiciels actifs récupérés avec succès")
}

// GetByID récupère un logiciel par son ID
// @Summary Récupérer un logiciel par ID
// @Description Récupère un logiciel par son identifiant (nécessite software.view)
// @Tags software
// @Security BearerAuth
// @Produce json
// @Param software_id path int true "ID du logiciel"
// @Success 200 {object} dto.SoftwareDTO
// @Failure 404 {object} utils.Response
// @Router /software/{software_id} [get]
func (h *SoftwareHandler) GetByID(c *gin.Context) {
	// Vérifier la permission
	if !utils.RequirePermission(c, "software.view") {
		utils.ForbiddenResponse(c, "Permission insuffisante: software.view")
		return
	}

	idParam := c.Param("software_id")
	if idParam == "" {
		idParam = c.Param("id") // Fallback pour compatibilité
	}
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	software, err := h.softwareService.GetByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Logiciel introuvable")
		return
	}

	utils.SuccessResponse(c, software, "Logiciel récupéré avec succès")
}

// GetByCode récupère un logiciel par son code
// @Summary Récupérer un logiciel par code
// @Description Récupère un logiciel par son code (nécessite software.view)
// @Tags software
// @Security BearerAuth
// @Produce json
// @Param code path string true "Code du logiciel"
// @Success 200 {object} dto.SoftwareDTO
// @Failure 404 {object} utils.Response
// @Router /software/code/{code} [get]
func (h *SoftwareHandler) GetByCode(c *gin.Context) {
	// Vérifier la permission
	if !utils.RequirePermission(c, "software.view") {
		utils.ForbiddenResponse(c, "Permission insuffisante: software.view")
		return
	}

	code := c.Param("code")
	if code == "" {
		utils.BadRequestResponse(c, "Code invalide")
		return
	}

	software, err := h.softwareService.GetByCode(code)
	if err != nil {
		utils.NotFoundResponse(c, "Logiciel introuvable")
		return
	}

	utils.SuccessResponse(c, software, "Logiciel récupéré avec succès")
}

// Update met à jour un logiciel
// @Summary Mettre à jour un logiciel
// @Description Met à jour un logiciel (nécessite software.update)
// @Tags software
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param software_id path int true "ID du logiciel"
// @Param request body dto.UpdateSoftwareRequest true "Données de mise à jour"
// @Success 200 {object} dto.SoftwareDTO
// @Failure 400 {object} utils.Response
// @Router /software/{software_id} [put]
func (h *SoftwareHandler) Update(c *gin.Context) {
	// Vérifier la permission
	if !utils.RequirePermission(c, "software.update") {
		utils.ForbiddenResponse(c, "Permission insuffisante: software.update")
		return
	}

	idParam := c.Param("software_id")
	if idParam == "" {
		idParam = c.Param("id") // Fallback pour compatibilité
	}
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	var req dto.UpdateSoftwareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	software, err := h.softwareService.Update(uint(id), req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, software, "Logiciel mis à jour avec succès")
}

// Delete supprime un logiciel
// @Summary Supprimer un logiciel
// @Description Supprime un logiciel (nécessite software.delete)
// @Tags software
// @Security BearerAuth
// @Produce json
// @Param software_id path int true "ID du logiciel"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /software/{software_id} [delete]
func (h *SoftwareHandler) Delete(c *gin.Context) {
	// Vérifier la permission
	if !utils.RequirePermission(c, "software.delete") {
		utils.ForbiddenResponse(c, "Permission insuffisante: software.delete")
		return
	}

	idParam := c.Param("software_id")
	if idParam == "" {
		idParam = c.Param("id") // Fallback pour compatibilité
	}
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	err = h.softwareService.Delete(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, nil, "Logiciel supprimé avec succès")
}
