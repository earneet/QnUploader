package utils

import (
	"mime"
	"net/http"
	"path/filepath"
	"strings"
)

// ValidateFile 验证上传文件
func ValidateFile(fileHeader *http.Request, maxSize int64, allowedTypes []string) error {
	// 检查文件大小
	if fileHeader.ContentLength > maxSize {
		return &ValidationError{
			Code:    "FILE_TOO_LARGE",
			Message: "文件大小超过限制",
		}
	}

	// 检查文件类型
	contentType := fileHeader.Header.Get("Content-Type")
	if !isAllowedType(contentType, allowedTypes) {
		return &ValidationError{
			Code:    "INVALID_FILE_TYPE",
			Message: "不支持的文件类型",
		}
	}

	return nil
}

// ValidateFilename 验证文件名
func ValidateFilename(filename string) error {
	// 检查文件扩展名
	ext := strings.ToLower(filepath.Ext(filename))
	allowedExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}

	for _, allowedExt := range allowedExts {
		if ext == allowedExt {
			return nil
		}
	}

	return &ValidationError{
		Code:    "INVALID_FILE_EXTENSION",
		Message: "不支持的文件扩展名",
	}
}

// isAllowedType 检查文件类型是否在允许列表中
func isAllowedType(contentType string, allowedTypes []string) bool {
	for _, allowedType := range allowedTypes {
		if contentType == allowedType {
			return true
		}
		// 检查MIME类型前缀（如image/*）
		if strings.HasSuffix(allowedType, "/*") {
			typePrefix := strings.TrimSuffix(allowedType, "/*")
			if strings.HasPrefix(contentType, typePrefix+"/") {
				return true
			}
		}
	}
	return false
}

// GetMimeTypeFromExtension 根据文件扩展名获取MIME类型
func GetMimeTypeFromExtension(filename string) string {
	ext := filepath.Ext(filename)
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		return "application/octet-stream"
	}
	return mimeType
}

// ValidationError 验证错误
type ValidationError struct {
	Code    string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}