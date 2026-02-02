package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/services"
	"github.com/mcicare/itsm-backend/internal/utils"
)

// AuthHandler gère les handlers d'authentification
type AuthHandler struct {
	authService services.AuthService
	userService services.UserService
}

// NewAuthHandler crée une nouvelle instance de AuthHandler
func NewAuthHandler(authService services.AuthService, userService services.UserService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		userService: userService,
	}
}

// Login gère la connexion d'un utilisateur
// @Summary Connexion utilisateur
// @Description Authentifie un utilisateur avec son email et mot de passe, retourne un token JWT
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Identifiants de connexion (email et mot de passe)"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données de connexion invalides", err.Error())
		return
	}

	response, err := h.authService.Login(req)
	if err != nil {
		utils.UnauthorizedResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, response, "Connexion réussie")
}

// RefreshToken gère le rafraîchissement d'un token
// @Summary Rafraîchir le token
// @Description Rafraîchit un token JWT expiré
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Token de rafraîchissement"
// @Success 200 {object} map[string]string
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	newToken, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		utils.UnauthorizedResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{"token": newToken}, "Token rafraîchi avec succès")
}

// Logout gère la déconnexion d'un utilisateur
// @Summary Déconnexion utilisateur
// @Description Déconnecte un utilisateur et invalide son token
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// Récupérer l'ID de l'utilisateur depuis le contexte (défini par AuthMiddleware)
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	// Récupérer le token depuis le header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		utils.UnauthorizedResponse(c, "Token manquant")
		return
	}

	// Extraire le token (format: "Bearer <token>")
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 {
		utils.UnauthorizedResponse(c, "Format de token invalide")
		return
	}
	token := parts[1]

	// Hasher le token pour l'invalider
	tokenHash := utils.HashString(token)

	// Appeler le service de déconnexion
	err := h.authService.Logout(userID.(uint), tokenHash)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Erreur lors de la déconnexion", err.Error())
		return
	}

	utils.SuccessResponse(c, nil, "Déconnexion réussie")
}

// GetMe retourne les informations de l'utilisateur connecté
// @Summary Informations utilisateur connecté
// @Description Retourne les informations de l'utilisateur actuellement connecté
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dto.UserDTO
// @Failure 401 {object} utils.Response
// @Router /auth/me [get]
func (h *AuthHandler) GetMe(c *gin.Context) {
	// Récupérer l'ID de l'utilisateur depuis le contexte
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	// Récupérer l'utilisateur complet depuis la base de données
	user, err := h.userService.GetByID(userID.(uint))
	if err != nil {
		utils.UnauthorizedResponse(c, "Utilisateur introuvable")
		return
	}

	utils.SuccessResponse(c, user, "Informations utilisateur récupérées")
}

// Register gère l'inscription d'un nouvel utilisateur
// @Summary Inscription utilisateur
// @Description Crée un nouveau compte utilisateur et connecte automatiquement l'utilisateur
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Données d'inscription"
// @Success 201 {object} dto.RegisterResponse
// @Failure 400 {object} utils.Response
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données d'inscription invalides", err.Error())
		return
	}

	response, err := h.authService.Register(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.CreatedResponse(c, response, "Inscription réussie")
}
