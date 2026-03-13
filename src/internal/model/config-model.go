package model

type CorsConfig struct {
	AllowedOrigins []string `json:"allowed_origins"`
	AllowedMethods []string `json:"allowed_methods"`
	AllowedHeaders []string `json:"allowed_headers"`
}

type LogConfig struct {
	FilePath   string `json:"file_path"`
	MaxSizeMB  int    `json:"max_size_mb"`
	MaxBackups int    `json:"max_backups"`
	MaxAgeDays int    `json:"max_age_days"`
	Compress   bool   `json:"compress"`
}

type CacheRoute struct {
	Path string `json:"path"`
	Ttl  int    `json:"ttl_sec"` // in seconds
}
type CacheConfig struct {
	Routes   []CacheRoute `json:"routes"`
	Capacity int          `json:"capacity"`
}

type Config struct {
	ListenPort          int         `json:"listen_port"`
	Backend             string      `json:"backend"`
	MaxIdleConns        int         `json:"max_idle_conns"`
	MaxIdleConnsPerHost int         `json:"max_idle_conns_per_host"`
	FlushIntervalMs     int         `json:"flush_interval_ms"`
	IdleConnTimeoutSec  int         `json:"idle_conn_timeout_sec"`
	CORS                CorsConfig  `json:"cors"`
	Log                 LogConfig   `json:"log"`
	Cache               CacheConfig `json:"cache"`
}
