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

func main_copy() {
	initialize.InitConfig("..")
	initialize.InitChatModel()
	ctx := context.Background()
	content, err := os.ReadFile("./muelsyse.txt")
	if err != nil {
		fmt.Println("read file error: ", err)
	}

	// 创建一个存储对话历史的切片
	var messages []schema.MessagesTemplate

	// 添加系统消息
	systemMessage := schema.SystemMessage("你的档案是{role},你是一个活泼的精灵")
	messages = append(messages, systemMessage)

	// 创建对话链
	chain := compose.NewChain[map[string]any, *schema.Message]().
		AppendChatTemplate(prompt.FromMessages(schema.FString, messages...)).
		AppendChatModel(global.ChatModel)

	runnable, err := chain.Compile(ctx)
	if err != nil {
		fmt.Println("compile error: ", err)
	}

	reader := bufio.NewReader(os.Stdin)

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

		// 添加用户消息到历史记录
		userMessage := schema.UserMessage(input)
		messages = append(messages, userMessage)

		// 更新对话链的模板，包含完整的对话历史
		chain = compose.NewChain[map[string]any, *schema.Message]().
			AppendChatTemplate(prompt.FromMessages(schema.FString, messages...)).
			AppendChatModel(global.ChatModel)

		runnable, err = chain.Compile(ctx)
		if err != nil {
			fmt.Println("compile error: ", err)
			continue
		}

		// 调用模型
		result, err := runnable.Invoke(ctx, map[string]any{"role": string(content)})
		if err != nil {
			fmt.Println("invoke error: ", err)
			continue
		}

		// 添加助手回复到历史记录
		messages = append(messages, result)

		// 打印助手回复
		fmt.Printf("\nAI: %s\n", result.Content)
	}
}
