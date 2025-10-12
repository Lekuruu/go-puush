package database

import (
	"crypto/md5"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type User struct {
	Id              int         `gorm:"primaryKey;autoIncrement;not null"`
	Name            string      `gorm:"size:16;not null"`
	Email           string      `gorm:"size:256;not null;unique"`
	Password        string      `gorm:"size:60;not null"`
	CreatedAt       time.Time   `gorm:"not null;CURRENT_TIMESTAMP"`
	LatestActivity  time.Time   `gorm:"not null;CURRENT_TIMESTAMP"`
	Active          bool        `gorm:"default:true;not null"`
	Type            AccountType `gorm:"not null;default:0"`
	ApiKey          string      `gorm:"size:64;not null;unique"`
	DiskUsage       int64       `gorm:"default:0;not null"`
	SubscriptionEnd *time.Time  `gorm:"default:NULL"`
	DefaultPoolId   int         `gorm:"default:NULL"`

	DefaultPool *Pool     `gorm:"foreignKey:DefaultPoolId;constraint:OnDelete:SET NULL"`
	Pools       []*Pool   `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE"`
	Uploads     []*Upload `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE"`
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

type Upload struct {
	Id        int       `gorm:"primaryKey;autoIncrement;not null"`
	UserId    int       `gorm:"not null"`
	PoolId    int       `gorm:"not null"`
	Filename  string    `gorm:"size:256;not null"`
	Filesize  int64     `gorm:"not null"`
	Checksum  string    `gorm:"size:32;not null"`
	CreatedAt time.Time `gorm:"not null;CURRENT_TIMESTAMP"`
	Views     int       `gorm:"default:0;not null"`
	MimeType  string    `gorm:"size:64;default:''"`

	User *User      `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE"`
	Pool *Pool      `gorm:"foreignKey:PoolId;constraint:OnDelete:CASCADE"`
	Link *ShortLink `gorm:"foreignKey:UploadId;constraint:OnDelete:CASCADE"`
}

func (upload *Upload) Key() string {
	return strconv.Itoa(upload.Id)
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

func (upload *Upload) SizeHumanReadable() string {
	return formatBytes(upload.Filesize)
}

func (upload *Upload) Url() string {
	if upload.Pool == nil {
		return ""
	}
	if upload.Pool.Type == PoolTypePasswordProtected && upload.Pool.Password != nil {
		return fmt.Sprintf("/%s/%s/%s", upload.Pool.Identifier, upload.Pool.PasswordHash(), upload.Filename)
	}
	return fmt.Sprintf("/%s/%s", upload.Pool.Identifier, upload.Filename)
}

func (upload *Upload) UrlEncoded() string {
	if upload.Pool == nil {
		return ""
	}
	if upload.Pool.Type == PoolTypePasswordProtected && upload.Pool.Password != nil {
		return fmt.Sprintf("/%s/%s/%s", upload.Pool.Identifier, upload.Pool.PasswordHash(), upload.FilenameEncoded())
	}
	return fmt.Sprintf("/%s/%s", upload.Pool.Identifier, upload.FilenameEncoded())
}

type Pool struct {
	Id         int       `gorm:"primaryKey;autoIncrement;not null"`
	UserId     int       `gorm:"not null"`
	Name       string    `gorm:"size:32;not null"`
	Identifier string    `gorm:"size:8;not null;unique"`
	Password   *string   `gorm:"size:32;default:NULL"`
	Type       PoolType  `gorm:"not null"`
	CreatedAt  time.Time `gorm:"not null;CURRENT_TIMESTAMP"`
	LastUpload time.Time `gorm:"not null;CURRENT_TIMESTAMP"`

	Uploads []*Upload `gorm:"foreignKey:PoolId;constraint:OnDelete:CASCADE"`
}

func (pool *Pool) PasswordHash() string {
	if pool.Password == nil {
		return "nil"
	}
	sum := md5.Sum([]byte(*pool.Password))
	return fmt.Sprintf("%x", sum)
}

type ShortLink struct {
	Identifier string     `gorm:"primaryKey;size:16;not null"`
	CreatedAt  time.Time  `gorm:"not null;CURRENT_TIMESTAMP"`
	ExpiresAt  *time.Time `gorm:"default:NULL"`
	UploadId   int        `gorm:"not null;unique"`

	Upload *Upload `gorm:"foreignKey:UploadId;constraint:OnDelete:CASCADE"`
}

func (shortlink *ShortLink) Url() string {
	return fmt.Sprintf("/%s", shortlink.Identifier)
}

func (shortlink *ShortLink) UrlEncoded() string {
	if shortlink.Upload == nil || shortlink.Upload.Pool == nil {
		return fmt.Sprintf("/%s", url.PathEscape(shortlink.Identifier))
	}
	if shortlink.Upload.Pool.Type == PoolTypePrivate {
		return fmt.Sprintf("/%s/%s", shortlink.Upload.Pool.Identifier, url.PathEscape(shortlink.Identifier))
	}
	return fmt.Sprintf("/%s", url.PathEscape(shortlink.Identifier))
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
