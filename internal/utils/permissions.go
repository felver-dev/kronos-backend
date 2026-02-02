package utils

import (
	"github.com/gin-gonic/gin"
)

// RequirePermission vérifie si l'utilisateur a une permission spécifique
// Retourne true si l'utilisateur a la permission, false sinon
// Si le scope n'est pas trouvé, retourne false (sécurité par défaut)
func RequirePermission(c *gin.Context, permission string) bool {
	queryScope := GetScopeFromContext(c)
	if queryScope == nil {
		return false
	}
	return queryScope.HasPermission(permission)
}

// RequireAnyPermission vérifie si l'utilisateur a au moins une des permissions spécifiées
func RequireAnyPermission(c *gin.Context, permissions ...string) bool {
	queryScope := GetScopeFromContext(c)
	if queryScope == nil {
		return false
	}
	return queryScope.HasAnyPermission(permissions...)
}

// RequireAllPermissions vérifie si l'utilisateur a toutes les permissions spécifiées
func RequireAllPermissions(c *gin.Context, permissions ...string) bool {
	queryScope := GetScopeFromContext(c)
	if queryScope == nil {
		return false
	}
	return queryScope.HasAllPermissions(permissions...)
}

// PermissionMiddleware crée un middleware qui vérifie une permission avant d'autoriser l'accès
func PermissionMiddleware(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !RequirePermission(c, permission) {
			ForbiddenResponse(c, "Permission insuffisante: "+permission)
			c.Abort()
			return
		}
		c.Next()
	}
}

// AnyPermissionMiddleware crée un middleware qui vérifie qu'au moins une des permissions est présente
func AnyPermissionMiddleware(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !RequireAnyPermission(c, permissions...) {
			ForbiddenResponse(c, "Permissions insuffisantes")
			c.Abort()
			return
		}
		c.Next()
	}
}
