package main

import "os"

var defaultEnvironment = `
# Web configuration
WEB_HOST=127.0.0.1
WEB_PORT=80
CDN_URL=http://localhost

# Storage configuration
STORAGE_TYPE=local
STORAGE_URI=./.data

# Database configuration
DB_PATH=./.data/puush.db

# Service configuration
SERVICE_URL=http://localhost
SERVICE_NAME=puush
SERVICE_EMAIL=puush@puush.me
TWITTER_HANDLE=@puushme
TWITTER_URL=https://twitter.com/puushme
DOWNLOAD_WINDOWS=/dl/puush-win.zip
DOWNLOAD_MAC=/dl/puush.zip
DOWNLOAD_IOS=https://itunes.apple.com/au/app/puush/id386524126
REGISTRATION_ENABLED=true
`

func CreateDefaultEnvironment() error {
	if _, err := os.Stat(".env"); !os.IsNotExist(err) {
		return nil
	}

	return os.WriteFile(".env", []byte(defaultEnvironment), 0644)
}
