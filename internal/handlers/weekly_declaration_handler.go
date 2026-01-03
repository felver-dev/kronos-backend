package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/services"
	"github.com/mcicare/itsm-backend/internal/utils"
)

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

