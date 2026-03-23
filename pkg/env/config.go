package env

import "os"

type Config struct {
	Port        string
	JWTSecret   string
	DatabaseURL string
	UploadDir   string
}

func Load() (*Config, error) {
	cfg := &Config{
		Port:      getEnv("CABIN_CHAT_PORT", "1111"),
		JWTSecret: getEnv("CABIN_CHAT_JWT_SECRET", "CabinCrew0726"),
		// This needs to be updated with the actual URL later
		DatabaseURL: getEnv("CABIN_CHAT_DATABASE_URL", "postgres://user:password@localhost:5432/dbname"),
		UploadDir:   getEnv("CABIN_CHAT_UPLOAD_DIR", "./uploads"),
	}
	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
