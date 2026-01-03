package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql" // Driver MySQL pour database/sql
	"github.com/mcicare/itsm-backend/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB est l'instance globale de GORM pour accéder à la base de données
var DB *gorm.DB

// Connect établit la connexion à la base de données MySQL
// Utilise les paramètres de configuration pour construire le DSN (Data Source Name)
// Crée la base de données si elle n'existe pas
func Connect() error {
	// D'abord, se connecter sans spécifier la base de données pour pouvoir la créer
	dsnWithoutDB := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=%s&parseTime=%t&loc=%s",
		config.AppConfig.DBUser,
		config.AppConfig.DBPassword,
		config.AppConfig.DBHost,
		config.AppConfig.DBPort,
		config.AppConfig.DBCharset,
		config.AppConfig.DBParseTime,
		config.AppConfig.DBLoc,
	)

	// Se connecter sans base de données spécifiée
	tempDB, err := sql.Open("mysql", dsnWithoutDB)
	if err != nil {
		return fmt.Errorf("erreur de connexion au serveur MySQL: %w", err)
	}
	defer tempDB.Close()

	// Vérifier la connexion
	if err = tempDB.Ping(); err != nil {
		return fmt.Errorf("impossible de se connecter au serveur MySQL: %w. Vérifiez que MySQL est démarré et que les identifiants sont corrects", err)
	}

	// Créer la base de données si elle n'existe pas
	createDBQuery := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", config.AppConfig.DBName)
	_, err = tempDB.Exec(createDBQuery)
	if err != nil {
		return fmt.Errorf("erreur lors de la création de la base de données: %w", err)
	}
	log.Printf("✅ Base de données '%s' vérifiée/créée", config.AppConfig.DBName)

	// Maintenant, se connecter à la base de données spécifiée
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

	// Configurer le pool de connexions
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("erreur lors de la récupération de l'instance SQL: %w", err)
	}

	// Configurer le pool de connexions
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	// Vérifier la connexion
	if err = sqlDB.Ping(); err != nil {
		return fmt.Errorf("erreur lors du ping de la base de données: %w", err)
	}

	log.Println("✅ Connexion à la base de données MySQL réussie")
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
