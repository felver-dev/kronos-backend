package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/services"
	"github.com/mcicare/itsm-backend/internal/utils"
)

// DailyDeclarationHandler gère les handlers des déclarations journalières
type DailyDeclarationHandler struct {
	dailyDeclarationService services.DailyDeclarationService
}

// NewDailyDeclarationHandler crée une nouvelle instance de DailyDeclarationHandler
func NewDailyDeclarationHandler(dailyDeclarationService services.DailyDeclarationService) *DailyDeclarationHandler {
	return &DailyDeclarationHandler{
		dailyDeclarationService: dailyDeclarationService,
	}
}

// GetByID récupère une déclaration par son ID
func (h *DailyDeclarationHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	declaration, err := h.dailyDeclarationService.GetByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Déclaration introuvable")
		return
	}

	utils.SuccessResponse(c, declaration, "Déclaration récupérée avec succès")
}

// GetByUserIDAndDate récupère une déclaration par utilisateur et date
func (h *DailyDeclarationHandler) GetByUserIDAndDate(c *gin.Context) {
	userIDParam := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID utilisateur invalide")
		return
	}

	dateParam := c.Query("date")
	if dateParam == "" {
		utils.BadRequestResponse(c, "Paramètre date manquant")
		return
	}

	date, err := time.Parse("2006-01-02", dateParam)
	if err != nil {
		utils.BadRequestResponse(c, "Format de date invalide (attendu: YYYY-MM-DD)")
		return
	}

	declaration, err := h.dailyDeclarationService.GetByUserIDAndDate(uint(userID), date)
	if err != nil {
		utils.NotFoundResponse(c, "Déclaration introuvable")
		return
	}

	utils.SuccessResponse(c, declaration, "Déclaration récupérée avec succès")
}

// GetByUserID récupère les déclarations d'un utilisateur
func (h *DailyDeclarationHandler) GetByUserID(c *gin.Context) {
	userIDParam := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID utilisateur invalide")
		return
	}

	declarations, err := h.dailyDeclarationService.GetByUserID(uint(userID))
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des déclarations")
		return
	}

	utils.SuccessResponse(c, declarations, "Déclarations récupérées avec succès")
}

// Validate valide une déclaration
func (h *DailyDeclarationHandler) Validate(c *gin.Context) {
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

	declaration, err := h.dailyDeclarationService.Validate(uint(id), validatedByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, declaration, "Déclaration validée avec succès")
}
