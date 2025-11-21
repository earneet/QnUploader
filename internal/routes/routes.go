package routes

import (
	"qiniu-uploader/internal/config"
	"qiniu-uploader/internal/handlers"
	"qiniu-uploader/internal/middleware"
	"qiniu-uploader/internal/services"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置路由
func SetupRoutes(router *gin.Engine, cfg *config.Config) {
	// 初始化服务
	qiniuService := services.NewQiniuService(cfg)
	uploadHandler := handlers.NewUploadHandler(cfg, qiniuService)

	// 全局中间件
	router.Use(middleware.CORS())

	// 静态文件服务
	router.Static("/static", "./web/static")

	// API路由组
	api := router.Group("/api")
	{
		// 上传相关路由
		upload := api.Group("/upload")
		{
			upload.POST("", uploadHandler.UploadImage)
		}

		// 图片相关路由
		images := api.Group("/images")
		{
			images.GET("", uploadHandler.GetImages)
		}
	}

	// 默认路由
	router.GET("/", func(c *gin.Context) {
		c.File("./web/static/index.html")
	})
}