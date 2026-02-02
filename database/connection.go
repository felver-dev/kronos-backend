package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql" // Driver MySQL pour database/sql
	"github.com/mcicare/itsm-backend/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB *gorm.DB
)

// InitDB initialise la connexion à la base de données MySQL avec GORM
func InitDB(cfg *config.Config) (*gorm.DB, error) {
	// D'abord, se connecter sans spécifier la base de données pour pouvoir la créer
	dsnWithoutDB := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
	)

	// Se connecter sans base de données spécifiée
	tempDB, err := sql.Open("mysql", dsnWithoutDB)
	if err != nil {
		return nil, fmt.Errorf("erreur de connexion au serveur MySQL: %w", err)
	}
	defer tempDB.Close()

	// Vérifier la connexion
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = tempDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("impossible de se connecter au serveur MySQL: %w. Vérifiez que MySQL est démarré et que les identifiants sont corrects", err)
	}

	// Créer la base de données si elle n'existe pas
	createDBQuery := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", cfg.Database.Name)
	_, err = tempDB.ExecContext(ctx, createDBQuery)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la création de la base de données: %w", err)
	}
	log.Printf("✅ Base de données '%s' vérifiée/créée", cfg.Database.Name)

	// Maintenant, se connecter à la base de données spécifiée
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)

	var logLevel logger.LogLevel

	switch cfg.App.LogLevel {
	case "debug":
		logLevel = logger.Info
	case "info":
		logLevel = logger.Warn
	case "warn":
		logLevel = logger.Error
	default:
		logLevel = logger.Silent
	}

	gormConfig := &gorm.Config{
		Logger:                           logger.Default.LogMode(logLevel),
		SkipDefaultTransaction:           true,
		PrepareStmt:                      true,
		DisableForeignKeyConstraintWhenMigrating: true, // Désactiver la création automatique de contraintes FK
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	}

	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("echec de la connexion à la base de données: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récuperation de l'instance DB : %w", err)
	}

	// Utiliser le contexte existant pour le ping
	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("échec du ping à la base de données: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(cfg.Database.ConnMaxIdleTime)

	DB = db

	log.Printf("✅ Connexion à MySQL réussie - Base: %s sur %s:%s",
		cfg.Database.Name,
		cfg.Database.Host,
		cfg.Database.Port,
	)
	return db, nil
}

// Connect établit la connexion à la base de données (compatibilité avec l'ancien code)
func Connect() error {
	if config.AppConfig == nil {
		return fmt.Errorf("configuration non chargée, appelez config.LoadConfig() d'abord")
	}
	_, err := InitDB(config.AppConfig)
	return err
}

// CloseDB ferme proprement la connexion à la base de données
func CloseDB() error {
	if DB == nil {
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("erreur lors de la récupération de l'instance: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("erreur lors de la fermeture de la base de données: %w", err)
	}

	log.Printf("✅ Connexion à la base de données fermée")
	return nil
}

// Close ferme la connexion à la base de données (compatibilité avec l'ancien code)
func Close() error {
	return CloseDB()
}

// GetDB retourne l'instance de la base de données
func GetDB() *gorm.DB {
	return DB
}

// HealthCheck vérifie que la connexion à la base de données est saine
func HealthCheck(ctx context.Context) error {
	if DB == nil {
		return fmt.Errorf("base de données non initialisée")
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("erreur lors de la récupération de l'instance DB: %w", err)
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("ping échoué: %w", err)
	}
	return nil
}
