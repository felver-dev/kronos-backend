package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/services"
	"github.com/mcicare/itsm-backend/internal/utils"
)

// ServiceRequestHandler gère les handlers des demandes de service
type ServiceRequestHandler struct {
	serviceRequestService services.ServiceRequestService
}

// NewServiceRequestHandler crée une nouvelle instance de ServiceRequestHandler
func NewServiceRequestHandler(serviceRequestService services.ServiceRequestService) *ServiceRequestHandler {
	return &ServiceRequestHandler{
		serviceRequestService: serviceRequestService,
	}
}

// Create crée une nouvelle demande de service
func (h *ServiceRequestHandler) Create(c *gin.Context) {
	var req dto.CreateServiceRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	createdByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	serviceRequest, err := h.serviceRequestService.Create(req, createdByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.CreatedResponse(c, serviceRequest, "Demande de service créée avec succès")
}

// GetByID récupère une demande de service par son ID
func (h *ServiceRequestHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	serviceRequest, err := h.serviceRequestService.GetByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Demande de service introuvable")
		return
	}

	utils.SuccessResponse(c, serviceRequest, "Demande de service récupérée avec succès")
}

// GetAll récupère toutes les demandes de service
func (h *ServiceRequestHandler) GetAll(c *gin.Context) {
	serviceRequests, err := h.serviceRequestService.GetAll()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des demandes de service")
		return
	}

	utils.SuccessResponse(c, serviceRequests, "Demandes de service récupérées avec succès")
}

// Validate valide une demande de service
func (h *ServiceRequestHandler) Validate(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	var req dto.ValidateServiceRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	validatedByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	serviceRequest, err := h.serviceRequestService.Validate(uint(id), req, validatedByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, serviceRequest, "Demande de service validée avec succès")
}

// Delete supprime une demande de service
func (h *ServiceRequestHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	err = h.serviceRequestService.Delete(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Demande de service introuvable")
		return
	}

	utils.SuccessResponse(c, nil, "Demande de service supprimée avec succès")
}

