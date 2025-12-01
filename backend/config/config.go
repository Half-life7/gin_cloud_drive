package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Server ServerConfig `json:"server"`
	User   UserConfig   `json:"user"`
	File   FileConfig   `json:"file"`
	System SystemConfig `json:"system"`
}

type ServerConfig struct {
	Port string `json:"port"`
}

type UserConfig struct {
	AdminUsername string `json:"admin_username"`
	AdminPassword string `json:"admin_password"`
}

type FileConfig struct {
	UploadPath string `json:"upload_path"`
	MaxSize    int64  `json:"max_size"`
}

type SystemConfig struct {
	DataFile string `json:"data_file"`
	Interval int    `json:"interval"`
}

var config *Config

// InitConfig 初始化配置
func InitConfig() {
	// 默认配置
	config = &Config{
		Server: ServerConfig{
			Port: "8080",
		},
		User: UserConfig{
			AdminUsername: "admin",
			AdminPassword: "admin123",
		},
		File: FileConfig{
			UploadPath: "./upload",
			MaxSize:    100 << 20, // 100MB
		},
		System: SystemConfig{
			DataFile: "./system/system_history.json",
			Interval: 60, // 1分钟
		},
	}

	// 从文件加载配置（如果存在）
	loadConfigFromFile()
}

// GetConfig 获取配置
func GetConfig() *Config {
	return config
}

// 从文件加载配置
func loadConfigFromFile() {
	file, err := os.ReadFile("./backend/config/config.json")
	if err != nil {
		// 如果文件不存在，使用默认配置
		return
	}

	err = json.Unmarshal(file, config)
	if err != nil {
		// 如果解析失败，使用默认配置
		return
	}
}

// 保存配置到文件
func SaveConfig() error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile("./backend/config/config.json", data, 0644)
}
