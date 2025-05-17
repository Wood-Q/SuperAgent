package main

import (
	"MoonAgent/internal/config"
	"MoonAgent/internal/pipeline"
	"context"

	"github.com/cloudwego/eino/schema"
)

func main() {
	ctx := context.Background()

	config.InitConfig("../configs")

	model, err := pipeline.NewChatModel(ctx)
	if err != nil {
		panic(err)
	}

	// 准备消息
	messages := []*schema.Message{
		schema.SystemMessage("你是一个助手"),
		schema.UserMessage("介绍一下火山引擎"),
	}

	// 生成回复
	response, err := model.Generate(ctx, messages)
	if err != nil {
		panic(err)
	}

	// 处理回复
	println(response.Content)

	// 获取 Token 使用情况
	if usage := response.ResponseMeta.Usage; usage != nil {
		println("提示 Tokens:", usage.PromptTokens)
		println("生成 Tokens:", usage.CompletionTokens)
		println("总 Tokens:", usage.TotalTokens)
	}
}
