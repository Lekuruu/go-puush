package database

// DatabaseConfig holds configuration for database connections
type DatabaseConfig struct {
	Path string `envconfig:"DB_PATH" default:"./.data/puush.db"`

	// Connection Pool Settings
	MaxOpenConns    int `envconfig:"DB_MAX_OPEN_CONNS" default:"10"`
	MaxIdleConns    int `envconfig:"DB_MAX_IDLE_CONNS" default:"5"`
	ConnMaxLifetime int `envconfig:"DB_CONN_MAX_LIFETIME" default:"3600"` // in seconds

	// SQLite Specific Settings
	BusyTimeout int    `envconfig:"DB_BUSY_TIMEOUT" default:"5000"`  // milliseconds
	JournalMode string `envconfig:"DB_JOURNAL_MODE" default:"WAL"`   // WAL, DELETE, TRUNCATE, PERSIST, MEMORY, OFF
	Synchronous string `envconfig:"DB_SYNCHRONOUS" default:"NORMAL"` // OFF, NORMAL, FULL, EXTRA
	CacheMode   string `envconfig:"DB_CACHE_MODE" default:"shared"`  // shared, private

	// Performance Tuning
	CacheSize         int  `envconfig:"DB_CACHE_SIZE" default:"-2000"` // negative = KB, positive = pages
	EnableWAL         bool `envconfig:"DB_ENABLE_WAL" default:"true"`
	WALAutoCheckpoint int  `envconfig:"DB_WAL_AUTOCHECKPOINT" default:"1000"` // pages
}
