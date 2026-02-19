package model

type RequestLog struct {
	Timestamp     string `json:"timestamp"`
	ClientIP      string `json:"client_ip"`
	Method        string `json:"method"`
	Path          string `json:"path"`
	Proto         string `json:"proto"`
	Status        int    `json:"status"`
	DurationMs    int64  `json:"duration_ms"`
	UserAgent     string `json:"user_agent"`
	ContentLength int64  `json:"content_length"`
	Referer       string `json:"referer"`
}
