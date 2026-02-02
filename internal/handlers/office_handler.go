package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/services"
	"github.com/mcicare/itsm-backend/internal/utils"
)

// OfficeHandler gère les handlers des sièges
type OfficeHandler struct {
	officeService services.OfficeService
}

// NewOfficeHandler crée une nouvelle instance de OfficeHandler
func NewOfficeHandler(officeService services.OfficeService) *OfficeHandler {
	return &OfficeHandler{
		officeService: officeService,
	}
}

// Create crée un nouveau siège
// @Summary Créer un siège
// @Description Crée un nouveau siège/bureau
// @Tags offices
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateOfficeRequest true "Données du siège"
// @Success 201 {object} dto.OfficeDTO
// @Failure 400 {object} utils.Response
// @Router /offices [post]
func (h *OfficeHandler) Create(c *gin.Context) {
	var req dto.CreateOfficeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	// Vérifier si l'utilisateur peut créer un siège dans n'importe quelle filiale
	canSelectAnyFiliale := utils.RequirePermission(c, "offices.create_any_filiale")

	// Si l'utilisateur n'a pas la permission de gérer les filiales,
	// il ne peut créer un siège que dans sa propre filiale
	if !canSelectAnyFiliale {
		// Récupérer le scope qui contient déjà la filiale de l'utilisateur
		scope := utils.GetScopeFromContext(c)
		if scope != nil && scope.FilialeID != nil {
			// Si une filiale est spécifiée et qu'elle est différente de celle du créateur, refuser
			if req.FilialeID != nil && *req.FilialeID != *scope.FilialeID {
				utils.ForbiddenResponse(c, "Vous ne pouvez créer un siège que dans votre propre filiale")
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

	office, err := h.officeService.Create(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.CreatedResponse(c, office, "Siège créé avec succès")
}

// GetAll récupère tous les sièges
// @Summary Récupérer tous les sièges
// @Description Récupère la liste des sièges selon le scope : offices.view_all = toutes les filiales ; offices.view_filiale ou offices.view = sa filiale uniquement.
// @Tags offices
// @Security BearerAuth
// @Produce json
// @Param active query bool false "Récupérer uniquement les sièges actifs"
// @Success 200 {array} dto.OfficeDTO
// @Failure 403 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /offices [get]
func (h *OfficeHandler) GetAll(c *gin.Context) {
	activeOnly := c.DefaultQuery("active", "false") == "true"
	scope := utils.GetScopeFromContext(c)
	if scope == nil {
		utils.InternalServerErrorResponse(c, "Contexte utilisateur introuvable")
		return
	}

	// Accès : au moins une permission de vue requise
	hasViewAll := scope.HasPermission("offices.view_all")
	hasViewFiliale := scope.HasPermission("offices.view_filiale") || scope.HasPermission("offices.view")
	if !hasViewAll && !hasViewFiliale {
		utils.ErrorResponse(c, http.StatusForbidden, "Vous n'avez pas la permission de voir les sièges", nil)
		return
	}

	// Vue globale : view_all ou rétrocompat (filiales.manage / create_any_filiale / update_any_filiale)
	canViewAll := hasViewAll ||
		scope.HasPermission("filiales.manage") ||
		scope.HasPermission("offices.create_any_filiale") ||
		scope.HasPermission("offices.update_any_filiale")

	var offices []dto.OfficeDTO
	var err error

	if canViewAll {
		offices, err = h.officeService.GetAll(activeOnly)
	} else {
		if scope.FilialeID != nil {
			offices, err = h.officeService.GetByFilialeID(*scope.FilialeID)
		} else {
			offices = []dto.OfficeDTO{}
		}
	}

	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des sièges")
		return
	}

	utils.SuccessResponse(c, offices, "Sièges récupérés avec succès")
}

// GetByID récupère un siège par son ID
// @Summary Récupérer un siège par ID
// @Description Récupère un siège par son identifiant (scope : view_all = tout, view_filiale/view = même filiale uniquement)
// @Tags offices
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID du siège"
// @Success 200 {object} dto.OfficeDTO
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /offices/{id} [get]
func (h *OfficeHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	office, err := h.officeService.GetByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Siège introuvable")
		return
	}

	scope := utils.GetScopeFromContext(c)
	if scope != nil && !scope.HasPermission("offices.view_all") &&
		!scope.HasPermission("filiales.manage") &&
		!scope.HasPermission("offices.create_any_filiale") &&
		!scope.HasPermission("offices.update_any_filiale") {
		if scope.FilialeID != nil && (office.FilialeID == nil || *office.FilialeID != *scope.FilialeID) {
			utils.ErrorResponse(c, http.StatusForbidden, "Vous n'avez pas accès à ce siège", nil)
			return
		}
	}

	utils.SuccessResponse(c, office, "Siège récupéré avec succès")
}

// Update met à jour un siège
// @Summary Mettre à jour un siège
// @Description Met à jour un siège
// @Tags offices
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID du siège"
// @Param request body dto.UpdateOfficeRequest true "Données de mise à jour"
// @Success 200 {object} dto.OfficeDTO
// @Failure 400 {object} utils.Response
// @Router /offices/{id} [put]
func (h *OfficeHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	var req dto.UpdateOfficeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	// Vérifier si l'utilisateur peut modifier un siège dans n'importe quelle filiale
	canSelectAnyFiliale := utils.RequirePermission(c, "offices.update_any_filiale")

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

	office, err := h.officeService.Update(uint(id), req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, office, "Siège mis à jour avec succès")
}

// Delete supprime un siège
// @Summary Supprimer un siège
// @Description Supprime un siège
// @Tags offices
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID du siège"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /offices/{id} [delete]
func (h *OfficeHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	err = h.officeService.Delete(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, nil, "Siège supprimé avec succès")
}

// GetByCountry récupère les sièges d'un pays
// @Summary Récupérer les sièges d'un pays
// @Description Récupère tous les sièges d'un pays spécifique
// @Tags offices
// @Security BearerAuth
// @Produce json
// @Param country path string true "Nom du pays"
// @Success 200 {array} dto.OfficeDTO
// @Failure 500 {object} utils.Response
// @Router /offices/country/{country} [get]
func (h *OfficeHandler) GetByCountry(c *gin.Context) {
	country := c.Param("country")

	offices, err := h.officeService.GetByCountry(country)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des sièges")
		return
	}

	utils.SuccessResponse(c, offices, "Sièges récupérés avec succès")
}

// GetByCity récupère les sièges d'une ville
// @Summary Récupérer les sièges d'une ville
// @Description Récupère tous les sièges d'une ville spécifique
// @Tags offices
// @Security BearerAuth
// @Produce json
// @Param city path string true "Nom de la ville"
// @Success 200 {array} dto.OfficeDTO
// @Failure 500 {object} utils.Response
// @Router /offices/city/{city} [get]
func (h *OfficeHandler) GetByCity(c *gin.Context) {
	city := c.Param("city")

	offices, err := h.officeService.GetByCity(city)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des sièges")
		return
	}

	utils.SuccessResponse(c, offices, "Sièges récupérés avec succès")
}
