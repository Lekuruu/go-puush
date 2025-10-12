package database

type AccountType int8

const (
	AccountTypeRegular AccountType = iota
	AccountTypePro
	AccountTypeUnlimited
)

func (at AccountType) String() string {
	switch at {
	case AccountTypeRegular:
		return "Free"
	case AccountTypePro:
		return "Pro"
	case AccountTypeUnlimited:
		return "Haxxor"
	default:
		return "Unknown"
	}
}

type PoolType int8

const (
	PoolTypePublic PoolType = iota
	PoolTypePrivate
	PoolTypePasswordProtected
	PoolTypeGallery
)

func (pt PoolType) String() string {
	switch pt {
	case PoolTypePublic:
		return "Public"
	case PoolTypePrivate:
		return "Private"
	case PoolTypePasswordProtected:
		return "Password Protected"
	case PoolTypeGallery:
		return "Gallery"
	default:
		return "Unknown"
	}
}

const (
	UploadLimitRegular = 200 * 1024 * 1024       // 200 MB
	UploadLimitPro     = 15 * 1000 * 1024 * 1024 // 15 GB
)
