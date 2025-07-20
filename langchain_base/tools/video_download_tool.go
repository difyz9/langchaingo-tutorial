package tools

import (
	"context"
	"fmt"
	"time"
)

// VideoDownloaderTool
type VideoDownloaderTool struct{}

func (t VideoDownloaderTool) Name() string {
	return "VideoDownloader"
}

func (t VideoDownloaderTool) Description() string {
	return "从给定的URL下载视频文件。输入应为视频URL。"
}

func (t VideoDownloaderTool) Call(ctx context.Context, input string) (string, error) {
	fmt.Printf("🎬 正在下载视频: %s...\n", input)
	time.Sleep(2 * time.Second) // 模拟网络延迟
	return fmt.Sprintf("视频 '%s' 下载成功，保存为 downloaded_video.mp4", input), nil
}
