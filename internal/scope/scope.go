package scope

import (
	"log"

	"github.com/mcicare/itsm-backend/internal/models"
)

// QueryScope représente le contexte de requête avec les permissions et attributs de l'utilisateur
// Ce scope est utilisé pour filtrer automatiquement les données selon les permissions
type QueryScope struct {
	UserID       uint     // ID de l'utilisateur connecté
	DepartmentID *uint    // ID du département de l'utilisateur (peut être nil)
	FilialeID    *uint    // ID de la filiale de l'utilisateur (peut être nil)
	Role         string   // Nom du rôle (ex: "DSI", "TECHNICIEN_IT")
	Permissions  []string // Liste des permissions de l'utilisateur
	// IsResolver indique si l'utilisateur est un résolveur (département IT de la filiale fournisseur de logiciels)
	// Les résolveurs ont les mêmes capacités que tickets.create_any_filiale pour la création de tickets
	IsResolver bool
	// DepartmentIsIT indique si le département de l'utilisateur est un département IT (is_it_department)
	// Utilisé pour afficher/créer les commentaires internes sur les tickets
	DepartmentIsIT bool
	// FilterUserID filtre optionnel par utilisateur (ex. pour delays: membre du département ou utilisateur quelconque selon la permission)
	FilterUserID *uint
	// FilterFilialeID filtre optionnel par filiale (pour les rapports et vues filtrées)
	FilterFilialeID *uint
	// DashboardScopeHint force le périmètre pour le tableau de bord : "department" | "filiale" | "global" (vide = comportement par permissions)
	DashboardScopeHint string
}

// NewQueryScopeFromUser crée un QueryScope à partir d'un modèle User
func NewQueryScopeFromUser(user *models.User) *QueryScope {
	isResolver := false
	if user.Department != nil && user.Department.IsITDepartment && user.Filiale != nil && user.Filiale.IsSoftwareProvider {
		isResolver = true
	}
	// Filiale : priorité user → rôle → département (DSI filiale peut n'avoir que son département rattaché)
	filialeID := user.FilialeID
	source := "nil"
	if filialeID != nil {
		source = "user"
	} else if user.Role.FilialeID != nil {
		filialeID = user.Role.FilialeID
		source = "role"
	} else if user.Department != nil && user.Department.FilialeID != nil {
		filialeID = user.Department.FilialeID
		source = "department"
	}
	if filialeID != nil {
		log.Printf("[scope] User %d: FilialeID=%d (source=%s)", user.ID, *filialeID, source)
	} else {
		log.Printf("[scope] User %d: FilialeID=nil (définir filiale sur l'utilisateur, le rôle ou un département avec filiale)", user.ID)
	}
	departmentIsIT := user.Department != nil && user.Department.IsITDepartment
	return &QueryScope{
		UserID:         user.ID,
		DepartmentID:   user.DepartmentID,
		FilialeID:      filialeID,
		Role:           user.Role.Name,
		Permissions:    GetPermissionsForRole(user.Role.Name),
		IsResolver:     isResolver,
		DepartmentIsIT: departmentIsIT,
	}
}

// HasPermission vérifie si le scope a une permission donnée
func (s *QueryScope) HasPermission(permission string) bool {
	for _, p := range s.Permissions {
		if p == permission {
			return true
		}
	}
	return false
}

// HasAnyPermission vérifie si le scope a au moins une des permissions données
func (s *QueryScope) HasAnyPermission(permissions ...string) bool {
	for _, perm := range permissions {
		if s.HasPermission(perm) {
			return true
		}
	}
	return false
}

// HasAllPermissions vérifie si le scope a toutes les permissions données
func (s *QueryScope) HasAllPermissions(permissions ...string) bool {
	for _, perm := range permissions {
		if !s.HasPermission(perm) {
			return false
		}
	}
	return true
}

// permissionsGetter est une fonction qui récupère les permissions d'un rôle par son nom
// Cette fonction peut être injectée depuis l'extérieur pour éviter les cycles d'importation
var permissionsGetter func(roleName string) []string

// SetPermissionsGetter définit la fonction de récupération des permissions
// Cette fonction doit être appelée au démarrage de l'application
func SetPermissionsGetter(getter func(roleName string) []string) {
	permissionsGetter = getter
}

// GetPermissionsForRole retourne la liste des permissions associées à un rôle donné
// Les permissions sont récupérées depuis la base de données via la table role_permissions
// Cette fonction doit être identique à celle dans auth_service.go pour la cohérence
func GetPermissionsForRole(roleName string) []string {
	if permissionsGetter == nil {
		log.Printf("⚠️  PermissionsGetter non initialisé pour le rôle '%s', retour de permissions minimales", roleName)
		// Retourner des permissions minimales par défaut si le getter n'est pas défini
		return []string{"tickets.view_own"}
	}
	return permissionsGetter(roleName)
}
