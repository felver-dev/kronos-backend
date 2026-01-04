package services

import (
	"errors"
	"time"

	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/models"
	"github.com/mcicare/itsm-backend/internal/repositories"
)

// AssetService interface pour les opérations sur les actifs IT
type AssetService interface {
	Create(req dto.CreateAssetRequest, createdByID uint) (*dto.AssetDTO, error)
	GetByID(id uint) (*dto.AssetDTO, error)
	GetAll() ([]dto.AssetDTO, error)
	GetByCategory(categoryID uint) ([]dto.AssetDTO, error)
	GetByStatus(status string) ([]dto.AssetDTO, error)
	GetByAssignedTo(userID uint) ([]dto.AssetDTO, error)
	Update(id uint, req dto.UpdateAssetRequest, updatedByID uint) (*dto.AssetDTO, error)
	Assign(id uint, req dto.AssignAssetRequest, assignedByID uint) (*dto.AssetDTO, error)
	Unassign(id uint, req dto.AssignAssetRequest, unassignedByID uint) (*dto.AssetDTO, error)
	GetInventory() (*dto.AssetInventoryDTO, error)
	GetLinkedTickets(assetID uint) ([]dto.TicketDTO, error)
	LinkTicket(assetID uint, ticketID uint, linkedByID uint) error
	UnlinkTicket(assetID uint, ticketID uint) error
	Delete(id uint) error
}

// assetService implémente AssetService
type assetService struct {
	assetRepo         repositories.AssetRepository
	assetCategoryRepo repositories.AssetCategoryRepository
	userRepo          repositories.UserRepository
	ticketAssetRepo   repositories.TicketAssetRepository
	ticketRepo        repositories.TicketRepository
}

// NewAssetService crée une nouvelle instance de AssetService
func NewAssetService(
	assetRepo repositories.AssetRepository,
	assetCategoryRepo repositories.AssetCategoryRepository,
	userRepo repositories.UserRepository,
	ticketAssetRepo repositories.TicketAssetRepository,
	ticketRepo repositories.TicketRepository,
) AssetService {
	return &assetService{
		assetRepo:         assetRepo,
		assetCategoryRepo: assetCategoryRepo,
		userRepo:          userRepo,
		ticketAssetRepo:   ticketAssetRepo,
		ticketRepo:        ticketRepo,
	}
}

// Create crée un nouvel actif
func (s *assetService) Create(req dto.CreateAssetRequest, createdByID uint) (*dto.AssetDTO, error) {
	// Vérifier que la catégorie existe
	_, err := s.assetCategoryRepo.FindByID(req.CategoryID)
	if err != nil {
		return nil, errors.New("catégorie introuvable")
	}

	// Vérifier que l'utilisateur assigné existe si fourni
	if req.AssignedTo != nil {
		_, err = s.userRepo.FindByID(*req.AssignedTo)
		if err != nil {
			return nil, errors.New("utilisateur assigné introuvable")
		}
	}

	// Parser les dates si fournies
	var purchaseDate, warrantyExpiry *time.Time
	if req.PurchaseDate != nil && *req.PurchaseDate != "" {
		parsed, err := time.Parse("2006-01-02", *req.PurchaseDate)
		if err == nil {
			purchaseDate = &parsed
		}
	}
	if req.WarrantyExpiry != nil && *req.WarrantyExpiry != "" {
		parsed, err := time.Parse("2006-01-02", *req.WarrantyExpiry)
		if err == nil {
			warrantyExpiry = &parsed
		}
	}

	// Définir le statut par défaut
	status := req.Status
	if status == "" {
		status = "available"
	}

	// Créer l'actif
	asset := &models.Asset{
		Name:           req.Name,
		SerialNumber:   req.SerialNumber,
		Model:          req.Model,
		Manufacturer:   req.Manufacturer,
		CategoryID:     req.CategoryID,
		AssignedToID:   req.AssignedTo,
		Status:         status,
		PurchaseDate:   purchaseDate,
		WarrantyExpiry: warrantyExpiry,
		Location:       req.Location,
		Notes:          req.Notes,
		CreatedByID:    &createdByID,
	}

	if err := s.assetRepo.Create(asset); err != nil {
		return nil, errors.New("erreur lors de la création de l'actif")
	}

	// Récupérer l'actif créé avec ses relations
	createdAsset, err := s.assetRepo.FindByID(asset.ID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération de l'actif créé")
	}

	assetDTO := s.assetToDTO(createdAsset)
	return &assetDTO, nil
}

// GetByID récupère un actif par son ID
func (s *assetService) GetByID(id uint) (*dto.AssetDTO, error) {
	asset, err := s.assetRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("actif introuvable")
	}

	assetDTO := s.assetToDTO(asset)
	return &assetDTO, nil
}

// GetAll récupère tous les actifs
func (s *assetService) GetAll() ([]dto.AssetDTO, error) {
	assets, err := s.assetRepo.FindAll()
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des actifs")
	}

	var assetDTOs []dto.AssetDTO
	for _, asset := range assets {
		assetDTOs = append(assetDTOs, s.assetToDTO(&asset))
	}

	return assetDTOs, nil
}

// GetByCategory récupère les actifs d'une catégorie
func (s *assetService) GetByCategory(categoryID uint) ([]dto.AssetDTO, error) {
	assets, err := s.assetRepo.FindByCategory(categoryID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des actifs")
	}

	var assetDTOs []dto.AssetDTO
	for _, asset := range assets {
		assetDTOs = append(assetDTOs, s.assetToDTO(&asset))
	}

	return assetDTOs, nil
}

// GetByStatus récupère les actifs par statut
func (s *assetService) GetByStatus(status string) ([]dto.AssetDTO, error) {
	assets, err := s.assetRepo.FindByStatus(status)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des actifs")
	}

	var assetDTOs []dto.AssetDTO
	for _, asset := range assets {
		assetDTOs = append(assetDTOs, s.assetToDTO(&asset))
	}

	return assetDTOs, nil
}

// GetByAssignedTo récupère les actifs assignés à un utilisateur
func (s *assetService) GetByAssignedTo(userID uint) ([]dto.AssetDTO, error) {
	assets, err := s.assetRepo.FindByAssignedTo(userID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des actifs")
	}

	var assetDTOs []dto.AssetDTO
	for _, asset := range assets {
		assetDTOs = append(assetDTOs, s.assetToDTO(&asset))
	}

	return assetDTOs, nil
}

// Update met à jour un actif
func (s *assetService) Update(id uint, req dto.UpdateAssetRequest, updatedByID uint) (*dto.AssetDTO, error) {
	asset, err := s.assetRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("actif introuvable")
	}

	// Mettre à jour les champs fournis
	if req.Name != "" {
		asset.Name = req.Name
	}
	if req.SerialNumber != "" {
		asset.SerialNumber = req.SerialNumber
	}
	if req.Model != "" {
		asset.Model = req.Model
	}
	if req.Manufacturer != "" {
		asset.Manufacturer = req.Manufacturer
	}
	if req.CategoryID != nil {
		// Vérifier que la catégorie existe
		_, err = s.assetCategoryRepo.FindByID(*req.CategoryID)
		if err != nil {
			return nil, errors.New("catégorie introuvable")
		}
		asset.CategoryID = *req.CategoryID
	}
	if req.AssignedTo != nil {
		// Vérifier que l'utilisateur existe si assigné
		if *req.AssignedTo != 0 {
			_, err = s.userRepo.FindByID(*req.AssignedTo)
			if err != nil {
				return nil, errors.New("utilisateur assigné introuvable")
			}
			asset.AssignedToID = req.AssignedTo
		} else {
			asset.AssignedToID = nil
		}
	}
	if req.Status != "" {
		asset.Status = req.Status
	}
	if req.PurchaseDate != nil && *req.PurchaseDate != "" {
		parsed, err := time.Parse("2006-01-02", *req.PurchaseDate)
		if err == nil {
			asset.PurchaseDate = &parsed
		}
	}
	if req.WarrantyExpiry != nil && *req.WarrantyExpiry != "" {
		parsed, err := time.Parse("2006-01-02", *req.WarrantyExpiry)
		if err == nil {
			asset.WarrantyExpiry = &parsed
		}
	}
	if req.Location != "" {
		asset.Location = req.Location
	}
	if req.Notes != "" {
		asset.Notes = req.Notes
	}

	if err := s.assetRepo.Update(asset); err != nil {
		return nil, errors.New("erreur lors de la mise à jour de l'actif")
	}

	// Récupérer l'actif mis à jour
	updatedAsset, err := s.assetRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération de l'actif mis à jour")
	}

	assetDTO := s.assetToDTO(updatedAsset)
	return &assetDTO, nil
}

// Assign assigne un actif à un utilisateur
func (s *assetService) Assign(id uint, req dto.AssignAssetRequest, assignedByID uint) (*dto.AssetDTO, error) {
	asset, err := s.assetRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("actif introuvable")
	}

	// Vérifier que l'utilisateur existe
	_, err = s.userRepo.FindByID(req.UserID)
	if err != nil {
		return nil, errors.New("utilisateur introuvable")
	}

	asset.AssignedToID = &req.UserID
	asset.Status = "in_use"

	if err := s.assetRepo.Update(asset); err != nil {
		return nil, errors.New("erreur lors de l'assignation de l'actif")
	}

	// Récupérer l'actif mis à jour
	updatedAsset, err := s.assetRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération de l'actif mis à jour")
	}

	assetDTO := s.assetToDTO(updatedAsset)
	return &assetDTO, nil
}

// Unassign retire l'assignation d'un actif
func (s *assetService) Unassign(id uint, req dto.AssignAssetRequest, unassignedByID uint) (*dto.AssetDTO, error) {
	asset, err := s.assetRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("actif introuvable")
	}

	asset.AssignedToID = nil
	asset.Status = "available"

	if err := s.assetRepo.Update(asset); err != nil {
		return nil, errors.New("erreur lors de la désassignation de l'actif")
	}

	updatedAsset, err := s.assetRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération de l'actif mis à jour")
	}

	assetDTO := s.assetToDTO(updatedAsset)
	return &assetDTO, nil
}

// GetInventory récupère l'inventaire des actifs
func (s *assetService) GetInventory() (*dto.AssetInventoryDTO, error) {
	allAssets, err := s.assetRepo.FindAll()
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des actifs")
	}

	total := len(allAssets)
	byStatus := make(map[string]int)
	byCategory := make(map[string]int)
	assigned := 0
	available := 0

	for _, asset := range allAssets {
		byStatus[asset.Status]++
		if asset.Category.ID != 0 {
			byCategory[asset.Category.Name]++
		}
		if asset.AssignedToID != nil {
			assigned++
		} else {
			available++
		}
	}

	return &dto.AssetInventoryDTO{
		Total:      total,
		ByStatus:   byStatus,
		ByCategory: byCategory,
		Assigned:   assigned,
		Available:  available,
	}, nil
}

// GetLinkedTickets récupère les tickets liés à un actif
func (s *assetService) GetLinkedTickets(assetID uint) ([]dto.TicketDTO, error) {
	// Vérifier que l'actif existe
	_, err := s.assetRepo.FindByID(assetID)
	if err != nil {
		return nil, errors.New("actif introuvable")
	}

	// Récupérer les associations
	ticketAssets, err := s.ticketAssetRepo.FindByAssetID(assetID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des tickets liés")
	}

	var ticketDTOs []dto.TicketDTO
	for _, ta := range ticketAssets {
		if ta.Ticket.ID != 0 {
			ticketDTO := s.ticketToDTO(&ta.Ticket)
			ticketDTOs = append(ticketDTOs, ticketDTO)
		}
	}

	return ticketDTOs, nil
}

// LinkTicket lie un ticket à un actif
func (s *assetService) LinkTicket(assetID uint, ticketID uint, linkedByID uint) error {
	// Vérifier que l'actif existe
	_, err := s.assetRepo.FindByID(assetID)
	if err != nil {
		return errors.New("actif introuvable")
	}

	// Vérifier que le ticket existe
	_, err = s.ticketRepo.FindByID(ticketID)
	if err != nil {
		return errors.New("ticket introuvable")
	}

	// Vérifier que l'association n'existe pas déjà
	existingLinks, _ := s.ticketAssetRepo.FindByAssetID(assetID)
	for _, link := range existingLinks {
		if link.TicketID == ticketID {
			return errors.New("le ticket est déjà lié à cet actif")
		}
	}

	// Créer l'association
	ticketAsset := &models.TicketAsset{
		TicketID: ticketID,
		AssetID:  assetID,
	}

	if err := s.ticketAssetRepo.Create(ticketAsset); err != nil {
		return errors.New("erreur lors de la liaison du ticket")
	}

	return nil
}

// UnlinkTicket supprime la liaison entre un ticket et un actif
func (s *assetService) UnlinkTicket(assetID uint, ticketID uint) error {
	// Vérifier que l'actif existe
	_, err := s.assetRepo.FindByID(assetID)
	if err != nil {
		return errors.New("actif introuvable")
	}

	// Supprimer l'association
	if err := s.ticketAssetRepo.DeleteByTicketAndAsset(ticketID, assetID); err != nil {
		return errors.New("erreur lors de la suppression de la liaison")
	}

	return nil
}

// ticketToDTO convertit un modèle Ticket en DTO (méthode helper)
func (s *assetService) ticketToDTO(ticket *models.Ticket) dto.TicketDTO {
	ticketDTO := dto.TicketDTO{
		ID:          ticket.ID,
		Title:       ticket.Title,
		Description: ticket.Description,
		Category:    ticket.Category,
		Source:      ticket.Source,
		Status:      ticket.Status,
		Priority:    ticket.Priority,
		CreatedAt:   ticket.CreatedAt,
		UpdatedAt:   ticket.UpdatedAt,
	}

	if ticket.EstimatedTime != nil {
		ticketDTO.EstimatedTime = ticket.EstimatedTime
	}
	if ticket.ActualTime != nil {
		ticketDTO.ActualTime = ticket.ActualTime
	}
	if ticket.ClosedAt != nil {
		ticketDTO.ClosedAt = ticket.ClosedAt
	}

	// Convertir CreatedBy
	if ticket.CreatedBy.ID != 0 {
		userDTO := s.userToDTO(&ticket.CreatedBy)
		ticketDTO.CreatedBy = userDTO
	}

	// Convertir AssignedTo
	if ticket.AssignedTo != nil && ticket.AssignedTo.ID != 0 {
		userDTO := s.userToDTO(ticket.AssignedTo)
		ticketDTO.AssignedTo = &userDTO
	}

	return ticketDTO
}

// Delete supprime un actif
func (s *assetService) Delete(id uint) error {
	_, err := s.assetRepo.FindByID(id)
	if err != nil {
		return errors.New("actif introuvable")
	}

	if err := s.assetRepo.Delete(id); err != nil {
		return errors.New("erreur lors de la suppression de l'actif")
	}

	return nil
}

// assetToDTO convertit un modèle Asset en DTO
func (s *assetService) assetToDTO(asset *models.Asset) dto.AssetDTO {
	assetDTO := dto.AssetDTO{
		ID:           asset.ID,
		Name:         asset.Name,
		SerialNumber: asset.SerialNumber,
		Model:        asset.Model,
		Manufacturer: asset.Manufacturer,
		CategoryID:   asset.CategoryID,
		Status:       asset.Status,
		Location:     asset.Location,
		Notes:        asset.Notes,
		CreatedAt:    asset.CreatedAt,
		UpdatedAt:    asset.UpdatedAt,
	}

	if asset.PurchaseDate != nil {
		assetDTO.PurchaseDate = asset.PurchaseDate
	}
	if asset.WarrantyExpiry != nil {
		assetDTO.WarrantyExpiry = asset.WarrantyExpiry
	}

	// Convertir la catégorie si présente
	if asset.Category.ID != 0 {
		categoryDTO := s.categoryToDTO(&asset.Category)
		assetDTO.Category = &categoryDTO
	}

	// Convertir l'utilisateur assigné si présent
	if asset.AssignedToID != nil {
		assetDTO.AssignedTo = asset.AssignedToID
		if asset.AssignedTo != nil && asset.AssignedTo.ID != 0 {
			userDTO := s.userToDTO(asset.AssignedTo)
			assetDTO.AssignedUser = &userDTO
		}
	}

	return assetDTO
}

// categoryToDTO convertit un modèle AssetCategory en DTO
func (s *assetService) categoryToDTO(category *models.AssetCategory) dto.AssetCategoryDTO {
	categoryDTO := dto.AssetCategoryDTO{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}

	if category.ParentID != nil {
		categoryDTO.ParentID = category.ParentID
	}

	return categoryDTO
}

// userToDTO convertit un modèle User en DTO (méthode helper)
func (s *assetService) userToDTO(user *models.User) dto.UserDTO {
	userDTO := dto.UserDTO{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Avatar:    user.Avatar,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	if user.RoleID != 0 {
		userDTO.Role = user.Role.Name
	}

	if user.LastLogin != nil {
		userDTO.LastLogin = user.LastLogin
	}

	return userDTO
}

// AssetCategoryService interface pour les opérations sur les catégories d'actifs
type AssetCategoryService interface {
	Create(req dto.CreateAssetCategoryRequest, createdByID uint) (*dto.AssetCategoryDTO, error)
	GetByID(id uint) (*dto.AssetCategoryDTO, error)
	GetAll() ([]dto.AssetCategoryDTO, error)
	GetByParentID(parentID uint) ([]dto.AssetCategoryDTO, error)
	Update(id uint, req dto.UpdateAssetCategoryRequest, updatedByID uint) (*dto.AssetCategoryDTO, error)
	Delete(id uint) error
}

// assetCategoryService implémente AssetCategoryService
type assetCategoryService struct {
	categoryRepo repositories.AssetCategoryRepository
	userRepo     repositories.UserRepository
}

// NewAssetCategoryService crée une nouvelle instance de AssetCategoryService
func NewAssetCategoryService(
	categoryRepo repositories.AssetCategoryRepository,
	userRepo repositories.UserRepository,
) AssetCategoryService {
	return &assetCategoryService{
		categoryRepo: categoryRepo,
		userRepo:     userRepo,
	}
}

// Create crée une nouvelle catégorie d'actif
func (s *assetCategoryService) Create(req dto.CreateAssetCategoryRequest, createdByID uint) (*dto.AssetCategoryDTO, error) {
	category := &models.AssetCategory{
		Name:        req.Name,
		Description: req.Description,
		ParentID:    req.ParentID,
	}

	if err := s.categoryRepo.Create(category); err != nil {
		return nil, errors.New("erreur lors de la création de la catégorie")
	}

	createdCategory, err := s.categoryRepo.FindByID(category.ID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération de la catégorie créée")
	}

	categoryDTO := s.categoryToDTO(createdCategory)
	return &categoryDTO, nil
}

// GetByID récupère une catégorie par son ID
func (s *assetCategoryService) GetByID(id uint) (*dto.AssetCategoryDTO, error) {
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("catégorie introuvable")
	}

	categoryDTO := s.categoryToDTO(category)
	return &categoryDTO, nil
}

// GetAll récupère toutes les catégories
func (s *assetCategoryService) GetAll() ([]dto.AssetCategoryDTO, error) {
	categories, err := s.categoryRepo.FindAll()
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des catégories")
	}

	var categoryDTOs []dto.AssetCategoryDTO
	for _, category := range categories {
		categoryDTOs = append(categoryDTOs, s.categoryToDTO(&category))
	}

	return categoryDTOs, nil
}

// GetByParentID récupère les catégories enfants d'une catégorie parente
func (s *assetCategoryService) GetByParentID(parentID uint) ([]dto.AssetCategoryDTO, error) {
	categories, err := s.categoryRepo.FindByParentID(parentID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des catégories")
	}

	var categoryDTOs []dto.AssetCategoryDTO
	for _, category := range categories {
		categoryDTOs = append(categoryDTOs, s.categoryToDTO(&category))
	}

	return categoryDTOs, nil
}

// Update met à jour une catégorie
func (s *assetCategoryService) Update(id uint, req dto.UpdateAssetCategoryRequest, updatedByID uint) (*dto.AssetCategoryDTO, error) {
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("catégorie introuvable")
	}

	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Description != "" {
		category.Description = req.Description
	}
	if req.ParentID != nil {
		category.ParentID = req.ParentID
	}

	if err := s.categoryRepo.Update(category); err != nil {
		return nil, errors.New("erreur lors de la mise à jour de la catégorie")
	}

	updatedCategory, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération de la catégorie mise à jour")
	}

	categoryDTO := s.categoryToDTO(updatedCategory)
	return &categoryDTO, nil
}

// Delete supprime une catégorie
func (s *assetCategoryService) Delete(id uint) error {
	_, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return errors.New("catégorie introuvable")
	}

	if err := s.categoryRepo.Delete(id); err != nil {
		return errors.New("erreur lors de la suppression de la catégorie")
	}

	return nil
}

// categoryToDTO convertit un modèle AssetCategory en DTO (pour assetCategoryService)
func (s *assetCategoryService) categoryToDTO(category *models.AssetCategory) dto.AssetCategoryDTO {
	categoryDTO := dto.AssetCategoryDTO{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}

	if category.ParentID != nil {
		categoryDTO.ParentID = category.ParentID
	}

	return categoryDTO
}
