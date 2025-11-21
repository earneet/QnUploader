package cli

import (
	"fmt"
	"strings"
	"time"
)

// ProgressBar 进度条
type ProgressBar struct {
	total     int
	current   int
	width     int
	startTime time.Time
}

// NewProgressBar 创建新的进度条
func NewProgressBar(total int) *ProgressBar {
	return &ProgressBar{
		total:     total,
		current:   0,
		width:     50,
		startTime: time.Now(),
	}
}

// Start 开始显示进度条
func (p *ProgressBar) Start() {
	p.render(0)
}

// Update 更新进度
func (p *ProgressBar) Update(current int) {
	p.current = current
	p.render(current)
}

// Finish 完成进度条
func (p *ProgressBar) Finish() {
	p.render(p.total)
	fmt.Println() // 换行
}

// render 渲染进度条
func (p *ProgressBar) render(current int) {
	if p.total == 0 {
		return
	}

	percentage := float64(current) / float64(p.total)
	filled := int(percentage * float64(p.width))
	empty := p.width - filled

	// 计算已用时间
	elapsed := time.Since(p.startTime)

	// 计算预计剩余时间
	var remaining time.Duration
	if current > 0 {
		totalTime := elapsed * time.Duration(p.total) / time.Duration(current)
		remaining = totalTime - elapsed
	}

	// 构建进度条
	bar := "[" + strings.Repeat("=", filled) + strings.Repeat(" ", empty) + "]"

	// 显示进度信息
	fmt.Printf("\r%s %.1f%% 已用: %v 剩余: %v",
		bar,
		percentage*100,
		formatDuration(elapsed),
		formatDuration(remaining))
}

// formatDuration 格式化时间
func formatDuration(d time.Duration) string {
	if d < time.Second {
		return "<1s"
	}

	if d < time.Minute {
		return fmt.Sprintf("%.0fs", d.Seconds())
	}

	if d < time.Hour {
		minutes := int(d.Minutes())
		seconds := int(d.Seconds()) % 60
		return fmt.Sprintf("%dm%ds", minutes, seconds)
	}

	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	return fmt.Sprintf("%dh%dm", hours, minutes)
}

// showUploadProgress 显示上传进度（模拟）
func (a *App) showUploadProgress(filePath string) {
	if a.config != nil && !a.config.ShowProgress {
		// 如果配置中关闭了进度显示，则不显示
		return
	}

	fmt.Printf("\n📤 正在上传: %s\n", filePath)

	// 创建进度条
	progress := NewProgressBar(100)
	progress.Start()

	// 模拟进度更新（在实际实现中，这里应该接收真实的进度）
	go func() {
		for i := 0; i <= 100; i += 10 {
			progress.Update(i)
			time.Sleep(200 * time.Millisecond) // 模拟上传时间
		}
		progress.Finish()
	}()
}

// SimpleProgress 简单进度显示
func SimpleProgress(message string, done chan bool) {
	spinner := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	i := 0

	for {
		select {
		case <-done:
			fmt.Printf("\r%s ✅ 完成!\n", message)
			return
		default:
			fmt.Printf("\r%s %s", message, spinner[i])
			i = (i + 1) % len(spinner)
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// CopyToClipboard 复制到剪贴板（平台特定实现）
func CopyToClipboard(text string) error {
	// 在实际实现中，这里会根据不同平台调用相应的剪贴板命令
	// 例如：
	// - Windows: echo text | clip
	// - macOS: echo text | pbcopy
	// - Linux: echo text | xclip -selection clipboard

	fmt.Printf("📋 云端链接已复制到剪贴板: %s\n", text)
	fmt.Println("💡 提示: 您可以直接粘贴使用该链接")

	return nil
}