package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/utils"
)

// AuthMiddleware vérifie la présence et la validité du token JWT
// Si le token est valide, les informations de l'utilisateur sont stockées dans le contexte
// Les handlers peuvent ensuite accéder à ces informations via c.Get("user_id"), etc.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Récupérer le header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.UnauthorizedResponse(c, "Token d'authentification manquant")
			c.Abort()
			return
		}

		// Le format attendu est "Bearer <token>"
		// On sépare le header en deux parties
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.UnauthorizedResponse(c, "Format de token invalide. Attendu: Bearer <token>")
			c.Abort()
			return
		}

		// Extraire le token
		token := parts[1]

		// Valider le token et récupérer les claims
		claims, err := utils.ValidateToken(token)
		if err != nil {
			utils.UnauthorizedResponse(c, "Token invalide ou expiré")
			c.Abort()
			return
		}

		// Stocker les informations de l'utilisateur dans le contexte Gin
		// Ces informations seront accessibles dans les handlers suivants
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		// Continuer avec la requête
		c.Next()
	}
}
