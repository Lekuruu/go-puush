package config

// DatabaseConfig holds configuration for database connections
type DatabaseConfig struct {
	Path string `env:"DB_PATH" envDefault:"./.data/puush.db"`

	// Connection Pool Settings
	MaxOpenConns    int `env:"DB_MAX_OPEN_CONNS" envDefault:"10"`
	MaxIdleConns    int `env:"DB_MAX_IDLE_CONNS" envDefault:"5"`
	ConnMaxLifetime int `env:"DB_CONN_MAX_LIFETIME" envDefault:"0"`

	// SQLite Specific Settings
	BusyTimeout int    `env:"DB_BUSY_TIMEOUT" envDefault:"5000"`
	JournalMode string `env:"DB_JOURNAL_MODE" envDefault:"WAL"`
	Synchronous string `env:"DB_SYNCHRONOUS" envDefault:"NORMAL"`
	CacheMode   string `env:"DB_CACHE_MODE" envDefault:"private"`

	// Performance Tuning
	CacheSize         int  `env:"DB_CACHE_SIZE" envDefault:"-2000"`
	EnableWAL         bool `env:"DB_ENABLE_WAL" envDefault:"true"`
	WALAutoCheckpoint int  `env:"DB_WAL_AUTOCHECKPOINT" envDefault:"1000"`
}
