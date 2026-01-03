package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/services"
	"github.com/mcicare/itsm-backend/internal/utils"
)

// ProjectHandler gère les handlers des projets
type ProjectHandler struct {
	projectService services.ProjectService
}

// NewProjectHandler crée une nouvelle instance de ProjectHandler
func NewProjectHandler(projectService services.ProjectService) *ProjectHandler {
	return &ProjectHandler{
		projectService: projectService,
	}
}

// Create crée un nouveau projet
func (h *ProjectHandler) Create(c *gin.Context) {
	var req struct {
		Name            string `json:"name" binding:"required"`
		Description     string `json:"description,omitempty"`
		TotalBudgetTime *int   `json:"total_budget_time,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	createdByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	project, err := h.projectService.Create(req.Name, req.Description, req.TotalBudgetTime, createdByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.CreatedResponse(c, project, "Projet créé avec succès")
}

// GetByID récupère un projet par son ID
func (h *ProjectHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	project, err := h.projectService.GetByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Projet introuvable")
		return
	}

	utils.SuccessResponse(c, project, "Projet récupéré avec succès")
}

// GetAll récupère tous les projets
func (h *ProjectHandler) GetAll(c *gin.Context) {
	projects, err := h.projectService.GetAll()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des projets")
		return
	}

	utils.SuccessResponse(c, projects, "Projets récupérés avec succès")
}
