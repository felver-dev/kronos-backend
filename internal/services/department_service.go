package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/models"
	"github.com/mcicare/itsm-backend/internal/repositories"
)

// DepartmentService interface pour les opérations sur les départements
type DepartmentService interface {
	Create(req dto.CreateDepartmentRequest) (*dto.DepartmentDTO, error)
	GetAll(activeOnly bool) ([]dto.DepartmentDTO, error)
	GetByID(id uint) (*dto.DepartmentDTO, error)
	GetByCode(code string) (*dto.DepartmentDTO, error)
	Update(id uint, req dto.UpdateDepartmentRequest) (*dto.DepartmentDTO, error)
	Delete(id uint) error
	GetByOfficeID(officeID uint) ([]dto.DepartmentDTO, error)
	GetByFilialeID(filialeID uint) ([]dto.DepartmentDTO, error)
}

// departmentService implémente DepartmentService
type departmentService struct {
	departmentRepo repositories.DepartmentRepository
	officeRepo     repositories.OfficeRepository
	filialeRepo    repositories.FilialeRepository
}

// NewDepartmentService crée une nouvelle instance de DepartmentService
func NewDepartmentService(
	departmentRepo repositories.DepartmentRepository,
	officeRepo repositories.OfficeRepository,
	filialeRepo repositories.FilialeRepository,
) DepartmentService {
	return &departmentService{
		departmentRepo: departmentRepo,
		officeRepo:     officeRepo,
		filialeRepo:    filialeRepo,
	}
}

func prefixCodeWithFiliale(filialeCode string, raw string) string {
	raw = strings.TrimSpace(raw)
	fc := strings.ToUpper(strings.TrimSpace(filialeCode))
	prefix := fc + "-"
	if raw == "" {
		return prefix
	}
	upperRaw := strings.ToUpper(raw)
	upperPrefix := strings.ToUpper(prefix)
	if strings.HasPrefix(upperRaw, upperPrefix) && len(raw) >= len(prefix) {
		// normaliser le préfixe (ex: niger-dev -> NIGER-dev)
		return prefix + raw[len(prefix):]
	}
	return prefix + raw
}

// Create crée un nouveau département
func (s *departmentService) Create(req dto.CreateDepartmentRequest) (*dto.DepartmentDTO, error) {
	// Vérifier que la filiale existe
	if req.FilialeID == nil {
		return nil, errors.New("la filiale est obligatoire")
	}
	filiale, err := s.filialeRepo.FindByID(*req.FilialeID)
	if err != nil {
		return nil, errors.New("filiale introuvable")
	}

	// Préfixer automatiquement le code avec le code filiale
	req.Code = prefixCodeWithFiliale(filiale.Code, req.Code)

	// Vérifier si le code existe déjà (après préfixage)
	existing, _ := s.departmentRepo.FindByCode(req.Code)
	if existing != nil {
		return nil, errors.New("un département avec ce code existe déjà")
	}

	// Vérifier si le siège existe (si fourni)
	if req.OfficeID != nil {
		_, err := s.officeRepo.FindByID(*req.OfficeID)
		if err != nil {
			return nil, errors.New("siège introuvable")
		}
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	// Vérifier que le code n'est pas vide après préfixage
	if strings.TrimSpace(req.Code) == "" {
		return nil, errors.New("le code du département ne peut pas être vide")
	}

	// Gérer IsITDepartment : uniquement pour la filiale fournisseur IT
	isITDepartment := false
	if req.IsITDepartment != nil && *req.IsITDepartment {
		// Vérifier que c'est bien la filiale fournisseur IT
		if !filiale.IsSoftwareProvider {
			return nil, errors.New("seuls les départements de la filiale fournisseur IT peuvent être marqués comme département IT")
		}
		isITDepartment = true
	}

	department := &models.Department{
		Name:           req.Name,
		Code:           req.Code,
		Description:    req.Description,
		FilialeID:      req.FilialeID,
		OfficeID:       req.OfficeID,
		IsActive:       isActive,
		IsITDepartment: isITDepartment,
	}

	if err := s.departmentRepo.Create(department); err != nil {
		// Analyser l'erreur pour donner un message plus précis
		errStr := err.Error()

		// Erreur de contrainte unique (code déjà existant)
		if strings.Contains(errStr, "Duplicate entry") || strings.Contains(errStr, "UNIQUE constraint") || strings.Contains(errStr, "duplicate key") {
			return nil, fmt.Errorf("un département avec le code '%s' existe déjà", req.Code)
		}

		// Erreur de contrainte de clé étrangère
		if strings.Contains(errStr, "foreign key constraint") {
			if strings.Contains(errStr, "filiale") {
				return nil, errors.New("la filiale spécifiée est invalide ou introuvable")
			}
			if strings.Contains(errStr, "office") || strings.Contains(errStr, "siège") {
				return nil, errors.New("le siège spécifié est invalide ou introuvable")
			}
			return nil, errors.New("référence invalide (filiale ou siège)")
		}

		// Erreur de contrainte NOT NULL
		if strings.Contains(errStr, "NOT NULL") {
			if strings.Contains(errStr, "code") {
				return nil, errors.New("le code du département est obligatoire")
			}
			if strings.Contains(errStr, "name") {
				return nil, errors.New("le nom du département est obligatoire")
			}
			return nil, errors.New("certains champs obligatoires sont manquants")
		}

		// Erreur de longueur de champ
		if strings.Contains(errStr, "Data too long") || strings.Contains(errStr, "value too long") {
			return nil, errors.New("le code ou le nom du département est trop long")
		}

		// Erreur générique avec détails pour le débogage
		return nil, fmt.Errorf("erreur lors de la création du département: %v", err)
	}

	// Récupérer le département créé
	createdDepartment, err := s.departmentRepo.FindByID(department.ID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération du département créé")
	}

	return s.departmentToDTO(createdDepartment), nil
}

// GetAll récupère tous les départements
func (s *departmentService) GetAll(activeOnly bool) ([]dto.DepartmentDTO, error) {
	var departments []models.Department
	var err error

	if activeOnly {
		departments, err = s.departmentRepo.FindActive()
	} else {
		departments, err = s.departmentRepo.FindAll()
	}

	if err != nil {
		return nil, errors.New("erreur lors de la récupération des départements")
	}

	var departmentDTOs []dto.DepartmentDTO
	for _, department := range departments {
		departmentDTOs = append(departmentDTOs, *s.departmentToDTO(&department))
	}

	return departmentDTOs, nil
}

// GetByID récupère un département par son ID
func (s *departmentService) GetByID(id uint) (*dto.DepartmentDTO, error) {
	department, err := s.departmentRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("département introuvable")
	}

	return s.departmentToDTO(department), nil
}

// GetByCode récupère un département par son code
func (s *departmentService) GetByCode(code string) (*dto.DepartmentDTO, error) {
	department, err := s.departmentRepo.FindByCode(code)
	if err != nil {
		return nil, errors.New("département introuvable")
	}

	return s.departmentToDTO(department), nil
}

// Update met à jour un département
func (s *departmentService) Update(id uint, req dto.UpdateDepartmentRequest) (*dto.DepartmentDTO, error) {
	department, err := s.departmentRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("département introuvable")
	}

	// Si le code est modifié, le préfixer avec le code filiale (celle du département ou celle demandée)
	if req.Code != nil {
		targetFilialeID := department.FilialeID
		if req.FilialeID != nil {
			targetFilialeID = req.FilialeID
		}
		if targetFilialeID == nil {
			return nil, errors.New("la filiale est obligatoire pour définir le code")
		}
		filiale, err := s.filialeRepo.FindByID(*targetFilialeID)
		if err != nil {
			return nil, errors.New("filiale introuvable")
		}
		prefixed := prefixCodeWithFiliale(filiale.Code, *req.Code)
		req.Code = &prefixed

		// Vérifier si le code existe déjà (si modifié)
		if *req.Code != department.Code {
			existing, _ := s.departmentRepo.FindByCode(*req.Code)
			if existing != nil && existing.ID != id {
				return nil, errors.New("un département avec ce code existe déjà")
			}
		}
	}

	// Vérifier si le siège existe (si fourni)
	if req.OfficeID != nil {
		_, err := s.officeRepo.FindByID(*req.OfficeID)
		if err != nil {
			return nil, errors.New("siège introuvable")
		}
	}

	// Gérer IsITDepartment : uniquement pour la filiale fournisseur IT
	if req.IsITDepartment != nil {
		// Déterminer la filiale à vérifier (celle fournie dans la requête ou celle actuelle du département)
		var filialeToCheck *models.Filiale
		if req.FilialeID != nil {
			filialeToCheck, err = s.filialeRepo.FindByID(*req.FilialeID)
			if err != nil {
				return nil, errors.New("filiale introuvable")
			}
		} else if department.FilialeID != nil {
			filialeToCheck, err = s.filialeRepo.FindByID(*department.FilialeID)
			if err != nil {
				return nil, errors.New("filiale introuvable")
			}
		}

		// Vérifier que c'est bien la filiale fournisseur IT si on veut marquer comme IT
		if *req.IsITDepartment && (filialeToCheck == nil || !filialeToCheck.IsSoftwareProvider) {
			return nil, errors.New("seuls les départements de la filiale fournisseur IT peuvent être marqués comme département IT")
		}
		department.IsITDepartment = *req.IsITDepartment
	}

	// Mettre à jour les champs fournis
	if req.Name != nil {
		department.Name = *req.Name
	}
	if req.Code != nil {
		department.Code = *req.Code
	}
	if req.Description != nil {
		department.Description = req.Description
	}
	if req.FilialeID != nil {
		department.FilialeID = req.FilialeID
	}
	if req.OfficeID != nil {
		department.OfficeID = req.OfficeID
	}
	if req.IsActive != nil {
		department.IsActive = *req.IsActive
	}

	if err := s.departmentRepo.Update(department); err != nil {
		return nil, errors.New("erreur lors de la mise à jour du département")
	}

	// Récupérer le département mis à jour
	updatedDepartment, err := s.departmentRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération du département mis à jour")
	}

	return s.departmentToDTO(updatedDepartment), nil
}

// Delete supprime un département
func (s *departmentService) Delete(id uint) error {
	_, err := s.departmentRepo.FindByID(id)
	if err != nil {
		return errors.New("département introuvable")
	}

	if err := s.departmentRepo.Delete(id); err != nil {
		return errors.New("erreur lors de la suppression du département")
	}

	return nil
}

// GetByOfficeID récupère les départements d'un siège
func (s *departmentService) GetByOfficeID(officeID uint) ([]dto.DepartmentDTO, error) {
	departments, err := s.departmentRepo.FindByOfficeID(officeID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des départements")
	}

	var departmentDTOs []dto.DepartmentDTO
	for _, department := range departments {
		departmentDTOs = append(departmentDTOs, *s.departmentToDTO(&department))
	}

	return departmentDTOs, nil
}

// GetByFilialeID récupère les départements actifs d'une filiale
func (s *departmentService) GetByFilialeID(filialeID uint) ([]dto.DepartmentDTO, error) {
	departments, err := s.departmentRepo.FindByFilialeID(filialeID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des départements")
	}

	var departmentDTOs []dto.DepartmentDTO
	for _, department := range departments {
		departmentDTOs = append(departmentDTOs, *s.departmentToDTO(&department))
	}

	return departmentDTOs, nil
}

// departmentToDTO convertit un modèle Department en DTO
func (s *departmentService) departmentToDTO(department *models.Department) *dto.DepartmentDTO {
	departmentDTO := &dto.DepartmentDTO{
		ID:             department.ID,
		Name:           department.Name,
		Code:           department.Code,
		Description:    department.Description,
		FilialeID:      department.FilialeID,
		OfficeID:       department.OfficeID,
		IsActive:       department.IsActive,
		IsITDepartment: department.IsITDepartment,
		CreatedAt:      department.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      department.UpdatedAt.Format(time.RFC3339),
	}

	// Inclure la filiale si présente
	if department.Filiale != nil {
		departmentDTO.Filiale = &dto.FilialeDTO{
			ID:          department.Filiale.ID,
			Code:        department.Filiale.Code,
			Name:        department.Filiale.Name,
			Country:     department.Filiale.Country,
			City:        department.Filiale.City,
			Address:     department.Filiale.Address,
			Phone:       department.Filiale.Phone,
			Email:       department.Filiale.Email,
			IsActive:    department.Filiale.IsActive,
			IsSoftwareProvider: department.Filiale.IsSoftwareProvider,
			CreatedAt:   department.Filiale.CreatedAt,
			UpdatedAt:   department.Filiale.UpdatedAt,
		}
	}

	// Inclure le siège si présent
	if department.Office != nil {
		departmentDTO.Office = &dto.OfficeDTO{
			ID:        department.Office.ID,
			Name:      department.Office.Name,
			Country:   department.Office.Country,
			City:      department.Office.City,
			Commune:   department.Office.Commune,
			Address:   department.Office.Address,
			Longitude: department.Office.Longitude,
			Latitude:  department.Office.Latitude,
			IsActive:  department.Office.IsActive,
			CreatedAt: department.Office.CreatedAt.Format(time.RFC3339),
			UpdatedAt: department.Office.UpdatedAt.Format(time.RFC3339),
		}
	}

	return departmentDTO
}
