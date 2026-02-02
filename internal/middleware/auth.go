package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/repositories"
	"github.com/mcicare/itsm-backend/internal/scope"
	"github.com/mcicare/itsm-backend/internal/utils"
)

// AuthMiddleware vérifie la présence et la validité du token JWT
// Si le token est valide, les informations de l'utilisateur sont stockées dans le contexte
// Les handlers peuvent ensuite accéder à ces informations via c.Get("user_id"), etc.
// Le middleware enrichit également le contexte avec un QueryScope pour le filtrage automatique des données
func AuthMiddleware() gin.HandlerFunc {
	// Créer le repository une seule fois (singleton)
	userRepo := repositories.NewUserRepository()

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

		// Récupérer l'utilisateur complet avec ses relations (rôle, département)
		// pour construire le QueryScope
		user, err := userRepo.FindByID(claims.UserID)
		if err != nil {
			utils.UnauthorizedResponse(c, "Utilisateur introuvable")
			c.Abort()
			return
		}

		// Vérifier que l'utilisateur est actif
		if !user.IsActive {
			utils.UnauthorizedResponse(c, "Compte utilisateur désactivé")
			c.Abort()
			return
		}

		// Créer le QueryScope avec les permissions et attributs de l'utilisateur
		queryScope := scope.NewQueryScopeFromUser(user)

		// Stocker les informations de l'utilisateur dans le contexte Gin
		// On utilise user.Username (DB) et non claims.Username (JWT) pour avoir la valeur à jour
		// (en cas de changement de username après connexion, ou refresh de session)
		c.Set("user_id", claims.UserID)
		c.Set("username", user.Username)
		c.Set("role", claims.Role)
		c.Set("scope", queryScope) // Ajouter le QueryScope au contexte

		// Continuer avec la requête
		c.Next()
	}
}
