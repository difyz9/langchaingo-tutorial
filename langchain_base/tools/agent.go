package tools

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// 智能意图识别引擎
type IntentEngine struct {
	tools map[string]Tool
}

func NewIntentEngine() *IntentEngine {
	tools := map[string]Tool{
		"VideoDownloader": VideoDownloaderTool{},
		"VideoTranslator": VideoTranslatorTool{},
		"COSUpTool":       COSUpTool{},
		"ReportGenerator": ReportGeneratorTool{},
		"EmailSender":     EmailSenderTool{},
		"VideoEditer":  VideoEditerTool{},
	}

	return &IntentEngine{tools: tools}
}

// 工作流步骤定义
type WorkflowStep struct {
	ToolName    string
	Input       string
	Description string
}

// 分析用户意图并生成工作流步骤
// 这里可以根据用户意图和工具的能力生成一系列工作流步骤
// 例如，如果用户意图是下载视频并翻译字幕，可以生成两个步骤：
// 1. 使用 VideoDownloaderTool 下载视频
// 2. 使用 VideoTranslatorTool 翻译视频字幕
// 具体实现可以根据实际需求进行调整
func (ie *IntentEngine) AnalyzeIntentAndPlan(userIntent string) []WorkflowStep {
	fmt.Printf("🧠 分析用户意图: %s\n", userIntent)
	input := strings.ToLower(userIntent)

	var steps []WorkflowStep

	// 提取视频URL
	videoURLRegex := regexp.MustCompile(`https?://[^\s]+`)
	videoURLs := videoURLRegex.FindAllString(userIntent, -1)

	// 提取邮箱地址
	emailRegex := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	emails := emailRegex.FindAllString(userIntent, -1)

	// 意图关键词检测
	needsDownload := strings.Contains(input, "下载") || strings.Contains(input, "download")
	needsTranslation := strings.Contains(input, "翻译") || strings.Contains(input, "translate") || strings.Contains(input, "字幕")
	needsUpload := strings.Contains(input, "上传") || strings.Contains(input, "upload") || strings.Contains(input, "cos") || strings.Contains(input, "存储")
	needsReport := strings.Contains(input, "报告") || strings.Contains(input, "report") || strings.Contains(input, "总结")
	needsEmail := strings.Contains(input, "邮件") || strings.Contains(input, "email") || strings.Contains(input, "发送") || strings.Contains(input, "通知")
	needEditer:= strings.Contains(input, "视频剪辑") || strings.Contains(input, "edit") || strings.Contains(input, "剪辑") 

	if len(videoURLs) > 0 {
		videoURL := videoURLs[0]

		// 第一步：下载视频
		if needsDownload || needsTranslation || needsUpload || needsReport || needsEmail {
			steps = append(steps, WorkflowStep{
				ToolName:    "VideoDownloader",
				Input:       videoURL,
				Description: fmt.Sprintf("下载视频: %s", videoURL),
			})
		}
	  
		if needEditer {

			steps = append(steps, WorkflowStep{
				ToolName:    "VideoEditer",
				Input:       "downloaded_video.mp4;00:00:00;00:01:00",
				Description: "编辑视频片段",
			})
		}



		// 第二步：翻译字幕（如果需要）
		if needsTranslation {
			steps = append(steps, WorkflowStep{
				ToolName:    "VideoTranslator",
				Input:       "downloaded_video.mp4",
				Description: "翻译视频字幕",
			})
		}

		// 第三步：上传到COS（如果需要）
		if needsUpload {
			steps = append(steps, WorkflowStep{
				ToolName:    "COSUpTool",
				Input:       "downloaded_video.mp4",
				Description: "上传视频到腾讯云COS",
			})
		}

		// 第四步：生成报告（如果需要）
		if needsReport {
			steps = append(steps, WorkflowStep{
				ToolName:    "ReportGenerator",
				Input:       "downloaded_video.mp4",
				Description: "生成视频处理报告",
			})
		}


		// 第五步：发送邮件通知（如果需要且有邮箱）
		if needsEmail && len(emails) > 0 {
			email := emails[0]
			subject := "视频处理完成通知"
			body := "您的视频处理任务已完成，详情请查看附件报告。"

			if needsDownload && !needsTranslation && !needsUpload && !needsReport {
				subject = "视频下载完成通知"
				body = "视频下载任务已完成。"
			}

			emailInput := fmt.Sprintf("%s;%s;%s", email, subject, body)
			steps = append(steps, WorkflowStep{
				ToolName:    "EmailSender",
				Input:       emailInput,
				Description: fmt.Sprintf("发送邮件通知到: %s", email),
			})
		}
	}

	return steps
}

// 执行工作流
func (ie *IntentEngine) ExecuteWorkflow(ctx context.Context, steps []WorkflowStep) error {
	if len(steps) == 0 {
		return fmt.Errorf("没有找到可执行的工作流步骤")
	}

	fmt.Printf("📋 工作流计划 (%d 个步骤):\n", len(steps))
	for i, step := range steps {
		fmt.Printf("  %d. %s\n", i+1, step.Description)
	}
	fmt.Println()

	// 按顺序执行每个步骤
	for i, step := range steps {
		fmt.Printf("--- 执行步骤 %d/%d ---\n", i+1, len(steps))

		tool, exists := ie.tools[step.ToolName]
		if !exists {
			return fmt.Errorf("工具 %s 不存在", step.ToolName)
		}

		result, err := tool.Call(ctx, step.Input)
		if err != nil {
			return fmt.Errorf("步骤 %d 执行失败: %w", i+1, err)
		}

		fmt.Printf("✅ %s\n\n", result)

		// 在步骤之间添加短暂延迟，模拟真实处理
		if i < len(steps)-1 {
			time.Sleep(500 * time.Millisecond)
		}
	}

	return nil
}

// 处理用户意图的主要入口
func (ie *IntentEngine) ProcessUserIntent(ctx context.Context, userIntent string) error {
	return ie.processUserIntentWithConfirmation(ctx, userIntent, true)
}

// 处理用户意图（交互模式，无需确认）
func (ie *IntentEngine) ProcessUserIntentInteractive(ctx context.Context, userIntent string) error {
	return ie.processUserIntentWithConfirmation(ctx, userIntent, false)
}

// 内部处理用户意图的方法
func (ie *IntentEngine) processUserIntentWithConfirmation(ctx context.Context, userIntent string, needConfirmation bool) error {
	fmt.Printf("🎯 收到用户请求: %s\n\n", userIntent)

	// 分析意图并制定执行计划
	steps := ie.AnalyzeIntentAndPlan(userIntent)

	if len(steps) == 0 {
		fmt.Println("❌ 抱歉，我无法理解您的请求。请提供更明确的指令，比如：")
		fmt.Println("   - 包含视频URL")
		fmt.Println("   - 明确指出需要的操作（下载、翻译、上传等）")
		fmt.Println("   - 提供接收通知的邮箱地址")
		fmt.Println("")
		fmt.Println("💡 建议格式:")
		fmt.Println("   • 下载 https://example.com/video.mp4 并发送到 user@email.com")
		fmt.Println("   • 翻译 https://site.com/video.mp4 字幕，上传到COS，生成报告")
		return fmt.Errorf("无法识别用户意图")
	}

	if needConfirmation {
		// 询问用户确认是否执行
		fmt.Printf("📋 计划执行以下 %d 个步骤:\n", len(steps))
		for i, step := range steps {
			fmt.Printf("  %d. %s\n", i+1, step.Description)
		}
		
		fmt.Print("❓ 是否继续执行工作流？(y/n): ")
		var confirmation string
		fmt.Scanln(&confirmation)
		
		if strings.ToLower(strings.TrimSpace(confirmation)) != "y" && 
		   strings.ToLower(strings.TrimSpace(confirmation)) != "yes" {
			fmt.Println("⏹️  用户取消执行")
			return nil
		}
	}

	// 执行工作流
	fmt.Println("🚀 开始执行工作流...")
	if err := ie.ExecuteWorkflow(ctx, steps); err != nil {
		return fmt.Errorf("工作流执行失败: %w", err)
	}

	fmt.Println("🎉 工作流执行完成！")
	return nil
}

// 展示系统能力
func (ie *IntentEngine) ShowCapabilities() {
	fmt.Println("🛠️  系统工具能力:")
	for _, tool := range ie.tools {
		fmt.Printf("  • %s: %s\n", tool.Name(), tool.Description())
	}

}
 