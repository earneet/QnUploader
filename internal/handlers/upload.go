package handlers

import (
	"io"
	"net/http"

	"qiniu-uploader/internal/config"
	"qiniu-uploader/internal/models"
	"qiniu-uploader/internal/services"
	"qiniu-uploader/internal/utils"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	config       *config.Config
	qiniuService *services.QiniuService
}

func NewUploadHandler(cfg *config.Config, qiniuService *services.QiniuService) *UploadHandler {
	return &UploadHandler{
		config:       cfg,
		qiniuService: qiniuService,
	}
}

// UploadImage 处理图片上传
func (h *UploadHandler) UploadImage(c *gin.Context) {
	// 获取上传的文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.UploadResponse{
			Success: false,
			Message: "请选择要上传的文件",
		})
		return
	}
	defer file.Close()

	// 验证文件类型
	if err := utils.ValidateFilename(header.Filename); err != nil {
		c.JSON(http.StatusBadRequest, models.UploadResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	// 验证文件大小和类型
	if err := utils.ValidateFile(c.Request, h.config.MaxFileSize, h.config.AllowedTypes); err != nil {
		c.JSON(http.StatusBadRequest, models.UploadResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	// 读取文件内容
	fileData, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.UploadResponse{
			Success: false,
			Message: "读取文件失败",
		})
		return
	}

	// 上传到七牛云
	response, err := h.qiniuService.UploadFile(fileData, header.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.UploadResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetImages 获取图片列表
func (h *UploadHandler) GetImages(c *gin.Context) {
	prefix := c.Query("prefix")
	limit := 50 // 默认限制50张图片

	images, err := h.qiniuService.GetFileList(prefix, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ImageListResponse{
			Success: false,
			Data:    []models.ImageInfo{},
			Total:   0,
		})
		return
	}

	c.JSON(http.StatusOK, models.ImageListResponse{
		Success: true,
		Data:    images,
		Total:   len(images),
	})
}