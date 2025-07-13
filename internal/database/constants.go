package database

type PoolType int8

const (
	PoolTypePublic PoolType = iota
	PoolTypePrivate
	PoolTypePasswordProtected
	PoolTypeGallery
)
