package database

import (
	"crypto/md5"
	"fmt"
	"strconv"
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

	User *User `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE"`
	Pool *Pool `gorm:"foreignKey:PoolId;constraint:OnDelete:CASCADE"`
}

func (upload *Upload) Key() string {
	return strconv.Itoa(upload.Id)
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
