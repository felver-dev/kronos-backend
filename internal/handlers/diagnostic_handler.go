package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/database"
	"github.com/mcicare/itsm-backend/internal/models"
	"github.com/mcicare/itsm-backend/internal/repositories"
	"github.com/mcicare/itsm-backend/internal/utils"
)

// DiagnosticHandler gère les endpoints de diagnostic
type DiagnosticHandler struct {
	filialeRepo repositories.FilialeRepository
}

// NewDiagnosticHandler crée une nouvelle instance de DiagnosticHandler
func NewDiagnosticHandler(filialeRepo repositories.FilialeRepository) *DiagnosticHandler {
	return &DiagnosticHandler{
		filialeRepo: filialeRepo,
	}
}

// GetITUsersInfo retourne des informations sur les utilisateurs IT de MCI CARE CI
// @Summary Informations sur les utilisateurs IT
// @Description Retourne la liste des utilisateurs IT de MCI CARE CI pour diagnostic
// @Tags diagnostic
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /diagnostic/it-users [get]
func (h *DiagnosticHandler) GetITUsersInfo(c *gin.Context) {
	// Trouver la filiale fournisseur de logiciels
	providerFiliale, err := h.filialeRepo.FindSoftwareProvider()
	if err != nil {
		utils.ErrorResponse(c, 500, "Filiale fournisseur de logiciels introuvable", err.Error())
		return
	}

	// Trouver les départements IT de cette filiale
	var itDepartments []models.Department
	err = database.DB.Where("filiale_id = ? AND is_it_department = ? AND is_active = ?", providerFiliale.ID, true, true).Find(&itDepartments).Error
	if err != nil {
		utils.ErrorResponse(c, 500, "Erreur lors de la récupération des départements IT", err.Error())
		return
	}

	departmentInfo := make([]map[string]interface{}, 0)
	for _, dept := range itDepartments {
		// Trouver les utilisateurs de ce département
		var users []models.User
		err = database.DB.Where("department_id = ? AND is_active = ?", dept.ID, true).Find(&users).Error
		if err != nil {
			continue
		}

		userInfo := make([]map[string]interface{}, 0)
		for _, user := range users {
			userInfo = append(userInfo, map[string]interface{}{
				"id":         user.ID,
				"username":   user.Username,
				"email":      user.Email,
				"first_name": user.FirstName,
				"last_name":  user.LastName,
			})
		}

		departmentInfo = append(departmentInfo, map[string]interface{}{
			"id":          dept.ID,
			"name":        dept.Name,
			"code":        dept.Code,
			"is_it":       dept.IsITDepartment,
			"users_count": len(users),
			"users":       userInfo,
		})
	}

	result := map[string]interface{}{
		"software_provider_filiale": map[string]interface{}{
			"id":   providerFiliale.ID,
			"name": providerFiliale.Name,
			"code": providerFiliale.Code,
		},
		"it_departments_count": len(itDepartments),
		"it_departments":       departmentInfo,
		"total_it_users": func() int {
			total := 0
			for _, dept := range departmentInfo {
				if count, ok := dept["users_count"].(int); ok {
					total += count
				}
			}
			return total
		}(),
	}

	utils.SuccessResponse(c, result, "Informations sur les utilisateurs IT récupérées")
}
