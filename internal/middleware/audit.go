package middleware

import (
	"log"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/models"
	"github.com/mcicare/itsm-backend/internal/repositories"
)

// AuditLogMiddleware enregistre un log d'audit pour les requêtes mutantes
func AuditLogMiddleware(auditLogRepo repositories.AuditLogRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Ne logger que les opérations de modification réussies
		method := c.Request.Method
		if method == "GET" || method == "HEAD" || method == "OPTIONS" {
			return
		}
		if c.Writer.Status() >= 400 {
			return
		}

		path := c.Request.URL.Path
		if strings.Contains(path, "/audit-logs") {
			return
		}

		action := resolveAction(method, path)
		entityType, entityID := resolveEntity(path)

		var userID *uint
		if id, exists := c.Get("user_id"); exists {
			if uid, ok := id.(uint); ok {
				userID = &uid
			}
		}

		auditLog := &models.AuditLog{
			UserID:      userID,
			Action:      action,
			EntityType:  entityType,
			EntityID:    entityID,
			IPAddress:   c.ClientIP(),
			UserAgent:   c.GetHeader("User-Agent"),
			Description: method + " " + path,
		}

		if err := auditLogRepo.Create(auditLog); err != nil {
			log.Printf("⚠️  Audit log non enregistré: %v (action=%s entity=%s)", err, action, entityType)
		}
	}
}

func resolveAction(method, path string) string {
	lowerPath := strings.ToLower(path)
	switch {
	case strings.Contains(lowerPath, "/assign"):
		return "assign"
	case strings.Contains(lowerPath, "/validate"):
		return "validate"
	case strings.Contains(lowerPath, "/close"):
		return "close"
	case strings.Contains(lowerPath, "/status"):
		return "status_change"
	case method == "POST":
		return "create"
	case method == "PUT" || method == "PATCH":
		return "update"
	case method == "DELETE":
		return "delete"
	default:
		return "update"
	}
}

func resolveEntity(path string) (string, *uint) {
	segments := strings.Split(path, "/")
	clean := make([]string, 0, len(segments))
	for _, seg := range segments {
		if seg != "" {
			clean = append(clean, seg)
		}
	}

	if len(clean) == 0 {
		return "unknown", nil
	}

	// Ignorer /api/v1
	if len(clean) >= 2 && clean[0] == "api" && clean[1] == "v1" {
		clean = clean[2:]
	}

	entity := "unknown"
	if len(clean) > 0 {
		entity = normalizeEntity(clean[0])
	}

	for _, seg := range clean[1:] {
		if id, err := strconv.ParseUint(seg, 10, 32); err == nil {
			uid := uint(id)
			return entity, &uid
		}
	}

	return entity, nil
}

func normalizeEntity(segment string) string {
	switch segment {
	case "tickets":
		return "ticket"
	case "users":
		return "user"
	case "assets":
		return "asset"
	case "incidents":
		return "incident"
	case "changes":
		return "change"
	case "service-requests":
		return "service_request"
	case "knowledge-base":
		return "knowledge"
	case "time-entries":
		return "time_entry"
	case "audit-logs":
		return "audit"
	default:
		return segment
	}
}
