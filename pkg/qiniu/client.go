package qiniu

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

// Client 七牛云客户端
type Client struct {
	bucketManager *storage.BucketManager
	formUploader  *storage.FormUploader
	config        *Config
}

// Config 七牛云配置
type Config struct {
	AccessKey string
	SecretKey string
	Bucket    string
	Domain    string
}

// UploadResult 上传结果
type UploadResult struct {
	Success  bool
	Message  string
	FileURL  string
	FileSize int64
	Key      string
	Hash     string
}

// NewClient 创建新的七牛云客户端
func NewClient(cfg *Config) *Client {
	mac := qbox.NewMac(cfg.AccessKey, cfg.SecretKey)

	qiniuConfig := storage.Config{
		Zone:          nil, // 自动检测区域
		UseHTTPS:      true,
		UseCdnDomains: true,
	}

	return &Client{
		bucketManager: storage.NewBucketManager(mac, &qiniuConfig),
		formUploader:  storage.NewFormUploader(&qiniuConfig),
		config:        cfg,
	}
}

// UploadFile 上传文件到七牛云
func (c *Client) UploadFile(filePath string) (*UploadResult, error) {
	// 检查文件是否存在
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return &UploadResult{
			Success: false,
			Message: fmt.Sprintf("文件不存在: %v", err),
		}, err
	}

	// 验证文件大小（最大10MB）
	if fileInfo.Size() > 10*1024*1024 {
		return &UploadResult{
			Success: false,
			Message: "文件大小超过10MB限制",
		}, fmt.Errorf("文件大小超过限制")
	}

	// 验证文件类型
	if !c.isImageFile(filePath) {
		return &UploadResult{
			Success: false,
			Message: "不支持的文件类型，仅支持图片文件",
		}, fmt.Errorf("不支持的文件类型")
	}

	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return &UploadResult{
			Success: false,
			Message: fmt.Sprintf("打开文件失败: %v", err),
		}, err
	}
	defer file.Close()

	// 生成存储key
	key := c.generateFileKey(filepath.Base(filePath))

	// 获取上传凭证
	putPolicy := storage.PutPolicy{
		Scope: c.config.Bucket,
	}
	upToken := putPolicy.UploadToken(c.mac())

	// 上传文件
	ret := storage.PutRet{}
	err = c.formUploader.Put(context.Background(), &ret, upToken, key, file, fileInfo.Size(), nil)
	if err != nil {
		return &UploadResult{
			Success: false,
			Message: fmt.Sprintf("上传失败: %v", err),
		}, err
	}

	// 生成访问URL
	fileURL := c.generateFileURL(ret.Key)

	return &UploadResult{
		Success:  true,
		Message:  "上传成功",
		FileURL:  fileURL,
		FileSize: fileInfo.Size(),
		Key:      ret.Key,
		Hash:     ret.Hash,
	}, nil
}

// ListFiles 获取文件列表
func (c *Client) ListFiles(prefix string, limit int) ([]FileInfo, error) {
	entries, _, _, hasNext, err := c.bucketManager.ListFiles(
		c.config.Bucket,
		prefix,
		"",
		"",
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("获取文件列表失败: %v", err)
	}

	var files []FileInfo
	for _, entry := range entries {
		if c.isImageFile(entry.Key) {
			file := FileInfo{
				Key:      entry.Key,
				URL:      c.generateFileURL(entry.Key),
				FileSize: entry.Fsize,
				MimeType: entry.MimeType,
				Uploaded: time.Unix(entry.PutTime/10000000, 0),
			}
			files = append(files, file)
		}
	}

	// 如果有更多文件，可以继续获取（这里简化处理）
	_ = hasNext

	return files, nil
}

// FileInfo 文件信息
type FileInfo struct {
	Key      string
	URL      string
	FileSize int64
	MimeType string
	Uploaded time.Time
}

// generateFileKey 生成文件存储key
func (c *Client) generateFileKey(filename string) string {
	ext := filepath.Ext(filename)
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("images/%d%s", timestamp, ext)
}

// generateFileURL 生成文件访问URL
func (c *Client) generateFileURL(key string) string {
	if c.config.Domain != "" {
		return fmt.Sprintf("https://%s/%s", c.config.Domain, key)
	}

	// 如果没有配置域名，使用七牛云默认域名格式
	// 注意：实际使用时应该配置正确的域名
	return fmt.Sprintf("https://example.com/%s", key)
}

// isImageFile 检查是否为图片文件
func (c *Client) isImageFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	imageExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp", ".bmp"}

	for _, imageExt := range imageExts {
		if ext == imageExt {
			return true
		}
	}
	return false
}

// mac 获取七牛云认证对象
func (c *Client) mac() *qbox.Mac {
	return qbox.NewMac(c.config.AccessKey, c.config.SecretKey)
}