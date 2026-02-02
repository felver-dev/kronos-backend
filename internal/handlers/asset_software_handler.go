package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/services"
	"github.com/mcicare/itsm-backend/internal/utils"
)

// AssetSoftwareHandler gère les handlers des logiciels installés sur les actifs
type AssetSoftwareHandler struct {
	assetSoftwareService services.AssetSoftwareService
}

// NewAssetSoftwareHandler crée une nouvelle instance de AssetSoftwareHandler
func NewAssetSoftwareHandler(assetSoftwareService services.AssetSoftwareService) *AssetSoftwareHandler {
	return &AssetSoftwareHandler{
		assetSoftwareService: assetSoftwareService,
	}
}

// Create crée un nouveau logiciel installé
// @Summary Créer un logiciel installé
// @Description Crée un nouveau logiciel installé sur un actif
// @Tags assets
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateAssetSoftwareRequest true "Données du logiciel installé"
// @Success 201 {object} dto.AssetSoftwareDTO
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /assets/software [post]
func (h *AssetSoftwareHandler) Create(c *gin.Context) {
	var req dto.CreateAssetSoftwareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	software, err := h.assetSoftwareService.Create(req)
	if err != nil {
		// Log l'erreur pour le débogage
		fmt.Printf("Erreur lors de la création du logiciel: %v\n", err)
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.CreatedResponse(c, software, "Logiciel installé créé avec succès")
}

// GetByID récupère un logiciel installé par son ID
// @Summary Récupérer un logiciel installé par ID
// @Description Récupère un logiciel installé par son identifiant
// @Tags assets
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID du logiciel installé"
// @Success 200 {object} dto.AssetSoftwareDTO
// @Failure 404 {object} utils.Response
// @Router /assets/software/{id} [get]
func (h *AssetSoftwareHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	software, err := h.assetSoftwareService.GetByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Logiciel installé introuvable")
		return
	}

	utils.SuccessResponse(c, software, "Logiciel installé récupéré avec succès")
}

// GetAll récupère tous les logiciels installés
// @Summary Récupérer tous les logiciels installés
// @Description Récupère la liste de tous les logiciels installés
// @Tags assets
// @Security BearerAuth
// @Produce json
// @Success 200 {array} dto.AssetSoftwareDTO
// @Failure 500 {object} utils.Response
// @Router /assets/software [get]
func (h *AssetSoftwareHandler) GetAll(c *gin.Context) {
	softwareList, err := h.assetSoftwareService.GetAll()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des logiciels")
		return
	}

	utils.SuccessResponse(c, softwareList, "Logiciels installés récupérés avec succès")
}

// GetByAssetID récupère tous les logiciels installés sur un actif
// @Summary Récupérer les logiciels d'un actif
// @Description Récupère tous les logiciels installés sur un actif spécifique
// @Tags assets
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID de l'actif"
// @Success 200 {array} dto.AssetSoftwareDTO
// @Failure 400 {object} utils.Response
// @Router /assets/{id}/software [get]
func (h *AssetSoftwareHandler) GetByAssetID(c *gin.Context) {
	assetIDParam := c.Param("id")
	assetID, err := strconv.ParseUint(assetIDParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID d'actif invalide")
		return
	}

	softwareList, err := h.assetSoftwareService.GetByAssetID(uint(assetID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, softwareList, "Logiciels installés récupérés avec succès")
}

// GetBySoftwareName récupère tous les actifs ayant un logiciel spécifique
// @Summary Récupérer les actifs par logiciel
// @Description Récupère tous les actifs ayant un logiciel spécifique installé
// @Tags assets
// @Security BearerAuth
// @Produce json
// @Param softwareName path string true "Nom du logiciel"
// @Success 200 {array} dto.AssetSoftwareDTO
// @Failure 400 {object} utils.Response
// @Router /assets/software/by-name/{softwareName} [get]
func (h *AssetSoftwareHandler) GetBySoftwareName(c *gin.Context) {
	softwareName := c.Param("softwareName")
	if softwareName == "" {
		utils.BadRequestResponse(c, "Nom du logiciel manquant")
		return
	}

	softwareList, err := h.assetSoftwareService.GetBySoftwareName(softwareName)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, softwareList, "Actifs récupérés avec succès")
}

// GetBySoftwareNameAndVersion récupère tous les actifs ayant un logiciel avec une version spécifique
// @Summary Récupérer les actifs par logiciel et version
// @Description Récupère tous les actifs ayant un logiciel avec une version spécifique installé
// @Tags assets
// @Security BearerAuth
// @Produce json
// @Param softwareName path string true "Nom du logiciel"
// @Param version path string true "Version du logiciel"
// @Success 200 {array} dto.AssetSoftwareDTO
// @Failure 400 {object} utils.Response
// @Router /assets/software/by-name/{softwareName}/version/{version} [get]
func (h *AssetSoftwareHandler) GetBySoftwareNameAndVersion(c *gin.Context) {
	softwareName := c.Param("softwareName")
	version := c.Param("version")
	if softwareName == "" || version == "" {
		utils.BadRequestResponse(c, "Nom du logiciel ou version manquant")
		return
	}

	softwareList, err := h.assetSoftwareService.GetBySoftwareNameAndVersion(softwareName, version)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, softwareList, "Actifs récupérés avec succès")
}

// Update met à jour un logiciel installé
// @Summary Mettre à jour un logiciel installé
// @Description Met à jour un logiciel installé sur un actif
// @Tags assets
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID du logiciel installé"
// @Param request body dto.UpdateAssetSoftwareRequest true "Données de mise à jour"
// @Success 200 {object} dto.AssetSoftwareDTO
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /assets/software/{id} [put]
func (h *AssetSoftwareHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	var req dto.UpdateAssetSoftwareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	software, err := h.assetSoftwareService.Update(uint(id), req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, software, "Logiciel installé mis à jour avec succès")
}

// Delete supprime un logiciel installé
// @Summary Supprimer un logiciel installé
// @Description Supprime un logiciel installé d'un actif
// @Tags assets
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID du logiciel installé"
// @Success 200 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /assets/software/{id} [delete]
func (h *AssetSoftwareHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	err = h.assetSoftwareService.Delete(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Logiciel installé introuvable")
		return
	}

	utils.SuccessResponse(c, nil, "Logiciel installé supprimé avec succès")
}

// GetStatistics récupère des statistiques sur les logiciels installés
// @Summary Statistiques des logiciels installés
// @Description Récupère des statistiques sur les logiciels installés (nombre d'actifs par logiciel, version, catégorie)
// @Tags assets
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dto.AssetSoftwareStatisticsDTO
// @Failure 500 {object} utils.Response
// @Router /assets/software/statistics [get]
func (h *AssetSoftwareHandler) GetStatistics(c *gin.Context) {
	stats, err := h.assetSoftwareService.GetStatistics()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des statistiques")
		return
	}

	utils.SuccessResponse(c, stats, "Statistiques récupérées avec succès")
}
