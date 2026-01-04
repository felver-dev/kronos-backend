package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/services"
	"github.com/mcicare/itsm-backend/internal/utils"
)

// ServiceRequestTypeHandler gère les handlers des types de demandes de service
type ServiceRequestTypeHandler struct {
	serviceRequestTypeService services.ServiceRequestTypeService
}

// NewServiceRequestTypeHandler crée une nouvelle instance de ServiceRequestTypeHandler
func NewServiceRequestTypeHandler(serviceRequestTypeService services.ServiceRequestTypeService) *ServiceRequestTypeHandler {
	return &ServiceRequestTypeHandler{
		serviceRequestTypeService: serviceRequestTypeService,
	}
}

// GetAll récupère tous les types de demandes de service
// @Summary Récupérer tous les types de demandes
// @Description Récupère la liste de tous les types de demandes de service paramétrables
// @Tags service-requests
// @Security BearerAuth
// @Produce json
// @Success 200 {array} dto.ServiceRequestTypeDTO
// @Failure 500 {object} utils.Response
// @Router /service-requests/types [get]
func (h *ServiceRequestTypeHandler) GetAll(c *gin.Context) {
	types, err := h.serviceRequestTypeService.GetAll()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la récupération des types")
		return
	}

	utils.SuccessResponse(c, types, "Types récupérés avec succès")
}

// GetByID récupère un type par son ID
// @Summary Récupérer un type par ID
// @Description Récupère un type de demande de service par son identifiant
// @Tags service-requests
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID du type"
// @Success 200 {object} dto.ServiceRequestTypeDTO
// @Failure 404 {object} utils.Response
// @Router /service-requests/types/{id} [get]
func (h *ServiceRequestTypeHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	serviceRequestType, err := h.serviceRequestTypeService.GetByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Type introuvable")
		return
	}

	utils.SuccessResponse(c, serviceRequestType, "Type récupéré avec succès")
}

// Create crée un nouveau type de demande de service
// @Summary Créer un type de demande
// @Description Crée un nouveau type de demande de service paramétrable
// @Tags service-requests
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateServiceRequestTypeRequest true "Données du type"
// @Success 201 {object} dto.ServiceRequestTypeDTO
// @Failure 400 {object} utils.Response
// @Router /service-requests/types [post]
func (h *ServiceRequestTypeHandler) Create(c *gin.Context) {
	var req dto.CreateServiceRequestTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	createdByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	serviceRequestType, err := h.serviceRequestTypeService.Create(req, createdByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.CreatedResponse(c, serviceRequestType, "Type créé avec succès")
}

// Update met à jour un type de demande de service
// @Summary Mettre à jour un type
// @Description Met à jour les informations d'un type de demande de service
// @Tags service-requests
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID du type"
// @Param request body dto.UpdateServiceRequestTypeRequest true "Données à mettre à jour"
// @Success 200 {object} dto.ServiceRequestTypeDTO
// @Failure 400 {object} utils.Response
// @Router /service-requests/types/{id} [put]
func (h *ServiceRequestTypeHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	var req dto.UpdateServiceRequestTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	updatedByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	serviceRequestType, err := h.serviceRequestTypeService.Update(uint(id), req, updatedByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, serviceRequestType, "Type mis à jour avec succès")
}

// Delete supprime un type de demande de service
// @Summary Supprimer un type
// @Description Supprime un type de demande de service du système
// @Tags service-requests
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID du type"
// @Success 200 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /service-requests/types/{id} [delete]
func (h *ServiceRequestTypeHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID invalide")
		return
	}

	err = h.serviceRequestTypeService.Delete(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Type introuvable")
		return
	}

	utils.SuccessResponse(c, nil, "Type supprimé avec succès")
}

