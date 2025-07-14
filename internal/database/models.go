package database

import "time"

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

	Uploads []*Upload `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE"`
	Pools   []*Pool   `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE"`
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

type Pool struct {
	Id         int       `gorm:"primaryKey;autoIncrement;not null"`
	UserId     int       `gorm:"not null"`
	Name       string    `gorm:"size:32;not null"`
	Type       PoolType  `gorm:"not null"`
	Password   *string   `gorm:"size:32;default:NULL"`
	CreatedAt  time.Time `gorm:"not null;CURRENT_TIMESTAMP"`
	LastUpload time.Time `gorm:"not null;CURRENT_TIMESTAMP"`

	Uploads []*Upload `gorm:"foreignKey:PoolId;constraint:OnDelete:CASCADE"`
}
