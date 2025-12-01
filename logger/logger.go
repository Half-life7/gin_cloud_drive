package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// 日志级别
const (
	LevelDebug = "DEBUG"
	LevelInfo  = "INFO"
	LevelWarn  = "WARN"
	LevelError = "ERROR"
	LevelFatal = "FATAL"
)

// 日志类型
const (
	TypeFile   = "FILE"   // 文件操作
	TypeUser   = "USER"   // 用户操作
	TypeSystem = "SYSTEM" // 系统操作
	TypeAccess = "ACCESS" // 访问记录
	TypeError  = "ERROR"  // 错误记录
)

// LogEntry 日志条目结构体
type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`      // 时间戳
	Level     string    `json:"level"`          // 日志级别
	Type      string    `json:"type"`           // 日志类型
	IP        string    `json:"ip"`             // 客户端IP
	UserAgent string    `json:"user_agent"`     // 用户代理
	Action    string    `json:"action"`         // 操作内容
	Details   string    `json:"details"`        // 详细信息
	File      string    `json:"file,omitempty"` // 涉及的文件（可选）
	Size      int64     `json:"size,omitempty"` // 文件大小（可选）
}

// Logger 日志管理器
type Logger struct {
	logFile  *os.File
	mutex    sync.Mutex
	logLevel string
	logPath  string
}

// 全局日志实例
var (
	globalLogger *Logger
	once         sync.Once
)

// InitLogger 初始化日志管理器
func InitLogger(logPath string, logLevel string) error {
	var err error
	once.Do(func() {
		globalLogger = &Logger{
			logLevel: logLevel,
			logPath:  logPath,
		}
		err = globalLogger.openLogFile()
	})
	return err
}

// openLogFile 打开日志文件
func (l *Logger) openLogFile() error {
	// 确保日志目录存在
	if err := os.MkdirAll(l.logPath, 0755); err != nil {
		return fmt.Errorf("创建日志目录失败: %v", err)
	}

	// 日志文件名格式：logs_20060102.log
	logFileName := fmt.Sprintf("logs_%s.log", time.Now().Format("20060102"))
	logFilePath := fmt.Sprintf("%s/%s", l.logPath, logFileName)

	// 打开或创建日志文件
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("打开日志文件失败: %v", err)
	}

	l.logFile = file
	return nil
}

// writeLog 写入日志
func (l *Logger) writeLog(entry LogEntry) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	// 检查日志文件是否需要轮转（每天一个文件）
	currentDate := time.Now().Format("20060102")

	if l.logFile == nil {
		if err := l.openLogFile(); err != nil {
			return err
		}
	} else {
		// 检查当前日志文件是否是今天的
		fileInfo, err := l.logFile.Stat()
		if err != nil {
			return err
		}

		fileDate := fileInfo.Name()[5:13] // 提取文件名中的日期部分
		if fileDate != currentDate {
			// 关闭旧文件，打开新文件
			l.logFile.Close()
			if err := l.openLogFile(); err != nil {
				return err
			}
		}
	}

	// 转换为JSON格式
	logJSON, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("日志转换为JSON失败: %v", err)
	}

	// 写入文件
	_, err = l.logFile.WriteString(fmt.Sprintf("%s\n", string(logJSON)))
	if err != nil {
		return fmt.Errorf("写入日志文件失败: %v", err)
	}

	// 同时输出到控制台
	fmt.Printf("%s [%s] [%s] %s - %s\n",
		entry.Timestamp.Format("2006-01-02 15:04:05"),
		entry.Level,
		entry.Type,
		entry.IP,
		entry.Action,
	)

	return nil
}

// log 通用日志记录函数
func log(level, logType, ip, userAgent, action, details string, file string, size int64) error {
	if globalLogger == nil {
		// 如果未初始化，使用默认配置
		if err := InitLogger("./logs", LevelInfo); err != nil {
			return err
		}
	}

	// 检查日志级别
	if !shouldLog(globalLogger.logLevel, level) {
		return nil
	}

	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Type:      logType,
		IP:        ip,
		UserAgent: userAgent,
		Action:    action,
		Details:   details,
		File:      file,
		Size:      size,
	}

	return globalLogger.writeLog(entry)
}

// shouldLog 检查是否应该记录该级别的日志
func shouldLog(currentLevel, logLevel string) bool {
	levelOrder := map[string]int{
		LevelDebug: 0,
		LevelInfo:  1,
		LevelWarn:  2,
		LevelError: 3,
		LevelFatal: 4,
	}

	return levelOrder[logLevel] >= levelOrder[currentLevel]
}

// Debug 记录调试日志
func Debug(logType, ip, userAgent, action, details string, file string, size int64) error {
	return log(LevelDebug, logType, ip, userAgent, action, details, file, size)
}

// Info 记录信息日志
func Info(logType, ip, userAgent, action, details string, file string, size int64) error {
	return log(LevelInfo, logType, ip, userAgent, action, details, file, size)
}

// Warn 记录警告日志
func Warn(logType, ip, userAgent, action, details string, file string, size int64) error {
	return log(LevelWarn, logType, ip, userAgent, action, details, file, size)
}

// Error 记录错误日志
func Error(logType, ip, userAgent, action, details string, file string, size int64) error {
	return log(LevelError, logType, ip, userAgent, action, details, file, size)
}

// Fatal 记录致命错误日志
func Fatal(logType, ip, userAgent, action, details string, file string, size int64) error {
	if err := log(LevelFatal, logType, ip, userAgent, action, details, file, size); err != nil {
		return err
	}
	os.Exit(1)
	return nil
}

// LogFileOperation 记录文件操作日志
func LogFileOperation(ip, userAgent, action, filePath string, fileSize int64) error {
	return Info(TypeFile, ip, userAgent, action, "", filePath, fileSize)
}

// LogUserOperation 记录用户操作日志
func LogUserOperation(ip, userAgent, action, details string) error {
	return Info(TypeUser, ip, userAgent, action, details, "", 0)
}

// LogSystemOperation 记录系统操作日志
func LogSystemOperation(ip, userAgent, action, details string) error {
	return Info(TypeSystem, ip, userAgent, action, details, "", 0)
}

// LogAccess 记录访问日志
func LogAccess(ip, userAgent, action, details string) error {
	return Info(TypeAccess, ip, userAgent, action, details, "", 0)
}

// LogError 记录错误日志
func LogError(ip, userAgent, action, details string) error {
	return Error(TypeError, ip, userAgent, action, details, "", 0)
}
