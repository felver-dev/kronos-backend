package services

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/mcicare/itsm-backend/config"
	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/models"
	"github.com/mcicare/itsm-backend/internal/repositories"
)

// TicketAttachmentService interface pour les opérations sur les pièces jointes de tickets
type TicketAttachmentService interface {
	UploadAttachment(ticketID uint, fileName, filePath, thumbnailPath string, fileSize int, mimeType string, isImage bool, description string, displayOrder int, userID uint) (*dto.TicketAttachmentDTO, error)
	GetByTicketID(ticketID uint, imagesOnly bool) ([]dto.TicketAttachmentDTO, error)
	GetByID(id uint) (*dto.TicketAttachmentDTO, error)
	GetImagesByTicketID(ticketID uint) ([]dto.TicketAttachmentDTO, error)
	GetFilePath(id uint) (string, error)
	GetThumbnailPath(id uint) (string, error)
	Update(id uint, req dto.UpdateTicketAttachmentRequest, updatedByID uint) (*dto.TicketAttachmentDTO, error)
	SetPrimary(ticketID, attachmentID uint, updatedByID uint) (*dto.TicketAttachmentDTO, error)
	Delete(id uint) error
	Reorder(ticketID uint, attachmentIDs []uint, updatedByID uint) error
}

// ticketAttachmentService implémente TicketAttachmentService
type ticketAttachmentService struct {
	attachmentRepo repositories.TicketAttachmentRepository
	ticketRepo     repositories.TicketRepository
	userRepo       repositories.UserRepository
}

// NewTicketAttachmentService crée une nouvelle instance de TicketAttachmentService
func NewTicketAttachmentService(
	attachmentRepo repositories.TicketAttachmentRepository,
	ticketRepo repositories.TicketRepository,
	userRepo repositories.UserRepository,
) TicketAttachmentService {
	return &ticketAttachmentService{
		attachmentRepo: attachmentRepo,
		ticketRepo:     ticketRepo,
		userRepo:       userRepo,
	}
}

// UploadAttachment upload une pièce jointe pour un ticket
func (s *ticketAttachmentService) UploadAttachment(ticketID uint, fileName, filePath, thumbnailPath string, fileSize int, mimeType string, isImage bool, description string, displayOrder int, userID uint) (*dto.TicketAttachmentDTO, error) {
	// Vérifier que le ticket existe
	_, err := s.ticketRepo.FindByID(ticketID)
	if err != nil {
		return nil, errors.New("ticket introuvable")
	}

	// Vérifier que l'utilisateur existe
	_, err = s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("utilisateur introuvable")
	}

	// Créer l'attachment
	attachment := &models.TicketAttachment{
		TicketID:      ticketID,
		UserID:        userID,
		FileName:      fileName,
		FilePath:      filePath,
		ThumbnailPath: thumbnailPath,
		FileSize:      &fileSize,
		MimeType:      mimeType,
		IsImage:       isImage,
		DisplayOrder:  displayOrder,
		Description:   description,
	}

	if err := s.attachmentRepo.Create(attachment); err != nil {
		return nil, errors.New("erreur lors de la création de la pièce jointe")
	}

	// Récupérer l'attachment créé
	createdAttachment, err := s.attachmentRepo.FindByID(attachment.ID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération de la pièce jointe créée")
	}

	attachmentDTO := s.attachmentToDTO(createdAttachment)
	return &attachmentDTO, nil
}

// GetByTicketID récupère toutes les pièces jointes d'un ticket
func (s *ticketAttachmentService) GetByTicketID(ticketID uint, imagesOnly bool) ([]dto.TicketAttachmentDTO, error) {
	// Vérifier que le ticket existe
	_, err := s.ticketRepo.FindByID(ticketID)
	if err != nil {
		return nil, errors.New("ticket introuvable")
	}

	var attachments []models.TicketAttachment
	if imagesOnly {
		attachments, err = s.attachmentRepo.FindImagesByTicketID(ticketID)
	} else {
		attachments, err = s.attachmentRepo.FindByTicketID(ticketID)
	}

	if err != nil {
		return nil, errors.New("erreur lors de la récupération des pièces jointes")
	}

	attachmentDTOs := make([]dto.TicketAttachmentDTO, len(attachments))
	for i, attachment := range attachments {
		attachmentDTOs[i] = s.attachmentToDTO(&attachment)
	}

	return attachmentDTOs, nil
}

// GetByID récupère une pièce jointe par son ID
func (s *ticketAttachmentService) GetByID(id uint) (*dto.TicketAttachmentDTO, error) {
	attachment, err := s.attachmentRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("pièce jointe introuvable")
	}

	attachmentDTO := s.attachmentToDTO(attachment)
	return &attachmentDTO, nil
}

// GetImagesByTicketID récupère toutes les images d'un ticket
func (s *ticketAttachmentService) GetImagesByTicketID(ticketID uint) ([]dto.TicketAttachmentDTO, error) {
	return s.GetByTicketID(ticketID, true)
}

// GetFilePath récupère le chemin complet d'une pièce jointe
func (s *ticketAttachmentService) GetFilePath(id uint) (string, error) {
	attachment, err := s.attachmentRepo.FindByID(id)
	if err != nil {
		return "", errors.New("pièce jointe introuvable")
	}

	fullPath := filepath.Join(config.AppConfig.TicketAttachmentsDir, attachment.FilePath)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", errors.New("fichier introuvable")
	}

	return fullPath, nil
}

// GetThumbnailPath récupère le chemin de la miniature
func (s *ticketAttachmentService) GetThumbnailPath(id uint) (string, error) {
	attachment, err := s.attachmentRepo.FindByID(id)
	if err != nil {
		return "", errors.New("pièce jointe introuvable")
	}

	if !attachment.IsImage {
		return "", errors.New("cette pièce jointe n'est pas une image")
	}

	if attachment.ThumbnailPath == "" {
		// Si pas de miniature, retourner l'image originale
		return s.GetFilePath(id)
	}

	fullPath := filepath.Join(config.AppConfig.TicketAttachmentsDir, attachment.ThumbnailPath)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		// Si la miniature n'existe pas, retourner l'image originale
		return s.GetFilePath(id)
	}

	return fullPath, nil
}

// Update met à jour une pièce jointe
func (s *ticketAttachmentService) Update(id uint, req dto.UpdateTicketAttachmentRequest, updatedByID uint) (*dto.TicketAttachmentDTO, error) {
	attachment, err := s.attachmentRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("pièce jointe introuvable")
	}

	if req.Description != "" {
		attachment.Description = req.Description
	}

	if req.DisplayOrder != 0 {
		attachment.DisplayOrder = req.DisplayOrder
	}

	if err := s.attachmentRepo.Update(attachment); err != nil {
		return nil, errors.New("erreur lors de la mise à jour de la pièce jointe")
	}

	updatedAttachment, err := s.attachmentRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération de la pièce jointe mise à jour")
	}

	attachmentDTO := s.attachmentToDTO(updatedAttachment)
	return &attachmentDTO, nil
}

// SetPrimary définit une image comme principale (déprécié, on utilise display_order)
func (s *ticketAttachmentService) SetPrimary(ticketID, attachmentID uint, updatedByID uint) (*dto.TicketAttachmentDTO, error) {
	// Vérifier que le ticket existe
	_, err := s.ticketRepo.FindByID(ticketID)
	if err != nil {
		return nil, errors.New("ticket introuvable")
	}

	// Vérifier que l'attachment existe et appartient au ticket
	attachment, err := s.attachmentRepo.FindByID(attachmentID)
	if err != nil {
		return nil, errors.New("pièce jointe introuvable")
	}

	if attachment.TicketID != ticketID {
		return nil, errors.New("cette pièce jointe n'appartient pas à ce ticket")
	}

	if !attachment.IsImage {
		return nil, errors.New("seules les images peuvent être définies comme principales")
	}

	// Mettre l'ordre d'affichage à 0 pour la rendre principale
	attachment.DisplayOrder = 0

	// Réorganiser les autres images
	attachments, _ := s.attachmentRepo.FindImagesByTicketID(ticketID)
	for _, att := range attachments {
		if att.ID != attachmentID && att.DisplayOrder == 0 {
			att.DisplayOrder = 1
			s.attachmentRepo.Update(&att)
		}
	}

	if err := s.attachmentRepo.Update(attachment); err != nil {
		return nil, errors.New("erreur lors de la mise à jour de la pièce jointe")
	}

	updatedAttachment, err := s.attachmentRepo.FindByID(attachmentID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération de la pièce jointe mise à jour")
	}

	attachmentDTO := s.attachmentToDTO(updatedAttachment)
	return &attachmentDTO, nil
}

// Delete supprime une pièce jointe
func (s *ticketAttachmentService) Delete(id uint) error {
	attachment, err := s.attachmentRepo.FindByID(id)
	if err != nil {
		return errors.New("pièce jointe introuvable")
	}

	// Supprimer le fichier
	filePath := filepath.Join(config.AppConfig.TicketAttachmentsDir, attachment.FilePath)
	if _, err := os.Stat(filePath); err == nil {
		os.Remove(filePath)
	}

	// Supprimer la miniature si elle existe
	if attachment.ThumbnailPath != "" {
		thumbnailPath := filepath.Join(config.AppConfig.TicketAttachmentsDir, attachment.ThumbnailPath)
		if _, err := os.Stat(thumbnailPath); err == nil {
			os.Remove(thumbnailPath)
		}
	}

	// Supprimer de la base de données
	if err := s.attachmentRepo.Delete(id); err != nil {
		return errors.New("erreur lors de la suppression de la pièce jointe")
	}

	return nil
}

// Reorder réorganise les pièces jointes d'un ticket
func (s *ticketAttachmentService) Reorder(ticketID uint, attachmentIDs []uint, updatedByID uint) error {
	// Vérifier que le ticket existe
	_, err := s.ticketRepo.FindByID(ticketID)
	if err != nil {
		return errors.New("ticket introuvable")
	}

	// Mettre à jour l'ordre d'affichage pour chaque attachment
	for order, attachmentID := range attachmentIDs {
		attachment, err := s.attachmentRepo.FindByID(attachmentID)
		if err != nil {
			continue // Ignorer les IDs invalides
		}

		if attachment.TicketID != ticketID {
			continue // Ignorer les attachments qui n'appartiennent pas au ticket
		}

		attachment.DisplayOrder = order
		s.attachmentRepo.Update(attachment)
	}

	return nil
}

// attachmentToDTO convertit un modèle TicketAttachment en DTO
func (s *ticketAttachmentService) attachmentToDTO(attachment *models.TicketAttachment) dto.TicketAttachmentDTO {
	userDTO := dto.UserDTO{
		ID:        attachment.User.ID,
		Username:  attachment.User.Username,
		Email:     attachment.User.Email,
		FirstName: attachment.User.FirstName,
		LastName:  attachment.User.LastName,
		Avatar:    attachment.User.Avatar,
		Role:      attachment.User.Role.Name,
		IsActive:  attachment.User.IsActive,
		CreatedAt: attachment.User.CreatedAt,
		UpdatedAt: attachment.User.UpdatedAt,
	}

	return dto.TicketAttachmentDTO{
		ID:            attachment.ID,
		TicketID:      attachment.TicketID,
		User:          userDTO,
		FileName:      attachment.FileName,
		FilePath:      attachment.FilePath,
		ThumbnailPath: attachment.ThumbnailPath,
		FileSize:      attachment.FileSize,
		MimeType:      attachment.MimeType,
		IsImage:       attachment.IsImage,
		DisplayOrder:  attachment.DisplayOrder,
		Description:   attachment.Description,
		CreatedAt:     attachment.CreatedAt,
	}
}

