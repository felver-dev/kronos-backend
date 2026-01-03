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
// @Summary Récupérer une déclaration journalière par ID
// @Description Récupère une déclaration journalière par son identifiant
// @Tags daily-declarations
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID de la déclaration"
// @Success 200 {object} dto.DailyDeclarationDTO
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /daily-declarations/{id} [get]
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
// @Summary Récupérer une déclaration journalière par utilisateur et date
// @Description Récupère une déclaration journalière pour un utilisateur et une date donnée
// @Tags daily-declarations
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param user_id path int true "ID de l'utilisateur"
// @Param date query string true "Date (format: YYYY-MM-DD)"
// @Success 200 {object} dto.DailyDeclarationDTO
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /daily-declarations/users/{user_id}/by-date [get]
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
// @Summary Récupérer les déclarations journalières d'un utilisateur
// @Description Récupère toutes les déclarations journalières d'un utilisateur
// @Tags daily-declarations
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param user_id path int true "ID de l'utilisateur"
// @Success 200 {array} dto.DailyDeclarationDTO
// @Failure 400 {object} utils.Response
// @Router /daily-declarations/users/{user_id} [get]
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
// @Summary Valider une déclaration journalière
// @Description Valide une déclaration journalière
// @Tags daily-declarations
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID de la déclaration"
// @Success 200 {object} dto.DailyDeclarationDTO
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /daily-declarations/{id}/validate [post]
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
