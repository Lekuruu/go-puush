package database

type AccountType int8

const (
	AccountTypeRegular AccountType = iota
	AccountTypePro
	AccountTypeUnlimited
)

type PoolType int8

const (
	PoolTypePublic PoolType = iota
	PoolTypePrivate
	PoolTypePasswordProtected
	PoolTypeGallery
)
