package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

// Config holds the application configuration loaded from the environment.
type Config struct {
	Cdn      CdnConfig
	Web      WebConfig
	Storage  StorageConfig
	Email    EmailConfig
	Service  ServiceConfig
	Database DatabaseConfig
}

type CdnConfig struct {
	Url            string `env:"CDN_URL" envDefault:"http://puu.sh"`
	RateLimit      int    `env:"CDN_RATE_LIMIT" envDefault:"10"`
	RateLimitBurst int    `env:"CDN_RATE_LIMIT_BURST" envDefault:"30"`
}

type WebConfig struct {
	Host string `env:"WEB_HOST" envDefault:"0.0.0.0"`
	Port int    `env:"WEB_PORT" envDefault:"8082"`
}

type StorageConfig struct {
	Type string `env:"STORAGE_TYPE" envDefault:"local"`
	Uri  string `env:"STORAGE_URI" envDefault:"./.data/"`
}

type ServiceConfig struct {
	Url                 string `env:"SERVICE_URL" envDefault:"http://puush.me"`
	Name                string `env:"SERVICE_NAME" envDefault:"puush"`
	Email               string `env:"SERVICE_EMAIL" envDefault:"puush@puush.me"`
	TwitterHandle       string `env:"TWITTER_HANDLE" envDefault:"@puushme"`
	TwitterUrl          string `env:"TWITTER_URL" envDefault:"https://twitter.com/puushme"`
	DownloadWindows     string `env:"DOWNLOAD_WINDOWS" envDefault:"/dl/puush-win.zip"`
	DownloadMac         string `env:"DOWNLOAD_MAC" envDefault:"/dl/puush.zip"`
	DownloadIOS         string `env:"DOWNLOAD_IOS" envDefault:"https://itunes.apple.com/au/app/puush/id386524126"`
	RegistrationEnabled bool   `env:"REGISTRATION_ENABLED" envDefault:"true"`
	RequireActivation   bool   `env:"REQUIRE_ACTIVATION" envDefault:"false"`
	RequireInvitation   bool   `env:"REQUIRE_INVITATION" envDefault:"false"`
	ExtendedGallery     bool   `env:"EXTENDED_GALLERY" envDefault:"false"`
}

// LoadConfig loads dotenv files, then parses the process environment.
func LoadConfig(environmentFiles ...string) (*Config, error) {
	// Try to apply .env file if it exists
	if len(environmentFiles) == 0 {
		environmentFiles = []string{".env"}
	}
	for _, file := range environmentFiles {
		godotenv.Load(file)
	}

	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return nil, fmt.Errorf("config: failed to parse environment: %w", err)
	}
	return &cfg, nil
}
