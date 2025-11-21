package services

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"qiniu-uploader/internal/config"
	"qiniu-uploader/internal/models"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

type QiniuService struct {
	config     *config.Config
	mac        *qbox.Mac
	bucket     string
	uploader   *storage.FormUploader
	bucketMgr  *storage.BucketManager
}

func NewQiniuService(cfg *config.Config) *QiniuService {
	mac := qbox.NewMac(cfg.QiniuAccessKey, cfg.QiniuSecretKey)

	// 配置上传区域（这里使用华东区域作为默认值）
	cfg := storage.Config{
		Zone:          &storage.ZoneHuadong,
		UseHTTPS:      true,
		UseCdnDomains: true,
	}

	uploader := storage.NewFormUploader(&cfg)
	bucketMgr := storage.NewBucketManager(mac, &cfg)

	return &QiniuService{
		config:    cfg,
		mac:       mac,
		bucket:    cfg.QiniuBucket,
		uploader:  uploader,
		bucketMgr: bucketMgr,
	}
}

// UploadFile 上传文件到七牛云
func (s *QiniuService) UploadFile(fileData []byte, filename string) (*models.UploadResponse, error) {
	// 生成存储key
	key := s.generateFileKey(filename)

	// 获取上传凭证
	putPolicy := storage.PutPolicy{
		Scope: s.bucket,
	}
	upToken := putPolicy.UploadToken(s.mac)

	// 上传文件
	ret := storage.PutRet{}
	err := s.uploader.Put(context.Background(), &ret, upToken, key, fileData, int64(len(fileData)), nil)
	if err != nil {
		return nil, fmt.Errorf("上传失败: %v", err)
	}

	// 构建返回结果
	response := &models.UploadResponse{
		Success: true,
		Message: "上传成功",
	}
	response.Data.Key = ret.Key
	response.Data.Hash = ret.Hash
	response.Data.URL = s.generateFileURL(ret.Key)
	response.Data.FileSize = int64(len(fileData))
	response.Data.MimeType = "image/jpeg" // 这里应该根据实际文件类型设置

	return response, nil
}

// GetFileList 获取文件列表
func (s *QiniuService) GetFileList(prefix string, limit int) ([]models.ImageInfo, error) {
	entries, _, _, hasNext, err := s.bucketMgr.ListFiles(
		s.bucket,
		prefix,
		"",
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("获取文件列表失败: %v", err)
	}

	var images []models.ImageInfo
	for _, entry := range entries {
		// 只处理图片文件
		if s.isImageFile(entry.Key) {
			image := models.ImageInfo{
				ID:       entry.Hash,
				Key:      entry.Key,
				URL:      s.generateFileURL(entry.Key),
				FileSize: entry.Fsize,
				MimeType: entry.MimeType,
				Uploaded: time.Unix(entry.PutTime/10000000, 0).Format(time.RFC3339),
			}
			images = append(images, image)
		}
	}

	// 如果有更多文件，可以继续获取（这里简化处理）
	_ = hasNext

	return images, nil
}

// generateFileKey 生成文件存储key
func (s *QiniuService) generateFileKey(filename string) string {
	ext := filepath.Ext(filename)
	timestamp := time.Now().Unix()
	return fmt.Sprintf("images/%d%s", timestamp, ext)
}

// generateFileURL 生成文件访问URL
func (s *QiniuService) generateFileURL(key string) string {
	if s.config.QiniuDomain != "" {
		return fmt.Sprintf("https://%s/%s", s.config.QiniuDomain, key)
	}

	// 如果没有配置域名，使用七牛云默认域名
	domain := storage.DefaultPubHosts[0]
	return fmt.Sprintf("https://%s/%s", domain, key)
}

// isImageFile 检查是否为图片文件
func (s *QiniuService) isImageFile(filename string) bool {
	ext := filepath.Ext(filename)
	imageExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp", ".bmp"}

	for _, imageExt := range imageExts {
		if strings.EqualFold(ext, imageExt) {
			return true
		}
	}
	return false
}