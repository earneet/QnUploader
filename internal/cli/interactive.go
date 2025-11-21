package cli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"qiniu-uploader/internal/config"
	"qiniu-uploader/pkg/qiniu"
)

// startInteractiveUpload 启动交互式上传
func (a *App) startInteractiveUpload() error {
	if a.client == nil {
		fmt.Println("❌ 七牛云客户端未初始化")
		fmt.Println("请先运行 'qiniu-uploader config init' 配置七牛云信息")
		return fmt.Errorf("七牛云客户端未初始化")
	}

	fmt.Println("🚀 七牛云上传工具 - 交互模式")
	fmt.Println("=" + strings.Repeat("=", 50))
	fmt.Println("支持以下操作:")
	fmt.Println("  1. 输入文件路径上传 (支持拖拽文件到终端)")
	fmt.Println("  2. 输入 'list' 查看已上传文件")
	fmt.Println("  3. 输入 'config' 显示当前配置")
	fmt.Println("  4. 输入 'quit' 或 'exit' 退出")
	fmt.Println("=" + strings.Repeat("=", 50))

	// 显示拖拽使用说明
	if a.dragDropHandler != nil {
		fmt.Println()
		fmt.Println(a.dragDropHandler.GetDragDropInstructions())
	}

	// 启动拖拽监听（如果支持）
	if a.dragDropHandler != nil && a.dragDropHandler.IsDragDropSupported() {
		if err := a.dragDropHandler.Start(); err != nil {
			fmt.Printf("⚠️  拖拽功能初始化失败: %v\n", err)
		}
	}

	// 启动输入循环
	return a.startInputLoop()
}

// startInputLoop 启动输入循环
func (a *App) startInputLoop() error {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("\n📁 请输入文件路径或命令: ")

		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())

		// 处理命令
		switch strings.ToLower(input) {
		case "", "quit", "exit":
			fmt.Println("👋 再见!")
			return nil
		case "list":
			a.listUploadedFiles()
		case "config":
			a.showConfig()
		default:
			// 处理文件上传
			if err := a.handleFileInput(input); err != nil {
				fmt.Printf("❌ 错误: %v\n", err)
			}
		}
	}

	return scanner.Err()
}

// handleFileInput 处理文件输入
func (a *App) handleFileInput(input string) error {
	// 如果拖拽处理器可用，使用它来处理文件路径（包括WSL路径转换）
	if a.dragDropHandler != nil {
		return a.dragDropHandler.HandleFileDrop(input)
	}

	// 回退到原始逻辑
	filePath := input

	// 如果输入包含引号，去除引号
	if strings.HasPrefix(filePath, "\"") && strings.HasSuffix(filePath, "\"") {
		filePath = filePath[1 : len(filePath)-1]
	}

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("文件不存在: %s", filePath)
	}

	// 上传文件
	return a.uploadFile(filePath)
}

// listUploadedFiles 列出已上传文件
func (a *App) listUploadedFiles() {
	if a.client == nil {
		fmt.Println("❌ 七牛云客户端未初始化")
		return
	}

	fmt.Println("\n📚 已上传文件列表:")
	fmt.Println("-" + strings.Repeat("-", 80))

	files, err := a.client.ListFiles("images/", 20)
	if err != nil {
		fmt.Printf("❌ 获取文件列表失败: %v\n", err)
		return
	}

	if len(files) == 0 {
		fmt.Println("  暂无上传文件")
		return
	}

	for i, file := range files {
		fmt.Printf("%2d. %s\n", i+1, filepath.Base(file.Key))
		fmt.Printf("    大小: %.2f MB | 上传时间: %s\n",
			float64(file.FileSize)/1024/1024,
			file.Uploaded.Format("2006-01-02 15:04:05"))
		fmt.Printf("    链接: %s\n", file.URL)
		if i < len(files)-1 {
			fmt.Println()
		}
	}
	fmt.Println("-" + strings.Repeat("-", 80))
}

// initConfig 初始化配置
func (a *App) initConfig() error {
	fmt.Println("🔧 初始化七牛云上传工具配置")
	fmt.Println("=" + strings.Repeat("=", 50))

	cfg := &config.Config{}

	// 获取七牛云配置
	fmt.Println("\n📋 请输入七牛云配置:")
	fmt.Print("Access Key: ")
	fmt.Scanln(&cfg.QiniuAccessKey)

	fmt.Print("Secret Key: ")
	fmt.Scanln(&cfg.QiniuSecretKey)

	fmt.Print("Bucket 名称: ")
	fmt.Scanln(&cfg.QiniuBucket)

	fmt.Print("域名 (可选): ")
	fmt.Scanln(&cfg.QiniuDomain)

	// 设置默认快捷键配置
	cfg.HotkeyKeys = []int{85} // U键
	cfg.HotkeyCtrl = true
	cfg.HotkeyShift = true
	cfg.HotkeyAlt = false
	cfg.AutoCopyURL = true
	cfg.ShowProgress = true

	// 保存配置
	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("保存配置失败: %v", err)
	}

	fmt.Println("\n✅ 配置保存成功!")
	fmt.Println("配置文件位置: ~/.config/qiniu-uploader/config.yaml")

	// 重新加载配置
	a.config = cfg

	// 重新初始化七牛云客户端
	qiniuConfig := &qiniu.Config{
		AccessKey: cfg.QiniuAccessKey,
		SecretKey: cfg.QiniuSecretKey,
		Bucket:    cfg.QiniuBucket,
		Domain:    cfg.QiniuDomain,
	}
	a.client = qiniu.NewClient(qiniuConfig)

	return nil
}

// showConfig 显示当前配置
func (a *App) showConfig() error {
	if a.config == nil {
		fmt.Println("❌ 配置未初始化")
		fmt.Println("请运行 'qiniu-uploader config init' 初始化配置")
		return nil
	}

	fmt.Println("\n🔧 当前配置:")
	fmt.Println("=" + strings.Repeat("=", 50))

	// 七牛云配置
	fmt.Println("📋 七牛云配置:")
	if a.config.QiniuAccessKey != "" {
		fmt.Printf("  Access Key: %s*** (已设置)\n", a.config.QiniuAccessKey[:4])
	} else {
		fmt.Println("  Access Key: 未设置")
	}

	if a.config.QiniuSecretKey != "" {
		fmt.Printf("  Secret Key: %s*** (已设置)\n", a.config.QiniuSecretKey[:4])
	} else {
		fmt.Println("  Secret Key: 未设置")
	}

	fmt.Printf("  Bucket: %s\n", a.config.QiniuBucket)
	fmt.Printf("  域名: %s\n", a.config.QiniuDomain)

	// 快捷键配置
	fmt.Println("\n⌨️  快捷键配置:")
	modifiers := []string{}
	if a.config.HotkeyCtrl {
		modifiers = append(modifiers, "Ctrl")
	}
	if a.config.HotkeyShift {
		modifiers = append(modifiers, "Shift")
	}
	if a.config.HotkeyAlt {
		modifiers = append(modifiers, "Alt")
	}

	if len(modifiers) > 0 {
		fmt.Printf("  快捷键: %s+U\n", strings.Join(modifiers, "+"))
	} else {
		fmt.Println("  快捷键: 未设置")
	}

	// UI配置
	fmt.Println("\n🎨 UI配置:")
	fmt.Printf("  自动复制链接: %v\n", a.config.AutoCopyURL)
	fmt.Printf("  显示进度条: %v\n", a.config.ShowProgress)

	fmt.Println("=" + strings.Repeat("=", 50))

	return nil
}

// startService 启动后台服务
func (a *App) startService() error {
	fmt.Println("🔧 后台服务功能正在开发中...")
	fmt.Println("当前版本暂不支持后台服务模式")
	fmt.Println("请使用 'qiniu-uploader upload' 进入交互模式")
	return nil
}