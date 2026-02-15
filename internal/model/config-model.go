package model

type Config struct {
	ListenPort          int    `json:"listen_port"`
	Backend             string `json:"backend"`
	MaxIdleConns        int    `json:"max_idle_conns"`
	MaxIdleConnsPerHost int    `json:"max_idle_conns_per_host"`
	FlushIntervalMs     int    `json:"flush_interval_ms"`
	IdleConnTimeoutSec  int    `json:"idle_conn_timeout_sec"`
}
