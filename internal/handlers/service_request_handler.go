package handlers

import (
	"net/http"
	"strconv"
	"time"

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
// @Summary Créer une demande de service
// @Description Crée une nouvelle demande de service dans le système
// @Tags service-requests
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateServiceRequestRequest true "Données de la demande de service"
// @Success 201 {object} dto.ServiceRequestDTO
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /service-requests [post]
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
// @Summary Récupérer une demande de service par ID
// @Description Récupère une demande de service par son identifiant
// @Tags service-requests
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID de la demande de service"
// @Success 200 {object} dto.ServiceRequestDTO
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /service-requests/{id} [get]
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
// @Summary Récupérer toutes les demandes de service
// @Description Récupère la liste de toutes les demandes de service
// @Tags service-requests
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {array} dto.ServiceRequestDTO
// @Failure 500 {object} utils.Response
// @Router /service-requests [get]
func (h *ServiceRequestHandler) GetAll(c *gin.Context) {
	serviceRequests, err := h.serviceRequestService.GetAll()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des demandes de service")
		return
	}

	utils.SuccessResponse(c, serviceRequests, "Demandes de service récupérées avec succès")
}

// Validate valide une demande de service
// @Summary Valider une demande de service
// @Description Valide une demande de service
// @Tags service-requests
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID de la demande de service"
// @Param request body dto.ValidateServiceRequestRequest true "Données de validation"
// @Success 200 {object} dto.ServiceRequestDTO
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /service-requests/{id}/validate [post]
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
// @Summary Supprimer une demande de service
// @Description Supprime une demande de service du système
// @Tags service-requests
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID de la demande de service"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /service-requests/{id} [delete]
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

// Update met à jour une demande de service
// @Summary Mettre à jour une demande de service
// @Description Met à jour les informations d'une demande de service
// @Tags service-requests
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID de la demande de service"
// @Param request body dto.UpdateServiceRequestRequest true "Données à mettre à jour"
// @Success 200 {object} dto.ServiceRequestDTO
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /service-requests/{id} [put]
func (h *ServiceRequestHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	var req dto.UpdateServiceRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	updatedByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	serviceRequest, err := h.serviceRequestService.Update(uint(id), req, updatedByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, serviceRequest, "Demande de service mise à jour avec succès")
}

// GetDeadline récupère le délai de traitement d'une demande de service
// @Summary Récupérer le délai de traitement
// @Description Récupère le délai de traitement d'une demande de service
// @Tags service-requests
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID de la demande de service"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} utils.Response
// @Router /service-requests/{id}/deadline [get]
func (h *ServiceRequestHandler) GetDeadline(c *gin.Context) {
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

	response := map[string]interface{}{
		"deadline":  serviceRequest.Deadline,
		"remaining": nil,
		"unit":      "days",
	}

	if serviceRequest.Deadline != nil {
		now := time.Now()
		remaining := serviceRequest.Deadline.Sub(now)
		remainingDays := int(remaining.Hours() / 24)
		response["remaining"] = remainingDays
		if remainingDays < 0 {
			response["status"] = "overdue"
		} else if remainingDays <= 1 {
			response["status"] = "urgent"
		} else {
			response["status"] = "on_time"
		}
	}

	utils.SuccessResponse(c, response, "Délai récupéré avec succès")
}

// GetValidationStatus récupère le statut de validation d'une demande de service
// @Summary Récupérer le statut de validation
// @Description Récupère le statut de validation d'une demande de service
// @Tags service-requests
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID de la demande de service"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} utils.Response
// @Router /service-requests/{id}/validation-status [get]
func (h *ServiceRequestHandler) GetValidationStatus(c *gin.Context) {
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

	response := map[string]interface{}{
		"validated": serviceRequest.Validated,
		"validator": serviceRequest.ValidatedBy,
		"date":      serviceRequest.ValidatedAt,
		"comment":   serviceRequest.ValidationComment,
	}

	utils.SuccessResponse(c, response, "Statut de validation récupéré avec succès")
}
