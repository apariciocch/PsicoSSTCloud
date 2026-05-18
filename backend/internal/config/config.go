package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config contiene toda la configuración de la aplicación
type Config struct {
	Server     ServerConfig
	PocketBase PocketBaseConfig
	JWT        JWTConfig
	Database   DatabaseConfig
	Security   SecurityConfig
	Logs       LogsConfig
	SMTP       SMTPConfig
	Files      FilesConfig
	Pagination PaginationConfig
}

// ServerConfig configuración del servidor
type ServerConfig struct {
	Host        string
	Port        string
	Environment string
}

// PocketBaseConfig configuración de PocketBase
type PocketBaseConfig struct {
	URL           string
	AdminEmail    string
	AdminPassword string
	DataDir       string
}

// JWTConfig configuración JWT
type JWTConfig struct {
	Secret            string
	Expiration        int
	RefreshExpiration int
}

// DatabaseConfig configuración de BD
type DatabaseConfig struct {
	URL string
}

// SecurityConfig configuración de seguridad
type SecurityConfig struct {
	CORSOrigins         []string
	RateLimitRequests   int
	RateLimitWindow     int
	MaxLoginAttempts    int
	LockDurationMinutes int
}

// LogsConfig configuración de logs
type LogsConfig struct {
	Level  string
	Format string
}

// SMTPConfig configuración de email
type SMTPConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	From     string
}

// FilesConfig configuración de archivos
type FilesConfig struct {
	MaxSizeMB    int
	AllowedTypes string
	StoragePath  string
}

// PaginationConfig configuración de paginación
type PaginationConfig struct {
	DefaultPageSize int
	MaxPageSize     int
}

// Load carga la configuración desde variables de entorno
func Load() (*Config, error) {
	// Cargar .env si existe
	_ = godotenv.Load()

	cfg := &Config{
		Server: ServerConfig{
			Host:        getEnv("SERVER_HOST", "0.0.0.0"),
			Port:        getEnv("SERVER_PORT", "8080"),
			Environment: getEnv("ENVIRONMENT", "development"),
		},
		PocketBase: PocketBaseConfig{
			URL:           getEnv("POCKETBASE_URL", "http://localhost:8090"),
			AdminEmail:    getEnv("POCKETBASE_ADMIN_EMAIL", "admin@example.com"),
			AdminPassword: getEnv("POCKETBASE_ADMIN_PASSWORD", ""),
			DataDir:       getEnv("POCKETBASE_DATA_DIR", "./pb_data"),
		},
		JWT: JWTConfig{
			Secret:            getEnv("JWT_SECRET", ""),
			Expiration:        getEnvInt("JWT_EXPIRATION", 900),
			RefreshExpiration: getEnvInt("JWT_REFRESH_EXPIRATION", 604800),
		},
		Database: DatabaseConfig{
			URL: getEnv("DATABASE_URL", "file:./pb_data/data.db"),
		},
		Security: SecurityConfig{
			CORSOrigins:         splitEnv("CORS_ORIGINS", "http://localhost:3000"),
			RateLimitRequests:   getEnvInt("RATE_LIMIT_REQUESTS", 100),
			RateLimitWindow:     getEnvInt("RATE_LIMIT_WINDOW", 60),
			MaxLoginAttempts:    getEnvInt("MAX_LOGIN_ATTEMPTS", 5),
			LockDurationMinutes: getEnvInt("LOCK_DURATION_MINUTES", 30),
		},
		Logs: LogsConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
		SMTP: SMTPConfig{
			Host:     getEnv("SMTP_HOST", "smtp.gmail.com"),
			Port:     getEnvInt("SMTP_PORT", 587),
			User:     getEnv("SMTP_USER", ""),
			Password: getEnv("SMTP_PASSWORD", ""),
			From:     getEnv("SMTP_FROM", "noreply@empresa.com"),
		},
		Files: FilesConfig{
			MaxSizeMB:    getEnvInt("MAX_FILE_SIZE_MB", 10),
			AllowedTypes: getEnv("ALLOWED_FILE_TYPES", "jpg,jpeg,png,gif,pdf"),
			StoragePath:  getEnv("STORAGE_PATH", "./storage/uploads"),
		},
		Pagination: PaginationConfig{
			DefaultPageSize: getEnvInt("DEFAULT_PAGE_SIZE", 20),
			MaxPageSize:     getEnvInt("MAX_PAGE_SIZE", 100),
		},
	}

	return cfg, nil
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func splitEnv(key, defaultValue string) []string {
	if value, exists := os.LookupEnv(key); exists {
		parts := strings.Split(value, ",")
		result := make([]string, 0)
		for _, part := range parts {
			if p := strings.TrimSpace(part); p != "" {
				result = append(result, p)
			}
		}
		if len(result) > 0 {
			return result
		}
	}
	return []string{defaultValue}
}
