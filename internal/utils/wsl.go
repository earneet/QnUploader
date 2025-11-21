package utils

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

// IsWSLEnvironment 检测是否运行在WSL环境中
func IsWSLEnvironment() bool {
	// 方法1: 检查 /proc/version 文件内容
	if runtime.GOOS == "linux" {
		if version, err := os.ReadFile("/proc/version"); err == nil {
			versionStr := string(version)
			return strings.Contains(strings.ToLower(versionStr), "microsoft") ||
				strings.Contains(strings.ToLower(versionStr), "wsl")
		}

		// 方法2: 检查 WSL_DISTRO_NAME 环境变量
		if os.Getenv("WSL_DISTRO_NAME") != "" {
			return true
		}

		// 方法3: 检查 WSL_INTEROP 环境变量
		if os.Getenv("WSL_INTEROP") != "" {
			return true
		}
	}

	return false
}

// ConvertWindowsPathToWSL 将Windows路径转换为WSL路径，带验证
func ConvertWindowsPathToWSL(windowsPath string) (string, error) {
	// 清理路径（去除引号和空格）
	cleanedPath := strings.TrimSpace(windowsPath)
	cleanedPath = strings.Trim(cleanedPath, "\"")

	// 检查是否为Windows路径格式
	if !isWindowsPath(cleanedPath) {
		return cleanedPath, nil // 不是Windows路径，直接返回
	}

	// 提取驱动器字母和路径部分
	driveLetter := strings.ToLower(string(cleanedPath[0]))
	pathPart := cleanedPath[3:] // 跳过 "C:\\"

	// 转换路径分隔符
	pathPart = strings.ReplaceAll(pathPart, "\\", "/")

	// 构建WSL路径
	wslPath := "/mnt/" + driveLetter + "/" + pathPart

	// 验证转换后的路径是否存在
	if _, err := os.Stat(wslPath); os.IsNotExist(err) {
		return "", fmt.Errorf("转换后的路径不存在: %s", wslPath)
	}

	return wslPath, nil
}

// isWindowsPath 检查是否为Windows路径格式
func isWindowsPath(path string) bool {
	// Windows路径格式: C:\\path\\to\\file 或 C:/path/to/file
	if len(path) < 3 {
		return false
	}

	// 检查驱动器字母
	if !((path[0] >= 'A' && path[0] <= 'Z') || (path[0] >= 'a' && path[0] <= 'z')) {
		return false
	}

	// 检查冒号和路径分隔符
	return path[1] == ':' && (path[2] == '\\' || path[2] == '/')
}

// NormalizePathForWSL 在WSL环境中规范化路径，带错误处理
func NormalizePathForWSL(path string) (string, error) {
	if !IsWSLEnvironment() {
		return path, nil // 不在WSL环境中，直接返回
	}

	// 如果是Windows路径，转换为WSL路径
	if isWindowsPath(path) {
		return ConvertWindowsPathToWSL(path)
	}

	// 已经是WSL或Linux路径，直接返回
	return path, nil
}