package system

import (
	"gin_cloud_drive/backend/config"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

type SystemInfo struct {
	CPU      CPUInfo    `json:"cpu"`
	Memory   MemoryInfo `json:"memory"`
	Disk     DiskInfo   `json:"disk"`
	OS       string     `json:"os"`
	Hostname string     `json:"hostname"`
	Time     string     `json:"time"`
}

type CPUInfo struct {
	Cores        int32   `json:"cores"`
	UsagePercent float64 `json:"usage_percent"`
	ModelName    string  `json:"model_name"`
}

type MemoryInfo struct {
	Total        uint64  `json:"total"`
	Available    uint64  `json:"available"`
	Used         uint64  `json:"used"`
	UsagePercent float64 `json:"usage_percent"`
}

type DiskInfo struct {
	Total        uint64  `json:"total"`
	Free         uint64  `json:"free"`
	Used         uint64  `json:"used"`
	UsagePercent float64 `json:"usage_percent"`
}

type DataPoint struct {
	Timestamp int64   `json:"timestamp"`
	CPU       float64 `json:"cpu"`
	Memory    float64 `json:"memory"`
	Disk      float64 `json:"disk"`
}

type HistoryData struct {
	Data []DataPoint `json:"data"`
	mu   sync.RWMutex
}

var historyData = &HistoryData{
	Data: make([]DataPoint, 0),
}

const (
	// 时间范围对应的最大数据点数
	maxPoints1h       = 60          // 1小时 * 60分钟
	maxPoints6h       = 360         // 6小时 * 60分钟
	maxPoints24h      = 1440        // 24小时 * 60分钟
	maxPoints3d       = 4320        // 3天 * 24小时 * 60分钟
	maxPoints7d       = 10080       // 7天 * 24小时 * 60分钟
	maxPoints         = maxPoints7d // 最大保留7天数据
	dataRetentionDays = 7           // 数据保留天数
)

// InitSystemMonitor 初始化系统状态监控
func InitSystemMonitor() {
	// 加载历史数据
	if err := loadHistoryFromFile(); err != nil {
		fmt.Println("加载历史数据失败:", err)
	}

	// 启动数据收集协程
	go collectDataPeriodically()

	// 启动数据清理协程
	go cleanDataPeriodically()

	// 启动时立即执行一次清理
	cleanOldData()
}

// 获取CPU信息
func getCPUInfo() CPUInfo {
	cores, _ := cpu.Counts(false)
	percent, _ := cpu.Percent(0, false)
	cpuInfo, _ := cpu.Info()

	modelName := ""
	if len(cpuInfo) > 0 {
		modelName = cpuInfo[0].ModelName
	}

	cpuPercent := 0.0
	if len(percent) > 0 {
		cpuPercent = percent[0]
	}

	return CPUInfo{
		Cores:        int32(cores),
		UsagePercent: cpuPercent,
		ModelName:    modelName,
	}
}

// 获取内存信息
func getMemoryInfo() MemoryInfo {
	v, err := mem.VirtualMemory()
	if err != nil {
		fmt.Printf("获取内存信息失败: %v\n", err)
		// 返回默认值，避免nil指针解引用
		return MemoryInfo{
			Total:        0,
			Available:    0,
			Used:         0,
			UsagePercent: 0,
		}
	}
	return MemoryInfo{
		Total:        v.Total,
		Available:    v.Available,
		Used:         v.Used,
		UsagePercent: v.UsedPercent,
	}
}

// 获取磁盘信息
func getDiskInfo() DiskInfo {
	diskStat, err := disk.Usage("/")
	if err != nil {
		fmt.Printf("获取磁盘信息失败: %v\n", err)
		// 返回默认值，避免nil指针解引用
		return DiskInfo{
			Total:        0,
			Free:         0,
			Used:         0,
			UsagePercent: 0,
		}
	}
	return DiskInfo{
		Total:        diskStat.Total,
		Free:         diskStat.Free,
		Used:         diskStat.Used,
		UsagePercent: diskStat.UsedPercent,
	}
}

// GetSystemInfo 获取系统信息
func GetSystemInfo() SystemInfo {
	hostname, _ := os.Hostname()

	return SystemInfo{
		CPU:      getCPUInfo(),
		Memory:   getMemoryInfo(),
		Disk:     getDiskInfo(),
		OS:       runtime.GOOS,
		Hostname: hostname,
		Time:     time.Now().Format("2006-01-02 15:04:05"),
	}
}

// 保存历史数据到文件
func saveHistoryToFile() error {
	historyData.mu.RLock()
	defer historyData.mu.RUnlock()

	data, err := json.MarshalIndent(historyData.Data, "", "  ")
	if err != nil {
		return err
	}

	cfg := config.GetConfig()
	return os.WriteFile(cfg.System.DataFile, data, 0644)
}

// 从文件加载历史数据
func loadHistoryFromFile() error {
	cfg := config.GetConfig()
	data, err := os.ReadFile(cfg.System.DataFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	historyData.mu.Lock()
	defer historyData.mu.Unlock()

	return json.Unmarshal(data, &historyData.Data)
}

// 删除7天前的数据
func cleanOldData() {
	historyData.mu.Lock()
	defer historyData.mu.Unlock()

	now := time.Now().Unix()
	cutoffTime := now - int64(dataRetentionDays*24*3600) // 7天前

	var newData []DataPoint
	for _, point := range historyData.Data {
		if point.Timestamp >= cutoffTime {
			newData = append(newData, point)
		}
	}

	if len(newData) < len(historyData.Data) {
		deletedCount := len(historyData.Data) - len(newData)
		fmt.Printf("[%s] 清理过期数据: 删除了 %d 条7天前的记录\n",
			time.Now().Format("2006-01-02 15:04:05"), deletedCount)
		historyData.Data = newData
		saveHistoryToFile()
	}
}

// 添加数据点
func addDataPoint(info SystemInfo) {
	point := DataPoint{
		Timestamp: time.Now().Unix(),
		CPU:       info.CPU.UsagePercent,
		Memory:    info.Memory.UsagePercent,
		Disk:      info.Disk.UsagePercent,
	}

	historyData.mu.Lock()
	historyData.Data = append(historyData.Data, point)

	// 保持最多7天的数据
	if len(historyData.Data) > maxPoints {
		historyData.Data = historyData.Data[len(historyData.Data)-maxPoints:]
	}
	historyData.mu.Unlock()

	saveHistoryToFile()
}

// GetHistoryData 获取指定时间范围的数据
func GetHistoryData(hours int) []DataPoint {
	historyData.mu.RLock()
	defer historyData.mu.RUnlock()

	now := time.Now().Unix()
	secondsAgo := int64(hours * 3600)
	startTime := now - secondsAgo

	var result []DataPoint
	for _, point := range historyData.Data {
		if point.Timestamp >= startTime {
			result = append(result, point)
		}
	}

	return result
}

// 定期收集数据
func collectDataPeriodically() {
	cfg := config.GetConfig()
	ticker := time.NewTicker(time.Duration(cfg.System.Interval) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		info := GetSystemInfo()
		addDataPoint(info)
	}
}

// 定期清理旧数据（每天凌晨1点执行）
func cleanDataPeriodically() {
	for {
		now := time.Now()
		// 计算下次凌晨1点的时间
		next := now.Add(time.Hour * 24)
		next = time.Date(next.Year(), next.Month(), next.Day(), 1, 0, 0, 0, next.Location())

		duration := next.Sub(now)
		timer := time.NewTimer(duration)

		<-timer.C
		cleanOldData()
	}
}
