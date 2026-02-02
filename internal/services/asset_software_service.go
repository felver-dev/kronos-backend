package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/models"
	"github.com/mcicare/itsm-backend/internal/repositories"
)

// AssetSoftwareService interface pour les opérations sur les logiciels installés
type AssetSoftwareService interface {
	Create(req dto.CreateAssetSoftwareRequest) (*dto.AssetSoftwareDTO, error)
	GetByID(id uint) (*dto.AssetSoftwareDTO, error)
	GetAll() ([]dto.AssetSoftwareDTO, error)
	GetByAssetID(assetID uint) ([]dto.AssetSoftwareDTO, error)
	GetBySoftwareName(softwareName string) ([]dto.AssetSoftwareDTO, error)
	GetBySoftwareNameAndVersion(softwareName, version string) ([]dto.AssetSoftwareDTO, error)
	Update(id uint, req dto.UpdateAssetSoftwareRequest) (*dto.AssetSoftwareDTO, error)
	Delete(id uint) error
	GetStatistics() (*dto.AssetSoftwareStatisticsDTO, error)
}

// assetSoftwareService implémente AssetSoftwareService
type assetSoftwareService struct {
	assetSoftwareRepo repositories.AssetSoftwareRepository
	assetRepo         repositories.AssetRepository
}

// NewAssetSoftwareService crée une nouvelle instance de AssetSoftwareService
func NewAssetSoftwareService(
	assetSoftwareRepo repositories.AssetSoftwareRepository,
	assetRepo repositories.AssetRepository,
) AssetSoftwareService {
	return &assetSoftwareService{
		assetSoftwareRepo: assetSoftwareRepo,
		assetRepo:         assetRepo,
	}
}

// Create crée un nouveau logiciel installé
func (s *assetSoftwareService) Create(req dto.CreateAssetSoftwareRequest) (*dto.AssetSoftwareDTO, error) {
	// Vérifier que l'actif existe seulement s'il est fourni
	if req.AssetID != nil && *req.AssetID > 0 {
		_, err := s.assetRepo.FindByID(*req.AssetID)
		if err != nil {
			return nil, errors.New("actif introuvable")
		}
	}

	// Parser la date d'installation si fournie
	var installationDate *time.Time
	if req.InstallationDate != nil && *req.InstallationDate != "" {
		parsed, err := time.Parse("2006-01-02", *req.InstallationDate)
		if err == nil {
			installationDate = &parsed
		}
	}

	// Créer le logiciel installé
	software := &models.AssetSoftware{
		AssetID:          req.AssetID,
		SoftwareName:     req.SoftwareName,
		Version:          req.Version,
		LicenseKey:       req.LicenseKey,
		InstallationDate: installationDate,
		Notes:            req.Notes,
	}

	err := s.assetSoftwareRepo.Create(software)
	if err != nil {
		// Retourner l'erreur brute pour faciliter le débogage
		// Si c'est une erreur de contrainte, elle sera visible dans les logs
		return nil, fmt.Errorf("erreur lors de la création du logiciel installé: %w", err)
	}

	// Recharger le logiciel avec les relations si asset_id est fourni
	if software.AssetID != nil && *software.AssetID > 0 {
		createdSoftware, err := s.assetSoftwareRepo.FindByID(software.ID)
		if err == nil {
			software = createdSoftware
		}
	}

	return s.assetSoftwareToDTO(software), nil
}

// GetByID récupère un logiciel installé par son ID
func (s *assetSoftwareService) GetByID(id uint) (*dto.AssetSoftwareDTO, error) {
	software, err := s.assetSoftwareRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("logiciel installé introuvable")
	}

	return s.assetSoftwareToDTO(software), nil
}

// GetAll récupère tous les logiciels installés
func (s *assetSoftwareService) GetAll() ([]dto.AssetSoftwareDTO, error) {
	softwareList, err := s.assetSoftwareRepo.FindAll()
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des logiciels")
	}

	dtos := make([]dto.AssetSoftwareDTO, len(softwareList))
	for i, sw := range softwareList {
		dtos[i] = *s.assetSoftwareToDTO(&sw)
	}

	return dtos, nil
}

// GetByAssetID récupère tous les logiciels installés sur un actif
func (s *assetSoftwareService) GetByAssetID(assetID uint) ([]dto.AssetSoftwareDTO, error) {
	softwareList, err := s.assetSoftwareRepo.FindByAssetID(assetID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des logiciels")
	}

	dtos := make([]dto.AssetSoftwareDTO, len(softwareList))
	for i, sw := range softwareList {
		dtos[i] = *s.assetSoftwareToDTO(&sw)
	}

	return dtos, nil
}

// GetBySoftwareName récupère tous les actifs ayant un logiciel spécifique installé
func (s *assetSoftwareService) GetBySoftwareName(softwareName string) ([]dto.AssetSoftwareDTO, error) {
	softwareList, err := s.assetSoftwareRepo.FindBySoftwareName(softwareName)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des actifs")
	}

	dtos := make([]dto.AssetSoftwareDTO, len(softwareList))
	for i, sw := range softwareList {
		dtos[i] = *s.assetSoftwareToDTO(&sw)
	}

	return dtos, nil
}

// GetBySoftwareNameAndVersion récupère tous les actifs ayant un logiciel avec une version spécifique
func (s *assetSoftwareService) GetBySoftwareNameAndVersion(softwareName, version string) ([]dto.AssetSoftwareDTO, error) {
	softwareList, err := s.assetSoftwareRepo.FindBySoftwareNameAndVersion(softwareName, version)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des actifs")
	}

	dtos := make([]dto.AssetSoftwareDTO, len(softwareList))
	for i, sw := range softwareList {
		dtos[i] = *s.assetSoftwareToDTO(&sw)
	}

	return dtos, nil
}

// Update met à jour un logiciel installé
func (s *assetSoftwareService) Update(id uint, req dto.UpdateAssetSoftwareRequest) (*dto.AssetSoftwareDTO, error) {
	software, err := s.assetSoftwareRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("logiciel installé introuvable")
	}

	// Mettre à jour les champs fournis
	if req.SoftwareName != "" {
		software.SoftwareName = req.SoftwareName
	}
	if req.Version != "" {
		software.Version = req.Version
	}
	if req.LicenseKey != "" {
		software.LicenseKey = req.LicenseKey
	}
	if req.Notes != "" {
		software.Notes = req.Notes
	}
	if req.InstallationDate != nil {
		if *req.InstallationDate == "" {
			software.InstallationDate = nil
		} else {
			parsed, err := time.Parse("2006-01-02", *req.InstallationDate)
			if err == nil {
				software.InstallationDate = &parsed
			}
		}
	}

	err = s.assetSoftwareRepo.Update(software)
	if err != nil {
		return nil, errors.New("erreur lors de la mise à jour du logiciel installé")
	}

	return s.assetSoftwareToDTO(software), nil
}

// Delete supprime un logiciel installé
func (s *assetSoftwareService) Delete(id uint) error {
	_, err := s.assetSoftwareRepo.FindByID(id)
	if err != nil {
		return errors.New("logiciel installé introuvable")
	}

	err = s.assetSoftwareRepo.Delete(id)
	if err != nil {
		return errors.New("erreur lors de la suppression du logiciel installé")
	}

	return nil
}

// GetStatistics récupère des statistiques sur les logiciels installés
func (s *assetSoftwareService) GetStatistics() (*dto.AssetSoftwareStatisticsDTO, error) {
	stats, err := s.assetSoftwareRepo.GetStatistics()
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des statistiques")
	}

	result := &dto.AssetSoftwareStatisticsDTO{}

	// Convertir les statistiques par logiciel et version
	if bySoftware, ok := stats["by_software"].([]struct {
		SoftwareName string
		Version      string
		Count        int64
		CategoryName string
	}); ok {
		result.BySoftware = make([]dto.SoftwareCountDTO, len(bySoftware))
		for i, s := range bySoftware {
			result.BySoftware[i] = dto.SoftwareCountDTO{
				SoftwareName: s.SoftwareName,
				Version:      s.Version,
				Count:        s.Count,
				CategoryName: s.CategoryName,
			}
		}
	}

	// Convertir les statistiques par nom de logiciel
	if bySoftwareName, ok := stats["by_software_name"].([]struct {
		SoftwareName string
		Count        int64
	}); ok {
		result.BySoftwareName = make([]dto.SoftwareNameCountDTO, len(bySoftwareName))
		for i, s := range bySoftwareName {
			result.BySoftwareName[i] = dto.SoftwareNameCountDTO{
				SoftwareName: s.SoftwareName,
				Count:        s.Count,
			}
		}
	}

	return result, nil
}

// assetSoftwareToDTO convertit un modèle AssetSoftware en DTO
func (s *assetSoftwareService) assetSoftwareToDTO(software *models.AssetSoftware) *dto.AssetSoftwareDTO {
	softwareDTO := &dto.AssetSoftwareDTO{
		ID:              software.ID,
		AssetID:         software.AssetID,
		SoftwareName:    software.SoftwareName,
		Version:          software.Version,
		LicenseKey:       software.LicenseKey,
		InstallationDate: software.InstallationDate,
		Notes:            software.Notes,
		CreatedAt:        software.CreatedAt,
		UpdatedAt:        software.UpdatedAt,
	}

	// Inclure l'actif si chargé (simplifié - juste les infos de base)
	// Vérifier que Asset n'est pas nil et a un ID valide
	if software.Asset != nil && software.Asset.ID != 0 {
		assetDTO := dto.AssetDTO{
			ID:           software.Asset.ID,
			Name:         software.Asset.Name,
			SerialNumber: software.Asset.SerialNumber,
			Model:        software.Asset.Model,
			Manufacturer: software.Asset.Manufacturer,
			CategoryID:   software.Asset.CategoryID,
			Status:       software.Asset.Status,
		}

		// Convertir l'utilisateur assigné si présent
		if software.Asset.AssignedToID != nil {
			assetDTO.AssignedTo = software.Asset.AssignedToID
			// Vérifier si AssignedTo est chargé (préchargé)
			if software.Asset.AssignedTo != nil && software.Asset.AssignedTo.ID != 0 {
				userDTO := dto.UserDTO{
					ID:        software.Asset.AssignedTo.ID,
					Username:  software.Asset.AssignedTo.Username,
					Email:     software.Asset.AssignedTo.Email,
					FirstName: software.Asset.AssignedTo.FirstName,
					LastName:  software.Asset.AssignedTo.LastName,
					IsActive:  software.Asset.AssignedTo.IsActive,
					CreatedAt: software.Asset.AssignedTo.CreatedAt,
					UpdatedAt: software.Asset.AssignedTo.UpdatedAt,
				}
				// Ajouter le rôle si chargé
				if software.Asset.AssignedTo.RoleID != 0 && software.Asset.AssignedTo.Role.ID != 0 {
					userDTO.Role = software.Asset.AssignedTo.Role.Name
				}
				assetDTO.AssignedUser = &userDTO
			}
		}

		if software.Asset.Category.ID != 0 {
			categoryDTO := dto.AssetCategoryDTO{
				ID:   software.Asset.Category.ID,
				Name: software.Asset.Category.Name,
			}
			assetDTO.Category = &categoryDTO
		}
		softwareDTO.Asset = &assetDTO
	}

	return softwareDTO
}
