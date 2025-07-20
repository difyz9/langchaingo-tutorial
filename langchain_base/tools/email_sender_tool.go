package tools

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// EmailSenderTool
type EmailSenderTool struct{}

func (t EmailSenderTool) Name() string {
	return "EmailSender"
}

func (t EmailSenderTool) Description() string {
	return "发送电子邮件。输入格式应为 'recipient;subject;body'，其中 recipient 是收件人邮箱地址，subject 是邮件主题，body 是邮件内容。"
}

func (t EmailSenderTool) Call(ctx context.Context, input string) (string, error) {
	parts := strings.SplitN(input, ";", 3)
	if len(parts) != 3 {
		return "", fmt.Errorf("邮件格式错误，正确格式: 'recipient;subject;body'")
	}

	recipient := parts[0]
	subject := parts[1]
	body := parts[2]

	fmt.Printf("📧 正在发送邮件到 '%s'，主题: '%s'，内容: %s...\n", recipient, subject, body)
	time.Sleep(1 * time.Second) // 模拟邮件发送
	return fmt.Sprintf("邮件发送成功！收件人: %s", recipient), nil
}
