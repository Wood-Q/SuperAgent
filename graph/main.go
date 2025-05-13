package main

import (
	"context"
	"fmt"

	"SuperAgent/global"
	"SuperAgent/initialize"

	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

func main() {
	initialize.InitConfig("..")
	initialize.InitChatModel()
	ctx := context.Background()
	chain := compose.NewChain[map[string]any, *schema.Message]().
		AppendChatTemplate(prompt.FromMessages(schema.FString, schema.SystemMessage("你的档案是{role},你要记住你自己的档案"), schema.UserMessage("你叫什么名字？"))).
		AppendChatModel(global.ChatModel)
	runnable, err := chain.Compile(ctx)
	if err != nil {
		fmt.Println("compile error: ", err)
	}
	result, err := runnable.Invoke(ctx, map[string]any{"role": "初音未来"})
	if err != nil {
		fmt.Println("invoke error: ", err)
	}
	fmt.Println(result)
}
