package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"qiniu-uploader/internal/config"
	"qiniu-uploader/pkg/qiniu"
)

// App 命令行应用
type App struct {
	rootCmd *cobra.Command
	client  *qiniu.Client
	config  *config.Config
	dragDropHandler *DragDropHandler
}

// NewApp 创建新的命令行应用
func NewApp() *App {
	app := &App{}

	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("警告: 加载配置失败: %v\n", err)
		fmt.Println("请运行 'qu config init' 初始化配置")
	}
	app.config = cfg

	// 初始化七牛云客户端
	if cfg != nil && cfg.QiniuAccessKey != "" && cfg.QiniuSecretKey != "" && cfg.QiniuBucket != "" {
		qiniuConfig := &qiniu.Config{
			AccessKey: cfg.QiniuAccessKey,
			SecretKey: cfg.QiniuSecretKey,
			Bucket:    cfg.QiniuBucket,
			Domain:    cfg.QiniuDomain,
		}
		app.client = qiniu.NewClient(qiniuConfig)
	}

	app.setupCommands()

	// 初始化拖拽处理器
	app.dragDropHandler = NewDragDropHandler(app)

	return app
}

// setupCommands 设置命令
func (a *App) setupCommands() {
	a.rootCmd = &cobra.Command{
		Use:   "qu",
		Short: "七牛云文件上传工具",
		Long:  "支持拖拽上传、快捷键操作的文件上传工具",
		RunE: func(cmd *cobra.Command, args []string) error {
			// 如果没有参数，显示帮助信息
			if len(args) == 0 {
				return cmd.Help()
			}
			return nil
		},
	}

	// 添加上传命令
	a.rootCmd.AddCommand(a.newUploadCommand())

	// 添加服务命令
	a.rootCmd.AddCommand(a.newServiceCommand())

	// 添加配置命令
	a.rootCmd.AddCommand(a.newConfigCommand())

	// 添加版本命令
	a.rootCmd.AddCommand(a.newVersionCommand())
}

// Run 运行应用
func (a *App) Run() error {
	return a.rootCmd.Execute()
}

// newUploadCommand 创建上传命令
func (a *App) newUploadCommand() *cobra.Command {
	var filePath string

	cmd := &cobra.Command{
		Use:   "upload",
		Short: "上传文件到七牛云",
		Long:  "支持交互式上传、拖拽上传和指定文件路径上传",
		RunE: func(cmd *cobra.Command, args []string) error {
			if filePath != "" {
				// 指定文件路径上传
				return a.uploadFile(filePath)
			}

			if len(args) > 0 {
				// 使用参数中的文件路径
				return a.uploadFile(args[0])
			}

			// 交互式上传
			return a.startInteractiveUpload()
		},
	}

	cmd.Flags().StringVarP(&filePath, "file", "f", "", "指定要上传的文件路径")

	return cmd
}

// newServiceCommand 创建服务命令
func (a *App) newServiceCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "service",
		Short: "启动后台服务",
		Long:  "启动后台服务，支持全局快捷键和系统托盘",
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.startService()
		},
	}
}

// newConfigCommand 创建配置命令
func (a *App) newConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "配置管理",
		Long:  "管理七牛云上传工具的配置",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "init",
		Short: "初始化配置",
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.initConfig()
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "show",
		Short: "显示当前配置",
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.showConfig()
		},
	})

	return cmd
}

// newVersionCommand 创建版本命令
func (a *App) newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "显示版本信息",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("七牛云上传工具 v1.0.0")
		},
	}
}

// uploadFile 上传单个文件
func (a *App) uploadFile(filePath string) error {
	if a.client == nil {
		return fmt.Errorf("七牛云客户端未初始化，请先运行 'qu config init' 配置七牛云信息")
	}

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("文件不存在: %s", filePath)
	}

	fmt.Printf("正在上传: %s\n", filepath.Base(filePath))

	result, err := a.client.UploadFile(filePath)
	if err != nil {
		return fmt.Errorf("上传失败: %v", err)
	}

	if result.Success {
		fmt.Printf("✅ 上传成功!\n")
		fmt.Printf("📁 文件名: %s\n", filepath.Base(filePath))
		fmt.Printf("📊 文件大小: %.2f MB\n", float64(result.FileSize)/1024/1024)
		fmt.Printf("🔗 访问链接: %s\n", result.FileURL)
		fmt.Printf("🔑 存储Key: %s\n", result.Key)
	} else {
		fmt.Printf("❌ 上传失败: %s\n", result.Message)
	}

	return nil
}