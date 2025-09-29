package app

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Api struct {
		Host string `envconfig:"API_HOST" default:"0.0.0.0"`
		Port int    `envconfig:"API_PORT" default:"8080"`
	}
	Cdn struct {
		Host string `envconfig:"CDN_HOST" default:"0.0.0.0"`
		Port int    `envconfig:"CDN_PORT" default:"8081"`
		Url  string `envconfig:"CDN_URL" default:"http://puu.sh"`
	}
	Database struct {
		Path string `envconfig:"DB_PATH" default:"./.data/puush.db"`
	}
	Storage struct {
		Type string `envconfig:"STORAGE_TYPE" default:"local"`
		Uri  string `envconfig:"STORAGE_URI" default:"./.data/"`
	}
	Service struct {
		Url                 string `envconfig:"SERVICE_URL" default:"http://puush.me"`
		Name                string `envconfig:"SERVICE_NAME" default:"puush"`
		Email               string `envconfig:"SERVICE_EMAIL" default:"puush@puush.me"`
		TwitterHandle       string `envconfig:"TWITTER_HANDLE" default:"@puushme"`
		TwitterUrl          string `envconfig:"TWITTER_URL" default:"https://twitter.com/puushme"`
		DownloadWindows     string `envconfig:"DOWNLOAD_WINDOWS" default:"/dl/puush-win.zip"`
		DownloadMac         string `envconfig:"DOWNLOAD_MAC" default:"/dl/puush.zip"`
		DownloadIOS         string `envconfig:"DOWNLOAD_IOS" default:"https://itunes.apple.com/au/app/puush/id386524126"`
		RegistrationEnabled bool   `envconfig:"REGISTRATION_ENABLED" default:"true"`
	}
}

func LoadConfig() (*Config, error) {
	var config Config
	if err := envconfig.Process("", &config); err != nil {
		return nil, err
	}
	return &config, nil
}
