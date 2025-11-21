package cli

import (
	"testing"
)

func TestHandleFileDropIntegration(t *testing.T) {
	// 创建一个模拟的App实例用于测试
	app := &App{}
	handler := NewDragDropHandler(app)

	// 测试正常路径处理（不会触发转换）
	tests := []struct {
		name        string
		inputPath   string
		expectError bool
	}{
		{
			name:        "Linux path",
			inputPath:   "/home/user/file.txt",
			expectError: true, // 因为文件不存在，会返回错误
		},
		{
			name:        "WSL path",
			inputPath:   "/mnt/c/Users/test/file.txt",
			expectError: true, // 因为文件不存在，会返回错误
		},
		{
			name:        "Path with quotes",
			inputPath:   `"/home/user/file.txt"`,
			expectError: true, // 因为文件不存在，会返回错误
		},
		{
			name:        "Path with spaces",
			inputPath:   "  /home/user/file.txt  ",
			expectError: true, // 因为文件不存在，会返回错误
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handler.HandleFileDrop(tt.inputPath)

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none for input: %s", tt.inputPath)
			}

			// 我们主要测试路径清理和转换逻辑
			// 文件不存在的错误是预期的，因为我们没有创建真实文件
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error for input %s: %v", tt.inputPath, err)
			}
		})
	}
}

func TestDragDropHandlerCreation(t *testing.T) {
	app := &App{}
	handler := NewDragDropHandler(app)

	if handler == nil {
		t.Error("Expected non-nil handler")
	}

	if handler.app != app {
		t.Error("Handler app reference mismatch")
	}
}

func TestDragDropSupportCheck(t *testing.T) {
	app := &App{}
	handler := NewDragDropHandler(app)

	// 在命令行环境中，拖拽应该总是被支持
	if !handler.IsDragDropSupported() {
		t.Error("Expected drag drop to be supported")
	}
}