package middleware

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Thashreef45/proxy-server/src/internal/model"
	"gopkg.in/natefinch/lumberjack.v2"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

type Logger struct {
	logger *log.Logger
	logChn chan model.RequestLog
}

func NewLogger(cfg model.LogConfig) *Logger {

	var output io.Writer

	if cfg.FilePath == "" {
		output = os.Stdout
	} else {

		output = &lumberjack.Logger{
			Filename:   cfg.FilePath,
			MaxSize:    cfg.MaxSizeMB,
			MaxAge:     cfg.MaxAgeDays,
			MaxBackups: cfg.MaxBackups,
			Compress:   cfg.Compress,
		}

	}

	// deafault buffer size for log channel
	bufferSize := 2000
	// default worker count for log processing
	workerCount := 5

	l := &Logger{
		logger: log.New(output, "", log.LstdFlags),
		logChn: make(chan model.RequestLog, bufferSize),
	}

	// spawn log keeper goroutines
	for range workerCount {
		go l.logKeeper()
	}

	return l
}

func (l *Logger) logKeeper() {
	// batch process may be needed for high traffic scenarios to reduce overhead of logging
	for log := range l.logChn {
		data, _ := json.Marshal(log)
		l.logger.Println(string(data))

		// can avoid json parsing (performance optimization)
		// l.logger.Println(log)
	}
}

func (l *Logger) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// log logic here

		lrw := &loggingResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		start := time.Now()

		// http.ResponseWriter
		next.ServeHTTP(lrw, r)

		duration := time.Since(start)

		reqLog := model.RequestLog{
			Timestamp:     time.Now().Format(time.RFC3339),
			ClientIP:      r.RemoteAddr,
			Method:        r.Method,
			Path:          r.URL.Path,
			Proto:         r.Proto,
			Status:        lrw.statusCode,
			DurationMs:    duration.Milliseconds(),
			UserAgent:     r.UserAgent(),
			ContentLength: r.ContentLength,
			Referer:       r.Referer(),
		}

		// log the request log to channel
		l.logChn <- reqLog

	})
}
