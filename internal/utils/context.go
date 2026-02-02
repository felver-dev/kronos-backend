package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/scope"
)

// ApplyDashboardScopeHint applique le paramètre query scope=department|filiale|global sur le QueryScope (pour le tableau de bord)
func ApplyDashboardScopeHint(c *gin.Context, queryScope *scope.QueryScope) {
	if queryScope == nil {
		return
	}
	switch c.Query("scope") {
	case "department", "filiale", "global":
		queryScope.DashboardScopeHint = c.Query("scope")
	}
}

// GetScopeFromContext extrait le QueryScope du contexte Gin
// Retourne nil si le scope n'est pas trouvé (ne devrait jamais arriver si AuthMiddleware est utilisé)
func GetScopeFromContext(c *gin.Context) *scope.QueryScope {
	scopeValue, exists := c.Get("scope")
	if !exists {
		return nil
	}

	queryScope, ok := scopeValue.(*scope.QueryScope)
	if !ok {
		return nil
	}

	return queryScope
}

// GetUserIDFromContext extrait l'ID utilisateur du contexte Gin
func GetUserIDFromContext(c *gin.Context) (uint, bool) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		return 0, false
	}

	return userID, true
}
