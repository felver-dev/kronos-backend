package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/services"
	"github.com/mcicare/itsm-backend/internal/utils"
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
