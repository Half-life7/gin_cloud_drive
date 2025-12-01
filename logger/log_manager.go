package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// LogQueryParams 日志查询参数
type LogQueryParams struct {
	StartDate time.Time `json:"start_date"` // 开始日期
	EndDate   time.Time `json:"end_date"`   // 结束日期
	Level     string    `json:"level"`      // 日志级别（可选）
	Type      string    `json:"type"`       // 日志类型（可选）
	IP        string    `json:"ip"`         // IP地址（可选）
	Action    string    `json:"action"`     // 操作内容（可选）
	File      string    `json:"file"`       // 文件名（可选）
	Page      int       `json:"page"`       // 页码
	PageSize  int       `json:"page_size"`  // 每页大小
}

// LogQueryResult 日志查询结果
type LogQueryResult struct {
	Total    int64      `json:"total"`     // 总记录数
	Page     int        `json:"page"`      // 当前页码
	PageSize int        `json:"page_size"` // 每页大小
	Logs     []LogEntry `json:"logs"`      // 日志列表
}

// QueryLogs 查询日志
func QueryLogs(params LogQueryParams) (*LogQueryResult, error) {
	// 默认值处理
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 20
	}
	if params.EndDate.IsZero() {
		params.EndDate = time.Now()
	}
	if params.StartDate.IsZero() {
		params.StartDate = params.EndDate.AddDate(0, 0, -7) // 默认查询最近7天
	}

	// 获取日志目录
	logPath := "./logs"
	if globalLogger != nil {
		logPath = globalLogger.logPath
	}

	// 读取所有日志文件
	files, err := os.ReadDir(logPath)
	if err != nil {
		return nil, fmt.Errorf("读取日志目录失败: %v", err)
	}

	// 筛选出日期范围内的日志文件
	var relevantFiles []string
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		filename := file.Name()
		if len(filename) < 13 || filename[:5] != "logs_" || filename[13:] != ".log" {
			continue
		}
		// 解析文件名中的日期
		fileDate, err := time.Parse("20060102", filename[5:13])
		if err != nil {
			continue
		}
		// 检查是否在查询日期范围内
		if (fileDate.Equal(params.StartDate) || fileDate.After(params.StartDate)) &&
			(fileDate.Equal(params.EndDate) || fileDate.Before(params.EndDate.AddDate(0, 0, 1))) {
			relevantFiles = append(relevantFiles, filepath.Join(logPath, filename))
		}
	}

	// 按日期排序，最新的文件排在前面
	sort.Slice(relevantFiles, func(i, j int) bool {
		return relevantFiles[i] > relevantFiles[j]
	})

	// 读取并筛选日志
	var allLogs []LogEntry
	for _, filePath := range relevantFiles {
		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}

		// 按行解析日志
		lines := splitLines(string(fileContent))
		for _, line := range lines {
			if line == "" {
				continue
			}

			var entry LogEntry
			if err := json.Unmarshal([]byte(line), &entry); err != nil {
				continue
			}

			// 筛选条件
			if !matchLogEntry(entry, params) {
				continue
			}

			allLogs = append(allLogs, entry)
		}
	}

	// 按时间排序，最新的日志排在前面
	sort.Slice(allLogs, func(i, j int) bool {
		return allLogs[i].Timestamp.After(allLogs[j].Timestamp)
	})

	// 分页处理
	total := int64(len(allLogs))
	start := (params.Page - 1) * params.PageSize
	end := start + params.PageSize
	if start >= len(allLogs) {
		return &LogQueryResult{
			Total:    total,
			Page:     params.Page,
			PageSize: params.PageSize,
			Logs:     []LogEntry{},
		}, nil
	}
	if end > len(allLogs) {
		end = len(allLogs)
	}

	return &LogQueryResult{
		Total:    total,
		Page:     params.Page,
		PageSize: params.PageSize,
		Logs:     allLogs[start:end],
	}, nil
}

// matchLogEntry 匹配日志条目
func matchLogEntry(entry LogEntry, params LogQueryParams) bool {
	// 时间范围匹配
	if entry.Timestamp.Before(params.StartDate) || entry.Timestamp.After(params.EndDate) {
		return false
	}

	// 日志级别匹配
	if params.Level != "" && entry.Level != params.Level {
		return false
	}

	// 日志类型匹配
	if params.Type != "" && entry.Type != params.Type {
		return false
	}

	// IP地址匹配
	if params.IP != "" && entry.IP != params.IP {
		return false
	}

	// 操作内容匹配
	if params.Action != "" && entry.Action != params.Action {
		return false
	}

	// 文件名匹配
	if params.File != "" && entry.File != params.File {
		return false
	}

	return true
}

// splitLines 按行分割字符串
func splitLines(s string) []string {
	var lines []string
	var currentLine string
	for _, r := range s {
		if r == '\n' {
			lines = append(lines, currentLine)
			currentLine = ""
		} else {
			currentLine += string(r)
		}
	}
	if currentLine != "" {
		lines = append(lines, currentLine)
	}
	return lines
}

// GetLogStats 获取日志统计信息
type LogStats struct {
	TotalLogs   int64            `json:"total_logs"`   // 总日志数
	TypeStats   map[string]int64 `json:"type_stats"`   // 按类型统计
	IPStats     map[string]int64 `json:"ip_stats"`     // 按IP统计
	ActionStats map[string]int64 `json:"action_stats"` // 按操作统计
	DailyStats  map[string]int64 `json:"daily_stats"`  // 按天统计
}

// GetLogStats 获取日志统计信息
func GetLogStats() (*LogStats, error) {
	// 获取日志目录
	logPath := "./logs"
	if globalLogger != nil {
		logPath = globalLogger.logPath
	}

	// 读取所有日志文件
	files, err := os.ReadDir(logPath)
	if err != nil {
		return nil, fmt.Errorf("读取日志目录失败: %v", err)
	}

	stats := &LogStats{
		TotalLogs:   0,
		TypeStats:   make(map[string]int64),
		IPStats:     make(map[string]int64),
		ActionStats: make(map[string]int64),
		DailyStats:  make(map[string]int64),
	}

	// 遍历所有日志文件
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		filename := file.Name()
		if len(filename) < 13 || filename[:5] != "logs_" || filename[13:] != ".log" {
			continue
		}

		filePath := filepath.Join(logPath, filename)
		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}

		// 按行解析日志
		lines := splitLines(string(fileContent))
		for _, line := range lines {
			if line == "" {
				continue
			}

			var entry LogEntry
			if err := json.Unmarshal([]byte(line), &entry); err != nil {
				continue
			}

			// 更新统计信息
			stats.TotalLogs++
			stats.TypeStats[entry.Type]++
			stats.IPStats[entry.IP]++
			stats.ActionStats[entry.Action]++
			stats.DailyStats[entry.Timestamp.Format("2006-01-02")]++
		}
	}

	return stats, nil
}

// ClearOldLogs 清理旧日志
func ClearOldLogs(days int) (int64, error) {
	if days <= 0 {
		days = 7 // 默认清理7天前的日志
	}

	// 获取日志目录
	logPath := "./logs"
	if globalLogger != nil {
		logPath = globalLogger.logPath
	}

	// 计算清理日期
	cleanDate := time.Now().AddDate(0, 0, -days)

	// 读取所有日志文件
	files, err := os.ReadDir(logPath)
	if err != nil {
		return 0, fmt.Errorf("读取日志目录失败: %v", err)
	}

	var deletedCount int64
	// 遍历所有日志文件
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		filename := file.Name()
		if len(filename) < 13 || filename[:5] != "logs_" || filename[13:] != ".log" {
			continue
		}

		// 解析文件名中的日期
		fileDate, err := time.Parse("20060102", filename[5:13])
		if err != nil {
			continue
		}

		// 检查是否需要清理
		if fileDate.Before(cleanDate) {
			filePath := filepath.Join(logPath, filename)
			if err := os.Remove(filePath); err != nil {
				continue
			}
			deletedCount++
		}
	}

	return deletedCount, nil
}
