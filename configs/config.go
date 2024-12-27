package config

import "os"

type R2 struct {
	AccountID  string
	AccessKey  string
	SecretKey  string
	BucketName string
}

type Config struct {
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURI  string
	PostgresURI        string
	DatabaseName       string
	FrontendURL        string
	FlaskURL           string
	R2                 R2
	SecretKey          string
	CookieName         string
}

func LoadConfig() *Config {
	return &Config{
		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectURI:  getEnv("GOOGLE_REDIRECT_URI", ""),
		PostgresURI:        getEnv("POSTGRES_URI", ""),
		DatabaseName:       getEnv("DATABASE_NAME", ""),
		FrontendURL:        getEnv("FRONTEND_URL", "http://localhost:5173"),
		FlaskURL:           getEnv("FLASK_URL", "http://localhost:5000"),
		R2: R2{
			AccountID:  getEnv("R2_ACCOUNT_ID", ""),
			AccessKey:  getEnv("R2_ACCESS_KEY", ""),
			SecretKey:  getEnv("R2_SECRET_KEY", ""),
			BucketName: getEnv("R2_BUCKET_NAME", ""),
		},
		SecretKey:  getEnv("SECRET_KEY", ""),
		CookieName: getEnv("COOKIE_NAME", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
