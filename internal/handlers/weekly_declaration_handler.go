package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/services"
	"github.com/mcicare/itsm-backend/internal/utils"
)

type weeklyDeclarationDTO = dto.WeeklyDeclarationDTO

// WeeklyDeclarationHandler gère les handlers des déclarations hebdomadaires
type WeeklyDeclarationHandler struct {
	weeklyDeclarationService services.WeeklyDeclarationService
}

// NewWeeklyDeclarationHandler crée une nouvelle instance de WeeklyDeclarationHandler
func NewWeeklyDeclarationHandler(weeklyDeclarationService services.WeeklyDeclarationService) *WeeklyDeclarationHandler {
	return &WeeklyDeclarationHandler{
		weeklyDeclarationService: weeklyDeclarationService,
	}
}

// GetByID récupère une déclaration par son ID
// @Summary Récupérer une déclaration hebdomadaire par ID
// @Description Récupère une déclaration hebdomadaire par son identifiant
// @Tags weekly-declarations
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID de la déclaration"
// @Success 200 {object} weeklyDeclarationDTO
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /weekly-declarations/{id} [get]
func (h *WeeklyDeclarationHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	declaration, err := h.weeklyDeclarationService.GetByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Déclaration introuvable")
		return
	}

	utils.SuccessResponse(c, declaration, "Déclaration récupérée avec succès")
}

// GetByUserIDAndWeek récupère une déclaration par utilisateur et semaine
// @Summary Récupérer une déclaration hebdomadaire par utilisateur et semaine
// @Description Récupère une déclaration hebdomadaire pour un utilisateur et une semaine donnée
// @Tags weekly-declarations
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param user_id path int true "ID de l'utilisateur"
// @Param week query string true "Semaine (format: YYYY-WW)"
// @Success 200 {object} weeklyDeclarationDTO
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /weekly-declarations/users/{user_id}/by-week [get]
func (h *WeeklyDeclarationHandler) GetByUserIDAndWeek(c *gin.Context) {
	userIDParam := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID utilisateur invalide")
		return
	}

	week := c.Query("week")
	if week == "" {
		utils.BadRequestResponse(c, "Paramètre week manquant")
		return
	}

	declaration, err := h.weeklyDeclarationService.GetByUserIDAndWeek(uint(userID), week)
	if err != nil {
		utils.NotFoundResponse(c, "Déclaration introuvable")
		return
	}

	utils.SuccessResponse(c, declaration, "Déclaration récupérée avec succès")
}

// GetByUserID récupère les déclarations d'un utilisateur
// @Summary Récupérer les déclarations hebdomadaires d'un utilisateur
// @Description Récupère toutes les déclarations hebdomadaires d'un utilisateur
// @Tags weekly-declarations
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param user_id path int true "ID de l'utilisateur"
// @Success 200 {array} weeklyDeclarationDTO
// @Failure 400 {object} utils.Response
// @Router /weekly-declarations/users/{user_id} [get]
func (h *WeeklyDeclarationHandler) GetByUserID(c *gin.Context) {
	userIDParam := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID utilisateur invalide")
		return
	}

	declarations, err := h.weeklyDeclarationService.GetByUserID(uint(userID))
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des déclarations")
		return
	}

	utils.SuccessResponse(c, declarations, "Déclarations récupérées avec succès")
}

// Validate valide une déclaration
// @Summary Valider une déclaration hebdomadaire
// @Description Valide une déclaration hebdomadaire
// @Tags weekly-declarations
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID de la déclaration"
// @Success 200 {object} weeklyDeclarationDTO
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /weekly-declarations/{id}/validate [post]
func (h *WeeklyDeclarationHandler) Validate(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	validatedByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	declaration, err := h.weeklyDeclarationService.Validate(uint(id), validatedByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, declaration, "Déclaration validée avec succès")
}

