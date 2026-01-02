package database

import (
	"fmt"
	"log"

	"github.com/mcicare/itsm-backend/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB est l'instance globale de GORM pour accéder à la base de données
var DB *gorm.DB

// Connect établit la connexion à la base de données MySQL
// Utilise les paramètres de configuration pour construire le DSN (Data Source Name)
func Connect() error {
	// Construction du DSN (Data Source Name) pour MySQL
	// Format: user:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%t&loc=%s",
		config.AppConfig.DBUser,
		config.AppConfig.DBPassword,
		config.AppConfig.DBHost,
		config.AppConfig.DBPort,
		config.AppConfig.DBName,
		config.AppConfig.DBCharset,
		config.AppConfig.DBParseTime,
		config.AppConfig.DBLoc,
	)

	var err error
	// Configuration de GORM avec le logger
	// En développement, on affiche les requêtes SQL (logger.Info)
	// En production, on peut utiliser logger.Silent pour désactiver les logs
	logLevel := logger.Info
	if config.AppConfig.AppEnv == "production" {
		logLevel = logger.Error // En production, on log seulement les erreurs
	}

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})

	if err != nil {
		return fmt.Errorf("erreur de connexion à la base de données: %w", err)
	}

	log.Println("Connexion à la base de données MySQL réussie")
	return nil
}

// Close ferme la connexion à la base de données
func Close() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

// GetDB retourne l'instance GORM (utile pour les tests ou accès direct)
func GetDB() *gorm.DB {
	return DB
}
