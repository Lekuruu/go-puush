package database

import (
	"crypto/md5"
	"fmt"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type User struct {
	Id                    int         `gorm:"primaryKey;autoIncrement;not null"`
	Name                  string      `gorm:"size:16;not null"`
	Email                 string      `gorm:"size:256;not null;unique"`
	Password              string      `gorm:"size:60;not null"`
	CreatedAt             time.Time   `gorm:"not null;CURRENT_TIMESTAMP"`
	LatestActivity        time.Time   `gorm:"not null;CURRENT_TIMESTAMP"`
	Active                bool        `gorm:"default:true;not null"`
	Type                  AccountType `gorm:"not null;default:0"`
	ApiKey                string      `gorm:"size:64;not null;unique"`
	DiskUsage             int64       `gorm:"default:0;not null"`
	ViewType              ViewType    `gorm:"size:16;default:'list';not null"`
	SubscriptionEnd       *time.Time  `gorm:"default:NULL"`
	DefaultPoolId         int         `gorm:"default:NULL"`
	UsernameSetupReminder bool        `gorm:"default:true;not null"`

	DefaultPool *Pool     `gorm:"foreignKey:DefaultPoolId;constraint:OnDelete:SET NULL"`
	Pools       []*Pool   `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE"`
	Uploads     []*Upload `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE"`
}

func (user *User) DiskUsageHumanReadable() string {
	return formatBytes(user.DiskUsage)
}

func (user *User) UploadLimit() int64 {
	switch user.Type {
	case AccountTypeRegular:
		return UploadLimitRegular
	case AccountTypePro:
		return UploadLimitPro
	default:
		return -1
	}
}

func (user *User) DisplayName() string {
	if !user.RequiresUsernameSetup() {
		return user.Name
	}
	return user.Email
}

func (user *User) RequiresUsernameSetup() bool {
	return user.Name == ""
}

type Upload struct {
	Id         int       `gorm:"primaryKey;autoIncrement;not null"`
	UserId     int       `gorm:"not null"`
	PoolId     int       `gorm:"not null"`
	Identifier string    `gorm:"size:16;not null;index"`
	Filename   string    `gorm:"size:256;not null"`
	Filesize   int64     `gorm:"not null"`
	Checksum   string    `gorm:"size:32;not null"`
	CreatedAt  time.Time `gorm:"not null;CURRENT_TIMESTAMP"`
	Views      int       `gorm:"default:0;not null"`
	MimeType   string    `gorm:"size:64;default:''"`

	User *User `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE"`
	Pool *Pool `gorm:"foreignKey:PoolId;constraint:OnDelete:CASCADE"`
}

func (upload *Upload) Key() string {
	return strconv.Itoa(upload.Id)
}

func (upload *Upload) FilenameExtension() string {
	return filepath.Ext(upload.Filename)
}

func (upload *Upload) FilenameEncoded() string {
	return url.PathEscape(upload.Filename)
}

func (upload *Upload) IsImage() bool {
	return strings.HasPrefix(upload.MimeType, "image")
}

func (upload *Upload) IsVideo() bool {
	return strings.HasPrefix(upload.MimeType, "video")
}

func (upload *Upload) IsAudio() bool {
	return strings.HasPrefix(upload.MimeType, "audio")
}

func (upload *Upload) SizeHumanReadable() string {
	return formatBytes(upload.Filesize)
}

func (upload *Upload) Url() string {
	if upload.Pool != nil && upload.Pool.Type == PoolTypePrivate {
		return fmt.Sprintf("/%s/%s", upload.Pool.Identifier, upload.Identifier) + upload.FilenameExtension()
	}
	return fmt.Sprintf("/%s", upload.Identifier) + upload.FilenameExtension()
}

func (upload *Upload) UrlEncoded() string {
	if upload.Pool != nil && upload.Pool.Type == PoolTypePrivate {
		return fmt.Sprintf("/%s/%s", url.PathEscape(upload.Pool.Identifier), url.PathEscape(upload.Identifier)) + upload.FilenameExtension()
	}
	return fmt.Sprintf("/%s", url.PathEscape(upload.Identifier)) + upload.FilenameExtension()
}

type Pool struct {
	Id          int       `gorm:"primaryKey;autoIncrement;not null"`
	UserId      int       `gorm:"not null;index"`
	Name        string    `gorm:"size:32;not null"`
	Identifier  string    `gorm:"size:8;not null;unique"`
	Password    *string   `gorm:"size:32;default:NULL"`
	Type        PoolType  `gorm:"not null"`
	CreatedAt   time.Time `gorm:"not null;CURRENT_TIMESTAMP"`
	LastUpload  time.Time `gorm:"not null;CURRENT_TIMESTAMP"`
	UploadCount int       `gorm:"default:0;not null"`

	Uploads []*Upload `gorm:"foreignKey:PoolId;constraint:OnDelete:CASCADE"`
	User    *User     `gorm:"foreignKey:UserId;references:Id;constraint:OnDelete:CASCADE"`
}

func (pool *Pool) PasswordHash() string {
	if pool.Password == nil {
		return "nil"
	}
	sum := md5.Sum([]byte(*pool.Password))
	return fmt.Sprintf("%x", sum)
}

func (pool *Pool) UploadIdentifierLength() int {
	switch pool.Type {
	case PoolTypePrivate:
		return 10
	case PoolTypePasswordProtected:
		return 10
	default:
		return 6
	}
}

type Session struct {
	Id        uint      `gorm:"primaryKey"`
	Token     string    `gorm:"uniqueIndex;not null"`
	UserId    int       `gorm:"index;not null"`
	ExpiresAt time.Time `gorm:"index;not null"`
	CreatedAt time.Time `gorm:"not null;CURRENT_TIMESTAMP"`

	User *User `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE"`
}

func (session *Session) IsExpired() bool {
	return time.Now().After(session.ExpiresAt)
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%dB", bytes)
	}

	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	units := []string{"KB", "MB", "GB", "TB", "PB", "EB"}
	return fmt.Sprintf("%.2f%s", float64(bytes)/float64(div), units[exp])
}

type InvitationKey struct {
	Id        int        `gorm:"primaryKey;autoIncrement;not null"`
	Key       string     `gorm:"size:16;not null;unique;index"`
	CreatedAt time.Time  `gorm:"not null;CURRENT_TIMESTAMP"`
	ExpiresAt *time.Time `gorm:"default:NULL"`
}

func (key *InvitationKey) IsExpired() bool {
	if key.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*key.ExpiresAt)
}

type EmailVerification struct {
	Id        int                     `gorm:"primaryKey;autoIncrement;not null"`
	Key       string                  `gorm:"size:32;not null;unique;index"`
	Action    EmailVerificationAction `gorm:"size:32;not null"`
	UserId    *int                    `gorm:"default:NULL"`
	CreatedAt time.Time               `gorm:"not null;CURRENT_TIMESTAMP"`
	ExpiresAt *time.Time              `gorm:"default:NULL"`

	User *User `gorm:"foreignKey:UserId;constraint:OnDelete:SET NULL"`
}

func (verification *EmailVerification) IsExpired() bool {
	if verification.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*verification.ExpiresAt)
}
