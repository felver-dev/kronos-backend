package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/services"
	"github.com/mcicare/itsm-backend/internal/utils"
)

// ReportHandler gère les handlers des rapports
type ReportHandler struct {
	reportService services.ReportService
}

// NewReportHandler crée une nouvelle instance de ReportHandler
func NewReportHandler(reportService services.ReportService) *ReportHandler {
	return &ReportHandler{
		reportService: reportService,
	}
}

// GetDashboard récupère le tableau de bord
func (h *ReportHandler) GetDashboard(c *gin.Context) {
	period := c.DefaultQuery("period", "month")

	dashboard, err := h.reportService.GetDashboard(period)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération du tableau de bord")
		return
	}

	utils.SuccessResponse(c, dashboard, "Tableau de bord récupéré avec succès")
}

// GetTicketCountReport récupère le rapport de nombre de tickets
func (h *ReportHandler) GetTicketCountReport(c *gin.Context) {
	period := c.DefaultQuery("period", "month")

	report, err := h.reportService.GetTicketCountReport(period)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la génération du rapport")
		return
	}

	utils.SuccessResponse(c, report, "Rapport récupéré avec succès")
}

// GetTicketTypeDistribution récupère la distribution des types de tickets
func (h *ReportHandler) GetTicketTypeDistribution(c *gin.Context) {
	distribution, err := h.reportService.GetTicketTypeDistribution()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération de la distribution")
		return
	}

	utils.SuccessResponse(c, distribution, "Distribution récupérée avec succès")
}

// GetAverageResolutionTime récupère le temps moyen de résolution
func (h *ReportHandler) GetAverageResolutionTime(c *gin.Context) {
	avgTime, err := h.reportService.GetAverageResolutionTime()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors du calcul du temps moyen")
		return
	}

	utils.SuccessResponse(c, avgTime, "Temps moyen récupéré avec succès")
}

// GenerateCustomReport génère un rapport personnalisé
func (h *ReportHandler) GenerateCustomReport(c *gin.Context) {
	var req dto.CustomReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "Données invalides", err.Error())
		return
	}

	report, err := h.reportService.GenerateCustomReport(req)
	if err != nil {
		utils.ErrorResponse(c, 400, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, report, "Rapport personnalisé généré avec succès")
}

