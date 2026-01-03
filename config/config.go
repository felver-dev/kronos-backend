package config

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config contient toute la configuration de l'application
type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	App      ApplicationConfig

	// Champs de compatibilité pour l'accès direct (deprecated, utiliser Database/Server/App)
	DBHost                   string
	DBPort                   string
	DBUser                   string
	DBPassword               string
	DBName                   string
	DBCharset                string
	DBParseTime              bool
	DBLoc                    string
	AppEnv                   string
	AppPort                  string
	AppName                  string
	AppURL                   string
	JWTSecret                string
	JWTExpirationHours       int
	JWTRefreshExpirationDays int
	UploadDir                string
	MaxUploadSize            int64
	AllowedImageTypes        []string
	AvatarMaxSize            int64
	AvatarDir                string
	TicketAttachmentsDir     string
}

// DatabaseConfig contient les paramètres de connexion à la base de données
type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	Name            string
	Charset         string
	ParseTime       bool
	Loc             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// ServerConfig contient la configuration du serveur HTTP
type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// ApplicationConfig contient la configuration générale de l'application
type ApplicationConfig struct {
	Name                     string
	Environment              string
	URL                      string
	LogLevel                 string
	JWTSecret                string
	JWTExpirationHours       int
	JWTRefreshExpirationDays int
	UploadDir                string
	MaxUploadSize            int64
	AllowedImageTypes        []string
	AvatarMaxSize            int64
	AvatarDir                string
	TicketAttachmentsDir     string
}

// AppConfig est l'instance globale de configuration
var AppConfig *Config

// loadEnvFile charge le fichier .env en gérant le BOM UTF-8
func loadEnvFile() {
	wd, _ := os.Getwd()
	envPaths := []string{
		filepath.Join(wd, ".env"),
		".env",
		filepath.Join(wd, "..", ".env"),
		"../.env",
	}

	for _, path := range envPaths {
		if _, err := os.Stat(path); err == nil {
			// Lire le fichier
			content, err := os.ReadFile(path)
			if err != nil {
				continue
			}

			// Supprimer le BOM UTF-8 s'il existe
			content = bytes.TrimPrefix(content, []byte("\xef\xbb\xbf"))

			// Parser et charger dans les variables d'environnement
			envMap, err := godotenv.UnmarshalBytes(content)
			if err == nil {
				for key, value := range envMap {
					os.Setenv(key, value)
				}
				log.Printf("✅ Fichier .env chargé depuis: %s", path)
				return
			}
		}
	}
}

// Load charge la configuration depuis les variables d'environnement
func Load() (*Config, error) {
	// Charger le fichier .env si présent
	loadEnvFile()

	env := getEnv("APP_ENV", "development")

	config := &Config{
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "127.0.0.1"),
			Port:            getEnv("DB_PORT", "3306"),
			User:            getEnv("DB_USER", "root"),
			Password:        getEnv("DB_PASSWORD", ""),
			Name:            getEnv("DB_NAME", "itsm_db"),
			Charset:         getEnv("DB_CHARSET", "utf8mb4"),
			ParseTime:       getEnvBool("DB_PARSE_TIME", true),
			Loc:             getEnv("DB_LOC", "Local"),
			MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 100),
			MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 10),
			ConnMaxLifetime: getEnvAsDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),
			ConnMaxIdleTime: getEnvAsDuration("DB_CONN_MAX_IDLE_TIME", 10*time.Minute),
		},
		Server: ServerConfig{
			Port:         getEnv("APP_PORT", "8080"),
			ReadTimeout:  getEnvAsDuration("SERVER_READ_TIMEOUT", 15*time.Second),
			WriteTimeout: getEnvAsDuration("SERVER_WRITE_TIMEOUT", 15*time.Second),
			IdleTimeout:  getEnvAsDuration("SERVER_IDLE_TIMEOUT", 60*time.Second),
		},
		App: ApplicationConfig{
			Name:                     getEnv("APP_NAME", "ITSM Backend"),
			Environment:              env,
			URL:                      getEnv("APP_URL", "http://localhost:8080"),
			LogLevel:                 getEnv("LOG_LEVEL", getDefaultLogLevel(env)),
			JWTSecret:                getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-in-production"),
			JWTExpirationHours:       getEnvAsInt("JWT_EXPIRATION_HOURS", 24),
			JWTRefreshExpirationDays: getEnvAsInt("JWT_REFRESH_EXPIRATION_DAYS", 7),
			UploadDir:                getEnv("UPLOAD_DIR", "./uploads"),
			MaxUploadSize:            getEnvAsInt64("MAX_UPLOAD_SIZE", 10485760), // 10 MB
			AllowedImageTypes:        getEnvSlice("ALLOWED_IMAGE_TYPES", []string{"jpg", "jpeg", "png", "gif", "webp"}),
			AvatarMaxSize:            getEnvAsInt64("AVATAR_MAX_SIZE", 2097152), // 2 MB
			AvatarDir:                getEnv("AVATAR_DIR", "./uploads/users"),
			TicketAttachmentsDir:     getEnv("TICKET_ATTACHMENTS_DIR", "./uploads/tickets"),
		},
	}

	// Remplir les champs de compatibilité pour l'accès direct
	config.DBHost = config.Database.Host
	config.DBPort = config.Database.Port
	config.DBUser = config.Database.User
	config.DBPassword = config.Database.Password
	config.DBName = config.Database.Name
	config.DBCharset = config.Database.Charset
	config.DBParseTime = config.Database.ParseTime
	config.DBLoc = config.Database.Loc
	config.AppEnv = config.App.Environment
	config.AppPort = config.Server.Port
	config.AppName = config.App.Name
	config.AppURL = config.App.URL
	config.JWTSecret = config.App.JWTSecret
	config.JWTExpirationHours = config.App.JWTExpirationHours
	config.JWTRefreshExpirationDays = config.App.JWTRefreshExpirationDays
	config.UploadDir = config.App.UploadDir
	config.MaxUploadSize = config.App.MaxUploadSize
	config.AllowedImageTypes = config.App.AllowedImageTypes
	config.AvatarMaxSize = config.App.AvatarMaxSize
	config.AvatarDir = config.App.AvatarDir
	config.TicketAttachmentsDir = config.App.TicketAttachmentsDir

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("configuration invalide: %w", err)
	}

	// Créer les dossiers d'upload si nécessaire
	createDirs(config)

	// Log de la configuration de la base de données (sans le mot de passe)
	log.Printf("📊 Configuration DB: Host=%s, Port=%s, User=%s, Database=%s",
		config.Database.Host, config.Database.Port, config.Database.User, config.Database.Name)

	return config, nil
}

// LoadConfig charge la configuration depuis les variables d'environnement (compatibilité)
func LoadConfig() {
	cfg, err := Load()
	if err != nil {
		log.Fatalf("❌ Erreur lors du chargement de la configuration: %v", err)
	}
	AppConfig = cfg
}

// Validate vérifie que la configuration est valide
func (c *Config) Validate() error {
	if c.Database.Host == "" {
		return fmt.Errorf("DB_HOST ne peut pas être vide")
	}
	if c.Database.Port == "" {
		return fmt.Errorf("DB_PORT ne peut pas être vide")
	}
	if c.Database.User == "" {
		return fmt.Errorf("DB_USER ne peut pas être vide")
	}
	if c.Database.Name == "" {
		return fmt.Errorf("DB_NAME ne peut pas être vide")
	}
	if c.Server.Port == "" {
		return fmt.Errorf("APP_PORT ne peut pas être vide")
	}
	if c.App.JWTSecret == "" || c.App.JWTSecret == "your-super-secret-jwt-key-change-in-production" {
		if c.IsProduction() {
			return fmt.Errorf("JWT_SECRET doit être défini en production")
		}
	}
	return nil
}

// IsDevelopment retourne true si l'application est en mode développement
func (c *Config) IsDevelopment() bool {
	return c.App.Environment == "dev" || c.App.Environment == "development"
}

// IsProduction retourne true si l'application est en mode production
func (c *Config) IsProduction() bool {
	return c.App.Environment == "production" || c.App.Environment == "prod"
}

// getEnv récupère une variable d'environnement ou retourne la valeur par défaut
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt récupère une variable d'environnement comme entier
func getEnvAsInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}

// getEnvAsInt64 récupère une variable d'environnement comme int64
func getEnvAsInt64(key string, defaultValue int64) int64 {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return defaultValue
	}
	return intValue
}

// getEnvAsDuration récupère une variable d'environnement comme durée
func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	duration, err := time.ParseDuration(value)
	if err != nil {
		return defaultValue
	}
	return duration
}

// getEnvBool récupère une variable d'environnement comme booléen
func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// getEnvSlice récupère une variable d'environnement comme slice de strings (séparée par des virgules)
func getEnvSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		// Split par virgule et nettoyer les espaces
		parts := strings.Split(value, ",")
		values := []string{}
		for _, part := range parts {
			trimmed := strings.TrimSpace(part)
			if trimmed != "" {
				values = append(values, trimmed)
			}
		}
		if len(values) > 0 {
			return values
		}
	}
	return defaultValue
}

// getDefaultLogLevel retourne le niveau de log par défaut selon l'environnement
func getDefaultLogLevel(env string) string {
	switch env {
	case "dev", "development":
		return "debug"
	case "staging":
		return "info"
	case "production", "prod":
		return "warn"
	default:
		return "info"
	}
}

// createDirs crée les dossiers nécessaires pour les uploads
func createDirs(cfg *Config) {
	dirs := []string{
		cfg.App.UploadDir,
		cfg.App.AvatarDir,
		cfg.App.TicketAttachmentsDir,
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Printf("Erreur lors de la création du dossier %s: %v", dir, err)
		}
	}
}
