package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/handlers"
	"github.com/mcicare/itsm-backend/internal/middleware"
)

// SetupTimesheetRoutes configure les routes des timesheets
func SetupTimesheetRoutes(router *gin.RouterGroup, timesheetHandler *handlers.TimesheetHandler) {
	timesheet := router.Group("/timesheet")
	timesheet.Use(middleware.AuthMiddleware())
	{
		// Saisie du temps par ticket
		timesheet.POST("/entries", timesheetHandler.CreateTimeEntry)
		timesheet.GET("/entries", timesheetHandler.GetTimeEntries)
		timesheet.GET("/entries/:id", timesheetHandler.GetTimeEntryByID)
		timesheet.PUT("/entries/:id", timesheetHandler.UpdateTimeEntry)
		timesheet.GET("/entries/by-date/:date", timesheetHandler.GetTimeEntriesByDate)
		timesheet.POST("/entries/:id/validate", timesheetHandler.ValidateTimeEntry)
		timesheet.GET("/entries/pending-validation", timesheetHandler.GetPendingValidationEntries)

		// Déclaration par jour
		timesheet.GET("/daily/:date", timesheetHandler.GetDailyDeclaration)
		timesheet.POST("/daily/:date", timesheetHandler.CreateOrUpdateDailyDeclaration)
		timesheet.GET("/daily/:date/tasks", timesheetHandler.GetDailyTasks)
		timesheet.POST("/daily/:date/tasks", timesheetHandler.CreateDailyTask)
		timesheet.DELETE("/daily/:date/tasks/:taskId", timesheetHandler.DeleteDailyTask)
		timesheet.GET("/daily/:date/summary", timesheetHandler.GetDailySummary)
		timesheet.GET("/daily/calendar", timesheetHandler.GetDailyCalendar)
		timesheet.GET("/daily/range", timesheetHandler.GetDailyRange)

		// Déclaration par semaine
		timesheet.GET("/weekly/:week", timesheetHandler.GetWeeklyDeclaration)
		timesheet.POST("/weekly/:week", timesheetHandler.CreateOrUpdateWeeklyDeclaration)
		timesheet.GET("/weekly/:week/tasks", timesheetHandler.GetWeeklyTasks)
		timesheet.GET("/weekly/:week/summary", timesheetHandler.GetWeeklySummary)
		timesheet.GET("/weekly/:week/daily-breakdown", timesheetHandler.GetWeeklyDailyBreakdown)
		timesheet.POST("/weekly/:week/validate", timesheetHandler.ValidateWeeklyDeclaration)
		timesheet.GET("/weekly/:week/validation-status", timesheetHandler.GetWeeklyValidationStatus)

		// Budget temps
		timesheet.GET("/budget-alerts", timesheetHandler.GetBudgetAlerts)
		timesheet.GET("/budget-status/:id", timesheetHandler.GetTicketBudgetStatus)

		// Validation
		timesheet.GET("/validation-history", timesheetHandler.GetValidationHistory)

		// Alertes
		timesheet.GET("/alerts/delays", timesheetHandler.GetDelayAlerts)
		timesheet.GET("/alerts/budget", timesheetHandler.GetBudgetAlertsForTimesheet)
		timesheet.GET("/alerts/overload", timesheetHandler.GetOverloadAlerts)
		timesheet.GET("/alerts/underload", timesheetHandler.GetUnderloadAlerts)
		timesheet.POST("/alerts/reminders", timesheetHandler.SendReminderAlerts)
		timesheet.GET("/alerts/justifications-pending", timesheetHandler.GetPendingJustificationAlerts)

		// Historique
		timesheet.GET("/history", timesheetHandler.GetTimesheetHistory)
		timesheet.GET("/history/:entryId", timesheetHandler.GetTimesheetHistoryEntry)
		timesheet.GET("/audit-trail", timesheetHandler.GetTimesheetAuditTrail)
		timesheet.GET("/modifications", timesheetHandler.GetTimesheetModifications)
	}
}

// SetupTicketTimesheetRoutes configure les routes de timesheet pour les tickets
func SetupTicketTimesheetRoutes(router *gin.RouterGroup, timesheetHandler *handlers.TimesheetHandler) {
	tickets := router.Group("/tickets")
	tickets.Use(middleware.AuthMiddleware())
	{
		// Routes spécifiques avec plus de segments - doivent être avant les routes génériques
		tickets.GET("/:id/time-entries", timesheetHandler.GetTimeEntriesByTicketID)
		tickets.POST("/:id/estimated-time", timesheetHandler.SetTicketEstimatedTime)
		tickets.GET("/:id/estimated-time", timesheetHandler.GetTicketEstimatedTime)
		tickets.PUT("/:id/estimated-time", timesheetHandler.UpdateTicketEstimatedTime)
		tickets.GET("/:id/time-comparison", timesheetHandler.GetTicketTimeComparison)
	}
}

// SetupUserTimesheetRoutes configure les routes de timesheet pour les utilisateurs
func SetupUserTimesheetRoutes(router *gin.RouterGroup, timesheetHandler *handlers.TimesheetHandler) {
	users := router.Group("/users")
	users.Use(middleware.AuthMiddleware())
	{
		// Route spécifique avec plus de segments - doit être avant les routes génériques
		users.GET("/:id/time-entries", timesheetHandler.GetTimeEntriesByUserID)
	}
}

// SetupProjectTimesheetRoutes configure les routes de timesheet pour les projets
func SetupProjectTimesheetRoutes(router *gin.RouterGroup, timesheetHandler *handlers.TimesheetHandler) {
	projects := router.Group("/projects")
	projects.Use(middleware.AuthMiddleware())
	{
		// Routes spécifiques avec plus de segments - doivent être avant les routes génériques
		projects.GET("/:id/time-budget", timesheetHandler.GetProjectTimeBudget)
		projects.POST("/:id/time-budget", timesheetHandler.SetProjectTimeBudget)
	}
}

