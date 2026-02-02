package services

import (
	"errors"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/mcicare/itsm-backend/database"
	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/models"
	"github.com/mcicare/itsm-backend/internal/repositories"
	"github.com/mcicare/itsm-backend/internal/scope"
)

// ReportService interface pour les opérations sur les rapports
type ReportService interface {
	GetDashboard(scope interface{}, period string) (*dto.DashboardDTO, error) // scope peut être *scope.QueryScope ou nil
	GetTicketCountReport(scope interface{}, period string) (*dto.TicketCountReportDTO, error)
	GetTicketTypeDistribution(scope interface{}) (*dto.TicketTypeDistributionDTO, error)
	GetAverageResolutionTime(scope interface{}) (*dto.AverageResolutionTimeDTO, error)
	GetWorkloadByAgent(scope interface{}, period string) ([]dto.WorkloadByAgentDTO, error)
	GetSLAComplianceReport(scope interface{}, period string) (*dto.SLAComplianceReportDTO, error)
	GetDelayedTicketsReport(scope interface{}, period string) ([]dto.DelayedTicketDTO, error)
	GetIndividualPerformanceReport(userID uint, period string) (*dto.IndividualPerformanceReportDTO, error)
	GetAssetSummary(scope interface{}, period string) (*dto.AssetReportDTO, error)
	GetKnowledgeSummary(scope interface{}, period string) (*dto.KnowledgeReportDTO, error)
	ExportReport(reportType, format, period string) (any, error)
	GenerateCustomReport(req dto.CustomReportRequest) (any, error)
}

// reportService implémente ReportService
type reportService struct {
	ticketRepo        repositories.TicketRepository
	ticketInternalRepo repositories.TicketInternalRepository
	slaRepo           repositories.SLARepository
	userRepo          repositories.UserRepository
}

// NewReportService crée une nouvelle instance de ReportService
func NewReportService(
	ticketRepo repositories.TicketRepository,
	ticketInternalRepo repositories.TicketInternalRepository,
	slaRepo repositories.SLARepository,
	userRepo repositories.UserRepository,
) ReportService {
	return &reportService{
		ticketRepo:        ticketRepo,
		ticketInternalRepo: ticketInternalRepo,
		slaRepo:           slaRepo,
		userRepo:          userRepo,
	}
}

func normalizePeriod(period string) string {
	switch period {
	case "week", "month", "quarter", "year":
		return period
	default:
		return "month"
	}
}

func periodStart(period string, now time.Time) time.Time {
	switch normalizePeriod(period) {
	case "week":
		return now.AddDate(0, 0, -7)
	case "quarter":
		return now.AddDate(0, -3, 0)
	case "year":
		return now.AddDate(-1, 0, 0)
	default: // month
		return now.AddDate(0, -1, 0)
	}
}

// GetDashboard récupère le tableau de bord
// Le scope est utilisé pour filtrer automatiquement selon les permissions de l'utilisateur
func (s *reportService) GetDashboard(scopeParam interface{}, period string) (*dto.DashboardDTO, error) {
	now := time.Now()
	start := periodStart(period, now)

	// Tableau de bord département demandé mais l'utilisateur n'a pas de département associé
	if scopeParam != nil {
		if queryScope, ok := scopeParam.(*scope.QueryScope); ok && queryScope.DashboardScopeHint == "department" && queryScope.DepartmentID == nil {
			return &dto.DashboardDTO{
				Tickets:     dto.TicketStatsDTO{},
				SLA:         dto.SLAStatsDTO{},
				Performance: dto.PerformanceStatsDTO{},
				Alerts:      []dto.AlertDTO{},
				Period:      normalizePeriod(period),
				Users:       dto.UserStatsDTO{},
				Assets:      dto.AssetStatsDTO{ByStatus: map[string]int{}, ByCategory: map[string]int{}},
				WorkedHours: dto.WorkedHoursStatsDTO{Period: normalizePeriod(period)},
				Message:     "Aucun département n'est associé à votre compte. Veuillez contacter l'administrateur pour associer votre compte à un département (ex: Direction des Systèmes d'Informations - MCI-DSI) afin d'afficher les données du tableau de bord.",
			}, nil
		}
	}

	totalTickets, err := s.countTicketsSince(scopeParam, start)
	if err != nil {
		return nil, err
	}

	byStatus, err := s.countTicketsByFieldSince(scopeParam, "status", start)
	if err != nil {
		return nil, err
	}

	byCategory, err := s.countTicketsByFieldSince(scopeParam, "category", start)
	if err != nil {
		return nil, err
	}

	byPriority, err := s.countTicketsByFieldSince(scopeParam, "priority", start)
	if err != nil {
		return nil, err
	}

	avgTime, err := s.GetAverageResolutionTime(scopeParam)
	if err != nil {
		return nil, err
	}

	slaReport, err := s.GetSLAComplianceReport(scopeParam, period)
	if err != nil {
		return nil, err
	}

	// Fusionner "resolu" dans "cloture" pour les graphiques (répartition par statut)
	if countResolu, ok := byStatus["resolu"]; ok {
		byStatus["cloture"] = byStatus["cloture"] + countResolu
		delete(byStatus, "resolu")
	}

	openTickets := 0
	closedTickets := 0
	if count, ok := byStatus["ouvert"]; ok {
		openTickets += count
	}
	if count, ok := byStatus["en_cours"]; ok {
		openTickets += count
	}
	if count, ok := byStatus["en_attente"]; ok {
		openTickets += count
	}
	if count, ok := byStatus["cloture"]; ok {
		closedTickets += count
	}

	// Statistiques utilisateurs (du département si scope département, sinon global)
	var totalUsers int64
	var activeUsers int64
	if deptIDs, ok := getDepartmentUserIDs(scopeParam); ok && len(deptIDs) > 0 {
		totalUsers = int64(len(deptIDs))
		activeUsers = int64(len(deptIDs))
	} else {
		if err := database.DB.Model(&models.User{}).Count(&totalUsers).Error; err != nil {
			log.Printf("Erreur lors du comptage des utilisateurs: %v", err)
		}
		if err := database.DB.Model(&models.User{}).Where("is_active = ?", true).Count(&activeUsers).Error; err != nil {
			log.Printf("Erreur lors du comptage des utilisateurs actifs: %v", err)
		}
	}

	// Statistiques actifs (tous les actifs, pas seulement ceux créés dans la période)
	assetSummary, err := s.GetAssetSummary(scopeParam, "year") // Utiliser "year" pour avoir tous les actifs
	if err != nil {
		log.Printf("Erreur lors de la récupération du résumé des actifs: %v", err)
		assetSummary = &dto.AssetReportDTO{Total: 0, ByStatus: make(map[string]int), ByCategory: make(map[string]int)}
	}

	// Heures travaillées (depuis time_entries) — conservées en API, non affichées au board
	var totalMinutes int64
	teQuery := database.DB.Model(&models.TimeEntry{}).Where("date >= ?", start)
	if deptIDs, ok := getDepartmentUserIDs(scopeParam); ok && len(deptIDs) > 0 {
		teQuery = teQuery.Where("user_id IN ?", deptIDs)
	}
	_ = teQuery.Select("COALESCE(SUM(time_spent), 0)").Scan(&totalMinutes).Error

	// Statistiques tickets internes (même périmètre que le tableau de bord) — uniquement si l'utilisateur a une permission tickets_internes.view_*
	var ticketInternes *dto.TicketInternalStatsDTO
	if scopeParam != nil {
		if queryScope, ok := scopeParam.(*scope.QueryScope); ok &&
			queryScope.HasAnyPermission("tickets_internes.view_all", "tickets_internes.view_filiale", "tickets_internes.view_department", "tickets_internes.view_own") {
			if totalTI, byStatusTI, openTI, closedTI, errTI := s.ticketInternalRepo.GetStatsForDashboard(scopeParam); errTI == nil {
				ticketInternes = &dto.TicketInternalStatsDTO{
					Total:    totalTI,
					ByStatus: byStatusTI,
					Open:     openTI,
					Closed:   closedTI,
				}
			}
		}
	}

	return &dto.DashboardDTO{
		Tickets: dto.TicketStatsDTO{
			Total:                 totalTickets,
			ByCategory:            byCategory,
			ByStatus:              byStatus,
			ByPriority:            byPriority,
			AverageResolutionTime: float64(avgTime.AverageTime),
			Delayed:               slaReport.TotalViolations,
			Open:                  openTickets,
			Closed:                closedTickets,
		},
		SLA: dto.SLAStatsDTO{
			OverallCompliance: slaReport.OverallCompliance,
			ByCategory:        slaReport.ByCategory,
			ByPriority:        slaReport.ByPriority,
			TotalViolations:   slaReport.TotalViolations,
			AtRisk:            0,
		},
		Performance: dto.PerformanceStatsDTO{},
		Alerts:      []dto.AlertDTO{},
		Period:      normalizePeriod(period),
		Users: dto.UserStatsDTO{
			Total:  int(totalUsers),
			Active: int(activeUsers),
			ByRole: make(map[string]int), // Peut être rempli plus tard si nécessaire
		},
		Assets: dto.AssetStatsDTO{
			Total:      assetSummary.Total,
			ByStatus:   assetSummary.ByStatus,
			ByCategory: assetSummary.ByCategory,
		},
		WorkedHours:    dto.WorkedHoursStatsDTO{
			TotalMinutes: int(totalMinutes),
			TotalHours:   float64(totalMinutes) / 60.0,
			Period:       normalizePeriod(period),
		},
		TicketInternes: ticketInternes,
	}, nil
}

// GetTicketCountReport récupère le rapport de nombre de tickets
// Le scope est utilisé pour filtrer automatiquement selon les permissions de l'utilisateur
func (s *reportService) GetTicketCountReport(scopeParam interface{}, period string) (*dto.TicketCountReportDTO, error) {
	now := time.Now()
	start := periodStart(period, now)
	// S'assurer que la date de début est au début de la journée (00:00:00)
	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())

	total, err := s.countTicketsSince(scopeParam, start)
	if err != nil {
		return nil, err
	}

	// Tableau de bord département sans département associé : ne pas renvoyer de breakdown
	if queryScope, ok := scopeParam.(*scope.QueryScope); ok && queryScope.DashboardScopeHint == "department" && queryScope.DepartmentID == nil {
		return &dto.TicketCountReportDTO{
			Period:    normalizePeriod(period),
			Count:     total,
			Breakdown: []dto.PeriodBreakdownDTO{},
		}, nil
	}

	type row struct {
		Period     string `gorm:"column:period"`
		Count      int    `gorm:"column:count"`
		Open       int    `gorm:"column:open"`
		InProgress int    `gorm:"column:in_progress"`
		Pending    int    `gorm:"column:pending"`
		Resolved   int    `gorm:"column:resolved"`
		Closed     int    `gorm:"column:closed"`
	}

	normalizedPeriod := normalizePeriod(period)

	dateLayout := "2006-01-02"
	var rows []row

	// Construire la requête de base
	baseQuery := database.DB.Table("tickets")

	if deptIDs, ok := getDepartmentUserIDs(scopeParam); ok && len(deptIDs) > 0 {
		baseQuery = baseQuery.Where("created_by_id IN ?", deptIDs)
	} else if scopeParam != nil {
		if queryScope, ok := scopeParam.(*scope.QueryScope); ok {
			baseQuery = scope.ApplyTicketScopeToTable(baseQuery, queryScope)
		}
	}

	// Utiliser une requête plus robuste qui fonctionne avec MySQL/MariaDB
	if normalizedPeriod == "quarter" || normalizedPeriod == "year" {
		dateLayout = "2006-01"
		query := baseQuery.Select(`
			DATE_FORMAT(created_at, '%Y-%m') as period,
			COUNT(*) as count,
			SUM(CASE WHEN status = 'ouvert' OR status = 'open' THEN 1 ELSE 0 END) as open,
			SUM(CASE WHEN status = 'en_cours' OR status = 'in_progress' THEN 1 ELSE 0 END) as in_progress,
			SUM(CASE WHEN status = 'en_attente' OR status = 'pending' THEN 1 ELSE 0 END) as pending,
			SUM(CASE WHEN status = 'resolu' THEN 1 ELSE 0 END) as resolved,
			SUM(CASE WHEN status = 'cloture' OR status = 'closed' THEN 1 ELSE 0 END) as closed
		`).
			Where("created_at >= ?", start).
			Group("DATE_FORMAT(created_at, '%Y-%m')").
			Order("period ASC")

		if err := query.Scan(&rows).Error; err != nil {
			return nil, fmt.Errorf("erreur lors de la récupération du breakdown: %w", err)
		}
	} else {
		// Pour week et month : mêmes colonnes que quarter/year (alignées sur les statuts en base)
		query := baseQuery.Select(`
			CAST(created_at AS DATE) as period,
			COUNT(*) as count,
			SUM(CASE WHEN status = 'ouvert' OR status = 'open' THEN 1 ELSE 0 END) as open,
			SUM(CASE WHEN status = 'en_cours' OR status = 'in_progress' THEN 1 ELSE 0 END) as in_progress,
			SUM(CASE WHEN status = 'en_attente' OR status = 'pending' THEN 1 ELSE 0 END) as pending,
			SUM(CASE WHEN status = 'resolu' THEN 1 ELSE 0 END) as resolved,
			SUM(CASE WHEN status = 'cloture' OR status = 'closed' THEN 1 ELSE 0 END) as closed
		`).
			Where("created_at >= ?", start).
			Group("CAST(created_at AS DATE)").
			Order("period ASC")

		if err := query.Scan(&rows).Error; err != nil {
			// Si CAST ne fonctionne pas, essayer avec DATE()
			log.Printf("Erreur avec CAST, essai avec DATE(): %v", err)
			query = baseQuery.Select(`
				DATE(created_at) as period,
				COUNT(*) as count,
				SUM(CASE WHEN status = 'ouvert' OR status = 'open' THEN 1 ELSE 0 END) as open,
				SUM(CASE WHEN status = 'en_cours' OR status = 'in_progress' THEN 1 ELSE 0 END) as in_progress,
				SUM(CASE WHEN status = 'en_attente' OR status = 'pending' THEN 1 ELSE 0 END) as pending,
				SUM(CASE WHEN status = 'resolu' THEN 1 ELSE 0 END) as resolved,
				SUM(CASE WHEN status = 'cloture' OR status = 'closed' THEN 1 ELSE 0 END) as closed
			`).
				Where("created_at >= ?", start).
				Group("DATE(created_at)").
				Order("period ASC")

			if err := query.Scan(&rows).Error; err != nil {
				return nil, fmt.Errorf("erreur lors de la récupération du breakdown: %w", err)
			}
		}
	}

	// Log pour déboguer
	log.Printf("GetTicketCountReport: period=%s, start=%v, total=%d, rows found=%d", normalizedPeriod, start, total, len(rows))
	for i, r := range rows {
		log.Printf("  Row %d: Period='%s', Count=%d", i, r.Period, r.Count)
	}

	breakdown := make([]dto.PeriodBreakdownDTO, 0, len(rows))
	for _, r := range rows {
		if r.Period == "" {
			continue
		}
		// Nettoyer la chaîne de date (enlever les espaces, etc.)
		periodStr := r.Period
		// Si la date contient un espace ou autre chose, prendre seulement la partie date
		if len(periodStr) > 10 {
			periodStr = periodStr[:10]
		}

		parsed, err := time.Parse(dateLayout, periodStr)
		if err != nil {
			// Si le parsing échoue avec le format attendu, essayer d'autres formats
			// Peut-être que MySQL retourne un format différent
			if altParsed, altErr := time.Parse("2006-01-02 15:04:05", r.Period); altErr == nil {
				parsed = altParsed
			} else if altParsed, altErr := time.Parse("2006-01-02T15:04:05Z", r.Period); altErr == nil {
				parsed = altParsed
			} else {
				// Si tous les formats échouent, on skip cette ligne
				continue
			}
		}
		breakdown = append(breakdown, dto.PeriodBreakdownDTO{
			Date:       parsed,
			Count:      r.Count,
			Open:       r.Open,
			InProgress: r.InProgress,
			Pending:    r.Pending,
			Resolved:   r.Resolved,
			Closed:     r.Closed,
		})
	}

	return &dto.TicketCountReportDTO{
		Period:    normalizedPeriod,
		Count:     total,
		Breakdown: breakdown,
	}, nil
}

// GetTicketTypeDistribution récupère la distribution des types de tickets
// Le scope est utilisé pour filtrer automatiquement selon les permissions de l'utilisateur
func (s *reportService) GetTicketTypeDistribution(scopeParam interface{}) (*dto.TicketTypeDistributionDTO, error) {
	byCategory, err := s.countTicketsByFieldSince(scopeParam, "category", time.Time{})
	if err != nil {
		return nil, err
	}

	// Helper function to safely get count from map
	getCount := func(key string) int {
		if count, ok := byCategory[key]; ok {
			return count
		}
		return 0
	}

	// Calculer le total pour vérifier
	totalCounted := getCount("incident") + getCount("incidents") +
		getCount("demande") + getCount("demandes") +
		getCount("changement") + getCount("changements") +
		getCount("developpement") + getCount("developpements") +
		getCount("assistance") +
		getCount("support")

	// Compter le total réel de tickets
	var totalTickets int64
	if err := database.DB.Model(&models.Ticket{}).Count(&totalTickets).Error; err != nil {
		log.Printf("Erreur lors du comptage total des tickets: %v", err)
	} else {
		if int64(totalCounted) != totalTickets {
			log.Printf("ATTENTION: Total compté (%d) ne correspond pas au total réel (%d). Catégories non reconnues: %v",
				totalCounted, totalTickets, byCategory)
		}
	}

	return &dto.TicketTypeDistributionDTO{
		Incidents:      getCount("incident") + getCount("incidents"), // Support both singular and plural
		Demandes:       getCount("demande") + getCount("demandes"),
		Changements:    getCount("changement") + getCount("changements"),
		Developpements: getCount("developpement") + getCount("developpements"),
		Assistance:     getCount("assistance"),
		Support:        getCount("support"),
	}, nil
}

// GetAverageResolutionTime récupère le temps moyen de résolution
// Calcule la moyenne du temps de résolution pour tous les tickets clôturés (status = cloture).
// Utilise actual_time si disponible et > 0, sinon durée entre created_at et COALESCE(closed_at, updated_at).
// Inclut les tickets clôturés même sans closed_at (fallback sur updated_at).
// Le scope est utilisé pour filtrer automatiquement selon les permissions de l'utilisateur
func (s *reportService) GetAverageResolutionTime(scopeParam interface{}) (*dto.AverageResolutionTimeDTO, error) {
	// Construire la requête de base : tous les tickets clôturés (avec ou sans closed_at)
	baseQuery := database.DB.Table("tickets").Where("status = ?", "cloture")

	// Appliquer le scope si fourni
	if scopeParam != nil {
		if queryScope, ok := scopeParam.(*scope.QueryScope); ok {
			baseQuery = scope.ApplyTicketScopeToTable(baseQuery, queryScope)
		}
	}

	// Compter le nombre de tickets clôturés pour le log
	var totalClosed int64
	if err := baseQuery.Count(&totalClosed).Error; err != nil {
		log.Printf("Erreur lors du comptage des tickets clôturés: %v", err)
	}

	// Reconstruire une requête identique pour l'AVG (éviter réutilisation après Count en GORM)
	avgQuery := database.DB.Table("tickets").Where("status = ?", "cloture")
	if scopeParam != nil {
		if queryScope, ok := scopeParam.(*scope.QueryScope); ok {
			avgQuery = scope.ApplyTicketScopeToTable(avgQuery, queryScope)
		}
	}

	type avgRow struct {
		Average float64 `gorm:"column:average_time"`
	}
	var row avgRow
	// actual_time si > 0, sinon durée created_at -> COALESCE(closed_at, updated_at)
	sel := "AVG(COALESCE(NULLIF(actual_time, 0), TIMESTAMPDIFF(MINUTE, created_at, COALESCE(closed_at, updated_at)))) as average_time"
	err := avgQuery.Select(sel).Scan(&row).Error
	if err != nil {
		return nil, err
	}

	// Si aucun ticket clôturé ou moyenne invalide, retourner 0
	if totalClosed == 0 || math.IsNaN(row.Average) || row.Average <= 0 {
		return &dto.AverageResolutionTimeDTO{
			AverageTime: 0,
			Unit:        "minutes",
			Breakdown:   map[string]float64{},
		}, nil
	}

	result := int(math.Round(row.Average))
	log.Printf("[Rapports] Temps moyen de résolution: %d minutes (basé sur %d tickets)", result, totalClosed)

	return &dto.AverageResolutionTimeDTO{
		AverageTime: result,
		Unit:        "minutes",
		Breakdown:   map[string]float64{},
	}, nil
}

// workloadAcc cumule les comptes tickets + ticket_internes par agent
type workloadAcc struct {
	TicketCount     int
	ResolvedCount   int
	InProgressCount int
	PendingCount    int
	OpenCount       int
	DelayedCount    int
	TotalTime       float64
}

// GetWorkloadByAgent récupère la charge de travail par agent (tickets normaux + tickets internes)
// Le scope est utilisé pour filtrer automatiquement selon les permissions de l'utilisateur
func (s *reportService) GetWorkloadByAgent(scopeParam interface{}, period string) ([]dto.WorkloadByAgentDTO, error) {
	now := time.Now()
	start := periodStart(period, now)

	type workloadRow struct {
		AssignedToID    *uint   `gorm:"column:assigned_to_id"`
		TicketCount     int     `gorm:"column:ticket_count"`
		ResolvedCount   int     `gorm:"column:resolved_count"`
		InProgressCount int     `gorm:"column:in_progress_count"`
		PendingCount    int     `gorm:"column:pending_count"`
		OpenCount       int     `gorm:"column:open_count"`
		DelayedCount    int     `gorm:"column:delayed_count"`
		AverageTime     float64 `gorm:"column:average_time"`
		TotalTime       float64 `gorm:"column:total_time"`
	}

	departmentUserIDs, _ := getDepartmentUserIDs(scopeParam)
	inDepartment := func(uid uint) bool {
		if len(departmentUserIDs) == 0 {
			return true
		}
		for _, id := range departmentUserIDs {
			if id == uid {
				return true
			}
		}
		return false
	}

	userAcc := make(map[uint]*workloadAcc)

	// 1) Charge depuis la table tickets
	baseQuery := database.DB.Table("tickets").
		Where("assigned_to_id IS NOT NULL AND created_at >= ?", start)
	if scopeParam != nil {
		if queryScope, ok := scopeParam.(*scope.QueryScope); ok {
			baseQuery = scope.ApplyTicketScopeToTable(baseQuery, queryScope)
		}
	}
	var rows []workloadRow
	// Résolus = cloture OU resolu (workflow : validation → resolu, fermeture → cloture). Même règle pour temps moyen/total.
	query := baseQuery.Select(`
		assigned_to_id,
		COUNT(*) as ticket_count,
		SUM(CASE WHEN status IN ('cloture', 'resolu') THEN 1 ELSE 0 END) as resolved_count,
		SUM(CASE WHEN status = 'en_cours' THEN 1 ELSE 0 END) as in_progress_count,
		SUM(CASE WHEN status = 'en_attente' THEN 1 ELSE 0 END) as pending_count,
		SUM(CASE WHEN status = 'ouvert' THEN 1 ELSE 0 END) as open_count,
		SUM(CASE WHEN status NOT IN ('cloture', 'resolu') AND closed_at IS NULL AND DATEDIFF(NOW(), created_at) > 7 THEN 1 ELSE 0 END) as delayed_count,
		AVG(CASE WHEN status IN ('cloture', 'resolu') THEN COALESCE(NULLIF(actual_time, 0), TIMESTAMPDIFF(MINUTE, created_at, COALESCE(closed_at, updated_at))) ELSE NULL END) as average_time,
		SUM(CASE WHEN status IN ('cloture', 'resolu') THEN COALESCE(NULLIF(actual_time, 0), TIMESTAMPDIFF(MINUTE, created_at, COALESCE(closed_at, updated_at))) ELSE 0 END) as total_time
	`).Group("assigned_to_id")
	if err := query.Scan(&rows).Error; err != nil {
		return nil, err
	}
	for _, row := range rows {
		if row.AssignedToID == nil || !inDepartment(*row.AssignedToID) {
			continue
		}
		uid := *row.AssignedToID
		userAcc[uid] = &workloadAcc{
			TicketCount:     row.TicketCount,
			ResolvedCount:   row.ResolvedCount,
			InProgressCount: row.InProgressCount,
			PendingCount:    row.PendingCount,
			OpenCount:       row.OpenCount,
			DelayedCount:    row.DelayedCount,
			TotalTime:       row.TotalTime,
		}
	}

	// 2) Charge depuis ticket_internes (pour voir la performance des départements non-IT)
	baseInternal := database.DB.Table("ticket_internes").
		Where("assigned_to_id IS NOT NULL AND created_at >= ?", start)
	if scopeParam != nil {
		if queryScope, ok := scopeParam.(*scope.QueryScope); ok {
			baseInternal = scope.ApplyTicketInternalScopeToTable(baseInternal, queryScope)
		}
	}
	var internalRows []workloadRow
	// Même règle : résolus = cloture ou resolu
	internalQuery := baseInternal.Select(`
		assigned_to_id,
		COUNT(*) as ticket_count,
		SUM(CASE WHEN status IN ('cloture', 'resolu') THEN 1 ELSE 0 END) as resolved_count,
		SUM(CASE WHEN status = 'en_cours' THEN 1 ELSE 0 END) as in_progress_count,
		SUM(CASE WHEN status = 'en_attente' THEN 1 ELSE 0 END) as pending_count,
		SUM(CASE WHEN status = 'ouvert' THEN 1 ELSE 0 END) as open_count,
		SUM(CASE WHEN status NOT IN ('cloture', 'resolu') AND closed_at IS NULL AND DATEDIFF(NOW(), created_at) > 7 THEN 1 ELSE 0 END) as delayed_count,
		AVG(CASE WHEN status IN ('cloture', 'resolu') THEN COALESCE(NULLIF(actual_time, 0), TIMESTAMPDIFF(MINUTE, created_at, COALESCE(closed_at, updated_at))) ELSE NULL END) as average_time,
		SUM(CASE WHEN status IN ('cloture', 'resolu') THEN COALESCE(NULLIF(actual_time, 0), TIMESTAMPDIFF(MINUTE, created_at, COALESCE(closed_at, updated_at))) ELSE 0 END) as total_time
	`).Group("assigned_to_id")
	if err := internalQuery.Scan(&internalRows).Error; err != nil {
		return nil, err
	}
	for _, row := range internalRows {
		if row.AssignedToID == nil || !inDepartment(*row.AssignedToID) {
			continue
		}
		uid := *row.AssignedToID
		if acc, ok := userAcc[uid]; ok {
			acc.TicketCount += row.TicketCount
			acc.ResolvedCount += row.ResolvedCount
			acc.InProgressCount += row.InProgressCount
			acc.PendingCount += row.PendingCount
			acc.OpenCount += row.OpenCount
			acc.DelayedCount += row.DelayedCount
			acc.TotalTime += row.TotalTime
		} else {
			userAcc[uid] = &workloadAcc{
				TicketCount:     row.TicketCount,
				ResolvedCount:   row.ResolvedCount,
				InProgressCount: row.InProgressCount,
				PendingCount:    row.PendingCount,
				OpenCount:       row.OpenCount,
				DelayedCount:    row.DelayedCount,
				TotalTime:       row.TotalTime,
			}
		}
	}

	// 3) Construire les DTO et calculer efficacité / temps moyen
	results := make([]dto.WorkloadByAgentDTO, 0, len(userAcc))
	for uid, acc := range userAcc {
		user, err := s.userRepo.FindByID(uid)
		if err != nil {
			continue
		}
		userDTO := userToDTO(user)
		efficiency := 0.0
		if acc.TicketCount > 0 {
			efficiency = (float64(acc.ResolvedCount) / float64(acc.TicketCount)) * 100
		}
		avgTime := 0.0
		if acc.ResolvedCount > 0 {
			avgTime = acc.TotalTime / float64(acc.ResolvedCount)
		}
		results = append(results, dto.WorkloadByAgentDTO{
			UserID:          uid,
			User:            &userDTO,
			TicketCount:     acc.TicketCount,
			ResolvedCount:   acc.ResolvedCount,
			InProgressCount: acc.InProgressCount,
			PendingCount:    acc.PendingCount,
			OpenCount:       acc.OpenCount,
			DelayedCount:    acc.DelayedCount,
			AverageTime:     avgTime,
			TotalTime:       int(math.Round(acc.TotalTime)),
			Efficiency:      efficiency,
		})
	}

	// 4) Ajouter les membres du département sans aucun ticket (ni normal ni interne)
	for _, uid := range departmentUserIDs {
		if _, ok := userAcc[uid]; ok {
			continue
		}
		user, err := s.userRepo.FindByID(uid)
		if err != nil {
			continue
		}
		userDTO := userToDTO(user)
		results = append(results, dto.WorkloadByAgentDTO{
			UserID:          uid,
			User:            &userDTO,
			TicketCount:     0,
			ResolvedCount:   0,
			InProgressCount: 0,
			PendingCount:    0,
			OpenCount:       0,
			DelayedCount:    0,
			AverageTime:     0,
			TotalTime:       0,
			Efficiency:      0,
		})
	}

	return results, nil
}

// GetSLAComplianceReport récupère le rapport de conformité SLA
// Le scope est utilisé pour filtrer automatiquement selon les permissions de l'utilisateur
func (s *reportService) GetSLAComplianceReport(scopeParam interface{}, period string) (*dto.SLAComplianceReportDTO, error) {
	now := time.Now()
	start := periodStart(period, now)

	type statusRow struct {
		Status string `gorm:"column:status"`
		Count  int    `gorm:"column:count"`
	}
	var statusRows []statusRow

	// Construire la requête de base
	baseQuery := database.DB.Table("ticket_sla").
		Select("ticket_sla.status, COUNT(*) as count").
		Joins("INNER JOIN tickets ON tickets.id = ticket_sla.ticket_id").
		Where("ticket_sla.created_at >= ?", start)

	// Appliquer le scope si fourni (sur les tickets)
	if scopeParam != nil {
		if queryScope, ok := scopeParam.(*scope.QueryScope); ok {
			baseQuery = scope.ApplyTicketScopeToTable(baseQuery, queryScope)
		}
	}

	if err := baseQuery.Group("ticket_sla.status").Scan(&statusRows).Error; err != nil {
		return nil, err
	}

	statusCounts := map[string]int{}
	totalTickets := 0
	for _, row := range statusRows {
		statusCounts[row.Status] = row.Count
		totalTickets += row.Count
	}

	totalViolations := statusCounts["violated"]
	overallCompliance := 0.0
	if totalTickets > 0 {
		overallCompliance = (float64(totalTickets-totalViolations) / float64(totalTickets)) * 100
	}

	byCategory := make(map[string]float64)
	type categoryRow struct {
		Category string `gorm:"column:category"`
		Total    int    `gorm:"column:total"`
		Viol     int    `gorm:"column:violations"`
	}
	var categoryRows []categoryRow
	if err := database.DB.
		Table("ticket_sla").
		Select("tickets.category as category, COUNT(*) as total, SUM(CASE WHEN ticket_sla.status = 'violated' THEN 1 ELSE 0 END) as violations").
		Joins("JOIN tickets ON tickets.id = ticket_sla.ticket_id").
		Where("ticket_sla.created_at >= ?", start).
		Group("tickets.category").
		Scan(&categoryRows).Error; err != nil {
		return nil, err
	}
	for _, row := range categoryRows {
		if row.Total > 0 {
			byCategory[row.Category] = (float64(row.Total-row.Viol) / float64(row.Total)) * 100
		}
	}

	byPriority := make(map[string]float64)
	type priorityRow struct {
		Priority string `gorm:"column:priority"`
		Total    int    `gorm:"column:total"`
		Viol     int    `gorm:"column:violations"`
	}
	var priorityRows []priorityRow
	if err := database.DB.
		Table("ticket_sla").
		Select("tickets.priority as priority, COUNT(*) as total, SUM(CASE WHEN ticket_sla.status = 'violated' THEN 1 ELSE 0 END) as violations").
		Joins("JOIN tickets ON tickets.id = ticket_sla.ticket_id").
		Where("ticket_sla.created_at >= ?", start).
		Group("tickets.priority").
		Scan(&priorityRows).Error; err != nil {
		return nil, err
	}
	for _, row := range priorityRows {
		if row.Total > 0 {
			byPriority[row.Priority] = (float64(row.Total-row.Viol) / float64(row.Total)) * 100
		}
	}

	return &dto.SLAComplianceReportDTO{
		OverallCompliance: overallCompliance,
		ByCategory:        byCategory,
		ByPriority:        byPriority,
		TotalTickets:      totalTickets,
		TotalViolations:   totalViolations,
		Period:            normalizePeriod(period),
		GeneratedAt:       time.Now(),
	}, nil
}

// GetDelayedTicketsReport récupère le rapport des tickets en retard
// Le scope est utilisé pour filtrer automatiquement selon les permissions de l'utilisateur
func (s *reportService) GetDelayedTicketsReport(scopeParam interface{}, period string) ([]dto.DelayedTicketDTO, error) {
	now := time.Now()
	start := periodStart(period, now)

	// Construire la requête de base
	baseQuery := database.DB.Model(&models.TicketSLA{}).
		Preload("Ticket").
		Where("status = ? AND created_at >= ?", "violated", start)

	// Appliquer le scope si fourni (filtrer sur les tickets)
	if scopeParam != nil {
		if queryScope, ok := scopeParam.(*scope.QueryScope); ok {
			// Joindre avec tickets pour appliquer le scope
			baseQuery = baseQuery.Joins("INNER JOIN tickets ON tickets.id = ticket_sla.ticket_id")
			baseQuery = scope.ApplyTicketScopeToTable(baseQuery, queryScope)
		}
	}

	var delayedTickets []models.TicketSLA
	if err := baseQuery.Find(&delayedTickets).Error; err != nil {
		return nil, err
	}

	results := make([]dto.DelayedTicketDTO, 0, len(delayedTickets))
	for _, tsla := range delayedTickets {
		results = append(results, dto.DelayedTicketDTO{
			TicketID:     tsla.TicketID,
			Ticket:       &dto.TicketDTO{ID: tsla.Ticket.ID, Title: tsla.Ticket.Title, Category: tsla.Ticket.Category, Status: tsla.Ticket.Status, CreatedAt: tsla.Ticket.CreatedAt},
			ExpectedDate: tsla.TargetTime,
			DelayedBy:    0,
			Priority:     tsla.Ticket.Priority,
			Category:     tsla.Ticket.Category,
		})
	}

	return results, nil
}

// GetIndividualPerformanceReport récupère le rapport de performance individuel
func (s *reportService) GetIndividualPerformanceReport(userID uint, period string) (*dto.IndividualPerformanceReportDTO, error) {
	// TODO: Implémenter le calcul du rapport de performance individuel
	return &dto.IndividualPerformanceReportDTO{
		UserID: userID,
		Period: normalizePeriod(period),
	}, nil
}

// GetAssetSummary récupère un résumé des actifs (filtré par scope si fourni)
func (s *reportService) GetAssetSummary(scopeParam interface{}, period string) (*dto.AssetReportDTO, error) {
	now := time.Now()
	start := periodStart(period, now)

	baseQuery := database.DB.Model(&models.Asset{}).Where("created_at >= ?", start)
	if scopeParam != nil {
		if queryScope, ok := scopeParam.(*scope.QueryScope); ok {
			baseQuery = scope.ApplyAssetScope(baseQuery, queryScope)
		}
	}

	var total int64
	if err := baseQuery.Count(&total).Error; err != nil {
		log.Printf("GetAssetSummary count error: %v", err)
		return nil, err
	}

	type statusRow struct {
		Status string `gorm:"column:status"`
		Count  int    `gorm:"column:count"`
	}
	statusQuery := database.DB.Model(&models.Asset{}).Where("created_at >= ?", start)
	if scopeParam != nil {
		if queryScope, ok := scopeParam.(*scope.QueryScope); ok {
			statusQuery = scope.ApplyAssetScope(statusQuery, queryScope)
		}
	}
	var statusRows []statusRow
	if err := statusQuery.Select("status, COUNT(*) as count").Group("status").Scan(&statusRows).Error; err != nil {
		log.Printf("GetAssetSummary status group error: %v", err)
		return nil, err
	}
	byStatus := map[string]int{}
	for _, row := range statusRows {
		byStatus[row.Status] = row.Count
	}

	type categoryRow struct {
		Name  string `gorm:"column:name"`
		Count int    `gorm:"column:count"`
	}
	catQuery := database.DB.Table("assets").Select("asset_categories.name as name, COUNT(*) as count").
		Joins("JOIN asset_categories ON asset_categories.id = assets.category_id").
		Where("assets.created_at >= ?", start)
	if scopeParam != nil {
		if queryScope, ok := scopeParam.(*scope.QueryScope); ok {
			catQuery = scope.ApplyAssetScope(catQuery, queryScope)
		}
	}
	var categoryRows []categoryRow
	if err := catQuery.Group("asset_categories.name").Scan(&categoryRows).Error; err != nil {
		log.Printf("GetAssetSummary category group error: %v", err)
		return nil, err
	}
	byCategory := map[string]int{}
	for _, row := range categoryRows {
		byCategory[row.Name] = row.Count
	}

	return &dto.AssetReportDTO{
		Period:     normalizePeriod(period),
		Total:      int(total),
		ByStatus:   byStatus,
		ByCategory: byCategory,
	}, nil
}

// GetKnowledgeSummary récupère un résumé de la base de connaissances (filtré par scope si fourni)
func (s *reportService) GetKnowledgeSummary(scopeParam interface{}, period string) (*dto.KnowledgeReportDTO, error) {
	now := time.Now()
	start := periodStart(period, now)

	baseQuery := database.DB.Model(&models.KnowledgeArticle{}).Where("created_at >= ?", start)
	if scopeParam != nil {
		if queryScope, ok := scopeParam.(*scope.QueryScope); ok {
			baseQuery = scope.ApplyKnowledgeScope(baseQuery, queryScope)
		}
	}

	var total int64
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, err
	}

	publishedQuery := database.DB.Model(&models.KnowledgeArticle{}).Where("created_at >= ? AND is_published = ?", start, true)
	if scopeParam != nil {
		if queryScope, ok := scopeParam.(*scope.QueryScope); ok {
			publishedQuery = scope.ApplyKnowledgeScope(publishedQuery, queryScope)
		}
	}
	var published int64
	if err := publishedQuery.Count(&published).Error; err != nil {
		return nil, err
	}

	draftQuery := database.DB.Model(&models.KnowledgeArticle{}).Where("created_at >= ? AND is_published = ?", start, false)
	if scopeParam != nil {
		if queryScope, ok := scopeParam.(*scope.QueryScope); ok {
			draftQuery = scope.ApplyKnowledgeScope(draftQuery, queryScope)
		}
	}
	var draft int64
	if err := draftQuery.Count(&draft).Error; err != nil {
		return nil, err
	}

	type categoryRow struct {
		Name  string `gorm:"column:name"`
		Count int    `gorm:"column:count"`
	}
	catQuery := database.DB.Table("knowledge_articles").
		Select("knowledge_categories.name as name, COUNT(*) as count").
		Joins("JOIN knowledge_categories ON knowledge_categories.id = knowledge_articles.category_id").
		Where("knowledge_articles.created_at >= ?", start)
	if scopeParam != nil {
		if queryScope, ok := scopeParam.(*scope.QueryScope); ok {
			catQuery = scope.ApplyKnowledgeScope(catQuery, queryScope)
		}
	}
	var categoryRows []categoryRow
	if err := catQuery.Group("knowledge_categories.name").Scan(&categoryRows).Error; err != nil {
		return nil, err
	}
	byCategory := map[string]int{}
	for _, row := range categoryRows {
		byCategory[row.Name] = row.Count
	}

	return &dto.KnowledgeReportDTO{
		Period:     normalizePeriod(period),
		Total:      int(total),
		Published:  int(published),
		Draft:      int(draft),
		ByCategory: byCategory,
	}, nil
}

// getDepartmentUserIDs retourne les IDs des utilisateurs actifs du département quand scope = tableau de bord département
func getDepartmentUserIDs(scopeParam interface{}) ([]uint, bool) {
	if scopeParam == nil {
		return nil, false
	}
	queryScope, ok := scopeParam.(*scope.QueryScope)
	if !ok || queryScope.DashboardScopeHint != "department" || queryScope.DepartmentID == nil {
		return nil, false
	}
	var ids []uint
	if err := database.DB.Model(&models.User{}).Where("department_id = ? AND is_active = ?", *queryScope.DepartmentID, true).Pluck("id", &ids).Error; err != nil {
		return nil, false
	}
	return ids, true
}

func (s *reportService) countTicketsSince(scopeParam interface{}, start time.Time) (int, error) {
	var total int64
	query := database.DB.Model(&models.Ticket{})

	if deptIDs, ok := getDepartmentUserIDs(scopeParam); ok && len(deptIDs) > 0 {
		// Tableau de bord département : compter les tickets créés par les membres du département
		query = query.Where("created_by_id IN ?", deptIDs)
	} else if scopeParam != nil {
		if queryScope, ok := scopeParam.(*scope.QueryScope); ok {
			query = scope.ApplyTicketScope(query, queryScope)
		}
	}

	if !start.IsZero() {
		query = query.Where("created_at >= ?", start)
	}
	if err := query.Count(&total).Error; err != nil {
		return 0, err
	}
	return int(total), nil
}

func (s *reportService) countTicketsByFieldSince(scopeParam interface{}, field string, start time.Time) (map[string]int, error) {
	type row struct {
		Value string `gorm:"column:value"`
		Count int    `gorm:"column:count"`
	}
	var rows []row
	query := database.DB.Model(&models.Ticket{}).Select(fmt.Sprintf("%s as value, COUNT(*) as count", field))

	if deptIDs, ok := getDepartmentUserIDs(scopeParam); ok && len(deptIDs) > 0 {
		query = query.Where("created_by_id IN ?", deptIDs)
	} else if scopeParam != nil {
		if queryScope, ok := scopeParam.(*scope.QueryScope); ok {
			query = scope.ApplyTicketScope(query, queryScope)
		}
	}

	if !start.IsZero() {
		query = query.Where("created_at >= ?", start)
	}
	if err := query.Group(field).Scan(&rows).Error; err != nil {
		return nil, err
	}
	results := make(map[string]int, len(rows))
	for _, r := range rows {
		results[r.Value] = r.Count
	}
	return results, nil
}

func userToDTO(user *models.User) dto.UserDTO {
	userDTO := dto.UserDTO{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		DepartmentID: user.DepartmentID,
		Avatar:       user.Avatar,
		Role:         user.Role.Name,
		// Pas besoin des permissions complètes dans les rapports
		IsActive:  user.IsActive,
		LastLogin: user.LastLogin,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	if user.Department != nil {
		departmentDTO := dto.DepartmentDTO{
			ID:          user.Department.ID,
			Name:        user.Department.Name,
			Code:        user.Department.Code,
			Description: user.Department.Description,
			OfficeID:    user.Department.OfficeID,
			IsActive:    user.Department.IsActive,
			CreatedAt:   user.Department.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   user.Department.UpdatedAt.Format(time.RFC3339),
		}
		if user.Department.Office != nil {
			departmentDTO.Office = &dto.OfficeDTO{
				ID:        user.Department.Office.ID,
				Name:      user.Department.Office.Name,
				Country:   user.Department.Office.Country,
				City:      user.Department.Office.City,
				Commune:   user.Department.Office.Commune,
				Address:   user.Department.Office.Address,
				Longitude: user.Department.Office.Longitude,
				Latitude:  user.Department.Office.Latitude,
				IsActive:  user.Department.Office.IsActive,
				CreatedAt: user.Department.Office.CreatedAt.Format(time.RFC3339),
				UpdatedAt: user.Department.Office.UpdatedAt.Format(time.RFC3339),
			}
		}
		userDTO.Department = &departmentDTO
	}

	return userDTO
}

// ExportReport exporte un rapport dans un format spécifique
func (s *reportService) ExportReport(reportType, format, period string) (interface{}, error) {
	// TODO: Implémenter l'export de rapport (PDF, Excel, CSV)
	return nil, errors.New("export de rapport non implémenté")
}

// GenerateCustomReport génère un rapport personnalisé
func (s *reportService) GenerateCustomReport(req dto.CustomReportRequest) (interface{}, error) {
	// TODO: Implémenter la génération de rapport personnalisé
	return nil, errors.New("rapport personnalisé non implémenté")
}
