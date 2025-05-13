package main

import (
	"SuperAgent/global"
	"SuperAgent/initialize"
	"context"
	"fmt"

	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

func graph() {
	initialize.InitConfig("..")
	initialize.InitChatModel()
	ctx := context.Background()
	g := compose.NewGraph[map[string]any, *schema.Message]()
	_ = g.AddChatTemplateNode("prompt", prompt.FromMessages(schema.FString, schema.SystemMessage("你的档案是{role},你是一个活泼的精灵"), schema.UserMessage("你叫什么名字？")))
	_ = g.AddChatModelNode("model", global.ChatModel)
	_ = g.AddEdge(compose.START, "prompt")
	_ = g.AddEdge("prompt", "model")
	_ = g.AddEdge("model", compose.END)
	runnable, _ := g.Compile(ctx)
	out, _ := runnable.Invoke(ctx, map[string]any{"role": "初音未来"})
	fmt.Println(out)
}
