package main

import (
	"bufio"
	"context"
	"fmt"
	"langchain_base/tools"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	engine := tools.NewIntentEngine()

	// 显示欢迎信息和系统能力
	showWelcomeMessage()
	engine.ShowCapabilities()
	
	// 提供运行模式选择
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📋 请选择运行模式:")
	fmt.Println("  1. 交互模式 - 实时输入命令并执行")
	fmt.Println("  2. 测试模式 - 运行预设测试场景")
	fmt.Print("请选择 (1/2): ")
	
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	choice := strings.TrimSpace(scanner.Text())
	
	switch choice {
	case "1":
		runInteractiveMode(engine)
	case "2":
		runTestMode(engine)
	default:
		fmt.Println("默认进入交互模式...")
		runInteractiveMode(engine)
	}
}

// 显示欢迎信息
func showWelcomeMessage() {
	fmt.Println("🎬 智能视频处理工作流引擎")
	fmt.Println("=" + strings.Repeat("=", 48) + "=")
	fmt.Println("🚀 支持视频下载、翻译、剪辑、上传、报告生成、邮件通知等功能")
	fmt.Println("💡 示例命令:")
	fmt.Println("   • 下载视频 https://example.com/video.mp4 并发送到 user@email.com")
	fmt.Println("   • 翻译 https://site.com/lecture.mp4 的字幕，上传到COS，生成报告")
	fmt.Println("   • 剪辑视频 https://video.com/long.mp4 并通知 admin@company.com")
}

// 交互模式
func runInteractiveMode(engine *tools.IntentEngine) {
	fmt.Println("\n🔄 进入交互模式 (输入 'quit' 或 'exit' 退出，输入 'help' 查看帮助)")
	fmt.Println(strings.Repeat("-", 60))
	
	// 设置信号处理，优雅退出
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	
	scanner := bufio.NewScanner(os.Stdin)
	
	for {
		fmt.Print("\n💬 请输入您的需求: ")
		
		// 监听用户输入和系统信号
		inputChan := make(chan string, 1)
		
		go func() {
			if scanner.Scan() {
				inputChan <- scanner.Text()
			}
		}()
		
		select {
		case input := <-inputChan:
			input = strings.TrimSpace(input)
			
			if input == "" {
				continue
			}
			
			// 处理特殊命令
			switch strings.ToLower(input) {
			case "quit", "exit", "q":
				fmt.Println("👋 再见！感谢使用智能视频处理工作流引擎")
				return
			case "help", "h":
				showHelpMessage()
				continue
			case "clear", "cls":
				clearScreen()
				continue
			case "capabilities", "cap":
				engine.ShowCapabilities()
				continue
			}
			
			// 执行用户意图
			fmt.Println("\n" + strings.Repeat("-", 60))
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			
			if err := engine.ProcessUserIntentInteractive(ctx, input); err != nil {
				fmt.Printf("⚠️  执行失败: %v\n", err)
			}
			
			cancel()
			fmt.Println(strings.Repeat("-", 60))
			
		case <-c:
			fmt.Println("\n\n👋 收到退出信号，正在优雅退出...")
			return
		}
	}
}

// 测试模式
func runTestMode(engine *tools.IntentEngine) {
	fmt.Println("\n🧪 进入测试模式")
	fmt.Println(strings.Repeat("-", 60))
	
	testScenarios := []struct {
		name   string
		intent string
	}{
		{
			name:   "完整视频处理流水线",
			intent: "请帮我下载视频 https://example.com/tutorial.mp4，翻译字幕，上传到COS存储，生成处理报告，并发送到 admin@company.com",
		},
		{
			name:   "下载并发送邮件的处理流水线",
			intent: "请帮我下载视频 https://example.com/tutorial.mp4,并发送到 admin@company.com",
		},
		{
			name:   "简单下载和邮件通知",
			intent: "下载这个视频 https://videos.site.com/meeting.mp4 完成后通知 user@example.com",
		},
		{
			name:   "视频翻译专项流程",
			intent: "我需要翻译 https://content.edu.com/lecture.mp4 的字幕，完成后发邮件给 translator@team.com",
		},
		{
			name:   "云存储备份流程,剪辑视频",
			intent: "下载 https://backup.com/archive.mp4,剪辑视频, 并上传到COS，然后生成报告发给 storage@admin.com",
		},
		{
			name:   "错误处理测试",
			intent: "帮我处理一下视频",
		},
	}
	
	// 逐个执行测试场景
	for i, scenario := range testScenarios {
		fmt.Printf("🎯 测试场景 %d: %s\n", i+1, scenario.name)
		fmt.Println(strings.Repeat("-", 60))

		if err := engine.ProcessUserIntent(context.Background(), scenario.intent); err != nil {
			fmt.Printf("⚠️  场景执行结果: %v\n", err)
		}

		if i < len(testScenarios)-1 {
			fmt.Println("\n" + strings.Repeat("=", 80) + "\n")
			time.Sleep(1 * time.Second) // 场景间暂停
		}
	}
}

// 显示帮助信息
func showHelpMessage() {
	fmt.Println("\n📚 帮助信息:")
	fmt.Println("  • 输入自然语言描述您的需求")
	fmt.Println("  • 支持的操作: 下载、翻译、剪辑、上传、报告生成、邮件通知")
	fmt.Println("  • 'help' 或 'h' - 显示此帮助")
	fmt.Println("  • 'capabilities' 或 'cap' - 显示系统工具能力")
	fmt.Println("  • 'clear' 或 'cls' - 清屏")
	fmt.Println("  • 'quit', 'exit' 或 'q' - 退出程序")
	fmt.Println("\n💡 命令示例:")
	fmt.Println("  • 下载视频 https://example.com/video.mp4")
	fmt.Println("  • 翻译 https://site.com/video.mp4 字幕并发邮件给 user@email.com")
	fmt.Println("  • 剪辑视频 https://video.com/content.mp4 上传到COS生成报告")
}

// 清屏函数
func clearScreen() {
	fmt.Print("\033[H\033[2J")
}
