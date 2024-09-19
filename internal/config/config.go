package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port int
	Env  string
	DB   struct {
		DSN     string
		Logging bool
	}
	CientURL string
	Cors     struct {
		TrustedOrigins []string
	}
	ServiceApis struct {
		Idenitity struct {
			URL string
		}
	}
	Mailgun struct {
		Domain string
		APIKey string
	}
}

func LoadConfig(cfg *Config) {

	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
	}

	// Load ENV
	env := os.Getenv("ENV")
	if env == "" {
		cfg.Env = "local"
	} else {
		cfg.Env = env
	}

	// Load PORT
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatalf("PORT not available in .env")
	}

	cfg.Port = port

	// Load CLIENT_URL
	client_url := os.Getenv("CLIENT_URL")
	if client_url == "" {
		log.Fatalf("CLIENT_URL not available in .env")
	}

	cfg.CientURL = client_url

	// Load DATABASE_URL
	postgres_url := os.Getenv("POSTGRES_URL")
	if postgres_url == "" {
		log.Fatalf("POSTGRES_URL not available in .env")
	}

	cfg.DB.DSN = postgres_url

	cfg.Cors.TrustedOrigins = []string{"http://localhost:3000"}

	identity_url := os.Getenv("IDENTITY_URL")
	if identity_url == "" {
		log.Fatalf("IDENTITY_URL not available in .env")
	}

	cfg.ServiceApis.Idenitity.URL = identity_url

	mailgun_domain := os.Getenv("MAILGUN_DOMAIN")
	if mailgun_domain == "" {
		log.Fatalf("MAILGUN_DOMAIN not available in .env")
	}

	cfg.Mailgun.Domain = mailgun_domain

	mailgun_api_key := os.Getenv("MAILGUN_API_KEY")
	if mailgun_api_key == "" {
		log.Fatalf("MAILGUN_API_KEY not available in .env")
	}

	cfg.Mailgun.APIKey = mailgun_api_key

}
