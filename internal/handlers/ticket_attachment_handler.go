package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/config"
	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/services"
	"github.com/mcicare/itsm-backend/internal/utils"
)

// TicketAttachmentHandler gère les handlers des pièces jointes de tickets
type TicketAttachmentHandler struct {
	attachmentService services.TicketAttachmentService
}

// NewTicketAttachmentHandler crée une nouvelle instance de TicketAttachmentHandler
func NewTicketAttachmentHandler(attachmentService services.TicketAttachmentService) *TicketAttachmentHandler {
	return &TicketAttachmentHandler{
		attachmentService: attachmentService,
	}
}

// UploadAttachment upload une pièce jointe pour un ticket
// @Summary Uploader une pièce jointe
// @Description Upload une pièce jointe (image ou document) pour un ticket
// @Tags tickets
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "ID du ticket"
// @Param file formData file true "Fichier à uploader"
// @Param description formData string false "Description de la pièce jointe"
// @Param display_order formData int false "Ordre d'affichage"
// @Success 201 {object} dto.TicketAttachmentDTO
// @Failure 400 {object} utils.Response
// @Router /tickets/{id}/attachments [post]
func (h *TicketAttachmentHandler) UploadAttachment(c *gin.Context) {
	ticketIDParam := c.Param("id")
	ticketID, err := strconv.ParseUint(ticketIDParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID de ticket invalide")
		return
	}

	// Récupérer le fichier
	file, err := c.FormFile("file")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Fichier manquant", err.Error())
		return
	}

	// Vérifier la taille
	if file.Size > config.AppConfig.MaxUploadSize {
		utils.ErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Fichier trop volumineux. Taille maximale: %d bytes", config.AppConfig.MaxUploadSize), nil)
		return
	}

	// Vérifier le type de fichier
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp", ".pdf", ".doc", ".docx", ".xls", ".xlsx", ".txt", ".zip"}
	isAllowed := false
	for _, allowedExt := range allowedExts {
		if ext == allowedExt {
			isAllowed = true
			break
		}
	}
	if !isAllowed {
		utils.ErrorResponse(c, http.StatusBadRequest, "Type de fichier non autorisé", nil)
		return
	}

	// Déterminer si c'est une image
	isImage := false
	imageExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	for _, imgExt := range imageExts {
		if ext == imgExt {
			isImage = true
			break
		}
	}

	// Créer le dossier de destination s'il n'existe pas
	attachmentsDir := config.AppConfig.TicketAttachmentsDir
	ticketDir := filepath.Join(attachmentsDir, fmt.Sprintf("ticket_%d", ticketID))
	if err := os.MkdirAll(ticketDir, 0755); err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la création du dossier de destination")
		return
	}

	// Générer un nom de fichier unique
	timestamp := time.Now().Unix()
	fileName := fmt.Sprintf("%d_%s", timestamp, file.Filename)
	filePath := filepath.Join(ticketDir, fileName)

	// Sauvegarder le fichier
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		utils.InternalServerErrorResponse(c, "Erreur lors de la sauvegarde du fichier")
		return
	}

	// Générer le chemin relatif
	relativePath := filepath.Join(fmt.Sprintf("ticket_%d", ticketID), fileName)

	// Générer une miniature si c'est une image (TODO: implémenter la génération de miniature)
	thumbnailPath := ""
	if isImage {
		// Pour l'instant, on utilise le fichier original comme miniature
		thumbnailPath = relativePath
	}

	// Récupérer les paramètres optionnels
	description := c.PostForm("description")
	displayOrderStr := c.PostForm("display_order")
	displayOrder := 0
	if displayOrderStr != "" {
		if order, err := strconv.Atoi(displayOrderStr); err == nil {
			displayOrder = order
		}
	}

	// Récupérer l'utilisateur
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	// Détecter le type MIME
	mimeType := file.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	// Créer l'attachment
	attachment, err := h.attachmentService.UploadAttachment(
		uint(ticketID),
		file.Filename,
		relativePath,
		thumbnailPath,
		int(file.Size),
		mimeType,
		isImage,
		description,
		displayOrder,
		userID.(uint),
	)
	if err != nil {
		// Supprimer le fichier en cas d'erreur
		os.Remove(filePath)
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.CreatedResponse(c, attachment, "Pièce jointe uploadée avec succès")
}

// GetAttachments récupère toutes les pièces jointes d'un ticket
// @Summary Récupérer les pièces jointes d'un ticket
// @Description Récupère la liste de toutes les pièces jointes d'un ticket
// @Tags tickets
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID du ticket"
// @Param images_only query bool false "Récupérer uniquement les images"
// @Success 200 {array} dto.TicketAttachmentDTO
// @Failure 400 {object} utils.Response
// @Router /tickets/{id}/attachments [get]
func (h *TicketAttachmentHandler) GetAttachments(c *gin.Context) {
	ticketIDParam := c.Param("id")
	ticketID, err := strconv.ParseUint(ticketIDParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID de ticket invalide")
		return
	}

	imagesOnly := c.Query("images_only") == "true"

	attachments, err := h.attachmentService.GetByTicketID(uint(ticketID), imagesOnly)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, attachments, "Pièces jointes récupérées avec succès")
}

// GetImages récupère toutes les images d'un ticket
// @Summary Récupérer les images d'un ticket
// @Description Récupère la liste de toutes les images d'un ticket
// @Tags tickets
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID du ticket"
// @Success 200 {array} dto.TicketAttachmentDTO
// @Failure 400 {object} utils.Response
// @Router /tickets/{id}/attachments/images [get]
func (h *TicketAttachmentHandler) GetImages(c *gin.Context) {
	ticketIDParam := c.Param("id")
	ticketID, err := strconv.ParseUint(ticketIDParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID de ticket invalide")
		return
	}

	images, err := h.attachmentService.GetImagesByTicketID(uint(ticketID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, images, "Images récupérées avec succès")
}

// GetByID récupère une pièce jointe par son ID
// @Summary Récupérer une pièce jointe
// @Description Récupère les informations d'une pièce jointe par son ID
// @Tags tickets
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID du ticket"
// @Param attachmentId path int true "ID de la pièce jointe"
// @Success 200 {object} dto.TicketAttachmentDTO
// @Failure 404 {object} utils.Response
// @Router /tickets/{id}/attachments/{attachmentId} [get]
func (h *TicketAttachmentHandler) GetByID(c *gin.Context) {
	attachmentIDParam := c.Param("attachmentId")
	attachmentID, err := strconv.ParseUint(attachmentIDParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID de pièce jointe invalide")
		return
	}

	attachment, err := h.attachmentService.GetByID(uint(attachmentID))
	if err != nil {
		utils.NotFoundResponse(c, "Pièce jointe introuvable")
		return
	}

	utils.SuccessResponse(c, attachment, "Pièce jointe récupérée avec succès")
}

// Download télécharge une pièce jointe
// @Summary Télécharger une pièce jointe
// @Description Télécharge une pièce jointe
// @Tags tickets
// @Security BearerAuth
// @Produce application/octet-stream
// @Param id path int true "ID du ticket"
// @Param attachmentId path int true "ID de la pièce jointe"
// @Success 200 {file} file "Fichier téléchargeable"
// @Failure 404 {object} utils.Response
// @Router /tickets/{id}/attachments/{attachmentId}/download [get]
func (h *TicketAttachmentHandler) Download(c *gin.Context) {
	attachmentIDParam := c.Param("attachmentId")
	attachmentID, err := strconv.ParseUint(attachmentIDParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID de pièce jointe invalide")
		return
	}

	filePath, err := h.attachmentService.GetFilePath(uint(attachmentID))
	if err != nil {
		utils.NotFoundResponse(c, err.Error())
		return
	}

	c.File(filePath)
}

// GetThumbnail récupère la miniature d'une image
// @Summary Récupérer la miniature d'une image
// @Description Récupère la miniature d'une image
// @Tags tickets
// @Security BearerAuth
// @Produce image/*
// @Param id path int true "ID du ticket"
// @Param attachmentId path int true "ID de la pièce jointe"
// @Success 200 {file} file "Miniature de l'image"
// @Failure 404 {object} utils.Response
// @Router /tickets/{id}/attachments/{attachmentId}/thumbnail [get]
func (h *TicketAttachmentHandler) GetThumbnail(c *gin.Context) {
	attachmentIDParam := c.Param("attachmentId")
	attachmentID, err := strconv.ParseUint(attachmentIDParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID de pièce jointe invalide")
		return
	}

	thumbnailPath, err := h.attachmentService.GetThumbnailPath(uint(attachmentID))
	if err != nil {
		utils.NotFoundResponse(c, err.Error())
		return
	}

	c.File(thumbnailPath)
}

// Update met à jour une pièce jointe
// @Summary Mettre à jour une pièce jointe
// @Description Met à jour les informations d'une pièce jointe
// @Tags tickets
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID du ticket"
// @Param attachmentId path int true "ID de la pièce jointe"
// @Param request body dto.UpdateTicketAttachmentRequest true "Données à mettre à jour"
// @Success 200 {object} dto.TicketAttachmentDTO
// @Failure 400 {object} utils.Response
// @Router /tickets/{id}/attachments/{attachmentId} [put]
func (h *TicketAttachmentHandler) Update(c *gin.Context) {
	attachmentIDParam := c.Param("attachmentId")
	attachmentID, err := strconv.ParseUint(attachmentIDParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID de pièce jointe invalide")
		return
	}

	var req dto.UpdateTicketAttachmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	updatedByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	attachment, err := h.attachmentService.Update(uint(attachmentID), req, updatedByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, attachment, "Pièce jointe mise à jour avec succès")
}

// SetPrimary définit une image comme principale
// @Summary Définir une image comme principale
// @Description Définit une image comme image principale du ticket
// @Tags tickets
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID du ticket"
// @Param attachmentId path int true "ID de la pièce jointe"
// @Success 200 {object} dto.TicketAttachmentDTO
// @Failure 400 {object} utils.Response
// @Router /tickets/{id}/attachments/{attachmentId}/set-primary [put]
func (h *TicketAttachmentHandler) SetPrimary(c *gin.Context) {
	ticketIDParam := c.Param("id")
	ticketID, err := strconv.ParseUint(ticketIDParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID de ticket invalide")
		return
	}

	attachmentIDParam := c.Param("attachmentId")
	attachmentID, err := strconv.ParseUint(attachmentIDParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID de pièce jointe invalide")
		return
	}

	updatedByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	attachment, err := h.attachmentService.SetPrimary(uint(ticketID), uint(attachmentID), updatedByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, attachment, "Image principale définie avec succès")
}

// Delete supprime une pièce jointe
// @Summary Supprimer une pièce jointe
// @Description Supprime une pièce jointe d'un ticket
// @Tags tickets
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID du ticket"
// @Param attachmentId path int true "ID de la pièce jointe"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /tickets/{id}/attachments/{attachmentId} [delete]
func (h *TicketAttachmentHandler) Delete(c *gin.Context) {
	attachmentIDParam := c.Param("attachmentId")
	attachmentID, err := strconv.ParseUint(attachmentIDParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID de pièce jointe invalide")
		return
	}

	err = h.attachmentService.Delete(uint(attachmentID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, nil, "Pièce jointe supprimée avec succès")
}

// Reorder réorganise les pièces jointes d'un ticket
// @Summary Réorganiser les pièces jointes
// @Description Réorganise l'ordre d'affichage des pièces jointes d'un ticket
// @Tags tickets
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID du ticket"
// @Param request body dto.ReorderTicketAttachmentsRequest true "Liste des IDs dans le nouvel ordre"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /tickets/{id}/attachments/reorder [put]
func (h *TicketAttachmentHandler) Reorder(c *gin.Context) {
	ticketIDParam := c.Param("id")
	ticketID, err := strconv.ParseUint(ticketIDParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "ID de ticket invalide")
		return
	}

	var req dto.ReorderTicketAttachmentsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Données invalides", err.Error())
		return
	}

	updatedByID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Utilisateur non authentifié")
		return
	}

	err = h.attachmentService.Reorder(uint(ticketID), req.AttachmentIDs, updatedByID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, nil, "Ordre des pièces jointes mis à jour avec succès")
}

