package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/services"
	"github.com/mcicare/itsm-backend/internal/utils"
)

// AssetCategoryHandler gère les handlers des catégories d'actifs IT
type AssetCategoryHandler struct {
	assetCategoryService services.AssetCategoryService
}

// NewAssetCategoryHandler crée une nouvelle instance de AssetCategoryHandler
func NewAssetCategoryHandler(assetCategoryService services.AssetCategoryService) *AssetCategoryHandler {
	return &AssetCategoryHandler{
		assetCategoryService: assetCategoryService,
	}
}

// GetAll récupère toutes les catégories d'actifs avec pagination
// @Summary Récupérer toutes les catégories d'actifs
// @Description Récupère la liste de toutes les catégories d'actifs IT avec pagination
// @Tags assets
// @Security BearerAuth
// @Produce json
// @Param page query int false "Numéro de page (défaut: 1)"
// @Param limit query int false "Nombre d'éléments par page (défaut: 25, max: 100)"
// @Success 200 {object} dto.AssetCategoryListResponse
// @Failure 500 {object} utils.Response
// @Router /assets/categories [get]
func (h *AssetCategoryHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "25"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 25
	}

	response, err := h.assetCategoryService.GetAllPaginated(page, limit)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des catégories")
		return
	}

	utils.SuccessResponse(c, response, "Catégories récupérées avec succès")
}

// GetByID récupère une catégorie par son ID
// @Summary Récupérer une catégorie par ID
// @Description Récupère une catégorie d'actif IT par son identifiant
// @Tags assets
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID de la catégorie"
// @Success 200 {object} dto.AssetCategoryDTO
// @Failure 404 {object} utils.Response
// @Router /assets/categories/{id} [get]
func (h *AssetCategoryHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	category, err := h.assetCategoryService.GetByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Catégorie introuvable")
		return
	}

	utils.SuccessResponse(c, category, "Catégorie récupérée avec succès")
}

// Create crée une nouvelle catégorie d'actif
// @Summary Créer une catégorie d'actif
// @Description Crée une nouvelle catégorie d'actif IT
// @Tags assets
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateAssetCategoryRequest true "Données de la catégorie"
// @Success 201 {object} dto.AssetCategoryDTO
// @Failure 400 {object} utils.Response
// @Router /assets/categories [post]
func (h *AssetCategoryHandler) Create(c *gin.Context) {
	var req dto.CreateAssetCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	createdByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	category, err := h.assetCategoryService.Create(req, createdByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.CreatedResponse(c, category, "Catégorie créée avec succès")
}

// Update met à jour une catégorie d'actif
// @Summary Mettre à jour une catégorie d'actif
// @Description Met à jour les informations d'une catégorie d'actif IT
// @Tags assets
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID de la catégorie"
// @Param request body dto.UpdateAssetCategoryRequest true "Données à mettre à jour"
// @Success 200 {object} dto.AssetCategoryDTO
// @Failure 400 {object} utils.Response
// @Router /assets/categories/{id} [put]
func (h *AssetCategoryHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	var req dto.UpdateAssetCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	updatedByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	category, err := h.assetCategoryService.Update(uint(id), req, updatedByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, category, "Catégorie mise à jour avec succès")
}

// Delete supprime une catégorie d'actif
// @Summary Supprimer une catégorie d'actif
// @Description Supprime une catégorie d'actif IT du système. Si la catégorie a des sous-catégories, le nom de confirmation doit être fourni pour supprimer en cascade.
// @Tags assets
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID de la catégorie"
// @Param request body dto.DeleteAssetCategoryRequest false "Requête de suppression (confirm_name requis si la catégorie a des enfants)"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /assets/categories/{id} [delete]
func (h *AssetCategoryHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	// Lire le body pour obtenir le nom de confirmation (optionnel)
	var req dto.DeleteAssetCategoryRequest
	confirmName := ""
	if c.Request.ContentLength > 0 {
		if err := c.ShouldBindJSON(&req); err == nil {
			confirmName = req.ConfirmName
		}
		// Si le body n'est pas du JSON valide, on continue sans erreur (pour compatibilité)
	}

	err = h.assetCategoryService.Delete(uint(id), confirmName)
	if err != nil {
		// Vérifier le type d'erreur pour retourner le bon code HTTP
		errMsg := err.Error()
		if errMsg == "catégorie introuvable" || errMsg == "Catégorie introuvable" {
			utils.NotFoundResponse(c, errMsg)
		} else {
			// Erreur de validation (sous-catégories ou actifs associés)
			utils.ErrorResponse(c, http.StatusBadRequest, errMsg, nil)
		}
		return
	}

	utils.SuccessResponse(c, nil, "Catégorie supprimée avec succès")
}

