package utils

import (
	"gin_cloud_drive/backend/config"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// FileInfo 文件信息
type FileInfo struct {
	Name         string    `json:"name"`
	Path         string    `json:"path"`
	Size         int64     `json:"size"`
	IsDirectory  bool      `json:"is_directory"`
	ModifiedTime time.Time `json:"modified_time"`
	Type         string    `json:"type"`
}

// ListFiles 列出文件
func ListFiles(path string, sortBy string, sortOrder string) ([]FileInfo, error) {
	uploadPath := config.GetConfig().File.UploadPath
	// 使用filepath.Join构建完整路径，确保路径分隔符正确
	fullPath := filepath.Join(uploadPath, path)

	// 检查路径是否存在
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return nil, err
	}

	// 读取目录
	entries, err := os.ReadDir(fullPath)
	if err != nil {
		return nil, err
	}

	var files []FileInfo
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		// 使用filepath.Join构建文件路径，确保路径分隔符正确
		filePath := filepath.Join(path, info.Name())
		// 将路径分隔符转换为正斜杠，确保URL兼容
		filePath = filepath.ToSlash(filePath)

		file := FileInfo{
			Name:         info.Name(),
			Path:         filePath,
			Size:         info.Size(),
			IsDirectory:  entry.IsDir(),
			ModifiedTime: info.ModTime(),
			Type:         getFileType(info.Name()),
		}

		files = append(files, file)
	}

	// 排序文件列表
	// 首先确保文件夹在文件前面
	// 然后根据指定字段排序
	ascending := sortOrder == "asc"

	// 定义比较函数
	compare := func(a, b FileInfo) bool {
		// 文件夹始终在文件前面
		if a.IsDirectory != b.IsDirectory {
			return a.IsDirectory
		}

		// 根据排序字段比较
		switch sortBy {
		case "name":
			if a.Name != b.Name {
				if ascending {
					return a.Name < b.Name
				} else {
					return a.Name > b.Name
				}
			}
		case "size":
			if a.Size != b.Size {
				if ascending {
					return a.Size < b.Size
				} else {
					return a.Size > b.Size
				}
			}
		case "time":
			if !a.ModifiedTime.Equal(b.ModifiedTime) {
				if ascending {
					return a.ModifiedTime.Before(b.ModifiedTime)
				} else {
					return a.ModifiedTime.After(b.ModifiedTime)
				}
			}
		case "type":
			if a.Type != b.Type {
				if ascending {
					return a.Type < b.Type
				} else {
					return a.Type > b.Type
				}
			}
		}

		// 最后按名称排序作为默认
		if ascending {
			return a.Name < b.Name
		} else {
			return a.Name > b.Name
		}
	}

	// 冒泡排序
	for i := 0; i < len(files)-1; i++ {
		for j := 0; j < len(files)-i-1; j++ {
			if compare(files[j+1], files[j]) {
				files[j], files[j+1] = files[j+1], files[j]
			}
		}
	}

	return files, nil
}

// getFileType 获取文件类型
func getFileType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	if ext == "" {
		return "file"
	}

	// 图片类型
	imageExts := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".svg"}
	for _, e := range imageExts {
		if ext == e {
			return "image"
		}
	}

	// 文档类型
	docExts := []string{".txt", ".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".md"}
	for _, e := range docExts {
		if ext == e {
			return "document"
		}
	}

	// 视频类型
	videoExts := []string{".mp4", ".avi", ".mov", ".wmv", ".flv"}
	for _, e := range videoExts {
		if ext == e {
			return "video"
		}
	}

	// 音频类型
	audioExts := []string{".mp3", ".wav", ".ogg", ".flac"}
	for _, e := range audioExts {
		if ext == e {
			return "audio"
		}
	}

	return "file"
}

// CreateDirectory 创建目录
func CreateDirectory(path string) error {
	uploadPath := config.GetConfig().File.UploadPath
	fullPath := filepath.Join(uploadPath, path)
	return os.MkdirAll(fullPath, 0755)
}

// RenameFile 重命名文件
func RenameFile(oldPath, newName string) error {
	uploadPath := config.GetConfig().File.UploadPath
	oldFullPath := filepath.Join(uploadPath, oldPath)
	newFullPath := filepath.Join(filepath.Dir(oldFullPath), newName)
	return os.Rename(oldFullPath, newFullPath)
}

// MoveFile 移动文件或目录
func MoveFile(oldPath, newPath string) error {
	uploadPath := config.GetConfig().File.UploadPath

	// 清理路径，移除可能的控制字符
	oldPath = strings.TrimSpace(oldPath)
	newPath = strings.TrimSpace(newPath)

	// 检查路径是否为空
	if oldPath == "" || newPath == "" {
		return fmt.Errorf("path cannot be empty")
	}

	// 构建完整路径，确保路径分隔符正确
	oldFullPath := filepath.Join(uploadPath, oldPath)
	newFullPath := filepath.Join(uploadPath, newPath)

	// 检查源文件/目录是否存在
	_, err := os.Stat(oldFullPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("source path does not exist: %s", oldFullPath)
	}

	// 检查目标路径是否存在
	newInfo, err := os.Stat(newFullPath)

	// 确定目标是否为目录
	isDir := false
	if err == nil {
		isDir = newInfo.IsDir()
	} else {
		// 目标路径不存在，默认视为目录
		isDir = true
	}

	// 如果目标是目录，将源移动到该目录下
	if isDir {
		// 获取源文件名
		sourceName := filepath.Base(oldFullPath)
		newFullPath = filepath.Join(newFullPath, sourceName)
		// 确保目标目录存在
		if err := os.MkdirAll(filepath.Dir(newFullPath), 0755); err != nil {
			return fmt.Errorf("failed to create target directory: %v", err)
		}
	} else if os.IsNotExist(err) {
		// 目标不存在，确保父目录存在
		if err := os.MkdirAll(filepath.Dir(newFullPath), 0755); err != nil {
			return fmt.Errorf("failed to create parent directory: %v", err)
		}
	}

	// 确保旧路径和新路径不相同
	if oldFullPath == newFullPath {
		return nil // 路径相同，无需移动
	}

	// 执行移动操作
	if err := os.Rename(oldFullPath, newFullPath); err != nil {
		return fmt.Errorf("failed to rename: %v, old: %s, new: %s", err, oldFullPath, newFullPath)
	}

	return nil
}

// DeleteFile 删除文件或文件夹
func DeleteFile(path string) error {
	uploadPath := config.GetConfig().File.UploadPath

	// 清理路径，移除可能的控制字符
	path = strings.TrimSpace(path)

	// 检查路径是否为空
	if path == "" {
		return fmt.Errorf("path cannot be empty")
	}

	// 构建完整路径
	fullPath := filepath.Join(uploadPath, path)

	// 检查路径是否存在
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return nil // 路径不存在，视为删除成功
	}

	// 执行删除操作
	return os.RemoveAll(fullPath)
}
