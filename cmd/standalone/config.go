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

# Email configuration
EMAIL_TYPE=noop
EMAIL_FROM=puush@puush.me

# SMTP configuration (used if EMAIL_TYPE is set to 'smtp')
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USERNAME=username
SMTP_PASSWORD=password
SMTP_USE_TLS=true
SMTP_SKIP_TLS_VERIFY=false

# Database configuration
DB_PATH=./.data/puush.db

# Service configuration
SERVICE_URL=http://localhost
SERVICE_NAME=puush
SERVICE_EMAIL=puush@puush.me

# Twitter handles used on the website
TWITTER_HANDLE=@puushme
TWITTER_URL=https://twitter.com/puushme

# Download location used on the front page
DOWNLOAD_WINDOWS=/dl/puush-win.zip
DOWNLOAD_MAC=/dl/puush.zip
DOWNLOAD_IOS=https://itunes.apple.com/au/app/puush/id386524126

# Enable or disable user registration
REGISTRATION_ENABLED=true

# Setting this to 'true' will require an invitation key for registration
REQUIRE_INVITATION=false
`

func CreateDefaultEnvironment() error {
	if _, err := os.Stat(".env"); !os.IsNotExist(err) {
		return nil
	}

	return os.WriteFile(".env", []byte(defaultEnvironment), 0644)
}
