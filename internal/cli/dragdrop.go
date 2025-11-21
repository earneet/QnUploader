package cli

import (
	"fmt"
	"runtime"
	"strings"
)

// DragDropHandler 拖拽处理器
type DragDropHandler struct {
	app *App
}

// NewDragDropHandler 创建新的拖拽处理器
func NewDragDropHandler(app *App) *DragDropHandler {
	return &DragDropHandler{
		app: app,
	}
}

// Start 启动拖拽监听
func (h *DragDropHandler) Start() error {
	// 在命令行工具中，我们无法直接监听系统级的拖拽事件
	// 但是我们可以提供友好的提示信息
	fmt.Println("💡 提示: 您可以直接拖拽文件到终端窗口，然后按回车上传")
	fmt.Println("   或者输入文件路径进行上传")

	// 对于不同平台，提供特定的拖拽提示
	switch runtime.GOOS {
	case "windows":
		fmt.Println("   Windows: 拖拽文件到终端窗口，路径会自动填充")
	case "darwin":
		fmt.Println("   macOS: 拖拽文件到终端窗口，路径会自动填充")
	case "linux":
		fmt.Println("   Linux: 拖拽文件到终端窗口，路径会自动填充")
	}

	// 在实际实现中，这里会启动平台特定的拖拽监听
	// 但由于命令行工具的限制，我们依赖终端本身的拖拽支持
	return nil
}

// HandleFileDrop 处理拖拽文件
func (h *DragDropHandler) HandleFileDrop(filePath string) error {
	// 清理文件路径（去除可能的引号和空格）
	filePath = strings.TrimSpace(filePath)
	filePath = strings.Trim(filePath, "\"")

	fmt.Printf("\n📁 检测到文件拖拽: %s\n", filePath)

	// 调用上传逻辑
	return h.app.uploadFile(filePath)
}

// getPlatformSpecificImplementation 获取平台特定的实现
func (h *DragDropHandler) getPlatformSpecificImplementation() string {
	switch runtime.GOOS {
	case "windows":
		return `
// Windows实现使用COM接口监听拖拽事件
// 需要导入: "github.com/lxn/walk"
// 或者使用Windows API
func (h *DragDropHandler) startWindows() error {
    // Windows特定的拖拽实现
    return nil
}`

	case "darwin":
		return `
// macOS实现使用Cocoa API
// 需要导入: "github.com/progrium/macdriver"
func (h *DragDropHandler) startMacOS() error {
    // macOS特定的拖拽实现
    return nil
}`

	case "linux":
		return `
// Linux实现使用GTK或X11
// 需要导入: "github.com/gotk3/gotk3"
func (h *DragDropHandler) startLinux() error {
    // Linux特定的拖拽实现
    return nil
}`

	default:
		return "// 不支持的操作系统"
	}
}

// IsDragDropSupported 检查当前平台是否支持拖拽
func (h *DragDropHandler) IsDragDropSupported() bool {
	// 在命令行环境中，我们依赖终端本身的拖拽支持
	// 大多数现代终端都支持拖拽文件到窗口
	return true
}

// GetDragDropInstructions 获取拖拽使用说明
func (h *DragDropHandler) GetDragDropInstructions() string {
	instructions := []string{
		"💡 拖拽上传使用说明:",
		"   1. 打开文件管理器",
		"   2. 选择要上传的文件",
		"   3. 拖拽文件到终端窗口",
		"   4. 文件路径会自动填充",
		"   5. 按回车键开始上传",
		"",
		"📝 注意:",
		"   - 支持拖拽多个文件",
		"   - 支持图片文件 (jpg, png, gif, webp, bmp)",
		"   - 文件大小限制: 10MB",
	}

	return strings.Join(instructions, "\n")
}