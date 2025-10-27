package app

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"

	"github.com/Lekuruu/go-puush/internal/database"
)

type Config struct {
	Api struct {
		Host string `env:"API_HOST" envDefault:"0.0.0.0"`
		Port int    `env:"API_PORT" envDefault:"8080"`
	}
	Cdn struct {
		Host string `env:"CDN_HOST" envDefault:"0.0.0.0"`
		Port int    `env:"CDN_PORT" envDefault:"8081"`
		Url  string `env:"CDN_URL" envDefault:"http://puu.sh"`
	}
	Web struct {
		Host string `env:"WEB_HOST" envDefault:"0.0.0.0"`
		Port int    `env:"WEB_PORT" envDefault:"8082"`
	}
	Storage struct {
		Type string `env:"STORAGE_TYPE" envDefault:"local"`
		Uri  string `env:"STORAGE_URI" envDefault:"./.data/"`
	}
	Email struct {
		Type string `env:"EMAIL_TYPE" envDefault:"noop"`
		From string `env:"EMAIL_FROM"`
	}
	Service struct {
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
	}
	Database database.DatabaseConfig
}

func LoadConfig() (*Config, error) {
	// Try to apply .env file if it exists
	godotenv.Load()

	var config Config
	if err := env.Parse(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
