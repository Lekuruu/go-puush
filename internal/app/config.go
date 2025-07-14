package app

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Api struct {
		Host                  string `envconfig:"API_HOST" default:"0.0.0.0"`
		Port                  int    `envconfig:"API_PORT" default:"8080"`
		ClientVersion         int    `envconfig:"API_CLIENT_VERSION" default:"93"`
		ClientDownloadWindows string `envconfig:"API_CLIENT_WINDOWS" default:"http://puush.me/dl/puush-win.zip"`
		ClientDownloadMacOS   string `envconfig:"API_CLIENT_MACOS" default:"http://puush.me/dl/puush.zip"`
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
}

func LoadConfig() (*Config, error) {
	var config Config
	if err := envconfig.Process("", &config); err != nil {
		return nil, err
	}
	return &config, nil
}
