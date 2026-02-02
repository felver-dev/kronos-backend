package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"github.com/mcicare/itsm-backend/database"
	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/models"
	"github.com/mcicare/itsm-backend/internal/repositories"
)

// TicketService interface pour les opérations sur les tickets
type TicketService interface {
	Create(req dto.CreateTicketRequest, createdByID uint) (*dto.TicketDTO, error)
	GetByID(id uint, includeDepartment bool) (*dto.TicketDTO, error)
	GetAll(scope interface{}, page, limit int) (*dto.TicketListResponse, error) // scope peut être *scope.QueryScope ou nil
	GetAllWithFilters(scope interface{}, page, limit int, status string, filialeID *uint, assigneeUserID *uint) (*dto.TicketListResponse, error)
	GetByStatus(scope interface{}, status string, page, limit int) (*dto.TicketListResponse, error)
	GetByCategory(scope interface{}, category string, page, limit int, status, priority string) (*dto.TicketListResponse, error)
	GetBySource(scope interface{}, source string, page, limit int) (*dto.TicketListResponse, error)
	GetByAssignedTo(userID uint, page, limit int) (*dto.TicketListResponse, error)
	GetPanier(userID uint, page, limit int) (*dto.TicketListResponse, error) // Panier: tickets assignés à l'utilisateur, non clôturés
	GetByCreatedBy(userID uint, page, limit int) (*dto.TicketListResponse, error)
	GetByUser(userID uint, page, limit int, status string) (*dto.TicketListResponse, error)
	GetByDepartment(departmentID uint, page, limit int) (*dto.TicketListResponse, error)
	GetHistory(ticketID uint) ([]dto.TicketHistoryDTO, error)
	Update(id uint, req dto.UpdateTicketRequest, updatedByID uint) (*dto.TicketDTO, error)
	Assign(id uint, req dto.AssignTicketRequest, assignedByID uint) (*dto.TicketDTO, error)
	ChangeStatus(id uint, status string, changedByID uint) (*dto.TicketDTO, error)
	Close(id uint, closedByID uint) (*dto.TicketDTO, error)
	ValidateTicket(id uint, validatedByID uint) (*dto.TicketDTO, error) // Valider un ticket résolu (le ferme automatiquement)
	Delete(id uint) error
	AddComment(ticketID uint, req dto.CreateTicketCommentRequest, userID uint) (*dto.TicketCommentDTO, error)
	GetComments(ticketID uint, canViewInternalComments bool) ([]dto.TicketCommentDTO, error)
	UpdateComment(ticketID uint, commentID uint, req dto.UpdateTicketCommentRequest, userID uint) (*dto.TicketCommentDTO, error)
	DeleteComment(ticketID uint, commentID uint, userID uint) error
}

// ticketService implémente TicketService
type ticketService struct {
	ticketRepo          repositories.TicketRepository
	userRepo            repositories.UserRepository
	commentRepo         repositories.TicketCommentRepository
	historyRepo         repositories.TicketHistoryRepository
	slaRepo             repositories.SLARepository
	ticketSLARepo       repositories.TicketSLARepository
	ticketCategoryRepo  repositories.TicketCategoryRepository
	notificationRepo    repositories.NotificationRepository
	notificationService NotificationService // Service de notifications pour WebSocket
	departmentRepo      repositories.DepartmentRepository
	filialeRepo         repositories.FilialeRepository
	timeEntryRepo       repositories.TimeEntryRepository // pour valider les entrées de temps quand le ticket est validé
}

// NewTicketService crée une nouvelle instance de TicketService
func NewTicketService(
	ticketRepo repositories.TicketRepository,
	userRepo repositories.UserRepository,
	commentRepo repositories.TicketCommentRepository,
	historyRepo repositories.TicketHistoryRepository,
	slaRepo repositories.SLARepository,
	ticketSLARepo repositories.TicketSLARepository,
	ticketCategoryRepo repositories.TicketCategoryRepository,
	notificationRepo repositories.NotificationRepository,
	notificationService NotificationService,
	departmentRepo repositories.DepartmentRepository,
	filialeRepo repositories.FilialeRepository,
	timeEntryRepo repositories.TimeEntryRepository,
) TicketService {
	return &ticketService{
		ticketRepo:          ticketRepo,
		userRepo:            userRepo,
		commentRepo:         commentRepo,
		historyRepo:         historyRepo,
		slaRepo:             slaRepo,
		ticketSLARepo:       ticketSLARepo,
		ticketCategoryRepo:  ticketCategoryRepo,
		notificationRepo:    notificationRepo,
		notificationService: notificationService,
		departmentRepo:      departmentRepo,
		filialeRepo:         filialeRepo,
		timeEntryRepo:       timeEntryRepo,
	}
}

// Create crée un nouveau ticket
func (s *ticketService) Create(req dto.CreateTicketRequest, createdByID uint) (*dto.TicketDTO, error) {
	// Vérifier que l'utilisateur créateur existe et récupérer sa filiale
	creator, err := s.userRepo.FindByID(createdByID)
	if err != nil {
		return nil, errors.New("utilisateur créateur introuvable")
	}

	// Définir automatiquement filiale_id depuis l'utilisateur créateur si non fourni (user → rôle → département)
	filialeID := req.FilialeID
	if filialeID == nil && creator.FilialeID != nil {
		filialeID = creator.FilialeID
	}
	if filialeID == nil && creator.Role.FilialeID != nil {
		filialeID = creator.Role.FilialeID
	}
	if filialeID == nil && creator.Department != nil && creator.Department.FilialeID != nil {
		filialeID = creator.Department.FilialeID
	}

	// Valider software_id si fourni
	if req.SoftwareID != nil && *req.SoftwareID != 0 {
		// Vérifier que le logiciel existe (sera implémenté avec le repository Software)
		// Pour l'instant, on accepte la valeur
	}

	if err := s.validateCategorySlug(req.Category); err != nil {
		return nil, err
	}

	assigneeIDs, leadID, err := normalizeAssignees(req.AssigneeIDs, req.LeadID)
	if err != nil {
		return nil, err
	}
	// Par défaut : aucun assigné. L'assignation se fera plus tard si besoin.
	// Valider que les utilisateurs assignés appartiennent au même département IT si le créateur est IT
	if err := s.validateAssigneesForITUser(assigneeIDs, createdByID); err != nil {
		return nil, err
	}
	if req.ParentID != nil {
		if *req.ParentID == 0 {
			return nil, errors.New("ticket parent invalide")
		}
		if _, err := s.ticketRepo.FindByID(*req.ParentID); err != nil {
			return nil, errors.New("ticket parent introuvable")
		}
	}

	// Générer le code du ticket (format: TKT-YYYY-NNNN)
	now := time.Now()
	year := now.Year()

	// Générer un code unique avec retry en cas de collision
	// On commence par vérifier tous les codes existants pour cette année
	var code string
	maxRetries := 50 // Augmenter significativement le nombre de tentatives

	// Obtenir le numéro de séquence suggéré
	sequenceNumber, err := s.ticketRepo.GetNextSequenceNumber(year)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la génération du code du ticket: %w", err)
	}

	// Essayer de trouver un code unique en vérifiant systématiquement
	for i := 0; i < maxRetries; i++ {
		code = fmt.Sprintf("TKT-%d-%04d", year, sequenceNumber)

		// Vérifier que le code n'existe pas déjà (y compris les tickets supprimés)
		exists, err := s.ticketRepo.CodeExists(code)
		if err != nil {
			return nil, fmt.Errorf("erreur lors de la vérification du code: %w", err)
		}
		if !exists {
			break // Code unique trouvé
		}
		// Si le code existe, on incrémente et on réessaie
		sequenceNumber++
		if i == maxRetries-1 {
			return nil, fmt.Errorf("impossible de générer un code unique après %d tentatives (dernier code testé: %s, numéro de séquence suggéré: %d)", maxRetries, code, sequenceNumber)
		}
	}

	// Vérifier et valider le requester_id si fourni
	if req.RequesterID != nil && *req.RequesterID != 0 {
		_, err := s.userRepo.FindByID(*req.RequesterID)
		if err != nil {
			return nil, errors.New("utilisateur demandeur introuvable")
		}
	}

	// Déterminer la source : si l'utilisateur n'est pas du département IT de la filiale fournisseur,
	// définir automatiquement "kronos" comme source
	source := req.Source
	isITSupplier, err := s.isUserITOfSupplierFiliale(createdByID)
	if err != nil {
		// En cas d'erreur, utiliser la source fournie (fallback)
		log.Printf("Erreur lors de la vérification du département IT: %v", err)
	} else if !isITSupplier {
		// L'utilisateur n'est pas du département IT de la filiale fournisseur
		// Définir automatiquement "kronos" comme source
		source = "kronos"
	}

	// Créer le ticket
	// S'assurer que Priority a une valeur par défaut si elle est vide
	priority := req.Priority
	if priority == "" {
		priority = "medium" // Valeur par défaut
	}

	ticket := &models.Ticket{
		Code:                code,
		Title:               req.Title,
		Description:         req.Description,
		Category:            req.Category,
		Source:              source,
		Status:              "ouvert", // Statut par défaut
		Priority:            priority,
		CreatedByID:         createdByID,
		RequesterID:         req.RequesterID,
		RequesterName:       req.RequesterName,
		RequesterDepartment: req.RequesterDepartment,
		FilialeID:           filialeID,      // Filiale de l'utilisateur créateur
		SoftwareID:          req.SoftwareID, // Logiciel concerné (optionnel)
		EstimatedTime:       req.EstimatedTime,
		ParentID:            req.ParentID,
	}
	if leadID != nil {
		ticket.AssignedToID = leadID
	} else if len(assigneeIDs) > 0 {
		first := assigneeIDs[0]
		ticket.AssignedToID = &first
	}

	if err := s.ticketRepo.Create(ticket); err != nil {
		// Retourner l'erreur réelle pour le débogage
		return nil, fmt.Errorf("erreur lors de la création du ticket: %w", err)
	}

	// Créer une entrée d'historique
	s.createHistory(ticket.ID, createdByID, "created", "", "", "Ticket créé")

	if len(assigneeIDs) > 0 {
		if err := s.replaceAssignees(ticket.ID, assigneeIDs, leadID); err != nil {
			return nil, err
		}
	}

	// Récupérer le ticket créé avec ses relations
	createdTicket, err := s.ticketRepo.FindByID(ticket.ID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération du ticket créé")
	}

	// Appliquer automatiquement un SLA si une règle correspondante existe
	s.applySLAIfApplicable(createdTicket)

	// Notification : Envoyer une notification à la DSI de MCI CARE CI lors de la création d'un ticket
	// Récupérer les informations du créateur et de la filiale pour le message
	// Note: creator est déjà récupéré au début de la fonction
	filialeName := "une filiale"
	if createdTicket.Filiale != nil {
		filialeName = createdTicket.Filiale.Name
	}
	requesterName := "Un utilisateur"
	if creator != nil {
		if creator.FirstName != "" || creator.LastName != "" {
			requesterName = fmt.Sprintf("%s %s", creator.FirstName, creator.LastName)
		} else {
			requesterName = creator.Username
		}
	} else if createdTicket.RequesterName != "" {
		requesterName = createdTicket.RequesterName
	}

	notificationTitle := fmt.Sprintf("Nouveau ticket : %s", createdTicket.Title)
	notificationMessage := fmt.Sprintf("Un nouveau ticket a été créé par %s (%s). Code: %s", requesterName, filialeName, createdTicket.Code)
	linkURL := fmt.Sprintf("/app/tickets/%d", createdTicket.ID)
	metadata := map[string]any{
		"ticket_id":     createdTicket.ID,
		"ticket_code":   createdTicket.Code,
		"filiale_id":    createdTicket.FilialeID,
		"created_by_id": createdByID,
	}
	s.notifyITDepartmentOfSoftwareProvider("ticket_created", notificationTitle, notificationMessage, linkURL, metadata)

	// Convertir en DTO
	ticketDTO := s.ticketToDTO(createdTicket)
	return &ticketDTO, nil
}

// GetByID récupère un ticket par son ID
func (s *ticketService) GetByID(id uint, includeDepartment bool) (*dto.TicketDTO, error) {
	var (
		ticket *models.Ticket
		err    error
	)
	if includeDepartment {
		ticket, err = s.ticketRepo.FindByID(id)
	} else {
		ticket, err = s.ticketRepo.FindByIDLean(id)
	}
	if err != nil {
		return nil, errors.New("ticket introuvable")
	}

	ticketDTO := s.ticketToDTOWithSubTickets(ticket, true)
	return &ticketDTO, nil
}

// GetAll récupère tous les tickets avec pagination
// Le scope est utilisé pour filtrer automatiquement selon les permissions de l'utilisateur
func (s *ticketService) GetAll(scopeParam interface{}, page, limit int) (*dto.TicketListResponse, error) {
	tickets, total, err := s.ticketRepo.FindAll(scopeParam, page, limit, nil)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des tickets")
	}

	ticketDTOs := make([]dto.TicketDTO, len(tickets))
	for i, ticket := range tickets {
		ticketDTOs[i] = s.ticketToDTO(&ticket)
	}

	return &dto.TicketListResponse{
		Tickets: ticketDTOs,
		Pagination: dto.PaginationDTO{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: calculateTotalPages(total, limit),
		},
	}, nil
}

// GetAllWithFilters récupère les tickets avec filtres optionnels (statut, filiale, assigné)
func (s *ticketService) GetAllWithFilters(scopeParam interface{}, page, limit int, status string, filialeID *uint, assigneeUserID *uint) (*dto.TicketListResponse, error) {
	tickets, total, err := s.ticketRepo.FindWithFilters(scopeParam, page, limit, status, filialeID, assigneeUserID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des tickets")
	}
	ticketDTOs := make([]dto.TicketDTO, len(tickets))
	for i, ticket := range tickets {
		ticketDTOs[i] = s.ticketToDTO(&ticket)
	}
	return &dto.TicketListResponse{
		Tickets: ticketDTOs,
		Pagination: dto.PaginationDTO{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: calculateTotalPages(total, limit),
		},
	}, nil
}

// GetByStatus récupère les tickets par statut
// Le scope est utilisé pour filtrer automatiquement selon les permissions de l'utilisateur
func (s *ticketService) GetByStatus(scopeParam interface{}, status string, page, limit int) (*dto.TicketListResponse, error) {
	tickets, total, err := s.ticketRepo.FindByStatus(scopeParam, status, page, limit)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des tickets")
	}

	ticketDTOs := make([]dto.TicketDTO, len(tickets))
	for i, ticket := range tickets {
		ticketDTOs[i] = s.ticketToDTO(&ticket)
	}

	return &dto.TicketListResponse{
		Tickets: ticketDTOs,
		Pagination: dto.PaginationDTO{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: calculateTotalPages(total, limit),
		},
	}, nil
}

// GetByCategory récupère les tickets par catégorie (avec filtres optionnels)
// Le scope est utilisé pour filtrer automatiquement selon les permissions de l'utilisateur
func (s *ticketService) GetByCategory(scopeParam interface{}, category string, page, limit int, status, priority string) (*dto.TicketListResponse, error) {
	tickets, total, err := s.ticketRepo.FindByCategory(scopeParam, category, page, limit, status, priority)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des tickets")
	}

	ticketDTOs := make([]dto.TicketDTO, len(tickets))
	for i, ticket := range tickets {
		ticketDTOs[i] = s.ticketToDTO(&ticket)
	}

	return &dto.TicketListResponse{
		Tickets: ticketDTOs,
		Pagination: dto.PaginationDTO{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: calculateTotalPages(total, limit),
		},
	}, nil
}

// GetByAssignedTo récupère les tickets assignés à un utilisateur
func (s *ticketService) GetByAssignedTo(userID uint, page, limit int) (*dto.TicketListResponse, error) {
	tickets, total, err := s.ticketRepo.FindByAssignedTo(userID, page, limit)
	if err != nil {
		log.Printf("❌ GetByAssignedTo error: %v", err)
		return nil, fmt.Errorf("erreur lors de la récupération des tickets: %w", err)
	}

	ticketDTOs := make([]dto.TicketDTO, len(tickets))
	for i, ticket := range tickets {
		ticketDTOs[i] = s.ticketToDTO(&ticket)
	}

	return &dto.TicketListResponse{
		Tickets: ticketDTOs,
		Pagination: dto.PaginationDTO{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: calculateTotalPages(total, limit),
		},
	}, nil
}

// GetPanier récupère le panier de l'utilisateur: tickets assignés et non clôturés
func (s *ticketService) GetPanier(userID uint, page, limit int) (*dto.TicketListResponse, error) {
	tickets, total, err := s.ticketRepo.FindPanierByUser(userID, page, limit)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération du panier: %w", err)
	}
	ticketDTOs := make([]dto.TicketDTO, len(tickets))
	for i, ticket := range tickets {
		ticketDTOs[i] = s.ticketToDTO(&ticket)
	}
	return &dto.TicketListResponse{
		Tickets: ticketDTOs,
		Pagination: dto.PaginationDTO{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: calculateTotalPages(total, limit),
		},
	}, nil
}

// GetByCreatedBy récupère les tickets créés par un utilisateur
func (s *ticketService) GetByCreatedBy(userID uint, page, limit int) (*dto.TicketListResponse, error) {
	tickets, total, err := s.ticketRepo.FindByCreatedBy(userID, page, limit)
	if err != nil {
		log.Printf("❌ GetByCreatedBy error: %v", err)
		return nil, fmt.Errorf("erreur lors de la récupération des tickets: %w", err)
	}

	ticketDTOs := make([]dto.TicketDTO, len(tickets))
	for i, ticket := range tickets {
		ticketDTOs[i] = s.ticketToDTO(&ticket)
	}

	return &dto.TicketListResponse{
		Tickets: ticketDTOs,
		Pagination: dto.PaginationDTO{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: calculateTotalPages(total, limit),
		},
	}, nil
}

// GetByUser récupère les tickets créés par l'utilisateur ou qui lui sont assignés
func (s *ticketService) GetByUser(userID uint, page, limit int, status string) (*dto.TicketListResponse, error) {
	tickets, total, err := s.ticketRepo.FindByUser(userID, page, limit, status)
	if err != nil {
		log.Printf("❌ GetByUser error: %v", err)
		return nil, fmt.Errorf("erreur lors de la récupération des tickets: %w", err)
	}

	ticketDTOs := make([]dto.TicketDTO, len(tickets))
	for i, ticket := range tickets {
		ticketDTOs[i] = s.ticketToDTO(&ticket)
	}

	return &dto.TicketListResponse{
		Tickets: ticketDTOs,
		Pagination: dto.PaginationDTO{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: calculateTotalPages(total, limit),
		},
	}, nil
}

// GetBySource récupère les tickets par source
// Le scope est utilisé pour filtrer automatiquement selon les permissions de l'utilisateur
func (s *ticketService) GetBySource(scopeParam interface{}, source string, page, limit int) (*dto.TicketListResponse, error) {
	tickets, total, err := s.ticketRepo.FindBySource(scopeParam, source, page, limit)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des tickets")
	}

	ticketDTOs := make([]dto.TicketDTO, len(tickets))
	for i, ticket := range tickets {
		ticketDTOs[i] = s.ticketToDTO(&ticket)
	}

	return &dto.TicketListResponse{
		Tickets: ticketDTOs,
		Pagination: dto.PaginationDTO{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: calculateTotalPages(total, limit),
		},
	}, nil
}

// GetByDepartment récupère les tickets par département du demandeur
func (s *ticketService) GetByDepartment(departmentID uint, page, limit int) (*dto.TicketListResponse, error) {
	tickets, total, err := s.ticketRepo.FindByDepartment(departmentID, page, limit)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération des tickets par département: %w", err)
	}

	ticketDTOs := make([]dto.TicketDTO, len(tickets))
	for i, ticket := range tickets {
		ticketDTOs[i] = s.ticketToDTO(&ticket)
	}

	return &dto.TicketListResponse{
		Tickets: ticketDTOs,
		Pagination: dto.PaginationDTO{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: calculateTotalPages(total, limit),
		},
	}, nil
}

// GetHistory récupère l'historique d'un ticket
func (s *ticketService) GetHistory(ticketID uint) ([]dto.TicketHistoryDTO, error) {
	histories, err := s.historyRepo.FindByTicketID(ticketID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération de l'historique")
	}

	historyDTOs := make([]dto.TicketHistoryDTO, len(histories))
	for i, history := range histories {
		userDTO := s.userToDTO(&history.User)
		historyDTOs[i] = dto.TicketHistoryDTO{
			ID:          history.ID,
			TicketID:    history.TicketID,
			User:        userDTO,
			Action:      history.Action,
			FieldName:   history.FieldName,
			OldValue:    history.OldValue,
			NewValue:    history.NewValue,
			Description: history.Description,
			CreatedAt:   history.CreatedAt,
		}
	}

	return historyDTOs, nil
}

// Update met à jour un ticket
func (s *ticketService) Update(id uint, req dto.UpdateTicketRequest, updatedByID uint) (*dto.TicketDTO, error) {
	// Récupérer le ticket existant
	start := time.Now()
	ticket, err := s.ticketRepo.FindByIDForUpdate(id)
	if err != nil {
		return nil, errors.New("ticket introuvable")
	}

	updates := map[string]interface{}{}

	// Mettre à jour les champs fournis
	if req.Title != "" {
		s.createHistory(id, updatedByID, "updated", "title", ticket.Title, req.Title)
		ticket.Title = req.Title
		updates["title"] = req.Title
	}

	if req.Description != "" {
		s.createHistory(id, updatedByID, "updated", "description", ticket.Description, req.Description)
		ticket.Description = req.Description
		updates["description"] = req.Description
	}

	// Mettre à jour la catégorie si elle est fournie
	if req.Category != "" {
		if err := s.validateCategorySlug(req.Category); err != nil {
			return nil, err
		}
		oldCategory := ticket.Category
		if oldCategory != req.Category {
			s.createHistory(id, updatedByID, "updated", "category", oldCategory, req.Category)
		}
		ticket.Category = req.Category
		updates["category"] = req.Category
	}

	if req.Status != "" {
		s.createHistory(id, updatedByID, "updated", "status", ticket.Status, req.Status)
		ticket.Status = req.Status
		updates["status"] = req.Status
	}

	if req.Priority != "" {
		s.createHistory(id, updatedByID, "updated", "priority", ticket.Priority, req.Priority)
		ticket.Priority = req.Priority
		updates["priority"] = req.Priority
	}

	// Gérer RequesterID (prioritaire sur RequesterName)
	if req.RequesterID != nil {
		// Si RequesterID est fourni, vérifier qu'il existe
		if *req.RequesterID != 0 {
			_, err := s.userRepo.FindByID(*req.RequesterID)
			if err != nil {
				return nil, errors.New("utilisateur demandeur introuvable")
			}
		}
		oldRequesterID := ticket.RequesterID
		if oldRequesterID == nil || *oldRequesterID != *req.RequesterID {
			oldValue := ""
			if oldRequesterID != nil {
				oldValue = fmt.Sprintf("%d", *oldRequesterID)
			}
			newValue := ""
			if *req.RequesterID != 0 {
				newValue = fmt.Sprintf("%d", *req.RequesterID)
			}
			s.createHistory(id, updatedByID, "updated", "requester_id", oldValue, newValue)
			ticket.RequesterID = req.RequesterID
			updates["requester_id"] = req.RequesterID
			// Si RequesterID est défini, on peut aussi mettre à jour RequesterName depuis l'utilisateur
			if *req.RequesterID != 0 {
				requesterUser, _ := s.userRepo.FindByID(*req.RequesterID)
				if requesterUser != nil {
					requesterName := fmt.Sprintf("%s %s", requesterUser.FirstName, requesterUser.LastName)
					if requesterUser.FirstName == "" && requesterUser.LastName == "" {
						requesterName = requesterUser.Username
					}
					ticket.RequesterName = requesterName
					updates["requester_name"] = requesterName
				}
			} else {
				// Si RequesterID est mis à 0/null, on garde RequesterName tel quel
			}
		}
	} else if req.RequesterName != "" {
		// Si RequesterID n'est pas fourni mais RequesterName oui, on met à jour seulement le nom
		s.createHistory(id, updatedByID, "updated", "requester_name", ticket.RequesterName, req.RequesterName)
		ticket.RequesterName = req.RequesterName
		updates["requester_name"] = req.RequesterName
	}

	if req.RequesterDepartment != "" {
		s.createHistory(id, updatedByID, "updated", "requester_department", ticket.RequesterDepartment, req.RequesterDepartment)
		ticket.RequesterDepartment = req.RequesterDepartment
		updates["requester_department"] = req.RequesterDepartment
	}

	// Gérer SoftwareID
	if req.SoftwareID != nil {
		oldSoftwareID := ticket.SoftwareID
		if oldSoftwareID == nil || (req.SoftwareID != nil && *oldSoftwareID != *req.SoftwareID) {
			oldValue := ""
			if oldSoftwareID != nil {
				oldValue = fmt.Sprintf("%d", *oldSoftwareID)
			}
			newValue := ""
			if *req.SoftwareID != 0 {
				newValue = fmt.Sprintf("%d", *req.SoftwareID)
				// Vérifier que le logiciel existe (sera implémenté avec le repository Software)
			}
			s.createHistory(id, updatedByID, "updated", "software_id", oldValue, newValue)
			ticket.SoftwareID = req.SoftwareID
			updates["software_id"] = req.SoftwareID
		}
	}

	if req.ParentID != nil {
		if *req.ParentID == 0 {
			return nil, errors.New("ticket parent invalide")
		}
		if *req.ParentID == ticket.ID {
			return nil, errors.New("un ticket ne peut pas être son propre parent")
		}
		if _, err := s.ticketRepo.FindByID(*req.ParentID); err != nil {
			return nil, errors.New("ticket parent introuvable")
		}
		ticket.ParentID = req.ParentID
		updates["parent_id"] = req.ParentID
	}

	assigneesStart := time.Now()
	if len(req.AssigneeIDs) > 0 || req.LeadID != nil {
		assigneeIDs, leadID, err := normalizeAssignees(req.AssigneeIDs, req.LeadID)
		if err != nil {
			return nil, err
		}
		// Valider que les utilisateurs assignés appartiennent au même département IT si l'assigneur est IT
		if err := s.validateAssigneesForITUser(assigneeIDs, updatedByID); err != nil {
			return nil, err
		}
		if leadID != nil {
			ticket.AssignedToID = leadID
		} else if len(assigneeIDs) > 0 {
			first := assigneeIDs[0]
			ticket.AssignedToID = &first
		}
		updates["assigned_to_id"] = ticket.AssignedToID
		if err := s.replaceAssignees(ticket.ID, assigneeIDs, leadID); err != nil {
			return nil, err
		}
	}
	assigneesDur := time.Since(assigneesStart)

	// Temps estimé (résolveurs)
	if req.EstimatedTime != nil {
		oldVal := ""
		if ticket.EstimatedTime != nil {
			oldVal = fmt.Sprintf("%d", *ticket.EstimatedTime)
		}
		newVal := fmt.Sprintf("%d", *req.EstimatedTime)
		s.createHistory(id, updatedByID, "updated", "estimated_time", oldVal, newVal)
		ticket.EstimatedTime = req.EstimatedTime
		updates["estimated_time"] = req.EstimatedTime
	}

	// Sauvegarder
	fmt.Printf("DEBUG: Avant sauvegarde - Catégorie du ticket: '%s'\n", ticket.Category)
	updateStart := time.Now()
	if err := s.ticketRepo.UpdateFields(ticket.ID, updates); err != nil {
		fmt.Printf("DEBUG: Erreur lors de la sauvegarde: %v\n", err)
		return nil, errors.New("erreur lors de la mise à jour du ticket")
	}
	fmt.Printf("DEBUG: Sauvegarde réussie\n")
	updateDur := time.Since(updateStart)

	ticketDTO := s.ticketToDTO(ticket)
	log.Printf("PERF Update ticket=%d assignees=%s update=%s total=%s", id, assigneesDur, updateDur, time.Since(start))
	return &ticketDTO, nil
}

// Assign assigne un ticket à un utilisateur
func (s *ticketService) Assign(id uint, req dto.AssignTicketRequest, assignedByID uint) (*dto.TicketDTO, error) {
	start := time.Now()
	// Récupérer le ticket (léger)
	ticket, err := s.ticketRepo.FindByIDForAssign(id)
	if err != nil {
		return nil, errors.New("ticket introuvable")
	}

	assigneeIDs := req.UserIDs
	if len(assigneeIDs) == 0 && req.UserID != 0 {
		assigneeIDs = []uint{req.UserID}
	}
	assigneeIDs, leadID, err := normalizeAssignees(assigneeIDs, req.LeadID)
	if err != nil {
		return nil, err
	}
	if len(assigneeIDs) == 0 {
		return nil, errors.New("aucun utilisateur assigné")
	}
	validateStart := time.Now()
	// Valider que les utilisateurs assignés appartiennent au même département IT si l'assigneur est IT
	if err := s.validateAssigneesForITUser(assigneeIDs, assignedByID); err != nil {
		return nil, err
	}
	validateDur := time.Since(validateStart)

	// Enregistrer l'ancien assigné pour l'historique
	oldAssignedID := ticket.AssignedToID
	var newAssignedID *uint
	if leadID != nil {
		newAssignedID = leadID
	} else {
		first := assigneeIDs[0]
		newAssignedID = &first
	}

	// Changer le statut si assigné
	status := ticket.Status
	if status == "ouvert" {
		status = "en_cours"
	}

	// Sauvegarder (update léger)
	updateStart := time.Now()
	if err := s.ticketRepo.UpdateAssignFields(ticket.ID, newAssignedID, req.EstimatedTime, status); err != nil {
		return nil, errors.New("erreur lors de l'assignation du ticket")
	}
	updateDur := time.Since(updateStart)

	replaceStart := time.Now()
	if err := s.replaceAssignees(ticket.ID, assigneeIDs, leadID); err != nil {
		return nil, err
	}
	replaceDur := time.Since(replaceStart)

	// Créer une entrée d'historique
	oldValue := ""
	newValue := ""
	if oldAssignedID != nil {
		oldValue = fmt.Sprintf("user#%d", *oldAssignedID)
	}
	if newAssignedID != nil {
		newValue = fmt.Sprintf("user#%d", *newAssignedID)
	}
	s.createHistory(id, assignedByID, "assigned", "assigned_to", oldValue, newValue)

	// Récupérer le ticket mis à jour
	fetchStart := time.Now()
	updatedTicket, err := s.ticketRepo.FindByIDLean(id)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération du ticket mis à jour")
	}
	fetchDur := time.Since(fetchStart)

	ticketDTO := s.ticketToDTO(updatedTicket)
	fmt.Printf("PERF AssignService ticket=%d users=%d validate=%s update=%s replace=%s fetch=%s total=%s\n",
		id, len(assigneeIDs), validateDur, updateDur, replaceDur, fetchDur, time.Since(start))
	return &ticketDTO, nil
}

// ChangeStatus change le statut d'un ticket
func (s *ticketService) ChangeStatus(id uint, status string, changedByID uint) (*dto.TicketDTO, error) {
	// Récupérer le ticket
	ticket, err := s.ticketRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("ticket introuvable")
	}

	// Valider le statut (ajouter "resolu" pour le workflow multi-filiales)
	validStatuses := []string{"ouvert", "en_cours", "en_attente", "resolu", "cloture"}
	valid := false
	for _, vs := range validStatuses {
		if status == vs {
			valid = true
			break
		}
	}
	if !valid {
		return nil, errors.New("statut invalide")
	}

	// Enregistrer l'ancien statut pour l'historique
	oldStatus := ticket.Status
	ticket.Status = status

	// Si le ticket est clôturé, enregistrer la date de clôture
	if status == "cloture" && ticket.ClosedAt == nil {
		now := time.Now()
		ticket.ClosedAt = &now
	}

	// Sauvegarder
	if err := s.ticketRepo.Update(ticket); err != nil {
		return nil, errors.New("erreur lors du changement de statut")
	}

	// Créer une entrée d'historique
	s.createHistory(id, changedByID, "status_changed", "status", oldStatus, status)

	// Notification : Si le ticket est soumis pour validation (en_attente), notifier le demandeur
	if status == "en_attente" && oldStatus != "en_attente" {
		ticketWithRelations, err := s.ticketRepo.FindByID(id)
		if err == nil {
			var requesterID uint
			if ticketWithRelations.RequesterID != nil && *ticketWithRelations.RequesterID != 0 {
				requesterID = *ticketWithRelations.RequesterID
			} else if ticketWithRelations.CreatedByID != 0 {
				requesterID = ticketWithRelations.CreatedByID
			}
			if requesterID != 0 {
				resolver, _ := s.userRepo.FindByID(changedByID)
				resolverName := "L'équipe IT"
				if resolver != nil {
					if resolver.FirstName != "" || resolver.LastName != "" {
						resolverName = fmt.Sprintf("%s %s", resolver.FirstName, resolver.LastName)
					} else {
						resolverName = resolver.Username
					}
				}
				notificationTitle := fmt.Sprintf("Ticket soumis pour validation : %s", ticketWithRelations.Title)
				notificationMessage := fmt.Sprintf("Votre ticket %s a été traité par %s et est en attente de votre validation. Veuillez valider si le problème est réglé, ou invalider le ticket si ce n'est pas le cas.", ticketWithRelations.Code, resolverName)
				linkURL := fmt.Sprintf("/app/tickets/%d", id)
				metadata := map[string]any{
					"ticket_id":      id,
					"ticket_code":    ticketWithRelations.Code,
					"resolved_by_id": changedByID,
				}
				s.createNotification(requesterID, "ticket_submitted_for_validation", notificationTitle, notificationMessage, linkURL, metadata)
			}
		}
	}

	// Notification : Si le ticket passe de "resolu" à un autre statut (invalidation), notifier la DSI MCI CARE CI
	if oldStatus == "resolu" && status != "resolu" && status != "cloture" {
		ticketWithRelations, err := s.ticketRepo.FindByID(id)
		if err == nil {
			// Récupérer les informations du validateur (celui qui invalide)
			invalidator, _ := s.userRepo.FindByID(changedByID)
			invalidatorName := "Un utilisateur"
			if invalidator != nil {
				if invalidator.FirstName != "" || invalidator.LastName != "" {
					invalidatorName = fmt.Sprintf("%s %s", invalidator.FirstName, invalidator.LastName)
				} else {
					invalidatorName = invalidator.Username
				}
			}

			filialeName := "une filiale"
			if ticketWithRelations.Filiale != nil {
				filialeName = ticketWithRelations.Filiale.Name
			}

			notificationTitle := fmt.Sprintf("Ticket invalidé : %s", ticketWithRelations.Title)
			notificationMessage := fmt.Sprintf("Le ticket %s (%s) a été invalidé par %s. Le ticket nécessite une nouvelle résolution.", ticketWithRelations.Code, filialeName, invalidatorName)
			linkURL := fmt.Sprintf("/app/tickets/%d", id)
			metadata := map[string]any{
				"ticket_id":         id,
				"ticket_code":       ticketWithRelations.Code,
				"invalidated_by_id": changedByID,
				"new_status":        status,
			}
			s.notifyITDepartmentOfSoftwareProvider("ticket_invalidated", notificationTitle, notificationMessage, linkURL, metadata)
		}
	}

	// Récupérer le ticket mis à jour
	updatedTicket, err := s.ticketRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération du ticket mis à jour")
	}

	ticketDTO := s.ticketToDTO(updatedTicket)
	return &ticketDTO, nil
}

// ValidateTicket valide un ticket soumis pour validation (en_attente) et le passe à "resolu"
// Seuls les utilisateurs avec tickets.validate peuvent valider
// Le ticket doit être en statut "en_attente" pour être validé
func (s *ticketService) ValidateTicket(id uint, validatedByID uint) (*dto.TicketDTO, error) {
	// Récupérer le ticket
	ticket, err := s.ticketRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("ticket introuvable")
	}

	// Vérifier que le ticket est en statut "en_attente" (soumis pour validation)
	if ticket.Status != "en_attente" {
		return nil, errors.New("seuls les tickets en attente de validation peuvent être validés")
	}

	// Vérifier que l'utilisateur validateur existe
	_, err = s.userRepo.FindByID(validatedByID)
	if err != nil {
		return nil, errors.New("utilisateur validateur introuvable")
	}

	// Note: La vérification de permission (tickets.validate, tickets.validate_own)
	// et la vérification si l'utilisateur est le créateur sont faites dans le handler

	// Marquer le ticket comme validé : statut passe à "resolu" (la fermeture se fait via le bouton Fermer le ticket)
	oldStatus := ticket.Status
	now := time.Now()
	ticket.ValidatedByUserID = &validatedByID
	ticket.ValidatedAt = &now
	ticket.Status = "resolu"

	// Sauvegarder
	if err := s.ticketRepo.Update(ticket); err != nil {
		return nil, errors.New("erreur lors de la validation du ticket")
	}

	// Créer une entrée d'historique
	s.createHistory(id, validatedByID, "validated", "status", oldStatus, "resolu")
	s.createHistory(id, validatedByID, "status_changed", "status", oldStatus, "resolu")

	// Mettre à jour le statut SLA si le ticket a un SLA associé
	s.updateSLAOnClose(id)

	// Valider automatiquement toutes les entrées de temps liées au ticket (cohérence Gestion du temps)
	if err := s.timeEntryRepo.ValidateByTicketID(id, validatedByID); err != nil {
		log.Printf("⚠️ Validation du ticket %d : erreur lors de la validation des entrées de temps : %v", id, err)
	}

	// Notification : Lors de la validation d'un ticket, notifier la DSI de MCI CARE CI
	ticketWithRelations, err := s.ticketRepo.FindByID(id)
	if err == nil {
		// Récupérer les informations du validateur
		validator, _ := s.userRepo.FindByID(validatedByID)
		validatorName := "Un utilisateur"
		if validator != nil {
			if validator.FirstName != "" || validator.LastName != "" {
				validatorName = fmt.Sprintf("%s %s", validator.FirstName, validator.LastName)
			} else {
				validatorName = validator.Username
			}
		}

		filialeName := "une filiale"
		if ticketWithRelations.Filiale != nil {
			filialeName = ticketWithRelations.Filiale.Name
		}

		notificationTitle := fmt.Sprintf("Ticket validé : %s", ticketWithRelations.Title)
		notificationMessage := fmt.Sprintf("Le ticket %s (%s) a été validé et clôturé par %s.", ticketWithRelations.Code, filialeName, validatorName)
		linkURL := fmt.Sprintf("/app/tickets/%d", id)
		metadata := map[string]any{
			"ticket_id":       id,
			"ticket_code":     ticketWithRelations.Code,
			"validated_by_id": validatedByID,
		}
		s.notifyITDepartmentOfSoftwareProvider("ticket_validated", notificationTitle, notificationMessage, linkURL, metadata)

		// Notifier aussi le créateur du ticket (demandeur) pour qu'il ait la confirmation
		creatorID := ticketWithRelations.CreatedByID
		if creatorID != 0 {
			creatorTitle := fmt.Sprintf("Votre ticket a été validé : %s", ticketWithRelations.Title)
			creatorMessage := fmt.Sprintf("Le ticket %s a été validé. Le problème est considéré comme résolu.", ticketWithRelations.Code)
			s.createNotification(creatorID, "ticket_validated", creatorTitle, creatorMessage, linkURL, metadata)
		}
	}

	// Récupérer le ticket mis à jour
	updatedTicket, err := s.ticketRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération du ticket mis à jour")
	}

	ticketDTO := s.ticketToDTO(updatedTicket)
	return &ticketDTO, nil
}

// Close ferme un ticket
func (s *ticketService) Close(id uint, closedByID uint) (*dto.TicketDTO, error) {
	ticketDTO, err := s.ChangeStatus(id, "cloture", closedByID)
	if err != nil {
		return nil, err
	}

	// Mettre à jour le statut SLA si le ticket a un SLA associé
	s.updateSLAOnClose(id)

	return ticketDTO, nil
}

// Delete supprime un ticket (soft delete)
func (s *ticketService) Delete(id uint) error {
	// Vérifier que le ticket existe
	_, err := s.ticketRepo.FindByID(id)
	if err != nil {
		return errors.New("ticket introuvable")
	}

	// Supprimer (soft delete)
	if err := s.ticketRepo.Delete(id); err != nil {
		return errors.New("erreur lors de la suppression du ticket")
	}

	return nil
}

// AddComment ajoute un commentaire à un ticket
func (s *ticketService) AddComment(ticketID uint, req dto.CreateTicketCommentRequest, userID uint) (*dto.TicketCommentDTO, error) {
	// Vérifier que le ticket existe (requête légère)
	exists, err := s.ticketRepo.ExistsByID(ticketID)
	if err != nil {
		return nil, errors.New("erreur lors de la vérification du ticket")
	}
	if !exists {
		return nil, errors.New("ticket introuvable")
	}

	// Créer le commentaire
	comment := &models.TicketComment{
		TicketID:   ticketID,
		UserID:     userID,
		Comment:    req.Comment,
		IsInternal: req.IsInternal,
	}

	if err := s.commentRepo.Create(comment); err != nil {
		return nil, errors.New("erreur lors de la création du commentaire")
	}

	// Créer une entrée d'historique
	s.createHistory(ticketID, userID, "comment_added", "", "", "Commentaire ajouté")

	// Récupérer le commentaire créé
	createdComment, err := s.commentRepo.FindByIDWithUser(comment.ID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération du commentaire créé")
	}

	// Convertir en DTO
	commentDTO := s.commentToDTO(createdComment)
	return &commentDTO, nil
}

// GetComments récupère tous les commentaires d'un ticket.
// Si canViewInternalComments est false, les commentaires internes sont exclus (visibles uniquement par l'IT).
func (s *ticketService) GetComments(ticketID uint, canViewInternalComments bool) ([]dto.TicketCommentDTO, error) {
	comments, err := s.commentRepo.FindByTicketID(ticketID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des commentaires")
	}

	commentDTOs := make([]dto.TicketCommentDTO, 0, len(comments))
	for i := range comments {
		dto := s.commentToDTO(&comments[i])
		if dto.IsInternal && !canViewInternalComments {
			continue
		}
		commentDTOs = append(commentDTOs, dto)
	}

	return commentDTOs, nil
}

// UpdateComment met à jour un commentaire. Seul l'auteur du commentaire peut le modifier.
func (s *ticketService) UpdateComment(ticketID uint, commentID uint, req dto.UpdateTicketCommentRequest, userID uint) (*dto.TicketCommentDTO, error) {
	comment, err := s.commentRepo.FindByIDWithUser(commentID)
	if err != nil || comment == nil {
		return nil, errors.New("commentaire introuvable")
	}
	if comment.TicketID != ticketID {
		return nil, errors.New("commentaire introuvable pour ce ticket")
	}
	if comment.UserID != userID {
		return nil, errors.New("seul l'auteur du commentaire peut le modifier")
	}
	comment.Comment = strings.TrimSpace(req.Comment)
	if comment.Comment == "" {
		return nil, errors.New("le commentaire ne peut pas être vide")
	}
	if err := s.commentRepo.Update(comment); err != nil {
		return nil, errors.New("erreur lors de la mise à jour du commentaire")
	}
	updated, err := s.commentRepo.FindByIDWithUser(commentID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération du commentaire")
	}
	dto := s.commentToDTO(updated)
	return &dto, nil
}

// DeleteComment supprime un commentaire (soft delete). Seul l'auteur du commentaire peut le supprimer.
func (s *ticketService) DeleteComment(ticketID uint, commentID uint, userID uint) error {
	comment, err := s.commentRepo.FindByID(commentID)
	if err != nil || comment == nil {
		return errors.New("commentaire introuvable")
	}
	if comment.TicketID != ticketID {
		return errors.New("commentaire introuvable pour ce ticket")
	}
	if comment.UserID != userID {
		return errors.New("seul l'auteur du commentaire peut le supprimer")
	}
	if err := s.commentRepo.Delete(commentID); err != nil {
		return errors.New("erreur lors de la suppression du commentaire")
	}
	return nil
}

func calculateTotalPages(total int64, limit int) int {
	if limit <= 0 {
		return 0
	}
	if total == 0 {
		return 0
	}
	return int(math.Ceil(float64(total) / float64(limit)))
}

func normalizeAssignees(assigneeIDs []uint, leadID *uint) ([]uint, *uint, error) {
	unique := make([]uint, 0, len(assigneeIDs))
	seen := map[uint]bool{}
	for _, id := range assigneeIDs {
		if id == 0 || seen[id] {
			continue
		}
		seen[id] = true
		unique = append(unique, id)
	}

	if leadID != nil {
		if *leadID == 0 {
			return nil, nil, errors.New("responsable invalide")
		}
		if len(unique) == 0 {
			unique = []uint{*leadID}
			return unique, leadID, nil
		}
		if !seen[*leadID] {
			return nil, nil, errors.New("le responsable doit faire partie des assignés")
		}
	}

	return unique, leadID, nil
}

func (s *ticketService) validateAssignees(assigneeIDs []uint) error {
	count, err := s.userRepo.CountByIDs(assigneeIDs)
	if err != nil {
		return errors.New("erreur lors de la vérification des assignés")
	}
	if int(count) != len(assigneeIDs) {
		return errors.New("utilisateur assigné introuvable")
	}
	return nil
}

// validateAssigneesForITUser valide que les utilisateurs assignés appartiennent au même département IT que l'utilisateur qui assigne
func (s *ticketService) validateAssigneesForITUser(assigneeIDs []uint, assignedByID uint) error {
	// Vérifier que l'utilisateur qui assigne est IT de la filiale fournisseur
	isITSupplier, err := s.isUserITOfSupplierFiliale(assignedByID)
	if err != nil {
		return fmt.Errorf("erreur lors de la vérification du département IT: %w", err)
	}

	// Si l'utilisateur n'est pas IT de la filiale fournisseur, pas de restriction
	if !isITSupplier {
		return s.validateAssignees(assigneeIDs)
	}

	// Récupérer le département de l'utilisateur qui assigne
	assigner, err := s.userRepo.FindByID(assignedByID)
	if err != nil {
		return errors.New("utilisateur assigneur introuvable")
	}

	if assigner.DepartmentID == nil {
		return errors.New("l'utilisateur assigneur n'appartient à aucun département")
	}

	assignerDeptID := *assigner.DepartmentID

	// Vérifier que tous les utilisateurs assignés appartiennent au même département
	for _, assigneeID := range assigneeIDs {
		assignee, err := s.userRepo.FindByID(assigneeID)
		if err != nil {
			return fmt.Errorf("utilisateur assigné (ID: %d) introuvable", assigneeID)
		}

		if assignee.DepartmentID == nil {
			return fmt.Errorf("l'utilisateur assigné (ID: %d) n'appartient à aucun département", assigneeID)
		}

		if *assignee.DepartmentID != assignerDeptID {
			return fmt.Errorf("l'utilisateur assigné (ID: %d) n'appartient pas au même département IT que vous", assigneeID)
		}
	}

	return nil
}

// validateCategorySlug vérifie que le slug existe dans ticket_categories et que la catégorie est active.
func (s *ticketService) validateCategorySlug(slug string) error {
	if slug == "" {
		return nil
	}
	cat, err := s.ticketCategoryRepo.FindBySlug(slug)
	if err != nil || cat == nil {
		return fmt.Errorf("catégorie inconnue : %q (utilisez un slug de catégorie existant, ex. incident, demande, changement)", slug)
	}
	if !cat.IsActive {
		return fmt.Errorf("la catégorie %q n'est plus active", slug)
	}
	return nil
}

func (s *ticketService) replaceAssignees(ticketID uint, assigneeIDs []uint, leadID *uint) error {
	if err := database.DB.Where("ticket_id = ?", ticketID).Delete(&models.TicketAssignee{}).Error; err != nil {
		return errors.New("erreur lors de la mise à jour des assignations")
	}
	if len(assigneeIDs) == 0 {
		return nil
	}

	assignees := make([]models.TicketAssignee, 0, len(assigneeIDs))
	for _, id := range assigneeIDs {
		isLead := leadID != nil && *leadID == id
		assignees = append(assignees, models.TicketAssignee{
			TicketID: ticketID,
			UserID:   id,
			IsLead:   isLead,
		})
	}

	if err := database.DB.Create(&assignees).Error; err != nil {
		return errors.New("erreur lors de la création des assignations")
	}
	return nil
}

// createHistory crée une entrée d'historique pour un ticket
func (s *ticketService) createHistory(ticketID, userID uint, action, fieldName, oldValue, newValue string) {
	history := &models.TicketHistory{
		TicketID:    ticketID,
		UserID:      userID,
		Action:      action,
		FieldName:   fieldName,
		OldValue:    oldValue,
		NewValue:    newValue,
		Description: "",
	}
	go func(h *models.TicketHistory) {
		if err := s.historyRepo.Create(h); err != nil {
			log.Printf("WARN history create ticket=%d action=%s err=%v", ticketID, action, err)
		}
	}(history)
}

// ticketToDTO convertit un modèle Ticket en DTO TicketDTO
func (s *ticketService) ticketToDTO(ticket *models.Ticket) dto.TicketDTO {
	return s.ticketToDTOWithSubTickets(ticket, false)
}

func (s *ticketService) ticketToDTOWithSubTickets(ticket *models.Ticket, includeSubTickets bool) dto.TicketDTO {
	// Convertir les utilisateurs en DTOs
	var assignedToDTO *dto.UserDTO
	if ticket.AssignedTo != nil {
		assignedDTO := s.userToDTO(ticket.AssignedTo)
		assignedToDTO = &assignedDTO
	}

	assigneesDTO := make([]dto.TicketAssigneeDTO, 0, len(ticket.Assignees))
	var leadDTO *dto.UserDTO
	for _, assignee := range ticket.Assignees {
		userDTO := s.userToDTO(&assignee.User)
		assigneesDTO = append(assigneesDTO, dto.TicketAssigneeDTO{
			User:   userDTO,
			IsLead: assignee.IsLead,
		})
		if assignee.IsLead {
			leadCopy := userDTO
			leadDTO = &leadCopy
		}
	}

	// CreatedBy : utiliser le Preload si chargé, sinon charger manuellement pour éviter "Utilisateur inconnu"
	var createdByDTO dto.UserDTO
	if ticket.CreatedBy.ID != 0 {
		createdByDTO = s.userToDTO(&ticket.CreatedBy)
	} else if ticket.CreatedByID != 0 {
		createdByUser, err := s.userRepo.FindByID(ticket.CreatedByID)
		if err == nil && createdByUser != nil {
			createdByDTO = s.userToDTO(createdByUser)
		} else {
			createdByDTO = dto.UserDTO{ID: ticket.CreatedByID, Username: "Utilisateur inconnu"}
		}
	} else {
		createdByDTO = dto.UserDTO{Username: "Utilisateur inconnu"}
	}

	// Gérer le Requester (prioritaire sur RequesterName)
	var requesterDTO *dto.UserDTO
	var requesterName string
	if ticket.Requester != nil && ticket.Requester.ID != 0 {
		// Si la relation Requester est chargée, l'utiliser
		reqDTO := s.userToDTO(ticket.Requester)
		requesterDTO = &reqDTO
		// Construire le nom depuis l'utilisateur
		requesterName = fmt.Sprintf("%s %s", ticket.Requester.FirstName, ticket.Requester.LastName)
		if ticket.Requester.FirstName == "" && ticket.Requester.LastName == "" {
			requesterName = ticket.Requester.Username
		}
	} else if ticket.RequesterID != nil && *ticket.RequesterID != 0 {
		// Si RequesterID est défini mais la relation n'est pas chargée, charger l'utilisateur
		requesterUser, err := s.userRepo.FindByID(*ticket.RequesterID)
		if err == nil && requesterUser != nil {
			reqDTO := s.userToDTO(requesterUser)
			requesterDTO = &reqDTO
			requesterName = fmt.Sprintf("%s %s", requesterUser.FirstName, requesterUser.LastName)
			if requesterUser.FirstName == "" && requesterUser.LastName == "" {
				requesterName = requesterUser.Username
			}
		} else {
			// Fallback sur RequesterName si l'utilisateur n'existe plus
			requesterName = ticket.RequesterName
		}
	} else {
		// Fallback sur RequesterName pour les demandeurs externes ou anciens tickets
		requesterName = ticket.RequesterName
	}

	var subTickets []dto.TicketDTO
	if includeSubTickets && len(ticket.SubTickets) > 0 {
		subTickets = make([]dto.TicketDTO, 0, len(ticket.SubTickets))
		for _, sub := range ticket.SubTickets {
			subTickets = append(subTickets, s.ticketToDTOWithSubTickets(&sub, false))
		}
	}

	// Gérer ValidatedBy
	var validatedByDTO *dto.UserDTO
	if ticket.ValidatedBy != nil && ticket.ValidatedBy.ID != 0 {
		validDTO := s.userToDTO(ticket.ValidatedBy)
		validatedByDTO = &validDTO
	} else if ticket.ValidatedByUserID != nil && *ticket.ValidatedByUserID != 0 {
		validatedUser, err := s.userRepo.FindByID(*ticket.ValidatedByUserID)
		if err == nil && validatedUser != nil {
			validDTO := s.userToDTO(validatedUser)
			validatedByDTO = &validDTO
		}
	}

	// Convertir Filiale en DTO si présent
	var filialeDTO *dto.FilialeDTO
	if ticket.Filiale != nil && ticket.Filiale.ID != 0 {
		filialeDTO = &dto.FilialeDTO{
			ID:                 ticket.Filiale.ID,
			Code:               ticket.Filiale.Code,
			Name:               ticket.Filiale.Name,
			Country:            ticket.Filiale.Country,
			City:               ticket.Filiale.City,
			Address:            ticket.Filiale.Address,
			Phone:              ticket.Filiale.Phone,
			Email:              ticket.Filiale.Email,
			IsActive:           ticket.Filiale.IsActive,
			IsSoftwareProvider: ticket.Filiale.IsSoftwareProvider,
			CreatedAt:          ticket.Filiale.CreatedAt,
			UpdatedAt:          ticket.Filiale.UpdatedAt,
		}
	}

	// Convertir Software en DTO si présent
	var softwareDTO *dto.SoftwareDTO
	if ticket.Software != nil && ticket.Software.ID != 0 {
		softwareDTO = &dto.SoftwareDTO{
			ID:          ticket.Software.ID,
			Code:        ticket.Software.Code,
			Name:        ticket.Software.Name,
			Description: ticket.Software.Description,
			Version:     ticket.Software.Version,
			IsActive:    ticket.Software.IsActive,
			CreatedAt:   ticket.Software.CreatedAt,
			UpdatedAt:   ticket.Software.UpdatedAt,
		}
	}

	return dto.TicketDTO{
		ID:                  ticket.ID,
		Code:                ticket.Code,
		Title:               ticket.Title,
		Description:         ticket.Description,
		Category:            ticket.Category,
		Source:              ticket.Source,
		Status:              ticket.Status,
		Priority:            ticket.Priority,
		AssignedTo:          assignedToDTO,
		Assignees:           assigneesDTO,
		Lead:                leadDTO,
		CreatedBy:           createdByDTO,
		RequesterID:         ticket.RequesterID,
		Requester:           requesterDTO,
		RequesterName:       requesterName,
		RequesterDepartment: ticket.RequesterDepartment,
		FilialeID:           ticket.FilialeID,
		Filiale:             filialeDTO,
		SoftwareID:          ticket.SoftwareID,
		Software:            softwareDTO,
		ValidatedByUserID:   ticket.ValidatedByUserID,
		ValidatedBy:         validatedByDTO,
		ValidatedAt:         ticket.ValidatedAt,
		EstimatedTime:       ticket.EstimatedTime,
		ActualTime:          ticket.ActualTime,
		ParentID:            ticket.ParentID,
		SubTickets:          subTickets,
		CreatedAt:           ticket.CreatedAt,
		UpdatedAt:           ticket.UpdatedAt,
		ClosedAt:            ticket.ClosedAt,
	}
}

// commentToDTO convertit un modèle TicketComment en DTO TicketCommentDTO
func (s *ticketService) commentToDTO(comment *models.TicketComment) dto.TicketCommentDTO {
	userDTO := s.userToDTO(&comment.User)
	return dto.TicketCommentDTO{
		ID:         comment.ID,
		TicketID:   comment.TicketID,
		User:       userDTO,
		Comment:    comment.Comment,
		IsInternal: comment.IsInternal,
		CreatedAt:  comment.CreatedAt,
		UpdatedAt:  comment.UpdatedAt,
	}
}

// applySLAIfApplicable applique automatiquement un SLA au ticket s'il existe une règle correspondante
func (s *ticketService) applySLAIfApplicable(ticket *models.Ticket) {
	// Vérifier si un SLA existe déjà pour ce ticket
	existingTicketSLA, err := s.ticketSLARepo.FindByTicketID(ticket.ID)
	if err == nil && existingTicketSLA != nil {
		// Un SLA existe déjà, ne pas en créer un autre
		return
	}

	// Chercher un SLA actif correspondant à la catégorie et priorité du ticket
	var sla *models.SLA
	var errSLA error

	// D'abord, chercher un SLA spécifique à la priorité (si le ticket a une priorité)
	if ticket.Priority != "" {
		sla, errSLA = s.slaRepo.FindByCategoryAndPriority(ticket.Category, ticket.Priority)
		if errSLA == nil && sla != nil {
			// SLA spécifique trouvé, l'utiliser
			log.Printf("SLA spécifique trouvé pour ticket %d: catégorie=%s, priorité=%s", ticket.ID, ticket.Category, ticket.Priority)
		}
	}

	// Si aucun SLA spécifique n'est trouvé, chercher un SLA général (sans priorité = NULL)
	if sla == nil || errSLA != nil {
		sla, errSLA = s.slaRepo.FindByCategoryAndPriority(ticket.Category, "")
		if errSLA == nil && sla != nil {
			log.Printf("SLA général trouvé pour ticket %d: catégorie=%s (sans priorité spécifique)", ticket.ID, ticket.Category)
		}
	}

	if errSLA != nil && sla == nil {
		log.Printf("Aucun SLA trouvé pour ticket %d: catégorie=%s, priorité=%s", ticket.ID, ticket.Category, ticket.Priority)
	}

	// Si un SLA est trouvé, créer l'association ticket-SLA
	if sla != nil && errSLA == nil {
		// Calculer la date cible : created_at + target_time
		var targetTime time.Time
		switch sla.Unit {
		case "minutes":
			targetTime = ticket.CreatedAt.Add(time.Duration(sla.TargetTime) * time.Minute)
		case "hours":
			targetTime = ticket.CreatedAt.Add(time.Duration(sla.TargetTime) * time.Hour)
		case "days":
			targetTime = ticket.CreatedAt.AddDate(0, 0, sla.TargetTime)
		default:
			// Par défaut, traiter comme minutes
			targetTime = ticket.CreatedAt.Add(time.Duration(sla.TargetTime) * time.Minute)
		}

		// Déterminer le statut initial
		status := "on_time"
		now := time.Now()
		if now.After(targetTime) {
			status = "violated"
		} else {
			// Vérifier si on est à risque (moins de 25% du temps restant)
			totalDuration := targetTime.Sub(ticket.CreatedAt)
			if totalDuration > 0 {
				remainingPercent := float64(targetTime.Sub(now)) / float64(totalDuration)
				if remainingPercent < 0.25 {
					status = "at_risk"
				}
			}
		}

		ticketSLA := &models.TicketSLA{
			TicketID:   ticket.ID,
			SLAID:      sla.ID,
			TargetTime: targetTime,
			Status:     status,
		}

		if err := s.ticketSLARepo.Create(ticketSLA); err != nil {
			log.Printf("Erreur lors de l'application du SLA au ticket %d: %v", ticket.ID, err)
		} else {
			log.Printf("SLA appliqué au ticket %d: %s (cible: %v, statut: %s)", ticket.ID, sla.Name, targetTime, status)
		}
	}
}

// updateSLAOnClose met à jour le statut SLA lorsqu'un ticket est clôturé
func (s *ticketService) updateSLAOnClose(ticketID uint) {
	ticketSLA, err := s.ticketSLARepo.FindByTicketID(ticketID)
	if err != nil || ticketSLA == nil {
		// Pas de SLA associé, rien à faire
		return
	}

	ticket, err := s.ticketRepo.FindByID(ticketID)
	if err != nil || ticket == nil {
		return
	}

	// Mettre à jour ActualTime et recalculer le statut
	now := time.Now()
	ticketSLA.ActualTime = &now

	// Recalculer le statut
	if now.After(ticketSLA.TargetTime) {
		ticketSLA.Status = "violated"
		// Calculer le temps de violation en minutes
		violationMinutes := int(now.Sub(ticketSLA.TargetTime).Minutes())
		ticketSLA.ViolationTime = &violationMinutes
	} else {
		ticketSLA.Status = "on_time"
		ticketSLA.ViolationTime = nil
	}

	if err := s.ticketSLARepo.Update(ticketSLA); err != nil {
		log.Printf("Erreur lors de la mise à jour du SLA pour le ticket %d: %v", ticketID, err)
	} else {
		log.Printf("SLA mis à jour pour le ticket %d: statut=%s", ticketID, ticketSLA.Status)
	}
}

// isUserITOfSupplierFiliale vérifie si l'utilisateur appartient au département IT de la filiale fournisseur
func (s *ticketService) isUserITOfSupplierFiliale(userID uint) (bool, error) {
	// Récupérer l'utilisateur avec son département et sa filiale
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return false, err
	}

	// Vérifier si l'utilisateur a un département
	if user.DepartmentID == nil {
		return false, nil
	}

	// Récupérer le département avec sa filiale
	dept, err := s.departmentRepo.FindByID(*user.DepartmentID)
	if err != nil {
		return false, err
	}

	// Vérifier si le département appartient à une filiale
	if dept.FilialeID == nil {
		return false, nil
	}

	// Récupérer la filiale
	filiale, err := s.filialeRepo.FindByID(*dept.FilialeID)
	if err != nil {
		return false, err
	}

	// Vérifier si c'est la filiale fournisseur ET si le département est IT
	return filiale.IsSoftwareProvider && dept.IsITDepartment, nil
}

// getITUsersOfSoftwareProvider récupère tous les utilisateurs actifs du département IT de la filiale fournisseur de logiciels
func (s *ticketService) getITUsersOfSoftwareProvider() ([]uint, error) {
	// Trouver la filiale fournisseur de logiciels
	providerFiliale, err := s.filialeRepo.FindSoftwareProvider()
	if err != nil {
		return nil, fmt.Errorf("filiale fournisseur de logiciels introuvable: %w", err)
	}

	// Trouver les départements IT de cette filiale
	var itDepartments []models.Department
	err = database.DB.Where("filiale_id = ? AND is_it_department = ? AND is_active = ?", providerFiliale.ID, true, true).Find(&itDepartments).Error
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération des départements IT: %w", err)
	}

	if len(itDepartments) == 0 {
		return []uint{}, nil // Aucun département IT trouvé
	}

	// Extraire les IDs des départements
	departmentIDs := make([]uint, len(itDepartments))
	for i, d := range itDepartments {
		departmentIDs[i] = d.ID
	}

	// Trouver tous les utilisateurs actifs de ces départements
	var users []models.User
	err = database.DB.Where("department_id IN ? AND is_active = ?", departmentIDs, true).Find(&users).Error
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération des utilisateurs IT: %w", err)
	}

	// Extraire les IDs des utilisateurs
	userIDs := make([]uint, len(users))
	for i, user := range users {
		userIDs[i] = user.ID
	}

	return userIDs, nil
}

// createNotification crée une notification pour un utilisateur via le NotificationService (pour WebSocket)
func (s *ticketService) createNotification(userID uint, notificationType string, title string, message string, linkURL string, metadata map[string]any) {
	if s.notificationService != nil {
		if err := s.notificationService.Create(userID, notificationType, title, message, linkURL, metadata); err != nil {
			log.Printf("Erreur lors de la création de la notification pour l'utilisateur %d: %v", userID, err)
		}
	} else {
		// Fallback si le service n'est pas disponible (ne devrait pas arriver)
		log.Printf("Warning: NotificationService non disponible, création directe via repository")
		notification := &models.Notification{
			UserID:  userID,
			Type:    notificationType,
			Title:   title,
			Message: message,
			LinkURL: linkURL,
			IsRead:  false,
		}

		if metadata != nil {
			metadataJSON, err := json.Marshal(metadata)
			if err == nil {
				notification.Metadata = metadataJSON
			}
		}

		if err := s.notificationRepo.Create(notification); err != nil {
			log.Printf("Erreur lors de la création de la notification pour l'utilisateur %d: %v", userID, err)
		}
	}
}

// notifyITDepartmentOfSoftwareProvider envoie une notification à tous les utilisateurs IT de la filiale fournisseur de logiciels
func (s *ticketService) notifyITDepartmentOfSoftwareProvider(notificationType string, title string, message string, linkURL string, metadata map[string]any) {
	itUserIDs, err := s.getITUsersOfSoftwareProvider()
	if err != nil {
		log.Printf("Erreur lors de la récupération des utilisateurs IT de la filiale fournisseur: %v", err)
		return
	}

	if len(itUserIDs) == 0 {
		log.Printf("⚠️  Aucun utilisateur IT trouvé dans la filiale fournisseur pour la notification de type: %s", notificationType)
		log.Printf("   Vérifiez qu'un département de la filiale fournisseur est marqué comme IT (is_it_department=true)")
		log.Printf("   et qu'au moins un utilisateur actif est assigné à ce département")
		return
	}

	log.Printf("✅ Envoi de notification '%s' à %d utilisateur(s) IT de la filiale fournisseur", notificationType, len(itUserIDs))
	for _, userID := range itUserIDs {
		s.createNotification(userID, notificationType, title, message, linkURL, metadata)
	}
}

// userToDTO convertit un modèle User en DTO UserDTO (méthode utilitaire)
func (s *ticketService) userToDTO(user *models.User) dto.UserDTO {
	roleName := ""
	if user.Role.ID != 0 {
		roleName = user.Role.Name
	}

	userDTO := dto.UserDTO{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		DepartmentID: user.DepartmentID,
		Avatar:       user.Avatar,
		Role:         roleName,
		// Pas besoin des permissions ici (contexte tickets uniquement)
		IsActive:  user.IsActive,
		LastLogin: user.LastLogin,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	// Inclure le département complet si présent
	if user.Department != nil {
		departmentDTO := dto.DepartmentDTO{
			ID:          user.Department.ID,
			Name:        user.Department.Name,
			Code:        user.Department.Code,
			Description: user.Department.Description,
			OfficeID:    user.Department.OfficeID,
			IsActive:    user.Department.IsActive,
			CreatedAt:   user.Department.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   user.Department.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
		// Inclure le siège si présent
		if user.Department.Office != nil {
			departmentDTO.Office = &dto.OfficeDTO{
				ID:        user.Department.Office.ID,
				Name:      user.Department.Office.Name,
				Country:   user.Department.Office.Country,
				City:      user.Department.Office.City,
				Commune:   user.Department.Office.Commune,
				Address:   user.Department.Office.Address,
				Longitude: user.Department.Office.Longitude,
				Latitude:  user.Department.Office.Latitude,
				IsActive:  user.Department.Office.IsActive,
				CreatedAt: user.Department.Office.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
				UpdatedAt: user.Department.Office.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
			}
		}
		userDTO.Department = &departmentDTO
	}

	return userDTO
}
