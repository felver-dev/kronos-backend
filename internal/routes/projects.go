package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/handlers"
	"github.com/mcicare/itsm-backend/internal/middleware"
)

// SetupProjectRoutes configure les routes des projets
func SetupProjectRoutes(router *gin.RouterGroup, projectHandler *handlers.ProjectHandler) {
	projects := router.Group("/projects")
	projects.Use(middleware.AuthMiddleware())
	{
		projects.GET("", projectHandler.GetAll)
		projects.GET("/:id", projectHandler.GetByID)
		projects.GET("/:id/budget-extensions", projectHandler.GetBudgetExtensions)
		projects.POST("", projectHandler.Create)
		projects.POST("/:id/budget-extensions", projectHandler.AddBudgetExtension)
		projects.PUT("/:id/budget-extensions/:extId", projectHandler.UpdateBudgetExtension)
		projects.DELETE("/:id/budget-extensions/:extId", projectHandler.DeleteBudgetExtension)
		projects.PUT("/:id", projectHandler.Update)
		projects.DELETE("/:id", projectHandler.Delete)

		// Phases — /reorder avant /:phaseId
		projects.GET("/:id/phases", projectHandler.GetPhases)
		projects.POST("/:id/phases", projectHandler.CreatePhase)
		projects.PUT("/:id/phases/reorder", projectHandler.ReorderPhases)
		projects.PUT("/:id/phases/:phaseId", projectHandler.UpdatePhase)
		projects.DELETE("/:id/phases/:phaseId", projectHandler.DeletePhase)

		// Functions
		projects.GET("/:id/functions", projectHandler.GetFunctions)
		projects.POST("/:id/functions", projectHandler.CreateFunction)
		projects.PUT("/:id/functions/:functionId", projectHandler.UpdateFunction)
		projects.DELETE("/:id/functions/:functionId", projectHandler.DeleteFunction)

		// Members
		projects.GET("/:id/members", projectHandler.GetMembers)
		projects.POST("/:id/members", projectHandler.AddMember)
		projects.DELETE("/:id/members/:userId", projectHandler.RemoveMember)
		projects.PUT("/:id/members/:userId/function", projectHandler.SetMemberFunction)
		projects.PUT("/:id/members/:userId/set-project-manager", projectHandler.SetProjectManager)
		projects.PUT("/:id/members/:userId/set-lead", projectHandler.SetLead)

		// Phase members
		projects.GET("/:id/phases/:phaseId/members", projectHandler.GetPhaseMembers)
		projects.POST("/:id/phases/:phaseId/members", projectHandler.AddPhaseMember)
		projects.DELETE("/:id/phases/:phaseId/members/:userId", projectHandler.RemovePhaseMember)
		projects.PUT("/:id/phases/:phaseId/members/:userId/function", projectHandler.SetPhaseMemberFunction)

		// Tasks
		projects.GET("/:id/tasks", projectHandler.GetTasks)
		projects.POST("/:id/tasks", projectHandler.CreateTask)
		projects.GET("/:id/phases/:phaseId/tasks", projectHandler.GetTasksByPhase)
		projects.PUT("/:id/tasks/:taskId", projectHandler.UpdateTask)
		projects.DELETE("/:id/tasks/:taskId", projectHandler.DeleteTask)
	}
}

