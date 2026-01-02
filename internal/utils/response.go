package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response est la structure standard pour toutes les réponses JSON
type Response struct {
	Success bool        `json:"success"`           // Indique si l'opération a réussi
	Message string      `json:"message,omitempty"` // Message optionnel
	Data    interface{} `json:"data,omitempty"`    // Données de la réponse
	Error   interface{} `json:"error,omitempty"`   // Erreur si échec
}

// PaginatedResponse est la structure pour les réponses paginées
type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

// Pagination contient les informations de pagination
type Pagination struct {
	Page       int   `json:"page"`        // Page actuelle
	Limit      int   `json:"limit"`       // Nombre d'éléments par page
	Total      int64 `json:"total"`       // Nombre total d'éléments
	TotalPages int   `json:"total_pages"` // Nombre total de pages
}

// SuccessResponse envoie une réponse de succès (200 OK)
func SuccessResponse(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// CreatedResponse envoie une réponse de création (201 Created)
func CreatedResponse(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse envoie une réponse d'erreur avec un code HTTP personnalisé
func ErrorResponse(c *gin.Context, statusCode int, message string, err interface{}) {
	c.JSON(statusCode, Response{
		Success: false,
		Message: message,
		Error:   err,
	})
}

// PaginatedSuccessResponse envoie une réponse paginée de succès
func PaginatedSuccessResponse(c *gin.Context, data interface{}, pagination Pagination) {
	c.JSON(http.StatusOK, PaginatedResponse{
		Success:    true,
		Data:       data,
		Pagination: pagination,
	})
}

// Fonctions helper pour les erreurs courantes

// BadRequestResponse envoie une erreur 400 Bad Request
func BadRequestResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusBadRequest, message, nil)
}

// UnauthorizedResponse envoie une erreur 401 Unauthorized
func UnauthorizedResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusUnauthorized, message, nil)
}

// ForbiddenResponse envoie une erreur 403 Forbidden
func ForbiddenResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusForbidden, message, nil)
}

// NotFoundResponse envoie une erreur 404 Not Found
func NotFoundResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusNotFound, message, nil)
}

// InternalServerErrorResponse envoie une erreur 500 Internal Server Error
func InternalServerErrorResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusInternalServerError, message, nil)
}
