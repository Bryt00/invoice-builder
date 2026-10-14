package systemlog

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// LogLevel represents severity of log entry.
type LogLevel string

const (
	LevelInfo  LogLevel = "INFO"
	LevelWarn  LogLevel = "WARN"
	LevelError LogLevel = "ERROR"
	LevelFatal LogLevel = "FATAL"
)

// SystemLog represents a single structured log entry stored in SQLite.
type SystemLog struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Timestamp  time.Time `gorm:"index:idx_timestamp" json:"timestamp"`
	Level      string    `gorm:"type:varchar(10);index:idx_level" json:"level"`
	Source     string    `gorm:"type:varchar(50);index:idx_source" json:"source"` // HTTP, AUTH, PAYMENT, DB, SYSTEM
	Message    string    `gorm:"type:text" json:"message"`
	Method     string    `gorm:"type:varchar(10)" json:"method,omitempty"`
	Path       string    `gorm:"type:varchar(255);index:idx_path" json:"path,omitempty"`
	StatusCode int       `gorm:"index:idx_status" json:"status_code,omitempty"`
	LatencyMs  int64     `json:"latency_ms,omitempty"`
	IPAddress  string    `gorm:"type:varchar(50)" json:"ip_address,omitempty"`
	UserID     string    `gorm:"type:varchar(50)" json:"user_id,omitempty"`
	Details    string    `gorm:"type:text" json:"details,omitempty"`
}

// Filter represents query options for searching system logs.
type Filter struct {
	Level      string
	Source     string
	Search     string
	StatusCode int
	StatusType string // "error" (>=400), "success" (<400)
	Page       int
	Limit      int
}

// Stats represents high-level metrics for the logs console.
type Stats struct {
	TotalLogs     int64   `json:"total_logs"`
	ErrorCount24h int64   `json:"error_count_24h"`
	WarnCount24h  int64   `json:"warn_count_24h"`
	InfoCount24h  int64   `json:"info_count_24h"`
	AvgLatencyMs  float64 `json:"avg_latency_ms"`
	DBSizeBytes   int64   `json:"db_size_bytes"`
}

// Logger handles asynchronous SQLite logging.
type Logger struct {
	db       *gorm.DB
	logChan  chan *SystemLog
	dbPath   string
	wg       sync.WaitGroup
	ctx      context.Context
	cancel   context.CancelFunc
	maxCount int // Maximum number of logs to retain (e.g. 30,000)
}

var (
	defaultLogger *Logger
	once          sync.Once
)

// Init initializes the global structured SQLite logger.
func Init(dbFilePath string) (*Logger, error) {
	var initErr error
	once.Do(func() {
		if dbFilePath == "" {
			dbFilePath = "logs/system_logs.db"
		}

		dir := filepath.Dir(dbFilePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			initErr = fmt.Errorf("failed to create log dir: %w", err)
			return
		}

		// Open SQLite with GORM, tuning WAL mode and busy timeout for high-concurrency writes
		dsn := fmt.Sprintf("%s?_journal_mode=WAL&_busy_timeout=5000&_synchronous=NORMAL", dbFilePath)
		db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		if err != nil {
			initErr = fmt.Errorf("failed to open sqlite log db: %w", err)
			return
		}

		// Configure pool: SQLite works best with 1 writer or limited pool
		sqlDB, err := db.DB()
		if err == nil {
			sqlDB.SetMaxOpenConns(1)
			sqlDB.SetMaxIdleConns(1)
			sqlDB.SetConnMaxLifetime(time.Hour)
		}

		// AutoMigrate schema
		if err := db.AutoMigrate(&SystemLog{}); err != nil {
			initErr = fmt.Errorf("failed to migrate system_logs table: %w", err)
			return
		}

		ctx, cancel := context.WithCancel(context.Background())
		l := &Logger{
			db:       db,
			logChan:  make(chan *SystemLog, 10000), // Buffered channel prevents blocking callers
			dbPath:   dbFilePath,
			ctx:      ctx,
			cancel:   cancel,
			maxCount: 30000,
		}

		// Start background log consumer worker
		l.wg.Add(1)
		go l.worker()

		// Start background log retention pruner (every 1 hour)
		l.wg.Add(1)
		go l.retentionWorker()

		defaultLogger = l
	})

	return defaultLogger, initErr
}

// Get returns the global Logger instance.
func Get() *Logger {
	return defaultLogger
}

// worker batches and writes incoming log records to SQLite.
func (l *Logger) worker() {
	defer l.wg.Done()

	batch := make([]*SystemLog, 0, 100)
	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()

	flush := func() {
		if len(batch) == 0 {
			return
		}
		if err := l.db.Create(&batch).Error; err != nil {
			log.Printf("[SystemLog] failed to write batch: %v", err)
		}
		batch = batch[:0]
	}

	for {
		select {
		case <-l.ctx.Done():
			// Drain remaining logs before exit
			for {
				select {
				case entry := <-l.logChan:
					batch = append(batch, entry)
					if len(batch) >= 100 {
						flush()
					}
				default:
					flush()
					return
				}
			}

		case entry := <-l.logChan:
			batch = append(batch, entry)
			if len(batch) >= 50 {
				flush()
			}

		case <-ticker.C:
			flush()
		}
	}
}

// retentionWorker automatically purges logs older than 14 days or beyond maxCount.
func (l *Logger) retentionWorker() {
	defer l.wg.Done()
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	prune := func() {
		cutoff := time.Now().Add(-14 * 24 * time.Hour)
		l.db.Where("timestamp < ?", cutoff).Delete(&SystemLog{})

		// Also cap total rows to maxCount if table grows too large
		var count int64
		l.db.Model(&SystemLog{}).Count(&count)
		if count > int64(l.maxCount) {
			excess := count - int64(l.maxCount)
			subQuery := l.db.Model(&SystemLog{}).Select("id").Order("id asc").Limit(int(excess))
			l.db.Where("id IN (?)", subQuery).Delete(&SystemLog{})
		}
	}

	// Run initial prune
	prune()

	for {
		select {
		case <-l.ctx.Done():
			return
		case <-ticker.C:
			prune()
		}
	}
}

// Close gracefully flushes pending logs and closes the database.
func (l *Logger) Close() {
	if l == nil {
		return
	}
	l.cancel()
	l.wg.Wait()
	if sqlDB, err := l.db.DB(); err == nil {
		sqlDB.Close()
	}
}

// Log queues a new structured log entry without blocking the caller.
func (l *Logger) Log(entry *SystemLog) {
	if l == nil {
		return
	}
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now()
	}

	// Always print to stdout for standard container / terminal monitoring
	fmt.Printf("[%s] %s | %s | %s\n",
		entry.Timestamp.Format("15:04:05"),
		entry.Level,
		entry.Source,
		entry.Message,
	)

	select {
	case l.logChan <- entry:
	default:
		// Drop if buffer full to guarantee zero backpressure on HTTP serving
		log.Printf("[SystemLog] queue full, dropped entry: %s", entry.Message)
	}
}

// Convenience log functions
func (l *Logger) Info(source, message string, details ...any) {
	l.Log(&SystemLog{
		Level:   string(LevelInfo),
		Source:  source,
		Message: message,
		Details: formatDetails(details...),
	})
}

func (l *Logger) Warn(source, message string, details ...any) {
	l.Log(&SystemLog{
		Level:   string(LevelWarn),
		Source:  source,
		Message: message,
		Details: formatDetails(details...),
	})
}

func (l *Logger) Error(source, message string, err error, details ...any) {
	msg := message
	if err != nil {
		msg = fmt.Sprintf("%s: %v", message, err)
	}
	l.Log(&SystemLog{
		Level:   string(LevelError),
		Source:  source,
		Message: msg,
		Details: formatDetails(details...),
	})
}

func (l *Logger) LogHTTP(method, path string, statusCode int, latency time.Duration, ip, userID string, details ...any) {
	level := string(LevelInfo)
	if statusCode >= 500 {
		level = string(LevelError)
	} else if statusCode >= 400 {
		level = string(LevelWarn)
	}

	l.Log(&SystemLog{
		Level:      level,
		Source:     "HTTP",
		Method:     method,
		Path:       path,
		StatusCode: statusCode,
		LatencyMs:  latency.Milliseconds(),
		IPAddress:  ip,
		UserID:     userID,
		Message:    fmt.Sprintf("%s %s -> %d (%dms)", method, path, statusCode, latency.Milliseconds()),
		Details:    formatDetails(details...),
	})
}

func formatDetails(details ...any) string {
	if len(details) == 0 {
		return ""
	}
	if len(details) == 1 {
		if s, ok := details[0].(string); ok {
			return s
		}
		if b, err := json.MarshalIndent(details[0], "", "  "); err == nil {
			return string(b)
		}
	}
	b, err := json.MarshalIndent(details, "", "  ")
	if err != nil {
		return fmt.Sprintf("%v", details)
	}
	return string(b)
}

// Query Methods for the Admin Logging Console

func (l *Logger) GetLogs(filter Filter) ([]SystemLog, int64, error) {
	var logs []SystemLog
	var total int64

	q := l.db.Model(&SystemLog{})

	if filter.Level != "" && filter.Level != "ALL" {
		q = q.Where("level = ?", strings.ToUpper(filter.Level))
	}
	if filter.Source != "" && filter.Source != "ALL" {
		q = q.Where("source = ?", strings.ToUpper(filter.Source))
	}
	if filter.StatusCode > 0 {
		q = q.Where("status_code = ?", filter.StatusCode)
	}
	if filter.StatusType == "error" {
		q = q.Where("status_code >= 400")
	} else if filter.StatusType == "success" {
		q = q.Where("status_code > 0 AND status_code < 400")
	}
	if filter.Search != "" {
		s := "%" + strings.ToLower(filter.Search) + "%"
		q = q.Where("LOWER(message) LIKE ? OR LOWER(path) LIKE ? OR LOWER(ip_address) LIKE ? OR LOWER(details) LIKE ?", s, s, s, s)
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	page := filter.Page
	if page < 1 {
		page = 1
	}
	limit := filter.Limit
	if limit < 1 || limit > 200 {
		limit = 50
	}
	offset := (page - 1) * limit

	err := q.Order("id desc").Offset(offset).Limit(limit).Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

func (l *Logger) GetStats() (Stats, error) {
	var stats Stats

	// Total logs
	l.db.Model(&SystemLog{}).Count(&stats.TotalLogs)

	// Past 24 hours counts
	yesterday := time.Now().Add(-24 * time.Hour)
	l.db.Model(&SystemLog{}).Where("timestamp >= ? AND level = ?", yesterday, LevelError).Count(&stats.ErrorCount24h)
	l.db.Model(&SystemLog{}).Where("timestamp >= ? AND level = ?", yesterday, LevelWarn).Count(&stats.WarnCount24h)
	l.db.Model(&SystemLog{}).Where("timestamp >= ? AND level = ?", yesterday, LevelInfo).Count(&stats.InfoCount24h)

	// Average latency for HTTP requests in past 24h
	var avgLatency struct {
		Avg float64
	}
	l.db.Model(&SystemLog{}).
		Select("COALESCE(AVG(latency_ms), 0) as avg").
		Where("timestamp >= ? AND source = 'HTTP' AND latency_ms > 0", yesterday).
		Scan(&avgLatency)
	stats.AvgLatencyMs = avgLatency.Avg

	// SQLite file size
	if fi, err := os.Stat(l.dbPath); err == nil {
		stats.DBSizeBytes = fi.Size()
	}

	return stats, nil
}

func (l *Logger) ClearLogs() error {
	return l.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&SystemLog{}).Error
}
