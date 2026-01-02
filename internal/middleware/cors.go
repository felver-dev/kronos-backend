package middleware

import (
	"github.com/gin-gonic/gin"
)

// CORSMiddleware configure les en-têtes CORS pour permettre les requêtes cross-origin
// CORS (Cross-Origin Resource Sharing) permet à un navigateur d'autoriser
// les requêtes HTTP depuis une origine différente (domaine, port, protocole)
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Autoriser toutes les origines (en production, spécifier les domaines autorisés)
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		
		// Autoriser l'envoi de cookies et credentials
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		
		// En-têtes autorisés dans les requêtes
		c.Writer.Header().Set("Access-Control-Allow-Headers", 
			"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		
		// Méthodes HTTP autorisées
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		// Si c'est une requête OPTIONS (préflight), répondre immédiatement avec 204 No Content
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		// Continuer avec la requête normale
		c.Next()
	}
}

