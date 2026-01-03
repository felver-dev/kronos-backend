package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/services"
	"github.com/mcicare/itsm-backend/internal/utils"
)

// UserHandler gère les handlers des utilisateurs
type UserHandler struct {
	userService services.UserService
}

// NewUserHandler crée une nouvelle instance de UserHandler
func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// Create crée un nouvel utilisateur
// @Summary Créer un utilisateur
// @Description Crée un nouvel utilisateur dans le système
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateUserRequest true "Données de l'utilisateur"
// @Success 201 {object} dto.UserDTO
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /users [post]
func (h *UserHandler) Create(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	// Récupérer l'ID de l'utilisateur créateur depuis le contexte
	createdByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	user, err := h.userService.Create(req, createdByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.CreatedResponse(c, user, "Utilisateur créé avec succès")
}

// GetByID récupère un utilisateur par son ID
// @Summary Récupérer un utilisateur
// @Description Récupère les informations d'un utilisateur par son ID
// @Tags users
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID de l'utilisateur"
// @Success 200 {object} dto.UserDTO
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /users/{id} [get]
func (h *UserHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	user, err := h.userService.GetByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Utilisateur introuvable")
		return
	}

	utils.SuccessResponse(c, user, "Utilisateur récupéré avec succès")
}

// GetAll récupère tous les utilisateurs
// @Summary Liste des utilisateurs
// @Description Récupère la liste de tous les utilisateurs
// @Tags users
// @Security BearerAuth
// @Produce json
// @Success 200 {array} dto.UserDTO
// @Failure 500 {object} utils.Response
// @Router /users [get]
func (h *UserHandler) GetAll(c *gin.Context) {
	users, err := h.userService.GetAll()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des utilisateurs")
		return
	}

	utils.SuccessResponse(c, users, "Utilisateurs récupérés avec succès")
}

// Update met à jour un utilisateur
// @Summary Mettre à jour un utilisateur
// @Description Met à jour les informations d'un utilisateur
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID de l'utilisateur"
// @Param request body dto.UpdateUserRequest true "Données à mettre à jour"
// @Success 200 {object} dto.UserDTO
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	// Récupérer l'ID de l'utilisateur qui effectue la mise à jour
	updatedByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	user, err := h.userService.Update(uint(id), req, updatedByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, user, "Utilisateur mis à jour avec succès")
}

// Delete supprime un utilisateur
// @Summary Supprimer un utilisateur
// @Description Supprime un utilisateur du système
// @Tags users
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID de l'utilisateur"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	err = h.userService.Delete(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Utilisateur introuvable")
		return
	}

	utils.SuccessResponse(c, nil, "Utilisateur supprimé avec succès")
}

// ChangePassword change le mot de passe d'un utilisateur
// @Summary Changer le mot de passe
// @Description Change le mot de passe d'un utilisateur
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID de l'utilisateur"
// @Param request body map[string]string true "Ancien et nouveau mot de passe"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /users/{id}/password [put]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	err = h.userService.ChangePassword(uint(id), req.OldPassword, req.NewPassword)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, nil, "Mot de passe modifié avec succès")
}

