package services

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/models"
	"github.com/mcicare/itsm-backend/internal/repositories"
)

// OfficeService interface pour les opérations sur les sièges
type OfficeService interface {
	Create(req dto.CreateOfficeRequest) (*dto.OfficeDTO, error)
	GetAll(activeOnly bool) ([]dto.OfficeDTO, error)
	GetByID(id uint) (*dto.OfficeDTO, error)
	GetByFilialeID(filialeID uint) ([]dto.OfficeDTO, error)
	Update(id uint, req dto.UpdateOfficeRequest) (*dto.OfficeDTO, error)
	Delete(id uint) error
	GetByCountry(country string) ([]dto.OfficeDTO, error)
	GetByCity(city string) ([]dto.OfficeDTO, error)
}

// officeService implémente OfficeService
type officeService struct {
	officeRepo  repositories.OfficeRepository
	filialeRepo repositories.FilialeRepository
}

// NewOfficeService crée une nouvelle instance de OfficeService
func NewOfficeService(officeRepo repositories.OfficeRepository, filialeRepo repositories.FilialeRepository) OfficeService {
	return &officeService{
		officeRepo:  officeRepo,
		filialeRepo: filialeRepo,
	}
}

// generateCodeFromName génère un code à partir du nom du siège
// Exemple: "Siège Plateau 1" -> "SP1" ou "SP-001"
func generateCodeFromName(name string) string {
	// Mots à ignorer (articles, prépositions, etc.)
	stopWords := map[string]bool{
		"le": true, "la": true, "les": true, "de": true, "du": true, "des": true,
		"siège": true, "siege": true, "bureau": true, "bureaux": true,
		"à": true, "au": true, "aux": true, "et": true, "ou": true,
	}

	// Nettoyer le nom et le diviser en mots
	name = strings.TrimSpace(name)
	words := regexp.MustCompile(`\s+`).Split(name, -1)

	var initials []string
	var number string

	// Extraire les initiales et les numéros
	for _, word := range words {
		wordLower := strings.ToLower(strings.TrimSpace(word))
		if stopWords[wordLower] {
			continue
		}
		if wordLower == "" {
			continue
		}

		// Vérifier si le mot contient un numéro
		numMatch := regexp.MustCompile(`\d+`).FindString(word)
		if numMatch != "" {
			number = numMatch
			// Extraire les lettres avant le numéro
			letters := regexp.MustCompile(`[a-zA-Z]+`).FindString(word)
			if letters != "" {
				initials = append(initials, strings.ToUpper(letters[:1]))
			}
		} else {
			// Prendre la première lettre en majuscule
			if len(word) > 0 {
				firstChar := strings.ToUpper(string(word[0]))
				initials = append(initials, firstChar)
			}
		}
	}

	// Construire le code de base
	codeBase := strings.Join(initials, "")
	if codeBase == "" {
		codeBase = "OFF" // Par défaut si aucun mot valide
	}

	// Ajouter le numéro si présent
	if number != "" {
		// Format: SP1 ou SP-001 selon la longueur du numéro
		if len(number) <= 2 {
			return codeBase + number
		}
		return codeBase + "-" + number
	}

	return codeBase
}

// findNextAvailableCode trouve le prochain code disponible en incrémentant le numéro
func (s *officeService) findNextAvailableCode(baseCode string, filialeCode string) string {
	// Le code de base doit déjà être préfixé avec la filiale
	// Exemple: NIGER-SP1, NIGER-SP-001

	// Extraire le préfixe filiale et le suffixe
	prefix := strings.ToUpper(filialeCode) + "-"
	if !strings.HasPrefix(strings.ToUpper(baseCode), strings.ToUpper(prefix)) {
		baseCode = prefix + baseCode
	}

	// Extraire la partie après le préfixe filiale
	suffix := baseCode[len(prefix):]

	// Vérifier si le suffixe contient un numéro
	numMatch := regexp.MustCompile(`(\d+)$`).FindString(suffix)
	if numMatch == "" {
		// Pas de numéro, ajouter -001
		suffix = suffix + "-001"
		baseCode = prefix + suffix
	}

	// Essayer le code tel quel
	testCode := baseCode
	counter := 0
	maxAttempts := 1000

	for counter < maxAttempts {
		_, err := s.officeRepo.FindByCode(testCode)
		if err != nil {
			// Code disponible
			return testCode
		}

		// Code existe, incrémenter
		counter++
		// Extraire le numéro final
		parts := regexp.MustCompile(`(.+?)(\d+)$`).FindStringSubmatch(suffix)
		if len(parts) == 3 {
			basePart := parts[1]
			numPart := parts[2]
			// Incrémenter le numéro
			var num int
			fmt.Sscanf(numPart, "%d", &num)
			num++
			suffix = basePart + fmt.Sprintf("%03d", num)
		} else {
			// Ajouter un numéro séquentiel
			suffix = suffix + fmt.Sprintf("-%03d", counter)
		}
		testCode = prefix + suffix
	}

	// Fallback: utiliser un timestamp
	return prefix + suffix + fmt.Sprintf("-%d", time.Now().Unix()%10000)
}

// Create crée un nouveau siège
func (s *officeService) Create(req dto.CreateOfficeRequest) (*dto.OfficeDTO, error) {
	// La filiale est obligatoire (un siège appartient forcément à une filiale)
	if req.FilialeID == nil {
		return nil, errors.New("la filiale est obligatoire")
	}
	filiale, err := s.filialeRepo.FindByID(*req.FilialeID)
	if err != nil {
		return nil, errors.New("filiale introuvable")
	}

	// Générer le code automatiquement si non fourni ou vide
	var code string
	if req.Code == "" || strings.TrimSpace(req.Code) == "" {
		// Générer le code à partir du nom
		generatedCode := generateCodeFromName(req.Name)
		// Trouver le prochain code disponible
		code = s.findNextAvailableCode(generatedCode, filiale.Code)
	} else {
		// Préfixer automatiquement le code avec le code filiale
		code = prefixCodeWithFiliale(filiale.Code, req.Code)
		// Vérifier si le code existe déjà
		existing, _ := s.officeRepo.FindByCode(code)
		if existing != nil {
			return nil, errors.New("un siège avec ce code existe déjà")
		}
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	office := &models.Office{
		Name:      req.Name,
		Code:      &code,
		Country:   req.Country,
		City:      req.City,
		Commune:   req.Commune,
		Address:   req.Address,
		FilialeID: req.FilialeID,
		Longitude: req.Longitude,
		Latitude:  req.Latitude,
		IsActive:  isActive,
	}

	if err := s.officeRepo.Create(office); err != nil {
		return nil, errors.New("erreur lors de la création du siège")
	}

	// Récupérer le siège créé
	createdOffice, err := s.officeRepo.FindByID(office.ID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération du siège créé")
	}

	return s.officeToDTO(createdOffice), nil
}

// GetAll récupère tous les sièges
func (s *officeService) GetAll(activeOnly bool) ([]dto.OfficeDTO, error) {
	var offices []models.Office
	var err error

	if activeOnly {
		offices, err = s.officeRepo.FindActive()
	} else {
		offices, err = s.officeRepo.FindAll()
	}

	if err != nil {
		return nil, errors.New("erreur lors de la récupération des sièges")
	}

	var officeDTOs []dto.OfficeDTO
	for _, office := range offices {
		officeDTOs = append(officeDTOs, *s.officeToDTO(&office))
	}

	return officeDTOs, nil
}

// GetByFilialeID récupère les sièges actifs d'une filiale
func (s *officeService) GetByFilialeID(filialeID uint) ([]dto.OfficeDTO, error) {
	offices, err := s.officeRepo.FindByFilialeID(filialeID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des sièges")
	}

	var officeDTOs []dto.OfficeDTO
	for _, office := range offices {
		officeDTOs = append(officeDTOs, *s.officeToDTO(&office))
	}

	return officeDTOs, nil
}

// GetByID récupère un siège par son ID
func (s *officeService) GetByID(id uint) (*dto.OfficeDTO, error) {
	office, err := s.officeRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("siège introuvable")
	}

	return s.officeToDTO(office), nil
}

// Update met à jour un siège
func (s *officeService) Update(id uint, req dto.UpdateOfficeRequest) (*dto.OfficeDTO, error) {
	office, err := s.officeRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("siège introuvable")
	}

	// Mettre à jour les champs fournis
	if req.Name != nil {
		office.Name = *req.Name
	}
	if req.Country != nil {
		office.Country = *req.Country
	}
	if req.City != nil {
		office.City = *req.City
	}
	if req.Commune != nil {
		office.Commune = req.Commune
	}
	if req.Address != nil {
		office.Address = req.Address
	}
	if req.FilialeID != nil {
		office.FilialeID = req.FilialeID
	}
	if req.Code != nil {
		// Déterminer la filiale à utiliser pour le préfixe (celle en place, ou celle demandée)
		targetFilialeID := office.FilialeID
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
		office.Code = &prefixed
	}
	if req.Longitude != nil {
		office.Longitude = req.Longitude
	}
	if req.Latitude != nil {
		office.Latitude = req.Latitude
	}
	if req.IsActive != nil {
		office.IsActive = *req.IsActive
	}

	if err := s.officeRepo.Update(office); err != nil {
		return nil, errors.New("erreur lors de la mise à jour du siège")
	}

	// Récupérer le siège mis à jour
	updatedOffice, err := s.officeRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération du siège mis à jour")
	}

	return s.officeToDTO(updatedOffice), nil
}

// Delete supprime un siège
func (s *officeService) Delete(id uint) error {
	_, err := s.officeRepo.FindByID(id)
	if err != nil {
		return errors.New("siège introuvable")
	}

	if err := s.officeRepo.Delete(id); err != nil {
		return errors.New("erreur lors de la suppression du siège")
	}

	return nil
}

// GetByCountry récupère les sièges d'un pays
func (s *officeService) GetByCountry(country string) ([]dto.OfficeDTO, error) {
	offices, err := s.officeRepo.FindByCountry(country)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des sièges")
	}

	var officeDTOs []dto.OfficeDTO
	for _, office := range offices {
		officeDTOs = append(officeDTOs, *s.officeToDTO(&office))
	}

	return officeDTOs, nil
}

// GetByCity récupère les sièges d'une ville
func (s *officeService) GetByCity(city string) ([]dto.OfficeDTO, error) {
	offices, err := s.officeRepo.FindByCity(city)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des sièges")
	}

	var officeDTOs []dto.OfficeDTO
	for _, office := range offices {
		officeDTOs = append(officeDTOs, *s.officeToDTO(&office))
	}

	return officeDTOs, nil
}

// officeToDTO convertit un modèle Office en DTO
func (s *officeService) officeToDTO(office *models.Office) *dto.OfficeDTO {
	officeDTO := &dto.OfficeDTO{
		ID:        office.ID,
		Name:      office.Name,
		Code:      office.Code,
		Country:   office.Country,
		City:      office.City,
		Commune:   office.Commune,
		Address:   office.Address,
		FilialeID: office.FilialeID,
		Longitude: office.Longitude,
		Latitude:  office.Latitude,
		IsActive:  office.IsActive,
		CreatedAt: office.CreatedAt.Format(time.RFC3339),
		UpdatedAt: office.UpdatedAt.Format(time.RFC3339),
	}

	// Inclure la filiale si présente
	if office.Filiale != nil {
		officeDTO.Filiale = &dto.FilialeDTO{
			ID:          office.Filiale.ID,
			Code:        office.Filiale.Code,
			Name:        office.Filiale.Name,
			Country:     office.Filiale.Country,
			City:        office.Filiale.City,
			Address:     office.Filiale.Address,
			Phone:       office.Filiale.Phone,
			Email:       office.Filiale.Email,
			IsActive:    office.Filiale.IsActive,
			IsSoftwareProvider: office.Filiale.IsSoftwareProvider,
			CreatedAt:   office.Filiale.CreatedAt,
			UpdatedAt:   office.Filiale.UpdatedAt,
		}
	}

	return officeDTO
}
