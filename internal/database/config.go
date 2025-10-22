package database

// DatabaseConfig holds configuration for database connections
type DatabaseConfig struct {
	Path string `env:"DB_PATH" envDefault:"./.data/puush.db"`

	// Connection Pool Settings
	MaxOpenConns    int `env:"DB_MAX_OPEN_CONNS" envDefault:"10"`
	MaxIdleConns    int `env:"DB_MAX_IDLE_CONNS" envDefault:"5"`
	ConnMaxLifetime int `env:"DB_CONN_MAX_LIFETIME" envDefault:"3600"` // in seconds

	// SQLite Specific Settings
	BusyTimeout int    `env:"DB_BUSY_TIMEOUT" envDefault:"5000"`  // milliseconds
	JournalMode string `env:"DB_JOURNAL_MODE" envDefault:"WAL"`   // WAL, DELETE, TRUNCATE, PERSIST, MEMORY, OFF
	Synchronous string `env:"DB_SYNCHRONOUS" envDefault:"NORMAL"` // OFF, NORMAL, FULL, EXTRA
	CacheMode   string `env:"DB_CACHE_MODE" envDefault:"shared"`  // shared, private

	// Performance Tuning
	CacheSize         int  `env:"DB_CACHE_SIZE" envDefault:"-2000"` // negative = KB, positive = pages
	EnableWAL         bool `env:"DB_ENABLE_WAL" envDefault:"true"`
	WALAutoCheckpoint int  `env:"DB_WAL_AUTOCHECKPOINT" envDefault:"1000"` // pages
}
