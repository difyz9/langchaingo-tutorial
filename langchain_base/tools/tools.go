package tools

import 	"context"


// 工具接口定义
type Tool interface {
	Name() string
	Description() string
	Call(ctx context.Context, input string) (string, error)
}
