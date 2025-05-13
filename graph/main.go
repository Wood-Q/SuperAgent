package main

import (
	"context"
	"fmt"
	"os"

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
	content, err := os.ReadFile("./muelsyse.txt")
	if err != nil {
		fmt.Println("read file error: ", err)
	}
	chain := compose.NewChain[map[string]any, *schema.Message]().
		AppendChatTemplate(prompt.FromMessages(schema.FString, schema.SystemMessage("你的档案是{role},你是一个活泼的精灵"), schema.UserMessage("你叫什么名字？"))).
		AppendChatModel(global.ChatModel)
	runnable, err := chain.Compile(ctx)
	if err != nil {
		fmt.Println("compile error: ", err)
	}
	result, err := runnable.Invoke(ctx, map[string]any{"role": string(content)})
	if err != nil {
		fmt.Println("invoke error: ", err)
	}
	fmt.Println(result)
}
