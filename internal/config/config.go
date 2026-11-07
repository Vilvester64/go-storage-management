package config

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

// Config menyimpan seluruh variabel konfigurasi aplikasi dari environment/.env.
type Config struct {
	Port               string
	GinMode            string
	DBUser             string
	DBPassword         string
	DBHost             string
	DBPort             string
	DBName             string
	JWTSecret          string
	JWTExpirationHours int
	StoragePath        string
	MaxUploadSize      int64
}

// LoadEnv membaca file .env dan memasukkan pasangan key-value ke environment process.
func LoadEnv(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		value = strings.Trim(value, `"'`)

		if _, exists := os.LookupEnv(key); !exists {
			_ = os.Setenv(key, value)
		}
	}

	return scanner.Err()
}

// LoadConfig memuat file .env (jika ada) dan mengembalikan struct Config.
func LoadConfig(envPaths ...string) (*Config, error) {
	path := ".env"
	if len(envPaths) > 0 && envPaths[0] != "" {
		path = envPaths[0]
	}

	if err := LoadEnv(path); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	cfg := &Config{
		Port:               GetEnv("PORT", "8080"),
		GinMode:            GetEnv("GIN_MODE", "debug"),
		DBUser:             GetEnv("DB_USER", "root"),
		DBPassword:         GetEnv("DB_PASSWORD", ""),
		DBHost:             GetEnv("DB_HOST", "127.0.0.1"),
		DBPort:             GetEnv("DB_PORT", "3306"),
		DBName:             GetEnv("DB_NAME", "storage_db"),
		JWTSecret:          GetEnv("JWT_SECRET", "super_secret_key_development_only"),
		JWTExpirationHours: GetEnvAsInt("JWT_EXPIRATION_HOURS", 24),
		StoragePath:        GetEnv("STORAGE_PATH", "./uploads"),
		MaxUploadSize:      GetEnvAsInt64("MAX_UPLOAD_SIZE", 52428800),
	}

	return cfg, nil
}

// GetEnv mengambil nilai environment variable atau mengembalikan fallback jika kosong.
func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return fallback
}

// GetEnvAsInt mengambil nilai environment variable sebagai int atau mengembalikan fallback.
func GetEnvAsInt(key string, fallback int) int {
	valStr := GetEnv(key, "")
	if valStr == "" {
		return fallback
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return fallback
	}
	return val
}

// GetEnvAsInt64 mengambil nilai environment variable sebagai int64 atau mengembalikan fallback.
func GetEnvAsInt64(key string, fallback int64) int64 {
	valStr := GetEnv(key, "")
	if valStr == "" {
		return fallback
	}
	val, err := strconv.ParseInt(valStr, 10, 64)
	if err != nil {
		return fallback
	}
	return val
}
