package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

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
	type roleState struct {
		history []string
	}

	gen := func(ctx context.Context) *roleState {
		return &roleState{}
	}

	// 创建图结构
	g := compose.NewGraph[map[string]any, *schema.Message](compose.WithGenLocalState(gen))

	// 添加系统消息和用户消息模板
	messages := []schema.MessagesTemplate{
		schema.SystemMessage("你的档案是{role},你是一个活泼的精灵"),
		schema.UserMessage("{input}"), // 使用input作为用户输入的占位符
	}

	// 添加对话模板节点
	err = g.AddChatTemplateNode("prompt",
		prompt.FromMessages(schema.FString, messages...))
	if err != nil {
		fmt.Println("add template node error:", err)
		return
	}

	// 添加模型节点
	err = g.AddChatModelNode("model", global.ChatModel)
	if err != nil {
		fmt.Println("add model node error:", err)
		return
	}

	// 构建图的边
	_ = g.AddEdge(compose.START, "prompt")
	_ = g.AddEdge("prompt", "model")
	_ = g.AddEdge("model", compose.END)

	// 编译图
	runnable, err := g.Compile(ctx)
	if err != nil {
		fmt.Println("compile error:", err)
		return
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("开始对话（输入'exit'退出）:")

	for {
		// 获取用户输入
		fmt.Print("\n您: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("读取输入错误: %v", err)
			continue
		}

		// 处理输入
		input = strings.TrimSpace(input)
		if strings.ToLower(input) == "exit" {
			fmt.Println("再见！")
			break
		}

		// 构建参数映射
		params := map[string]any{
			"role":  string(content),
			"input": input,
		}

		// 调用模型
		result, err := runnable.Invoke(ctx, params)
		if err != nil {
			fmt.Printf("对话错误: %v\n", err)
			continue
		}

		// 打印助手回复
		fmt.Printf("\nAI: %s\n", result.Content)
	}
}
