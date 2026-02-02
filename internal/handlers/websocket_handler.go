package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	gorillaWS "github.com/gorilla/websocket"
	"github.com/mcicare/itsm-backend/internal/utils"
	ws "github.com/mcicare/itsm-backend/internal/websocket"
)

var upgrader = gorillaWS.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// En production, vérifier l'origine de la requête
		// Pour le développement, accepter toutes les origines
		return true
	},
}

// WebSocketHandler gère les connexions WebSocket
type WebSocketHandler struct {
	Hub *ws.Hub
}

// NewWebSocketHandler crée une nouvelle instance de WebSocketHandler
func NewWebSocketHandler(hub *ws.Hub) *WebSocketHandler {
	return &WebSocketHandler{
		Hub: hub,
	}
}

// HandleWebSocket gère les connexions WebSocket
// @Summary Connexion WebSocket pour les notifications
// @Description Établit une connexion WebSocket pour recevoir les notifications en temps réel
// @Tags websocket
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 101 {object} string "Switching Protocols"
// @Failure 401 {object} utils.Response
// @Router /ws [get]
func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	// Extraire le token de la query string ou du header Authorization
	token := c.Query("token")
	if token == "" {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		}
	}

	if token == "" {
		utils.UnauthorizedResponse(c, "Token d'authentification requis")
		return
	}

	// Valider le token en utilisant utils.ValidateToken
	claims, err := utils.ValidateToken(token)
	if err != nil {
		utils.UnauthorizedResponse(c, "Token invalide ou expiré")
		return
	}

	userID := claims.UserID
	username := claims.Username
	if username == "" {
		username = "unknown"
	}

	// Mettre à niveau la connexion HTTP vers WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Erreur lors de la mise à niveau WebSocket: %v", err)
		utils.ErrorResponse(c, http.StatusBadRequest, "Erreur lors de la connexion WebSocket", nil)
		return
	}

	// Enregistrer le client dans le hub
	ws.ServeWs(h.Hub, conn, userID, username)
}
