package database

import (
	"fmt"
	"log"
	"time"

	"github.com/mcicare/itsm-backend/internal/models"
	"github.com/mcicare/itsm-backend/internal/utils"
)

// SeedDemoData génère des données de démonstration
func SeedDemoData() error {
	if DB == nil {
		return fmt.Errorf("la base de données n'est pas initialisée")
	}

	log.Println("🌱 Démarrage du seeding des données de démonstration...")

	// 1. Créer des utilisateurs supplémentaires
	users := []struct {
		Username  string
		Email     string
		FirstName string
		LastName  string
		RoleName  string
		Password  string
	}{
		{"tech", "tech@kronos.com", "Thomas", "Technicien", "USER", "kronos123"},
		{"user", "user@kronos.com", "Alice", "Utilisateur", "USER", "kronos123"},
	}

	for _, u := range users {
		var user models.User
		// Vérifier si l'utilisateur existe déjà
		if err := DB.Where("username = ?", u.Username).First(&user).Error; err == nil {
			log.Printf("   ℹ️  Utilisateur %s existe déjà", u.Username)
			continue
		}

		// Récupérer le rôle
		var role models.Role
		if err := DB.Where("name = ?", u.RoleName).First(&role).Error; err != nil {
			log.Printf("   ⚠️  Rôle %s non trouvé pour l'utilisateur %s", u.RoleName, u.Username)
			continue
		}

		// Hasher le mot de passe
		hashedPassword, _ := utils.HashPassword(u.Password)

		newUser := models.User{
			Username:     u.Username,
			Email:        u.Email,
			FirstName:    u.FirstName,
			LastName:     u.LastName,
			PasswordHash: hashedPassword,
			RoleID:       role.ID,
			IsActive:     true,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		if err := DB.Create(&newUser).Error; err != nil {
			log.Printf("   ⚠️  Erreur lors de la création de l'utilisateur %s: %v", u.Username, err)
		} else {
			log.Printf("   ✅ Utilisateur %s créé (pass: %s)", u.Username, u.Password)
		}
	}

	// 2. Créer quelques tickets de démo
	// Nécessite de récupérer un utilisateur et une catégorie
	var adminUser models.User
	DB.Where("username = ?", "admin").First(&adminUser)

	var incidentCat models.TicketCategory
	DB.Where("slug = ?", "incident").First(&incidentCat)

	if adminUser.ID != 0 && incidentCat.ID != 0 {
		tickets := []struct {
			Title       string
			Description string
			Priority    string
			Status      string
		}{
			{"Imprimante HS", "L'imprimante du 2ème étage ne répond plus.", "HIGH", "OPEN"},
			{"Wifi lent", "La connexion wifi est très lente dans la salle de réunion.", "MEDIUM", "IN_PROGRESS"},
			{"Demande de licence", "Besoin d'une licence Photoshop pour le marketing.", "LOW", "OPEN"},
		}

		for _, t := range tickets {
			ticket := models.Ticket{
				Title:       t.Title,
				Description: t.Description,
				Status:      t.Status,
				Priority:    t.Priority,
				CategoryID:  &incidentCat.ID,
				RequesterID: &adminUser.ID,
				CreatedByID: adminUser.ID,
				Code:        fmt.Sprintf("INC-%d", time.Now().UnixNano()%10000), // Code temporaire simple
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}

			if err := DB.Create(&ticket).Error; err != nil {
				log.Printf("   ⚠️  Erreur création ticket %s: %v", t.Title, err)
			} else {
				log.Printf("   ✅ Ticket créé: %s", t.Title)
			}
		}
	}

	log.Println("✅ Données de démonstration générées")
	return nil
}
