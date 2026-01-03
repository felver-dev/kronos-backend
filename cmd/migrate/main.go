package main

import (
	"flag"
	"log"

	"github.com/mcicare/itsm-backend/config"
	"github.com/mcicare/itsm-backend/database"
	"github.com/mcicare/itsm-backend/migrations"
	"gorm.io/gorm/logger"
)

func main() {
	// Parse des flags
	seed := flag.Bool("seed", false, "Exécuter le seeding des données initiales après les migrations")
	flag.Parse()

	// Charger la configuration
	config.LoadConfig()

	// Se connecter à la base de données
	if err := database.Connect(); err != nil {
		log.Fatalf("❌ Erreur de connexion à la base de données: %v", err)
	}
	defer database.Close()

	// Désactiver complètement le logging pendant les migrations pour éviter le spam
	database.DB.Logger = logger.Default.LogMode(logger.Silent)

	// Exécuter les migrations
	if err := migrations.RunMigrations(); err != nil {
		log.Fatalf("❌ Erreur lors des migrations: %v", err)
	}

	// Exécuter le seeding si demandé
	if *seed {
		if err := migrations.SeedData(); err != nil {
			log.Fatalf("❌ Erreur lors du seeding: %v", err)
		}
	}

	log.Println("✨ Migrations terminées avec succès!")
}
