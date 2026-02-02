package main

import (
	"flag"
	"log"

	"github.com/mcicare/itsm-backend/config"
	"github.com/mcicare/itsm-backend/database"
)

func main() {
	// Parse des flags
	seed := flag.Bool("seed", false, "Exécuter le seeding des données initiales après les migrations")
	reset := flag.Bool("reset", false, "Supprimer et recréer toutes les tables (ATTENTION: supprime toutes les données!)")
	flag.Parse()

	// Charger la configuration
	config.LoadConfig()

	// Se connecter à la base de données
	if err := database.Connect(); err != nil {
		log.Fatalf("❌ Erreur de connexion à la base de données: %v", err)
	}
	defer database.Close()

	// Reset si demandé
	if *reset {
		log.Println("🔄 Réinitialisation de la base de données...")
		if err := database.ResetDatabase(); err != nil {
			log.Fatalf("❌ Erreur lors de la réinitialisation: %v", err)
		}
		log.Println("✅ Base de données réinitialisée")
		return
	}

	// Exécuter les migrations
	if err := database.AutoMigrate(); err != nil {
		log.Fatalf("❌ Erreur lors des migrations: %v", err)
	}

	// Exécuter le seeding si demandé
	if *seed {
		log.Println("🌱 Exécution du seeding...")
		if err := database.SeedDemoData(); err != nil {
			log.Printf("⚠️  Erreur lors du seeding: %v", err)
		}
	}

	log.Println("✨ Migrations terminées avec succès!")
}
