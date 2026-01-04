package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/services"
	"github.com/mcicare/itsm-backend/internal/utils"
)

// StatisticsHandler gère les handlers des statistiques
type StatisticsHandler struct {
	statisticsService services.StatisticsService
}

// NewStatisticsHandler crée une nouvelle instance de StatisticsHandler
func NewStatisticsHandler(statisticsService services.StatisticsService) *StatisticsHandler {
	return &StatisticsHandler{
		statisticsService: statisticsService,
	}
}

// GetOverview récupère la vue d'ensemble des statistiques
// @Summary Vue d'ensemble des statistiques
// @Description Récupère une vue d'ensemble des statistiques du système
// @Tags statistics
// @Security BearerAuth
// @Produce json
// @Param period query string false "Période (week, month, quarter, year) - défaut: month"
// @Success 200 {object} dto.StatisticsOverviewDTO
// @Failure 500 {object} utils.Response
// @Router /stats/overview [get]
func (h *StatisticsHandler) GetOverview(c *gin.Context) {
	period := c.DefaultQuery("period", "month")

	overview, err := h.statisticsService.GetOverview(period)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des statistiques")
		return
	}

	utils.SuccessResponse(c, overview, "Statistiques récupérées avec succès")
}

// GetWorkload récupère les statistiques de charge de travail
// @Summary Statistiques de charge de travail
// @Description Récupère les statistiques de charge de travail
// @Tags statistics
// @Security BearerAuth
// @Produce json
// @Param period query string false "Période (week, month, quarter, year) - défaut: month"
// @Param userId query int false "ID de l'utilisateur (optionnel)"
// @Success 200 {object} dto.WorkloadStatisticsDTO
// @Failure 500 {object} utils.Response
// @Router /stats/workload [get]
func (h *StatisticsHandler) GetWorkload(c *gin.Context) {
	period := c.DefaultQuery("period", "month")
	userIDStr := c.Query("userId")

	var userID *uint
	if userIDStr != "" {
		id, err := strconv.ParseUint(userIDStr, 10, 32)
		if err == nil {
			uid := uint(id)
			userID = &uid
		}
	}

	workload, err := h.statisticsService.GetWorkload(period, userID)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération de la charge de travail")
		return
	}

	utils.SuccessResponse(c, workload, "Statistiques de charge de travail récupérées avec succès")
}

// GetPerformance récupère les statistiques de performance
// @Summary Statistiques de performance
// @Description Récupère les statistiques de performance globales
// @Tags statistics
// @Security BearerAuth
// @Produce json
// @Param period query string false "Période (week, month, quarter, year) - défaut: month"
// @Success 200 {object} dto.PerformanceStatisticsDTO
// @Failure 500 {object} utils.Response
// @Router /stats/performance [get]
func (h *StatisticsHandler) GetPerformance(c *gin.Context) {
	period := c.DefaultQuery("period", "month")

	performance, err := h.statisticsService.GetPerformance(period)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des statistiques de performance")
		return
	}

	utils.SuccessResponse(c, performance, "Statistiques de performance récupérées avec succès")
}

// GetTrends récupère les tendances
// @Summary Tendances
// @Description Récupère les tendances pour une métrique donnée
// @Tags statistics
// @Security BearerAuth
// @Produce json
// @Param metric query string true "Métrique (tickets, resolution_time, sla_compliance, etc.)"
// @Param period query string false "Période (1month, 3months, 6months, year) - défaut: 3months"
// @Success 200 {object} dto.TrendsStatisticsDTO
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /stats/trends [get]
func (h *StatisticsHandler) GetTrends(c *gin.Context) {
	metric := c.Query("metric")
	if metric == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Paramètre 'metric' manquant", nil)
		return
	}

	period := c.DefaultQuery("period", "3months")

	trends, err := h.statisticsService.GetTrends(metric, period)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des tendances")
		return
	}

	utils.SuccessResponse(c, trends, "Tendances récupérées avec succès")
}

// GetKPI récupère les indicateurs de succès (KPI)
// @Summary Indicateurs de succès (KPI)
// @Description Récupère les indicateurs de succès (KPI) du système
// @Tags statistics
// @Security BearerAuth
// @Produce json
// @Param period query string false "Période (week, month, quarter, year) - défaut: month"
// @Success 200 {object} dto.KPIStatisticsDTO
// @Failure 500 {object} utils.Response
// @Router /stats/kpi [get]
func (h *StatisticsHandler) GetKPI(c *gin.Context) {
	period := c.DefaultQuery("period", "month")

	kpi, err := h.statisticsService.GetKPI(period)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des KPI")
		return
	}

	utils.SuccessResponse(c, kpi, "KPI récupérés avec succès")
}

