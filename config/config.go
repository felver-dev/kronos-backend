package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config contient toute la configuration de l'application
type Config struct {
	// Application
	AppName string
	AppEnv  string // development, production
	AppPort string
	AppURL  string

	// Base de données MySQL
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	DBCharset   string
	DBParseTime bool
	DBLoc       string

	// JWT
	JWTSecret                string
	JWTExpirationHours       int
	JWTRefreshExpirationDays int

	// Uploads
	UploadDir         string
	MaxUploadSize     int64 // en bytes
	AllowedImageTypes []string

	// Avatars
	AvatarMaxSize int64 // en bytes
	AvatarDir     string

	// Tickets
	TicketAttachmentsDir string
}

// AppConfig est l'instance globale de configuration
var AppConfig *Config

// LoadConfig charge la configuration depuis les variables d'environnement
func LoadConfig() {
	// Charger le fichier .env si présent (ignoré si le fichier n'existe pas)
	godotenv.Load()

	AppConfig = &Config{
		// Application
		AppName: getEnv("APP_NAME", "ITSM Backend"),
		AppEnv:  getEnv("APP_ENV", "development"),
		AppPort: getEnv("APP_PORT", "8080"),
		AppURL:  getEnv("APP_URL", "http://localhost:8080"),

		// Base de données
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "3306"),
		DBUser:      getEnv("DB_USER", "root"),
		DBPassword:  getEnv("DB_PASSWORD", ""),
		DBName:      getEnv("DB_NAME", "itsm_db"),
		DBCharset:   getEnv("DB_CHARSET", "utf8mb4"),
		DBParseTime: getEnvBool("DB_PARSE_TIME", true),
		DBLoc:       getEnv("DB_LOC", "Local"),

		// JWT
		JWTSecret:                getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-in-production"),
		JWTExpirationHours:       getEnvInt("JWT_EXPIRATION_HOURS", 24),
		JWTRefreshExpirationDays: getEnvInt("JWT_REFRESH_EXPIRATION_DAYS", 7),

		// Uploads
		UploadDir:         getEnv("UPLOAD_DIR", "./uploads"),
		MaxUploadSize:     getEnvInt64("MAX_UPLOAD_SIZE", 10485760), // 10 MB
		AllowedImageTypes: getEnvSlice("ALLOWED_IMAGE_TYPES", []string{"jpg", "jpeg", "png", "gif", "webp"}),

		// Avatars
		AvatarMaxSize: getEnvInt64("AVATAR_MAX_SIZE", 2097152), // 2 MB
		AvatarDir:     getEnv("AVATAR_DIR", "./uploads/users"),

		// Tickets
		TicketAttachmentsDir: getEnv("TICKET_ATTACHMENTS_DIR", "./uploads/tickets"),
	}

	// Créer les dossiers d'upload si nécessaire
	createDirs()
}

// getEnv récupère une variable d'environnement ou retourne la valeur par défaut
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt récupère une variable d'environnement comme entier
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvInt64 récupère une variable d'environnement comme int64
func getEnvInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
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

// createDirs crée les dossiers nécessaires pour les uploads
func createDirs() {
	dirs := []string{
		AppConfig.UploadDir,
		AppConfig.AvatarDir,
		AppConfig.TicketAttachmentsDir,
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Printf("Erreur lors de la création du dossier %s: %v", dir, err)
		}
	}
}
