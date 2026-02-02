package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/models"
	"github.com/mcicare/itsm-backend/internal/scope"
	"github.com/mcicare/itsm-backend/internal/services"
	"github.com/mcicare/itsm-backend/internal/utils"
)

// Alias pour la doc Swagger (évite "cannot find type definition: models.Project")
type project = models.Project

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
// @Summary Créer un projet
// @Description Crée un nouveau projet dans le système
// @Tags projects
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object true "Données du projet" SchemaExample({"name":"string","description":"string","total_budget_time":0})
// @Success 201 {object} project
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /projects [post]
func (h *ProjectHandler) Create(c *gin.Context) {
	var req struct {
		Name            string  `json:"name" binding:"required"`
		Description     string  `json:"description,omitempty"`
		TotalBudgetTime *int    `json:"total_budget_time,omitempty"`
		StartDate       *string `json:"start_date,omitempty"`
		EndDate         *string `json:"end_date,omitempty"`
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

	project, err := h.projectService.Create(req.Name, req.Description, req.TotalBudgetTime, req.StartDate, req.EndDate, createdByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.CreatedResponse(c, project, "Projet créé avec succès")
}

// GetByID récupère un projet par son ID
// @Summary Récupérer un projet par ID
// @Description Récupère un projet par son identifiant
// @Tags projects
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID du projet"
// @Success 200 {object} project
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /projects/{id} [get]
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
// @Summary Récupérer tous les projets
// @Description Récupère la liste de tous les projets. Query ?scope=own pour « Mon tableau de bord » (uniquement les projets où l'utilisateur est impliqué).
// @Tags projects
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param scope query string false "scope=own pour limiter aux projets de l'utilisateur connecté (Mon tableau de bord)"
// @Success 200 {array} project
// @Failure 500 {object} utils.Response
// @Router /projects [get]
func (h *ProjectHandler) GetAll(c *gin.Context) {
	queryScope := utils.GetScopeFromContext(c)
	switch c.Query("scope") {
	case "own":
		if queryScope != nil {
			scopeOwn := &scope.QueryScope{
				UserID:       queryScope.UserID,
				DepartmentID: queryScope.DepartmentID,
				FilialeID:    queryScope.FilialeID,
				Permissions:  []string{"projects.view_own"},
			}
			queryScope = scopeOwn
		}
	case "department":
		if queryScope != nil && queryScope.DepartmentID != nil {
			scopeDept := &scope.QueryScope{
				UserID:       queryScope.UserID,
				DepartmentID: queryScope.DepartmentID,
				FilialeID:    queryScope.FilialeID,
				Permissions:  []string{"projects.view_team"},
				DashboardScopeHint: "department",
			}
			queryScope = scopeDept
		}
	default:
		utils.ApplyDashboardScopeHint(c, queryScope)
	}

	projects, err := h.projectService.GetAll(queryScope)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des projets")
		return
	}

	utils.SuccessResponse(c, projects, "Projets récupérés avec succès")
}

// Update met à jour un projet
// @Summary Mettre à jour un projet
// @Description Met à jour les informations d'un projet
// @Tags projects
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID du projet"
// @Param request body object true "Données à mettre à jour" SchemaExample({"name":"string","description":"string","total_budget_time":0,"status":"active"})
// @Success 200 {object} project
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /projects/{id} [put]
func (h *ProjectHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	var req struct {
		Name            *string `json:"name,omitempty"`
		Description     *string `json:"description,omitempty"`
		TotalBudgetTime *int    `json:"total_budget_time,omitempty"`
		Status          *string `json:"status,omitempty"`
		StartDate       *string `json:"start_date,omitempty"`
		EndDate         *string `json:"end_date,omitempty"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	updatedByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	nameStr := ""
	if req.Name != nil {
		nameStr = *req.Name
	}
	descriptionStr := ""
	if req.Description != nil {
		descriptionStr = *req.Description
	}
	statusStr := ""
	if req.Status != nil {
		statusStr = *req.Status
	}

	project, err := h.projectService.Update(uint(id), nameStr, descriptionStr, req.TotalBudgetTime, statusStr, req.StartDate, req.EndDate, updatedByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, project, "Projet mis à jour avec succès")
}

// Delete supprime un projet
// @Summary Supprimer un projet
// @Description Supprime un projet. Le corps doit contenir le nom d'utilisateur pour confirmer.
// @Tags projects
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID du projet"
// @Param request body object true "Confirmation" SchemaExample({"username":"string"})
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /projects/{id} [delete]
func (h *ProjectHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	var req struct {
		Username string `json:"username" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Le nom d'utilisateur est requis pour confirmer la suppression", err.Error())
		return
	}

	ctxUsernameVal, exists := c.Get("username")
	if !exists || ctxUsernameVal == nil {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}
	ctxUsername, ok := ctxUsernameVal.(string)
	if !ok {
		utils.ErrorResponse(c, http.StatusBadRequest, "Session invalide. Veuillez vous reconnecter puis réessayer.", nil)
		return
	}
	// Comparaison insensible à la casse et aux espaces (trim)
	want := strings.TrimSpace(ctxUsername)
	got := strings.TrimSpace(req.Username)
	if want == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Session invalide. Veuillez vous reconnecter puis réessayer.", nil)
		return
	}
	if !strings.EqualFold(want, got) {
		utils.ErrorResponse(c, http.StatusBadRequest, "Le nom d'utilisateur saisi ne correspond pas à votre compte. Utilisez exactement le nom d'utilisateur de connexion (pas l'email).", nil)
		return
	}

	if err := h.projectService.Delete(uint(id)); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, nil, "Projet supprimé avec succès")
}

// AddBudgetExtension ajoute une extension au budget temps du projet (temps + justification)
func (h *ProjectHandler) AddBudgetExtension(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	var req struct {
		AdditionalMinutes int     `json:"additional_minutes" binding:"required,gt=0"`
		Justification     string  `json:"justification" binding:"required,min=3"`
		StartDate         *string `json:"start_date,omitempty"` // Période de l'extension (optionnel)
		EndDate           *string `json:"end_date,omitempty"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides : temps strictement positif et justification d'au moins 3 caractères requis", err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	ext, err := h.projectService.AddBudgetExtension(uint(id), req.AdditionalMinutes, strings.TrimSpace(req.Justification), req.StartDate, req.EndDate, userID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.CreatedResponse(c, ext, "Budget étendu avec succès")
}

// GetBudgetExtensions retourne l'historique des extensions de budget d'un projet
func (h *ProjectHandler) GetBudgetExtensions(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	list, err := h.projectService.GetBudgetExtensions(uint(id))
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des extensions")
		return
	}

	utils.SuccessResponse(c, list, "Extensions récupérées")
}

// UpdateBudgetExtension modifie une extension de budget (temps, justification, période)
func (h *ProjectHandler) UpdateBudgetExtension(c *gin.Context) {
	if !utils.RequireAnyPermission(c, "projects.budget.extensions.update", "projects.budget.manage") {
		utils.ForbiddenResponse(c, "Permission insuffisante: projects.budget.extensions.update ou projects.budget.manage")
		return
	}
	idParam := c.Param("id")
	extIdParam := c.Param("extId")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID projet invalide")
		return
	}
	extID, err := strconv.ParseUint(extIdParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID extension invalide")
		return
	}
	var req struct {
		AdditionalMinutes int     `json:"additional_minutes" binding:"required,gt=0"`
		Justification     string  `json:"justification" binding:"required,min=3"`
		StartDate         *string `json:"start_date,omitempty"`
		EndDate           *string `json:"end_date,omitempty"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides : temps strictement positif et justification d'au moins 3 caractères requis", err.Error())
		return
	}
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}
	ext, err := h.projectService.UpdateBudgetExtension(uint(id), uint(extID), req.AdditionalMinutes, strings.TrimSpace(req.Justification), req.StartDate, req.EndDate, userID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	utils.SuccessResponse(c, ext, "Extension modifiée avec succès")
}

// DeleteBudgetExtension supprime une extension de budget
func (h *ProjectHandler) DeleteBudgetExtension(c *gin.Context) {
	if !utils.RequireAnyPermission(c, "projects.budget.extensions.delete", "projects.budget.manage") {
		utils.ForbiddenResponse(c, "Permission insuffisante: projects.budget.extensions.delete ou projects.budget.manage")
		return
	}
	idParam := c.Param("id")
	extIdParam := c.Param("extId")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID projet invalide")
		return
	}
	extID, err := strconv.ParseUint(extIdParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID extension invalide")
		return
	}
	if err := h.projectService.DeleteBudgetExtension(uint(id), uint(extID)); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	utils.SuccessResponse(c, nil, "Extension supprimée")
}

// --- Phases ---
func (h *ProjectHandler) GetPhases(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	list, err := h.projectService.GetPhases(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	utils.SuccessResponse(c, list, "")
}

func (h *ProjectHandler) CreatePhase(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var req struct {
		Name         string `json:"name" binding:"required"`
		Description  string `json:"description"`
		DisplayOrder int    `json:"display_order"`
		Status       string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Données invalides")
		return
	}
	if req.Status == "" {
		req.Status = "not_started"
	}
	p, err := h.projectService.CreatePhase(uint(id), req.Name, req.Description, req.DisplayOrder, req.Status)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	utils.CreatedResponse(c, p, "Étape créée")
}

func (h *ProjectHandler) UpdatePhase(c *gin.Context) {
	pid, _ := strconv.ParseUint(c.Param("phaseId"), 10, 32)
	var req struct {
		Name         string `json:"name"`
		Description  string `json:"description"`
		DisplayOrder *int   `json:"display_order"`
		Status       string `json:"status"`
	}
	_ = c.ShouldBindJSON(&req)
	p, err := h.projectService.UpdatePhase(uint(pid), req.Name, req.Description, req.DisplayOrder, req.Status)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	utils.SuccessResponse(c, p, "Étape mise à jour")
}

func (h *ProjectHandler) DeletePhase(c *gin.Context) {
	pid, _ := strconv.ParseUint(c.Param("phaseId"), 10, 32)
	if err := h.projectService.DeletePhase(uint(pid)); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	utils.SuccessResponse(c, nil, "Étape supprimée")
}

func (h *ProjectHandler) ReorderPhases(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var req struct {
		Order []uint `json:"order" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "order (tableau d'IDs) requis")
		return
	}
	if err := h.projectService.ReorderPhases(uint(id), req.Order); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	utils.SuccessResponse(c, nil, "Ordre enregistré")
}

// --- Functions ---
func (h *ProjectHandler) GetFunctions(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	list, err := h.projectService.GetFunctions(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	utils.SuccessResponse(c, list, "")
}

func (h *ProjectHandler) CreateFunction(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var req struct {
		Name         string `json:"name" binding:"required"`
		Type         string `json:"type"` // "direction" | "execution"
		DisplayOrder int    `json:"display_order"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Données invalides")
		return
	}
	typeStr := req.Type
	if typeStr != "direction" && typeStr != "execution" {
		typeStr = "execution"
	}
	f, err := h.projectService.CreateFunction(uint(id), req.Name, typeStr, req.DisplayOrder)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	utils.CreatedResponse(c, f, "Fonction créée")
}

func (h *ProjectHandler) UpdateFunction(c *gin.Context) {
	fid, _ := strconv.ParseUint(c.Param("functionId"), 10, 32)
	var req struct {
		Name         string  `json:"name"`
		Type         *string `json:"type"` // "direction" | "execution"
		DisplayOrder *int    `json:"display_order"`
	}
	_ = c.ShouldBindJSON(&req)
	f, err := h.projectService.UpdateFunction(uint(fid), req.Name, req.Type, req.DisplayOrder)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	utils.SuccessResponse(c, f, "Fonction mise à jour")
}

func (h *ProjectHandler) DeleteFunction(c *gin.Context) {
	fid, _ := strconv.ParseUint(c.Param("functionId"), 10, 32)
	if err := h.projectService.DeleteFunction(uint(fid)); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	utils.SuccessResponse(c, nil, "Fonction supprimée")
}

// --- Members ---
func (h *ProjectHandler) GetMembers(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	list, err := h.projectService.GetMembers(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	utils.SuccessResponse(c, list, "")
}

func (h *ProjectHandler) AddMember(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var req struct {
		UserID            uint   `json:"user_id" binding:"required"`
		FunctionIDs       []uint `json:"function_ids"`
		ProjectFunctionID *uint  `json:"project_function_id"` // rétrocompat: si function_ids vide, on utilise celui-ci
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "user_id requis")
		return
	}
	if req.FunctionIDs == nil {
		req.FunctionIDs = []uint{}
	}
	if len(req.FunctionIDs) == 0 && req.ProjectFunctionID != nil {
		req.FunctionIDs = []uint{*req.ProjectFunctionID}
	}
	m, err := h.projectService.AddMember(uint(id), req.UserID, req.FunctionIDs)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	utils.CreatedResponse(c, m, "Membre ajouté")
}

func (h *ProjectHandler) RemoveMember(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	uid, _ := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err := h.projectService.RemoveMember(uint(id), uint(uid)); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	utils.SuccessResponse(c, nil, "Membre retiré")
}

func (h *ProjectHandler) SetMemberFunction(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	uid, _ := strconv.ParseUint(c.Param("userId"), 10, 32)
	var req struct {
		FunctionIDs []uint `json:"function_ids"`
	}
	_ = c.ShouldBindJSON(&req)
	if req.FunctionIDs == nil {
		req.FunctionIDs = []uint{}
	}
	if err := h.projectService.SetMemberFunctions(uint(id), uint(uid), req.FunctionIDs); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	utils.SuccessResponse(c, nil, "Fonctions mises à jour")
}

func (h *ProjectHandler) SetProjectManager(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	uid, _ := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err := h.projectService.SetProjectManager(uint(id), uint(uid)); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	utils.SuccessResponse(c, nil, "Chef de projet mis à jour")
}

func (h *ProjectHandler) SetLead(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	uid, _ := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err := h.projectService.SetLead(uint(id), uint(uid)); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	utils.SuccessResponse(c, nil, "Lead mis à jour")
}

// --- Phase members ---
func (h *ProjectHandler) GetPhaseMembers(c *gin.Context) {
	pid, _ := strconv.ParseUint(c.Param("phaseId"), 10, 32)
	list, err := h.projectService.GetPhaseMembers(uint(pid))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	utils.SuccessResponse(c, list, "")
}

func (h *ProjectHandler) AddPhaseMember(c *gin.Context) {
	pid, _ := strconv.ParseUint(c.Param("phaseId"), 10, 32)
	var req struct {
		UserID             uint  `json:"user_id" binding:"required"`
		ProjectFunctionID  *uint `json:"project_function_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "user_id requis")
		return
	}
	m, err := h.projectService.AddPhaseMember(uint(pid), req.UserID, req.ProjectFunctionID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	utils.CreatedResponse(c, m, "Membre ajouté à l'étape")
}

func (h *ProjectHandler) RemovePhaseMember(c *gin.Context) {
	pid, _ := strconv.ParseUint(c.Param("phaseId"), 10, 32)
	uid, _ := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err := h.projectService.RemovePhaseMember(uint(pid), uint(uid)); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	utils.SuccessResponse(c, nil, "Membre retiré de l'étape")
}

func (h *ProjectHandler) SetPhaseMemberFunction(c *gin.Context) {
	pid, _ := strconv.ParseUint(c.Param("phaseId"), 10, 32)
	uid, _ := strconv.ParseUint(c.Param("userId"), 10, 32)
	var req struct {
		ProjectFunctionID *uint `json:"project_function_id"`
	}
	_ = c.ShouldBindJSON(&req)
	if err := h.projectService.SetPhaseMemberFunction(uint(pid), uint(uid), req.ProjectFunctionID); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	utils.SuccessResponse(c, nil, "Fonction mise à jour")
}

// --- Tasks ---
func (h *ProjectHandler) GetTasks(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	list, err := h.projectService.GetTasks(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	utils.SuccessResponse(c, list, "")
}

func (h *ProjectHandler) GetTasksByPhase(c *gin.Context) {
	pid, _ := strconv.ParseUint(c.Param("phaseId"), 10, 32)
	list, err := h.projectService.GetTasksByPhase(uint(pid))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	utils.SuccessResponse(c, list, "")
}

func (h *ProjectHandler) CreateTask(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	userID, _ := c.Get("user_id")
	var req struct {
		ProjectPhaseID uint     `json:"project_phase_id" binding:"required"`
		Title          string   `json:"title" binding:"required"`
		Description    string   `json:"description"`
		Status         string   `json:"status"`
		Priority       string   `json:"priority"`
		AssigneeIDs    []uint   `json:"assignee_ids"`
		EstimatedTime  *int     `json:"estimated_time"`
		DueDate        *string  `json:"due_date"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "project_phase_id et title requis")
		return
	}
	if req.AssigneeIDs == nil {
		req.AssigneeIDs = []uint{}
	}
	t, err := h.projectService.CreateTask(uint(id), req.ProjectPhaseID, userID.(uint), req.Title, req.Description, req.Status, req.Priority, req.AssigneeIDs, req.EstimatedTime, req.DueDate)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	utils.CreatedResponse(c, t, "Tâche créée")
}

func (h *ProjectHandler) UpdateTask(c *gin.Context) {
	tid, _ := strconv.ParseUint(c.Param("taskId"), 10, 32)
	var req struct {
		Title          string   `json:"title"`
		Description    string   `json:"description"`
		Status         string   `json:"status"`
		Priority       string   `json:"priority"`
		AssigneeIDs    *[]uint  `json:"assignee_ids"`
		EstimatedTime  *int     `json:"estimated_time"`
		ActualTime     *int     `json:"actual_time"`
		DueDate        *string  `json:"due_date"`
		ProjectPhaseID *uint    `json:"project_phase_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}
	t, err := h.projectService.UpdateTask(uint(tid), req.Title, req.Description, req.Status, req.Priority, req.AssigneeIDs, req.EstimatedTime, req.ActualTime, req.DueDate, req.ProjectPhaseID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	utils.SuccessResponse(c, t, "Tâche mise à jour")
}

func (h *ProjectHandler) DeleteTask(c *gin.Context) {
	tid, _ := strconv.ParseUint(c.Param("taskId"), 10, 32)
	if err := h.projectService.DeleteTask(uint(tid)); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	utils.SuccessResponse(c, nil, "Tâche supprimée")
}
