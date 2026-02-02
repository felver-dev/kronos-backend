package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/services"
	"github.com/mcicare/itsm-backend/internal/utils"
)

// Aliases pour la doc Swagger (évite "cannot find type definition: dto.XXX")
type (
	performanceDTO     = dto.PerformanceDTO
	efficiencyDTO      = dto.EfficiencyDTO
	productivityDTO    = dto.ProductivityDTO
	performanceRankingDTO = dto.PerformanceRankingDTO
)

// PerformanceHandler gère les handlers des performances
type PerformanceHandler struct {
	performanceService services.PerformanceService
}

// NewPerformanceHandler crée une nouvelle instance de PerformanceHandler
func NewPerformanceHandler(performanceService services.PerformanceService) *PerformanceHandler {
	return &PerformanceHandler{
		performanceService: performanceService,
	}
}

// GetPerformanceByUserID récupère les métriques de performance d'un utilisateur
// @Summary Récupérer les métriques de performance d'un utilisateur
// @Description Récupère les métriques de performance complètes d'un utilisateur
// @Tags performance
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param user_id path int true "ID de l'utilisateur"
// @Success 200 {object} performanceDTO
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /performance/users/{user_id} [get]
func (h *PerformanceHandler) GetPerformanceByUserID(c *gin.Context) {
	userIDParam := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID utilisateur invalide")
		return
	}

	performance, err := h.performanceService.GetPerformanceByUserID(uint(userID))
	if err != nil {
		utils.NotFoundResponse(c, "Utilisateur introuvable")
		return
	}

	utils.SuccessResponse(c, performance, "Métriques de performance récupérées avec succès")
}

// GetEfficiencyByUserID récupère les métriques d'efficacité d'un utilisateur
// @Summary Récupérer les métriques d'efficacité d'un utilisateur
// @Description Récupère les métriques d'efficacité d'un utilisateur
// @Tags performance
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param user_id path int true "ID de l'utilisateur"
// @Success 200 {object} efficiencyDTO
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /performance/users/{user_id}/efficiency [get]
func (h *PerformanceHandler) GetEfficiencyByUserID(c *gin.Context) {
	userIDParam := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID utilisateur invalide")
		return
	}

	efficiency, err := h.performanceService.GetEfficiencyByUserID(uint(userID))
	if err != nil {
		utils.NotFoundResponse(c, "Utilisateur introuvable")
		return
	}

	utils.SuccessResponse(c, efficiency, "Métriques d'efficacité récupérées avec succès")
}

// GetProductivityByUserID récupère les métriques de productivité d'un utilisateur
// @Summary Récupérer les métriques de productivité d'un utilisateur
// @Description Récupère les métriques de productivité d'un utilisateur
// @Tags performance
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param user_id path int true "ID de l'utilisateur"
// @Success 200 {object} productivityDTO
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /performance/users/{user_id}/productivity [get]
func (h *PerformanceHandler) GetProductivityByUserID(c *gin.Context) {
	userIDParam := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID utilisateur invalide")
		return
	}

	productivity, err := h.performanceService.GetProductivityByUserID(uint(userID))
	if err != nil {
		utils.NotFoundResponse(c, "Utilisateur introuvable")
		return
	}

	utils.SuccessResponse(c, productivity, "Métriques de productivité récupérées avec succès")
}

// GetPerformanceRanking récupère le classement des performances
// @Summary Récupérer le classement des performances
// @Description Récupère le classement des performances des utilisateurs
// @Tags performance
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param limit query int false "Nombre de résultats (défaut: 10, max: 100)"
// @Success 200 {array} performanceRankingDTO
// @Failure 500 {object} utils.Response
// @Router /performance/ranking [get]
func (h *PerformanceHandler) GetPerformanceRanking(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 || limit > 100 {
		limit = 10
	}

	ranking, err := h.performanceService.GetPerformanceRanking(limit)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération du classement")
		return
	}

	utils.SuccessResponse(c, ranking, "Classement récupéré avec succès")
}
