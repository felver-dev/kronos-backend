package services

import (
	"errors"
	"time"

	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/models"
	"github.com/mcicare/itsm-backend/internal/repositories"
)

// FilialeSoftwareService interface pour les opérations sur les déploiements de logiciels
type FilialeSoftwareService interface {
	Create(req dto.CreateFilialeSoftwareRequest) (*dto.FilialeSoftwareDTO, error)
	GetByID(id uint) (*dto.FilialeSoftwareDTO, error)
	GetAll() ([]dto.FilialeSoftwareDTO, error)
	GetByFilialeID(filialeID uint) ([]dto.FilialeSoftwareDTO, error)
	GetBySoftwareID(softwareID uint) ([]dto.FilialeSoftwareDTO, error)
	GetByFilialeAndSoftware(filialeID, softwareID uint) (*dto.FilialeSoftwareDTO, error)
	GetActive() ([]dto.FilialeSoftwareDTO, error)
	GetActiveByFiliale(filialeID uint) ([]dto.FilialeSoftwareDTO, error)
	GetActiveBySoftware(softwareID uint) ([]dto.FilialeSoftwareDTO, error)
	Update(id uint, req dto.UpdateFilialeSoftwareRequest) (*dto.FilialeSoftwareDTO, error)
	Delete(id uint) error
}

// filialeSoftwareService implémente FilialeSoftwareService
type filialeSoftwareService struct {
	deploymentRepo repositories.FilialeSoftwareRepository
	filialeRepo    repositories.FilialeRepository
	softwareRepo   repositories.SoftwareRepository
}

// NewFilialeSoftwareService crée une nouvelle instance de FilialeSoftwareService
func NewFilialeSoftwareService(
	deploymentRepo repositories.FilialeSoftwareRepository,
	filialeRepo repositories.FilialeRepository,
	softwareRepo repositories.SoftwareRepository,
) FilialeSoftwareService {
	return &filialeSoftwareService{
		deploymentRepo: deploymentRepo,
		filialeRepo:    filialeRepo,
		softwareRepo:   softwareRepo,
	}
}

// Create crée un nouveau déploiement
func (s *filialeSoftwareService) Create(req dto.CreateFilialeSoftwareRequest) (*dto.FilialeSoftwareDTO, error) {
	// Vérifier que FilialeID est fourni
	if req.FilialeID == 0 {
		return nil, errors.New("filiale_id est obligatoire")
	}

	// Vérifier que la filiale existe
	_, err := s.filialeRepo.FindByID(req.FilialeID)
	if err != nil {
		return nil, errors.New("filiale introuvable")
	}

	// Vérifier que le logiciel existe
	_, err = s.softwareRepo.FindByID(req.SoftwareID)
	if err != nil {
		return nil, errors.New("logiciel introuvable")
	}

	deployedAt := req.DeployedAt
	if deployedAt == nil {
		now := time.Now()
		deployedAt = &now
	}

	deployment := &models.FilialeSoftware{
		FilialeID:  req.FilialeID,
		SoftwareID: req.SoftwareID,
		Version:    "", // Non utilisé - la version du logiciel est dans la table Software
		DeployedAt: deployedAt,
		IsActive:   true,
		Notes:      req.Notes,
	}

	if err := s.deploymentRepo.Create(deployment); err != nil {
		return nil, errors.New("erreur lors de la création du déploiement")
	}

	// Récupérer le déploiement créé
	createdDeployment, err := s.deploymentRepo.FindByID(deployment.ID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération du déploiement créé")
	}

	return s.deploymentToDTO(createdDeployment), nil
}

// GetByID récupère un déploiement par son ID
func (s *filialeSoftwareService) GetByID(id uint) (*dto.FilialeSoftwareDTO, error) {
	deployment, err := s.deploymentRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("déploiement introuvable")
	}

	return s.deploymentToDTO(deployment), nil
}

// GetAll récupère tous les déploiements
func (s *filialeSoftwareService) GetAll() ([]dto.FilialeSoftwareDTO, error) {
	deployments, err := s.deploymentRepo.FindAll()
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des déploiements")
	}

	var deploymentDTOs []dto.FilialeSoftwareDTO
	for _, deployment := range deployments {
		deploymentDTOs = append(deploymentDTOs, *s.deploymentToDTO(&deployment))
	}

	return deploymentDTOs, nil
}

// GetByFilialeID récupère tous les déploiements d'une filiale
func (s *filialeSoftwareService) GetByFilialeID(filialeID uint) ([]dto.FilialeSoftwareDTO, error) {
	deployments, err := s.deploymentRepo.FindByFilialeID(filialeID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des déploiements")
	}

	var deploymentDTOs []dto.FilialeSoftwareDTO
	for _, deployment := range deployments {
		deploymentDTOs = append(deploymentDTOs, *s.deploymentToDTO(&deployment))
	}

	return deploymentDTOs, nil
}

// GetBySoftwareID récupère tous les déploiements d'un logiciel
func (s *filialeSoftwareService) GetBySoftwareID(softwareID uint) ([]dto.FilialeSoftwareDTO, error) {
	deployments, err := s.deploymentRepo.FindBySoftwareID(softwareID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des déploiements")
	}

	var deploymentDTOs []dto.FilialeSoftwareDTO
	for _, deployment := range deployments {
		deploymentDTOs = append(deploymentDTOs, *s.deploymentToDTO(&deployment))
	}

	return deploymentDTOs, nil
}

// GetByFilialeAndSoftware récupère un déploiement spécifique
func (s *filialeSoftwareService) GetByFilialeAndSoftware(filialeID, softwareID uint) (*dto.FilialeSoftwareDTO, error) {
	deployment, err := s.deploymentRepo.FindByFilialeAndSoftware(filialeID, softwareID)
	if err != nil {
		return nil, errors.New("déploiement introuvable")
	}

	return s.deploymentToDTO(deployment), nil
}

// GetActive récupère tous les déploiements actifs
func (s *filialeSoftwareService) GetActive() ([]dto.FilialeSoftwareDTO, error) {
	deployments, err := s.deploymentRepo.FindActive()
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des déploiements actifs")
	}

	var deploymentDTOs []dto.FilialeSoftwareDTO
	for _, deployment := range deployments {
		deploymentDTOs = append(deploymentDTOs, *s.deploymentToDTO(&deployment))
	}

	return deploymentDTOs, nil
}

// GetActiveByFiliale récupère les déploiements actifs d'une filiale
func (s *filialeSoftwareService) GetActiveByFiliale(filialeID uint) ([]dto.FilialeSoftwareDTO, error) {
	deployments, err := s.deploymentRepo.FindActiveByFiliale(filialeID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des déploiements actifs")
	}

	var deploymentDTOs []dto.FilialeSoftwareDTO
	for _, deployment := range deployments {
		deploymentDTOs = append(deploymentDTOs, *s.deploymentToDTO(&deployment))
	}

	return deploymentDTOs, nil
}

// GetActiveBySoftware récupère les déploiements actifs d'un logiciel
func (s *filialeSoftwareService) GetActiveBySoftware(softwareID uint) ([]dto.FilialeSoftwareDTO, error) {
	deployments, err := s.deploymentRepo.FindActiveBySoftware(softwareID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des déploiements actifs")
	}

	var deploymentDTOs []dto.FilialeSoftwareDTO
	for _, deployment := range deployments {
		deploymentDTOs = append(deploymentDTOs, *s.deploymentToDTO(&deployment))
	}

	return deploymentDTOs, nil
}

// Update met à jour un déploiement
func (s *filialeSoftwareService) Update(id uint, req dto.UpdateFilialeSoftwareRequest) (*dto.FilialeSoftwareDTO, error) {
	deployment, err := s.deploymentRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("déploiement introuvable")
	}

	// Mettre à jour les champs fournis
	// Note: Version (numéro de déploiement) n'est pas modifiable car calculé automatiquement lors de la création
	if req.DeployedAt != nil {
		deployment.DeployedAt = req.DeployedAt
	}
	if req.IsActive != nil {
		deployment.IsActive = *req.IsActive
	}
	if req.Notes != nil {
		deployment.Notes = req.Notes
	}

	if err := s.deploymentRepo.Update(deployment); err != nil {
		return nil, errors.New("erreur lors de la mise à jour du déploiement")
	}

	// Récupérer le déploiement mis à jour
	updatedDeployment, err := s.deploymentRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération du déploiement mis à jour")
	}

	return s.deploymentToDTO(updatedDeployment), nil
}

// Delete supprime un déploiement
func (s *filialeSoftwareService) Delete(id uint) error {
	_, err := s.deploymentRepo.FindByID(id)
	if err != nil {
		return errors.New("déploiement introuvable")
	}

	if err := s.deploymentRepo.Delete(id); err != nil {
		return errors.New("erreur lors de la suppression du déploiement")
	}

	return nil
}

// deploymentToDTO convertit un modèle FilialeSoftware en DTO
func (s *filialeSoftwareService) deploymentToDTO(deployment *models.FilialeSoftware) *dto.FilialeSoftwareDTO {
	deploymentDTO := &dto.FilialeSoftwareDTO{
		ID:         deployment.ID,
		FilialeID:  deployment.FilialeID,
		SoftwareID: deployment.SoftwareID,
		Version:    deployment.Version,
		DeployedAt: deployment.DeployedAt,
		IsActive:   deployment.IsActive,
		Notes:      deployment.Notes,
		CreatedAt:  deployment.CreatedAt,
		UpdatedAt:  deployment.UpdatedAt,
	}

	// Inclure la filiale si présente
	if deployment.Filiale.ID != 0 {
		deploymentDTO.Filiale = dto.FilialeDTO{
			ID:          deployment.Filiale.ID,
			Code:        deployment.Filiale.Code,
			Name:        deployment.Filiale.Name,
			Country:     deployment.Filiale.Country,
			City:        deployment.Filiale.City,
			Address:     deployment.Filiale.Address,
			Phone:       deployment.Filiale.Phone,
			Email:       deployment.Filiale.Email,
			IsActive:    deployment.Filiale.IsActive,
			IsSoftwareProvider: deployment.Filiale.IsSoftwareProvider,
			CreatedAt:   deployment.Filiale.CreatedAt,
			UpdatedAt:   deployment.Filiale.UpdatedAt,
		}
	}

	// Inclure le logiciel si présent
	if deployment.Software.ID != 0 {
		deploymentDTO.Software = dto.SoftwareDTO{
			ID:          deployment.Software.ID,
			Code:        deployment.Software.Code,
			Name:        deployment.Software.Name,
			Description: deployment.Software.Description,
			Version:     deployment.Software.Version,
			IsActive:    deployment.Software.IsActive,
			CreatedAt:   deployment.Software.CreatedAt,
			UpdatedAt:   deployment.Software.UpdatedAt,
		}
	}

	return deploymentDTO
}
