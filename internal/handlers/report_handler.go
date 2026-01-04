package handlers

import (
	"strconv"

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
// @Summary Récupérer le tableau de bord
// @Description Récupère les statistiques du tableau de bord
// @Tags reports
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param period query string false "Période (défaut: month)"
// @Success 200 {object} dto.DashboardDTO
// @Failure 500 {object} utils.Response
// @Router /reports/dashboard [get]
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
// @Summary Récupérer le rapport de nombre de tickets
// @Description Récupère le rapport sur le nombre de tickets
// @Tags reports
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param period query string false "Période (défaut: month)"
// @Success 200 {object} dto.TicketCountReportDTO
// @Failure 500 {object} utils.Response
// @Router /reports/tickets/count [get]
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
// @Summary Récupérer la distribution des types de tickets
// @Description Récupère la distribution des types de tickets
// @Tags reports
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} dto.TicketTypeDistributionDTO
// @Failure 500 {object} utils.Response
// @Router /reports/tickets/distribution [get]
func (h *ReportHandler) GetTicketTypeDistribution(c *gin.Context) {
	distribution, err := h.reportService.GetTicketTypeDistribution()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération de la distribution")
		return
	}

	utils.SuccessResponse(c, distribution, "Distribution récupérée avec succès")
}

// GetAverageResolutionTime récupère le temps moyen de résolution
// @Summary Récupérer le temps moyen de résolution
// @Description Récupère le temps moyen de résolution des tickets
// @Tags reports
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} dto.AverageResolutionTimeDTO
// @Failure 500 {object} utils.Response
// @Router /reports/tickets/average-resolution-time [get]
func (h *ReportHandler) GetAverageResolutionTime(c *gin.Context) {
	avgTime, err := h.reportService.GetAverageResolutionTime()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors du calcul du temps moyen")
		return
	}

	utils.SuccessResponse(c, avgTime, "Temps moyen récupéré avec succès")
}

// GenerateCustomReport génère un rapport personnalisé
// @Summary Générer un rapport personnalisé
// @Description Génère un rapport personnalisé selon les critères spécifiés
// @Tags reports
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CustomReportRequest true "Critères du rapport"
// @Success 200 {object} any
// @Failure 400 {object} utils.Response
// @Router /reports/custom [post]
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

// GetWorkloadByAgent récupère la charge de travail par agent
// @Summary Récupérer la charge de travail par agent
// @Description Récupère la charge de travail (nombre de tickets) par agent
// @Tags reports
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {array} dto.WorkloadByAgentDTO
// @Failure 500 {object} utils.Response
// @Router /reports/workload/by-agent [get]
func (h *ReportHandler) GetWorkloadByAgent(c *gin.Context) {
	workload, err := h.reportService.GetWorkloadByAgent()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération de la charge de travail")
		return
	}

	utils.SuccessResponse(c, workload, "Charge de travail récupérée avec succès")
}

// GetSLAComplianceReport récupère le rapport de conformité SLA
// @Summary Récupérer le rapport de conformité SLA
// @Description Récupère le rapport de conformité des SLA
// @Tags reports
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param period query string false "Période (défaut: month)"
// @Success 200 {object} dto.SLAComplianceReportDTO
// @Failure 500 {object} utils.Response
// @Router /reports/sla/compliance [get]
func (h *ReportHandler) GetSLAComplianceReport(c *gin.Context) {
	period := c.DefaultQuery("period", "month")

	report, err := h.reportService.GetSLAComplianceReport(period)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la génération du rapport")
		return
	}

	utils.SuccessResponse(c, report, "Rapport de conformité SLA récupéré avec succès")
}

// GetDelayedTicketsReport récupère le rapport des tickets en retard
// @Summary Récupérer le rapport des tickets en retard
// @Description Récupère le rapport des tickets qui ont dépassé leur délai
// @Tags reports
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param period query string false "Période (défaut: month)"
// @Success 200 {array} dto.DelayedTicketDTO
// @Failure 500 {object} utils.Response
// @Router /reports/tickets/delayed [get]
func (h *ReportHandler) GetDelayedTicketsReport(c *gin.Context) {
	period := c.DefaultQuery("period", "month")

	report, err := h.reportService.GetDelayedTicketsReport(period)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la génération du rapport")
		return
	}

	utils.SuccessResponse(c, report, "Rapport des tickets en retard récupéré avec succès")
}

// GetIndividualPerformanceReport récupère le rapport de performance individuel
// @Summary Récupérer le rapport de performance individuel
// @Description Récupère le rapport de performance d'un utilisateur spécifique
// @Tags reports
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param userId path int true "ID de l'utilisateur"
// @Param period query string false "Période (défaut: month)"
// @Success 200 {object} dto.IndividualPerformanceReportDTO
// @Failure 500 {object} utils.Response
// @Router /reports/performance/individual/{userId} [get]
func (h *ReportHandler) GetIndividualPerformanceReport(c *gin.Context) {
	userIDParam := c.Param("userId")
	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID utilisateur invalide")
		return
	}

	period := c.DefaultQuery("period", "month")

	report, err := h.reportService.GetIndividualPerformanceReport(uint(userID), period)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la génération du rapport")
		return
	}

	utils.SuccessResponse(c, report, "Rapport de performance récupéré avec succès")
}

// ExportReport exporte un rapport dans un format spécifique
// @Summary Exporter un rapport
// @Description Exporte un rapport dans un format spécifique (PDF, Excel, CSV)
// @Tags reports
// @Security BearerAuth
// @Accept json
// @Produce application/pdf,application/vnd.openxmlformats-officedocument.spreadsheetml.sheet,text/csv
// @Param format path string true "Format d'export (pdf, excel, csv)"
// @Param reportType query string true "Type de rapport (dashboard, tickets, sla, performance)"
// @Param period query string false "Période (défaut: month)"
// @Success 200 {file} file
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /reports/export/{format} [get]
func (h *ReportHandler) ExportReport(c *gin.Context) {
	format := c.Param("format")
	reportType := c.Query("reportType")
	period := c.DefaultQuery("period", "month")

	if reportType == "" {
		utils.BadRequestResponse(c, "Type de rapport requis")
		return
	}

	file, err := h.reportService.ExportReport(reportType, format, period)
	if err != nil {
		utils.ErrorResponse(c, 400, err.Error(), nil)
		return
	}

	// Pour l'instant, on retourne les données JSON
	// TODO: Implémenter la génération de fichiers PDF/Excel/CSV
	utils.SuccessResponse(c, file, "Rapport exporté avec succès")
}
